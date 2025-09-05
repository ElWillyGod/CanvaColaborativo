# Letra del proyecto

Desarrollar un **servidor TCP** (telnet-friendly) en **Go** que expone un **canvas ASCII multiusuario** en tiempo real. Las personas se conectan con telnet host puerto, ejecutan **comandos de dibujo** (puntos, líneas, rectángulos, texto), ven las acciones de otros en vivo. El sistema debe manejar **concurrencia**, **resolución de conflictos** de edición, **rate-limiting** (anti-flood) y **persistencia** de snapshots para restaurar el estado. 


El foco no debería ser en la parte visual, sino enfocarse en las funcionalidades core: rate-limiting, persistencia, concurrencia, etc.


# Desarrollo de la letra

Todo en Go **tengo que orientar esto a miles de usuarios**

Canvas ASCII multiusuario en tiempo real

1. **Conectividad y protocolo:** servidor TCP line-oriented compatible con clientes telnet; ayuda mínima integrada. 
A que se refiere con ayuda mínima integrada?

2. **Funciones esenciales:** operaciones básicas de dibujo (puntos/figuras simples), mensajes de chat y limpieza del canvas con confirmación. Fusión entre Canvas?

3. **Concurrencia y coherencia:** difusión en tiempo real de cambios; política simple de resolución de conflictos y **rate limiting** por usuario. (como aplico memoria dinamica?)

4. **Persistencia:** snapshots para restaurar estado (podría ser por modificación como un Ctrl+Z)

5. **Configuración:** tamaño de canvas, puerto, y límites ajustables por variables de entorno. 

Variables de entorno de Telnet y del servidor.

Donde puedo aplicar estructuras? algo con skips list y memoria dinamica. Para la gestion de clientes?
slices para buffers de red y rate limiting?

el historial de comandos puede ser guardado mediante una lista enlazada.

[x] Puedo usar listas circulares para gestionar lo de rate limiting? 

manejar tiempos de espera, nose pueden quedar esperando si hay muchos clientes. noon

Grupos de canvas en paralelo:
las personas se pueden unir y crear su propio canvas con un id unico
/////////////////////////////////////////////////
Motor de canvas **tileado + esparso** (alto rendimiento)

* **Idea:** dividir el canvas en **tiles** (p. ej. 64×32). Solo asignás memoria para tiles “tocados”.
* **Estructuras:**

  * **HashMap → Tile** (map\[TileID]\*Tile) con **pooling** (`sync.Pool`) para reciclar tiles y buffers.
  * Dentro de cada tile, **matriz de run-length (RLE) por filas** para comprimir secuencias de caracteres iguales (ideal para ASCII).
  * **Copy-on-write** para snapshots (ver §4).
* **Beneficio:** baja uso de RAM, updates localizados, snapshots y “undos” baratos.



Compresión y optimización de tráfico

Implementar compresión delta: en lugar de mandar todo el canvas, enviar solo las diferencias (patches ASCII).

Podés mostrar cómo reducís ancho de banda para miles de usuarios concurrentes.

esto se puede hacer con telnet?

Opciones para “verlo” bien en Telnet

Modo textual (más simple):

Los usuarios ven mensajes tipo UPDATE 10 5 X.

No se actualiza el canvas en pantalla automáticamente, el usuario interpreta.

Esto es 100% compatible con Telnet.

Modo gráfico ASCII (más desafiante):

Podés usar códigos ANSI de terminal (telnet lo soporta) para mover el cursor a (10,5) y dibujar la X.

Entonces el usuario ve cómo el canvas cambia en vivo sin redibujar todo.

Esto sí da el efecto “canvas en tiempo real con deltas”.
///////////////////////////////////////////////////

Indexación Espacial con Quadtrees para Operaciones de Área
En lugar de (o además de) tu map[TileID]*Tile, puedes implementar un Quadtree para indexar los tiles que contienen datos.

Qué es: Un Quadtree es una estructura de datos en árbol usada para particionar un espacio 2D, subdividiendo recursivamente una región en cuatro cuadrantes.

////////////////////////////////////////////////////
Agrupación de Paquetes (Packet Batching) y Ticks del Servidor
Enviar cada pequeña actualización en su propio paquete TCP es extremadamente ineficiente debido a la sobrecarga de las cabeceras TCP/IP (40-60 bytes por paquete para enviar a veces un solo byte de datos).

Qué es: En lugar de que el goroutine de un cliente envíe datos inmediatamente después de una modificación, los coloca en un buffer de salida compartido o en un canal. Un único goroutine "broadcaster" se despierta a intervalos fijos (por ejemplo, cada 20-50 milisegundos, un "tick"), recoge todas las actualizaciones pendientes para cada cliente y las envía en un solo paquete grande.
Por qué sorprende: Es una implementación a nivel de aplicación del Algoritmo de Nagle, una optimización de red fundamental. Muestra una comprensión profunda de cómo funciona TCP y cómo evitar la congestión y el overhead.
Implementación:
Cada cliente tiene un canal de salida (chan []byte).
Cuando el canvas se modifica, se generan los "deltas" y se envían a los canales de todos los clientes suscritos.
El goroutine de escritura de cada cliente no envía inmediatamente. Intenta leer del canal en un bucle, agrupando todos los mensajes que pueda durante un breve período de tiempo (o hasta un tamaño máximo) antes de hacer una única llamada a conn.Write().
3
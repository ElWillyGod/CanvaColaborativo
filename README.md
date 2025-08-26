# Letra del proyecto

Desarrollar un **servidor TCP** (telnet-friendly) en **Go** que expone un **canvas ASCII multiusuario** en tiempo real. Las personas se conectan con telnet host puerto, ejecutan **comandos de dibujo** (puntos, líneas, rectángulos, texto), ven las acciones de otros en vivo. El sistema debe manejar **concurrencia**, **resolución de conflictos** de edición, **rate-limiting** (anti-flood) y **persistencia** de snapshots para restaurar el estado. 


El foco no debería ser en la parte visual, sino enfocarse en las funcionalidades core: rate-limiting, persistencia, concurrencia, etc.


# Desarrollo de la letra

Todo en Go

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

Puedo usar listas circulares para gestionar lo de rate limiting?

manejar tiempos de espera, nose pueden quedar esperando si hay muchos clientes.
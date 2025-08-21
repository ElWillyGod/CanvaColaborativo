# CanvaColaborativo
canvas ASCII multiusuario en tiempo real.

# Letra del proyecto

Desarrollar un **servidor TCP** (telnet-friendly) en **Go** que expone un **canvas ASCII multiusuario** en tiempo real. Las personas se conectan con telnet host puerto, ejecutan **comandos de dibujo** (puntos, líneas, rectángulos, texto), ven las acciones de otros en vivo. El sistema debe manejar **concurrencia**, **resolución de conflictos** de edición, **rate-limiting** (anti-flood) y **persistencia** de snapshots para restaurar el estado. 


El foco no debería ser en la parte visual, sino enfocarse en las funcionalidades core: rate-limiting, persistencia, concurrencia, etc.


# Desarrollo de la letra

Todo en Go

Canvas ASCII multiusuario en tiempo real

1. **Conectividad y protocolo:** servidor TCP line-oriented compatible con clientes telnet; ayuda mínima integrada. 
A que se refiere con ayuda mínima integrada?

2. **Funciones esenciales:** operaciones básicas de dibujo (puntos/figuras simples), mensajes de chat y limpieza del canvas con confirmación.

3. **Concurrencia y coherencia:** difusión en tiempo real de cambios; política simple de resolución de conflictos y **rate limiting** por usuario.

4. **Persistencia:** snapshots para restaurar estado (podría ser por modificación como un Ctrl+Z)

5. **Configuración:** tamaño de canvas, puerto, y límites ajustables por variables de entorno.



Seguramente use Function Pointer Para los comandos de dibujo.

Ya tengo lo distribuir los mensajes.

Falta lo de persistencia. y rate-limiting.
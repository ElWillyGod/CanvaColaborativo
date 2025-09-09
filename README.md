# 🎨 Canva Colaborativo en Go

Este proyecto es un servidor de canvas colaborativo en tiempo real implementado en Go. Permite que múltiples usuarios se conecten a través de Telnet y dibujen simultáneamente en un lienzo de texto compartido. El estado de cada lienzo se guarda de forma persistente en una base de datos **Valkey** (o Redis).

## ✨ Características Principales

-   **Colaboración en Tiempo Real**: Múltiples usuarios pueden dibujar en el mismo lienzo y ver las actualizaciones de los demás al instante.
-   **Persistencia Eficiente**: Los lienzos se guardan en Valkey utilizando una estrategia optimizada. En lugar de guardar todo el lienzo con cada cambio, solo se actualizan las "baldosas" (tiles) modificadas, reduciendo drásticamente la carga en la base de datos.
-   **Lienzos Múltiples**: Los usuarios pueden crear nuevos lienzos o unirse a lienzos existentes utilizando su ID único (UUID).
-   **Interfaz por Comandos**: La interacción se realiza a través de comandos de texto simples e intuitivos.
-   **Seguridad y Estabilidad**:
    -   **Rate Limiting**: Implementa un limitador de tasa por conexión para prevenir el abuso y los ataques de inundación (flooding).
    -   **Manejo Concurrente Seguro**: Utiliza mutex y canales de Go para gestionar de forma segura el estado compartido entre múltiples goroutines de clientes.
-   **Acciones Grupales**: Funciones como limpiar el lienzo requieren la confirmación de todos los usuarios conectados, promoviendo un entorno colaborativo.

## 🚀 Cómo Empezar

### Prerrequisitos

-   [Go](https://golang.org/dl/) (versión 1.18 o superior)
-   [Valkey](https://valkey.io/) o [Redis](https://redis.io/)
-   Un cliente Telnet.

### Instalación

1.  **Clona el repositorio:**
    ```sh
    git clone <URL_DEL_REPOSITORIO>
    cd CanvaColaborativo
    ```

2.  **Instala las dependencias:**
    ```sh
    go mod tidy
    ```

### Configuración

El servidor se puede configurar mediante variables de entorno:

-   `PORT`: El puerto en el que se ejecutará el servidor (por defecto: `8080`).
-   `VALKEY_ADDR`: La dirección del servidor Valkey/Redis (por defecto: `localhost:6379`).
-   `CANVAS_WIDTH`: Ancho del lienzo en caracteres (por defecto: `80`).
-   `CANVAS_HEIGHT`: Alto del lienzo en caracteres (por defecto: `40`).

### Ejecución

1.  Asegúrate de que tu servidor Valkey/Redis esté en funcionamiento.
2.  Inicia el servidor de canvas:
    ```sh
    go run .
    ```

## ✍️ Cómo Usar

1.  **Conéctate al servidor** usando un cliente Telnet:
    ```sh
    telnet localhost 8080
    ```

2.  **Únete a un lienzo**:
    -   Para crear un lienzo nuevo, simplemente conéctate. Se te asignará un nuevo ID de lienzo.
    -   Para unirte a un lienzo existente, usa el comando `/load <ID_DEL_LIENZO>`.

### Comandos Disponibles

-   `/p <x> <y> <char>`: Dibuja un carácter (`char`) en la coordenada (`x`, `y`).
    -   Ejemplo: `/p 10 5 X`
-   `/load <canvas_id>`: Carga un lienzo existente o cambia a él.
-   `/id`: Muestra el ID del lienzo actual.
-   `/clear`: Inicia una votación para limpiar el lienzo. Todos los usuarios conectados deben confirmar.
-   `/clear yes`: Emite tu voto para confirmar la limpieza del lienzo.
-   `/help`: Muestra una lista de los comandos disponibles.

## 🛠️ Detalles Técnicos

### Concurrencia

El servidor está diseñado para ser altamente concurrente. Cada conexión de cliente se maneja en su propia goroutine. El estado compartido (como la lista de clientes en un `CanvasGroup`) está protegido por un `sync.RWMutex` para permitir múltiples lecturas concurrentes (broadcasts) pero escrituras exclusivas (añadir/eliminar clientes).

### Persistencia Optimizada

Para minimizar la latencia y la carga en la base de datos, el lienzo no se guarda como un único blob. En su lugar, se divide en "baldosas" de 16x8 caracteres. Cuando un usuario modifica un carácter, solo la baldosa afectada se serializa y se guarda en un **Hash** de Valkey.

-   **Clave en Valkey**: `canvas:<canvas_id>`
-   **Tipo**: `Hash`
-   **Campo del Hash**: Coordenada de la baldosa (ej: `"0,1"`)
-   **Valor del Hash**: Datos binarios de la baldosa serializados con `gob`.

Este enfoque permite actualizaciones atómicas y muy rápidas, siendo ideal para un entorno colaborativo.

### Rate Limiting

Para proteger el servidor, se implementa un algoritmo de **Ventana Deslizante** utilizando un búfer circular por conexión. Este sistema limita el número de comandos que un usuario puede enviar en un período de tiempo determinado, previniendo el spam y asegurando un uso justo de los recursos. La gestión de memoria está controlada, ya que el limitador de un usuario se elimina del mapa de seguimiento en cuanto este se desconecta.


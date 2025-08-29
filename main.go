package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

/*
Servidor TCP

Con soporte para telnet
*/

/*
Mapeo de conexiones de clientes
que tan bueno puede ser esto con miles de usuarios???
canviarlo a channels con worker o hacer sharding map
*/

var (
	clients   = make(map[net.Conn]bool)
	clientsMu sync.RWMutex
)

// Dimensiones del canvas

var canvasWidth = 40
var canvasHeight = 40
var PORT = ":8080"

// Canvas: matriz de caracteres

var canvas = make([][]rune, canvasHeight)

/*
Iniciar canvas Nuevo

Esto se tiene que arreglar, podemos hacer que no se genere todo el canvas
sino que asignar memoria dinamica solo a los bloques se escriban una especie de hashMap
que pasa si dos usuarios quieren escribir al mismo tiempo en el mismo bloque, es posible implementar
un sistema que bloquee solamente esas secciones del canvas.
Como afecta esto al guardado de la informacion en redis?????
Podria aplicar compresion de bloques de ASCII para las secciones???
*/

func initCanvas() {
	for i := 0; i < canvasHeight; i++ {
		canvas[i] = make([]rune, canvasWidth)
		for j := 0; j < canvasWidth; j++ {
			canvas[i][j] = ' '
		}
	}
}

/*
Render canvas
*/

func renderCanvas() string {
	clientsMu.RLock()
	defer clientsMu.RUnlock()

	var output string
	for i := 0; i < canvasHeight; i++ {
		for j := 0; j < canvasWidth; j++ {
			output += string(canvas[i][j])
		}
		output += "\n"
	}
	return output
}

/*
Reenvío de mensajes a todos los clientes conectados.
*/
func broadcast(message string, sender net.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for client := range clients {
		if client != sender {
			client.Write([]byte(message))
		}
	}
}

/*
Para aceptar conexiones entrantes y manejar la comunicación con los clientes.
*/

func handleConnection(conn net.Conn) {
	defer conn.Close()

	//conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	fmt.Println("Nueva conexión desde", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Recibido:", line)

		if !allowCommand(conn) {
			fmt.Println("Demasiados comandos enviados")
			conn.Write([]byte("afloja la moto flaco\n"))
			continue
		}

		if isCommand(line, nil) == 0 {
			fmt.Println("Comando no reconocido")
			conn.Write([]byte("fijate bien que pusiste algo mal\n"))
		}

		//broadcast(conn.RemoteAddr().String()+" :"+line+"\n", conn)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer del cliente:", err)
	}
	fmt.Println("Conexión cerrada desde", conn.RemoteAddr())
}

func main() {

	// Creamos el listener
	initCanvas()
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor escuchando en ", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error al aceptar la conexión:", err)
			continue
		}

		/*
			logica de preguntar si quiere un nuevo canva o unirse a uno existente
			esto implica saber cuales estan cargados y cuales estan en memoria
		*/

		conn.Write([]byte(renderCanvas()))

		clientsMu.Lock()
		clients[conn] = true
		clientsMu.Unlock()

		go handleConnection(conn)
	}
}

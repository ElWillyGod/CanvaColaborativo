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
*/
var (
	clients   = make(map[net.Conn]bool)
	clientsMu sync.RWMutex
)

func main() {

	// Creamos el listener
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor escuchando en :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error al aceptar la conexión:", err)
			continue
		}

		clientsMu.Lock()
		clients[conn] = true
		clientsMu.Unlock()

		go handleConnection(conn)
	}
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

		isCommand(line, nil)

		//broadcast(conn.RemoteAddr().String()+" :"+line+"\n", conn)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer del cliente:", err)
	}
	fmt.Println("Conexión cerrada desde", conn.RemoteAddr())
}

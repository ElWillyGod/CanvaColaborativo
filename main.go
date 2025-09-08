package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

/*
Servidor TCP

Con soporte para telnet
*/

// dimensiones del canvas
var canvasWidth = 80
var canvasHeight = 80
var PORT = ":8080"

/*
Mapeo de conexiones de clientes
que tan bueno puede ser esto con miles de usuarios???
canviarlo a channels con worker o hacer sharding map
*/
var (
	clients   = make(map[net.Conn]bool)
	clientsMu sync.RWMutex
)

type Delta struct {
	X, Y int
	Char rune
}

/*
Estructura del canvas

type Canvas struct {
	ID     string
	Matrix [][]rune
}
*/
/*
Iniciar canvas Nuevo

func initCanvas(id string) *Canvas {
	matrix := make([][]rune, canvasHeight)
	for i := range matrix {
		matrix[i] = make([]rune, canvasWidth)
		for j := range matrix[i] {
			matrix[i][j] = ' '
		}
	}
	return &Canvas{
		ID:     id,
		Matrix: matrix,
	}
}

/*
Reenvío de mensajes a todos los clientes conectados.
*/

func broadcast(message string, sender net.Conn) {
	clientsMu.RLock()
	defer clientsMu.RUnlock()
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
	var canvasGroup *CanvasGroup

	defer func() {
		if canvasGroup != nil {
			canvasGroup.removeClient(conn)
		}
		conn.Close()
		fmt.Println("Conexión cerrada desde", conn.RemoteAddr())
	}()

	fmt.Println("Nueva conexión desde", conn.RemoteAddr())
	conn.Write([]byte("ID del canvas o escribe 'nuevo': "))

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}

	input := strings.TrimSpace(scanner.Text())
	var canvasID string

	if input == "nuevo" {
		canvasID = generateCanvasID()
		canvasGroup = gestCanvas(canvasID)
		//saveCanvasValkey(canvasGroup.Canvas)
		conn.Write([]byte("Canvas creado con ID: " + canvasID + "\n"))
	} else {
		///////////////////////////////////////////////////
		/*
			Implementar el modoelo de guardado hibrido, para guardar las cosas en archivos
			binarios Protobuf o XML
		*/
		///////////////////////////////////////////////////
		canvas, err := loadCanvasFromValkey(input)
		if err != nil {
			canvasID = generateCanvasID()
			canvasGroup = gestCanvas(canvasID)
		} else {
			canvasGroup = gestCanvas(input)
			canvasGroup.Canvas = canvas
		}
		conn.Write([]byte("Canvas ID: " + canvasGroup.Canvas.ID + "\n"))
	}
	fmt.Println("ID: " + canvasGroup.Canvas.ID + "\n")

	canvasGroup.addClient(conn)
	conn.Write([]byte(canvasGroup.renderCanvas()))

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Recibido:", line)

		if !allowCommand(conn) {
			fmt.Println("Demasiados comandos enviados")
			conn.Write([]byte("afloja la moto flaco\n"))
			continue
		}

		// isCommand devuelve 1 si fue un comando de dibujo que modificó el canvas.
		deltas, ok := isCommand(line, []string{conn.RemoteAddr().String()}, canvasGroup)

		if !ok {
			// Si es un mensaje de chat (o comando que no modifica), solo hacer broadcast.
			canvasGroup.broadcast(line+"\n", conn)
		}
		if deltas != nil {
			// Si fue un comando de dibujo (resultado 1):
			// 1. Renderizar el nuevo estado del canvas.
			updateString := deltasAnsi(deltas)

			// 2. Difundir el canvas actualizado a TODOS los clientes.
			canvasGroup.broadcast(updateString, nil) // nil para enviar a todos

			// 3. Guardar el estado en la base de datos.
			err := saveCanvasValkey(canvasGroup.Canvas)
			if err != nil {
				fmt.Println("Error al autoguardar el canvas:", err)
				// Opcional: notificar al usuario del error de guardado.
				conn.Write([]byte("Error al guardar el canvas en la base de datos.\n"))
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer del cliente:", err)
	}
}

func main() {
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

		go handleConnection(conn)
	}
}

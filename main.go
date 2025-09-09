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
var canvasWidth int
var canvasHeight int
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
	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	go func() {
		for message := range client.send {
			client.conn.Write(message)
		}
	}()

	defer func() {

		removeLimiter(conn)

		if canvasGroup != nil {
			canvasGroup.removeClient(client)
			fmt.Printf("%s eliminado del grupo %s.\n", conn.RemoteAddr(), canvasGroup.Canvas.ID)

			canvasGroup.Mutex.RLock()
			isGroupEmpty := len(canvasGroup.Clients) == 0
			canvasGroup.Mutex.RUnlock()

			if isGroupEmpty {
				/*
					fmt.Printf("grupo %s vacio. guardando Valkey \n", canvasGroup.Canvas.ID)

					err := saveCanvasValkey(canvasGroup.Canvas)
					if err != nil {
						fmt.Printf("Error al guardar el canvas %s al cerrar el grupo: %v\n", canvasGroup.Canvas.ID, err)
					}
				*/

				canvasesMu.Lock()
				delete(canvasGroups, canvasGroup.Canvas.ID)
				canvasesMu.Unlock()
			}
		}
		conn.Close()
		fmt.Println("Conexion cerrada desde", conn.RemoteAddr())
	}()

	fmt.Println("Nueva conexion desde", conn.RemoteAddr())

SESSON_LOOP:
	for {

		if canvasGroup != nil {
			canvasGroup.removeClient(client)
		}

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

		canvasGroup.addClient(client)
		conn.Write([]byte(canvasGroup.renderCanvas()))

		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("Recibido:", line)

			if !allowCommand(conn) {
				fmt.Println("Demasiados comandos enviados")
				conn.Write([]byte("afloja la moto flaco\n"))
				continue
			}

			commandResult := isCommand(line, []string{conn.RemoteAddr().String()}, canvasGroup)

			if commandResult == 0 {
				canvasGroup.broadcast([]byte(line+"\n"), client)
			}
			if commandResult == 1 {

				//canvasRendered := canvasGroup.renderCanvas()
				canvasGroup.broadcast([]byte(canvasGroup.renderCanvas()), nil)
			}
			if commandResult == 2 {
				//saveCanvasValkey(canvasGroup.Canvas)
				continue SESSON_LOOP
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error al leer del cliente:", err)
		}
	}
}

func main() {

	canvasWidth = getEnv("CANVAS_WIDTH", 80)
	canvasHeight = getEnv("CANVAS_HEIGHT", 40)

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

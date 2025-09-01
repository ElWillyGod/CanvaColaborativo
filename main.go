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
var canvasWidth = 40
var canvasHeight = 40
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

/*
Canvas activo global - por ahora solo uno
*/
var (
	currentCanvas *Canvas
	canvasMu      sync.RWMutex
)

/*
Estructura del canvas
*/
type Canvas struct {
	ID     string
	Matrix [][]rune
}

/*
Iniciar canvas Nuevo
*/
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
Render canvas
*/
func renderCanvas() string {
	canvasMu.RLock()
	defer canvasMu.RUnlock()

	if currentCanvas == nil {
		return "No hay canvas activo\n"
	}

	var output string
	for i := 0; i < canvasHeight; i++ {
		for j := 0; j < canvasWidth; j++ {
			output += string(currentCanvas.Matrix[i][j])
		}
		output += "\n"
	}
	return output
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
	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
		fmt.Println("Conexión cerrada desde", conn.RemoteAddr())
	}()

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

		if isCommand(line, []string{conn.RemoteAddr().String()}) == 0 {
			fmt.Println("Comando no reconocido")
			conn.Write([]byte("fijate bien que pusiste algo mal\n"))
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer del cliente:", err)
	}
}

func main() {
	// Creamos el listener
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

		// Solicitar ID de canvas o crear uno nuevo
		conn.Write([]byte("ID del canvas o escribe 'nuevo' para crear uno nuevo: "))

		scanner := bufio.NewScanner(conn)
		if !scanner.Scan() {
			conn.Close()
			continue
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "nuevo" {
			// Crear nuevo canvas
			id := generateCanvasID()
			canvasMu.Lock()
			currentCanvas = initCanvas(id)
			canvasMu.Unlock()

			// Guardar en Valkey
			if err := saveCanvasValkey(currentCanvas); err != nil {
				conn.Write([]byte("Error al guardar el canvas: " + err.Error() + "\n"))
				conn.Close()
				continue
			}

			conn.Write([]byte("Canvas nuevo creado con ID: " + id + "\n"))
		} else {
			// Intentar cargar canvas existente
			canvas, err := loadCanvasFromValkey(input)
			if err != nil {
				conn.Write([]byte("No se pudo cargar el canvas. Creando uno nuevo...\n"))
				id := generateCanvasID()
				canvasMu.Lock()
				currentCanvas = initCanvas(id)
				canvasMu.Unlock()
				saveCanvasValkey(currentCanvas)
				conn.Write([]byte("Canvas nuevo creado con ID: " + id + "\n"))
			} else {
				canvasMu.Lock()
				currentCanvas = canvas
				canvasMu.Unlock()
				conn.Write([]byte("Canvas cargado con ID: " + currentCanvas.ID + "\n"))
			}
		}

		// Mostrar el canvas al usuario
		conn.Write([]byte(renderCanvas()))

		// Registrar cliente y manejar conexión
		clientsMu.Lock()
		clients[conn] = true
		clientsMu.Unlock()

		go handleConnection(conn)
	}
}

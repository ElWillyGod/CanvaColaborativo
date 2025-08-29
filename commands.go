package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
	Comandos para la edicion del Canvas
	Seguramente sea una function pointer
	operaciones básicas de dibujo (puntos/figuras simples) y limpieza
*/

type Command struct {
	Name       string
	Parameters []string
	Execute    func(args []string) int
}

var commands = map[string]func(args []string) int{
	"/triangle": triangleCommand,
	"/line":     lineCommand,
	//////////////////////////////////////////////////
	"/set":   setEnvironment,
	"/save":  saveCanvas,
	"/load":  loadCanvas,
	"/clear": clearCanvas,
}

///////////////////////////////////////////////////////
// cositas para el clear con confirmacion

var (
	pendingClear       = false
	clearConfirmations = make(map[string]bool)
	clearStartTime     time.Time
	clearDuration      = 10 * time.Second
	clearMu            sync.Mutex
)

func isCommand(command string, args []string) int {
	// Aquí se puede implementar la lógica para verificar si el comando es válido
	// y devolver y ejecutarlo

	parts := strings.Fields(command)
	if len(parts) == 0 {
		return 0
	}
	cmd := parts[0]
	if fn, ok := commands[cmd]; ok {
		return fn(parts[1:])
	}
	//fmt.Println("Comando no reconocido")
	return 0
}

func triangleCommand(args []string) int {
	// Comando de Triángulo
	if len(args) < 7 {
		return 0
	}
	x1, _ := strconv.Atoi(args[0])
	y1, _ := strconv.Atoi(args[1])
	x2, _ := strconv.Atoi(args[2])
	y2, _ := strconv.Atoi(args[3])
	x3, _ := strconv.Atoi(args[4])
	y3, _ := strconv.Atoi(args[5])
	char := []rune(args[6])[0]

	drawLine(x1, y1, x2, y2, char)
	drawLine(x2, y2, x3, y3, char)
	drawLine(x3, y3, x1, y1, char)

	broadcast(renderCanvas(), nil)
	return 1
}

func lineCommand(args []string) int {
	// Comando de Línea
	fmt.Println("llego a line")
	if len(args) < 5 {

		fmt.Println("Error con los parametros de line")
		return 0
	}
	x1, _ := strconv.Atoi(args[0])
	y1, _ := strconv.Atoi(args[1])
	x2, _ := strconv.Atoi(args[2])
	y2, _ := strconv.Atoi(args[3])
	char := []rune(args[4])[0]

	//fmt.Println("anda?")
	drawLine(x1, y1, x2, y2, char)
	broadcast(renderCanvas(), nil)
	return 1
}

func setEnvironment(args []string) int {
	// Aquí se puede implementar la lógica para establecer el entorno
	return 0
}

func saveCanvas(args []string) int {
	// Aquí se puede implementar la lógica para guardar el canvas
	return 0
}

func loadCanvas(args []string) int {
	// Aquí se puede implementar la lógica para cargar el canvas
	return 0
}

/*
	Aquí se puede implementar la lógica para limpiar el canvas
	necesitamos confirmacion de todos los clientes conectados
	no quiero modificar el handleConnection, entonces voy a meter un timeout
*/

func clearCanvas(args []string) int {

	clearMu.Lock()
	defer clearMu.Unlock()

	return 0
}

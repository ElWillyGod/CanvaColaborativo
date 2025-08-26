package main

import (
	"fmt"
	"strconv"
	"strings"
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
	"/pointer": pointerCommand,
	"/line":    lineCommand,
	"/set":     setEnvironment,
}

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

func pointerCommand(args []string) int {
	// Comando de Punto
	return 0
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

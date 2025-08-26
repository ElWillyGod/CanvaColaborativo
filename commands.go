package main

import (
	"fmt"
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
	"pointer": pointerCommand,
	"line":    lineCommand,
	"/set":    setEnvironment,
}

func isCommand(command string, args []string) int {
	// Aquí se puede implementar la lógica para verificar si el comando es válido
	// y devolver y ejecutarlo
	if cmd, ok := commands[command]; ok {
		return cmd(args)
	}
	return 0
}

func pointerCommand(args []string) int {
	// Comando de Punto
	return 0
}

func lineCommand(args []string) int {
	// Comando de Línea
	fmt.Printf("Dibujando línea con args: %v\n", args)
	return 0
}

func setEnvironment(args []string) int {
	// Aquí se puede implementar la lógica para establecer el entorno
	return 0
}

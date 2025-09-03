package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

/*
	Comandos para la edicion del Canvas
	Seguramente sea una function pointer
	operaciones básicas de dibujo (puntos/figuras simples) y limpieza

	Me falta el comando para el quit del user jaja
*/

type Command struct {
	Name       string
	Parameters []string
	Execute    func(args []string) int
}

var commands = map[string]func(args []string, canvasGroup *CanvasGroup) int{
	"/triangle": triangleCommand,
	"/line":     lineCommand,
	//////////////////////////////////////////////////
	"/set":   setEnvironment,
	"/save":  saveCanvas,
	"/load":  loadCanvas,
	"/clear": clearCanvas,
	"/help":  helpCanvas,
}

///////////////////////////////////////////////////////
// cositas para el clear con confirmacion

var clearDuration = 10 * time.Second

func isCommand(command string, extraArgs []string, canvasGroup *CanvasGroup) int {
	// Aquí se puede implementar la lógica para verificar si el comando es válido
	// y devolver y ejecutarlo

	parts := strings.Fields(command)
	if len(parts) == 0 {
		return 0
	}
	cmd := parts[0]
	if fn, ok := commands[cmd]; ok {
		args := append(parts[1:], extraArgs...)
		return fn(args, canvasGroup) // Pasar el grupo como parámetro
	}
	//fmt.Println("Comando no reconocido")
	return 0
}

func triangleCommand(args []string, canvasGroup *CanvasGroup) int {
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

	drawLine(x1, y1, x2, y2, char, canvasGroup)
	drawLine(x2, y2, x3, y3, char, canvasGroup)
	drawLine(x3, y3, x1, y1, char, canvasGroup)

	canvasGroup.broadcast(canvasGroup.renderCanvas(), nil)
	return 1
}

func lineCommand(args []string, canvasGroup *CanvasGroup) int {
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
	drawLine(x1, y1, x2, y2, char, canvasGroup)
	canvasGroup.broadcast(canvasGroup.renderCanvas(), nil)
	return 1
}

func setEnvironment(args []string, canvasGroup *CanvasGroup) int {
	// Aquí se puede implementar la lógica para establecer el entorno
	return 0
}

func saveCanvas(args []string, canvasGroup *CanvasGroup) int {
	if canvasGroup.Canvas == nil {
		return 0
	}
	err := saveCanvasValkey(canvasGroup.Canvas)
	if err != nil {
		return 0
	}
	return 1
}

func loadCanvas(args []string, canvasGroup *CanvasGroup) int {
	if len(args) < 1 {
		return 0
	}
	id := args[0]
	canvas, err := loadCanvasFromValkey(id)
	if err != nil {
		return 0
	}
	canvasGroup.Canvas = canvas
	return 1
}

/*
	Aquí se puede implementar la lógica para limpiar el canvas
	necesitamos confirmacion de todos los clientes conectados
	no quiero modificar el handleConnection, entonces voy a meter un timeout
*/

func clearCanvas(args []string, canvasGroup *CanvasGroup) int {
	canvasGroup.Mutex.Lock()
	defer canvasGroup.Mutex.Unlock()

	// Si no hay limpieza pendiente, iniciar una nueva
	if !canvasGroup.PendingClear {
		canvasGroup.PendingClear = true
		canvasGroup.ClearConfirmations = make(map[string]bool)
		canvasGroup.ClearStartTime = time.Now()

		go func() {
			canvasGroup.broadcast("Limpieza de canvas iniciada. Todos los usuarios deben confirmar con /clear yes en los proximos 10 segundos.\n", nil)
		}()

		go waitForClearConfirmations(canvasGroup)
		return 1
	}

	// Si hay confirmación pendiente y el usuario responde "yes"
	if len(args) > 0 && args[0] == "yes" && len(args) > 1 {
		userID := args[1]
		canvasGroup.ClearConfirmations[userID] = true
		return 1
	}

	return 0
}

/*
La func de help
*/
func helpCanvas(args []string, canvasGrup *CanvasGroup) int {
	return 0
}

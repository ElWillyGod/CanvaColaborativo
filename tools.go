package main

import "time"

/*
	Algoritmo de Bresenham para dibujar líneas

	Modificar el algoritmo para meter lo de los characters ANSI
*/

func drawLine(x1, y1, x2, y2 int, char rune, canvasGroup *CanvasGroup) {
	dx := abs(x2 - x1)
	dy := -abs(y2 - y1)
	sx := 1
	if x1 >= x2 {
		sx = -1
	}
	sy := 1
	if y1 >= y2 {
		sy = -1
	}
	err := dx + dy

	for {
		if x1 >= 0 && x1 < canvasWidth && y1 >= 0 && y1 < canvasHeight {
			canvasGroup.Canvas.setChar(x1, y1, char)
		}
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x1 += sx
		}
		if e2 <= dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func resetCanvas(canvasGroup *CanvasGroup) {
	canvas := canvasGroup.Canvas
	canvas.mutexCanvas.Lock()
	defer canvas.mutexCanvas.Unlock()

	canvas.tiles = map[TileID]*Tile{}
}

func waitForClearConfirmations(canvasGroup *CanvasGroup) {
	time.Sleep(clearDuration)

	// Versión simple sin mutex anidados
	canvasGroup.Mutex.Lock()
	defer canvasGroup.Mutex.Unlock()

	if !canvasGroup.PendingClear {
		canvasGroup.Mutex.Unlock()
		return
	}

	numClients := len(canvasGroup.Clients)
	numConfirmations := len(canvasGroup.ClearConfirmations)
	canvasGroup.PendingClear = false

	shouldClear := numConfirmations == numClients && numClients > 0

	// Preparar datos necesarios para después del unlock
	var canvasRendered string
	if shouldClear {
		resetCanvas(canvasGroup)
		canvasRendered = canvasGroup.renderCanvas()
	}

	// defer unlock se ejecutará aca
	// Pero necesitamos hacer los broadcasts DESPUÉS del unlock
	// Solución: usar una goroutine
	go func() {
		if shouldClear {
			canvasGroup.broadcast(canvasRendered, nil)
			canvasGroup.broadcast("Canvas limpiado.\n", nil)
			//saveCanvasValkey(canvasGroup.Canvas)
		} else {
			canvasGroup.broadcast("Limpieza cancelada.\n", nil)
		}
	}()
}

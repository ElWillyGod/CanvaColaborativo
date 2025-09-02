package main

import "time"

/*
	Algoritmo de Bresenham para dibujar líneas
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

	canvasGroup.Mutex.Lock()
	defer canvasGroup.Mutex.Unlock()
	for {
		if x1 >= 0 && x1 < canvasWidth && y1 >= 0 && y1 < canvasHeight {
			canvasGroup.Canvas.Matrix[y1][x1] = char
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
	canvasGroup.Mutex.Lock()
	defer canvasGroup.Mutex.Unlock()
	if canvasGroup.Canvas != nil {
		for i := range canvasGroup.Canvas.Matrix {
			for j := range canvasGroup.Canvas.Matrix[i] {
				canvasGroup.Canvas.Matrix[i][j] = ' '
			}
		}
	}
}

func waitForClearConfirmations(canvasGroup *CanvasGroup) {
	time.Sleep(clearDuration)
	clearMu.Lock()
	defer clearMu.Unlock()

	if !pendingClear {
		return
	}

	canvasGroup.Mutex.RLock()
	numClients := len(canvasGroup.Clients)
	canvasGroup.Mutex.RUnlock()

	if len(clearConfirmations) == numClients {
		resetCanvas(canvasGroup)
		canvasGroup.broadcast(canvasGroup.renderCanvas(), nil)
		canvasGroup.broadcast("Canvas limpiado por todos los usuarios.\n", nil)
	} else {
		broadcast("No todos los usuarios han confirmado la limpieza del canvas.\n", nil)
	}
	pendingClear = false
}

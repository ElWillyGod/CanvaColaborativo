package main

import "time"

/*
	Algoritmo de Bresenham para dibujar líneas
*/

func drawLine(x1, y1, x2, y2 int, char rune) {
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
			canvas[y1][x1] = char
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

func waitForClearConfirmations() {
	time.Sleep(clearDuration)
	clearMu.Lock()
	defer clearMu.Unlock()

	if !pendingClear {
		return
	}

	if len(clearConfirmations) == len(clients) {
		initCanvas() //crea un canva nuevo (Canviar por un reset)
		broadcast(renderCanvas(), nil)
		broadcast("Canvas limpiado por todos los usuarios.\n", nil)
	} else {
		broadcast("No todos los usuarios han confirmado la limpieza del canvas.\n", nil)
	}
	pendingClear = false
}

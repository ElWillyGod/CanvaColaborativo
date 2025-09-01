package main

import (
	"net"
	"sync"
)

/*
Canvas group.
esctructuras y logica
*/

type CanvasGroup struct {
	Canvas  *Canvas
	Clients map[net.Conn]bool
	Mutex   sync.RWMutex
}

var (
	canvasGroup = make(map[string]*CanvasGroup)
	canvasesMu  sync.RWMutex
)

func gestCanvas(canvasID string) *CanvasGroup {
	canvasesMu.Lock()
	defer canvasesMu.Unlock()

	if group, exits := canvasGroup[canvasID]; exits {
		return group
	}

	group := &CanvasGroup{
		Canvas:  initCanvas(canvasID),
		Clients: make(map[net.Conn]bool),
	}
	canvasGroup[canvasID] = group
	return group
}

func (cg *CanvasGroup) addClient(conn net.Conn) {
	cg.Mutex.Lock()
	defer cg.Mutex.Unlock()
	cg.Clients[conn] = true
}

func (cg *CanvasGroup) removeClient(conn net.Conn) {
	cg.Mutex.Lock()
	defer cg.Mutex.Unlock()
	delete(cg.Clients, conn)
}

func (cg *CanvasGroup) broadcast(message string, sender net.Conn) {
	cg.Mutex.RLock()
	defer cg.Mutex.RUnlock()
	for client := range cg.Clients {
		if client != sender {
			client.Write([]byte(message))
		}
	}
}

// Renderizar canvas del grupo
func (cg *CanvasGroup) renderCanvas() string {
	if cg.Canvas == nil {
		return "No hay canvas activo\n"
	}

	var output string
	for i := 0; i < canvasHeight; i++ {
		for j := 0; j < canvasWidth; j++ {
			output += string(cg.Canvas.Matrix[i][j])
		}
		output += "\n"
	}
	return output
}

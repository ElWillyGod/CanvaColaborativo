package main

import (
	"net"
	"sync"
	"time"
)

/*
Canvas group.
esctructuras y logica
*/

type CanvasGroup struct {
	Canvas             *Canvas
	Clients            map[net.Conn]bool
	Mutex              sync.RWMutex
	PendingClear       bool
	ClearConfirmations map[string]bool
	ClearStartTime     time.Time
}

var (
	canvasGroups = make(map[string]*CanvasGroup)
	canvasesMu   sync.RWMutex
)

func gestCanvas(canvasID string) *CanvasGroup {
	canvasesMu.Lock()
	defer canvasesMu.Unlock()

	if group, exits := canvasGroups[canvasID]; exits {
		return group
	}

	group := &CanvasGroup{
		Canvas:             newCanvas(canvasID),
		Clients:            make(map[net.Conn]bool),
		PendingClear:       false,
		ClearConfirmations: make(map[string]bool),
	}
	canvasGroups[canvasID] = group
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
// Esto puede cambiar para los ANSI
func (cg *CanvasGroup) renderCanvas() string {
	if cg.Canvas == nil {
		return "No hay canvas activo\n"
	}
	return cg.Canvas.render(canvasWidth, canvasHeight)
}

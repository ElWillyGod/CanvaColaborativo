package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

/*
Canvas group.
esctructuras y logica
*/

const MaxOper = 2

type Client struct {
	conn net.Conn
	send chan []byte
}

type CanvasGroup struct {
	Canvas             *Canvas
	Clients            map[*Client]bool
	Mutex              sync.RWMutex
	Oper               [][]*Delta
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
		Clients:            make(map[*Client]bool),
		Oper:               make([][]*Delta, 0, MaxOper),
		PendingClear:       false,
		ClearConfirmations: make(map[string]bool),
	}
	canvasGroups[canvasID] = group
	return group
}

func (cg *CanvasGroup) addClient(client *Client) {
	cg.Mutex.Lock()
	defer cg.Mutex.Unlock()
	cg.Clients[client] = true
}

func (cg *CanvasGroup) removeClient(client *Client) {
	cg.Mutex.Lock()
	defer cg.Mutex.Unlock()
	delete(cg.Clients, client)
}

func (cg *CanvasGroup) broadcast(message []byte, sender *Client) {
	cg.Mutex.RLock()
	defer cg.Mutex.RUnlock()

	for client := range cg.Clients {
		if client != sender {
			select {
			case client.send <- message:
			default:

				fmt.Printf("chan client %s lleno.\n", client.conn.RemoteAddr())
			}
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

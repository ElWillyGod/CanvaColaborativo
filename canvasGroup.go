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

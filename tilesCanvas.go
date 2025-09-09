package main

import (
	"bytes"
	"fmt"
	"sync"
)

/*
	Idea: dividir el canvas en **tiles** (p. ej. 64×32). Solo asignás memoria para tiles “tocados”.

  	HashMap → Tile (map\[TileID]\*Tile) con pooling (`sync.Pool`) para reciclar tiles y buffers.
  	Dentro de cada tile,
   	matriz de run-length (RLE) por filas para comprimir secuencias de caracteres iguales (ideal para ASCII).
  	Copy-on-write para snapshots (ver §4).
	Beneficio: baja uso de RAM, updates localizados, snapshots y “undos” baratos.

	esto en palabras normales consiste en dividir la matriz en partes y
	actualizar secciones particulares de esta matris, bloqueando esas micro partes

*/

const (
	TileWidth  = 16
	TileHeight = 8
)

type TileID struct {
	X, Y int
}

type Tile struct {
	mutex sync.RWMutex
	Data  []rune
}

type Canvas struct {
	ID          string
	mutexCanvas sync.RWMutex
	tiles       map[TileID]*Tile
}

func newCanvas(id string) *Canvas {
	return &Canvas{
		ID:    id,
		tiles: make(map[TileID]*Tile),
	}
}

func newTile() *Tile {
	t := &Tile{
		Data: make([]rune, TileWidth*TileHeight),
	}
	for i := range t.Data {
		t.Data[i] = ' '
	}

	return t
}

func (c *Canvas) setChar(x, y int, char rune) {
	if x < 0 || y < 0 {
		return
	}

	tileID := TileID{X: x / TileWidth, Y: y / TileHeight}

	c.mutexCanvas.RLock()
	tile, ok := c.tiles[tileID]
	c.mutexCanvas.RUnlock()

	if !ok {
		c.mutexCanvas.Lock()

		tile, ok = c.tiles[tileID]

		if !ok {
			tile = newTile()
			c.tiles[tileID] = tile
		}
		c.mutexCanvas.Unlock()
	}

	tile.mutex.Lock()
	localX := x % TileWidth
	localY := y % TileHeight

	index := localY*TileWidth + localX
	if index < len(tile.Data) {
		tile.Data[index] = char
	}
	tile.mutex.Unlock()

	go func() {
		if err := saveTileToValkey(c.ID, tileID, tile); err != nil {
			fmt.Printf("error en setchar %v: %v\n", tileID, err)
		}
	}()
}

func (c *Canvas) getChar(x, y int) rune {
	if x < 0 || y < 0 {
		return ' '
	}

	tileID := TileID{X: x / TileWidth, Y: y / TileHeight}

	c.mutexCanvas.RLock()
	tile, ok := c.tiles[tileID]
	c.mutexCanvas.RUnlock()

	if !ok {
		return ' '
	}

	tile.mutex.RLock()
	defer tile.mutex.RUnlock()

	localX := x % TileWidth
	localY := y % TileHeight
	index := localY*TileWidth + localX
	if index < len(tile.Data) {
		return tile.Data[index]
	}

	return ' '
}

func (c *Canvas) render(width, height int) string {
	c.mutexCanvas.RLock()
	defer c.mutexCanvas.RUnlock()
	////////////////////////////////////////////////

	var buf bytes.Buffer
	buf.Grow(width * height * 2)

	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)

		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	for id, tile := range c.tiles {
		starX := id.X * TileWidth
		starY := id.Y * TileHeight

		tile.mutex.RLock()

		for i := 0; i < TileHeight; i++ {
			for j := 0; j < TileWidth; j++ {

				absY, absX := starY+i, starX+j

				if absY < height && absX < width {

					grid[absY][absX] = tile.Data[i*TileWidth+j]
				}
			}
		}
		tile.mutex.RUnlock()
	}

	for y := 0; y < height; y++ {
		buf.WriteString(string(grid[y]))
		buf.WriteRune('\n')
	}

	return buf.String()
}

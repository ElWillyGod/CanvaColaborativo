package main

/*
	Persistencia, que linda historia como podemos implementar esto para que me quede encaminado
	para lo de los grupos de canvas?
	Los grupos de canvas no importa aca.
	Politica de guardado:
	se van a crear binarios con el estado actual del canvas con cada comando aplicado
	al vaciarse la sala del canvas se guardara en memoria

	base de datos a usar para la persistencia:
	no SQL porque va a guadar el DI del canvas y el ultimo estado
	puede ser en JSON, BSON, XML Protobuf

*/
import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

/*
Esto tengo que modificarlo para que ande con los tiles nuevos
*/
func generateCanvasID() string {
	return uuid.New().String()
}

func saveTileToValkey(canvasID string, tileID TileID, tile *Tile) error {
	ctx := context.Background()

	tileKey := fmt.Sprintf("%d,%d", tileID.X, tileID.Y)

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(tile.Data); err != nil {
		return fmt.Errorf("error al codificar la baldosa con gob: %w", err)
	}

	err := rdb.HSet(ctx, "canvas:"+canvasID, tileKey, buffer.Bytes()).Err()
	if err != nil {
		return fmt.Errorf("error al guardar la baldosa en redis (HSET): %w", err)
	}

	return nil
}

func loadCanvasFromValkey(id string) (*Canvas, error) {
	ctx := context.Background()
	//HGETALL
	tilesData, err := rdb.HGetAll(ctx, "canvas:"+id).Result()
	if err != nil {
		return nil, err
	}
	if len(tilesData) == 0 {
		return newCanvas(id), nil
	}

	canvas := newCanvas(id)

	for tileKey, tileData := range tilesData {
		var tileID TileID
		if _, err := fmt.Sscanf(tileKey, "%d,%d", &tileID.X, &tileID.Y); err != nil {
			fmt.Printf("error parsear tile key '%s', saltando: %v\n", tileKey, err)
			continue
		}

		var data []rune
		buffer := bytes.NewBufferString(tileData)
		if err := gob.NewDecoder(buffer).Decode(&data); err != nil {
			fmt.Printf("error decodificar tile data para key '%s', saltando: %v\n", tileKey, err)
			continue
		}

		canvas.tiles[tileID] = &Tile{Data: data}
	}

	fmt.Printf("%s cargado con %d baldosas.\n", id, len(canvas.tiles))
	return canvas, nil
}

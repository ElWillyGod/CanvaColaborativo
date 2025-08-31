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
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func generateCanvasID() string {
	return uuid.New().String()
}

func saveCanvasValkey(canvas *Canvas) error {
	ctx := context.Background()
	data, err := json.Marshal(canvas)
	if err != nil {
		return err
	}
	err = rdb.Set(ctx, "canvas:"+canvas.ID, data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func loadCanvasFromValkey(id string) (*Canvas, error) {
	ctx := context.Background()
	data, err := rdb.Get(ctx, "canvas:"+id).Result()
	if err != nil {
		return nil, err
	}

	var canvas Canvas
	err = json.Unmarshal([]byte(data), &canvas)
	if err != nil {
		return nil, err
	}
	return &canvas, nil
}

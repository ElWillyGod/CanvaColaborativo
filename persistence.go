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

func saveCanvasValkey(canvas *Canvas) error {
	ctx := context.Background()

	canvas.mutex.RLock()
	defer canvas.mutex.RUnlock()

	// Usamos un buffer para la codificación binaria
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	// Codificamos el mapa de tiles. Gob sí soporta struct como clave.
	if err := encoder.Encode(canvas.tiles); err != nil {
		return fmt.Errorf("error al codificar con gob: %w", err)
	}

	// Guardamos los bytes del buffer en Redis.
	err := rdb.Set(ctx, "canvas:"+canvas.ID, buffer.Bytes(), 0).Err()
	if err != nil {
		return fmt.Errorf("error al guardar en redis: %w", err)
	}
	fmt.Printf("Canvas %s guardado.\n", canvas.ID)
	return nil
}

func loadCanvasFromValkey(id string) (*Canvas, error) {
	ctx := context.Background()
	// Obtenemos los datos como bytes
	data, err := rdb.Get(ctx, "canvas:"+id).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("canvas con ID '%s' no encontrado", id)
		}
		return nil, err
	}

	canvas := newCanvas(id)

	// Creamos un buffer a partir de los bytes y un decodificador gob
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	// Decodificamos los datos en el mapa de tiles.
	if err := decoder.Decode(&canvas.tiles); err != nil {
		return nil, fmt.Errorf("error al decodificar con gob: %w", err)
	}

	fmt.Printf("Canvas %s cargado.\n", id)
	return canvas, nil
}

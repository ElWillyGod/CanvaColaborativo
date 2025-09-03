package main

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

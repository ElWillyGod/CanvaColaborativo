package canvacolaborativo

/*
	Comandos para la edicion del Canvas
	Seguramente sea una function pointer
	operaciones básicas de dibujo (puntos/figuras simples) y limpieza
*/

type Command struct {
	Name       string
	Parameters []string
	Execute    func() int
}

func isCommand(command string) int {
	// Aquí se puede implementar la lógica para verificar si el comando es válido
	// y devolver y ejecutarlo

	return 0
}

func pointerCommand(command string) int {
	// Aquí se puede implementar la lógica para ejecutar el comando
	return 0
}

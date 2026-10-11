// Comando docgen: genera la referencia de comandos de dots en Markdown.
//
// Recorre el árbol de cobra del motor (el mismo que se ejecuta de verdad, ya
// en castellano) y escribe una página por comando, más un índice, en
// docs/reference/commands/. Así la referencia sale de la propia ayuda y no se
// desfasa: el CI la regenera y falla si no coincide con lo versionado.
//
// No usa cobra/doc porque sus encabezados ("Synopsis", "Options", "SEE ALSO"...)
// vienen fijos en inglés y añade una fecha de generación que haría cambiar los
// ficheros cada día. Generarlo a mano son pocas líneas y queda a nuestro gusto.
//
// Uso (desde la raíz del repo): go run ./tools/docgen
package main

//go:generate go run . -out ../../docs/reference/commands

import (
	"flag"
	"fmt"
	"os"

	"go.kreikol.dev/dotsmithy/internal/cli"
)

func main() {
	out := flag.String("out", "docs/reference/commands", "directorio donde escribir la referencia")
	flag.Parse()

	if err := generate(cli.Root(), *out); err != nil {
		fmt.Fprintln(os.Stderr, "docgen:", err)
		os.Exit(1)
	}
}

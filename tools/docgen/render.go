package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// header encabeza cada fichero generado: avisa de que no se edita a mano.
const header = "<!-- Generado por tools/docgen: no lo edites a mano. -->\n\n"

// generate escribe en dir una página por comando documentable del árbol de
// root, más un README.md con el índice. Antes borra los .md que hubiera, para
// que un comando eliminado no deje su página huérfana (el directorio entero es
// generado).
func generate(root *cobra.Command, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	old, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return err
	}
	for _, f := range old {
		if err := os.Remove(f); err != nil {
			return err
		}
	}

	cmds := collect(root)
	for _, c := range cmds {
		if err := os.WriteFile(filepath.Join(dir, fileName(c)), []byte(page(c)), 0o644); err != nil {
			return err
		}
	}
	return os.WriteFile(filepath.Join(dir, "README.md"), []byte(index(cmds)), 0o644)
}

// collect devuelve, en orden de recorrido (padre antes que hijos, hermanos por
// nombre como los ordena cobra), todos los comandos que se documentan.
func collect(c *cobra.Command) []*cobra.Command {
	out := []*cobra.Command{c}
	for _, sub := range c.Commands() {
		if documented(sub) {
			out = append(out, collect(sub)...)
		}
	}
	return out
}

// documented indica si un subcomando tiene página: fuera los ocultos, los
// obsoletos y el propio help (su ayuda ya es esta referencia).
func documented(c *cobra.Command) bool {
	return c.IsAvailableCommand() && c.Name() != "help"
}

// fileName da el nombre de fichero de un comando: «dots.md» para la raíz y la
// ruta sin la raíz, unida con guiones, para el resto («init.md»,
// «completion-bash.md»).
func fileName(c *cobra.Command) string {
	if !c.HasParent() {
		return c.Name() + ".md"
	}
	path := strings.Fields(c.CommandPath())[1:]
	return strings.Join(path, "-") + ".md"
}

// page monta la página Markdown de un comando, con la misma estructura que su
// --help: descripción, uso, ejemplos, subcomandos, opciones y opciones globales.
func page(c *cobra.Command) string {
	var b strings.Builder
	b.WriteString(header)
	fmt.Fprintf(&b, "# %s\n\n", c.CommandPath())
	lead := sentence(c.Short)
	fmt.Fprintf(&b, "%s\n\n", lead)
	if long := longHelp(c.Long, lead); long != "" {
		b.WriteString(prose(long) + "\n\n")
	}

	b.WriteString("## Uso\n\n```\n")
	if c.Runnable() {
		b.WriteString(useLine(c) + "\n")
	}
	if c.HasAvailableSubCommands() {
		b.WriteString(c.CommandPath() + " [comando]\n")
	}
	b.WriteString("```\n\n")

	if c.HasExample() {
		fmt.Fprintf(&b, "## Ejemplos\n\n```sh\n%s\n```\n\n", strings.TrimRight(c.Example, "\n"))
	}

	var subs []*cobra.Command
	for _, sub := range c.Commands() {
		if documented(sub) {
			subs = append(subs, sub)
		}
	}
	if len(subs) > 0 {
		b.WriteString("## Comandos\n\n| Comando | Qué hace |\n| --- | --- |\n")
		for _, sub := range subs {
			fmt.Fprintf(&b, "| [`%s`](%s) | %s |\n", sub.Name(), fileName(sub), cell(sub.Short))
		}
		b.WriteString("\n")
	}

	if c.HasAvailableLocalFlags() {
		b.WriteString("## Opciones\n\n" + flagTable(c.LocalFlags()) + "\n")
	}
	if c.HasAvailableInheritedFlags() {
		b.WriteString("## Opciones globales\n\n" + flagTable(c.InheritedFlags()) + "\n")
	}

	b.WriteString("## Ver también\n\n")
	if c.HasParent() {
		p := c.Parent()
		fmt.Fprintf(&b, "- [`%s`](%s): %s\n", p.CommandPath(), fileName(p), p.Short)
	}
	b.WriteString("- [Todos los comandos](README.md)\n")
	return b.String()
}

// index monta el README.md del directorio: una tabla con todos los comandos.
func index(cmds []*cobra.Command) string {
	var b strings.Builder
	b.WriteString(header)
	b.WriteString("# Comandos\n\n")
	b.WriteString("Todos los comandos de `dots`, con sus opciones. Es lo mismo que ves con\n")
	b.WriteString("`dots <comando> --help`, aquí a mano para consultarlo.\n\n")
	b.WriteString("| Comando | Qué hace |\n| --- | --- |\n")
	for _, c := range cmds {
		fmt.Fprintf(&b, "| [`%s`](%s) | %s |\n", c.CommandPath(), fileName(c), cell(c.Short))
	}
	return b.String()
}

// useLine es la línea de uso con «[opciones]» en vez del «[flags]» de cobra,
// igual que en la ayuda de la terminal.
func useLine(c *cobra.Command) string {
	return strings.Replace(c.UseLine(), " [flags]", " [opciones]", 1)
}

// flagTable pinta las opciones como tabla: atajo y nombre (con el tipo de valor
// si no es booleana) y su descripción, más el valor por defecto si no es el
// vacío.
func flagTable(fs *pflag.FlagSet) string {
	var b strings.Builder
	b.WriteString("| Opción | Qué hace |\n| --- | --- |\n")
	fs.VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		varName, usage := pflag.UnquoteUsage(f)
		name := "--" + f.Name
		if varName != "" {
			name += " <" + valueName(varName) + ">"
		}
		opt := "`" + name + "`"
		if f.Shorthand != "" {
			opt = "`-" + f.Shorthand + "`, " + opt
		}
		if !zeroDefault(f) {
			usage += fmt.Sprintf(" (por defecto: `%s`)", f.DefValue)
		}
		fmt.Fprintf(&b, "| %s | %s |\n", opt, cell(usage))
	})
	return b.String()
}

// valueName traduce el tipo que pflag pone al valor de una opción («string»...)
// a algo legible en castellano.
func valueName(t string) string {
	switch t {
	case "string":
		return "texto"
	case "strings", "stringSlice", "stringArray":
		return "textos"
	case "int", "int64", "uint":
		return "número"
	case "duration":
		return "duración"
	}
	return t
}

// zeroDefault indica si el valor por defecto de una opción es el «vacío» de su
// tipo (y, por tanto, no merece la pena enseñarlo).
func zeroDefault(f *pflag.Flag) bool {
	switch f.DefValue {
	case "", "false", "0", "[]", "<nil>":
		return true
	}
	return false
}

// longHelp devuelve la ayuda larga sin espacios sobrantes y sin su primer
// párrafo si este repite la frase de entrada (pasa en los subcomandos de
// completion, que empiezan su ayuda con la misma frase que su Short).
func longHelp(long, lead string) string {
	long = strings.TrimSpace(long)
	first, rest, _ := strings.Cut(long, "\n\n")
	if escape(strings.TrimSpace(first)) == lead {
		return strings.TrimSpace(rest)
	}
	return long
}

// sentence pone en mayúscula la primera letra y cierra con punto. Los Short de
// cobra van en minúscula y sin punto (así quedan bien en las listas de la
// terminal), pero como frase suelta de una página quedan mejor así.
func sentence(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	// El nombre del producto va siempre en minúscula, también al empezar.
	if !strings.HasPrefix(s, "dotsmithy") {
		r, size := utf8.DecodeRuneInString(s)
		s = string(unicode.ToUpper(r)) + s[size:]
	}
	if !strings.HasSuffix(s, ".") {
		s += "."
	}
	return escape(s)
}

// prose convierte la ayuda larga de un comando en Markdown: las líneas
// sangradas con tabulador (los ejemplos de comandos de la ayuda) pasan a un
// bloque de código; el resto es texto, con los «<» y «>» escapados para que
// cosas como «machines/<máquina>» no se tomen por HTML.
func prose(s string) string {
	var out []string
	inCode := false
	for _, line := range strings.Split(s, "\n") {
		code := strings.HasPrefix(line, "\t")
		switch {
		case code && !inCode:
			out = append(out, "```sh")
			inCode = true
		case !code && inCode && strings.TrimSpace(line) != "":
			out = append(out, "```")
			inCode = false
		}
		if code {
			out = append(out, strings.TrimPrefix(line, "\t"))
		} else if inCode {
			// Línea en blanco dentro de un bloque: puede que siga más código.
			out = append(out, line)
		} else {
			out = append(out, escape(line))
		}
	}
	if inCode {
		out = append(out, "```")
	}
	return tidyCode(strings.Join(out, "\n"))
}

// tidyCode quita las líneas en blanco que quedan justo antes del cierre de un
// bloque de código (las que separaban el bloque del texto siguiente en la ayuda).
func tidyCode(s string) string {
	for strings.Contains(s, "\n\n```\n") {
		s = strings.ReplaceAll(s, "\n\n```\n", "\n```\n\n")
	}
	if strings.HasSuffix(s, "\n\n```") {
		s = strings.TrimSuffix(s, "\n\n```") + "\n```"
	}
	return s
}

// cell prepara un texto para una celda de tabla: escapado y sin «|».
func cell(s string) string {
	return strings.ReplaceAll(escape(s), "|", `\|`)
}

// escape protege los «<» y «>» del texto, que Markdown tomaría por etiquetas
// HTML (y GitHub se las comería).
func escape(s string) string {
	return strings.NewReplacer("<", "&lt;", ">", "&gt;").Replace(s)
}

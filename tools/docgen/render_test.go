package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.kreikol.dev/dotsmithy/internal/cli"
)

// TestGenerate genera la referencia del árbol real en un directorio temporal y
// comprueba lo básico: una página por comando con nombre limpio, cabecera de
// fichero generado, nada en inglés de cobra y enlaces que apuntan a ficheros
// que existen.
func TestGenerate(t *testing.T) {
	dir := t.TempDir()
	// Un fichero viejo tiene que desaparecer: el directorio es todo generado.
	if err := os.WriteFile(filepath.Join(dir, "stale.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := generate(cli.Root(), dir); err != nil {
		t.Fatalf("generate: %v", err)
	}

	for _, want := range []string{"README.md", "dots.md", "init.md", "link.md", "add.md", "completion-bash.md"} {
		if _, err := os.Stat(filepath.Join(dir, want)); err != nil {
			t.Errorf("falta %s: %v", want, err)
		}
	}
	for _, unwanted := range []string{"stale.md", "help.md"} {
		if _, err := os.Stat(filepath.Join(dir, unwanted)); err == nil {
			t.Errorf("sobra %s", unwanted)
		}
	}

	files, _ := filepath.Glob(filepath.Join(dir, "*.md"))
	english := []string{"[flags]", "Synopsis", "SEE ALSO", "Auto generated", "Options inherited", "(default "}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		text := string(data)
		if !strings.HasPrefix(text, header) {
			t.Errorf("%s: falta la cabecera de fichero generado", filepath.Base(f))
		}
		for _, e := range english {
			if strings.Contains(text, e) {
				t.Errorf("%s: contiene %q", filepath.Base(f), e)
			}
		}
		for _, link := range links(text) {
			if _, err := os.Stat(filepath.Join(dir, link)); err != nil {
				t.Errorf("%s: enlace roto a %s", filepath.Base(f), link)
			}
		}
	}
}

// links saca los destinos de los enlaces Markdown relativos a .md de un texto.
func links(text string) []string {
	var out []string
	for _, part := range strings.Split(text, "](")[1:] {
		target, _, ok := strings.Cut(part, ")")
		if ok && strings.HasSuffix(target, ".md") && !strings.Contains(target, "/") {
			out = append(out, target)
		}
	}
	return out
}

// TestProse comprueba la conversión de la ayuda larga: el texto escapado y las
// líneas sangradas con tabulador en un bloque de código, sin escapar.
func TestProse(t *testing.T) {
	in := "Va en machines/<máquina>.\n\nPara cargarlo:\n\n\tsource <(dots completion bash)\n\nY listo."
	want := "Va en machines/&lt;máquina&gt;.\n\nPara cargarlo:\n\n```sh\nsource <(dots completion bash)\n```\n\nY listo."
	if got := prose(in); got != want {
		t.Errorf("prose:\n got %q\nwant %q", got, want)
	}
}

// TestSentence comprueba la frase de entrada: mayúscula y punto, salvo el nombre
// del producto, que va siempre en minúscula.
func TestSentence(t *testing.T) {
	cases := map[string]string{
		"trae los cambios":            "Trae los cambios.",
		"ya acaba en punto.":          "Ya acaba en punto.",
		"dotsmithy: la herrería":      "dotsmithy: la herrería.",
		"despliega a machines/<algo>": "Despliega a machines/&lt;algo&gt;.",
	}
	for in, want := range cases {
		if got := sentence(in); got != want {
			t.Errorf("sentence(%q) = %q, want %q", in, got, want)
		}
	}
}

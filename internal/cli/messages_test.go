package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// runRoot ejecuta la raíz con los argumentos dados y devuelve lo que ha sacado
// por pantalla y el error. Después deja --help a false: la raíz es global y, si
// no, un test que pide ayuda contaminaría a los siguientes.
func runRoot(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs(args)
	t.Cleanup(func() {
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
		rootCmd.SetArgs(nil)
		_ = rootCmd.PersistentFlags().Set("help", "false")
	})
	err := rootCmd.Execute()
	return out.String(), err
}

// TestUsageErrorsInSpanish comprueba que los errores de uso (argumentos y
// opciones) salen en castellano, con la pista de dónde mirar, y sin restos del
// inglés de cobra/pflag.
func TestUsageErrorsInSpanish(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{"faltan argumentos", []string{"add", "tldr"},
			"a «dots add» le faltan argumentos: necesita al menos 2 y le has pasado 1"},
		{"comando desconocido", []string{"instal"},
			`no conozco el comando "instal" (mira «dots --help»)`},
		{"comando parecido", []string{"lnk"},
			"¿querías decir «dots link»?"},
		{"argumento de más", []string{"link", "sobra"},
			`«dots link» no admite argumentos y le has pasado "sobra"`},
		{"opción larga desconocida", []string{"link", "--foo"},
			"no conozco la opción --foo (mira «dots link --help»)"},
		{"atajo desconocido", []string{"link", "-x"},
			"no conozco la opción -x"},
		{"opción sin valor", []string{"init", "--profile"},
			"a la opción --profile le falta el valor"},
		{"valor no válido", []string{"update", "--backup=quizas"},
			`el valor "quizas" no vale para --backup: tiene que ser true o false`},
		{"help de algo que no existe", []string{"help", "nada"},
			`no conozco el comando "nada"`},
		{"argumento de más en completion", []string{"completion", "zsh", "sobra"},
			"«dots completion zsh» no admite argumentos"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := runRoot(t, tc.args...)
			if err == nil {
				t.Fatalf("dots %v: esperaba un error y no lo hay", tc.args)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("dots %v:\n  quiero que contenga: %s\n  tengo:               %s", tc.args, tc.want, err)
			}
			assertNoEnglish(t, err.Error())
		})
	}
}

// TestHelpOfEveryCommandInSpanish recorre todo el árbol de comandos (incluidos
// help y completion, que crea cobra) y comprueba que ninguna ayuda deja textos
// en inglés de cobra. Es la red de seguridad: si un comando nuevo o una versión
// nueva de cobra cuela algo, salta aquí.
func TestHelpOfEveryCommandInSpanish(t *testing.T) {
	// cobra añade el comando help al ejecutar; lo forzamos para que esté en el
	// árbol pase lo que pase con el orden de los tests.
	rootCmd.InitDefaultHelpCmd()
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		if c.Hidden {
			return
		}
		var out bytes.Buffer
		c.SetOut(&out)
		if err := c.Help(); err != nil {
			t.Errorf("%s: la ayuda ha fallado: %v", c.CommandPath(), err)
		}
		c.SetOut(nil)
		help := out.String()
		if !strings.Contains(help, "Uso:") {
			t.Errorf("%s: la ayuda no tiene la sección «Uso:»:\n%s", c.CommandPath(), help)
		}
		assertNoEnglish(t, help)
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}
	walk(rootCmd)
}

// TestRootWithoutArgsShowsHelp: «dots» a secas enseña la ayuda, sin error.
func TestRootWithoutArgsShowsHelp(t *testing.T) {
	out, err := runRoot(t)
	if err != nil {
		t.Fatalf("dots sin argumentos no debería fallar: %v", err)
	}
	if !strings.Contains(out, "Comandos:") {
		t.Errorf("dots sin argumentos debería enseñar la ayuda; ha sacado:\n%s", out)
	}
}

// englishMarkers son textos de cobra/pflag que delatan que algo se ha quedado
// sin traducir.
var englishMarkers = []string{
	"Usage:", "Available Commands", "Flags:", "Global Flags", "[flags]",
	"[command]", "for more information", "help for", "version for",
	"Help about", "Generate the autocompletion", "disable completion",
	"unknown command", "unknown flag", "unknown shorthand", "requires at least",
	"accepts ", "flag needs an argument", "invalid argument", "Did you mean",
}

// assertNoEnglish falla si el texto contiene alguno de los marcadores en inglés.
func assertNoEnglish(t *testing.T, text string) {
	t.Helper()
	for _, m := range englishMarkers {
		if strings.Contains(text, m) {
			t.Errorf("queda inglés de cobra/pflag (%q) en:\n%s", m, text)
		}
	}
}

package cli

// Todo lo que cobra y pflag le dicen a la usuaria por su cuenta, pasado a
// castellano y al tono del proyecto.
//
// Estrategia global (en vez de parchear comando a comando):
//   - validadores de argumentos propios (noArgs, minArgs, subcommandArgs) en
//     lugar de los de cobra (cobra.NoArgs, cobra.MinimumNArgs…);
//   - una función de errores de flags que traduce los errores tipados de pflag;
//   - plantilla de uso, flags -h/--version, comando help y comando completion
//     con textos propios.
//
// Todo se aplica una vez sobre la raíz con localize(); los subcomandos lo
// heredan (la plantilla de uso y la función de errores de flags se buscan
// subiendo por el árbol).

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// usageTemplate es la plantilla de uso de cobra traducida. Mantiene la misma
// estructura que la original (salvo grupos de comandos y temas de ayuda
// adicionales, que no usamos).
const usageTemplate = `Uso:{{if .Runnable}}
  {{useLine .}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [comando]{{end}}{{if gt (len .Aliases) 0}}

Alias:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Ejemplos:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Comandos:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Opciones:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Opciones globales:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableSubCommands}}

Usa «{{.CommandPath}} [comando] --help» para ver la ayuda de un comando.{{end}}
`

// localize aplica los textos en castellano a la raíz (y, por herencia, a todo el
// árbol). Hay que llamarla después de registrar los subcomandos, porque el
// comando completion solo se crea si la raíz ya tiene hijos.
func localize(root *cobra.Command) {
	// La plantilla usa useLine en vez de .UseLine (que pega un «[flags]» fijo).
	cobra.AddTemplateFunc("useLine", useLine)
	root.SetUsageTemplate(usageTemplate)

	// Distancia para sugerir comandos parecidos (la de cobra por defecto; él solo
	// la fija en su propio mensaje de error, que ya no usamos).
	root.SuggestionsMinimumDistance = 2
	root.SetVersionTemplate("dots {{.Version}}\n")
	root.SetFlagErrorFunc(flagError)

	// cobra solo crea sus flags -h y --version si no existen ya; declarándolos
	// nosotros antes, nos quedamos con el texto. -h va persistente para que valga
	// en todos los subcomandos. --version va sin atajo (-v es --verbose).
	root.PersistentFlags().BoolP("help", "h", false, "enseña esta ayuda")
	root.Flags().Bool("version", false, "enseña la versión de dots")

	root.SetHelpCommand(newHelpCmd())
	localizeCompletion(root)
}

// helpHint es la coletilla que cierra los errores de uso: dónde mirar.
func helpHint(cmd *cobra.Command) string {
	return fmt.Sprintf(" (mira «%s --help»)", cmd.CommandPath())
}

// useLine devuelve la línea de uso de un comando con «[opciones]» en lugar del
// «[flags]» que pega cobra. Ej.: "dots add <gestor> <paquete>... [opciones]".
func useLine(cmd *cobra.Command) string {
	return strings.Replace(cmd.UseLine(), " [flags]", " [opciones]", 1)
}

// noArgs es el validador de los comandos que no admiten argumentos sueltos.
// Sustituye a cobra.NoArgs y, sobre todo, evita que un argumento de más se
// acepte en silencio (lo que hace cobra si no se declara Args).
func noArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("«%s» no admite argumentos y le has pasado %q%s",
			cmd.CommandPath(), strings.Join(args, " "), helpHint(cmd))
	}
	return nil
}

// minArgs devuelve un validador que exige al menos n argumentos (sustituye a
// cobra.MinimumNArgs). El mensaje enseña el uso para que se vea qué falta.
func minArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			return fmt.Errorf("a «%s» le faltan argumentos: necesita al menos %d y le has pasado %d, el uso es «%s»%s",
				cmd.CommandPath(), n, len(args), useLine(cmd), helpHint(cmd))
		}
		return nil
	}
}

// subcommandArgs es el validador de los comandos que solo agrupan subcomandos
// (la raíz): cualquier argumento suelto es un comando que no existe. Si se
// parece a uno de verdad, lo sugerimos.
func subcommandArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	msg := fmt.Sprintf("no conozco el comando %q", args[0])
	if sugg := cmd.SuggestionsFor(args[0]); len(sugg) > 0 {
		msg += fmt.Sprintf("; ¿querías decir «%s %s»?", cmd.CommandPath(), sugg[0])
	}
	return fmt.Errorf("%s%s", msg, helpHint(cmd))
}

// flagError traduce los errores de pflag al parsear las opciones. pflag los da
// tipados, así que no hace falta comparar textos. Lo que no reconozcamos se deja
// tal cual, envuelto en una frase nuestra.
func flagError(cmd *cobra.Command, err error) error {
	var (
		notExist *pflag.NotExistError
		noValue  *pflag.ValueRequiredError
		invalid  *pflag.InvalidValueError
		syntax   *pflag.InvalidSyntaxError
		msg      string
	)
	switch {
	case errors.As(err, &notExist):
		// Para un atajo (-x), pflag deja en «shortnames» lo que falla; si no, es
		// una opción larga (--foo).
		if notExist.GetSpecifiedShortnames() != "" {
			msg = fmt.Sprintf("no conozco la opción -%s", notExist.GetSpecifiedName())
		} else {
			msg = fmt.Sprintf("no conozco la opción --%s", notExist.GetSpecifiedName())
		}
	case errors.As(err, &noValue):
		msg = fmt.Sprintf("a la opción --%s le falta el valor", noValue.GetFlag().Name)
	case errors.As(err, &invalid):
		f := invalid.GetFlag()
		msg = fmt.Sprintf("el valor %q no vale para --%s", invalid.GetValue(), f.Name)
		if f.Value.Type() == "bool" {
			msg += ": tiene que ser true o false"
		}
	case errors.As(err, &syntax):
		msg = fmt.Sprintf("opción mal escrita: %s", syntax.GetSpecifiedFlag())
	default:
		msg = fmt.Sprintf("algo no cuadra en las opciones (%v)", err)
	}
	return fmt.Errorf("%s%s", msg, helpHint(cmd))
}

// newHelpCmd es nuestro comando «help», que sustituye al de cobra. Enseña la
// ayuda del comando que se le pida y, si no existe, da un error en castellano
// (el de cobra imprime "Unknown help topic" y el uso).
func newHelpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "help [comando]",
		Short: "enseña la ayuda de cualquier comando",
		Long: `help enseña la ayuda de cualquier comando.
Escribe «dots help <comando>» para ver los detalles de ese comando.`,
		// Autocompletado: los nombres de los comandos del nivel que toque.
		ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			parent, _, err := c.Root().Find(args)
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			var out []cobra.Completion
			for _, sub := range parent.Commands() {
				if sub.IsAvailableCommand() && strings.HasPrefix(sub.Name(), toComplete) {
					out = append(out, cobra.CompletionWithDesc(sub.Name(), sub.Short))
				}
			}
			return out, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(c *cobra.Command, args []string) error {
			target, rest, err := c.Root().Find(args)
			if err != nil || len(rest) > 0 {
				return fmt.Errorf("no conozco el comando %q%s", strings.Join(args, " "), helpHint(c.Root()))
			}
			return target.Help()
		},
	}
}

// completionLongs son las ayudas largas, en castellano, del comando completion y
// de sus subcomandos por shell. %[1]s es el nombre del binario. Solo damos las
// instrucciones de Linux: el motor es Fedora-first.
var completionLongs = map[string]string{
	"completion": `Genera el script de autocompletado de %[1]s para la shell que elijas.
Mira la ayuda de cada subcomando para ver cómo usar el script generado.
`,
	"bash": `Genera el script de autocompletado para bash.

Necesita el paquete «bash-completion»; si no lo tienes, instálalo con el gestor
de paquetes de tu sistema.

Para cargarlo en la sesión actual:

	source <(%[1]s completion bash)

Para tenerlo en todas las sesiones nuevas, ejecuta una vez:

	%[1]s completion bash > /etc/bash_completion.d/%[1]s

Abre una shell nueva para que surta efecto.
`,
	"zsh": `Genera el script de autocompletado para zsh.

Si no tienes el autocompletado activado en tu entorno, actívalo una vez con:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

Para cargarlo en la sesión actual:

	source <(%[1]s completion zsh)

Para tenerlo en todas las sesiones nuevas, ejecuta una vez:

	%[1]s completion zsh > "${fpath[1]}/_%[1]s"

Abre una shell nueva para que surta efecto.
`,
	"fish": `Genera el script de autocompletado para fish.

Para cargarlo en la sesión actual:

	%[1]s completion fish | source

Para tenerlo en todas las sesiones nuevas, ejecuta una vez:

	%[1]s completion fish > ~/.config/fish/completions/%[1]s.fish

Abre una shell nueva para que surta efecto.
`,
	"powershell": `Genera el script de autocompletado para powershell.

Para cargarlo en la sesión actual:

	%[1]s completion powershell | Out-String | Invoke-Expression

Para tenerlo en todas las sesiones nuevas, añade la salida del comando anterior
a tu perfil de powershell.
`,
}

// localizeCompletion crea ya el comando completion de cobra (si no, lo crearía
// él al ejecutar, con sus textos) y le pone los nuestros: descripciones, ayudas
// largas, el flag --no-descriptions y validadores de argumentos en castellano.
func localizeCompletion(root *cobra.Command) {
	root.InitDefaultCompletionCmd()

	var completion *cobra.Command
	for _, c := range root.Commands() {
		if c.Name() == "completion" {
			completion = c
			break
		}
	}
	if completion == nil {
		return
	}

	name := root.Name()
	completion.Short = "genera el script de autocompletado para tu shell"
	completion.Long = fmt.Sprintf(completionLongs["completion"], name)
	completion.Args = subcommandArgs

	for _, sh := range completion.Commands() {
		sh.Short = fmt.Sprintf("genera el script de autocompletado para %s", sh.Name())
		if long, ok := completionLongs[sh.Name()]; ok {
			sh.Long = fmt.Sprintf(long, name)
		}
		sh.Args = noArgs
		if f := sh.Flags().Lookup("no-descriptions"); f != nil {
			f.Usage = "sin descripciones en las sugerencias"
		}
	}
}

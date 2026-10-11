package cli

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"go.kreikol.dev/dotsmithy/internal/manifest"
	"go.kreikol.dev/dotsmithy/internal/packages"
)

// newAddCmd: da de alta uno o varios paquetes: los registra en la lista de una
// capa del perfil activo (por defecto la de su máquina) y los instala. Así el
// alta queda declarada en el repo y aplicada en la máquina de una vez.
func newAddCmd() *cobra.Command {
	var layer string
	var contentDir string
	var profile string

	cmd := &cobra.Command{
		Use:   "add <gestor> <paquete>...",
		Short: "registra e instala uno o varios paquetes",
		Long: `add coge uno o varios paquetes de un gestor (ej. dnf), los registra en la
lista de una capa del contenido y los instala. Deja el alta declarada en el repo
y aplicada en la máquina.

Por defecto los registra en la capa de la máquina del perfil activo
(machines/<máquina>), que es lo más prudente: así no se cuelan en otras
máquinas. Con --layer eliges otra (shared, ux/terminal...), pero tiene que ser
una de las capas del perfil activo, porque registrar en una capa que este
perfil no aplica no instalaría nada aquí.

Ejemplos:
  dots add dnf ripgrep fzf
  dots add dnf tmux --layer shared`,
		Args: minArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			mgrName, pkgs := args[0], args[1:]

			mgr, ok := packages.Get(mgrName)
			if !ok {
				return fmt.Errorf("el gestor %q no está soportado todavía", mgrName)
			}

			profile, content, err := resolveProfileAndContent(profile, contentDir)
			if err != nil {
				return err
			}
			m, err := loadManifestForProfile(content, profile)
			if err != nil {
				return err
			}
			target, err := addLayer(m, profile, layer)
			if err != nil {
				return err
			}

			if dryRun {
				fmt.Printf("dry-run: registraría en %s/packages/%s.txt e instalaría: %v\n", target, mgrName, pkgs)
				return nil
			}

			// 1) Registrar en la lista de la capa (declarativo).
			added, err := packages.AddToList(content, target, mgrName, pkgs)
			if err != nil {
				return err
			}
			if len(added) > 0 {
				fmt.Printf("registrados en %s/packages/%s.txt: %v\n", target, mgrName, added)
			} else {
				fmt.Println("ya estaban todos registrados; no toco la lista.")
			}

			// 2) Instalar (dnf es idempotente: los que ya estén no se tocan).
			fmt.Printf("instalando: %v\n", pkgs)
			return mgr.Install(pkgs)
		},
	}

	cmd.Flags().StringVar(&layer, "layer", "", "capa del perfil activo donde registrar (por defecto: machines/<máquina>)")
	cmd.Flags().StringVar(&contentDir, "content", "", "ruta al repo de contenido (por defecto: la del estado local, o el directorio actual)")
	cmd.Flags().StringVar(&profile, "profile", "", "perfil activo (por defecto: el del estado local)")
	return cmd
}

// addLayer decide en qué capa registra add: sin --layer, la de la máquina del
// perfil; con --layer, la pedida, siempre que esté en la pila del perfil. Si no
// está, el error enseña las capas que sí valen.
func addLayer(m *manifest.Manifest, profile, layer string) (string, error) {
	if layer == "" {
		return manifest.MachinesDir + "/" + m.Profiles[profile].Machine, nil
	}
	layer = strings.Trim(layer, "/")
	stack := m.Stack(profile)
	if !slices.Contains(stack, layer) {
		return "", fmt.Errorf("la capa %q no es del perfil %q, así que aquí no se aplicaría nunca (puedes usar: %s)",
			layer, profile, strings.Join(stack, ", "))
	}
	return layer, nil
}

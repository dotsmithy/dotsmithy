package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Points son los puntos de ciclo de vida válidos, en el orden en que se disparan
// durante un init (ADR 0007). update reutiliza post-link y post-packages.
var Points = []string{"post-link", "post-packages", "post-init"}

// Context es lo que se pasa a cada hook (y al post de los externals) por
// variables de entorno.
type Context struct {
	Profile    string            // perfil activo (DOTS_PROFILE)
	Machine    string            // máquina del perfil (DOTS_MACHINE)
	Params     map[string]string // valor del perfil en cada dimensión (DOTS_PARAM_<DIMENSIÓN>)
	ContentDir string            // raíz del contenido (DOTS_CONTENT_DIR y cwd de los hooks)
	Target     string            // destino del stow ya expandido, ej. $HOME (DOTS_TARGET)
	Layers     []string          // pila de capas resuelta (DOTS_LAYERS, separadas por espacio)
	HooksDir   string            // subdirectorio de hooks dentro de cada capa (ej. "hooks")
}

// script es un hook encontrado en alguna capa de la pila.
type script struct {
	name  string // nombre del fichero, que fija el orden (NN-nombre.sh)
	layer int    // posición de la capa en la pila, para desempatar
	rel   string // ruta relativa al contenido, para los mensajes
	path  string // ruta absoluta, para ejecutarlo
}

// Run ejecuta los hooks de un punto. Los busca en <capa>/<HooksDir>/<punto>/*.sh
// de cada capa de la pila, los junta todos y los ejecuta por orden léxico del
// nombre del fichero (convención NN-nombre.sh), así que el número decide el
// orden aunque vengan de capas distintas. Si dos se llaman igual, va antes el de
// la capa anterior en la pila. Se ejecutan con bash y cwd = raíz del contenido.
// Fail-fast: si uno sale con código != 0, se aborta nombrándolo. Si no hay
// hooks para el punto, no hace nada.
//
// En dry-run NO se ejecuta nada: solo se informa de qué correría. Es la opción
// segura (el motor no puede garantizar que un hook respete el modo). log puede
// ser nil.
func Run(point string, ctx Context, dryRun bool, log func(string)) error {
	if log == nil {
		log = func(string) {}
	}

	scripts, err := collect(point, ctx)
	if err != nil {
		return err
	}

	for _, s := range scripts {
		if dryRun {
			log(fmt.Sprintf("· (dry-run) ejecutaría %s", s.rel))
			continue
		}
		log(fmt.Sprintf("· ejecutando %s", s.rel))
		cmd := exec.Command("bash", s.path)
		cmd.Dir = ctx.ContentDir
		cmd.Env = append(os.Environ(), ctx.Env(point)...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("el hook %s falló: %w", s.rel, err) // fail-fast
		}
	}
	return nil
}

// collect junta los scripts del punto de todas las capas de la pila, ya
// ordenados por nombre y, a igual nombre, por posición de la capa.
func collect(point string, ctx Context) ([]script, error) {
	var out []script
	for i, layer := range ctx.Layers {
		relDir := filepath.Join(filepath.FromSlash(layer), ctx.HooksDir, point)
		dir := filepath.Join(ctx.ContentDir, relDir)
		entries, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue // esta capa no tiene hooks para este punto: normal
			}
			return nil, fmt.Errorf("no he podido leer los hooks de %s en la capa %s: %w", point, layer, err)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".sh") {
				continue
			}
			out = append(out, script{
				name:  e.Name(),
				layer: i,
				rel:   filepath.Join(relDir, e.Name()),
				path:  filepath.Join(dir, e.Name()),
			})
		}
	}
	sort.SliceStable(out, func(a, b int) bool {
		if out[a].name != out[b].name {
			return out[a].name < out[b].name
		}
		return out[a].layer < out[b].layer
	})
	return out, nil
}

// Env arma el contrato de variables de entorno que reciben los hooks (y el post
// de los externals, con point = "externals"). Los params van como
// DOTS_PARAM_<DIMENSIÓN> en mayúsculas, ordenados para que sea determinista.
func (ctx Context) Env(point string) []string {
	env := []string{
		"DOTS_HOOK=" + point,
		"DOTS_PROFILE=" + ctx.Profile,
		"DOTS_MACHINE=" + ctx.Machine,
		"DOTS_CONTENT_DIR=" + ctx.ContentDir,
		"DOTS_TARGET=" + ctx.Target,
		"DOTS_LAYERS=" + strings.Join(ctx.Layers, " "),
	}
	dims := make([]string, 0, len(ctx.Params))
	for dim := range ctx.Params {
		dims = append(dims, dim)
	}
	sort.Strings(dims)
	for _, dim := range dims {
		env = append(env, ParamVar(dim)+"="+ctx.Params[dim])
	}
	return env
}

// ParamVar devuelve el nombre de la variable de entorno de una dimensión:
// DOTS_PARAM_ + la dimensión en mayúsculas, con cualquier carácter que no valga
// en un nombre de variable cambiado por "_" (ej. "ux-mode" -> DOTS_PARAM_UX_MODE).
func ParamVar(dim string) string {
	var b strings.Builder
	b.WriteString("DOTS_PARAM_")
	for _, r := range strings.ToUpper(dim) {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}

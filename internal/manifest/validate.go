package manifest

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// modeRe valida un modo octal de 3 o 4 dígitos (ej. "644", "0440").
var modeRe = regexp.MustCompile(`^[0-7]{3,4}$`)

// reservedNames son nombres de carpeta que usa el motor en la raíz del
// contenido, así que no pueden ser el nombre de una dimensión.
var reservedNames = []string{SharedLayer, MachinesDir, HooksDir}

// problems acumula las pegas de una validación para devolverlas todas juntas.
type problems []string

func (p *problems) add(format string, args ...any) {
	*p = append(*p, fmt.Sprintf(format, args...))
}

// err devuelve las pegas como un único error, o nil si no hay ninguna.
func (p problems) err(what string) error {
	if len(p) == 0 {
		return nil
	}
	return fmt.Errorf("%s tiene %d problema(s):\n  - %s",
		what, len(p), strings.Join(p, "\n  - "))
}

// Validate comprueba que el manifiesto es coherente en sí mismo, sin mirar el
// disco (eso lo hace ValidateLayout). Acumula TODOS los problemas y los devuelve
// juntos, para que la usuaria los vea de una sola pasada en vez de arreglarlos de
// uno en uno. Devuelve nil si todo está bien.
func (m *Manifest) Validate() error {
	var p problems

	// --- version ---
	if m.Version != SchemaVersion {
		p.add("version: es %d, pero este motor solo entiende la versión %d", m.Version, SchemaVersion)
	}

	// --- dimensions ---
	seenDims := make(map[string]bool)
	for i, dim := range m.Dimensions {
		if msg := badSegment(dim); msg != "" {
			p.add("dimensions[%d]: %s", i, msg)
			continue
		}
		if contains(reservedNames, dim) {
			p.add("dimensions[%d]: %q es un nombre reservado del motor (reservados: %s)", i, dim, strings.Join(reservedNames, ", "))
		}
		if seenDims[dim] {
			p.add("dimensions[%d]: %q está repetida", i, dim)
		}
		seenDims[dim] = true
	}

	// --- profiles ---
	if len(m.Profiles) == 0 {
		p.add("profiles: hace falta al menos un perfil")
	}
	for _, name := range sortedKeys(m.Profiles) {
		prof := m.Profiles[name]
		if strings.TrimSpace(name) == "" {
			p.add("profiles: hay un perfil con nombre vacío")
			continue
		}
		if prof.Machine == "" {
			p.add("profiles.%s.machine: falta (es obligatorio: la carpeta de la máquina en %s/)", name, MachinesDir)
		} else if msg := badSegment(prof.Machine); msg != "" {
			p.add("profiles.%s.machine: %s", name, msg)
		}
		// Los params cubren exactamente las dimensiones: ni una más ni una menos.
		for _, dim := range m.Dimensions {
			val, ok := prof.Params[dim]
			if !ok {
				p.add("profiles.%s.params: falta el valor de la dimensión %q", name, dim)
				continue
			}
			if msg := badSegment(val); msg != "" {
				p.add("profiles.%s.params.%s: %s", name, dim, msg)
			}
		}
		for _, key := range sortedKeys(prof.Params) {
			if !contains(m.Dimensions, key) {
				p.add("profiles.%s.params.%s: no es ninguna dimensión declarada (dimensions: [%s])", name, key, strings.Join(m.Dimensions, ", "))
			}
		}
	}

	// --- packages ---
	for i, mgr := range m.Packages.Managers {
		if !contains(KnownManagers, mgr) {
			p.add("packages.managers[%d]: gestor desconocido %q (conocidos: %s)", i, mgr, strings.Join(KnownManagers, ", "))
		}
	}

	// --- hooks ---
	for i, pt := range m.Hooks.Points {
		if !contains(KnownHookPoints, pt) {
			p.add("hooks.points[%d]: punto desconocido %q (conocidos: %s)", i, pt, strings.Join(KnownHookPoints, ", "))
		}
	}

	// --- layers ---
	resolvable := m.resolvableLayers()
	for _, key := range sortedKeys(m.Layers) {
		if !resolvable[key] {
			p.add("layers.%s: ningún perfil usa esta capa (¿una errata?). Capas posibles: %s", key, strings.Join(sortedKeys(resolvable), ", "))
		}
		m.validateLayer(&p, "layers."+key, m.Layers[key])
	}

	return p.err("el manifiesto")
}

// validateLayer comprueba lo que declara una capa. prefix es la ruta de la capa
// en el manifiesto, para que los mensajes digan dónde está la pega.
func (m *Manifest) validateLayer(p *problems, prefix string, l Layer) {
	// COPR solo tiene sentido si dnf está activo.
	if len(l.COPR) > 0 && !m.HasManager("dnf") {
		p.add("%s.copr: hay repos COPR declarados pero dnf no está en packages.managers", prefix)
	}

	for i, s := range l.System {
		if strings.TrimSpace(s.Src) == "" {
			p.add("%s.system[%d]: falta src", prefix, i)
		}
		if strings.TrimSpace(s.Dest) == "" {
			p.add("%s.system[%d]: falta dest", prefix, i)
		}
		if !contains(knownSystemTypes, s.Type) {
			p.add("%s.system[%d]: type %q inválido (usa: %s)", prefix, i, s.Type, strings.Join(knownSystemTypes, ", "))
		}
		if s.Mode != "" && !modeRe.MatchString(s.Mode) {
			p.add("%s.system[%d]: mode %q no es un octal válido (ej. \"0440\")", prefix, i, s.Mode)
		}
	}

	for i, e := range l.Externals {
		if strings.TrimSpace(e.Dest) == "" {
			p.add("%s.externals[%d]: falta dest", prefix, i)
		}
		switch e.Type {
		case "git":
			if strings.TrimSpace(e.Repo) == "" {
				p.add("%s.externals[%d]: type git requiere repo", prefix, i)
			}
		case "file":
			if strings.TrimSpace(e.URL) == "" {
				p.add("%s.externals[%d]: type file requiere url", prefix, i)
			}
		default:
			p.add("%s.externals[%d]: type %q inválido (usa: %s)", prefix, i, e.Type, strings.Join(knownExternalTypes, ", "))
		}
	}
}

// resolvableLayers devuelve el conjunto de capas que algún perfil puede
// resolver: shared, cada <dimensión>/<valor> usado y cada machines/<machine>.
func (m *Manifest) resolvableLayers() map[string]bool {
	out := map[string]bool{SharedLayer: true}
	for name := range m.Profiles {
		for _, layer := range m.Stack(name) {
			out[layer] = true
		}
	}
	return out
}

// ValidateLayout comprueba contra el disco que existen las capas que resuelven
// los perfiles: machines/<machine> y cada <dimensión>/<valor> (aunque estén
// vacías). Así una errata tipo identity/wrok no pasa sin avisar. shared/ es
// opcional. Se llama después de Validate, con el manifiesto ya coherente.
func (m *Manifest) ValidateLayout(contentDir string) error {
	var p problems
	checked := make(map[string]bool)
	for _, name := range sortedKeys(m.Profiles) {
		// La primera capa de la pila es shared, que es opcional.
		for _, layer := range m.Stack(name)[1:] {
			if checked[layer] {
				continue
			}
			checked[layer] = true
			dir := filepath.Join(contentDir, filepath.FromSlash(layer))
			info, err := os.Stat(dir)
			switch {
			case err == nil && info.IsDir():
				// Bien.
			case err == nil:
				p.add("la capa %q del perfil %s no es una carpeta (%s)", layer, name, dir)
			case os.IsNotExist(err):
				p.add("falta la carpeta de la capa %q, que usa el perfil %s (créala, aunque sea vacía)", layer, name)
			default:
				p.add("no he podido mirar la capa %q: %v", layer, err)
			}
		}
	}
	return p.err("el contenido")
}

// badSegment comprueba que s vale como un nombre de carpeta suelto (una
// dimensión, un valor o una máquina). Devuelve por qué no vale, o "" si vale.
func badSegment(s string) string {
	switch {
	case strings.TrimSpace(s) == "":
		return "está vacío"
	case s != strings.TrimSpace(s):
		return fmt.Sprintf("%q tiene espacios al principio o al final", s)
	case strings.ContainsAny(s, `/\`):
		return fmt.Sprintf("%q no puede llevar barras (es un nombre de carpeta, no una ruta)", s)
	case s == "." || s == ".." || path.Clean(s) != s:
		return fmt.Sprintf("%q no es un nombre de carpeta válido", s)
	}
	return ""
}

// sortedKeys devuelve las claves de un mapa ordenadas, para que la validación
// y los mensajes sean deterministas.
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// contains indica si el slice contiene el valor dado.
func contains(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

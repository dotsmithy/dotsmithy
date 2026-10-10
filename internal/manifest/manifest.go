package manifest

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SchemaVersion es la única versión de esquema que este motor entiende. Si el
// manifiesto declara otra, se rechaza (evita aplicar un contenido pensado para
// otro motor). La 2 introduce los perfiles compuestos (capa común, dimensiones
// del contenido y máquina); la 1 ya no se acepta.
const SchemaVersion = 2

// Capas que pone el motor en las dos puntas de la pila. Sus nombres de carpeta
// son fijos: cualquier contenido tiene algo común y algo propio de cada máquina.
const (
	// SharedLayer es la capa común, la primera de la pila. Es opcional en disco.
	SharedLayer = "shared"
	// MachinesDir es la carpeta que aloja una capa por máquina, la última de la
	// pila: machines/<machine>.
	MachinesDir = "machines"
	// HooksDir es la carpeta de hooks dentro de cada capa: <capa>/hooks/<punto>/.
	HooksDir = "hooks"
)

// Manifest es el modelo tipado de dots.yaml.
type Manifest struct {
	Version int `yaml:"version"`
	// Dimensions son los ejes en los que varía un perfil, declarados por el
	// contenido (ej. [ux, identity]). El orden de la lista es el orden de
	// apilado. El motor no sabe qué significa cada una: solo que cada perfil da
	// un valor y que la capa vive en <dimensión>/<valor>.
	Dimensions []string           `yaml:"dimensions"`
	Profiles   map[string]Profile `yaml:"profiles"`
	Stow       Stow               `yaml:"stow"`
	Packages   Packages           `yaml:"packages"`
	Hooks      Hooks              `yaml:"hooks"`
	// Layers lleva lo que cada capa declara en el manifiesto (COPR, externals,
	// system), indexado por la ruta de la capa (ej. "shared", "ux/desktop",
	// "machines/minipc"). Una capa sin nada que declarar no aparece.
	Layers map[string]Layer `yaml:"layers"`
}

// Profile describe un perfil: el contenedor que junta las piezas con las que se
// compone una máquina. Se elige explícitamente en `dots init`.
type Profile struct {
	Description string `yaml:"description"`
	// Machine es la carpeta de la máquina dentro de machines/. Obligatorio. Dos
	// perfiles pueden compartir máquina.
	Machine string `yaml:"machine"`
	// Params son los valores del perfil en cada dimensión. Son del contenido y
	// van agrupados para no mezclarse con los campos del motor.
	Params map[string]string `yaml:"params"`
}

// Stow describe qué se despliega por symlinks y dónde. Las capas ya no se
// declaran aquí: las resuelve el motor a partir del perfil (ver Stack).
type Stow struct {
	// Target es el destino del stow (normalmente "~", que es $HOME).
	Target string `yaml:"target"`
	// Subdir es el subdirectorio dentro de cada capa que se stowea (ej. "home").
	Subdir string `yaml:"subdir"`
}

// Packages agrupa la configuración global de paquetes. Las listas viven en
// <capa>/packages/<gestor>.txt y los COPR en la sección layers.
type Packages struct {
	// Managers son los gestores activos (ej. dnf). El motor une
	// <capa>/packages/<gestor>.txt de cada capa activa.
	Managers []string `yaml:"managers"`
}

// Layer es lo que una capa declara en el manifiesto. Los ficheros de la capa
// (home/, packages/, hooks/ y los orígenes de system) viven en su carpeta; aquí
// va lo que no es un fichero.
type Layer struct {
	// COPR son los repos COPR (dnf) a habilitar antes de instalar.
	COPR []string `yaml:"copr"`
	// Externals son los recursos externos que trae esta capa.
	Externals []External `yaml:"externals"`
	// System son los ficheros de sistema que coloca esta capa.
	System []SystemFile `yaml:"system"`
}

// SystemFile describe un fichero de sistema (fuera de $HOME, requiere root).
type SystemFile struct {
	Src      string `yaml:"src"`      // ruta relativa a la carpeta de la capa (system/...)
	Dest     string `yaml:"dest"`     // ruta absoluta de destino (/etc/...)
	Type     string `yaml:"type"`     // "symlink" | "copy"
	Mode     string `yaml:"mode"`     // opcional, octal (ej. "0440")
	Owner    string `yaml:"owner"`    // opcional, "user:group" (default root:root)
	Validate string `yaml:"validate"` // opcional, comando de validación sobre copia temporal
}

// Hooks describe los hooks de ciclo de vida.
type Hooks struct {
	// Points son los puntos de ciclo de vida activos (subconjunto de KnownHookPoints).
	Points []string `yaml:"points"`
}

// External describe un recurso externo que el motor mantiene (git o fichero).
type External struct {
	Dest string `yaml:"dest"` // destino directo (NO se stowea)
	Type string `yaml:"type"` // "git" | "file"
	Repo string `yaml:"repo"` // requerido si type == git
	URL  string `yaml:"url"`  // requerido si type == file
	Ref  string `yaml:"ref"`  // rama/tag/commit para type == git
	Post string `yaml:"post"` // script post a ejecutar (relativo a la carpeta de la capa)
}

// Conjuntos de valores conocidos, usados por la validación.
var (
	// KnownManagers son los gestores de paquetes que el esquema admite. Solo dnf
	// está implementado en v1, pero el esquema reconoce los demás para no romper
	// manifiestos que los declaren de cara a futuras fases.
	KnownManagers = []string{"dnf", "pip", "npm", "cargo"}
	// KnownHookPoints son los puntos de ciclo de vida válidos.
	KnownHookPoints = []string{"post-link", "post-packages", "post-init"}
	// knownSystemTypes son los tipos válidos de fichero de system.
	knownSystemTypes = []string{"symlink", "copy"}
	// knownExternalTypes son los tipos válidos de external.
	knownExternalTypes = []string{"git", "file"}
)

// Load lee, parsea, aplica defaults y valida un dots.yaml desde disco. Es el
// punto de entrada normal del motor. Devuelve el manifiesto ya listo para usar,
// o un error explicando qué falla (agregado, todas las pegas de una vez).
func Load(manifestPath string) (*Manifest, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("no he podido leer el manifiesto %q: %w", manifestPath, err)
	}
	m, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("en %q: %w", manifestPath, err)
	}
	// El contenido es la carpeta donde vive dots.yaml: ahí deben existir las
	// capas que resuelven los perfiles.
	if err := m.ValidateLayout(filepath.Dir(manifestPath)); err != nil {
		return nil, fmt.Errorf("en %q: %w", manifestPath, err)
	}
	return m, nil
}

// Parse hace lo mismo que Load pero sobre bytes en memoria (sin tocar disco).
// Separado para poder testear el parseo y la validación sin ficheros.
func Parse(data []byte) (*Manifest, error) {
	var m Manifest
	dec := yaml.NewDecoder(strings.NewReader(string(data)))
	// KnownFields hace que un campo desconocido (ej. una errata) sea un error,
	// en vez de ignorarse en silencio.
	dec.KnownFields(true)
	if err := dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("el YAML no es válido: %w", err)
	}
	m.applyDefaults()
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return &m, nil
}

// applyDefaults rellena los valores por defecto de los campos opcionales, para
// que el resto del motor no tenga que repetir esa lógica.
func (m *Manifest) applyDefaults() {
	if m.Stow.Target == "" {
		m.Stow.Target = "~"
	}
	if m.Stow.Subdir == "" {
		m.Stow.Subdir = "home"
	}
	// Si no se declaran puntos de hook, se asumen todos los conocidos.
	if len(m.Hooks.Points) == 0 {
		m.Hooks.Points = append([]string(nil), KnownHookPoints...)
	}
}

// Stack devuelve la pila de capas del perfil, en orden de aplicación:
//
//	shared  →  <dimensión>/<valor> (por cada dimensión, en orden)  →  machines/<machine>
//
// Es una operación pura derivada del manifiesto: no mira el disco (las capas que
// no existan simplemente no aportan nada) ni comprueba que el perfil exista, que
// es responsabilidad de quien llama.
func (m *Manifest) Stack(profile string) []string {
	p := m.Profiles[profile]
	out := make([]string, 0, len(m.Dimensions)+2)
	out = append(out, SharedLayer)
	for _, dim := range m.Dimensions {
		out = append(out, path.Join(dim, p.Params[dim]))
	}
	return append(out, path.Join(MachinesDir, p.Machine))
}

// COPR devuelve la unión (sin duplicados, en orden de capa) de los repos COPR
// que declaran las capas de la pila.
func (m *Manifest) COPR(stack []string) []string {
	var out []string
	seen := make(map[string]bool)
	for _, layer := range stack {
		for _, repo := range m.Layers[layer].COPR {
			if !seen[repo] {
				seen[repo] = true
				out = append(out, repo)
			}
		}
	}
	return out
}

// Externals devuelve los externals de las capas de la pila, en orden de capa.
// El post de cada uno se devuelve ya relativo a la raíz del contenido (la capa
// lo declara relativo a su carpeta).
func (m *Manifest) Externals(stack []string) []External {
	var out []External
	for _, layer := range stack {
		for _, e := range m.Layers[layer].Externals {
			if e.Post != "" {
				e.Post = path.Join(layer, e.Post)
			}
			out = append(out, e)
		}
	}
	return out
}

// SystemFiles devuelve los ficheros de sistema de las capas de la pila, en orden
// de capa. El src de cada uno se devuelve ya relativo a la raíz del contenido
// (la capa lo declara relativo a su carpeta).
func (m *Manifest) SystemFiles(stack []string) []SystemFile {
	var out []SystemFile
	for _, layer := range stack {
		for _, f := range m.Layers[layer].System {
			f.Src = path.Join(layer, f.Src)
			out = append(out, f)
		}
	}
	return out
}

// HasManager indica si un gestor de paquetes está activo en el manifiesto.
func (m *Manifest) HasManager(name string) bool {
	for _, mgr := range m.Packages.Managers {
		if mgr == name {
			return true
		}
	}
	return false
}

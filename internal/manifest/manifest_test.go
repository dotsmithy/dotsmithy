package manifest

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// fullValid es un manifiesto válido que ejercita todas las secciones.
const fullValid = `
version: 2
dimensions: [ux, identity]
profiles:
  desk:
    description: "Escritorio de casa"
    machine: desk
    params: { ux: desktop, identity: personal }
  laptop:
    description: "Portátil de trabajo"
    machine: laptop
    params: { ux: desktop, identity: work }
  devbox:
    description: "Máquina remota"
    machine: devbox
    params: { ux: terminal, identity: work }
stow:
  target: ~
  subdir: home
packages:
  managers: [dnf, cargo]
hooks:
  points: [post-link, post-init]
layers:
  shared:
    copr: [owner/common]
    externals:
      - dest: ~/.config/nvim
        type: git
        repo: git@github.com:example/nvim-config.git
        ref: main
        post: hooks/externals/nvim-sync.sh
  ux/desktop:
    copr: [owner/desktop-shell, owner/common]
    externals:
      - dest: ~/.zim/zimfw.zsh
        type: file
        url: https://example.com/zimfw.zsh
  identity/work:
    system:
      - src: system/sudoers/me
        dest: /etc/sudoers.d/me
        type: copy
        mode: "0440"
        validate: "visudo -cf {file}"
  machines/laptop:
    copr: [owner/laptop-only]
`

// mustParse parsea un manifiesto que el test da por válido.
func mustParse(t *testing.T, yaml string) *Manifest {
	t.Helper()
	m, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("no esperaba error, tengo: %v", err)
	}
	return m
}

func TestParseValidFull(t *testing.T) {
	m := mustParse(t, fullValid)
	if m.Version != 2 {
		t.Errorf("version: quiero 2, tengo %d", m.Version)
	}
	if !reflect.DeepEqual(m.Dimensions, []string{"ux", "identity"}) {
		t.Errorf("dimensions: tengo %v", m.Dimensions)
	}
	if len(m.Profiles) != 3 {
		t.Errorf("profiles: quiero 3, tengo %d", len(m.Profiles))
	}
	desk := m.Profiles["desk"]
	if desk.Description != "Escritorio de casa" || desk.Machine != "desk" || desk.Params["identity"] != "personal" {
		t.Errorf("perfil desk mal parseado: %+v", desk)
	}
	if len(m.Layers) != 4 {
		t.Errorf("layers: quiero 4, tengo %d", len(m.Layers))
	}
}

func TestDefaultsApplied(t *testing.T) {
	// Manifiesto mínimo: sin dimensiones y con un único perfil.
	const minimal = `
version: 2
profiles:
  desk: { machine: desk }
`
	m := mustParse(t, minimal)
	if m.Stow.Target != "~" {
		t.Errorf("default target: quiero \"~\", tengo %q", m.Stow.Target)
	}
	if m.Stow.Subdir != "home" {
		t.Errorf("default subdir: quiero \"home\", tengo %q", m.Stow.Subdir)
	}
	if len(m.Hooks.Points) != len(KnownHookPoints) {
		t.Errorf("default hooks.points: quiero %v, tengo %v", KnownHookPoints, m.Hooks.Points)
	}
}

func TestStack(t *testing.T) {
	m := mustParse(t, fullValid)
	cases := map[string][]string{
		"desk":   {"shared", "ux/desktop", "identity/personal", "machines/desk"},
		"laptop": {"shared", "ux/desktop", "identity/work", "machines/laptop"},
		"devbox": {"shared", "ux/terminal", "identity/work", "machines/devbox"},
	}
	for profile, want := range cases {
		if got := m.Stack(profile); !reflect.DeepEqual(got, want) {
			t.Errorf("pila de %s: quiero %v, tengo %v", profile, want, got)
		}
	}
}

func TestStackWithoutDimensions(t *testing.T) {
	m := mustParse(t, "version: 2\nprofiles:\n  desk: { machine: desk }\n")
	want := []string{"shared", "machines/desk"}
	if got := m.Stack("desk"); !reflect.DeepEqual(got, want) {
		t.Errorf("pila sin dimensiones: quiero %v, tengo %v", want, got)
	}
}

func TestCOPRUnionByStack(t *testing.T) {
	m := mustParse(t, fullValid)
	cases := map[string][]string{
		// owner/common aparece en shared y en ux/desktop: sale una vez.
		"desk":   {"owner/common", "owner/desktop-shell"},
		"laptop": {"owner/common", "owner/desktop-shell", "owner/laptop-only"},
		"devbox": {"owner/common"},
	}
	for profile, want := range cases {
		if got := m.COPR(m.Stack(profile)); !reflect.DeepEqual(got, want) {
			t.Errorf("COPR de %s: quiero %v, tengo %v", profile, want, got)
		}
	}
}

func TestExternalsByStackResolvePost(t *testing.T) {
	m := mustParse(t, fullValid)

	desk := m.Externals(m.Stack("desk"))
	if len(desk) != 2 {
		t.Fatalf("externals de desk: quiero 2, tengo %d", len(desk))
	}
	// El post lo declara la capa relativo a su carpeta; sale relativo al contenido.
	if desk[0].Post != "shared/hooks/externals/nvim-sync.sh" {
		t.Errorf("post resuelto: tengo %q", desk[0].Post)
	}
	if desk[1].Post != "" {
		t.Errorf("un external sin post debe seguir sin post, tengo %q", desk[1].Post)
	}
	// La resolución no debe tocar el manifiesto original.
	if got := m.Layers["shared"].Externals[0].Post; got != "hooks/externals/nvim-sync.sh" {
		t.Errorf("el manifiesto original no debería cambiar, tengo %q", got)
	}

	if got := len(m.Externals(m.Stack("devbox"))); got != 1 {
		t.Errorf("externals de devbox: quiero 1 (solo shared), tengo %d", got)
	}
}

func TestSystemFilesByStackResolveSrc(t *testing.T) {
	m := mustParse(t, fullValid)
	if got := m.SystemFiles(m.Stack("desk")); len(got) != 0 {
		t.Errorf("desk no usa identity/work: quiero 0 ficheros, tengo %v", got)
	}
	files := m.SystemFiles(m.Stack("devbox"))
	if len(files) != 1 || files[0].Src != "identity/work/system/sudoers/me" {
		t.Errorf("system de devbox: tengo %+v", files)
	}
}

func TestHasManager(t *testing.T) {
	m := &Manifest{Packages: Packages{Managers: []string{"dnf", "npm"}}}
	if !m.HasManager("dnf") {
		t.Error("debería tener dnf")
	}
	if m.HasManager("cargo") {
		t.Error("no debería tener cargo")
	}
}

func TestParseUnknownFieldFails(t *testing.T) {
	// Una errata en un campo (typo) debe fallar, no ignorarse. Incluye los
	// campos de la versión 1 que ya no existen.
	for name, yaml := range map[string]string{
		"errata":          "version: 2\nprofiles:\n  d: { machine: d }\nstow:\n  targett: ~\n",
		"stow.layers":     "version: 2\nprofiles:\n  d: { machine: d }\nstow:\n  layers: [shared]\n",
		"externals raíz":  "version: 2\nprofiles:\n  d: { machine: d }\nexternals: []\n",
		"campo en perfil": "version: 2\nprofiles:\n  d: { machine: d, ux: desktop }\n",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := Parse([]byte(yaml)); err == nil {
				t.Fatal("esperaba error por campo desconocido, no lo hubo")
			}
		})
	}
}

func TestValidationErrors(t *testing.T) {
	// base es un manifiesto válido al que cada caso le añade la pega.
	const base = "version: 2\ndimensions: [ux]\nprofiles:\n  d: { machine: d, params: { ux: desktop } }\n"
	cases := []struct {
		name string
		yaml string
		want string // fragmento que debe aparecer en el error
	}{
		{
			name: "versión 1 ya no se acepta",
			yaml: "version: 1\nprofiles:\n  d: { machine: d }\n",
			want: "solo entiende la versión 2",
		},
		{
			name: "sin perfiles",
			yaml: "version: 2\n",
			want: "al menos un perfil",
		},
		{
			name: "perfil sin machine",
			yaml: "version: 2\nprofiles:\n  d: { description: x }\n",
			want: "profiles.d.machine: falta",
		},
		{
			name: "machine con barra",
			yaml: "version: 2\nprofiles:\n  d: { machine: a/b }\n",
			want: "no puede llevar barras",
		},
		{
			name: "falta un param",
			yaml: "version: 2\ndimensions: [ux, identity]\nprofiles:\n  d: { machine: d, params: { ux: desktop } }\n",
			want: "falta el valor de la dimensión \"identity\"",
		},
		{
			name: "param de más",
			yaml: "version: 2\ndimensions: [ux]\nprofiles:\n  d: { machine: d, params: { ux: desktop, color: red } }\n",
			want: "params.color: no es ninguna dimensión",
		},
		{
			name: "valor vacío",
			yaml: "version: 2\ndimensions: [ux]\nprofiles:\n  d: { machine: d, params: { ux: \"\" } }\n",
			want: "params.ux: está vacío",
		},
		{
			name: "dimensión reservada",
			yaml: "version: 2\ndimensions: [machines]\nprofiles:\n  d: { machine: d, params: { machines: x } }\n",
			want: "nombre reservado",
		},
		{
			name: "dimensión repetida",
			yaml: "version: 2\ndimensions: [ux, ux]\nprofiles:\n  d: { machine: d, params: { ux: desktop } }\n",
			want: "está repetida",
		},
		{
			name: "gestor desconocido",
			yaml: base + "packages:\n  managers: [brew]\n",
			want: "gestor desconocido",
		},
		{
			name: "capa que ningún perfil usa",
			yaml: base + "layers:\n  ux/terminal:\n    copr: [a/b]\npackages:\n  managers: [dnf]\n",
			want: "layers.ux/terminal: ningún perfil usa esta capa",
		},
		{
			name: "copr sin dnf",
			yaml: base + "packages:\n  managers: [cargo]\nlayers:\n  shared:\n    copr: [atim/lazygit]\n",
			want: "layers.shared.copr: hay repos COPR declarados pero dnf no está",
		},
		{
			name: "system sin type válido",
			yaml: base + "layers:\n  machines/d:\n    system:\n      - { src: a, dest: /b, type: hardlink }\n",
			want: "layers.machines/d.system[0]: type \"hardlink\" inválido",
		},
		{
			name: "system con mode inválido",
			yaml: base + "layers:\n  shared:\n    system:\n      - { src: a, dest: /b, type: copy, mode: \"999x\" }\n",
			want: "octal válido",
		},
		{
			name: "hook point desconocido",
			yaml: base + "hooks:\n  points: [pre-link]\n",
			want: "punto desconocido",
		},
		{
			name: "external git sin repo",
			yaml: base + "layers:\n  ux/desktop:\n    externals:\n      - { dest: ~/x, type: git }\n",
			want: "layers.ux/desktop.externals[0]: type git requiere repo",
		},
		{
			name: "external file sin url",
			yaml: base + "layers:\n  shared:\n    externals:\n      - { dest: ~/x, type: file }\n",
			want: "type file requiere url",
		},
		{
			name: "external type inválido",
			yaml: base + "layers:\n  shared:\n    externals:\n      - { dest: ~/x, type: svn }\n",
			want: "inválido",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse([]byte(tc.yaml))
			if err == nil {
				t.Fatalf("esperaba error con %q, pero parseó bien", tc.want)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("el error no menciona %q; error: %v", tc.want, err)
			}
		})
	}
}

func TestValidationAggregatesProblems(t *testing.T) {
	// Manifiesto con varias pegas a la vez: deben salir todas juntas.
	const many = `
version: 3
dimensions: [ux]
profiles:
  d: { params: { ux: desktop } }
`
	_, err := Parse([]byte(many))
	if err == nil {
		t.Fatal("esperaba error")
	}
	for _, frag := range []string{"version", "machine: falta"} {
		if !strings.Contains(err.Error(), frag) {
			t.Errorf("el error agregado no menciona %q; error: %v", frag, err)
		}
	}
}

func TestLoadFromDisk(t *testing.T) {
	m, err := Load(filepath.Join("testdata", "content", "dots.yaml"))
	if err != nil {
		t.Fatalf("no esperaba error cargando el fixture: %v", err)
	}
	want := []string{"shared", "ux/terminal", "identity/work", "machines/devbox"}
	if got := m.Stack("devbox"); !reflect.DeepEqual(got, want) {
		t.Errorf("pila del fixture: quiero %v, tengo %v", want, got)
	}
}

func TestLoadMissingFile(t *testing.T) {
	if _, err := Load(filepath.Join("testdata", "no-existe.yaml")); err == nil {
		t.Fatal("esperaba error al cargar un fichero inexistente")
	}
}

func TestValidateLayout(t *testing.T) {
	const yaml = `
version: 2
dimensions: [identity]
profiles:
  a: { machine: one, params: { identity: personal } }
  b: { machine: two, params: { identity: wrok } }
`
	m := mustParse(t, yaml)
	dir := t.TempDir()
	// shared/ no existe a propósito: es opcional.
	for _, d := range []string{"machines/one", "identity/personal"} {
		if err := os.MkdirAll(filepath.Join(dir, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	// machines/two es un fichero, no una carpeta.
	if err := os.WriteFile(filepath.Join(dir, "machines", "two"), nil, 0o644); err != nil {
		t.Fatal(err)
	}

	err := m.ValidateLayout(dir)
	if err == nil {
		t.Fatal("esperaba error: falta identity/wrok y machines/two no es carpeta")
	}
	for _, frag := range []string{`falta la carpeta de la capa "identity/wrok"`, `"machines/two" del perfil b no es una carpeta`} {
		if !strings.Contains(err.Error(), frag) {
			t.Errorf("el error no menciona %q; error: %v", frag, err)
		}
	}
	if strings.Contains(err.Error(), "shared") {
		t.Errorf("shared es opcional y no debería salir; error: %v", err)
	}

	// Con todo en su sitio, pasa.
	if err := os.Remove(filepath.Join(dir, "machines", "two")); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"machines/two", "identity/wrok"} {
		if err := os.MkdirAll(filepath.Join(dir, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := m.ValidateLayout(dir); err != nil {
		t.Errorf("con todas las capas creadas no esperaba error: %v", err)
	}
}

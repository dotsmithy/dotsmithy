package cli

import (
	"strings"
	"testing"

	"go.kreikol.dev/dotsmithy/internal/manifest"
)

// addManifest es un manifiesto con una dimensión y dos perfiles, para ver qué
// capas acepta add en cada uno.
func addManifest() *manifest.Manifest {
	return &manifest.Manifest{
		Version:    manifest.SchemaVersion,
		Dimensions: []string{"ux"},
		Profiles: map[string]manifest.Profile{
			"remote": {Machine: "vm", Params: map[string]string{"ux": "terminal"}},
			"desk":   {Machine: "pc", Params: map[string]string{"ux": "desktop"}},
		},
	}
}

func TestAddLayerDefaultsToMachine(t *testing.T) {
	got, err := addLayer(addManifest(), "remote", "")
	if err != nil || got != "machines/vm" {
		t.Errorf("sin --layer: quiero machines/vm, tengo %q (err %v)", got, err)
	}
}

func TestAddLayerAcceptsProfileLayers(t *testing.T) {
	for _, layer := range []string{"shared", "ux/terminal", "machines/vm", "ux/terminal/"} {
		got, err := addLayer(addManifest(), "remote", layer)
		if err != nil {
			t.Errorf("--layer %s debería valer en remote: %v", layer, err)
			continue
		}
		if want := strings.Trim(layer, "/"); got != want {
			t.Errorf("--layer %s: quiero %q, tengo %q", layer, want, got)
		}
	}
}

func TestAddLayerRejectsForeignLayers(t *testing.T) {
	// Capas de otro perfil o que no existen: no se aplicarían en remote.
	for _, layer := range []string{"ux/desktop", "machines/pc", "identity/work"} {
		_, err := addLayer(addManifest(), "remote", layer)
		if err == nil {
			t.Errorf("--layer %s NO debería valer en remote", layer)
			continue
		}
		if !strings.Contains(err.Error(), "shared, ux/terminal, machines/vm") {
			t.Errorf("el error debería listar las capas válidas: %v", err)
		}
	}
}

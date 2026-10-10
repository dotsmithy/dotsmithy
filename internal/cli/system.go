package cli

import (
	"fmt"

	"go.kreikol.dev/dotsmithy/internal/manifest"
	"go.kreikol.dev/dotsmithy/internal/system"
)

// applySystem coloca los ficheros de sistema que declaran las capas del perfil
// (/etc, systemd, sudoers). Es la única parte que necesita root (sudo). Respeta
// --dry-run. Se aplica como parte de la fase link, antes del hook post-link.
func applySystem(m *manifest.Manifest, contentDir, profile string) error {
	files := m.SystemFiles(m.Stack(profile))
	if len(files) == 0 {
		return nil
	}
	fmt.Printf("· system: %d fichero(s) de sistema…\n", len(files))
	return system.Apply(files, contentDir, dryRun, func(s string) {
		fmt.Println("  " + s)
	})
}

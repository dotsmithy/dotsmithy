package cli

import (
	"fmt"

	"go.kreikol.dev/dotsmithy/internal/hooks"
	"go.kreikol.dev/dotsmithy/internal/manifest"
	"go.kreikol.dev/dotsmithy/internal/stow"
)

// applyHook ejecuta los hooks de un punto, si el manifiesto lo declara activo.
// Los busca en cada capa de la pila (ver hooks.Run) y respeta --dry-run. Lo usan init y update.
func applyHook(m *manifest.Manifest, contentDir, profile, point string) error {
	if !hookEnabled(m, point) {
		return nil
	}
	target, err := stow.ExpandTarget(m.Stow.Target)
	if err != nil {
		return err
	}
	return hooks.Run(point, hookContext(m, contentDir, profile, target), dryRun, func(s string) {
		fmt.Println("  " + s)
	})
}

// hookContext arma el contexto que reciben los hooks y el post de los externals
// (perfil, máquina, params, contenido, destino expandido y pila de capas).
func hookContext(m *manifest.Manifest, contentDir, profile, target string) hooks.Context {
	p := m.Profiles[profile]
	return hooks.Context{
		Profile:    profile,
		Machine:    p.Machine,
		Params:     p.Params,
		ContentDir: contentDir,
		Target:     target,
		Layers:     m.Stack(profile),
		HooksDir:   manifest.HooksDir,
	}
}

// hookEnabled indica si el manifiesto declara activo un punto de hook.
func hookEnabled(m *manifest.Manifest, point string) bool {
	for _, p := range m.Hooks.Points {
		if p == point {
			return true
		}
	}
	return false
}

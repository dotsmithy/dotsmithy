#!/usr/bin/env bash
# Ejemplo de uso de las variables que dotsmithy pasa a los hooks.
#
# Guarda el perfil, la máquina y el valor de cada dimensión (DOTS_PARAM_*) en
# ~/.cache/dotsmithy/profile.env, para que tu shell pueda saber dónde está sin
# preguntarle al motor. Lo usa la bienvenida de ~/.bashrc.d/10-dotsmithy-welcome.sh.
#
# Corre en post-link, que lanzan `dots init` y `dots update` (un `dots link`
# suelto no ejecuta hooks): si cambias de perfil, el fichero cambia con él.
set -euo pipefail

cache_dir="${XDG_CACHE_HOME:-$HOME/.cache}/dotsmithy"
mkdir -p "$cache_dir"

{
  printf 'DOTS_PROFILE=%q\n' "$DOTS_PROFILE"
  printf 'DOTS_MACHINE=%q\n' "$DOTS_MACHINE"
  # Una línea por dimensión: DOTS_PARAM_UX, DOTS_PARAM_IDENTITY...
  for var in $(compgen -v DOTS_PARAM_ || true); do
    printf '%s=%q\n' "$var" "${!var}"
  done
} >"$cache_dir/profile.env"

echo "dotsmithy: perfil guardado en $cache_dir/profile.env"

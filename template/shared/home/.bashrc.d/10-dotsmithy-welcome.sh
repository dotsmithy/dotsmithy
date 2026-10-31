# shellcheck shell=bash
# Bienvenida al abrir una terminal: qué versión de dotsmithy tienes y con qué
# perfil está montada esta máquina. Por ejemplo:
#   🔨 dotsmithy 1.0.0 · perfil laptop (máquina laptop, ux=desktop)
#
# El perfil sale de ~/.cache/dotsmithy/profile.env, que escribe el hook
# shared/hooks/post-link/10-profile-env.sh con las variables DOTS_* que reciben
# los hooks. Si ese fichero aún no existe, enseña solo la versión.

# Solo en shells interactivas: un script no necesita saludos.
[[ $- == *i* ]] || return 0

# Todo en una subshell, para no dejarte variables sueltas en la sesión.
(
  line="🔨 dotsmithy"
  if command -v dots >/dev/null 2>&1; then
    version=$(dots --version 2>/dev/null) || version=""
    [[ -n $version ]] && line+=" ${version#dots }"
  fi

  env_file="${XDG_CACHE_HOME:-$HOME/.cache}/dotsmithy/profile.env"
  if [[ -r $env_file ]]; then
    # shellcheck source=/dev/null
    source "$env_file"
    details="máquina ${DOTS_MACHINE:-?}"
    for var in $(compgen -v DOTS_PARAM_); do
      name=${var#DOTS_PARAM_}
      details+=", ${name,,}=${!var}"
    done
    line+=" · perfil ${DOTS_PROFILE:-?} ($details)"
  fi

  echo "$line"
)

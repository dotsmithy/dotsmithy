#!/usr/bin/env bash
#
# Integración de perfiles compuestos: se ejecuta DENTRO del contenedor Fedora.
# Con un contenido de dos perfiles que comparten shared pero difieren en la
# dimensión ux y en la máquina, comprueba que init aplica SOLO las capas del
# perfil elegido: dotfiles y hooks de shared + ux/terminal + machines/vm, y nada
# de ux/desktop ni de machines/pc. Comprueba también el orden de los hooks entre
# capas, el contrato DOTS_MACHINE/DOTS_PARAM_* y el aviso por hooks en la raíz.
set -euo pipefail

fail() {
	printf 'FALLO: %s\n' "$*" >&2
	exit 1
}

C=/tmp/profiles-content
mkdir -p "$C/shared/home" "$C/shared/hooks/post-link" \
	"$C/ux/terminal/hooks/post-link" \
	"$C/ux/desktop/home" "$C/ux/desktop/hooks/post-link" \
	"$C/machines/vm/home" "$C/machines/vm/hooks/post-link" \
	"$C/machines/pc" \
	"$C/hooks/post-link"
cat >"$C/dots.yaml" <<'EOF'
version: 1
dimensions: [ux]
profiles:
  remote: { description: "sin escritorio", machine: vm, params: { ux: terminal } }
  desk: { description: "con escritorio", machine: pc, params: { ux: desktop } }
hooks: { points: [post-link] }
EOF

# Dotfiles: uno por capa, para ver cuáles se enlazan.
echo common >"$C/shared/home/.dots-common"
echo desktop >"$C/ux/desktop/home/.dots-desktop"
echo vm >"$C/machines/vm/home/.dots-vm"

# Hooks: cada uno apunta su nombre en un log. Manda el número del nombre, no la
# capa: el de la máquina (05) va antes que el de shared (10).
cat >"$C/machines/vm/hooks/post-link/05-vm.sh" <<'EOF'
echo vm >> "$DOTS_TARGET/.hooks.log"
EOF
cat >"$C/shared/hooks/post-link/10-shared.sh" <<'EOF'
echo "shared $DOTS_MACHINE $DOTS_PARAM_UX" >> "$DOTS_TARGET/.hooks.log"
EOF
cat >"$C/ux/terminal/hooks/post-link/20-terminal.sh" <<'EOF'
echo terminal >> "$DOTS_TARGET/.hooks.log"
EOF
cat >"$C/ux/desktop/hooks/post-link/20-desktop.sh" <<'EOF'
echo desktop >> "$DOTS_TARGET/.hooks.log"
EOF
# Un hook en la raíz (sitio antiguo): no debe ejecutarse, pero sí avisarse.
cat >"$C/hooks/post-link/10-root.sh" <<'EOF'
echo root >> "$DOTS_TARGET/.hooks.log"
EOF

echo ">> init del perfil remote aplica solo sus capas"
err="$(dots init --profile remote --content "$C" 2>&1 >/dev/null)"

[ -L "$HOME/.dots-common" ] || fail "no se enlazó el dotfile de shared"
[ -L "$HOME/.dots-vm" ] || fail "no se enlazó el dotfile de la máquina"
if [ -e "$HOME/.dots-desktop" ]; then fail "NO debería enlazarse el dotfile de ux/desktop"; fi

echo ">> los hooks de las capas del perfil corren en orden de nombre"
want="$(printf 'vm\nshared vm terminal\nterminal')"
got="$(cat "$HOME/.hooks.log")"
[ "$got" = "$want" ] || fail "orden/contenido de hooks inesperado:
--- quiero
$want
--- tengo
$got"

echo ">> los hooks de la raíz no corren, pero se avisa"
grep -q "ahora los hooks van por capa" <<<"$err" || fail "esperaba el aviso por hooks en la raíz; stderr: $err"

echo "OK: el perfil compuesto aplica solo sus capas, con hooks por capa en orden y su contrato de entorno."

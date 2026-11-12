#!/usr/bin/env bash
#
# Levanta el lab: una Fedora con systemd y dotsmithy aplicado, para trastear a mano.
#
#   1) compila el motor (linux/amd64, estático) en tests/lab/bin/dots;
#   2) construye la imagen (tests/lab/Containerfile);
#   3) arranca el contenedor con systemd como PID 1 y espera a que esté listo;
#   4) aplica el contenido de muestra como el usuario "lab" (dots init).
#
# Con LAB_DOTS_VERSION (p.ej. v1.0.0-rc.2) no compila nada: la imagen sale sin
# dots y, ya con el lab arrancado, instala esa release publicada con el
# instalador (script/install.sh) como el usuario "lab", igual que en una máquina
# real. Sirve para probar una versión candidata antes de publicarla. Se usa el
# instalador del propio checkout (el que saldrá con esa versión): el de la URL raw
# de main es el de la última versión publicada, que puede ir por detrás.
#
# Después, según LAB_REMOTE:
#   - LAB_REMOTE=1 (lo usa el workflow Lab): abre una sesión SSH remota con upterm
#     a la que solo pueden entrar las claves públicas en GitHub de LAB_GITHUB_USER
#     (obligatorio), y espera hasta que esa sesión termine (saliendo de la shell)
#     o alguien haga `touch /tmp/close-lab` dentro del lab; al salir, borra el
#     contenedor.
#   - sin definir (uso local): deja el contenedor corriendo y te dice cómo entrar
#     y cómo tirarlo.
#
# Necesita podman en el host: el arranque con systemd se apoya en
# `podman run --systemd=always`, que en docker no existe tal cual. Por defecto va
# rootless; con LAB_PODMAN se cambia el comando (el workflow usa "sudo podman",
# porque el podman rootless de los runners de Ubuntu no deja leer /etc/shadow a
# sudo dentro del contenedor).
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$HERE/../.." && pwd)"
readonly IMAGE="dotsmithy-lab:latest"
readonly NAME="${LAB_NAME:-dotsmithy-lab}"
# Release publicada a instalar en vez de compilar (vacío = compilar).
readonly DOTS_VERSION_WANTED="${LAB_DOTS_VERSION:-}"

# Comando de podman a usar, troceado en palabras (p.ej. "sudo podman").
read -r -a PODMAN <<<"${LAB_PODMAN:-podman}"

if ! command -v podman >/dev/null 2>&1; then
	echo "Necesito podman en el sistema para el lab (arranca Fedora con systemd) y no lo encuentro." >&2
	exit 1
fi

# En modo remoto, sin una usuaria de GitHub a quien autorizar, cualquiera con el
# enlace podría entrar. Eso no lo hacemos, y lo miramos antes de montar nada.
if [ "${LAB_REMOTE:-}" = "1" ] && [ -z "${LAB_GITHUB_USER:-}" ]; then
	echo "Para abrir el acceso remoto necesito LAB_GITHUB_USER (la cuenta de GitHub cuyas claves SSH pueden entrar)." >&2
	exit 1
fi

# lab_exec ejecuta un comando dentro del lab como el usuario "lab", con shell de
# login (para que tenga el PATH y el entorno de una sesión normal).
lab_exec() {
	"${PODMAN[@]}" exec -u lab -w /home/lab "$NAME" bash -lc "$1"
}

# El binario va a bin/, que la imagen copia tal cual: con release, vacío.
rm -rf "${HERE:?}/bin"
mkdir -p "$HERE/bin"
if [ -z "$DOTS_VERSION_WANTED" ]; then
	echo ">> compilando el motor (linux/amd64, estático)…"
	(
		cd "$ROOT"
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
			go build -trimpath -ldflags "-s -w" -o "$HERE/bin/dots" ./cmd/dots
	)
else
	echo ">> nada que compilar: el lab va con la release $DOTS_VERSION_WANTED"
fi

echo ">> construyendo la imagen del lab…"
# La salida del build es muy ruidosa (scriptlets de dnf): solo la enseñamos si falla.
build_log="$(mktemp)"
if ! "${PODMAN[@]}" build -t "$IMAGE" "$HERE" >"$build_log" 2>&1; then
	tail -40 "$build_log" >&2
	echo "no he podido construir la imagen del lab (arriba, el final del build)." >&2
	exit 1
fi
rm -f "$build_log"

echo ">> arrancando Fedora con systemd…"
"${PODMAN[@]}" rm -f "$NAME" >/dev/null 2>&1 || true
"${PODMAN[@]}" run -d --name "$NAME" --systemd=always "$IMAGE" >/dev/null

# Esperamos a que systemd termine de arrancar. Recién lanzado el contenedor puede
# que aún no responda (bus sin levantar), así que reintentamos un rato en vez de
# rendirnos a la primera. En un contenedor lo normal es "degraded" (hay unidades
# que no pintan nada ahí y fallan): eso también vale.
state=""
for _ in $(seq 1 45); do
	state="$("${PODMAN[@]}" exec "$NAME" systemctl is-system-running 2>&1 || true)"
	case "$state" in running | degraded) break ;; esac
	sleep 2
done
case "$state" in
running | degraded) echo ">> systemd listo (estado: $state)" ;;
*)
	# Volcamos todo lo que ayude a ver qué ha pasado.
	echo "systemd no ha arrancado bien dentro del lab (última respuesta: '${state:-nada}')." >&2
	echo "--- estado del contenedor:" >&2
	"${PODMAN[@]}" ps -a --filter "name=$NAME" --format '{{.Status}}' >&2 || true
	echo "--- PID 1 dentro:" >&2
	"${PODMAN[@]}" exec "$NAME" ps -p 1 -o comm= >&2 || true
	echo "--- logs del contenedor:" >&2
	"${PODMAN[@]}" logs "$NAME" 2>&1 | tail -30 >&2 || true
	exit 1
	;;
esac

if [ -n "$DOTS_VERSION_WANTED" ]; then
	echo ">> instalando la release $DOTS_VERSION_WANTED con el instalador, como el usuario lab…"
	# Igual que el «curl … | bash» de una máquina real, pero con el instalador del
	# checkout: se le pasa por la entrada estándar.
	if ! "${PODMAN[@]}" exec -i -u lab -w /home/lab -e DOTS_VERSION="$DOTS_VERSION_WANTED" \
		"$NAME" bash -s <"$ROOT/script/install.sh"; then
		echo "no he podido instalar la release $DOTS_VERSION_WANTED (arriba, lo que dijo el instalador)." >&2
		exit 1
	fi
fi

# Que quede claro en el log qué dots se está probando.
echo ">> versión de dots en el lab: $(lab_exec "dots --version" 2>&1)"

echo ">> aplicando el contenido de muestra como el usuario lab (dots init)…"
lab_exec "dots init --profile minipc --content ~/content" || echo "(init terminó con avisos)"

if [ "${LAB_REMOTE:-}" != "1" ]; then
	cat <<MSG
=================== LAB LISTO ===================
Fedora con systemd y dots aplicado, usuario "lab" (sudo sin contraseña).
  Entrar:  ${PODMAN[*]} exec -it -u lab -w /home/lab $NAME bash -l
  Tirarlo: ${PODMAN[*]} rm -f $NAME
================================================
MSG
	exit 0
fi

# El contenedor se borra al salir pase lo que pase.
trap '"${PODMAN[@]}" rm -f "$NAME" >/dev/null 2>&1 || true' EXIT

echo ">> abriendo sesión SSH remota (upterm), solo para github:$LAB_GITHUB_USER…"
# --detach deja la sesión corriendo en segundo plano e imprime cómo entrar. La
# sesión es una shell de login del usuario lab; al salir de ella, termina.
# --skip-host-key-check acepta la clave del servidor de upterm la primera vez:
# la máquina es efímera y no guarda nada sensible. Con tope de 2 minutos para que,
# si el servidor no responde, el job falle rápido en vez de comerse las 2 h.
if ! lab_exec "timeout 120 upterm host --detach --accept --skip-host-key-check \
	--name lab --authorized-user 'github:$LAB_GITHUB_USER' -- bash -l"; then
	echo "upterm no ha conseguido abrir la sesión (¿hay salida a uptermd.upterm.dev?)." >&2
	exit 1
fi
cat <<'MSG'
=================== LAB LISTO ===================
Entra con el comando ssh de arriba. Estarás como el usuario lab (sudo sin
contraseña) en una Fedora con systemd y dots aplicado.
Prueba: ds --help, systemctl status, ls -la ~
Para cerrar el lab: sal de la shell (exit) o haz touch /tmp/close-lab
================================================
MSG
# Esperamos hasta que la sesión termine o se pida el cierre. Miramos si sigue vivo
# el proceso de upterm: `upterm session info` devuelve 0 incluso con la sesión ya
# terminada (la enseña con "Status: ended").
while lab_exec "pgrep -u lab -x upterm >/dev/null" &&
	! "${PODMAN[@]}" exec "$NAME" test -e /tmp/close-lab; do
	sleep 10
done
echo ">> sesión terminada o cierre solicitado, saliendo."

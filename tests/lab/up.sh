#!/usr/bin/env bash
#
# Levanta el lab: una Fedora con systemd y dotsmithy aplicado, para trastear a mano.
#
#   1) compila el motor (linux/amd64, estático) en tests/lab/dots;
#   2) construye la imagen (tests/lab/Containerfile);
#   3) arranca el contenedor con systemd como PID 1 y espera a que esté listo;
#   4) aplica el contenido de muestra como el usuario "lab" (dots init).
#
# Después, según LAB_TMATE:
#   - LAB_TMATE=1 (lo usa el workflow Lab): abre una sesión tmate y se queda
#     esperando hasta que alguien haga `touch /tmp/close-lab` dentro del lab;
#     al salir, borra el contenedor.
#   - sin definir (uso local): deja el contenedor corriendo y te dice cómo entrar
#     y cómo tirarlo.
#
# Necesita podman en el host: el arranque con systemd se apoya en
# `podman run --systemd=always`, que en docker no existe tal cual. Por defecto va
# rootless; con LAB_PODMAN se cambia el comando (el workflow usa "sudo podman",
# porque el podman rootless de los runners de Ubuntu no deja leer /etc/shadow a
# sudo ni da red al contenedor como hace falta).
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$HERE/../.." && pwd)"
readonly IMAGE="dotsmithy-lab:latest"
readonly NAME="${LAB_NAME:-dotsmithy-lab}"

# Comando de podman a usar, troceado en palabras (p.ej. "sudo podman").
read -r -a PODMAN <<<"${LAB_PODMAN:-podman}"

if ! command -v podman >/dev/null 2>&1; then
	echo "Necesito podman en el sistema para el lab (arranca Fedora con systemd) y no lo encuentro." >&2
	exit 1
fi

# lab_exec ejecuta un comando dentro del lab como el usuario "lab", con shell de
# login (para que tenga el PATH y el entorno de una sesión normal).
lab_exec() {
	"${PODMAN[@]}" exec -u lab -w /home/lab "$NAME" bash -lc "$1"
}

echo ">> compilando el motor (linux/amd64, estático)…"
(
	cd "$ROOT"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -trimpath -ldflags "-s -w" -o "$HERE/dots" ./cmd/dots
)

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

echo ">> aplicando el contenido de muestra como el usuario lab (dots init)…"
lab_exec "dots init --profile minipc --content ~/content" || echo "(init terminó con avisos)"

if [ "${LAB_TMATE:-}" != "1" ]; then
	cat <<MSG
=================== LAB LISTO ===================
Fedora con systemd y dots aplicado, usuario "lab" (sudo sin contraseña).
  Entrar:  ${PODMAN[*]} exec -it -u lab -w /home/lab $NAME bash -l
  Tirarlo: ${PODMAN[*]} rm -f $NAME
================================================
MSG
	exit 0
fi

# Modo workflow: el contenedor se borra al salir pase lo que pase.
trap '"${PODMAN[@]}" rm -f "$NAME" >/dev/null 2>&1 || true' EXIT

echo ">> abriendo sesión SSH (tmate)…"
# Arrancamos el servidor de tmate (con log detallado en la home de lab, -v) y
# preguntamos por la dirección SSH hasta que aparezca. No usamos `wait
# tmate-ready`: si la sesión ya estaba lista cuando se lanza, se queda esperando
# un evento que ya pasó. Con tope de 2 minutos, para no comernos las 2 h del job.
lab_exec "tmate -v -S /tmp/t.sock new-session -d"
ssh_line=""
for _ in $(seq 1 24); do
	ssh_line="$(lab_exec 'tmate -S /tmp/t.sock display -p "#{tmate_ssh}"' 2>/dev/null || true)"
	[ -n "$ssh_line" ] && break
	sleep 5
done
if [ -z "$ssh_line" ]; then
	echo "tmate no ha conseguido conectar con su servidor en 2 minutos (¿hay salida a tmate.io?)." >&2
	echo "--- log de tmate:" >&2
	lab_exec 'tail -n 40 ~/tmate-*.log' >&2 || true
	exit 1
fi
echo "=================== LAB LISTO ==================="
echo "SSH:  $ssh_line"
echo "WEB:  $(lab_exec 'tmate -S /tmp/t.sock display -p "#{tmate_web}"')"
echo "Entras como el usuario lab (sudo sin contraseña) en una Fedora con systemd y dots aplicado."
echo "Prueba: ds --help, systemctl status, ls -la ~"
echo "Para cerrar el lab antes de tiempo: touch /tmp/close-lab"
echo "================================================"
until "${PODMAN[@]}" exec "$NAME" test -e /tmp/close-lab; do sleep 10; done
echo ">> cierre solicitado, saliendo."

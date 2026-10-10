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
# `podman run --systemd=always`, que en docker no existe tal cual.
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$HERE/../.." && pwd)"
readonly IMAGE="dotsmithy-lab:latest"
readonly NAME="${LAB_NAME:-dotsmithy-lab}"

if ! command -v podman >/dev/null 2>&1; then
	echo "Necesito podman en el sistema para el lab (arranca Fedora con systemd) y no lo encuentro." >&2
	exit 1
fi

# lab_exec ejecuta un comando dentro del lab como el usuario "lab", con shell de
# login (para que tenga el PATH y el entorno de una sesión normal).
lab_exec() {
	podman exec -u lab -w /home/lab "$NAME" bash -lc "$1"
}

echo ">> compilando el motor (linux/amd64, estático)…"
(
	cd "$ROOT"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -trimpath -ldflags "-s -w" -o "$HERE/dots" ./cmd/dots
)

echo ">> construyendo la imagen del lab…"
podman build -q -t "$IMAGE" "$HERE" >/dev/null

echo ">> arrancando Fedora con systemd…"
podman rm -f "$NAME" >/dev/null 2>&1 || true
podman run -d --name "$NAME" --systemd=always "$IMAGE" >/dev/null

# En un contenedor lo normal es "degraded" (hay unidades que no pintan nada ahí y
# fallan), así que no miramos el código de salida, solo que haya terminado de
# arrancar. Si systemd no fuera el PID 1, esto fallaría y lo cantamos.
state="$(podman exec "$NAME" systemctl is-system-running --wait 2>/dev/null || true)"
case "$state" in
running | degraded) echo ">> systemd listo (estado: $state)" ;;
*)
	echo "systemd no ha arrancado bien dentro del lab (estado: '${state:-desconocido}')." >&2
	podman logs "$NAME" 2>&1 | tail -20 >&2
	exit 1
	;;
esac

echo ">> aplicando el contenido de muestra como el usuario lab (dots init)…"
lab_exec "dots init --profile minipc --content ~/content" || echo "(init terminó con avisos)"

if [ "${LAB_TMATE:-}" != "1" ]; then
	cat <<MSG
=================== LAB LISTO ===================
Fedora con systemd y dots aplicado, usuario "lab" (sudo sin contraseña).
  Entrar:  podman exec -it -u lab -w /home/lab $NAME bash -l
  Tirarlo: podman rm -f $NAME
================================================
MSG
	exit 0
fi

# Modo workflow: el contenedor se borra al salir pase lo que pase.
trap 'podman rm -f "$NAME" >/dev/null 2>&1 || true' EXIT

echo ">> abriendo sesión SSH (tmate)…"
# Si tmate no llega a su servidor, `wait tmate-ready` espera para siempre y el job
# se comería las 2 h sin decir nada: le ponemos un tope y lo cantamos.
if ! lab_exec "tmate -S /tmp/t.sock new-session -d && timeout 120 tmate -S /tmp/t.sock wait tmate-ready"; then
	echo "tmate no ha conseguido conectar con su servidor en 2 minutos (¿hay salida a tmate.io?)." >&2
	exit 1
fi
echo "=================== LAB LISTO ==================="
echo "SSH:  $(lab_exec 'tmate -S /tmp/t.sock display -p "#{tmate_ssh}"')"
echo "WEB:  $(lab_exec 'tmate -S /tmp/t.sock display -p "#{tmate_web}"')"
echo "Entras como el usuario lab (sudo sin contraseña) en una Fedora con systemd y dots aplicado."
echo "Prueba: ds --help, systemctl status, ls -la ~"
echo "Para cerrar el lab antes de tiempo: touch /tmp/close-lab"
echo "================================================"
until podman exec "$NAME" test -e /tmp/close-lab; do sleep 10; done
echo ">> cierre solicitado, saliendo."

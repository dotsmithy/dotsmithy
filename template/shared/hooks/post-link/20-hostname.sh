#!/usr/bin/env bash
# Pone el nombre del equipo según la máquina del perfil.
#
# El nombre vive en machines/<máquina>/system/hostname, y dots.yaml lo copia a
# /etc/hostname (eso vale para el próximo arranque). Este hook lo aplica además
# en caliente con hostnamectl, para no tener que reiniciar. Va en shared y usa
# DOTS_MACHINE, así que sirve para todas las máquinas: la que no tenga fichero
# de hostname, simplemente se salta.
#
# OJO: cambia el nombre del equipo de verdad. Ajusta los ficheros de hostname a
# tus máquinas o borra este hook y las entradas de dots.yaml.
set -euo pipefail

file="$DOTS_CONTENT_DIR/machines/$DOTS_MACHINE/system/hostname"
if [[ ! -r $file ]]; then
  echo "hostname: la máquina $DOTS_MACHINE no define hostname, nada que hacer."
  exit 0
fi
wanted=$(<"$file")

# hostnamectl necesita systemd (en un contenedor, por ejemplo, no lo hay).
if [[ ! -d /run/systemd/system ]]; then
  echo "hostname: sin systemd corriendo, me salto el cambio en caliente."
  exit 0
fi

current=$(hostnamectl hostname 2>/dev/null || hostname)
if [[ $current == "$wanted" ]]; then
  echo "hostname: ya es $wanted."
  exit 0
fi

sudo=()
[[ $EUID -eq 0 ]] || sudo=(sudo)
"${sudo[@]}" hostnamectl set-hostname "$wanted"
echo "hostname: $current → $wanted."

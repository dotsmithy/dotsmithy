# El entorno de desarrollo

Cómo montar el entorno para trabajar en el motor, pasar los tests y probar `dots` a
mano en una Fedora limpia.

## El devbox

Todo lo necesario vive en un [devbox](https://www.jetify.com/devbox) local: nada se
instala a nivel global en tu máquina. El `devbox.json` y el `devbox.lock` van en el
repo, así que tienes las mismas versiones que el CI.

```sh
devbox shell           # entra al entorno (go, shellcheck, goreleaser)
devbox run build       # compila el binario en dist/dots
devbox run test        # go test ./...
devbox run lint-shim   # shellcheck del shim de bootstrap
devbox run docs        # regenera la referencia de comandos
devbox run site-serve  # la web de la documentación, en local (localhost:8000)
```

## Los tests

### Unitarios

```sh
devbox run test
```

### Integración

Corren el binario real en una Fedora limpia, dentro de un contenedor:

```sh
./tests/integration/run.sh
```

Necesitan un motor de contenedores en el sistema (`podman` o `docker`), que **no**
va en el devbox a propósito: son herramientas con mucha integración con el host
(rootless, subuid/subgid, cgroups) y vía Nix dan guerra. Puedes elegir cuál con
`DOTS_CONTAINER_ENGINE=podman|docker`.

### El shim de bootstrap

El instalador (`script/install.sh`) es lo único del motor escrito en bash. Lo revisa
`shellcheck`, y un test simula una instalación completa (con una descarga de
mentira, sin red):

```sh
devbox run lint-shim
./tests/shim/run.sh
```

En cada PR, el CI pasa todo esto (`shellcheck`, además, sobre todos los `.sh` del
repo) y alguna comprobación más: lo cuenta
[cómo se trabaja en el repo](workflow.md#los-checks).

## El Lab

Una Fedora con systemd y `dots` ya aplicado, en local o en GitHub, para trastear a
mano: [tests/lab](../tests/lab/README.md).

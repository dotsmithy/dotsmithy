# Contribuir

> [!NOTE]
> Esta parte está en obras. Lo que ya hay: cómo se trabaja en el repo, cómo montar
> el entorno y pasar los tests, y cómo se publica una versión.

## Cómo se trabaja en el repo

Ramas (`main`, `next`, `feat/`, `hotfix/`), pull requests, checks, Draft, issues y la
plantilla de contenido: [flujo de trabajo](flujo-de-trabajo.md).

## Lo que vendrá

- **Entorno de desarrollo**: el devbox y qué trae, con más detalle.
- **Tests**: unitarios, integración en Fedora y el shim de bootstrap.
- **El Lab**: una Fedora con systemd para probar a mano.

## Mientras tanto

### Entorno

El toolchain vive en un [devbox](https://www.jetify.com/devbox) local, nada se
instala a nivel global:

```sh
devbox shell           # entra al entorno (go, shellcheck, goreleaser)
devbox run build       # compila el binario en dist/dots
devbox run test        # go test ./...
devbox run lint-shim   # shellcheck del shim de bootstrap
devbox run docs        # regenera la referencia de comandos
devbox run site-serve  # la web de la documentación, en local (localhost:8000)
```

### Tests de integración

Corren el binario real en una Fedora limpia, dentro de un contenedor:

```sh
./tests/integration/run.sh
```

Necesitan un motor de contenedores en el sistema (`podman` o `docker`), que **no**
va en el devbox a propósito: son herramientas con mucha integración con el host
(rootless, subuid/subgid, cgroups) y vía Nix dan guerra. Puedes elegir cuál con
`DOTS_CONTAINER_ENGINE=podman|docker`.

En CI, la integración corre en cada PR.

### Más

- [El Lab](../tests/lab/README.md): una Fedora con systemd y `dots` aplicado, en
  local o en GitHub, para trastear a mano.
- [Publicar una versión](../RELEASING.md).

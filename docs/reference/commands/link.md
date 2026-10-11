<!-- Generado por tools/docgen: no lo edites a mano. -->

# dots link

Despliega los dotfiles a $HOME (symlinks, modelo Stow).

link resuelve la pila de capas del perfil (shared, una por dimensión y la
máquina) y crea los symlinks en tu $HOME apuntando a los ficheros reales del
repo. Es idempotente: si ya está enlazado, no hace nada; si falta, lo crea. Con
--dry-run enseña el plan sin tocar nada.

Si hay un fichero real donde iría un symlink, por defecto aborta sin tocar nada.
Con --backup, mueve ese fichero a un «.dots-bak» y crea el symlink.

Los flags mandan; lo que no indiques, se coge del estado local (lo que dejó
init). Si no hay estado ni flags, no puede saber qué perfil aplicar.

## Uso

```
dots link [opciones]
```

## Opciones

| Opción | Qué hace |
| --- | --- |
| `--backup` | si hay un fichero real en medio, respáldalo (.dots-bak) y enlaza |
| `--content <texto>` | ruta al repo de contenido (por defecto: la del estado local, o el directorio actual) |
| `--profile <texto>` | perfil de máquina a aplicar (por defecto: el del estado local) |

## Opciones globales

| Opción | Qué hace |
| --- | --- |
| `-n`, `--dry-run` | enseña qué haría sin tocar nada |
| `-h`, `--help` | enseña esta ayuda |
| `-v`, `--verbose` | saca más detalle de lo que va pasando |
| `-y`, `--yes` | no pregunta nada (modo no interactivo) |

## Ver también

- [`dots`](dots.md): dotsmithy: la herrería de tus dotfiles
- [Todos los comandos](README.md)

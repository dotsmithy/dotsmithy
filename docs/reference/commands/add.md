<!-- Generado por tools/docgen: no lo edites a mano. -->

# dots add

Registra e instala uno o varios paquetes.

add coge uno o varios paquetes de un gestor (ej. dnf), los registra en la
lista de una capa del contenido y los instala. Deja el alta declarada en el repo
y aplicada en la máquina.

Por defecto los registra en la capa de la máquina del perfil activo
(machines/&lt;máquina&gt;), que es lo más prudente: así no se cuelan en otras
máquinas. Con --layer eliges otra (shared, ux/terminal...), pero tiene que ser
una de las capas del perfil activo, porque registrar en una capa que este
perfil no aplica no instalaría nada aquí.

Ejemplos:
  dots add dnf ripgrep fzf
  dots add dnf tmux --layer shared

## Uso

```
dots add <gestor> <paquete>... [opciones]
```

## Opciones

| Opción | Qué hace |
| --- | --- |
| `--content <texto>` | ruta al repo de contenido (por defecto: la del estado local, o el directorio actual) |
| `--layer <texto>` | capa del perfil activo donde registrar (por defecto: machines/&lt;máquina&gt;) |
| `--profile <texto>` | perfil activo (por defecto: el del estado local) |

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

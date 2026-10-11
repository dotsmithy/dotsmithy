<!-- Generado por tools/docgen: no lo edites a mano. -->

# dots init

Arranca una máquina desde cero (perfil, contenido, estado, aplica).

init deja una máquina lista: elige el perfil, consigue tu repo de contenido
(clonándolo desde --remote por SSH, o usando el de --content si ya lo tienes en
local), guarda el estado local y aplica todo (link).

Si clonas por SSH y no tienes clave, init te genera una (ed25519) y confía en el
host; solo tienes que registrar esa clave en GitHub y reejecutar. Si el clone
falla por permisos, te enseña la clave pública y cómo añadirla.

## Uso

```
dots init [opciones]
```

## Opciones

| Opción | Qué hace |
| --- | --- |
| `--backup` | si hay ficheros reales donde irían symlinks, respáldalos (.dots-bak) y enlaza |
| `--content <texto>` | ruta a un repo de contenido ya presente en local (alternativa a --remote) |
| `--profile <texto>` | perfil de máquina a instalar (ej. minipc) |
| `--remote <texto>` | URL (SSH) del repo de contenido a clonar |

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

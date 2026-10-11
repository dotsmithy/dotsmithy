<!-- Generado por tools/docgen: no lo edites a mano. -->

# dots update

Trae los cambios del contenido y los aplica.

update hace pull de tu repo de contenido y vuelve a aplicar: enlaza lo
nuevo (link + post-link), sincroniza paquetes (+ post-packages) y deja todo al
día. No ejecuta post-init (eso es del primer init).

Perfil y contenido se resuelven como en link: flags o el estado local.

## Uso

```
dots update [opciones]
```

## Opciones

| Opción | Qué hace |
| --- | --- |
| `--backup` | si hay ficheros reales donde irían symlinks, respáldalos (.dots-bak) y enlaza |
| `--content <texto>` | ruta al repo de contenido (por defecto: la del estado local) |
| `--profile <texto>` | perfil de máquina (por defecto: el del estado local) |

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

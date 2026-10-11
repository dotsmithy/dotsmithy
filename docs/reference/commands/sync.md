<!-- Generado por tools/docgen: no lo edites a mano. -->

# dots sync

Informa de qué paquetes declarados faltan por instalar (solo lectura).

sync mira los paquetes que declaras (unión de las capas activas) y los
compara con lo que hay instalado, contándote qué falta. De momento es solo
lectura: no instala nada (eso vendrá con add/instalación).

Perfil y contenido se resuelven igual que en link: flags, o el estado local.

## Uso

```
dots sync [opciones]
```

## Opciones

| Opción | Qué hace |
| --- | --- |
| `--content <texto>` | ruta al repo de contenido (por defecto: la del estado local, o el directorio actual) |
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

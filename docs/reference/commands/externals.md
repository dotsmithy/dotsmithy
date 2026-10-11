<!-- Generado por tools/docgen: no lo edites a mano. -->

# dots externals

Trae y prepara los recursos externos declarados.

externals baja los recursos externos que declaras (repos git o ficheros)
a su destino y ejecuta su post cuando los trae o cuando cambian. Va en modo
best-effort: si alguno falla, te avisa y sigue con el resto.

Perfil y contenido se resuelven como en link: flags o el estado local.

## Uso

```
dots externals [opciones]
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

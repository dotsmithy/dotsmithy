<!-- Generado por tools/docgen: no lo edites a mano. -->

# dots completion powershell

Genera el script de autocompletado para powershell.

Para cargarlo en la sesión actual:

```sh
dots completion powershell | Out-String | Invoke-Expression
```

Para tenerlo en todas las sesiones nuevas, añade la salida del comando anterior
a tu perfil de powershell.

## Uso

```
dots completion powershell [opciones]
```

## Opciones

| Opción | Qué hace |
| --- | --- |
| `--no-descriptions` | sin descripciones en las sugerencias |

## Opciones globales

| Opción | Qué hace |
| --- | --- |
| `-n`, `--dry-run` | enseña qué haría sin tocar nada |
| `-h`, `--help` | enseña esta ayuda |
| `-v`, `--verbose` | saca más detalle de lo que va pasando |
| `-y`, `--yes` | no pregunta nada (modo no interactivo) |

## Ver también

- [`dots completion`](completion.md): genera el script de autocompletado para tu shell
- [Todos los comandos](README.md)

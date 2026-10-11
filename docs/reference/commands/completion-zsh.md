<!-- Generado por tools/docgen: no lo edites a mano. -->

# dots completion zsh

Genera el script de autocompletado para zsh.

Si no tienes el autocompletado activado en tu entorno, actívalo una vez con:

```sh
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

Para cargarlo en la sesión actual:

```sh
source <(dots completion zsh)
```

Para tenerlo en todas las sesiones nuevas, ejecuta una vez:

```sh
dots completion zsh > "${fpath[1]}/_dots"
```

Abre una shell nueva para que surta efecto.

## Uso

```
dots completion zsh [opciones]
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

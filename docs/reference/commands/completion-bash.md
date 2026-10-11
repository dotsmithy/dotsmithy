<!-- Generado por tools/docgen: no lo edites a mano. -->

# dots completion bash

Genera el script de autocompletado para bash.

Necesita el paquete «bash-completion»; si no lo tienes, instálalo con el gestor
de paquetes de tu sistema.

Para cargarlo en la sesión actual:

```sh
source <(dots completion bash)
```

Para tenerlo en todas las sesiones nuevas, ejecuta una vez:

```sh
dots completion bash > /etc/bash_completion.d/dots
```

Abre una shell nueva para que surta efecto.

## Uso

```
dots completion bash
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

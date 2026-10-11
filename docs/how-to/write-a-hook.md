# Escribir un hook

Un hook es un script tuyo que `dots` ejecuta en un momento concreto: después de
enlazar los dotfiles, después de instalar paquetes o al terminar de montar una
máquina. Sirve para lo que no es un fichero ni un paquete: activar un servicio,
cambiar tu shell por defecto, preparar algo de un programa.

Como ejemplo, activaremos el servicio de Syncthing para tu usuario después de
instalar los paquetes.

## 1. Elige el momento y la capa

| Punto | Cuándo |
| --- | --- |
| `post-link` | Después de enlazar los dotfiles. |
| `post-packages` | Después de instalar paquetes. |
| `post-init` | Solo al final del primer `dots init`. |

Nuestro ejemplo necesita que Syncthing ya esté instalado, así que va en
`post-packages`. Y como lo queremos en todas las máquinas, en la capa `shared`.

## 2. Escribe el script

```sh
mkdir -p shared/hooks/post-packages
```

`shared/hooks/post-packages/20-syncthing.sh`:

```bash
#!/usr/bin/env bash
# Activa Syncthing para el usuario. Si ya está activo, no hace nada.
set -euo pipefail

if ! systemctl --user is-enabled --quiet syncthing.service; then
  systemctl --user enable --now syncthing.service
  echo "syncthing activado"
fi
```

Tres cosas que importan:

- **El nombre empieza por un número.** Los hooks de todas las capas se ejecutan
  por orden de nombre, así que `10-` va antes que `20-`, vengan de la capa que
  vengan.
- **Tiene que poder repetirse.** Se ejecuta en cada `dots update`, así que
  comprueba antes de cambiar nada, como hace el `if` del ejemplo.
- **No puede preguntar nada.** Corre sin teclado; si necesita confirmación, se
  queda colgado.

## 3. Pruébalo

Para ver qué hooks correrían, sin ejecutar ninguno:

```sh
dots update --dry-run
```

Y para ejecutarlo de verdad, `dots update`. Si el hook falla, `dots` se para ahí y
te dice cuál ha sido.

> [!TIP]
> Un `dots link` suelto no ejecuta hooks: solo lo hacen `dots init` y
> `dots update`.

## Usar el contexto

Cada hook recibe variables con el perfil activo, la máquina, el valor de cada
dimensión y más. Por ejemplo, para hacer algo solo en las máquinas con escritorio
sin crear un hook en otra capa:

```bash
if [[ "${DOTS_PARAM_UX:-}" == desktop ]]; then
  # ...
fi
```

Aunque si es algo solo de las máquinas con escritorio, lo más limpio suele ser
poner el hook en la capa `ux/desktop`.

## Ver también

- [Todas las variables y reglas de los hooks](../reference/content-layout.md#hooks-scripts-en-el-momento-justo)
- [`hooks.points` en el manifiesto](../reference/manifest.md#hooks), para
  desactivar un punto entero.

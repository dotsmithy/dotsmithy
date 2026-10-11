# Cómo se organiza el contenido

Qué carpetas y ficheros espera el motor en tu repo de contenido. Lo que no está
aquí, el motor no lo mira: puedes tener un `README.md`, notas o scripts tuyos sin
que estorben.

## Vista general

```
dots.yaml                    el manifiesto (obligatorio, en la raíz)
shared/                      capa común (opcional)
<dimensión>/<valor>/         una capa por cada valor de cada dimensión
machines/<máquina>/          una capa por máquina
```

Todas las capas tienen la misma forma por dentro:

```
<capa>/
  home/                      dotfiles que se enlazan en $HOME
  packages/<gestor>.txt      listas de paquetes
  hooks/<punto>/NN-*.sh      scripts que se ejecutan en cada punto
  system/                    ficheros de sistema (los declara dots.yaml)
  ...                        lo que necesiten tus externals o hooks
```

Ninguna subcarpeta es obligatoria: una capa puede estar vacía.

## Las capas

| Capa | Ruta | Cuándo existe |
| --- | --- | --- |
| Común | `shared/` | Opcional. Si no está, la pila empieza por las dimensiones. |
| Dimensión | `<dimensión>/<valor>/` | Una por cada valor que use algún perfil (`ux/desktop`, `identity/work`...). Tiene que existir aunque esté vacía. |
| Máquina | `machines/<máquina>/` | Una por cada `machine` de algún perfil. Tiene que existir aunque esté vacía. |

Una carpeta vacía no se guarda en git: mete dentro un `.gitkeep` para que exista.

Las dimensiones, los perfiles y sus valores se declaran en
[`dots.yaml`](manifest.md#dimensions). El orden en que se apilan las capas y qué
pasa cuando coinciden está en [cómo se combinan las capas](manifest.md#cómo-se-combinan-las-capas).

## `home/`: los dotfiles

Un espejo de tu `$HOME`. Cada fichero se enlaza en la misma ruta relativa:

| En el contenido | Symlink en |
| --- | --- |
| `shared/home/.bashrc.d/aliases.sh` | `~/.bashrc.d/aliases.sh` |
| `machines/laptop/home/.config/foo/config.toml` | `~/.config/foo/config.toml` |

- Se enlazan **ficheros**, nunca carpetas enteras. Las carpetas intermedias se
  crean de verdad en el destino, así que varias capas pueden poner ficheros en el
  mismo `~/.config`.
- La carpeta se llama `home` y el destino es `~` salvo que digas otra cosa en
  [`stow`](manifest.md#stow).

## `packages/`: las listas de paquetes

Un fichero por gestor: `packages/dnf.txt`. El gestor tiene que estar activo en
[`packages.managers`](manifest.md#packages).

```
# Un paquete por línea.
ripgrep
fd-find      # también valen comentarios al final de la línea
tmux
```

- Las líneas vacías y lo que va detrás de `#` se ignoran.
- Se juntan las listas de todas las capas de la pila.
- `dots add` escribe aquí por ti (en la capa de la máquina, salvo `--layer`).

## `hooks/`: scripts en el momento justo

Scripts en `hooks/<punto>/`, donde el punto es uno de estos:

| Punto | Cuándo se ejecuta |
| --- | --- |
| `post-link` | Después de enlazar los dotfiles, en `init` y `update`. Un `dots link` suelto no ejecuta hooks. |
| `post-packages` | Después de instalar paquetes, en `init` y `update`. |
| `post-init` | Al final del primer `dots init`, y solo ahí. |

- Solo cuentan los ficheros `.sh` y se ejecutan con `bash`. No hace falta que sean
  ejecutables.
- Nombre con prefijo numérico: `NN-nombre.sh`. En cada punto se juntan los de
  todas las capas y se ejecutan por orden de nombre. Si dos se llaman igual, va
  antes el de la capa anterior en la pila.
- Corren con el directorio de trabajo en la raíz del contenido.
- Si uno falla, `dots` se para ahí y dice cuál.
- Tienen que ser idempotentes (se repiten en cada `update`) y no pueden
  pedir nada por teclado.
- Los hooks en `hooks/` en la raíz del contenido, fuera de una capa, no se
  ejecutan. `dots` avisa si los encuentra.
- Con `hooks.points` en el manifiesto puedes desactivar puntos enteros.

### Variables que recibe cada hook

| Variable | Valor |
| --- | --- |
| `DOTS_HOOK` | Punto que se está ejecutando (ej. `post-link`). |
| `DOTS_PROFILE` | Perfil activo. |
| `DOTS_MACHINE` | Máquina del perfil. |
| `DOTS_PARAM_<DIM>` | Valor del perfil en cada dimensión, con el nombre en mayúsculas (ej. `DOTS_PARAM_UX=terminal`). Lo que no vale en un nombre de variable pasa a `_`. |
| `DOTS_CONTENT_DIR` | Raíz del contenido. |
| `DOTS_TARGET` | Destino de los symlinks, ya expandido (ej. tu `$HOME`). |
| `DOTS_LAYERS` | Pila de capas, separadas por espacios. |

El script `post` de un external recibe las mismas, con `DOTS_HOOK=externals`, y
además `DOTS_EXTERNAL_DEST` con el destino del recurso.

## `system/`: ficheros de sistema

Ficheros que van fuera de tu `$HOME` (`/etc`, unidades de systemd...). El nombre
`system/` es una convención: lo que manda es la entrada `system` de la capa en
[`dots.yaml`](manifest.md#layers), que dice qué fichero va a dónde.

```yaml
layers:
  machines/laptop:
    system:
      - src: system/hostname     # relativo a machines/laptop/
        dest: /etc/hostname
        type: copy
```

## Otros ficheros de la capa

Los scripts `post` de los externals y cualquier fichero que usen tus hooks pueden
vivir donde quieras dentro de la capa (por ejemplo, `scripts/`). Las rutas en
`dots.yaml` son relativas a la carpeta de la capa.

## Dónde deja las cosas `dots`

No es parte de tu contenido, pero viene bien saberlo:

| Qué | Dónde |
| --- | --- |
| El contenido, si lo clonó `dots init --remote` | `~/.local/share/dots/content` (respeta `XDG_DATA_HOME`) |
| El estado local (perfil y ruta del contenido) | `~/.config/dots/state.yaml` (respeta `XDG_CONFIG_HOME`) |
| Respaldos de `dots link --backup` | Junto al original, como `<fichero>.dots-bak` |

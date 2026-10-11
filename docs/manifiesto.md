# El manifiesto: `dots.yaml`

`dots.yaml` vive en la raíz de tu repo de contenido y es el único contrato entre
el motor y tus ficheros. Dice qué perfiles hay, cómo se compone cada uno y lo poco
que no se puede expresar con ficheros (repos COPR, recursos externos, ficheros de
sistema). Todo lo demás se saca de las carpetas.

Esta página describe la **versión 1** del esquema. El motor solo entiende esa: un
manifiesto con otra `version` se rechaza.

## La idea: perfiles compuestos

Un perfil no es una carpeta, sino una **pila de capas**. El motor la construye
siempre igual:

```
shared  →  <dimensión>/<valor> (una por dimensión, en orden)  →  machines/<máquina>
```

- **`shared`**: lo común a todas las máquinas. Es opcional.
- **Dimensiones**: los ejes en los que varían tus máquinas, y los declaras tú. Por
  ejemplo `ux` (con escritorio o solo terminal) e `identity` (personal o trabajo).
  Cada perfil elige un valor en cada dimensión, y la capa vive en
  `<dimensión>/<valor>` (`ux/desktop`, `identity/work`...). El motor no sabe qué
  significa cada una: solo las apila.
- **`machines/<máquina>`**: lo propio de una máquina concreta. Siempre está y
  siempre va la última.

Así, una máquina sin escritorio y con identidad de trabajo se monta con
`shared + ux/terminal + identity/work + machines/<la-suya>`, sin copiar ficheros
de ningún sitio.

## Ejemplo completo

```yaml
version: 1

dimensions: [ux, identity]

profiles:
  desk:
    description: "Escritorio de casa"
    machine: desk
    params: { ux: desktop, identity: personal }
  devbox:
    description: "Máquina remota de desarrollo"
    machine: devbox
    params: { ux: terminal, identity: work }

stow:
  target: ~
  subdir: home

packages:
  managers: [dnf]

hooks:
  points: [post-link, post-packages, post-init]

layers:
  shared:
    externals:
      - dest: ~/.config/nvim
        type: git
        repo: https://github.com/example/nvim-config
        ref: main
        post: scripts/nvim-post.sh
  ux/desktop:
    copr: [owner/desktop-shell]
  machines/desk:
    system:
      - src: system/example.conf
        dest: /etc/example.conf
        type: copy
        mode: "0644"
        validate: "grep -q example {file}"
```

Y el contenido que le corresponde:

```
dots.yaml
shared/
  home/              # dotfiles que se enlazan en $HOME
  packages/dnf.txt   # un paquete por línea
  hooks/post-link/10-algo.sh
  scripts/nvim-post.sh
ux/
  desktop/
    home/
    packages/dnf.txt
  terminal/          # puede estar vacía, pero tiene que existir
identity/
  personal/
  work/
machines/
  desk/
    system/example.conf
  devbox/
```

## Campos

### `version`

Obligatorio, y tiene que ser `1`.

### `dimensions`

Lista de dimensiones, en el orden en que se apilan. Puede estar vacía: entonces la
pila es solo `shared → machines/<máquina>`.

- No se pueden repetir.
- No pueden llamarse `shared`, `machines` ni `hooks` (son nombres del motor).
- Sin barras: cada dimensión es un nombre de carpeta.

### `profiles`

Mapa de perfiles. Hace falta al menos uno. El perfil se elige explícitamente al
hacer `dots init --profile <nombre>`, y luego queda guardado en el estado local,
así que el resto de comandos ya lo saben.

| Campo         | Obligatorio | Qué es                                                     |
| ------------- | ----------- | ---------------------------------------------------------- |
| `description` | no          | Texto libre para ti.                                       |
| `machine`     | sí          | Carpeta de la máquina dentro de `machines/`. No tiene valor por defecto. Dos perfiles pueden compartir máquina. |
| `params`      | sí, si hay dimensiones | Un valor por dimensión, **exactamente**: ni falta ninguna ni sobra ninguna. |

### `stow`

Qué se enlaza por symlinks y dónde.

| Campo    | Por defecto | Qué es                                              |
| -------- | ----------- | --------------------------------------------------- |
| `target` | `~`         | Destino de los symlinks.                            |
| `subdir` | `home`      | Subcarpeta de cada capa que se enlaza en `target`. |

### `packages`

| Campo      | Qué es                                                              |
| ---------- | ------------------------------------------------------------------- |
| `managers` | Gestores activos. Hoy solo `dnf` está implementado (el esquema reconoce también `pip`, `npm` y `cargo`, para más adelante). |

Las listas viven en `<capa>/packages/<gestor>.txt`: un paquete por línea, con
comentarios con `#`. Se juntan las de todas las capas de la pila.

### `hooks`

| Campo    | Por defecto | Qué es                                                       |
| -------- | ----------- | ------------------------------------------------------------ |
| `points` | todos       | Puntos activos: `post-link`, `post-packages`, `post-init`. |

Los hooks no se declaran aquí: son ficheros en `<capa>/hooks/<punto>/*.sh`
(detalle en la sección de hooks, más abajo).

### `layers`

Lo que cada capa declara y que no es un fichero. La clave es la ruta de la capa
(`shared`, `ux/desktop`, `machines/devbox`...), y una capa sin nada que declarar
simplemente no aparece. Cada clave tiene que ser una capa que algún perfil pueda
resolver: si no, es casi seguro una errata y el motor lo dice.

Cada capa admite:

- **`copr`**: repos COPR que habilitar antes de instalar paquetes. Exige `dnf` en
  `packages.managers`.
- **`externals`**: recursos que el motor descarga y mantiene al día (fuera del
  modelo de symlinks).

  | Campo  | Qué es                                                     |
  | ------ | ---------------------------------------------------------- |
  | `dest` | Destino directo (se admite `~`).                           |
  | `type` | `git` o `file`.                                            |
  | `repo` | Repo a clonar (obligatorio con `git`).                     |
  | `ref`  | Rama, tag o commit (con `git`).                            |
  | `url`  | URL a descargar (obligatorio con `file`).                  |
  | `post` | Script que corre cuando el recurso cambia. Ruta relativa a la carpeta de la capa. |

- **`system`**: ficheros de sistema fuera de `$HOME` (`/etc`, systemd...). Se
  colocan con `sudo` (o directamente si ya eres root).

  | Campo      | Qué es                                                   |
  | ---------- | -------------------------------------------------------- |
  | `src`      | Origen, relativo a la carpeta de la capa (ej. `system/x.conf`). |
  | `dest`     | Ruta absoluta de destino.                                |
  | `type`     | `symlink` o `copy`.                                      |
  | `mode`     | Opcional, en octal (ej. `"0440"`).                       |
  | `owner`    | Opcional, `usuario:grupo` (por defecto `root:root`).     |
  | `validate` | Opcional, comando que valida una copia temporal antes de reemplazar (`{file}` es esa copia, ej. `visudo -cf {file}`). |

## Cómo se combinan las capas

- **Dotfiles** (`<capa>/home/`): si dos capas tienen el mismo fichero, gana la
  última de la pila. La máquina pisa a todo lo demás.
- **Paquetes, COPR, externals y system**: se unen los de todas las capas.
- **Hooks**: se juntan todos y se ejecutan por nombre de fichero (ver abajo).

Las capas que no son del perfil no aportan nada: lo de `ux/desktop` no llega
nunca a una máquina con `ux: terminal`.

## Hooks

Un hook es un script bash en `<capa>/hooks/<punto>/NN-nombre.sh`.

- En cada punto se juntan los hooks de **todas** las capas de la pila y se
  ejecutan por **orden de nombre**, así que el número `NN-` manda aunque vengan de
  capas distintas. Si dos se llaman igual, va antes el de la capa anterior.
- Corren con el directorio de trabajo en la raíz del contenido.
- Si uno falla, se para ahí (y se dice cuál). En `--dry-run` no se ejecuta
  ninguno: solo se enseña qué correría.
- Los hooks en `hooks/` en la raíz del contenido (donde vivían antes) **no se
  ejecutan**. El motor avisa si los encuentra.

Cada hook recibe estas variables:

| Variable             | Valor                                                  |
| -------------------- | ------------------------------------------------------ |
| `DOTS_HOOK`          | Punto que se está ejecutando (ej. `post-link`).        |
| `DOTS_PROFILE`       | Perfil activo.                                         |
| `DOTS_MACHINE`       | Máquina del perfil.                                    |
| `DOTS_PARAM_<DIM>`   | Valor del perfil en cada dimensión, con el nombre en mayúsculas (ej. `DOTS_PARAM_UX=terminal`). Lo que no vale en un nombre de variable pasa a `_`. |
| `DOTS_CONTENT_DIR`   | Raíz del contenido.                                    |
| `DOTS_TARGET`        | Destino de los symlinks, ya expandido (ej. tu `$HOME`). |
| `DOTS_LAYERS`        | Pila de capas, separadas por espacios.                 |

El `post` de un external recibe las mismas variables (con `DOTS_HOOK=externals`)
y además `DOTS_EXTERNAL_DEST`, con el destino del recurso.

## `dots add` y las capas

`dots add <gestor> <paquete>...` registra los paquetes en una lista del contenido
y los instala. Por defecto los apunta en la capa de la **máquina** del perfil
activo, que es lo más prudente: no aparecen en otras máquinas sin que lo decidas.

Para registrarlos en otra capa, `--layer`:

```sh
dots add dnf ripgrep               # machines/<máquina>/packages/dnf.txt
dots add dnf tmux --layer shared   # shared/packages/dnf.txt
```

La capa tiene que ser una de las del perfil activo. Registrar en una capa que
este perfil no aplica (por ejemplo `ux/desktop` en una máquina sin escritorio) no
instalaría nada aquí, así que se rechaza.

## Validación

Al cargar el manifiesto, el motor comprueba todo esto y, si algo falla, enseña
todas las pegas de una vez:

- `version` es `1` y no hay campos desconocidos (una errata en un nombre de campo
  es un error, no se ignora).
- Las dimensiones no se repiten ni usan nombres reservados.
- Cada perfil tiene `machine` y sus `params` cubren exactamente las dimensiones.
- Los nombres de máquina y los valores de las dimensiones no llevan barras.
- Cada clave de `layers` es una capa que algún perfil resuelve, y lo que declara
  es coherente (tipos válidos, campos obligatorios, `mode` en octal...).
- Existen en disco las carpetas de todas las capas que resuelven los perfiles
  (`shared/` es la única opcional). Pueden estar vacías.

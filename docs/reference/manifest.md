# El manifiesto: `dots.yaml`

`dots.yaml` vive en la raíz de tu repo de contenido y es el único contrato entre
el motor y tus ficheros. Dice qué perfiles hay, cómo se compone cada uno y lo poco
que no se puede expresar con ficheros (repos COPR, recursos externos, ficheros de
sistema). Todo lo demás se saca de las carpetas.

Esta página describe la **versión 1** del esquema. El motor solo entiende esa: un
manifiesto con otra `version` se rechaza.

## La pila de capas

Cada perfil se monta apilando capas, siempre en este orden:

```
shared  →  <dimensión>/<valor> (una por dimensión, en orden)  →  machines/<máquina>
```

`shared` es opcional; la de la máquina siempre está y siempre va la última. El
porqué de este modelo está en [perfiles y capas](../explanation/profiles-and-layers.md).

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

Y el contenido que le corresponde (detalle en
[cómo se organiza el contenido](content-layout.md)):

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

Las listas viven en `<capa>/packages/<gestor>.txt`
([formato](content-layout.md#packages-las-listas-de-paquetes)). Se juntan las de
todas las capas de la pila.

### `hooks`

| Campo    | Por defecto | Qué es                                                       |
| -------- | ----------- | ------------------------------------------------------------ |
| `points` | todos       | Puntos activos: `post-link`, `post-packages`, `post-init`. |

Los hooks no se declaran aquí: son ficheros en `<capa>/hooks/<punto>/*.sh`. Cómo
se nombran, en qué orden corren y qué variables reciben está en
[cómo se organiza el contenido](content-layout.md#hooks-scripts-en-el-momento-justo).

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
  | `post` | Script que corre cuando el recurso cambia. Ruta relativa a la carpeta de la capa. Recibe las [mismas variables que un hook](content-layout.md#variables-que-recibe-cada-hook) más `DOTS_EXTERNAL_DEST`. |

  Si un external falla, `dots` avisa y sigue con los demás.

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
- **Hooks**: se juntan todos y se ejecutan por nombre de fichero
  ([detalle](content-layout.md#hooks-scripts-en-el-momento-justo)).

Las capas que no son del perfil no aportan nada: lo de `ux/desktop` no llega
nunca a una máquina con `ux: terminal`.

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

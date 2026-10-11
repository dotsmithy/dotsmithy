# Añadir un dotfile

Tienes un fichero de configuración que quieres llevar a todas tus máquinas (o solo
a algunas). Aquí lo metes en tu contenido y dejas que `dots` lo enlace.

Usaremos como ejemplo `~/.config/starship.toml`. Cambia la ruta por la de tu
fichero.

## 1. Elige la capa

La capa decide qué máquinas reciben el fichero:

| Si lo quieres en... | Va en la capa... |
| --- | --- |
| todas tus máquinas | `shared` |
| las que comparten un rasgo (por ejemplo, las que tienen escritorio) | la de esa dimensión, como `ux/desktop` |
| solo esta máquina | `machines/<tu-máquina>` |

Si dudas, empieza por la de la máquina: siempre puedes moverlo a `shared` cuando
lo eches de menos en otra. Lo de las capas lo tienes contado con calma en
[perfiles y capas](../explanation/profiles-and-layers.md).

## 2. Mueve el fichero a tu contenido

Dentro de cada capa, la carpeta `home/` es un espejo de tu `$HOME`: el fichero va
en la misma ruta relativa que tiene en tu `$HOME`.

```sh
cd ~/.local/share/dots/content
mkdir -p shared/home/.config
mv ~/.config/starship.toml shared/home/.config/starship.toml
```

> [!TIP]
> `~/.local/share/dots/content` es donde `dots init --remote` clona tu contenido.
> Si lo tienes en otro sitio (lo inicializaste con `--content`), usa esa ruta.

Muévelo, no lo copies. Si dejas el original en su sitio, `dots` no lo pisa y se
para con un conflicto (lo ves en el paso siguiente).

## 3. Enlázalo

Primero mira qué va a hacer, sin tocar nada:

```console
$ dots link --dry-run
dry-run del perfil "laptop" (capas: shared, machines/laptop):
  + /home/tu/.config/starship.toml -> /home/tu/.local/share/dots/content/shared/home/.config/starship.toml
dry-run: 1 cambio(s), 0 ya ok, 0 conflicto(s).
```

Si te cuadra, enlaza de verdad:

```sh
dots link
```

Ahora `~/.config/starship.toml` es un symlink al fichero de tu contenido. Lo
editas como siempre y el cambio ya está en el repo.

### Si te sale un conflicto

```console
ups, algo ha fallado: hay 1 conflicto(s): tienes ficheros reales donde irían symlinks. Quítalos, muévelos o usa --backup:
  - /home/tu/.config/starship.toml
```

Hay un fichero de verdad donde iría el symlink, y `dots` nunca borra nada tuyo.
Quítalo tú o deja que lo aparte:

```sh
dots link --backup
```

El original se queda al lado como `starship.toml.dots-bak`, por si acaso.

## 4. Guárdalo en el repo

```sh
cd ~/.local/share/dots/content
git add shared/home/.config/starship.toml
git commit -m "Añade la config de starship"
git push
```

En tus otras máquinas, un `dots update` lo trae y lo enlaza.

## Ver también

- [El día a día: actualizar y sincronizar](update-and-sync.md)
- [Cómo se organiza el contenido](../reference/content-layout.md)
- [`stow` en el manifiesto](../reference/manifest.md#stow), si quieres cambiar la
  carpeta `home/` o el destino.

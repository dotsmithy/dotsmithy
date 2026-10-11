# Instalar paquetes

Quieres que tus máquinas tengan instalados ciertos programas. En dotsmithy, cada
capa tiene sus listas de paquetes y cada máquina instala los de su pila.

> [!NOTE]
> De momento el gestor soportado es `dnf`.

## La forma rápida: `dots add`

```sh
dots add dnf ripgrep fzf
```

Hace dos cosas a la vez: apunta los paquetes en la lista de tu contenido y los
instala en esta máquina. Por defecto los apunta en la capa de la **máquina**, así
que no aparecen en las demás sin que lo decidas.

Si los quieres en otra capa, dilo con `--layer`:

```sh
dots add dnf tmux --layer shared        # en todas tus máquinas
dots add dnf gimp --layer ux/desktop    # solo en las que tienen escritorio
```

La capa tiene que ser una de las del perfil de esta máquina. Lo de `ux/desktop` no
se puede registrar desde una máquina sin escritorio, porque aquí no instalaría
nada.

Después, guarda el cambio en el repo:

```sh
cd ~/.local/share/dots/content
git add . && git commit -m "Añade ripgrep y fzf" && git push
```

## A mano: editando la lista

Las listas son ficheros de texto en `<capa>/packages/dnf.txt`, un paquete por
línea:

```
# Herramientas de terminal
ripgrep
fzf        # buscador difuso
```

Edita la que quieras, guarda y aplica:

```sh
dots update
```

## Paquetes de un repo COPR

Si un paquete viene de [COPR](https://copr.fedorainfracloud.org/), declara el repo
en la capa que lo necesita, en `dots.yaml`:

```yaml
layers:
  ux/desktop:
    copr: [owner/proyecto]
```

`dots` lo habilita antes de instalar los paquetes de esa capa. Solo las máquinas
con esa capa en su pila lo habilitan.

## ¿Falta algo por instalar?

```sh
dots sync
```

Compara lo que declaras con lo que hay instalado y te dice qué falta. Solo mira:
no instala nada. Para instalarlo, `dots update`.

## Ver también

- [Formato de las listas](../reference/content-layout.md#packages-las-listas-de-paquetes)
- [`packages` y `copr` en el manifiesto](../reference/manifest.md#packages)

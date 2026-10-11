# Preparar una máquina nueva

Tienes un equipo nuevo (o uno que aún no gestionabas con dotsmithy) y quieres que
reciba tu configuración. Necesita su **perfil**: la receta que dice qué máquina es
y qué valor toma en cada dimensión.

Usaremos como ejemplo un portátil, `portatil`, con escritorio y para el trabajo.

## 1. Crea la capa de la máquina

En tu repo de contenido, cada máquina tiene su carpeta en `machines/`, aunque de
momento esté vacía:

```sh
cd ~/.local/share/dots/content
mkdir -p machines/portatil
touch machines/portatil/.gitkeep
```

> [!TIP]
> El `.gitkeep` está solo para que git guarde la carpeta vacía. Cuando metas algo
> dentro, puedes borrarlo.

## 2. Declara el perfil en `dots.yaml`

Añade una entrada en `profiles`:

```yaml
profiles:
  # ...los que ya tenías...
  portatil:
    description: "Portátil del trabajo"
    machine: portatil
    params: { ux: desktop, identity: work }
```

- `machine` es el nombre de la carpeta que acabas de crear.
- `params` lleva **exactamente** un valor por cada dimensión de tu `dots.yaml`: ni
  falta ninguna ni sobra ninguna. Si no tienes dimensiones, no pongas `params`.

Si algún valor es nuevo (por ejemplo, es tu primera máquina de trabajo y no
existía `identity/work`), crea también su carpeta:

```sh
mkdir -p identity/work && touch identity/work/.gitkeep
```

## 3. Comprueba que todo cuadra

`dots` valida el manifiesto cada vez que lo carga. Un `link` en modo prueba,
indicando el perfil nuevo, te dice si falta algo sin tocar nada:

```sh
dots link --dry-run --profile portatil
```

Si hay alguna pega (una carpeta que no existe, un `params` incompleto), te las
enseña todas de una vez.

## 4. Súbelo y estrena la máquina

```sh
git add dots.yaml machines/portatil identity/work
git commit -m "Añade el perfil del portátil"
git push
```

Y en el portátil, después de instalar `dots`:

```sh
dots init --remote git@github.com:<tu-usuario>/<tu-contenido>.git --profile portatil
```

A partir de ahí, lo que quieras solo en esa máquina va en `machines/portatil/`.

## Ver también

- [Perfiles y capas](../explanation/profiles-and-layers.md)
- [`profiles` y `dimensions` en el manifiesto](../reference/manifest.md#profiles)
- [Añadir un dotfile](add-a-dotfile.md)

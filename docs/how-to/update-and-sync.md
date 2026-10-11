# El día a día: actualizar y sincronizar

Has cambiado algo en una máquina y lo quieres en las demás. O al revés: llegas a
una máquina que lleva días sin tocar y quieres ponerla al día.

## En la máquina donde cambias algo

Tu contenido es un repo git normal. Cuando cambias un dotfile (o `dots add`
apunta un paquete), el cambio queda en el repo; súbelo como siempre:

```sh
cd ~/.local/share/dots/content
git status
git add .
git commit -m "Cambia los alias de git"
git push
```

> [!TIP]
> Como tus dotfiles son symlinks al repo, editar `~/.bashrc.d/aliases.sh` ya es
> editar el fichero del repo. Solo falta el commit.

## En las demás máquinas

```sh
dots update
```

Trae los cambios del repo y vuelve a aplicarlo todo, en este orden:

1. Enlaza los dotfiles nuevos y coloca los ficheros de sistema.
2. Ejecuta los hooks `post-link`.
3. Instala los paquetes que falten (y habilita sus repos COPR).
4. Ejecuta los hooks `post-packages`.
5. Trae los externals.

Todo es repetible: si algo ya está bien, no lo toca. Si quieres ver qué haría sin
hacer nada, añade `--dry-run`.

> [!WARNING]
> `dots update` hace un `git pull` que solo avanza (`--ff-only`). Si tienes
> commits locales sin subir y en el remoto hay otros, se para. Resuélvelo con git
> (un `git pull --rebase` en el contenido suele bastar) y repite.

## Ver cómo está una máquina

Para saber si falta algún paquete, sin instalar nada:

```sh
dots sync
```

Y para ver qué symlinks cambiarían:

```sh
dots link --dry-run
```

## Ver también

- [Añadir un dotfile](add-a-dotfile.md)
- [Instalar paquetes](manage-packages.md)
- [Comandos](../reference/commands/README.md)

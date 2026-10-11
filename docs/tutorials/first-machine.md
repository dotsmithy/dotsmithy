# Tu primera máquina con dotsmithy

En este tutorial vas a montar una máquina desde cero con dotsmithy. Al acabar
tendrás tu propio repo de contenido, `dots` instalado y una terminal que te saluda
diciéndote con qué perfil está montada. Y sabrás cómo encaja todo.

Calcula unos 15 minutos.

## Lo que necesitas

- Una **Fedora**. Para esta primera vez, mejor una máquina virtual o un equipo que
  no te importe trastear: vamos a instalar paquetes y a cambiar el nombre del
  equipo.
- Una cuenta de **GitHub**.
- Ganas de cacharrear un rato. 🙂

## 1. Crea tu repo de contenido

Tu contenido es el repo donde vivirán tus dotfiles. No hace falta montarlo a mano:
parte de la [plantilla](https://github.com/kreikol/dotsmithy-content-template).

Entra en la plantilla, pulsa **Use this template** → **Create a new repository**,
ponle nombre (usaremos `mis-dots`) y márcalo como **Private**: tus dotfiles son
cosa tuya.

> [!TIP]
> Si usas la CLI de GitHub, lo mismo en una línea:
>
> ```sh
> gh repo create mis-dots --template kreikol/dotsmithy-content-template --private
> ```

Echa un vistazo a lo que trae. Lo más importante es `dots.yaml`, que define dos
perfiles:

```yaml
profiles:
  laptop:
    description: "Portátil con escritorio"
    machine: laptop
    params: { ux: desktop }
  server:
    description: "Servidor o VM sin escritorio"
    machine: server
    params: { ux: terminal }
```

Cada carpeta (`shared/`, `ux/desktop/`, `machines/laptop/`...) es una **capa**. Un
perfil es una pila de capas: el `laptop` recibe lo de `shared`, lo de `ux/desktop`
y lo de `machines/laptop`. De momento no hace falta más: lo vas a ver funcionar.

## 2. Ponle nombre a tu máquina

La plantilla trae un regalo que cambia el nombre del equipo según la máquina. Para
el perfil `laptop`, el nombre está en `machines/laptop/system/hostname` y es
`mi-laptop`.

Si quieres otro, edita ese fichero en GitHub (o en tu clon) y pon el que prefieras.
Si no quieres que toque el nombre del equipo, bórralo junto con su entrada en
`dots.yaml`.

> [!WARNING]
> No es de mentira: al aplicar el perfil, el equipo pasa a llamarse como diga ese
> fichero.

## 3. Instala `dots`

Ya en tu Fedora, una línea:

```sh
curl -fsSL https://raw.githubusercontent.com/kreikol/dotsmithy/main/script/install.sh | bash
```

Instala el comando `dots` (y su atajo, `ds`). Compruébalo:

```sh
dots --version
```

## 4. Monta la máquina

Ahora le dices a `dots` de dónde sale tu contenido y qué perfil es esta máquina:

```sh
dots init --remote git@github.com:<tu-usuario>/mis-dots.git --profile laptop
```

Como tu repo es privado, `dots` lo clona por SSH. Si esta máquina aún no tiene
clave SSH, `dots` te crea una y te la enseña. Añádela en GitHub (**Settings** →
**SSH and GPG keys** → **New SSH key**), vuelve a lanzar el mismo comando y esta
vez clonará.

Mientras trabaja, te va contando lo que hace. Es esto, en orden:

1. **Clona tu contenido** en `~/.local/share/dots/content` y apunta el perfil, para
   que el resto de comandos ya lo sepan.
2. **Enlaza tus dotfiles**: cada fichero de `home/` de tus capas pasa a ser un
   symlink en tu `$HOME`.
3. **Coloca los ficheros de sistema**, como el nombre del equipo. Para eso usa
   `sudo`, así que puede pedirte la contraseña.
4. **Ejecuta los hooks** de después de enlazar: uno guarda tu perfil para la
   bienvenida y otro aplica el nombre del equipo.
5. **Instala los paquetes** de las listas de tus capas (`tmux`, `tree`,
   `wl-clipboard`...).
6. **Trae los recursos externos**, si declaras alguno (la plantilla no trae
   ninguno activo).

## 5. Mira lo que ha pasado

Abre una terminal nueva. Te recibe esto:

```
🔨 dotsmithy 1.0.0 · perfil laptop (máquina laptop, ux=desktop)
```

Esa línea la pone un dotfile de tu contenido, y el perfil lo ha guardado un hook.
Mira de dónde sale:

```console
$ ls -l ~/.bashrc.d/10-dotsmithy-welcome.sh
lrwxrwxrwx. 1 tu tu 79 ... /home/tu/.bashrc.d/10-dotsmithy-welcome.sh -> /home/tu/.local/share/dots/content/shared/home/.bashrc.d/10-dotsmithy-welcome.sh
```

Es un **symlink** al fichero de tu repo. Lo mismo con la config de tmux o los
alias. Si editas cualquiera de ellos, estás editando el repo directamente.

Y el equipo ya tiene su nombre:

```sh
hostnamectl hostname
```

## 6. Hazlo tuyo

Vamos a añadir un paquete. Por ejemplo, `ripgrep`:

```sh
dots add dnf ripgrep
```

`dots` lo apunta en la lista de tu máquina (`machines/laptop/packages/dnf.txt`) y lo
instala. Ahora guárdalo en tu repo:

```sh
cd ~/.local/share/dots/content
git add machines/laptop/packages/dnf.txt
git commit -m "Añade ripgrep"
git push
```

Ya está en GitHub. En cualquier otra máquina con el perfil `laptop`, un
`dots update` lo instalaría también.

## Lo que has aprendido

- Tu **contenido** es un repo tuyo, y la plantilla te da un punto de partida.
- Cada máquina tiene un **perfil**, que apila **capas**: lo común, lo de su tipo y
  lo suyo.
- Tus dotfiles son **symlinks** a tu repo: editarlos es editar el repo.
- `dots init` monta una máquina de cero; `dots add` y `dots update` te acompañan
  en el día a día.

## Y ahora qué

- Mete tus propios dotfiles: [añadir un dotfile](../how-to/add-a-dotfile.md).
- Prepara tu segunda máquina: [preparar una máquina nueva](../how-to/create-a-profile.md).
- Entiende por qué funciona así: [perfiles y capas](../explanation/profiles-and-layers.md).

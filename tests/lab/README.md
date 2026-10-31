# Lab

Entorno para trastear dotsmithy a mano en una Fedora limpia **con systemd** y el
contenido de muestra ya aplicado.

- `up.sh` lo monta todo: compila el motor (estático), construye la imagen
  (`Containerfile`), arranca el contenedor con systemd como PID 1 y aplica el
  contenido con `dots init` como el usuario `lab` (no root, sudo sin contraseña),
  igual que en una máquina real.
- Necesita **podman** en el host (`podman run --systemd=always`).
- `content/` es un repo de contenido de mentira, mínimo pero representativo: un
  dotfile (bajo `.config/` para no chocar con el `.bashrc` de la imagen base),
  una lista de paquetes dnf y un hook `post-link`. Se **copia** a
  `/home/lab/content` dentro de la imagen, escribible, para que funcionen los
  comandos que escriben en el contenido (`dots add`, `dots sync`…). El fixture
  del repo no se toca.
- El binario se compila a `tests/lab/bin/dots` (ignorado por git).
- Con `LAB_DOTS_VERSION` no se compila nada: el lab instala esa **release
  publicada** con el instalador (`script/install.sh`) como el usuario `lab`, igual
  que en una máquina real. Es la forma de probar una versión candidata. Se usa el
  instalador del propio checkout, que es el que saldrá con esa versión.

## Uso

En el workflow manual **Lab** (`.github/workflows/lab.yml`), `up.sh` corre con
`LAB_REMOTE=1` y `LAB_GITHUB_USER` = quien lo lanza: abre una sesión SSH remota
con [upterm](https://github.com/owenthereal/upterm) a la que **solo pueden entrar
las claves SSH que esa cuenta tiene en GitHub**. El comando `ssh` para entrar sale
en el log del job. El lab dura hasta 2 h; se cierra antes saliendo de la shell
(`exit`) o con `touch /tmp/close-lab`.

Para probar una release publicada (por ejemplo una candidata) en lugar de compilar,
lánzalo con el input `dots_version`, desde la pestaña Actions o con:

```bash
gh workflow run lab.yml -f dots_version=v1.0.0-rc.2
```

En el log del job verás la línea `versión de dots en el lab: …` con lo que se está
probando.

En local (desde el devbox del proyecto, que aporta Go; podman sale del sistema):

```bash
devbox run -- ./tests/lab/up.sh
podman exec -it -u lab -w /home/lab dotsmithy-lab bash -l   # entrar
podman rm -f dotsmithy-lab                                 # tirarlo

# con una release publicada en vez de compilar:
LAB_DOTS_VERSION=v1.0.0-rc.2 ./tests/lab/up.sh
```

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
- El binario se compila a `tests/lab/dots` (ignorado por git).

## Uso

En el workflow manual **Lab** (`.github/workflows/lab.yml`), `up.sh` corre con
`LAB_TMATE=1`: abre una sesión SSH por tmate y espera hasta 2 h, o hasta que
hagas `touch /tmp/close-lab` dentro.

En local (desde el devbox del proyecto, que aporta Go; podman sale del sistema):

```bash
devbox run -- ./tests/lab/up.sh
podman exec -it -u lab -w /home/lab dotsmithy-lab bash -l   # entrar
podman rm -f dotsmithy-lab                                 # tirarlo
```

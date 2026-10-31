# dotsmithy

**La herrería de tus dotfiles.** Tus ficheros de configuración, tus paquetes y los
ajustes de cada máquina, guardados en un repo tuyo y aplicados con un solo comando.

dotsmithy tiene dos piezas: un **motor** (este repo, el comando `dots`) que sabe
hacer el trabajo, y tu **contenido** (un repo privado con tus dotfiles de verdad)
que dice qué hay que hacer. Se entienden a través de un fichero, `dots.yaml`. Tú
editas tus ficheros como siempre y `dots` se encarga de que cada máquina tenga lo
suyo.

> [!NOTE]
> dotsmithy está pensado para **Fedora**. Funciona y está cubierto por tests, pero
> todavía no ha llegado a la v1.0, así que puede cambiar algún detalle.

## Pruébalo en un minuto

En una Fedora limpia, esta línea instala `dots` (y su atajo `ds`):

```sh
curl -fsSL https://raw.githubusercontent.com/dotsmithy/dotsmithy/main/script/install.sh | bash
```

Luego le dices de dónde sale tu contenido y qué perfil quieres para esta máquina:

```sh
dots init --remote git@github.com:<tu-usuario>/<tu-contenido>.git --profile <perfil>
```

¿Aún no tienes repo de contenido? Créalo desde la
[plantilla](https://github.com/dotsmithy/content-template): el
[tutorial](docs/tutorials/first-machine.md) te lleva de la mano.

## ¿Qué quieres hacer?

| | |
| --- | --- |
| 🧭 **Aprender desde cero** | Monta tu primera máquina paso a paso: [tutoriales](docs/tutorials/README.md) |
| 🛠️ **Resolver algo concreto** | Añadir un dotfile, un paquete, un hook...: [guías](docs/how-to/README.md) |
| 📖 **Consultar un detalle** | Comandos, `dots.yaml`, cómo se organiza el contenido: [referencia](docs/reference/README.md) |
| 💡 **Entender cómo funciona** | Por qué está hecho así: [explicaciones](docs/explanation/README.md) |

Todo junto, en el [índice de la documentación](docs/README.md).

## ¿Quieres meterle mano al código?

Mira [CONTRIBUTING.md](CONTRIBUTING.md).

## Licencia

[MIT](LICENSE).

# Cómo se trabaja en el repo

Esta guía cuenta cómo van las ramas, los pull requests y las issues en dotsmithy, y
por qué están así. Lo básico para abrir tu primer PR está en
[CONTRIBUTING.md](../CONTRIBUTING.md).

## Las ramas

| Rama | Qué contiene |
| --- | --- |
| `main` | Solo versiones publicadas, cada una con su tag (`v1.0.0`, `v1.1.0`...). Es lo que instala el instalador y lo que enseña la web de la documentación |
| `next` | La integración: aquí se juntan los cambios de la próxima versión. Es la rama por defecto |
| `feat/<slug>` | Tu trabajo: funcionalidades, arreglos, documentación, CI... |

Tu rama sale de **`main`** y tu PR va a **`next`**. Salir de `main` hace que tu rama
lleve solo tu cambio, sin arrastrar lo que hay en `next` y aún no se ha publicado.
Si necesitas algo que solo está en `next`, mergea `next` en tu rama
(`git merge origin/next`) cuando te haga falta.

El tipo de cambio no va en el nombre de la rama: lo dice el título del PR. Por eso
todo es `feat/`.

## Los pull requests

### Cómo se mergean

Los PRs a `next` se mergean con **squash**: todos los commits de tu rama se quedan en
uno. Así que no te preocupes por los commits intermedios ("wip", "arreglo el
test"...): lo que cuenta es el PR.

### El título

El título del PR es **el commit que queda** en la historia, y de él sale el
changelog de cada versión. Por eso tiene que ser un Conventional Commit, en inglés:

```text
tipo(ámbito opcional): descripción en minúscula
```

Por ejemplo, `feat(cli): add status command` o `docs: fix a typo in the manifest
reference`. Un check lo comprueba en cada PR. Tipos admitidos: `feat`, `fix`,
`docs`, `chore`, `test`, `refactor`, `ci`, `build`, `perf`, `style` y `revert`.

### La descripción

La descripción del PR, en castellano, es para quien revisa: qué cambia, por qué y
cómo lo has comprobado. **No** pasa a la historia (el commit del squash es solo el
título), así que puedes extenderte lo que haga falta.

Si el PR resuelve una issue, pon la palabra clave en la descripción: `Fix #N` para
un fallo, `Resolve #N` para una mejora. La issue se cierra sola al mergear en
`next`.

### Draft: cuando aún no está para revisar

Si tu PR todavía no está listo, ábrelo como **Draft** y cuenta en la descripción en
qué punto está. Un Draft no se puede mergear y no pide revisión. Cuando esté listo,
pulsa **Ready for review**.

Quien mantiene el repo puede ponerle un label para ordenarlo, por ejemplo `on-hold`
si está terminado pero va en una versión posterior.

### Los checks

Para mergear, el PR tiene que estar en verde en todos estos checks:

- `build + test`: compila, `go vet` y `go test`.
- `integración (Fedora)`: el binario real en una Fedora limpia.
- `shellcheck (shim)`: el instalador y los scripts del repo.
- `plantilla`: el `dots` de tu PR valida la plantilla de `template/`.
- `título Conventional Commit`: el título del PR.
- `origen de PR a main`: a `main` no se llega desde una rama `feat/`.

Cómo pasar los tests en local, antes de abrir el PR:
[el entorno de desarrollo](development.md).

Si es tu primer PR en el repo, los checks esperan a que alguien del proyecto les dé
permiso para arrancar.

Además, la rama tiene que estar **al día** con `next`. Si alguien mergea antes que
tú, GitHub te pedirá ponerla al día: si no hay conflictos, basta con el botón
**Update branch**; si los hay, se resuelven en tu rama.

### Labels de área

Un labeler pone solo labels según los ficheros que toca el PR: `documentation`,
`ci`, `installer` (el instalador) y `tests`. Sirven para filtrar; no hace falta que
los pongas tú.

## Las issues

Si encuentras un fallo o se te ocurre una mejora, abre una
[issue](https://github.com/dotsmithy/dotsmithy/issues). Cuéntalo con un ejemplo: qué
comando, qué esperabas y qué ha pasado.

Desde ahí se clasifica: entra en el **milestone** de la versión en la que va a salir
(`v1.0`, `v1.1`...) o se queda en `backlog` si aún no tiene versión. Las
funcionalidades grandes son **épicas**: issues con el label `epic` y sus
**sub-issues** colgadas.

## La plantilla de contenido

La plantilla que se usa para empezar un contenido propio
([dotsmithy/content-template](https://github.com/dotsmithy/content-template)) se
mantiene **aquí**, en la carpeta [`template/`](../template/). Ese otro repo es un
espejo: en cada versión publicada se sincroniza solo desde `template/`. No se edita
a mano y no admite issues.

Si tu cambio afecta al manifiesto o a cómo se organiza el contenido, actualiza
`template/` en el mismo PR. El check `plantilla` comprueba que el `dots` de tu PR la
entiende.

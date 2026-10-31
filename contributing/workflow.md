# Cómo se trabaja en el repo

Esta guía cuenta cómo van las ramas, los pull requests y las issues en dotsmithy.
Si vienes a contribuir, con leer la primera sección ya puedes abrir tu PR; el resto
es para entender por qué las cosas están como están.

## Lo básico, en corto

- Tu rama sale de **`next`** y tu PR va a **`next`** (es la rama por defecto, así
  que GitHub te la propone sola).
- Llama a tu rama **`feat/<algo-descriptivo>`**.
- El **título del PR** es un [Conventional Commit](https://www.conventionalcommits.org/es/)
  en inglés: `feat(cli): add status command`, `fix: ...`, `docs: ...`.
- Si aún no está listo para revisar, ábrelo como **Draft**.
- Si la plantilla de contenido tiene que cambiar, el cambio va en **`template/`**,
  en este mismo repo y en el mismo PR.

## Las ramas

| Rama | Qué contiene |
|---|---|
| `main` | Solo versiones publicadas. Cada una lleva su tag (`v1.0.0`, `v1.1.0`...). Es lo que instala el instalador y lo que enseña la web de la documentación |
| `next` | La integración: aquí se juntan los cambios de la próxima versión. Es la rama por defecto |
| `feat/<slug>` | El trabajo del día a día: funcionalidades, arreglos, documentación, CI... Sale de `next` y vuelve a `next` |
| `hotfix/<slug>` | Solo para arreglar algo urgente de una versión ya publicada. Sale de `main` y vuelve a `main` |

El tipo de cambio no va en el nombre de la rama: lo dice el título del PR. Por eso
casi todo es `feat/`.

Una vez publicada una versión, `next` sigue acumulando cambios hasta la siguiente.
Nada llega a `main` sin pasar antes por `next`, salvo los hotfixes.

## Los pull requests

### A dónde van y cómo se mergean

| PR | Cómo se mergea |
|---|---|
| `feat/*` → `next` | **Squash**: queda un solo commit por PR |
| `hotfix/*` → `main` | **Squash** |
| `next` → `main` (publicar una versión) | **Merge commit**, nunca squash |
| `main` → `next` (después de un hotfix) | **Merge commit**, nunca squash |

Con squash, todos los commits de tu rama se quedan en uno. Así que no te preocupes
por los commits intermedios ("wip", "arreglo el test"...): lo que cuenta es el PR.

Los PRs entre `next` y `main` van con merge commit para que las dos ramas compartan
la misma historia; con squash divergirían y en cada versión saldrían conflictos que
no son reales.

### El título

El título del PR es **el commit que queda** en la historia, y de él sale el
changelog de cada versión. Por eso tiene que ser un Conventional Commit, en inglés:

```text
tipo(ámbito opcional): descripción en minúscula
```

Un check lo comprueba en cada PR. Tipos admitidos: `feat`, `fix`, `docs`, `chore`,
`test`, `refactor`, `ci`, `build`, `perf`, `style`, `revert` y `release` (este
último solo para el PR que publica una versión).

### La descripción

La descripción del PR, en castellano, es para quien revisa: qué cambia, por qué y
cómo lo has comprobado. **No** pasa a la historia (el commit del squash es solo el
título), así que puedes extenderte lo que haga falta.

Si el PR resuelve una issue, pon la palabra clave en la descripción: `Fix #N` para
un fallo, `Resolve #N` para una mejora. La issue se cierra sola al mergear en
`next`.

### Draft: cuando aún no está para revisar

Un PR en **Draft** no se puede mergear y no pide revisión. Para que se sepa por qué
está así, lleva uno de estos labels:

| Label | Qué significa |
|---|---|
| `wip` | Se está trabajando en él |
| `on-hold` | Está terminado, pero va en una versión posterior: espera su turno |
| `discussion` | Se ha abierto para comentar un enfoque antes de seguir |

"Para revisar" es simplemente un PR abierto que **no** está en Draft. Cuando el tuyo
esté listo, quita el label y pulsa **Ready for review**.

Un PR `on-hold` se revisa cuando le llega su turno, sobre su versión final y al día
con `next`, no antes.

### Los checks

Para mergear, el PR tiene que estar en verde en todos estos checks:

- `build + test`: compila, `go vet` y `go test`.
- `integración (Fedora)`: el binario real en una Fedora limpia.
- `shellcheck (shim)`: el instalador y los scripts del repo.
- `plantilla`: el `dots` de tu PR valida la plantilla de `template/`.
- `título Conventional Commit`: el título del PR.
- `origen de PR a main`: a `main` solo se llega desde `next` o desde un `hotfix/*`.

Además, la rama tiene que estar **al día** con su destino. Si alguien mergea antes
que tú, GitHub te pedirá ponerla al día: si no hay conflictos, basta con el botón
**Update branch**; si los hay, se resuelven en tu rama.

### Labels de área

Un labeler pone solo labels según los ficheros que toca el PR: `documentation`,
`ci`, `installer` (el instalador) y `tests`. Sirven para filtrar; no hace falta que
los pongas tú y nunca quita los que pongas a mano.

## Las issues y el tablero

- **Milestone = versión.** Cada issue (y cada PR) planificada lleva el milestone de
  la versión en la que sale: `v1.0`, `v1.1`... La barra del milestone dice cuánto
  falta para publicarla. Los parches (`v1.0.1`) no tienen milestone propio.
- **Épica = funcionalidad grande.** Es una issue con el label `epic` y sus
  **sub-issues** colgadas (lo verás en un bloque dentro de la propia issue, con su
  progreso). La épica no lleva milestone; cada sub-issue lleva el de la versión en la
  que sale, así que una funcionalidad puede repartirse entre varias versiones.
- **`triage`**: issue nueva que aún no se ha clasificado.
- **`backlog`**: decidida, pero sin versión asignada todavía.

Si encuentras un fallo o se te ocurre una mejora, abre una
[issue](https://github.com/dotsmithy/dotsmithy/issues): se clasifica desde ahí.

## La plantilla de contenido

La plantilla que se usa para empezar un contenido propio
([dotsmithy/content-template](https://github.com/dotsmithy/content-template)) se
mantiene **aquí**, en la carpeta [`template/`](../template/). Ese otro repo es un
espejo: en cada versión publicada se sincroniza solo desde `template/`, con el
mismo tag. No se edita a mano y no admite issues.

Si tu cambio afecta al manifiesto o a cómo se organiza el contenido, actualiza
`template/` en el mismo PR. El check `plantilla` comprueba que el `dots` de tu PR la
entiende.

## Publicar una versión

Lo cuenta [RELEASING.md](../RELEASING.md).

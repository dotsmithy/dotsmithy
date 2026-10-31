# Publicar una versión

Una versión sale de juntar en `next` todo lo de su milestone y llevarlo a `main` con
un tag. Al empujar el tag, todo lo demás se publica solo:

- **Binarios**: goreleaser (`.github/workflows/release.yml`) cross-compila (linux
  amd64 + arm64), arma los `.tar.gz` con el binario `dots`, genera checksums y crea
  la **GitHub Release** con un changelog agrupado por tipo de cambio. El instalador
  (`script/install.sh`) baja de ahí.
- **Plantilla de contenido**: `.github/workflows/template-sync.yml` copia
  `template/` a [dotsmithy/content-template](https://github.com/dotsmithy/content-template)
  y le pone el mismo tag. Solo en versiones finales, no en las candidatas.
- **Web de la documentación**: se publica desde `main`
  (`.github/workflows/docs.yml`), así que enseña siempre lo publicado.

Cómo se trabajan las ramas y los PRs: [contributing/flujo-de-trabajo.md](contributing/flujo-de-trabajo.md).

## El tag

Tiene que ser **semver válido** (goreleaser lo exige): `vX.Y.Z`, con sufijo de
pre-release opcional para las candidatas (`-rc.1`, `-rc.2`...).

## La web de la documentación

Se publica desde `main` con GitHub Pages (origen **GitHub Actions**). El entorno
`github-pages` del repo solo admite despliegues desde `main`. Ojo: si alguien
desactiva y vuelve a crear Pages, GitHub deja el entorno restringido a la rama por
defecto (`next`) y hay que volver a cambiarlo a `main` (Settings → Environments →
github-pages).

## Candidatas (opcional): `vX.Y.Z-rc.N` sobre `next`

Para probar una versión en máquinas reales antes de publicarla:

```sh
git switch next && git pull
git tag -s v1.1.0-rc.1 -m "v1.1.0-rc.1"
git push origin v1.1.0-rc.1
```

Sale como **pre-release**: GitHub no la marca como la última, así que el instalador
no la usa por defecto. Para instalarla:

```sh
curl -fsSL https://raw.githubusercontent.com/dotsmithy/dotsmithy/main/script/install.sh | DOTS_VERSION=v1.1.0-rc.1 bash
```

Lo que salga de las pruebas se arregla con PRs normales a `next` y, si hace falta,
otra candidata (`-rc.2`).

## Publicar la versión

1. **El milestone está completo** y `next` en verde. Lo que no vaya a llegar se
   pasa al milestone siguiente (o a `backlog`), no se queda colgando.
2. **PR de `next` a `main`**, titulado `release: vX.Y.Z`. Pasa el CI sobre el
   conjunto de la versión.
3. **Mergear con "Create a merge commit"**, nunca con squash: con squash, `main` y
   `next` dejarían de compartir historia.
4. **Tag firmado sobre `main`**:

   ```sh
   git switch main && git pull
   git tag -s v1.1.0 -m "v1.1.0"
   git push origin v1.1.0
   ```

5. **Revisar la Release** cuando termine el workflow. Si la versión lo merece,
   añadir arriba a mano un párrafo con lo destacado; el changelog de debajo sale
   solo.
6. **Cerrar el milestone.**

El changelog de una versión final se calcula desde la versión final anterior,
saltándose las candidatas, así que cuenta todo lo que trae la versión.

## Hotfix: arreglar algo urgente de lo publicado

Solo si el fallo no puede esperar a la siguiente versión. Si puede esperar, es un
PR normal a `next`.

1. Rama `hotfix/<slug>` desde `main`.
2. PR a **`main`** (ojo: GitHub propone `next`, hay que cambiar la base). Se mergea
   con **squash**.
3. Tag de parche firmado sobre `main` (`v1.1.1`), como en el paso 4 de arriba.
4. **PR de `main` a `next`**, titulado por ejemplo
   `chore: sync main into next after v1.1.1`, mergeado con **merge commit**. Así el
   arreglo llega también a `next` y la siguiente versión no da conflictos. De
   momento este PR se abre a mano.

## Comprobar el empaquetado sin publicar

En local (`goreleaser release --snapshot --clean`, dentro del devbox) o en GitHub
lanzando a mano el workflow **Release** (`gh workflow run release.yml`): hace lo
mismo en modo snapshot y no publica nada.

La sincronización de la plantilla también se puede ensayar a mano, sin publicar:
workflow **Publicar plantilla** con el tag de una versión final ya publicada y
`dry_run` activado.

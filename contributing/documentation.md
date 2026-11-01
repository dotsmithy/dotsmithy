# Escribir documentación

Esta guía es para ti si vas a escribir o cambiar documentación de dotsmithy: una
página nueva de uso, un arreglo en una guía o una parte de cómo contribuir. Cuenta
dónde va cada cosa, cómo se escribe y qué hay que tocar para que se vea bien en
GitHub y en la web.

Si es un cambio pequeño (una errata, una frase que no se entiende), no hace falta
que te la leas entera: con el [tono](#tono-y-estilo) y la
[checklist](#checklist-antes-de-abrir-la-pr) vas bien.

## Dónde va cada cosa

| Qué | Dónde | Para quién |
| --- | --- | --- |
| Cómo se usa dotsmithy | `docs/` | quien instala y usa `dots` |
| Cómo se trabaja en el motor | `contributing/` | quien quiere tocar el código o la doc |
| La puerta de entrada | `README.md` | quien llega al repo por primera vez |
| El aviso de cómo contribuir | `CONTRIBUTING.md` | GitHub lo enlaza en issues y PRs |
| Cómo se publica una versión | `RELEASING.md` | quien saca una release |

`contributing/` va fuera de `docs/` a propósito: no es documentación de uso, y así
no se cuela en la web ni en el índice de quien solo quiere usar `dots`.

Los **nombres de ficheros y carpetas van siempre en inglés** (`add-a-dotfile.md`,
`profiles-and-layers.md`), aunque el contenido esté en castellano. Toda la
documentación se escribe **en castellano**, el `README.md` incluido.

Las ramas, los commits y las pull requests se cuentan en el
[flujo de trabajo](workflow.md); aquí solo lo que es propio de la documentación.

## Los cuatro cuadrantes

La documentación de uso sigue [Diátaxis](https://diataxis.fr/): cuatro tipos de
página, cada uno para un momento distinto de quien lee. Las carpetas llevan los
nombres canónicos: `tutorials/`, `how-to/`, `reference/` y `explanation/`.

### Tutoriales: aprender haciendo

Para quien llega de cero. Es un camino guiado, de principio a fin, que siempre
funciona si se siguen los pasos.

- Se apoya en la plantilla de contenido (`template/`), así que lo que se enseña
  tiene que existir de verdad en ella.
- Sin alternativas ni desvíos: un solo camino. Las variantes van en las guías.
- Explica lo justo para seguir adelante; el porqué largo, a las explicaciones.
- Estructura: qué vas a conseguir y cuánto se tarda, `## Lo que necesitas`, pasos
  numerados (`## 1. Crea tu repo de contenido`), `## Lo que has aprendido` y
  `## Y ahora qué`.

Ejemplo: [Tu primera máquina con dotsmithy](../docs/tutorials/first-machine.md).

### Guías: resolver una tarea concreta

Para quien ya sabe lo básico y quiere hacer algo en particular.

- El título es la tarea: "Añadir un dotfile", "Escribir un hook".
- Empieza diciendo en una o dos frases para qué sirve y cuándo la necesitas.
- Pasos numerados (`## 1. Elige la capa`) cuando hay un orden; secciones por caso
  (`## Un repo git`, `## Un fichero suelto`) cuando son recetas sueltas.
- Da por hecho lo básico y enlaza en vez de explicarlo otra vez.
- Termina en `## Ver también`.

Ejemplo: [Añadir un dotfile](../docs/how-to/add-a-dotfile.md).

### Referencia: el dato exacto

Para consultar, no para leer de corrido.

- Seca, precisa y completa: cada campo, cada opción, cada valor posible.
- Ordenada como la cosa que describe (los campos del manifiesto en su orden, las
  carpetas de una capa una por una).
- Nada de instrucciones ni de opiniones: eso va en guías y explicaciones.
- La de comandos no se escribe a mano: se genera (ver
  [la referencia de comandos](#la-referencia-de-comandos)).

Ejemplo: [El manifiesto: `dots.yaml`](../docs/reference/manifest.md).

### Explicaciones: el porqué

Para leer con calma cuando quieres entender cómo funciona y por qué es así.

- Prosa tranquila, sin pasos.
- Aquí sí caben opiniones, alternativas descartadas e historia: por qué symlinks y
  no plantillas, por qué el perfil se elige y no se adivina.
- Parte de un problema reconocible ("casi nadie tiene todas sus máquinas iguales")
  antes de contar la solución.
- Termina en `## Ver también`.

Ejemplo: [Perfiles y capas](../docs/explanation/profiles-and-layers.md).

### ¿Dónde va esto?

| Si lo que escribes... | Va en... |
| --- | --- |
| lleva a alguien nuevo de la mano hasta un resultado | tutoriales |
| responde a "¿cómo hago X?" | guías |
| responde a "¿qué hace exactamente X?" o "¿qué valores admite?" | referencia |
| responde a "¿por qué es así?" | explicaciones |

Si una página empieza a mezclar dos (una guía que se para a explicar el porqué,
una referencia que da instrucciones), pártela y enlaza una con otra.

## Tono y estilo

- **Tuteo**, cercano y directo, como si se lo explicaras a una amiga que sabe de
  sistemas. Nada de "el usuario deberá proceder a...".
- **Lenguaje neutro**: no supongas el género de quien lee. "Si es tu primera vez"
  mejor que "si eres nuevo".
- **Frases cortas** y un ejemplo concreto mejor que una abstracción: una ruta de
  verdad (`~/.config/starship.toml`) antes que "un fichero de configuración".
- **Sin guiones largos** para incisos: usa dos puntos, comas o paréntesis, o
  reformula la frase.
- Los **avisos y errores de `dots`** siguen el mismo tono: si la doc cita un
  mensaje del motor, que suene igual.
- **Comandos, rutas, ficheros y campos** siempre en código: `dots init`,
  `dots.yaml`, `home/`, `dimensions`.

### Glosario

Usa siempre estas palabras para estas cosas, y solo estas. Que la misma cosa se
llame de tres formas en tres páginas despista más que cualquier explicación mala.

| Término | Qué es | No lo llames |
| --- | --- | --- |
| **motor** | el comando `dots` y este repo: lo genérico, igual para todo el mundo | herramienta, framework, gestor |
| **contenido** | el repo de cada persona con sus dotfiles, paquetes y hooks | repo de dotfiles, config, datos |
| **manifiesto** | `dots.yaml`, el único contrato entre motor y contenido | config, fichero de configuración |
| **capa** | una carpeta del contenido que se apila: `shared`, una por valor de dimensión, `machines/<máquina>` | overlay, nivel, módulo |
| **dimensión** | un eje en el que se diferencian las máquinas (`ux`, `identity`); cada valor es una capa | categoría, tipo, grupo |
| **máquina** | un equipo concreto; tiene su propia capa, la última de la pila | host, nodo |
| **perfil** | la receta de una pila: qué máquina es y qué valor toma en cada dimensión | configuración, entorno |
| **hook** | un script del contenido que `dots` ejecuta en un momento fijo (`post-link`, `post-packages`) | script de instalación, trigger |
| **recurso externo** | un repo git o un fichero que no vive en el contenido y `dots externals` trae | dependencia, submódulo |
| **plantilla** | el punto de partida para crear un contenido nuevo | ejemplo, demo, fork |

## Formato

- **Un `README.md` por carpeta**, que hace de índice: en GitHub se ve al abrir la
  carpeta y en la web es la página de la sección.
- **Enlaces relativos** a ficheros `.md` dentro de la documentación
  (`../reference/manifest.md`). Funcionan igual en GitHub y en la web. Lo que está
  fuera del repo (otra web, la plantilla publicada) va con su URL completa.
- **Desde `docs/` no se enlaza fuera de `docs/`**: la web solo publica esa carpeta
  y el enlace quedaría roto. Si hace falta, URL completa al fichero en GitHub.
- **Avisos** con la sintaxis de GitHub, que la web también entiende:

  ```markdown
  > [!TIP]
  > Un atajo o un truco que ahorra tiempo.
  ```

  `NOTE` para un dato a tener en cuenta, `TIP` para un atajo, `WARNING` para algo
  que puede salir mal o tocar lo que no esperas (cambiar el hostname, sobrescribir
  un fichero). Con moderación: si todo es un aviso, nada lo es.
- **Bloques de código siempre con lenguaje** (`sh`, `yaml`, `toml`...). Si el
  bloque es solo para copiar y pegar, los comandos van sin `$` delante; el `$` se
  usa solo cuando el bloque mezcla un comando con su salida, para distinguirlos.
- **Títulos con tilde sin miedo**: las anclas (`#qué-pasa-cuando-dos-capas-dicen-lo-mismo`)
  salen igual en GitHub y en la web.
- **Emojis solo en los encabezados de los índices** (`docs/README.md`), como guía
  visual. En el cuerpo de las páginas, como mucho alguno suelto y con intención.

## La referencia de comandos

Las páginas de `docs/reference/commands/` se generan desde el propio CLI con
`tools/docgen`: la descripción, el uso, los ejemplos y las opciones salen del
código de cada comando. **No se editan a mano**; llevan un aviso arriba.

Si cambias un comando (una opción nueva, otro texto de ayuda, un comando nuevo):

```sh
devbox run docs
```

y commitea lo que cambie junto con el código. El CI comprueba que la referencia
está al día y falla si no lo está.

Para que salga bien, escribe la ayuda del comando ya pensando en la
documentación: la descripción corta (`Short`) en minúscula y sin punto final,
porque va en la tabla del índice, y los ejemplos en el campo `Example`, no dentro
de la descripción larga, para que salgan como un bloque de código.

## La web

La documentación se publica también como web en GitHub Pages
([dotsmithy.github.io/dotsmithy](https://dotsmithy.github.io/dotsmithy/)), hecha con
[zensical](https://zensical.org/) a partir de `docs/`. La fuente es la misma: lo
que se lee en GitHub es lo que sale en la web, y `zensical.toml` solo dice cómo se
presenta. Se publica desde `main`, así que la web siempre enseña la última versión
publicada.

Si añades una página nueva, enlázala en tres sitios:

1. El `README.md` de su cuadrante.
2. El índice general, `docs/README.md`.
3. El `nav` de `zensical.toml`, en su sección.

La excepción es la referencia de comandos: el `nav` solo enlaza su índice, así
que un comando nuevo no obliga a tocar `zensical.toml`.

Para verla en local mientras escribes, y para comprobarla como lo hace el CI
(modo estricto: un enlace roto es un error, no un aviso):

```sh
devbox run site-serve   # en http://localhost:8000, se recarga sola
devbox run site-build   # la misma comprobación que el CI
```

## Checklist antes de abrir la PR

- [ ] Cada página está en su cuadrante y hace una sola cosa.
- [ ] Ficheros y carpetas con nombre en inglés; el texto, en castellano.
- [ ] Tuteo, lenguaje neutro, sin guiones largos y con los términos del glosario.
- [ ] Enlaces relativos dentro de `docs/`, ninguno que salga de `docs/`.
- [ ] Bloques de código con lenguaje, y sin `$` si son para copiar.
- [ ] La página nueva está en el `README.md` de su cuadrante, en `docs/README.md` y
      en el `nav` de `zensical.toml`.
- [ ] Si has tocado un comando, has pasado `devbox run docs`.
- [ ] `devbox run site-build` pasa sin errores.

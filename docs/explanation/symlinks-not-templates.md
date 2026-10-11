# Symlinks y no plantillas

Con dotsmithy, el fichero que editas es el fichero de verdad. Tu `~/.bashrc.d/aliases.sh`
es un symlink al que tienes en tu repo de contenido: lo abres, lo cambias, y el
cambio ya está en el repo listo para un commit. No hay ningún paso intermedio.

Parece lo obvio, pero no todos los gestores de dotfiles lo hacen así.

## La alternativa: plantillas

Otra forma muy extendida es guardar en el repo **plantillas** y que la herramienta
genere los ficheros finales. Algo así:

```
# en el repo: .gitconfig.tmpl
[user]
    email = {{ if eq .machine "trabajo" }}yo@empresa.com{{ else }}yo@casa.net{{ end }}
```

Es potente: un único fichero sirve para todas las máquinas. Pero tiene un precio
que se paga todos los días:

- **El fichero que usas no es el que guardas.** Si cambias tu `~/.gitconfig` a
  mano, ese cambio no está en el repo. Tienes que acordarte de editar la plantilla
  y regenerar, o la próxima vez se pisa.
- **Las plantillas se leen mal.** Un fichero de configuración lleno de condiciones
  deja de parecerse al fichero que es. Y tu editor ya no lo reconoce bien.
- **Hay que aprender un lenguaje más** solo para escribir tu configuración.

## La opción de dotsmithy: un fichero, una verdad

En dotsmithy no hay plantillas. Cada fichero existe tal cual en tu repo y se enlaza
con un symlink en tu `$HOME`, como hace desde hace años GNU Stow. Eso significa:

- **Lo que ves es lo que hay.** El fichero del repo y el que usa tu programa son el
  mismo.
- **Editar es guardar.** Cualquier cambio, lo hagas tú o el propio programa, ya
  está en el repo. `git diff` te dice qué ha cambiado.
- **Nada que regenerar.** No hay un paso de "aplicar" que se te pueda olvidar
  después de cada cambio.

## ¿Y si un fichero tiene que ser distinto en cada máquina?

Es la pregunta que resuelven las plantillas, y aquí se resuelve con
[capas](profiles-and-layers.md): en vez de un fichero con condiciones, tienes
varias versiones, cada una en la capa que le toca. El `.gitconfig` con el correo
del trabajo vive en la capa `identity/work`, y el de casa, en `identity/personal`.
Cada máquina enlaza el suyo.

Muchas herramientas, además, permiten partir la configuración en trozos
(`~/.bashrc.d/`, `includeIf` en git, `conf.d/` en tantos programas). Con eso, la
parte común va en `shared` y solo el trozo que cambia va en otra capa. Tienes lo
mejor de las dos cosas sin escribir una sola condición.

## Lo que no encaja con symlinks

Hay casos en los que el modelo no sirve tal cual:

- **Ficheros del sistema** (`/etc/...`): no siempre se pueden o se deben enlazar.
  Para eso está la sección `system` del manifiesto, que los copia o enlaza con
  cuidado y los valida antes de reemplazarlos.
- **Cosas que no escribes tú**, como un plugin de terceros: no tiene sentido
  guardarlas en tu repo. Para eso están los externals, que el motor descarga y
  mantiene al día.
- **Programas que reescriben su propia configuración**: con un symlink, escribirían
  directamente en tu repo. Hoy se resuelve con un hook en tu contenido.

## Ver también

- [Perfiles y capas](profiles-and-layers.md)
- [Añadir un dotfile](../how-to/add-a-dotfile.md)
- [El manifiesto: `system` y `externals`](../reference/manifest.md#layers)

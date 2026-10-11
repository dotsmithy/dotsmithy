# El motor y el contenido

dotsmithy son dos repos que no se parecen en nada: el **motor**, que es este, y tu
**contenido**, que es tuyo. Esta página cuenta por qué están separados y qué pone
cada uno.

## Dos piezas, dos trabajos

El **motor** es el comando `dots`. Sabe hacer cosas: enlazar ficheros, instalar
paquetes, ejecutar scripts en orden, colocar ficheros en `/etc` sin romper nada,
clonar repos. Es igual para todo el mundo, está versionado y tiene tests. No sabe
nada de ti.

El **contenido** es tu repo de dotfiles: tu `.bashrc`, tu config del editor, las
listas de paquetes que quieres, los scripts que te hacen la vida fácil. Es solo
tuyo, normalmente privado, y no tiene ni una línea de lógica del motor.

Si lo piensas como una librería y una aplicación: el motor es la librería, que
sirve para cualquiera, y tu contenido es la aplicación que la usa.

## El manifiesto, el único punto de contacto

Entre los dos solo hay un fichero: `dots.yaml`, en la raíz de tu contenido. Dice
cómo son tus perfiles y lo poco que no se puede deducir mirando las carpetas
(repos COPR, recursos externos, ficheros de sistema).

Todo lo demás se saca de la estructura de carpetas. No le dices al motor "enlaza
este fichero": lo pones en `home/` de la capa que toca y ya está. Por eso el
manifiesto se queda corto y casi no cambia, aunque tu contenido crezca mucho.

Que el contrato sea un solo fichero tiene una ventaja que se nota con el tiempo:
el motor puede mejorar por dentro sin que tengas que tocar tu contenido, y tú
puedes reorganizar tus dotfiles sin esperar nada del motor.

## Por qué no todo en un repo

Mucha gente tiene un único repo de dotfiles con sus scripts de instalación
dentro, y funciona. Pero con el tiempo esos scripts crecen: que si comprobar si
ya está instalado, que si no pisar un fichero que ya existe, que si esta máquina
es distinta... Acabas manteniendo un pequeño programa sin tests mezclado con tus
ficheros de configuración.

Separarlo tiene unas cuantas ventajas:

- **Lo difícil se hace una vez y bien.** No pisar un fichero tuyo, no dejar un
  `/etc` roto, poder repetir todo sin miedo: eso vive en el motor, con tests, y no
  lo reescribes en cada repo.
- **Tu contenido se queda limpio.** Ficheros de verdad y poco más. Se lee de un
  vistazo.
- **Lo público y lo privado no se mezclan.** El motor es público; tu contenido no
  tiene por qué serlo.
- **Se puede actualizar cada uno a su ritmo.** Una versión nueva de `dots` no te
  obliga a cambiar tus ficheros.

## Dónde está la frontera

Una regla sencilla para saber dónde va cada cosa: lo **genérico** (lo que le
serviría a cualquiera) va en el motor; lo **tuyo** va en el contenido.

Tus scripts siguen teniendo sitio: los [hooks](../reference/content-layout.md#hooks-scripts-en-el-momento-justo)
son el lugar para lo que es solo tuyo, como configurar un servicio a tu manera o
preparar algo de una herramienta concreta. Si te encuentras escribiendo en un
hook algo que le serviría a cualquiera, probablemente pertenezca al motor.

## Ver también

- [Perfiles y capas](profiles-and-layers.md)
- [Symlinks y no plantillas](symlinks-not-templates.md)
- [El manifiesto: `dots.yaml`](../reference/manifest.md)

# Perfiles y capas

Casi nadie tiene todas sus máquinas iguales. El portátil del trabajo lleva la
identidad de la empresa y la VPN, el sobremesa de casa tiene escritorio y juegos,
y el servidor que tienes por ahí no tiene ni pantalla. Pero las tres comparten un
montón de cosas: tu shell, tu editor, tus alias de siempre.

Las **capas** son la forma que tiene dotsmithy de expresar eso: lo común se escribe
una vez y cada máquina recibe además solo lo suyo.

## El problema de copiar

Lo primero que se le ocurre a cualquiera es una carpeta por máquina. Funciona el
primer día. El segundo cambias un alias, lo cambias en una carpeta, se te olvida en
las otras dos, y a la semana tienes tres versiones de tu `.bashrc` que ya no sabes
cuál es la buena.

Otras salidas tienen pegas parecidas: una rama de git por máquina acaba igual de
desperdigada, y generar los ficheros desde plantillas con condiciones te obliga a
editar algo que no es tu fichero de verdad (de eso va
[symlinks y no plantillas](symlinks-not-templates.md)).

## Una pila de capas

En dotsmithy, cada máquina se monta apilando capas, siempre en el mismo orden:

```
shared  →  una capa por cada dimensión  →  machines/<máquina>
```

- **`shared`** tiene lo que va en todas partes.
- **Las dimensiones** son los ejes en los que se diferencian tus máquinas, y te las
  inventas tú. Por ejemplo, `ux` (con escritorio o solo terminal) e `identity`
  (personal o trabajo). Cada máquina elige un valor en cada eje, y cada valor es una
  capa: `ux/desktop`, `identity/work`...
- **`machines/<máquina>`** tiene lo que es solo de esa máquina, y siempre va la
  última.

Un **perfil** no es más que la receta de una pila: qué máquina es y qué valor toma
en cada dimensión. El portátil del trabajo sería `shared + ux/desktop +
identity/work + machines/portatil`, y el servidor, `shared + ux/terminal +
identity/personal + machines/servidor`. Ningún fichero se copia: cada uno vive en
una sola capa y llega a todas las máquinas que la usan.

## Las dimensiones son tuyas

El motor no sabe qué significa `ux` ni `identity`. Para él son nombres de carpeta
que tiene que apilar en orden. Eso es a propósito: tus máquinas se diferencian en
lo que se diferencian, y no tendría sentido que el motor te impusiera unas
categorías. Si mañana necesitas un eje nuevo (por ejemplo, `gpu` para las que
tienen gráfica dedicada), lo añades a tu `dots.yaml` y listo.

Y si tus máquinas no se diferencian en nada más que en la propia máquina, no
declaras ninguna dimensión: la pila se queda en `shared + machines/<máquina>`.

## Qué pasa cuando dos capas dicen lo mismo

Depende de qué sea:

- **Dotfiles**: si el mismo fichero está en dos capas, gana la que va más tarde en la
  pila. Por eso la máquina va la última: es la que tiene la última palabra. Aun
  así, lo sano es que cada fichero viva en una sola capa.
- **Paquetes, repos COPR, recursos externos y ficheros de sistema**: se suman. Una
  máquina instala los paquetes de todas sus capas.
- **Hooks**: se juntan los de todas las capas y se ejecutan por el nombre del
  fichero, así que el número del principio (`10-`, `20-`...) decide el orden,
  vengan de la capa que vengan.

Lo que no es de la pila no cuenta. Lo que pongas en `ux/desktop` no llega nunca a
una máquina sin escritorio, aunque esté en el mismo repo.

## El perfil se elige, no se adivina

Cuando montas una máquina le dices qué perfil es (`dots init --profile ...`) y
`dots` se lo apunta. No intenta deducirlo del nombre del equipo ni del hardware:
esas pistas fallan justo cuando menos te lo esperas (un hostname que cambia, una
máquina nueva con el mismo modelo), y equivocarse de perfil es instalar lo que no
toca. Decirlo una vez es más sencillo y más fiable.

Por la misma prudencia, cuando registras un paquete con `dots add` va a la capa de
la máquina salvo que digas otra cosa: así no aparece por sorpresa en todas las
demás.

## Ver también

- [El manifiesto: `dots.yaml`](../reference/manifest.md), con los campos de
  `dimensions` y `profiles`.
- [Preparar una máquina nueva](../how-to/create-a-profile.md)
- [El motor y el contenido](engine-and-content.md)

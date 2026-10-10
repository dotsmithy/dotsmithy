// Paquete manifest: parseo y validación del manifiesto dots.yaml, el único
// contrato entre el motor y el contenido.
//
// Responsabilidad: leer dots.yaml, validarlo y exponer un modelo tipado
// (perfiles, dimensiones, paquetes, hooks y lo que declara cada capa) al resto
// del motor. Es lógica pura y sin efectos: no toca el disco salvo por Load, que
// lee el fichero y comprueba que existen las carpetas de las capas
// (ValidateLayout).
//
// También vive aquí la resolución de la pila de capas de un perfil (Stack:
// shared, una capa por dimensión y la máquina) y la unión de lo que declaran
// esas capas (COPR, Externals, SystemFiles), por ser operaciones puramente
// derivadas del manifiesto.
package manifest

// Paquete hooks: ejecución de los hooks del contenido en los momentos clave
// del ciclo de vida (ADR 0007).
//
// Responsabilidad: dado un punto (post-link/post-packages/post-init), juntar
// los scripts <capa>/hooks/<punto>/*.sh de todas las capas de la pila del perfil
// y ejecutarlos en orden léxico del nombre, con bash, con cwd = raíz del
// contenido, pasando el contexto por variables de entorno (DOTS_*), y política
// fail-fast (si un hook falla, se aborta nombrándolo). Los hooks son
// responsabilidad del contenido: deben ser idempotentes y no interactivos.
package hooks

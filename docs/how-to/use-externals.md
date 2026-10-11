# Traer recursos externos

Hay cosas que quieres en tus máquinas pero que no escribes tú: un plugin de un
tercero, un tema, un repo tuyo que tiene vida propia. No tiene sentido copiarlas a
tu contenido. Con un **external**, declaras de dónde salen y `dots` las descarga y
las mantiene al día.

## Un repo git

Por ejemplo, el gestor de plugins de tmux. En `dots.yaml`, dentro de la capa que lo
quiera:

```yaml
layers:
  shared:
    externals:
      - dest: ~/.tmux/plugins/tpm
        type: git
        repo: https://github.com/tmux-plugins/tpm
        ref: master
```

`dots` clona el repo en `dest` y, en cada actualización, trae lo nuevo de `ref`
(una rama, un tag o un commit).

## Un fichero suelto

```yaml
      - dest: ~/.local/bin/mi-script
        type: file
        url: https://example.com/mi-script
```

## Hacer algo cuando cambia

Si al traerlo hay que preparar algo (compilar, instalar dependencias, dar permisos),
añade un script `post`, con la ruta relativa a la carpeta de la capa:

```yaml
      - dest: ~/.local/bin/mi-script
        type: file
        url: https://example.com/mi-script
        post: scripts/mi-script-post.sh
```

`shared/scripts/mi-script-post.sh`:

```bash
#!/usr/bin/env bash
chmod +x "$DOTS_EXTERNAL_DEST"
```

Solo se ejecuta cuando el recurso es nuevo o ha cambiado, no en cada actualización.
Recibe las mismas variables que un hook, más `DOTS_EXTERNAL_DEST` con el destino.

## Traerlos

Se traen solos al final de `dots init` y de `dots update`. Si solo quieres
actualizar los externals:

```sh
dots externals
```

Si uno falla (el repo no responde, por ejemplo), `dots` te avisa y sigue con los
demás.

## Ver también

- [Campos de `externals` en el manifiesto](../reference/manifest.md#layers)

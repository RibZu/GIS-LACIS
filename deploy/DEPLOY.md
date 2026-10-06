# Publicar GIS-LACIS en Railway

Esta carpeta es una copia depurada del proyecto: tiene solo lo necesario para publicarlo. Seguí los pasos en orden. Tiempo estimado: 30 a 45 minutos la primera vez.

## 1. Qué vas a crear

Un proyecto de Railway con tres piezas:

1. **Servicio web**: el sitio, construido a partir de este código con el `Dockerfile`.
2. **Base de datos PostgreSQL**: la administra Railway. Vas a cargarle `db/lacis.sql`.
3. **Volumen**: un disco que conserva las imágenes, CVs y PDFs que se suben desde el panel, aunque el sitio se reinicie o se vuelva a publicar.

**Costo aproximado** (precios de Railway al 2026-10-05; confirmalos en [docs.railway.com/pricing/plans](https://docs.railway.com/pricing/plans)):

- Plan **Hobby**: USD 5 por mes, que ya incluyen USD 5 de consumo.
- Se cobra lo que usa el proyecto encendido: memoria (USD 10 por GB por mes), procesador (USD 20 por vCPU por mes), volumen (USD 0,15 por GB por mes) y tráfico de salida (USD 0,05 por GB).
- Un sitio chico como este, con poco tráfico, debería quedar cerca de los USD 5 por mes. No es una garantía: el consumo real se ve en la pestaña *Usage* del proyecto.
- Volver a publicar **no cuesta nada extra**.
- Como el servicio tiene un volumen, **cada publicación tiene unos segundos sin servicio**. Si una publicación falla, sigue funcionando la versión anterior.

## 2. Subir esta copia a GitHub

Creá un repositorio **privado** nuevo en GitHub (vacío, sin README). Después, desde esta carpeta:

```powershell
git init
git add .
git commit -m "Version para publicar"
git branch -M main
git remote add origin https://github.com/TU-USUARIO/TU-REPOSITORIO.git
git push -u origin main
```

## 3. Crear la base de datos y cargar los datos (una sola vez)

1. En [railway.com](https://railway.com) creá un proyecto nuevo (*New Project*) y agregá una base de datos **PostgreSQL**. El servicio se llama `Postgres`.
2. Abrí el servicio `Postgres`, pestaña *Variables*, y copiá el valor de `DATABASE_PUBLIC_URL`.
3. Desde esta carpeta, cargá los datos (reemplazá la dirección por la que copiaste):

```powershell
$env:PGCLIENTENCODING = "UTF8"
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" "PEGAR-AQUI-DATABASE_PUBLIC_URL" -v ON_ERROR_STOP=1 -f db/lacis.sql
```

Vas a ver muchas líneas (`SET`, `CREATE TABLE`, `COPY 36`, algunos números): es normal. Lo importante es que **no aparezca ninguna línea que empiece con `ERROR`**. Si no hay ninguna, la base está lista. **No repitas este paso**: si lo corrés de nuevo, falla porque las tablas ya existen, y si lo forzaras pisaría lo que se cargó después desde el panel.

## 4. Crear el servicio web

1. En el mismo proyecto: *New → GitHub Repo*, y elegí el repositorio del paso 2. Railway detecta el `Dockerfile` y empieza a construir.
2. Abrí el servicio, pestaña *Variables*, y agregá estas tres:

| Variable | Valor |
|----------|-------|
| `DATABASE_URL` | `${{Postgres.DATABASE_URL}}` (escribilo tal cual; `Postgres` es el nombre del servicio de base de datos) |
| `SESSION_SECRET` | una clave larga y al azar (comando abajo) |
| `GIN_MODE` | `release` |

Para generar la clave de `SESSION_SECRET`:

```powershell
$b = New-Object byte[] 32; [Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($b); [Convert]::ToBase64String($b)
```

Guardá esa clave. Si la cambiás, se cierran todas las sesiones abiertas. `PORT` no hace falta: lo define Railway.

## 5. Agregar el volumen

En el servicio web, agregá un volumen (*Add Volume*; también se encuentra con Ctrl+K) y montalo en esta ruta, **exactamente**:

```text
/app/ui/static/assets
```

Al arrancar, el sitio copia ahí las imágenes y PDFs que ya vienen con el proyecto y no pisa nunca lo que se subió después.

## 6. Dirección pública y chequeo de salud

1. En el servicio web: *Settings → Networking → Generate Domain*.
2. En *Settings → Deploy*, poné la ruta de *Healthcheck* en `/`.
3. Publicá de nuevo (*Deploy*) si Railway no lo hizo solo.

## 7. Verificar que todo funciona

1. Abrí la dirección pública: se ve la página de inicio y, en `/lacis`, los productos de software con sus acentos.
2. Entrá a `/login` con el usuario `admin` y la contraseña `admin123`. Tenés que llegar al Panel de Administración con todos los módulos.
3. Probá la subida de archivos: en el panel, creá un integrante de prueba con una imagen, reiniciá el servicio (*Restart*) y comprobá que la imagen sigue viéndose. Después borrá ese integrante.

Si algo falla, mirá la pestaña *Logs* del servicio web:

- **"revisá la variable DATABASE_URL"**: la variable falta o está mal escrita.
- **Se ve el sitio pero no las imágenes**: falta el volumen o su ruta no es exactamente la del paso 5.
- **Siempre pide volver a iniciar sesión**: falta `SESSION_SECRET`.

## 8. Actualizar el sitio más adelante

Subí los cambios con `git push` a `main`. Railway construye y publica solo. **No vuelvas a cargar `db/lacis.sql`.** Si un cambio agrega una tabla o columna nueva, hay que aplicar ese cambio a mano en la base de Railway con `psql` (como en el paso 3, pero con el archivo del cambio).

### Si el sitio ya estaba publicado: descripción de Simon Riberi

La descripción de Simon Riberi cambió en `db/lacis.sql`, pero ese archivo no se vuelve a cargar. Si ya cargaste los datos antes, corregí esa fila a mano: guardá este comando en un archivo `actualizar-descripcion.sql` (codificación UTF-8) y ejecutalo con `psql "<DATABASE_PUBLIC_URL>" -f actualizar-descripcion.sql`. Tiene que responder `UPDATE 1`; después abrí `/integrantes` y revisá la tarjeta.

```sql
UPDATE integrante SET descripcion = 'Técnico Universitario en Web (UNSL). Desarrollador web con experiencia en proyectos reales y académicos utilizando diversas tecnologías entre las que se incluyen Go, PostgreSQL, JavaScript, PHP, React y Node.js. Participó en el diseño, desarrollo y mantenimiento del sitio web oficial de LaCIS como Práctica.' WHERE nombre = 'Tec. Riberi Zunino' AND apellido = 'Simon Martin';
```

## 9. Aviso importante

`admin` / `admin123` es la contraseña de la **versión de prueba**. Cualquier persona que la adivine puede entrar al panel y modificar o borrar contenido. Antes de usar el sitio como definitivo, entrá como `admin`, abrí *Administradores*, editá el usuario `admin` y poné una contraseña nueva.

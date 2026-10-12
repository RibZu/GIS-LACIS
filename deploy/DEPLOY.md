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

### 8.1 Antes de empezar

Necesitás la copia de esta carpeta en tu PC (la que se creó en el paso 2), Git instalado y acceso propio al repositorio de GitHub. Para actualizar el sitio **no** hace falta volver a cargar la base de datos, ni tocar el volumen, ni cambiar variables en Railway. Volver a publicar no cuesta nada extra.

### 8.2 Elegir el camino

- **Camino corto**: el cambio es solo de archivos estáticos (hojas de estilo, scripts, imágenes) o de plantillas HTML. No necesita la base de datos local.
- **Camino completo**: cualquier otro cambio (código en Go, datos, esquema) o si tenés dudas. Regenera toda la copia y necesita la base de datos local.

### 8.3 Camino corto: copiar solo lo que cambió

Copiá cada archivo modificado desde el proyecto de desarrollo a la misma ruta dentro de esta carpeta. Por ejemplo, para el cambio de la leyenda institucional en el celular:

```powershell
$Dev = "H:\GIS-LACIS-Simon"
$Copia = "H:\GIS-LACIS-Deploy"
Copy-Item "$Dev\ui\static\css-GIS\style.css" "$Copia\ui\static\css-GIS\style.css"
Copy-Item "$Dev\ui\html\html-GIS\index.html" "$Copia\ui\html\html-GIS\index.html"
Copy-Item "$Dev\deploy\DEPLOY.md" "$Copia\DEPLOY.md"
```

Después seguí en el punto 8.6. Si esta carpeta ya tiene su carpeta oculta `.git`, saltá el punto 8.5.

### 8.4 Camino completo: regenerar la copia

1. Levantá la base de datos local con `pg_ctl` (no uses `iniciar.bat`: cierra todos los procesos de PostgreSQL de la PC):

```powershell
& "C:\Program Files\PostgreSQL\18\bin\pg_ctl.exe" -D "H:\GIS-LACIS-Simon\pg_data" -o "-p 5433" -l "$env:TEMP\pg.log" start
```

2. Desde el proyecto de desarrollo, regenerá la copia:

```powershell
powershell -ExecutionPolicy Bypass -File deploy\generar_copia_deploy.ps1 -DesarrollosPrueba '172,173' -UsuariosPrueba '57' -Reemplazar
```

3. Apagá la base de datos:

```powershell
& "C:\Program Files\PostgreSQL\18\bin\pg_ctl.exe" -D "H:\GIS-LACIS-Simon\pg_data" stop
```

Regenerar con `-Reemplazar` borra la carpeta entera, **incluida su carpeta `.git`**. Por eso hay que conectarla de nuevo (punto 8.5).

### 8.5 Conectar la copia con el repositorio publicado

Hacelo solo si la carpeta no tiene `.git` (después del camino completo). Estos comandos van **dentro de la copia** y toda la configuración es local a esta carpeta: nunca uses `--global`.

```powershell
cd H:\GIS-LACIS-Deploy
git init
git branch -M main
git config core.autocrlf false
git config user.name gislacis-byte
git config user.email 338333962+gislacis-byte@users.noreply.github.com
git remote add origin https://gislacis-byte@github.com/gislacis-byte/GIS-LACIS.git
git fetch origin
git reset origin/main
```

`git reset origin/main` deja tus archivos como están y hace que Git compare contra lo ya publicado: así solo aparecen las diferencias reales y no se crea un historial paralelo. El usuario va dentro de la dirección (`https://gislacis-byte@...`) para que Git no use el inicio de sesión guardado de otra cuenta.

Si Git avisa `detected dubious ownership` (pasa en el disco H:), agregá `-c safe.directory=H:/GIS-LACIS-Deploy` justo después de `git` en cada comando, o pedí que se agregue esa carpeta como segura.

### 8.6 Revisar antes de subir

```powershell
git status
git diff --stat
```

Tienen que aparecer **solo** los archivos que quisiste cambiar (en el ejemplo, `style.css`, `index.html` y `DEPLOY.md`). Si aparecen muchos más, o archivos "borrados" que el repositorio publicado tiene y tu copia no, **no subas nada**: significa que el repositorio publicado tiene cambios que tu copia no incluye. Traelos primero con `git fetch origin` y revisalos. No uses `git push --force` bajo ningún concepto.

### 8.7 Confirmar y subir

```powershell
git add .
git commit -m "Leyenda institucional visible en el inicio del celular"
```

```powershell
git push origin main
```

Git te va a pedir tus credenciales de GitHub en ese momento: ingresalas vos, esta guía nunca las incluye. Antes del commit podés comprobar la identidad con `git config user.email`: tiene que devolver la dirección `noreply` de la cuenta de la materia.

### 8.8 Verificar

1. En Railway, abrí el servicio web y esperá a que el último despliegue figure como terminado. Durante unos segundos el sitio no responde; si el despliegue falla, sigue funcionando la versión anterior.
2. Abrí la dirección pública desde un celular (o con el navegador en ancho de celular) y comprobá que debajo de "Grupo de Ingeniería de Software" se ve, al mismo tiempo, "Departamento de Informática · Universidad Nacional de San Luis".
3. Revisá que `/lacis`, el login y una imagen subida desde el panel sigan igual.

### 8.9 Problemas conocidos

- **La carpeta no tiene `.git` y `git status` dice "not a git repository"**: se regeneró con `-Reemplazar`. Seguí el punto 8.5.
- **`detected dubious ownership`**: el disco H: no guarda propietarios. Usá `-c safe.directory=H:/GIS-LACIS-Deploy` como se explica en 8.5.
- **El `push` pide las credenciales de otra cuenta, o responde que no tenés permiso**: el remoto tiene que llevar el usuario en la dirección (`https://gislacis-byte@github.com/...`). Corregilo con `git remote set-url origin https://gislacis-byte@github.com/gislacis-byte/GIS-LACIS.git`.
- **`git status` muestra todos los archivos como modificados**: faltó `git config core.autocrlf false`. Configuralo y volvé a mirar.
- **Los PDF o las imágenes aparecen como modificados sin que los hayas tocado**: es la misma causa; no confirmes hasta que `git status` quede limpio.
- **Railway tarda o el despliegue falla**: mirá la pestaña *Logs* del servicio. Mientras tanto sigue publicada la versión anterior.

### Si el sitio ya estaba publicado: descripción de Simon Riberi

La descripción de Simon Riberi cambió en `db/lacis.sql`, pero ese archivo no se vuelve a cargar. Si ya cargaste los datos antes, corregí esa fila a mano: guardá este comando en un archivo `actualizar-descripcion.sql` (codificación UTF-8) y ejecutalo con `psql "<DATABASE_PUBLIC_URL>" -f actualizar-descripcion.sql`. Tiene que responder `UPDATE 1`; después abrí `/integrantes` y revisá la tarjeta.

```sql
UPDATE integrante SET descripcion = 'Técnico Universitario en Web (UNSL). Desarrollador web con experiencia en proyectos reales y académicos utilizando diversas tecnologías entre las que se incluyen Go, PostgreSQL, JavaScript, PHP, React y Node.js. Participó en el diseño, desarrollo y mantenimiento del sitio web oficial de LaCIS como Práctica.' WHERE nombre = 'Tec. Riberi Zunino' AND apellido = 'Simon Martin';
```

## 9. Aviso importante

`admin` / `admin123` es la contraseña de la **versión de prueba**. Cualquier persona que la adivine puede entrar al panel y modificar o borrar contenido. Antes de usar el sitio como definitivo, entrá como `admin`, abrí *Administradores*, editá el usuario `admin` y poné una contraseña nueva.

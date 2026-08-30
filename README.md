# GIS-LACIS Project

Bienvenido al proyecto GIS-LACIS. Este proyecto es una aplicación web desarrollada en **Go** que utiliza **PostgreSQL** como base de datos y sirve páginas dinámicas generadas desde el servidor (Server-Side Rendering).

## Tecnologías Necesarias

Para que otra persona pueda correr este proyecto en su computadora, deberá instalar:

1. **Go (Golang)**: Versión 1.20 o superior. [Descargar Go](https://go.dev/dl/)
2. **PostgreSQL**: Versión 18 o la más reciente. [Descargar PostgreSQL](https://www.postgresql.org/download/)

## Cómo Iniciar el Proyecto

El repositorio incluye un script llamado `iniciar.bat` que automatiza el proceso de despliegue. Está diseñado para facilitar la ejecución del entorno local a cualquier desarrollador del equipo.

Para ejecutar la aplicación, haz doble clic en el archivo **`iniciar.bat`**. 
- En la primera ejecución, el script detectará la instalación de PostgreSQL, inicializará una base de datos local automáticamente, creará las tablas necesarias e insertará los datos iniciales.
- Posteriormente, iniciará el servidor de Go en una nueva ventana de consola.
- La aplicación estará disponible ingresando a `http://localhost:8080` en el navegador web.

*(Nota: Es necesario mantener abierta la ventana de la consola. Al cerrarla, tanto la base de datos como el servidor web se detendrán).*

### Base de Datos: Sincronización de Cambios del Equipo

**Importante:** En caso de que el equipo necesite realizar cambios en la estructura de la base de datos (agregar tablas, columnas, etc.), **deben sobreescribir y actualizar el archivo `Modelo Base de datos.sql`**. El archivo debe mantener **exactamente ese mismo nombre**. No se deben crear archivos nuevos (como "v2.sql") ya que el script de inicio solo reconoce y ejecuta el archivo original.

Si se descargan cambios estructurales desde el repositorio, la base de datos local **no se actualizará automáticamente**. Esta medida existe para proteger los datos de prueba locales de cada desarrollador.

Para sincronizar la base de datos local con los últimos cambios estructurales del equipo, se debe ejecutar el script **`resetear_db.bat`**:
1. Ejecuta el archivo `resetear_db.bat`.
2. El script eliminará la base de datos local anterior (los datos de prueba locales se perderán).
3. Inmediatamente después, creará una nueva base de datos basada en el archivo `.sql` más reciente.
4. El proyecto se iniciará automáticamente de forma habitual.

#### Últimos cambios en el esquema (`usuario_gestor`)

- Se agregó la columna `modulos VARCHAR(255) NOT NULL DEFAULT ''`: lista de módulos separados por coma (ej. `"integrantes,proyectos"`) que el usuario puede gestionar. El campo `rol` ahora se **calcula** a partir de `modulos` (un solo módulo → `"GESTOR <MODULO>"`, dos o más → `"GESTOR"`), no se carga a mano desde el formulario.
- El usuario `admin` sembrado por defecto ahora se crea con `rol = 'ADMIN'` (antes `'GESTOR'`) — es el único rol reservado que nunca se calcula desde módulos.

Si tu base local es de antes de este cambio, corré `resetear_db.bat` para traerla al día.

---

## Guía de Arquitectura y Estructura de Carpetas

El proyecto sigue una **arquitectura en capas** clásica (MVC / Clean Architecture simplificada). Para mantener un orden estricto, la estructura del frontend (vistas) es un reflejo ("espejo") de la estructura del backend (lógica de negocio). 

### Organización "Espejo"

Cada dominio principal de la base de datos tiene su propio espacio. Aquí es donde debes colocar cada cosa según lo que vayas a programar:

1. **Proyectos y Desarrollos**
   - Backend: `internal/proyecto/` (Lógica en Go, consultas SQL).
   - Frontend: `ui/html/admin/proyecto/` (Vistas HTML y Panel).

2. **Premios y Reconocimientos**
   - Backend: `internal/reconocimiento/`
   - Frontend: `ui/html/admin/reconocimiento/`

3. **Empresas Colaboradoras**
   - Backend: `internal/colaborador/`
   - Frontend: `ui/html/admin/colaborador/`

4. **Gestión de Usuarios y Autenticación**
   - Backend: `internal/usuario/`
   - Frontend: `ui/html/admin/usuario/`

5. **Tesis Académicas**
   - Backend: `internal/tesis/`
   - Frontend: `ui/html/admin/tesis/`

6. **Integrantes**
   - Backend: `internal/integrante/`
   - Frontend: `ui/html/admin/` (Módulo inicial).

### Otras Carpetas Importantes

- **`api/`**: Contiene la configuración de rutas web (Gin) y los `handlers` (controladores) que comunican las peticiones HTTP de las vistas HTML con las carpetas `internal/`.
- **`cmd/`**: Punto de entrada de la aplicación. `cmd/api/main.go` es el archivo principal.
- **`Modelo Base de datos/`**: Scripts de SQL (DDL) para la estructura relacional.
- **`ui/static/`**: Archivos estáticos públicos (CSS, JS, Imágenes, CVs).
  - **`ui/static/mockups/`**: Prototipos visuales de rediseño (HTML standalone, con estilos embebidos) — para comparar propuestas de estilo lado a lado antes de aplicarlas a las plantillas reales. No están conectados a la app en vivo ni a datos reales.
- **`test/`**: Pruebas unitarias y de integración en Go.

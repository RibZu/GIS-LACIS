# GIS-LACIS Project

Bienvenido al proyecto GIS-LACIS. Este proyecto es una aplicación web desarrollada en **Go** que utiliza **PostgreSQL** como base de datos y sirve páginas dinámicas generadas desde el servidor (Server-Side Rendering).

## Tecnologías Necesarias

Para que otra persona pueda correr este proyecto en su computadora, deberá instalar:

1. **Go (Golang)**: Versión 1.20 o superior. [Descargar Go](https://go.dev/dl/)
2. **PostgreSQL**: Versión 18 o la más reciente. [Descargar PostgreSQL](https://www.postgresql.org/download/)

## Cómo Levantar el Proyecto

He creado un archivo llamado `iniciar.bat` que automatiza todo el proceso. **¡Está pensado para que cualquier desarrollador pueda descargarlo de GitHub y hacerlo andar al instante!**

Simplemente haz doble clic en el archivo **`iniciar.bat`**. 
- Si es la primera vez que lo corres, el script detectará tu instalación de PostgreSQL, inicializará una base de datos local automáticamente (sin pedir contraseñas extrañas), creará las tablas necesarias e insertará todo.
- Luego, iniciará el servidor de Go en una nueva ventana.
- Podrás acceder a la página ingresando a `http://localhost:8080` en tu navegador.

*(Nota: Al cerrar la ventana del script, la base de datos se detendrá sola de manera limpia).*

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
- **`test/`**: Pruebas unitarias y de integración en Go.

---

## Panel de Administración y Autenticación

El proyecto incluye un sistema de Login y un **Dashboard de Administración** para gestionar el sitio (historias de usuario SCRUM).
- **Ruta**: `/login` (Accesible desde el botón del menú superior).
- **Usuario de prueba**: `admin`
- **Contraseña de prueba**: `admin123`

Desde el Dashboard, los administradores tendrán acceso preventivo (botones) para dirigirse a los módulos de ABM (Alta, Baja, Modificación) que irás desarrollando en sus respectivas carpetas "espejo".

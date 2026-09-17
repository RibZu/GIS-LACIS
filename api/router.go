package api

import (
	"PaginaSEG/api/handler"
	"PaginaSEG/internal/colaborador"
	"PaginaSEG/internal/desarrollo"
	"PaginaSEG/internal/integrante"
	"PaginaSEG/internal/proyecto"
	"PaginaSEG/internal/reconocimiento"
	"PaginaSEG/internal/tesis"
	"PaginaSEG/internal/usuario"
	"database/sql"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func InitRoutes(e *gin.Engine) {

	e.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	logger, err := zap.NewProduction()
	defer logger.Sync()

	dsn := "postgres://postgres:isma_mesa22@localhost:5433/lacis?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatal("No se pudo abrir conexión a PostgreSQL", zap.Error(err))
	}

	for i := 0; i < 20; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		logger.Warn("Esperando a que la base de datos inicie...", zap.Error(err))
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		logger.Fatal("No se pudo hacer Ping a PostgreSQL después de varios intentos", zap.Error(err))
	}

	e.Static("/static", "ui/static")
	e.Static("/ui/static", "ui/static")

	e.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int { return a + b },
	})
	e.LoadHTMLGlob("ui/html/**/*.html")

	integranteStorage := integrante.NewPostgresStorage(db)
	integranteService := integrante.NewService(integranteStorage, logger)
	integranteHandler := handler.NewIntegranteHandler(integranteService, logger)

	tesisStorage := tesis.NewPostgresStorage(db)
	tesisService := tesis.NewService(tesisStorage, logger)
	tesisHandler := handler.NewTesisHandler(tesisService, integranteService, logger)

	desarrolloStorage := desarrollo.NewPostgresStorage(db)
	desarrolloService := desarrollo.NewService(desarrolloStorage, logger)
	desarrolloHandler := handler.NewDesarrolloHandler(desarrolloService, integranteService, logger)

	usuarioStorage := usuario.NewPostgressStorage(db)
	usuarioService := usuario.NewService(usuarioStorage, logger)
	authHandler := handler.NewAuthHandler(usuarioService, integranteService, logger)
	usuarioHandler := handler.NewUsuarioHandler(usuarioService, logger)

	proyectoStorage := proyecto.NewPostgresStorage(db)
	proyectoService := proyecto.NewService(proyectoStorage, logger)
	proyectoHandler := handler.NewProyectoHandler(proyectoService, logger)

	reconocimientoStorage := reconocimiento.NewPostgresStorage(db)
	reconocimientoService := reconocimiento.NewService(reconocimientoStorage, logger)
	reconocimientoHandler := handler.NewReconocimientoHandler(reconocimientoService, logger)

	e.GET("/", func(c *gin.Context) {
		_, loggedIn := handler.CurrentUserID(c)
		c.HTML(http.StatusOK, "index.html", gin.H{"LoggedIn": loggedIn})
	})

	e.GET("/integrantes", func(c *gin.Context) {
		_, loggedIn := handler.CurrentUserID(c)
		c.HTML(http.StatusOK, "integrantes.html", gin.H{"LoggedIn": loggedIn})
	})

	e.GET("/lacis", desarrolloHandler.ViewLacis)

	e.GET("/proyectos", func(c *gin.Context) {
		_, loggedIn := handler.CurrentUserID(c)
		c.HTML(http.StatusOK, "Proyecto.html", gin.H{"LoggedIn": loggedIn})
	})

	e.GET("/doctorado-ing-software", func(c *gin.Context) {
		_, loggedIn := handler.CurrentUserID(c)
		c.HTML(http.StatusOK, "DrIngSoft.html", gin.H{"LoggedIn": loggedIn})
	})

	e.GET("/especializacion-ing-software", func(c *gin.Context) {
		_, loggedIn := handler.CurrentUserID(c)
		c.HTML(http.StatusOK, "EspIngSoft.html", gin.H{"LoggedIn": loggedIn})
	})

	e.GET("/maestria-calidad-software", func(c *gin.Context) {
		_, loggedIn := handler.CurrentUserID(c)
		c.HTML(http.StatusOK, "MgCalSoft.html", gin.H{"LoggedIn": loggedIn})
	})

	e.GET("/maestria-ing-software", func(c *gin.Context) {
		_, loggedIn := handler.CurrentUserID(c)
		c.HTML(http.StatusOK, "MgIngSoft.html", gin.H{"LoggedIn": loggedIn})
	})

	e.GET("/desarrollos", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/lacis#productos-software")
	})

	e.GET("/tesis", tesisHandler.ViewPublica)

	colaboradorStorage := colaborador.NewPostgresStorage(db)
	colaboradorService := colaborador.NewService(colaboradorStorage, logger)
	colaboradorHandler := handler.NewColaboradorHandler(colaboradorService, logger)

	e.GET("/login", authHandler.ShowLogin)
	e.POST("/login", authHandler.ProcessLogin)
	e.GET("/logout", authHandler.Logout)

	v1Admin := e.Group("/admin", handler.RequireLogin(usuarioService))
	v1Admin.GET("/dashboard", authHandler.ShowDashboard)

	integrantesAdmin := v1Admin.Group("")
	integrantesAdmin.Use(handler.RequireModule(usuarioService, "integrantes"))
	integrantesAdmin.GET("/integrantes", integranteHandler.Lista)
	integrantesAdmin.GET("/crear-integrante", integranteHandler.Crear)
	integrantesAdmin.POST("/insertar-integrante", integranteHandler.Insertar)
	integrantesAdmin.GET("/editar-integrante", integranteHandler.Editar)
	integrantesAdmin.POST("/actualizar-integrante", integranteHandler.Actualizar)
	integrantesAdmin.GET("/borrar-integrante", integranteHandler.Borrar)

	tesisAdmin := v1Admin.Group("")
	tesisAdmin.Use(handler.RequireModule(usuarioService, "tesis"))
	tesisAdmin.GET("/tesis", tesisHandler.Lista)
	tesisAdmin.GET("/crear-tesis", tesisHandler.Crear)
	tesisAdmin.POST("/insertar-tesis", tesisHandler.Insertar)
	tesisAdmin.GET("/editar-tesis", tesisHandler.Editar)
	tesisAdmin.POST("/actualizar-tesis", tesisHandler.Actualizar)
	tesisAdmin.GET("/borrar-tesis", tesisHandler.Borrar)

	desarrollosAdmin := v1Admin.Group("")
	desarrollosAdmin.Use(handler.RequireModule(usuarioService, "proyectos"))
	desarrollosAdmin.GET("/desarrollos", desarrolloHandler.Lista)
	desarrollosAdmin.GET("/crear-desarrollo", desarrolloHandler.Crear)
	desarrollosAdmin.POST("/insertar-desarrollo", desarrolloHandler.Insertar)
	desarrollosAdmin.GET("/editar-desarrollo", desarrolloHandler.Editar)
	desarrollosAdmin.POST("/actualizar-desarrollo", desarrolloHandler.Actualizar)
	desarrollosAdmin.GET("/borrar-desarrollo", desarrolloHandler.Borrar)

	usuariosAdmin := v1Admin.Group("")
	usuariosAdmin.Use(handler.RequireAdmin(usuarioService))
	usuariosAdmin.GET("/usuarios", usuarioHandler.Lista)
	usuariosAdmin.GET("/crear-usuario", usuarioHandler.Crear)
	usuariosAdmin.POST("/insertar-usuario", usuarioHandler.Insertar)
	usuariosAdmin.GET("/editar-usuario", usuarioHandler.Editar)
	usuariosAdmin.POST("/actualizar-usuario", usuarioHandler.Actualizar)
	usuariosAdmin.GET("/borrar-usuario", usuarioHandler.Borrar)

	v1API := e.Group("/api/v1")
	v1API.GET("/integrantes", integranteHandler.API_GetAll)
	v1API.GET("/integrantes/:id", integranteHandler.API_Read)
	v1API.GET("/tesis", tesisHandler.API_GetAll)
	v1API.GET("/tesis/:id", tesisHandler.API_Read)

	proyectosAdmin := v1Admin.Group("")
	proyectosAdmin.Use(handler.RequireModule(usuarioService, "proyectos"))

	proyectosAdmin.GET("/proyectos", proyectoHandler.View_ProyectosAdmin)

	v1API.GET("/proyectos", proyectoHandler.API_GetAll)
	v1API.GET("/reconocimientos", reconocimientoHandler.API_GetAll)
	v1API.GET("/colaboradores", colaboradorHandler.API_GetAll)

	proyectosAPIAdmin := v1API.Group("")
	proyectosAPIAdmin.Use(handler.RequireModuleAPI(usuarioService, "proyectos"))
	proyectosAPIAdmin.GET("/admin/proyectos-todos", proyectoHandler.API_GetAllAdmin)
	proyectosAPIAdmin.PATCH("/admin/proyectos/:id/restaurar", proyectoHandler.API_Restaurar)
	proyectosAPIAdmin.POST("/admin/proyectos", proyectoHandler.API_Create)
	proyectosAPIAdmin.GET("/admin/proyectos/:id", proyectoHandler.API_Read)
	proyectosAPIAdmin.PUT("/admin/proyectos/:id", proyectoHandler.API_Update)
	proyectosAPIAdmin.DELETE("/admin/proyectos/:id", proyectoHandler.API_Delete)
	proyectosAPIAdmin.POST("/admin/proyectos/:id/equipo", proyectoHandler.API_AgregarMiembro)
	proyectosAPIAdmin.DELETE("/admin/proyectos/:id/equipo/:integrante_id", proyectoHandler.API_QuitarMiembro)

	// Módulo "reconocimientos": requiere que el usuario logueado tenga ese módulo asignado (o sea ADMIN)
	reconocimientosAdmin := v1Admin.Group("")
	reconocimientosAdmin.Use(handler.RequireModule(usuarioService, "reconocimientos"))
	reconocimientosAdmin.GET("/reconocimientos", reconocimientoHandler.View_ReconocimientosAdmin)

	// Rutas API de administración de Reconocimientos (CRUD placeholders)
	reconocimientosAPIAdmin := v1API.Group("")
	reconocimientosAPIAdmin.Use(handler.RequireModuleAPI(usuarioService, "reconocimientos"))
	reconocimientosAPIAdmin.GET("/admin/reconocimientos-todos", reconocimientoHandler.API_GetAllAdmin)
	reconocimientosAPIAdmin.POST("/admin/reconocimientos", reconocimientoHandler.API_Create)
	reconocimientosAPIAdmin.GET("/admin/reconocimientos/:id", reconocimientoHandler.API_Read)
	reconocimientosAPIAdmin.PUT("/admin/reconocimientos/:id", reconocimientoHandler.API_Update)
	reconocimientosAPIAdmin.DELETE("/admin/reconocimientos/:id", reconocimientoHandler.API_Delete)
	reconocimientosAPIAdmin.PATCH("/admin/reconocimientos/:id/restaurar", reconocimientoHandler.API_Restaurar)

	// Módulo "empresas" (Colaboradores): requiere que el usuario logueado tenga ese módulo asignado (o sea ADMIN)
	colaboradoresAdmin := v1Admin.Group("")
	colaboradoresAdmin.Use(handler.RequireModule(usuarioService, "empresas"))
	colaboradoresAdmin.GET("/colaboradores", colaboradorHandler.View_ColaboradoresAdmin)

	// Rutas API de administración de Colaboradores
	colaboradoresAPIAdmin := v1API.Group("")
	colaboradoresAPIAdmin.Use(handler.RequireModuleAPI(usuarioService, "empresas"))
	colaboradoresAPIAdmin.GET("/admin/colaboradores-todos", colaboradorHandler.API_GetAllAdmin)
	colaboradoresAPIAdmin.POST("/admin/colaboradores", colaboradorHandler.API_Create)
	colaboradoresAPIAdmin.GET("/admin/colaboradores/:id", colaboradorHandler.API_Read)
	colaboradoresAPIAdmin.PUT("/admin/colaboradores/:id", colaboradorHandler.API_Update)
	colaboradoresAPIAdmin.DELETE("/admin/colaboradores/:id", colaboradorHandler.API_Delete)
	colaboradoresAPIAdmin.PATCH("/admin/colaboradores/:id/restaurar", colaboradorHandler.API_Restaurar)

	// Público: consultar equipo de un proyecto
	v1API.GET("/proyectos/:id/equipo", proyectoHandler.API_GetEquipo)

	// Nuevo grupo admin para integrantes
	integrantesAPIAdmin := v1API.Group("")
	integrantesAPIAdmin.Use(handler.RequireModuleAPI(usuarioService, "integrantes"))
	integrantesAPIAdmin.POST("/admin/integrantes/mini", integranteHandler.API_CreateMinimo)

}

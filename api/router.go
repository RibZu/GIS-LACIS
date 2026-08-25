package api

import (
	"PaginaSEG/api/handler"
	"PaginaSEG/internal/integrante"
	"database/sql"
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
	//  Conexión a PostgreSQL
	dsn := "postgres://postgres:isma_mesa22@localhost:5433/lacis?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatal("No se pudo abrir conexión a PostgreSQL", zap.Error(err))
	}
	
	// Reintentar conexión hasta 20 veces (esperando a que inicie la DB tras crash recovery)
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

	// Servir archivos estáticos (CSS, JS, Imágenes)
	e.Static("/static", "ui/static")
	e.Static("/ui/static", "ui/static")
	// Cargar las plantillas HTML
	e.LoadHTMLGlob("ui/html/**/*.html")

	// RUTAS PARA PÁGINAS WEB (HTML) estatico

	// 1. Inicio
	e.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	// 2. Integrantes
	e.GET("/integrantes", func(c *gin.Context) {
		c.HTML(http.StatusOK, "integrantes.html", nil)
	})
	// 3. LaCIS
	e.GET("/lacis", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Lacis.html", nil)
	})
	// 4. Proyectos
	e.GET("/proyectos", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Proyecto.html", nil)
	})
	// 5. Posgrado: Doctorado en Ingeniería de Software
	e.GET("/doctorado-ing-software", func(c *gin.Context) {
		c.HTML(http.StatusOK, "DrIngSoft.html", nil)
	})
	// 6. Posgrado: Especialización en Ingeniería de Software
	e.GET("/especializacion-ing-software", func(c *gin.Context) {
		c.HTML(http.StatusOK, "EspIngSoft.html", nil)
	})
	// 7. Posgrado: Maestría en Calidad de Software
	e.GET("/maestria-calidad-software", func(c *gin.Context) {
		c.HTML(http.StatusOK, "MgCalSoft.html", nil)
	})
	// 8. Posgrado: Maestría en Ingeniería de Software
	e.GET("/maestria-ing-software", func(c *gin.Context) {
		c.HTML(http.StatusOK, "MgIngSoft.html", nil)
	})

	// Al usar plantillas de GO se hace de esta manera

	// RUTAS DE ADMINISTRACIÓN CON PLANTILLAS GO (SERVER-SIDE RENDERED)

	integranteStorage := integrante.NewPostgresStorage(db)
	integranteService := integrante.NewService(integranteStorage, logger)
	integranteHandler := handler.NewIntegranteHandler(integranteService, logger)
	authHandler := handler.NewAuthHandler()

	e.GET("/login", authHandler.ShowLogin)
	e.POST("/login", authHandler.ProcessLogin)
	e.GET("/logout", authHandler.Logout)

	v1Admin := e.Group("/admin")
	v1Admin.GET("/dashboard", authHandler.ShowDashboard)
	v1Admin.GET("/integrantes", integranteHandler.Lista)
	v1Admin.GET("/crear-integrante", integranteHandler.Crear)
	v1Admin.POST("/insertar-integrante", integranteHandler.Insertar)
	v1Admin.GET("/editar-integrante", integranteHandler.Editar)
	v1Admin.POST("/actualizar-integrante", integranteHandler.Actualizar)
	v1Admin.GET("/borrar-integrante", integranteHandler.Borrar)

	// RUTAS PARA LA API REST JSON (OPCIONALES PARA POSTMAN),

	v1API := e.Group("/api/v1")
	v1API.GET("/integrantes", integranteHandler.API_GetAll)
	v1API.GET("/integrantes/:id", integranteHandler.API_Read)

}

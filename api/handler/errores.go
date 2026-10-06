package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var avisosOK = map[string]string{
	"producto-creado":      "Producto guardado.",
	"producto-actualizado": "Cambios guardados.",
	"producto-eliminado":   "Producto eliminado.",
	"usuario-creado":       "Usuario creado.",
	"usuario-actualizado":  "Cambios guardados.",
	"usuario-eliminado":    "Usuario eliminado.",
}

var avisosError = map[string]string{
	"no-encontrado":   "No encontramos ese registro.",
	"error-eliminar":  "No se pudo eliminar. Probá de nuevo en unos minutos.",
	"autoeliminacion": "No podés eliminar tu propia cuenta.",
	"ultimo-admin":    "Tiene que quedar al menos un administrador.",
}

func avisosDesdeQuery(c *gin.Context) (ok string, aviso string) {
	return avisosOK[c.Query("ok")], avisosError[c.Query("error")]
}

func paginaError(c *gin.Context, estado int, titulo, mensaje string) {
	c.HTML(estado, "ErrorPagina.html", gin.H{
		"Titulo":  titulo,
		"Mensaje": mensaje,
	})
}

func NuevoRecuperador(logger *zap.Logger) gin.RecoveryFunc {
	return func(c *gin.Context, err any) {
		logger.Error("Pánico en un handler",
			zap.Any("panic", err),
			zap.String("ruta", c.Request.URL.Path),
			zap.Stack("stack"),
		)
		paginaError(c, http.StatusInternalServerError, "Algo salió mal", "Algo salió mal. Probá de nuevo en unos minutos.")
		c.Abort()
	}
}

func PaginaNoEncontrada(c *gin.Context) {
	paginaError(c, http.StatusNotFound, "Página no encontrada", "No encontramos esa página.")
}

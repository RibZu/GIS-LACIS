package handler

import (
	"net/http"

	"PaginaSEG/internal/estadistica"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// EstadisticaHandler maneja la vista del módulo de estadísticas (solo lectura).
type EstadisticaHandler struct {
	service *estadistica.Service
	logger  *zap.Logger
}

func NewEstadisticaHandler(s *estadistica.Service, l *zap.Logger) *EstadisticaHandler {
	return &EstadisticaHandler{service: s, logger: l}
}

// Ver renderiza /admin/estadisticas. Si no se pueden leer los datos muestra un aviso y un 500,
// sin ningún número.
func (h *EstadisticaHandler) Ver(c *gin.Context) {
	est, err := h.service.Obtener()
	if err != nil {
		h.logger.Error("Error al obtener estadísticas", zap.Error(err))
		c.HTML(http.StatusInternalServerError, "Estadisticas.html", gin.H{
			"Error":    "No se pudieron cargar las estadísticas. Intentá de nuevo en unos minutos.",
			"LoggedIn": true,
		})
		return
	}
	c.HTML(http.StatusOK, "Estadisticas.html", gin.H{"E": est, "LoggedIn": true})
}

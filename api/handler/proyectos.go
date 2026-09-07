package handler

import (
	"errors"
	"net/http"
	"strconv"

	"PaginaSEG/internal/proyecto"
	"PaginaSEG/internal/reconocimiento"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProyectoHandler struct {
	service *proyecto.Service
	logger  *zap.Logger
}

func NewProyectoHandler(s *proyecto.Service, l *zap.Logger) *ProyectoHandler {
	return &ProyectoHandler{
		service: s,
		logger:  l,
	}
}

// API_GetAll retorna la lista de proyectos en JSON
func (h *ProyectoHandler) API_GetAll(c *gin.Context) {
	proyectos, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Error al obtener proyectos", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener proyectos"})
		return
	}
	c.JSON(http.StatusOK, proyectos)
}

// API_Create crea un nuevo proyecto desde payload JSON
func (h *ProyectoHandler) API_Create(c *gin.Context) {
	var p proyecto.Proyecto
	if err := c.ShouldBindJSON(&p); err != nil {
		h.logger.Warn("JSON inválido al crear proyecto", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if err := h.service.Create(&p); err != nil {
		h.logger.Warn("Error de validación al crear proyecto", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

// API_Read obtiene un proyecto por ID
func (h *ProyectoHandler) API_Read(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	p, err := h.service.Read(id)
	if err != nil {
		if errors.Is(err, proyecto.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "proyecto no encontrado"})
			return
		}
		h.logger.Error("Error al leer proyecto", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al consultar proyecto"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// API_Update actualiza campos de un proyecto por ID
func (h *ProyectoHandler) API_Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var fields proyecto.UpdateFields
	if err := c.ShouldBindJSON(&fields); err != nil {
		h.logger.Warn("JSON inválido al actualizar proyecto", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if err := h.service.Update(id, fields); err != nil {
		if errors.Is(err, proyecto.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "proyecto no encontrado"})
			return
		}
		if errors.Is(err, proyecto.ErrTituloRequerido) || errors.Is(err, proyecto.ErrAniosInvalidos) || errors.Is(err, proyecto.ErrIDInvalido) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("Error al actualizar proyecto", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar proyecto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "proyecto actualizado exitosamente"})
}

// API_Delete elimina un proyecto por ID
func (h *ProyectoHandler) API_Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, proyecto.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "proyecto no encontrado"})
			return
		}
		h.logger.Error("Error al eliminar proyecto", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al eliminar proyecto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "proyecto eliminado exitosamente"})
}

// Vistas HTML para el Panel de Administración

// View_ProyectosAdmin renderiza la vista de proyectos
func (h *ProyectoHandler) View_ProyectosAdmin(c *gin.Context) {
	c.HTML(http.StatusOK, "EditarProyectos.html", nil)
}

// Handler de Reconocimientos
type ReconocimientoHandler struct {
	service *reconocimiento.Service
	logger  *zap.Logger
}

func NewReconocimientoHandler(s *reconocimiento.Service, l *zap.Logger) *ReconocimientoHandler {
	return &ReconocimientoHandler{
		service: s,
		logger:  l,
	}
}

// API_GetAll retorna la lista de reconocimientos en JSON
func (h *ReconocimientoHandler) API_GetAll(c *gin.Context) {
	recs, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Error al obtener reconocimientos", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener reconocimientos"})
		return
	}
	c.JSON(http.StatusOK, recs)
}

// API_GetAllAdmin retorna todos los proyectos (incluyendo inactivos) para el panel admin
func (h *ProyectoHandler) API_GetAllAdmin(c *gin.Context) {
	proyectos, err := h.service.GetAllAdmin()
	if err != nil {
		h.logger.Error("Error al obtener proyectos admin", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener proyectos"})
		return
	}
	c.JSON(http.StatusOK, proyectos)
}

// API_Restaurar vuelve a activar un proyecto dado de baja lógica
func (h *ProyectoHandler) API_Restaurar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.service.Restaurar(id); err != nil {
		if errors.Is(err, proyecto.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "proyecto no encontrado"})
			return
		}
		h.logger.Error("Error al restaurar proyecto", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al restaurar proyecto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "proyecto restaurado exitosamente"})
}

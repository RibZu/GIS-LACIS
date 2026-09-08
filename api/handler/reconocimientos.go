package handler

import (
	"errors"
	"net/http"
	"strconv"

	"PaginaSEG/internal/reconocimiento"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ReconocimientoHandler maneja las peticiones HTTP relacionadas a reconocimientos
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

// API_GetAll retorna la lista de reconocimientos activos en JSON (endpoint público)
func (h *ReconocimientoHandler) API_GetAll(c *gin.Context) {
	recs, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Error al obtener reconocimientos", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener reconocimientos"})
		return
	}
	c.JSON(http.StatusOK, recs)
}

// API_GetAllAdmin retorna todos los reconocimientos (incluye inactivos) para el panel admin
func (h *ReconocimientoHandler) API_GetAllAdmin(c *gin.Context) {
	recs, err := h.service.GetAllAdmin()
	if err != nil {
		h.logger.Error("Error al obtener reconocimientos (admin)", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener reconocimientos"})
		return
	}
	c.JSON(http.StatusOK, recs)
}

// View_ReconocimientosAdmin renderiza la vista HTML del panel admin de reconocimientos
func (h *ReconocimientoHandler) View_ReconocimientosAdmin(c *gin.Context) {
	c.HTML(http.StatusOK, "EditarReconocimientos.html", nil)
}

// API_Create crea un nuevo reconocimiento desde payload JSON
func (h *ReconocimientoHandler) API_Create(c *gin.Context) {
	var r reconocimiento.Reconocimiento
	if err := c.ShouldBindJSON(&r); err != nil {
		h.logger.Warn("JSON inválido al crear reconocimiento", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if err := h.service.Create(&r); err != nil {
		h.logger.Warn("Error de validación al crear reconocimiento", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, r)
}

// API_Read obtiene un reconocimiento por ID
func (h *ReconocimientoHandler) API_Read(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	r, err := h.service.Read(id)
	if err != nil {
		if errors.Is(err, reconocimiento.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "reconocimiento no encontrado"})
			return
		}
		h.logger.Error("Error al leer reconocimiento", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al consultar reconocimiento"})
		return
	}

	c.JSON(http.StatusOK, r)
}

// API_Update actualiza campos de un reconocimiento por ID
func (h *ReconocimientoHandler) API_Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var fields reconocimiento.UpdateFields
	if err := c.ShouldBindJSON(&fields); err != nil {
		h.logger.Warn("JSON inválido al actualizar reconocimiento", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if err := h.service.Update(id, fields); err != nil {
		if errors.Is(err, reconocimiento.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "reconocimiento no encontrado"})
			return
		}
		h.logger.Error("Error al actualizar reconocimiento", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar reconocimiento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "reconocimiento actualizado exitosamente"})
}

// API_Delete realiza la baja lógica de un reconocimiento (activo = false)
func (h *ReconocimientoHandler) API_Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, reconocimiento.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "reconocimiento no encontrado"})
			return
		}
		h.logger.Error("Error al dar de baja reconocimiento", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al dar de baja reconocimiento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "reconocimiento dado de baja exitosamente"})
}

// API_Restaurar reactiva un reconocimiento dado de baja lógica
func (h *ReconocimientoHandler) API_Restaurar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.service.Restaurar(id); err != nil {
		if errors.Is(err, reconocimiento.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "reconocimiento no encontrado"})
			return
		}
		h.logger.Error("Error al restaurar reconocimiento", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al restaurar reconocimiento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "reconocimiento restaurado exitosamente"})
}

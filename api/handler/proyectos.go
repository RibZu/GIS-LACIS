package handler

import (
	"errors"
	"net/http"
	"strconv"

	"PaginaSEG/internal/proyecto"

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

type miembroEquipoInput struct {
	IntegranteID  int    `json:"integrante_id" binding:"required"`
	RolEnProyecto string `json:"rol_en_proyecto"`
}

// API_GetEquipo devuelve el equipo actual de un proyecto (endpoint público)
func (h *ProyectoHandler) API_GetEquipo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	equipo, err := h.service.GetEquipo(id)
	if err != nil {
		h.logger.Error("Error al obtener equipo", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener el equipo del proyecto"})
		return
	}
	c.JSON(http.StatusOK, equipo)
}

// API_AgregarMiembro agrega (o actualiza el rol de) un integrante al equipo del proyecto
func (h *ProyectoHandler) API_AgregarMiembro(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var input miembroEquipoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}
	if err := h.service.AgregarMiembro(id, input.IntegranteID, input.RolEnProyecto); err != nil {
		h.logger.Error("Error al agregar miembro al proyecto", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al agregar el integrante al equipo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje": "integrante agregado al equipo"})
}

// API_QuitarMiembro quita un integrante del equipo del proyecto
func (h *ProyectoHandler) API_QuitarMiembro(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	integranteID, err := strconv.Atoi(c.Param("integrante_id"))
	if err != nil || integranteID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de integrante inválido"})
		return
	}
	if err := h.service.QuitarMiembro(id, integranteID); err != nil {
		if errors.Is(err, proyecto.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "el integrante no pertenece al equipo de este proyecto"})
			return
		}
		h.logger.Error("Error al quitar miembro del proyecto", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al quitar el integrante del equipo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje": "integrante quitado del equipo"})
}

package handler

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"PaginaSEG/internal/colaborador"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ColaboradorHandler maneja las peticiones HTTP de empresas colaboradoras
type ColaboradorHandler struct {
	service *colaborador.Service
	logger  *zap.Logger
}

func NewColaboradorHandler(s *colaborador.Service, l *zap.Logger) *ColaboradorHandler {
	return &ColaboradorHandler{
		service: s,
		logger:  l,
	}
}

// API_GetAll retorna la lista de colaboradores activos en JSON (endpoint público)
func (h *ColaboradorHandler) API_GetAll(c *gin.Context) {
	colabs, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Error al obtener colaboradores", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener colaboradores"})
		return
	}
	c.JSON(http.StatusOK, colabs)
}

// API_GetAllAdmin retorna todos los colaboradores (incluye inactivos) para el panel admin
func (h *ColaboradorHandler) API_GetAllAdmin(c *gin.Context) {
	colabs, err := h.service.GetAllAdmin()
	if err != nil {
		h.logger.Error("Error al obtener colaboradores (admin)", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener colaboradores"})
		return
	}
	c.JSON(http.StatusOK, colabs)
}

// View_ColaboradoresAdmin renderiza la vista HTML del panel admin de colaboradores
func (h *ColaboradorHandler) View_ColaboradoresAdmin(c *gin.Context) {
	c.HTML(http.StatusOK, "EditarColaboradores.html", nil)
}

// Helper to handle file upload
func (h *ColaboradorHandler) saveLogoFile(c *gin.Context) (string, error) {
	file, err := c.FormFile("logo")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return "", nil // No file provided
		}
		return "", err
	}

	// Save to ui/static/assets/img-GIS
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("colab_%d%s", time.Now().UnixNano(), ext)
	uploadPath := filepath.Join("ui", "static", "assets", "img-GIS", filename)

	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		return "", err
	}

	return "/static/assets/img-GIS/" + filename, nil
}

// API_Create crea un nuevo colaborador desde form-data
func (h *ColaboradorHandler) API_Create(c *gin.Context) {
	descripcion := c.PostForm("descripcion")

	logoURL, err := h.saveLogoFile(c)
	if err != nil {
		h.logger.Warn("Error al procesar archivo", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al subir la imagen"})
		return
	}

	item := colaborador.Colaborador{
		Descripcion: descripcion,
		LogoURL:     logoURL,
	}

	if err := h.service.Create(&item); err != nil {
		h.logger.Warn("Error de validación al crear colaborador", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// API_Read obtiene un colaborador por ID
func (h *ColaboradorHandler) API_Read(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	colab, err := h.service.Read(id)
	if err != nil {
		if errors.Is(err, colaborador.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "colaborador no encontrado"})
			return
		}
		h.logger.Error("Error al leer colaborador", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al consultar colaborador"})
		return
	}

	c.JSON(http.StatusOK, colab)
}

// API_Update actualiza campos de un colaborador por ID (form-data)
func (h *ColaboradorHandler) API_Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var fields colaborador.UpdateFields

	if descripcion, ok := c.GetPostForm("descripcion"); ok {
		fields.Descripcion = &descripcion
	}

	logoURL, err := h.saveLogoFile(c)
	if err != nil {
		h.logger.Warn("Error al procesar archivo", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al subir la imagen"})
		return
	}
	if logoURL != "" {
		fields.LogoURL = &logoURL
	}

	if err := h.service.Update(id, fields); err != nil {
		if errors.Is(err, colaborador.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "colaborador no encontrado"})
			return
		}
		h.logger.Error("Error al actualizar colaborador", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar colaborador"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "colaborador actualizado exitosamente"})
}

func (h *ColaboradorHandler) API_Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, colaborador.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "colaborador no encontrado"})
			return
		}
		h.logger.Error("Error al dar de baja colaborador", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al dar de baja colaborador"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "colaborador dado de baja exitosamente"})
}

// API_Restaurar reactiva un colaborador dado de baja
func (h *ColaboradorHandler) API_Restaurar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.service.Restaurar(id); err != nil {
		if errors.Is(err, colaborador.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "colaborador no encontrado"})
			return
		}
		h.logger.Error("Error al restaurar colaborador", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al restaurar colaborador"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "colaborador restaurado exitosamente"})
}

package handler

import (
	"errors"
	"net/http"
	"strconv"

	"PaginaSEG/internal/integrante"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IntegranteHandler struct {
	service *integrante.Service
	logger  *zap.Logger
}

// NewIntegranteHandler crea un nuevo controlador de integrantes
func NewIntegranteHandler(s *integrante.Service, l *zap.Logger) *IntegranteHandler {
	return &IntegranteHandler{
		service: s,
		logger:  l,
	}
}

// 1. Ver Lista de Integrantes de PostgreSQL en la Plantilla HTML
func (h *IntegranteHandler) Lista(c *gin.Context) {
	integrantes, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Error al obtener integrantes para plantilla Lista.html", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error al cargar la lista de integrantes")
		return
	}

	c.HTML(http.StatusOK, "Lista.html", gin.H{
		"Integrantes": integrantes,
	})
}

// 2. Mostrar Formulario de Crear
func (h *IntegranteHandler) Crear(c *gin.Context) {
	c.HTML(http.StatusOK, "Crear.html", nil)
}

// Procesar el Formulario Crear (POST)
func (h *IntegranteHandler) Insertar(c *gin.Context) {
	rolID, _ := strconv.Atoi(c.PostForm("rol_id"))

	req := integrante.Integrante{
		Nombre:          c.PostForm("nombre"),
		Apellido:        c.PostForm("apellido"),
		Contacto:        c.PostForm("contacto"),
		Especializacion: c.PostForm("especializacion"),
		Descripcion:     c.PostForm("descripcion"),
		RolID:           rolID,
	}

	if fileImg, err := c.FormFile("imagen"); err == nil && fileImg != nil {
		destPath := "ui/static/assets/img-GIS/IMGintegrantes/" + fileImg.Filename
		if err := c.SaveUploadedFile(fileImg, destPath); err == nil {
			req.Imagen = "../../static/assets/img-GIS/IMGintegrantes/" + fileImg.Filename
		}
	}

	if fileCV, err := c.FormFile("cv"); err == nil && fileCV != nil {
		destPath := "ui/static/assets/CVs-GIS/" + fileCV.Filename
		if err := c.SaveUploadedFile(fileCV, destPath); err == nil {
			req.CV = "../../static/assets/CVs-GIS/" + fileCV.Filename
		}
	}

	err := h.service.Create(&req)
	if err != nil {
		h.logger.Error("Error al crear integrante desde formulario HTML", zap.Error(err))
		c.HTML(http.StatusBadRequest, "Crear.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/integrantes")
}

// Mostrar Formulario de Editar precargado
func (h *IntegranteHandler) Editar(c *gin.Context) {
	idParam := c.Query("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.Redirect(http.StatusSeeOther, "/admin/integrantes")
		return
	}

	integranteObj, err := h.service.Read(id)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/admin/integrantes")
		return
	}

	c.HTML(http.StatusOK, "Editar.html", gin.H{
		"Integrante": integranteObj,
	})
}

// Procesar la Actualización (POST)
func (h *IntegranteHandler) Actualizar(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	if id <= 0 {
		c.Redirect(http.StatusSeeOther, "/admin/integrantes")
		return
	}

	rolID, _ := strconv.Atoi(c.PostForm("rol_id"))
	nombre := c.PostForm("nombre")
	apellido := c.PostForm("apellido")
	contacto := c.PostForm("contacto")
	especializacion := c.PostForm("especializacion")
	descripcion := c.PostForm("descripcion")

	fields := integrante.UpdateFields{
		Nombre:          &nombre,
		Apellido:        &apellido,
		Contacto:        &contacto,
		Especializacion: &especializacion,
		Descripcion:     &descripcion,
		RolID:           &rolID,
	}

	if fileImg, err := c.FormFile("imagen"); err == nil && fileImg != nil {
		destPath := "ui/static/assets/img-GIS/IMGintegrantes/" + fileImg.Filename
		if err := c.SaveUploadedFile(fileImg, destPath); err == nil {
			imgPath := "../../static/assets/img-GIS/IMGintegrantes/" + fileImg.Filename
			fields.Imagen = &imgPath
		}
	}

	if fileCV, err := c.FormFile("cv"); err == nil && fileCV != nil {
		destPath := "ui/static/assets/CVs-GIS/" + fileCV.Filename
		if err := c.SaveUploadedFile(fileCV, destPath); err == nil {
			cvPath := "../../static/assets/CVs-GIS/" + fileCV.Filename
			fields.CV = &cvPath
		}
	}

	err := h.service.Update(id, fields)
	if err != nil {
		h.logger.Error("Error al actualizar integrante desde formulario HTML", zap.Int("id", id), zap.Error(err))
		integranteObj, _ := h.service.Read(id)
		c.HTML(http.StatusBadRequest, "Editar.html", gin.H{
			"Error":      err.Error(),
			"Integrante": integranteObj,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/integrantes")
}

// Eliminar Registro
func (h *IntegranteHandler) Borrar(c *gin.Context) {
	idParam := c.Query("id")
	id, err := strconv.Atoi(idParam)
	if err == nil && id > 0 {
		if errDel := h.service.Delete(id); errDel != nil {
			h.logger.Error("Error al eliminar integrante desde la plantilla", zap.Int("id", id), zap.Error(errDel))
		}
	}

	c.Redirect(http.StatusSeeOther, "/admin/integrantes")
}

// REST API JSON ENDPOINTS (Opcionales para Postman)

func (ih *IntegranteHandler) API_GetAll(c *gin.Context) {
	integrantes, err := ih.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	if integrantes == nil {
		integrantes = []integrante.Integrante{}
	}
	c.JSON(http.StatusOK, integrantes)
}

func (ih *IntegranteHandler) API_Read(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El ID debe ser un número entero válido"})
		return
	}
	i, err := ih.service.Read(id)
	if err != nil {
		if errors.Is(err, integrante.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Integrante no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	c.JSON(http.StatusOK, i)
}

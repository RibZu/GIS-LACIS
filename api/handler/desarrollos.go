package handler

import (
	"net/http"
	"strconv"

	"PaginaSEG/internal/desarrollo"
	"PaginaSEG/internal/integrante"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DesarrolloHandler struct {
	service           *desarrollo.Service
	integranteService *integrante.Service
	logger            *zap.Logger
}

func NewDesarrolloHandler(s *desarrollo.Service, is *integrante.Service, l *zap.Logger) *DesarrolloHandler {
	return &DesarrolloHandler{service: s, integranteService: is, logger: l}
}

// DesarrolloConParticipantes es lo que consume ListaDesarrollos.html: el
// desarrollo con sus integrantes registrados y sus participantes externos
// ya resueltos, listos para pintar el acordeón.
type DesarrolloConParticipantes struct {
	desarrollo.Desarrollo
	Integrantes           []integrante.Integrante
	ParticipantesExternos []string
}

// 1. Ver "Gestión de Desarrollos" (acordeón)
func (h *DesarrolloHandler) Lista(c *gin.Context) {
	desarrollos, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Error al obtener desarrollos para plantilla ListaDesarrollos.html", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error al cargar la lista de desarrollos")
		return
	}

	var lista []DesarrolloConParticipantes
	for _, d := range desarrollos {
		vinculados, _ := h.service.ObtenerIntegrantes(d.ID)
		externos, _ := h.service.ObtenerParticipantesExternos(d.ID)
		lista = append(lista, DesarrolloConParticipantes{
			Desarrollo:            d,
			Integrantes:           vinculados,
			ParticipantesExternos: externos,
		})
	}

	c.HTML(http.StatusOK, "ListaDesarrollos.html", gin.H{
		"Desarrollos": lista,
		"LoggedIn":    true,
	})
}

// 2. Mostrar Formulario de Carga
func (h *DesarrolloHandler) Crear(c *gin.Context) {
	integrantes, _ := h.integranteService.GetAll()
	c.HTML(http.StatusOK, "CrearDesarrollo.html", gin.H{
		"Integrantes": integrantes,
		"LoggedIn":    true,
	})
}

// leerParticipantesDelForm interpreta los dos campos que genera
// participantesWidget.js: "integrantes" (IDs de integrantes ya registrados)
// y "externos" (nombres sueltos de participantes sin ficha).
func leerParticipantesDelForm(c *gin.Context) (integranteIDs []int, externos []string) {
	for _, idStr := range c.PostFormArray("integrantes") {
		if id, err := strconv.Atoi(idStr); err == nil {
			integranteIDs = append(integranteIDs, id)
		}
	}
	externos = c.PostFormArray("externos")
	return integranteIDs, externos
}

// Procesar el Formulario de Carga (POST)
func (h *DesarrolloHandler) Insertar(c *gin.Context) {
	anio, _ := strconv.Atoi(c.PostForm("anio"))

	req := desarrollo.Desarrollo{
		Titulo:      c.PostForm("titulo"),
		Anio:        anio,
		URL:         c.PostForm("url"),
		Contacto:    c.PostForm("contacto"),
		Descripcion: c.PostForm("descripcion"),
	}

	err := h.service.Create(&req)
	if err != nil {
		h.logger.Error("Error al crear desarrollo desde formulario HTML", zap.Error(err))
		integrantes, _ := h.integranteService.GetAll()
		c.HTML(http.StatusBadRequest, "CrearDesarrollo.html", gin.H{
			"Error":       err.Error(),
			"Integrantes": integrantes,
			"LoggedIn":    true,
		})
		return
	}

	integranteIDs, externos := leerParticipantesDelForm(c)
	if err := h.service.VincularIntegrantes(req.ID, integranteIDs); err != nil {
		h.logger.Error("Error al vincular integrantes al desarrollo", zap.Int("desarrollo_id", req.ID), zap.Error(err))
	}
	if err := h.service.VincularParticipantesExternos(req.ID, externos); err != nil {
		h.logger.Error("Error al vincular participantes externos al desarrollo", zap.Int("desarrollo_id", req.ID), zap.Error(err))
	}

	c.Redirect(http.StatusSeeOther, "/admin/desarrollos")
}

// Mostrar Formulario de Editar precargado
func (h *DesarrolloHandler) Editar(c *gin.Context) {
	idParam := c.Query("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.Redirect(http.StatusSeeOther, "/admin/desarrollos")
		return
	}

	d, err := h.service.Read(id)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/admin/desarrollos")
		return
	}

	todos, _ := h.integranteService.GetAll()
	vinculados, _ := h.service.ObtenerIntegrantes(id)
	externos, _ := h.service.ObtenerParticipantesExternos(id)

	c.HTML(http.StatusOK, "EditarDesarrollo.html", gin.H{
		"Desarrollo":            d,
		"Integrantes":           todos,
		"Vinculados":            vinculados,
		"ParticipantesExternos": externos,
		"LoggedIn":              true,
	})
}

// Procesar la Actualización (POST)
func (h *DesarrolloHandler) Actualizar(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	if id <= 0 {
		c.Redirect(http.StatusSeeOther, "/admin/desarrollos")
		return
	}

	titulo := c.PostForm("titulo")
	anio, _ := strconv.Atoi(c.PostForm("anio"))
	url := c.PostForm("url")
	contacto := c.PostForm("contacto")
	descripcion := c.PostForm("descripcion")

	fields := desarrollo.UpdateFields{
		Titulo:      &titulo,
		Anio:        &anio,
		URL:         &url,
		Contacto:    &contacto,
		Descripcion: &descripcion,
	}

	err := h.service.Update(id, fields)
	if err != nil {
		h.logger.Error("Error al actualizar desarrollo desde formulario HTML", zap.Int("id", id), zap.Error(err))
		d, _ := h.service.Read(id)
		todos, _ := h.integranteService.GetAll()
		vinculados, _ := h.service.ObtenerIntegrantes(id)
		externos, _ := h.service.ObtenerParticipantesExternos(id)
		c.HTML(http.StatusBadRequest, "EditarDesarrollo.html", gin.H{
			"Error":                 err.Error(),
			"Desarrollo":            d,
			"Integrantes":           todos,
			"Vinculados":            vinculados,
			"ParticipantesExternos": externos,
			"LoggedIn":              true,
		})
		return
	}

	integranteIDs, externos := leerParticipantesDelForm(c)
	if err := h.service.VincularIntegrantes(id, integranteIDs); err != nil {
		h.logger.Error("Error al vincular integrantes al desarrollo", zap.Int("desarrollo_id", id), zap.Error(err))
	}
	if err := h.service.VincularParticipantesExternos(id, externos); err != nil {
		h.logger.Error("Error al vincular participantes externos al desarrollo", zap.Int("desarrollo_id", id), zap.Error(err))
	}

	c.Redirect(http.StatusSeeOther, "/admin/desarrollos")
}

// Eliminar Registro
func (h *DesarrolloHandler) Borrar(c *gin.Context) {
	idParam := c.Query("id")
	id, err := strconv.Atoi(idParam)
	if err == nil && id > 0 {
		if errDel := h.service.Delete(id); errDel != nil {
			h.logger.Error("Error al eliminar desarrollo desde la plantilla", zap.Int("id", id), zap.Error(errDel))
		}
	}
	c.Redirect(http.StatusSeeOther, "/admin/desarrollos")
}

// ViewPublica renderiza la vista pública de desarrollos en Desarrollos.html
func (h *DesarrolloHandler) ViewPublica(c *gin.Context) {
	desarrollos, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Error al obtener desarrollos para vista pública Desarrollos.html", zap.Error(err))
		desarrollos = []desarrollo.Desarrollo{}
	}

	_, loggedIn := CurrentUserID(c)

	c.HTML(http.StatusOK, "Desarrollos.html", gin.H{
		"Desarrollos": desarrollos,
		"LoggedIn":    loggedIn,
	})
}

func (h *DesarrolloHandler) ViewPublica(c *gin.Context) {
	lista, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Error al obtener desarrollos para la página pública", zap.Error(err))
		lista = []desarrollo.Desarrollo{}
	}

	_, loggedIn := CurrentUserID(c)

	c.HTML(http.StatusOK, "Desarrollos.html", gin.H{
		"Desarrollos": lista,
		"LoggedIn":    loggedIn,
	})
}

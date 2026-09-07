package handler

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"PaginaSEG/internal/integrante"
	"PaginaSEG/internal/tesis"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TesisHandler struct {
	service           *tesis.Service
	integranteService *integrante.Service
	logger            *zap.Logger
}

func NewTesisHandler(s *tesis.Service, is *integrante.Service, l *zap.Logger) *TesisHandler {
	return &TesisHandler{
		service:           s,
		integranteService: is,
		logger:            l,
	}
}

// 1. Ver Lista de Tesis en la Plantilla HTML
func (h *TesisHandler) Lista(c *gin.Context) {
	lista, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Error al obtener tesis para plantilla ListaTesis.html", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error al cargar la lista de tesis")
		return
	}

	c.HTML(http.StatusOK, "ListaTesis.html", gin.H{
		"Tesis":    lista,
		"LoggedIn": true,
	})
}

// 2. Mostrar Formulario de Crear
func (h *TesisHandler) Crear(c *gin.Context) {
	integrantes, err := h.integranteService.GetAll()
	if err != nil {
		h.logger.Warn("No se pudieron cargar integrantes para el selector de autoría", zap.Error(err))
		integrantes = []integrante.Integrante{}
	}

	c.HTML(http.StatusOK, "CrearTesis.html", gin.H{
		"Integrantes": integrantes,
		"LoggedIn":    true,
	})
}

// 3. Procesar el Formulario Crear (POST)
func (h *TesisHandler) Insertar(c *gin.Context) {
	titulo := strings.TrimSpace(c.PostForm("titulo"))
	nivel := strings.TrimSpace(c.PostForm("nivel"))
	carrera := strings.TrimSpace(c.PostForm("carrera_origen"))
	palabrasClave := strings.TrimSpace(c.PostForm("palabras_clave"))
	resumen := strings.TrimSpace(c.PostForm("resumen"))

	var anioPtr *int
	if anioVal, err := strconv.Atoi(c.PostForm("anio")); err == nil && anioVal > 0 {
		anioPtr = &anioVal
	}

	var autorIDPtr *int
	autorTipo := c.PostForm("autor_tipo")
	autorHistorico := strings.TrimSpace(c.PostForm("autor_historico"))
	if autorTipo == "registrado" {
		if val, err := strconv.Atoi(c.PostForm("autor_id")); err == nil && val > 0 {
			autorIDPtr = &val
			autorHistorico = ""
		}
	}

	var dirIDPtr *int
	dirTipo := c.PostForm("director_tipo")
	dirHistorico := strings.TrimSpace(c.PostForm("director_historico"))
	if dirTipo == "registrado" {
		if val, err := strconv.Atoi(c.PostForm("director_id")); err == nil && val > 0 {
			dirIDPtr = &val
			dirHistorico = ""
		}
	}

	var coodirIDPtr *int
	coodirTipo := c.PostForm("coodirector_tipo")
	coodirHistorico := strings.TrimSpace(c.PostForm("coodirector_historico"))
	if coodirTipo == "registrado" {
		if val, err := strconv.Atoi(c.PostForm("coodirector_id")); err == nil && val > 0 {
			coodirIDPtr = &val
			coodirHistorico = ""
		}
	}

	req := tesis.Tesis{
		Titulo:               titulo,
		Nivel:                nivel,
		CarreraOrigen:        carrera,
		Anio:                 anioPtr,
		PalabrasClave:        palabrasClave,
		Resumen:              resumen,
		AutorID:              autorIDPtr,
		AutorHistorico:       autorHistorico,
		DirectorID:           dirIDPtr,
		DirectorHistorico:    dirHistorico,
		CoodirectorID:        coodirIDPtr,
		CoodirectorHistorico: coodirHistorico,
	}

	// Manejo de archivo PDF
	if filePDF, err := c.FormFile("archivo_pdf"); err == nil && filePDF != nil {
		destDir := "ui/static/assets/tesis-GIS"
		_ = os.MkdirAll(destDir, os.ModePerm)
		filename := filepath.Base(filePDF.Filename)
		destPath := filepath.Join(destDir, filename)
		if err := c.SaveUploadedFile(filePDF, destPath); err == nil {
			req.ArchivoPDF = "/ui/static/assets/tesis-GIS/" + filename
		} else {
			h.logger.Error("Error al guardar archivo PDF de tesis", zap.Error(err))
		}
	}

	err := h.service.Create(&req)
	if err != nil {
		h.logger.Error("Error al crear tesis desde formulario HTML", zap.Error(err))
		integrantes, _ := h.integranteService.GetAll()
		c.HTML(http.StatusBadRequest, "CrearTesis.html", gin.H{
			"Error":       err.Error(),
			"Integrantes": integrantes,
			"Tesis":       req,
			"LoggedIn":    true,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/tesis")
}

// 4. Mostrar Formulario de Editar
func (h *TesisHandler) Editar(c *gin.Context) {
	idParam := c.Query("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.Redirect(http.StatusSeeOther, "/admin/tesis")
		return
	}

	tesisObj, err := h.service.Read(id)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/admin/tesis")
		return
	}

	integrantes, _ := h.integranteService.GetAll()

	c.HTML(http.StatusOK, "EditarTesis.html", gin.H{
		"Tesis":       tesisObj,
		"Integrantes": integrantes,
		"LoggedIn":    true,
	})
}

// 5. Procesar la Actualización (POST)
func (h *TesisHandler) Actualizar(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	if id <= 0 {
		c.Redirect(http.StatusSeeOther, "/admin/tesis")
		return
	}

	titulo := strings.TrimSpace(c.PostForm("titulo"))
	nivel := strings.TrimSpace(c.PostForm("nivel"))
	carrera := strings.TrimSpace(c.PostForm("carrera_origen"))
	palabrasClave := strings.TrimSpace(c.PostForm("palabras_clave"))
	resumen := strings.TrimSpace(c.PostForm("resumen"))

	var anioPtr *int
	if anioVal, err := strconv.Atoi(c.PostForm("anio")); err == nil && anioVal > 0 {
		anioPtr = &anioVal
	} else {
		valZero := 0
		anioPtr = &valZero
	}

	autorTipo := c.PostForm("autor_tipo")
	autorHistorico := strings.TrimSpace(c.PostForm("autor_historico"))
	var autorIDPtr *int
	if autorTipo == "registrado" {
		if val, err := strconv.Atoi(c.PostForm("autor_id")); err == nil && val > 0 {
			autorIDPtr = &val
			autorHistorico = ""
		}
	} else {
		valZero := 0
		autorIDPtr = &valZero
	}

	dirTipo := c.PostForm("director_tipo")
	dirHistorico := strings.TrimSpace(c.PostForm("director_historico"))
	var dirIDPtr *int
	if dirTipo == "registrado" {
		if val, err := strconv.Atoi(c.PostForm("director_id")); err == nil && val > 0 {
			dirIDPtr = &val
			dirHistorico = ""
		}
	} else {
		valZero := 0
		dirIDPtr = &valZero
	}

	coodirTipo := c.PostForm("coodirector_tipo")
	coodirHistorico := strings.TrimSpace(c.PostForm("coodirector_historico"))
	var coodirIDPtr *int
	if coodirTipo == "registrado" {
		if val, err := strconv.Atoi(c.PostForm("coodirector_id")); err == nil && val > 0 {
			coodirIDPtr = &val
			coodirHistorico = ""
		}
	} else {
		valZero := 0
		coodirIDPtr = &valZero
	}

	fields := tesis.UpdateFields{
		Titulo:               &titulo,
		Nivel:                &nivel,
		CarreraOrigen:        &carrera,
		Anio:                 anioPtr,
		PalabrasClave:        &palabrasClave,
		Resumen:              &resumen,
		AutorID:              autorIDPtr,
		AutorHistorico:       &autorHistorico,
		DirectorID:           dirIDPtr,
		DirectorHistorico:    &dirHistorico,
		CoodirectorID:        coodirIDPtr,
		CoodirectorHistorico: &coodirHistorico,
	}

	// Si se subió un nuevo PDF
	if filePDF, err := c.FormFile("archivo_pdf"); err == nil && filePDF != nil {
		destDir := "ui/static/assets/tesis-GIS"
		_ = os.MkdirAll(destDir, os.ModePerm)
		filename := filepath.Base(filePDF.Filename)
		destPath := filepath.Join(destDir, filename)
		if err := c.SaveUploadedFile(filePDF, destPath); err == nil {
			pdfPath := "/ui/static/assets/tesis-GIS/" + filename
			fields.ArchivoPDF = &pdfPath
		}
	}

	err := h.service.Update(id, fields)
	if err != nil {
		h.logger.Error("Error al actualizar tesis desde formulario HTML", zap.Int("id", id), zap.Error(err))
		tesisObj, _ := h.service.Read(id)
		integrantes, _ := h.integranteService.GetAll()
		c.HTML(http.StatusBadRequest, "EditarTesis.html", gin.H{
			"Error":       err.Error(),
			"Tesis":       tesisObj,
			"Integrantes": integrantes,
			"LoggedIn":    true,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/tesis")
}

// 6. Eliminar Registro
func (h *TesisHandler) Borrar(c *gin.Context) {
	idParam := c.Query("id")
	id, err := strconv.Atoi(idParam)
	if err == nil && id > 0 {
		if errDel := h.service.Delete(id); errDel != nil {
			h.logger.Error("Error al eliminar tesis desde la plantilla", zap.Int("id", id), zap.Error(errDel))
		}
	}

	c.Redirect(http.StatusSeeOther, "/admin/tesis")
}

// REST API JSON ENDPOINTS
func (h *TesisHandler) API_GetAll(c *gin.Context) {
	carrera := c.Query("carrera")
	nivel := c.Query("nivel")

	var lista []tesis.Tesis
	var err error

	if carrera != "" {
		lista, err = h.service.GetByCarrera(carrera)
	} else if nivel != "" {
		lista, err = h.service.GetByNivel(nivel)
	} else {
		lista, err = h.service.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al consultar tesis"})
		return
	}
	if lista == nil {
		lista = []tesis.Tesis{}
	}
	c.JSON(http.StatusOK, lista)
}

func (h *TesisHandler) API_Read(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El ID debe ser un número entero válido"})
		return
	}
	t, err := h.service.Read(id)
	if err != nil {
		if errors.Is(err, tesis.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tesis no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	c.JSON(http.StatusOK, t)
}

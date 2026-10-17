package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"PaginaSEG/internal/desarrollo"
	"PaginaSEG/internal/integrante"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	urlListaDesarrollos       = "/admin/desarrollos"
	mensajeDesarrolloGenerico = "No se pudo guardar. Probá de nuevo en unos minutos."
	mensajeLecturaProductos   = "No pudimos cargar los productos. Probá de nuevo en unos minutos."
)

func mensajeYEstadoDesarrollo(err error) (string, int) {
	switch {
	case errors.Is(err, desarrollo.ErrAnioInvalido):
		return "Ingresá un año válido.", http.StatusBadRequest
	case errors.Is(err, desarrollo.ErrTituloRequerido):
		return "Ingresá un título.", http.StatusBadRequest
	case errors.Is(err, desarrollo.ErrTituloLargo):
		return "El título no puede superar los 255 caracteres.", http.StatusBadRequest
	case errors.Is(err, desarrollo.ErrURLLarga):
		return "El enlace no puede superar los 500 caracteres.", http.StatusBadRequest
	case errors.Is(err, desarrollo.ErrContactoLargo):
		return "El contacto no puede superar los 255 caracteres.", http.StatusBadRequest
	case errors.Is(err, desarrollo.ErrParticipanteLargo):
		return "Cada participante externo puede tener hasta 150 caracteres.", http.StatusBadRequest
	case errors.Is(err, desarrollo.ErrIntegranteInexistente):
		return "Alguno de los integrantes elegidos ya no existe. Revisá la lista de participantes.", http.StatusBadRequest
	}
	return mensajeDesarrolloGenerico, http.StatusInternalServerError
}

type DesarrolloHandler struct {
	service           *desarrollo.Service
	integranteService *integrante.Service
	logger            *zap.Logger
}

func NewDesarrolloHandler(s *desarrollo.Service, is *integrante.Service, l *zap.Logger) *DesarrolloHandler {
	return &DesarrolloHandler{service: s, integranteService: is, logger: l}
}

type DesarrolloConParticipantes struct {
	desarrollo.Desarrollo
	Integrantes           []integrante.Integrante
	ParticipantesExternos []string
}

type formularioDesarrollo struct {
	Desarrollo    desarrollo.Desarrollo
	AnioTexto     string
	AnioNumerico  bool
	IntegranteIDs []int
	Externos      []string
}

func leerParticipantesDelForm(c *gin.Context) (integranteIDs []int, externos []string) {
	for _, idStr := range c.PostFormArray("integrantes") {
		if id, err := strconv.Atoi(idStr); err == nil {
			integranteIDs = append(integranteIDs, id)
		}
	}
	for _, nombre := range c.PostFormArray("externos") {
		if strings.TrimSpace(nombre) != "" {
			externos = append(externos, nombre)
		}
	}
	return integranteIDs, externos
}

func leerFormularioDesarrollo(c *gin.Context) formularioDesarrollo {
	texto := strings.TrimSpace(c.PostForm("anio"))
	anio, err := strconv.Atoi(texto)
	ids, externos := leerParticipantesDelForm(c)
	return formularioDesarrollo{
		Desarrollo: desarrollo.Desarrollo{
			Titulo:      c.PostForm("titulo"),
			Anio:        anio,
			URL:         c.PostForm("url"),
			Contacto:    c.PostForm("contacto"),
			Descripcion: c.PostForm("descripcion"),
		},
		AnioTexto:     texto,
		AnioNumerico:  err == nil,
		IntegranteIDs: ids,
		Externos:      externos,
	}
}

func vinculadosDesdeIDs(todos []integrante.Integrante, ids []int) []integrante.Integrante {
	porID := make(map[int]integrante.Integrante, len(todos))
	for _, it := range todos {
		porID[it.ID] = it
	}
	var vinculados []integrante.Integrante
	vistos := map[int]bool{}
	for _, id := range ids {
		if it, existe := porID[id]; existe && !vistos[id] {
			vistos[id] = true
			vinculados = append(vinculados, it)
		}
	}
	return vinculados
}

func (h *DesarrolloHandler) integrantesParaFormulario(c *gin.Context) ([]integrante.Integrante, bool) {
	todos, err := h.integranteService.GetAll()
	if err != nil {
		h.logger.Error("Error al obtener los integrantes para el formulario de productos", zap.Error(err))
		paginaError(c, http.StatusInternalServerError, "No pudimos cargar el formulario", "No pudimos cargar la lista de integrantes. Probá de nuevo en unos minutos.")
		return nil, false
	}
	return todos, true
}

func (h *DesarrolloHandler) productosConParticipantes() ([]DesarrolloConParticipantes, error) {
	desarrollos, err := h.service.GetAll()
	if err != nil {
		return nil, err
	}
	integrantes, externos, err := h.service.ParticipantesPorDesarrollo()
	if err != nil {
		return nil, err
	}
	lista := make([]DesarrolloConParticipantes, 0, len(desarrollos))
	for _, d := range desarrollos {
		lista = append(lista, DesarrolloConParticipantes{
			Desarrollo:            d,
			Integrantes:           integrantes[d.ID],
			ParticipantesExternos: externos[d.ID],
		})
	}
	return lista, nil
}

func (h *DesarrolloHandler) Lista(c *gin.Context) {
	avisoOK, avisoError := avisosDesdeQuery(c)

	lista, err := h.productosConParticipantes()
	if err != nil {
		h.logger.Error("Error al obtener productos para plantilla ListaDesarrollos.html", zap.Error(err))
		c.HTML(http.StatusInternalServerError, "ListaDesarrollos.html", gin.H{
			"Error":    mensajeLecturaProductos,
			"LoggedIn": true,
		})
		return
	}

	c.HTML(http.StatusOK, "ListaDesarrollos.html", gin.H{
		"Desarrollos": lista,
		"LoggedIn":    true,
		"Aviso":       avisoOK,
		"ErrorAviso":  avisoError,
	})
}

func (h *DesarrolloHandler) Crear(c *gin.Context) {
	integrantes, ok := h.integrantesParaFormulario(c)
	if !ok {
		return
	}
	c.HTML(http.StatusOK, "CrearDesarrollo.html", gin.H{
		"Integrantes": integrantes,
		"LoggedIn":    true,
	})
}

func (h *DesarrolloHandler) rechazarCreacion(c *gin.Context, err error, form formularioDesarrollo) {
	mensaje, estado := mensajeYEstadoDesarrollo(err)
	if estado == http.StatusBadRequest {
		h.logger.Warn("Alta de producto rechazada", zap.String("motivo", mensaje))
	}
	todos, ok := h.integrantesParaFormulario(c)
	if !ok {
		return
	}
	d := form.Desarrollo
	c.HTML(estado, "CrearDesarrollo.html", gin.H{
		"Error":                 mensaje,
		"Desarrollo":            &d,
		"AnioTexto":             form.AnioTexto,
		"Integrantes":           todos,
		"Vinculados":            vinculadosDesdeIDs(todos, form.IntegranteIDs),
		"ParticipantesExternos": form.Externos,
		"LoggedIn":              true,
	})
}

func (h *DesarrolloHandler) Insertar(c *gin.Context) {
	form := leerFormularioDesarrollo(c)
	if !form.AnioNumerico {
		h.rechazarCreacion(c, desarrollo.ErrAnioInvalido, form)
		return
	}

	req := form.Desarrollo
	if err := h.service.CrearCompleto(&req, form.IntegranteIDs, form.Externos); err != nil {
		h.rechazarCreacion(c, err, form)
		return
	}

	c.Redirect(http.StatusSeeOther, urlListaDesarrollos+"?ok=producto-creado")
}

func (h *DesarrolloHandler) Editar(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		c.Redirect(http.StatusSeeOther, urlListaDesarrollos+"?error=no-encontrado")
		return
	}

	d, err := h.service.Read(id)
	if err != nil {
		if errors.Is(err, desarrollo.ErrNotFound) || errors.Is(err, desarrollo.ErrIDInvalido) {
			c.Redirect(http.StatusSeeOther, urlListaDesarrollos+"?error=no-encontrado")
			return
		}
		paginaError(c, http.StatusInternalServerError, "No pudimos cargar el producto", "No pudimos cargar el producto. Probá de nuevo en unos minutos.")
		return
	}

	todos, ok := h.integrantesParaFormulario(c)
	if !ok {
		return
	}
	vinculados, err := h.service.ObtenerIntegrantes(id)
	if err != nil {
		h.logger.Error("Error al obtener los integrantes del producto", zap.Int("id", id), zap.Error(err))
		paginaError(c, http.StatusInternalServerError, "No pudimos cargar el producto", "No pudimos cargar el producto. Probá de nuevo en unos minutos.")
		return
	}
	externos, err := h.service.ObtenerParticipantesExternos(id)
	if err != nil {
		h.logger.Error("Error al obtener los participantes externos del producto", zap.Int("id", id), zap.Error(err))
		paginaError(c, http.StatusInternalServerError, "No pudimos cargar el producto", "No pudimos cargar el producto. Probá de nuevo en unos minutos.")
		return
	}

	c.HTML(http.StatusOK, "EditarDesarrollo.html", gin.H{
		"Desarrollo":            d,
		"Integrantes":           todos,
		"Vinculados":            vinculados,
		"ParticipantesExternos": externos,
		"LoggedIn":              true,
	})
}

func (h *DesarrolloHandler) rechazarEdicion(c *gin.Context, err error, id int, form formularioDesarrollo) {
	mensaje, estado := mensajeYEstadoDesarrollo(err)
	if estado == http.StatusBadRequest {
		h.logger.Warn("Edición de producto rechazada", zap.Int("id", id), zap.String("motivo", mensaje))
	}
	todos, ok := h.integrantesParaFormulario(c)
	if !ok {
		return
	}
	d := form.Desarrollo
	d.ID = id
	c.HTML(estado, "EditarDesarrollo.html", gin.H{
		"Error":                 mensaje,
		"Desarrollo":            &d,
		"AnioTexto":             form.AnioTexto,
		"Integrantes":           todos,
		"Vinculados":            vinculadosDesdeIDs(todos, form.IntegranteIDs),
		"ParticipantesExternos": form.Externos,
		"LoggedIn":              true,
	})
}

func (h *DesarrolloHandler) Actualizar(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	if id <= 0 {
		c.Redirect(http.StatusSeeOther, urlListaDesarrollos+"?error=no-encontrado")
		return
	}

	form := leerFormularioDesarrollo(c)
	if !form.AnioNumerico {
		h.rechazarEdicion(c, desarrollo.ErrAnioInvalido, id, form)
		return
	}

	fields := desarrollo.UpdateFields{
		Titulo:      &form.Desarrollo.Titulo,
		Anio:        &form.Desarrollo.Anio,
		URL:         &form.Desarrollo.URL,
		Contacto:    &form.Desarrollo.Contacto,
		Descripcion: &form.Desarrollo.Descripcion,
	}

	if err := h.service.ActualizarCompleto(id, fields, form.IntegranteIDs, form.Externos); err != nil {
		if errors.Is(err, desarrollo.ErrNotFound) {
			c.Redirect(http.StatusSeeOther, urlListaDesarrollos+"?error=no-encontrado")
			return
		}
		h.rechazarEdicion(c, err, id, form)
		return
	}

	c.Redirect(http.StatusSeeOther, urlListaDesarrollos+"?ok=producto-actualizado")
}

func (h *DesarrolloHandler) Borrar(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil || id <= 0 {
		c.Redirect(http.StatusSeeOther, urlListaDesarrollos+"?error=no-encontrado")
		return
	}

	err = h.service.Delete(id)
	switch {
	case err == nil:
		c.Redirect(http.StatusSeeOther, urlListaDesarrollos+"?ok=producto-eliminado")
	case errors.Is(err, desarrollo.ErrNotFound), errors.Is(err, desarrollo.ErrIDInvalido):
		c.Redirect(http.StatusSeeOther, urlListaDesarrollos+"?error=no-encontrado")
	default:
		h.logger.Error("Error al eliminar producto desde la lista", zap.Int("id", id), zap.Error(err))
		c.Redirect(http.StatusSeeOther, urlListaDesarrollos+"?error=error-eliminar")
	}
}

func (h *DesarrolloHandler) BorrarGet(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, urlListaDesarrollos)
}

func (h *DesarrolloHandler) ViewLacis(c *gin.Context) {
	lista, err := h.productosConParticipantes()
	_, loggedIn := CurrentUserID(c)

	datos := gin.H{"LoggedIn": loggedIn}
	if err != nil {
		h.logger.Error("Error al obtener productos para la sección de productos de software en Lacis.html", zap.Error(err))
		datos["ErrorProductos"] = true
	} else {
		datos["Desarrollos"] = lista
	}

	c.HTML(http.StatusOK, "Lacis.html", datos)
}

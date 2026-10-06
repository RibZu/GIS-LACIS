package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"PaginaSEG/internal/usuario"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var errContrasenasDistintas = errors.New("las contraseñas no coinciden")

const (
	urlListaUsuarios       = "/admin/usuarios"
	mensajeUsuarioGenerico = "No se pudo guardar. Probá de nuevo en unos minutos."
)

func mensajeYEstadoUsuario(err error) (string, int) {
	switch {
	case errors.Is(err, usuario.ErrUsernameRequerido):
		return "Ingresá un nombre de usuario.", http.StatusBadRequest
	case errors.Is(err, usuario.ErrUsernameInvalido):
		return "El nombre de usuario no puede tener espacios.", http.StatusBadRequest
	case errors.Is(err, usuario.ErrEmailRequerido):
		return "Ingresá un email.", http.StatusBadRequest
	case errors.Is(err, usuario.ErrPasswordRequerido):
		return "Ingresá una contraseña.", http.StatusBadRequest
	case errors.Is(err, errContrasenasDistintas):
		return "Las contraseñas no coinciden.", http.StatusBadRequest
	case errors.Is(err, usuario.ErrModulosRequerido):
		return "Elegí al menos un módulo.", http.StatusBadRequest
	case errors.Is(err, usuario.ErrUsernameDuplicado):
		return "Ya existe un usuario con ese nombre de usuario.", http.StatusBadRequest
	case errors.Is(err, usuario.ErrEmailDuplicado):
		return "Ya existe un usuario con ese email.", http.StatusBadRequest
	case errors.Is(err, usuario.ErrAutoDegradacion):
		return "No podés quitarte tu propio rol de administrador.", http.StatusBadRequest
	case errors.Is(err, usuario.ErrUltimoAdmin):
		return "Tiene que quedar al menos un administrador.", http.StatusBadRequest
	}
	return mensajeUsuarioGenerico, http.StatusInternalServerError
}

func seleccionadosDeLista(modulos []string) map[string]bool {
	seleccionados := map[string]bool{}
	for _, m := range modulos {
		if m = strings.TrimSpace(m); m != "" {
			seleccionados[m] = true
		}
	}
	return seleccionados
}

func seleccionadosDeTexto(modulos *string) map[string]bool {
	if modulos == nil {
		return map[string]bool{}
	}
	return seleccionadosDeLista(strings.Split(*modulos, ","))
}

type UsuarioHandler struct {
	service *usuario.Service
	logger  *zap.Logger
}

func NewUsuarioHandler(s *usuario.Service, l *zap.Logger) *UsuarioHandler {
	return &UsuarioHandler{
		service: s,
		logger:  l,
	}
}

func (h *UsuarioHandler) Lista(c *gin.Context) {
	avisoOK, avisoError := avisosDesdeQuery(c)

	usuarios, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Error al obtener usuarios para plantilla ListaUsuarios.html", zap.Error(err))
		c.HTML(http.StatusInternalServerError, "ListaUsuarios.html", gin.H{
			"Error":    "No pudimos cargar los usuarios. Probá de nuevo en unos minutos.",
			"LoggedIn": true,
		})
		return
	}

	currentID, _ := CurrentUserID(c)

	c.HTML(http.StatusOK, "ListaUsuarios.html", gin.H{
		"Usuarios":      usuarios,
		"CurrentUserID": currentID,
		"LoggedIn":      true,
		"Aviso":         avisoOK,
		"ErrorAviso":    avisoError,
	})
}

func (h *UsuarioHandler) Crear(c *gin.Context) {
	c.HTML(http.StatusOK, "CrearUsuario.html", gin.H{
		"Seleccionados": map[string]bool{},
		"LoggedIn":      true,
	})
}

func (h *UsuarioHandler) rechazarCreacion(c *gin.Context, err error, username, email string, modulos []string) {
	mensaje, estado := mensajeYEstadoUsuario(err)
	if estado == http.StatusBadRequest {
		h.logger.Warn("Alta de usuario rechazada", zap.String("motivo", mensaje))
	}
	c.HTML(estado, "CrearUsuario.html", gin.H{
		"Error":         mensaje,
		"Username":      username,
		"Email":         email,
		"Seleccionados": seleccionadosDeLista(modulos),
		"LoggedIn":      true,
	})
}

func (h *UsuarioHandler) Insertar(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	modulosElegidos := c.PostFormArray("modulos")

	if password != c.PostForm("password_confirm") {
		h.rechazarCreacion(c, errContrasenasDistintas, username, email, modulosElegidos)
		return
	}

	modulos := strings.Join(modulosElegidos, ",")
	req := usuario.UsuarioGestor{
		Username:     &username,
		PasswordHash: &password,
		Email:        &email,
		Modulos:      &modulos,
	}

	if err := h.service.Create(&req); err != nil {
		h.rechazarCreacion(c, err, username, email, modulosElegidos)
		return
	}

	c.Redirect(http.StatusSeeOther, urlListaUsuarios+"?ok=usuario-creado")
}

func (h *UsuarioHandler) Editar(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id <= 0 {
		c.Redirect(http.StatusSeeOther, urlListaUsuarios+"?error=no-encontrado")
		return
	}

	usuarioObj, err := h.service.Read(id)
	if err != nil {
		if errors.Is(err, usuario.ErrNotFound) || errors.Is(err, usuario.ErrIDInvalido) {
			c.Redirect(http.StatusSeeOther, urlListaUsuarios+"?error=no-encontrado")
			return
		}
		paginaError(c, http.StatusInternalServerError, "No pudimos cargar el usuario", "No pudimos cargar el usuario. Probá de nuevo en unos minutos.")
		return
	}

	c.HTML(http.StatusOK, "EditarUsuario.html", gin.H{
		"Usuario":       usuarioObj,
		"Seleccionados": seleccionadosDeTexto(usuarioObj.Modulos),
		"EsAdmin":       esAdmin(usuarioObj),
		"LoggedIn":      true,
	})
}

func (h *UsuarioHandler) rechazarEdicion(c *gin.Context, err error, destino *usuario.UsuarioGestor, username, email string, modulos []string) {
	mensaje, estado := mensajeYEstadoUsuario(err)
	if estado == http.StatusBadRequest {
		h.logger.Warn("Edición de usuario rechazada", zap.Int("id", destino.ID), zap.String("motivo", mensaje))
	}
	copia := *destino
	copia.Username = &username
	copia.Email = &email
	seleccionados := seleccionadosDeTexto(destino.Modulos)
	if !esAdmin(destino) {
		seleccionados = seleccionadosDeLista(modulos)
	}
	c.HTML(estado, "EditarUsuario.html", gin.H{
		"Error":         mensaje,
		"Usuario":       &copia,
		"Seleccionados": seleccionados,
		"EsAdmin":       esAdmin(destino),
		"LoggedIn":      true,
	})
}

func (h *UsuarioHandler) Actualizar(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	if id <= 0 {
		c.Redirect(http.StatusSeeOther, urlListaUsuarios+"?error=no-encontrado")
		return
	}

	actor, ok := currentUsuario(c, h.service)
	if !ok {
		redirigirAlLogin(c)
		return
	}

	destino, err := h.service.Read(id)
	if err != nil {
		if errors.Is(err, usuario.ErrNotFound) {
			c.Redirect(http.StatusSeeOther, urlListaUsuarios+"?error=no-encontrado")
			return
		}
		paginaError(c, http.StatusInternalServerError, "No pudimos cargar el usuario", "No pudimos cargar el usuario. Probá de nuevo en unos minutos.")
		return
	}

	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	modulosElegidos := c.PostFormArray("modulos")

	if password != "" && password != c.PostForm("password_confirm") {
		h.rechazarEdicion(c, errContrasenasDistintas, destino, username, email, modulosElegidos)
		return
	}

	fields := usuario.UpdateFieldGestor{
		Username: &username,
		Email:    &email,
	}

	if !(esAdmin(destino) && len(modulosElegidos) == 0) {
		modulos := strings.Join(modulosElegidos, ",")
		fields.Modulos = &modulos
	}

	if password != "" {
		fields.PasswordHash = &password
	}

	if err := h.service.Update(actor.ID, id, &fields); err != nil {
		if errors.Is(err, usuario.ErrNotFound) {
			c.Redirect(http.StatusSeeOther, urlListaUsuarios+"?error=no-encontrado")
			return
		}
		h.rechazarEdicion(c, err, destino, username, email, modulosElegidos)
		return
	}

	c.Redirect(http.StatusSeeOther, urlListaUsuarios+"?ok=usuario-actualizado")
}

func (h *UsuarioHandler) Borrar(c *gin.Context) {
	actor, ok := currentUsuario(c, h.service)
	if !ok {
		redirigirAlLogin(c)
		return
	}

	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil || id <= 0 {
		c.Redirect(http.StatusSeeOther, urlListaUsuarios+"?error=no-encontrado")
		return
	}

	err = h.service.Delete(actor.ID, id)
	switch {
	case err == nil:
		c.Redirect(http.StatusSeeOther, urlListaUsuarios+"?ok=usuario-eliminado")
	case errors.Is(err, usuario.ErrAutoEliminacion):
		h.logger.Warn("Intento de autoeliminación bloqueado", zap.Int("id", id))
		c.Redirect(http.StatusSeeOther, urlListaUsuarios+"?error=autoeliminacion")
	case errors.Is(err, usuario.ErrUltimoAdmin):
		h.logger.Warn("Intento de borrar al último administrador bloqueado", zap.Int("id", id))
		c.Redirect(http.StatusSeeOther, urlListaUsuarios+"?error=ultimo-admin")
	case errors.Is(err, usuario.ErrNotFound), errors.Is(err, usuario.ErrIDInvalido):
		c.Redirect(http.StatusSeeOther, urlListaUsuarios+"?error=no-encontrado")
	default:
		h.logger.Error("Error al eliminar usuario desde la lista", zap.Int("id", id), zap.Error(err))
		c.Redirect(http.StatusSeeOther, urlListaUsuarios+"?error=error-eliminar")
	}
}

func (h *UsuarioHandler) BorrarGet(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, urlListaUsuarios)
}

package handler

import (
	"net/http"
	"strconv"
	"strings"

	"PaginaSEG/internal/usuario"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UsuarioHandler struct {
	service *usuario.Service
	logger  *zap.Logger
}

// NewUsuarioHandler crea un nuevo controlador de usuarios gestores
func NewUsuarioHandler(s *usuario.Service, l *zap.Logger) *UsuarioHandler {
	return &UsuarioHandler{
		service: s,
		logger:  l,
	}
}

// 1. Ver Lista de Usuarios Gestores en la Plantilla HTML
func (h *UsuarioHandler) Lista(c *gin.Context) {
	usuarios, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("Error al obtener usuarios para plantilla ListaUsuarios.html", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error al cargar la lista de usuarios")
		return
	}

	currentID, _ := CurrentUserID(c)

	c.HTML(http.StatusOK, "ListaUsuarios.html", gin.H{
		"Usuarios":      usuarios,
		"CurrentUserID": currentID,
		"LoggedIn":      true,
	})
}

// 2. Mostrar Formulario de Crear
func (h *UsuarioHandler) Crear(c *gin.Context) {
	c.HTML(http.StatusOK, "CrearUsuario.html", gin.H{"LoggedIn": true})
}

// Procesar el Formulario Crear (POST)
func (h *UsuarioHandler) Insertar(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	modulos := strings.Join(c.PostFormArray("modulos"), ",")

	// El Rol NO se pide en el form: lo calcula el service a partir de los módulos elegidos.
	req := usuario.UsuarioGestor{
		Username:     &username,
		PasswordHash: &password,
		Email:        &email,
		Modulos:      &modulos,
	}

	err := h.service.Create(&req)
	if err != nil {
		h.logger.Error("Error al crear usuario desde formulario HTML", zap.Error(err))
		c.HTML(http.StatusBadRequest, "CrearUsuario.html", gin.H{
			"Error":    err.Error(),
			"LoggedIn": true,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/usuarios")
}

// Mostrar Formulario de Editar precargado
func (h *UsuarioHandler) Editar(c *gin.Context) {
	idParam := c.Query("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.Redirect(http.StatusSeeOther, "/admin/usuarios")
		return
	}

	usuarioObj, err := h.service.Read(id)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/admin/usuarios")
		return
	}

	seleccionados := map[string]bool{}
	if usuarioObj.Modulos != nil {
		for _, m := range strings.Split(*usuarioObj.Modulos, ",") {
			m = strings.TrimSpace(m)
			if m != "" {
				seleccionados[m] = true
			}
		}
	}

	c.HTML(http.StatusOK, "EditarUsuario.html", gin.H{
		"Usuario":       usuarioObj,
		"Seleccionados": seleccionados,
		"EsAdmin":       esAdmin(usuarioObj),
		"LoggedIn":      true,
	})
}

// Procesar la Actualización (POST)
func (h *UsuarioHandler) Actualizar(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	if id <= 0 {
		c.Redirect(http.StatusSeeOther, "/admin/usuarios")
		return
	}

	username := c.PostForm("username")
	email := c.PostForm("email")

	fields := usuario.UpdateFieldGestor{
		Username: &username,
		Email:    &email,
	}

	// El usuario ADMIN no tiene checkboxes de módulos en el form (ver EditarUsuario.html):
	// dejamos fields.Modulos en nil para que el service no toque ni módulos ni rol.
	if c.PostForm("es_admin") != "true" {
		modulos := strings.Join(c.PostFormArray("modulos"), ",")
		fields.Modulos = &modulos
	}

	// La contraseña es opcional al editar: solo se manda (y se hashea en el service) si se tipeó algo.
	if password := c.PostForm("password"); password != "" {
		fields.PasswordHash = &password
	}

	err := h.service.Update(id, &fields)
	if err != nil {
		h.logger.Error("Error al actualizar usuario desde formulario HTML", zap.Int("id", id), zap.Error(err))
		usuarioObj, _ := h.service.Read(id)
		c.HTML(http.StatusBadRequest, "EditarUsuario.html", gin.H{
			"Error":    err.Error(),
			"Usuario":  usuarioObj,
			"LoggedIn": true,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/usuarios")
}

// Eliminar Registro
func (h *UsuarioHandler) Borrar(c *gin.Context) {
	idParam := c.Query("id")
	id, err := strconv.Atoi(idParam)
	if err == nil && id > 0 {
		if currentID, ok := CurrentUserID(c); ok && currentID == id {
			h.logger.Warn("Intento de autoeliminación bloqueado", zap.Int("id", id))
			c.Redirect(http.StatusSeeOther, "/admin/usuarios")
			return
		}
		if errDel := h.service.Delete(id); errDel != nil {
			h.logger.Error("Error al eliminar usuario desde la plantilla", zap.Int("id", id), zap.Error(errDel))
		}
	}

	c.Redirect(http.StatusSeeOther, "/admin/usuarios")
}

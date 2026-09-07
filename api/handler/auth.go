package handler

import (
	"PaginaSEG/internal/integrante"
	"PaginaSEG/internal/usuario" // Importamos tu nuevo paquete de usuarios
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	usuarioService    *usuario.Service
	integranteService *integrante.Service
	logger            *zap.Logger
}

func NewAuthHandler(us *usuario.Service, is *integrante.Service, l *zap.Logger) *AuthHandler {
	return &AuthHandler{
		usuarioService:    us,
		integranteService: is,
		logger:            l,
	}
}

func (h *AuthHandler) ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "Login.html", gin.H{})
}

func (h *AuthHandler) ProcessLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user, err := h.usuarioService.ReadByUsername(username)

	if err != nil {
		h.logger.Warn("Usuario no existe", zap.Error(err))
		c.HTML(http.StatusOK, "Login.html", gin.H{
			"Error": "Usuario o contrasna incorrecto",
		})
		return

	}

	if user.PasswordHash != nil {
		err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password))
		if err == nil {
			h.logger.Warn("Usuario correcto", zap.String("username", username))
			c.SetCookie("session", strconv.Itoa(user.ID), 3600, "/", "", false, true)
			c.Redirect(http.StatusFound, "/admin/dashboard")
			return
		}
		h.logger.Warn("Intento de login fallido: contrasena incorrecta", zap.String("username", username))
		c.HTML(http.StatusUnauthorized, "Login.html", gin.H{
			"Error": "Intento de login fallido",
		})
		return
	}

	h.logger.Warn("Intento de login fallido: contrasena incorrecta", zap.String("username", username))
	c.HTML(http.StatusUnauthorized, "Login.html", gin.H{
		"Error": "Usuario o contraseña incorrectos",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("session", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

func (h *AuthHandler) ShowDashboard(c *gin.Context) {
	u, ok := currentUsuario(c, h.usuarioService)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	numUsuarios := 0
	numAdmins := 0
	if usuarios, err := h.usuarioService.GetAll(); err == nil {
		numUsuarios = len(usuarios)
		for _, us := range usuarios {
			if us.Rol != nil && *us.Rol == "ADMIN" {
				numAdmins++
			}
		}
	}

	numIntegrantes := 0
	if integrantes, err := h.integranteService.GetAll(); err == nil {
		numIntegrantes = len(integrantes)
	}

	c.HTML(http.StatusOK, "Dashboard.html", gin.H{
		"LoggedIn":             true,
		"EsAdmin":              esAdmin(u),
		"TieneIntegrantes":     tieneModulo(u, "integrantes"),
		"TieneProyectos":       tieneModulo(u, "proyectos"),
		"TieneReconocimientos": tieneModulo(u, "reconocimientos"),
		"TieneEmpresas":        tieneModulo(u, "empresas"),
		"NumUsuarios":          numUsuarios,
		"NumAdmins":            numAdmins,
		"NumIntegrantes":       numIntegrantes,
		// Módulos realmente activos hoy: Integrantes y Administradores.
		// Proyectos, Reconocimientos y Empresas siguen "Próximamente".
		"NumModulosActivos": 2,
	})
}

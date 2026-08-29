package handler

import (
	"PaginaSEG/internal/usuario" // Importamos tu nuevo paquete de usuarios
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	usuarioService *usuario.Service
	logger         *zap.Logger
}

func NewAuthHandler(us *usuario.Service, l *zap.Logger) *AuthHandler {
	return &AuthHandler{
		usuarioService: us,
		logger:         l,
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

	if user.PasswordHash != nil && *user.PasswordHash == password {
		h.logger.Info("Login Existoso", zap.String("username", username))
		c.SetCookie("session", "auth", 3600, "/", "", false, true)
		c.Redirect(http.StatusFound, "/admin/dashboard")
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
	cookie, err := c.Cookie("session")
	if err != nil || cookie != "auth" {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "Dashboard.html", gin.H{})
}

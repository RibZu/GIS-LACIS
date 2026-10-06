package handler

import (
	"PaginaSEG/internal/integrante"
	"PaginaSEG/internal/usuario"
	"crypto/rand"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	mensajeLoginFallido      = "Usuario o contraseña incorrectos."
	mensajeLoginNoDisponible = "No pudimos iniciar sesión en este momento. Probá de nuevo en unos minutos."
)

var hashRelleno []byte

func init() {
	aleatorio := make([]byte, 16)
	if _, err := rand.Read(aleatorio); err != nil {
		panic("no se pudo generar el hash de relleno: " + err.Error())
	}
	hash, err := bcrypt.GenerateFromPassword(aleatorio, bcrypt.DefaultCost)
	if err != nil {
		panic("no se pudo generar el hash de relleno: " + err.Error())
	}
	hashRelleno = hash
}

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
	if _, ok := CurrentUserID(c); ok {
		c.Redirect(http.StatusFound, "/admin/dashboard")
		return
	}
	datos := gin.H{}
	if c.Query("sesion") == "expirada" {
		datos["Aviso"] = "Tu sesión expiró. Volvé a ingresar."
	}
	c.HTML(http.StatusOK, "Login.html", datos)
}

func (h *AuthHandler) loginFallido(c *gin.Context, username string) {
	h.logger.Warn("Intento de login fallido", zap.String("username", username))
	c.HTML(http.StatusUnauthorized, "Login.html", gin.H{"Error": mensajeLoginFallido})
}

func (h *AuthHandler) ProcessLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		_ = bcrypt.CompareHashAndPassword(hashRelleno, []byte(password))
		h.loginFallido(c, username)
		return
	}

	user, err := h.usuarioService.ReadByUsername(username)
	if err != nil && !errors.Is(err, usuario.ErrNotFound) {
		c.HTML(http.StatusServiceUnavailable, "Login.html", gin.H{"Error": mensajeLoginNoDisponible})
		return
	}
	if err != nil || user == nil || user.PasswordHash == nil {
		_ = bcrypt.CompareHashAndPassword(hashRelleno, []byte(password))
		h.loginFallido(c, username)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		h.loginFallido(c, username)
		return
	}

	h.logger.Info("Inicio de sesión", zap.String("username", username), zap.Int("id", user.ID))
	IniciarSesion(c, user.ID)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	CerrarSesion(c)
	c.Redirect(http.StatusFound, "/login")
}

func (h *AuthHandler) ShowDashboard(c *gin.Context) {

	u, ok := autenticar(c, h.usuarioService, false)
	if !ok {
		return
	}

	c.HTML(http.StatusOK, "Dashboard.html", gin.H{
		"LoggedIn":             true,
		"EsAdmin":              esAdmin(u),
		"TieneIntegrantes":     tieneModulo(u, "integrantes"),
		"TieneTesis":           tieneModulo(u, "tesis"),
		"TieneProyectos":       tieneModulo(u, "proyectos"),
		"TieneReconocimientos": tieneModulo(u, "reconocimientos"),
		"TieneEmpresas":        tieneModulo(u, "empresas"),
		"TieneEstadisticas":    tieneModulo(u, "estadisticas"),
	})
}

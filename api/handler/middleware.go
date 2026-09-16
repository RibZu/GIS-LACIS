package handler

import (
	"net/http"
	"strings"

	"PaginaSEG/internal/usuario"

	"github.com/gin-gonic/gin"
)

func currentUsuario(c *gin.Context, s *usuario.Service) (*usuario.UsuarioGestor, bool) {
	id, ok := CurrentUserID(c)
	if !ok {
		return nil, false
	}
	u, err := s.Read(id)
	if err != nil {
		return nil, false
	}
	return u, true
}

func esAdmin(u *usuario.UsuarioGestor) bool {
	return u.Rol != nil && *u.Rol == "ADMIN"
}

func tieneModulo(u *usuario.UsuarioGestor, modulo string) bool {
	if esAdmin(u) {
		return true
	}
	if u.Modulos == nil {
		return false
	}
	for _, m := range strings.Split(*u.Modulos, ",") {
		if strings.TrimSpace(m) == modulo {
			return true
		}
	}
	return false
}

// RequireLogin bloquea la ruta salvo que haya una sesión válida (un usuario logueado que
// realmente exista). No chequea módulo ni rol: es el gate de entrada del grupo /admin completo,
// para que ninguna ruta —presente o futura— quede alcanzable sin sesión aunque a alguien se le
// olvide agregarle RequireModule/RequireAdmin. El chequeo de módulo/rol lo siguen resolviendo
// esos dos middlewares.
func RequireLogin(s *usuario.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := currentUsuario(c, s)
		if !ok {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireModule bloquea la ruta salvo que el usuario logueado tenga ese módulo
// asignado (o sea ADMIN, que tiene acceso a todo).
func RequireModule(s *usuario.Service, modulo string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := currentUsuario(c, s)
		if !ok {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		if !tieneModulo(u, modulo) {
			c.Redirect(http.StatusFound, "/admin/dashboard")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireAdmin bloquea la ruta salvo que el usuario logueado tenga rol ADMIN.
// Se usa para /admin/usuarios*, ya que "administradores" nunca es un módulo asignable.
func RequireAdmin(s *usuario.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := currentUsuario(c, s)
		if !ok {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		if !esAdmin(u) {
			c.Redirect(http.StatusFound, "/admin/dashboard")
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireModuleAPI(s *usuario.Service, modulo string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := currentUsuario(c, s)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No autenticado"})
			c.Abort()
			return
		}
		if !tieneModulo(u, modulo) {
			c.JSON(http.StatusForbidden, gin.H{"error": "No tenés permiso para este módulo"})
			c.Abort()
			return
		}
		c.Next()
	}
}

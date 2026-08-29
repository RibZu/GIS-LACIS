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

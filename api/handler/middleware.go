package handler

import (
	"errors"
	"net/http"
	"strings"

	"PaginaSEG/internal/usuario"

	"github.com/gin-gonic/gin"
)

const claveUsuarioContexto = "usuario"

const mensajeServicioNoDisponible = "No pudimos verificar tu sesión en este momento. Probá de nuevo en unos minutos."

type resultadoUsuario int

const (
	usuarioEncontrado resultadoUsuario = iota
	usuarioSinSesion
	usuarioErrorDeLectura
)

func buscarUsuario(c *gin.Context, s *usuario.Service) (*usuario.UsuarioGestor, resultadoUsuario) {
	if v, existe := c.Get(claveUsuarioContexto); existe {
		if u, esUsuario := v.(*usuario.UsuarioGestor); esUsuario {
			return u, usuarioEncontrado
		}
	}
	id, ok := CurrentUserID(c)
	if !ok {
		return nil, usuarioSinSesion
	}
	u, err := s.Read(id)
	if err != nil {
		if errors.Is(err, usuario.ErrNotFound) {
			return nil, usuarioSinSesion
		}
		return nil, usuarioErrorDeLectura
	}
	c.Set(claveUsuarioContexto, u)
	return u, usuarioEncontrado
}

func currentUsuario(c *gin.Context, s *usuario.Service) (*usuario.UsuarioGestor, bool) {
	u, resultado := buscarUsuario(c, s)
	return u, resultado == usuarioEncontrado
}

func redirigirAlLogin(c *gin.Context) {
	if _, estado := leerSesion(c); estado == sesionVencida {
		c.Redirect(http.StatusFound, "/login?sesion=expirada")
	} else {
		c.Redirect(http.StatusFound, "/login")
	}
	c.Abort()
}

func autenticar(c *gin.Context, s *usuario.Service, esAPI bool) (*usuario.UsuarioGestor, bool) {
	u, resultado := buscarUsuario(c, s)
	switch resultado {
	case usuarioEncontrado:
		return u, true
	case usuarioErrorDeLectura:
		if esAPI {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Servicio no disponible"})
		} else {
			paginaError(c, http.StatusServiceUnavailable, "Servicio no disponible", mensajeServicioNoDisponible)
		}
		c.Abort()
	default:
		if esAPI {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No autenticado"})
			c.Abort()
		} else {
			redirigirAlLogin(c)
		}
	}
	return nil, false
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

func RequireLogin(s *usuario.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		if _, ok := autenticar(c, s, false); !ok {
			return
		}
		c.Next()
	}
}

func RequireModule(s *usuario.Service, modulo string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := autenticar(c, s, false)
		if !ok {
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

func RequireAdmin(s *usuario.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := autenticar(c, s, false)
		if !ok {
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
		u, ok := autenticar(c, s, true)
		if !ok {
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

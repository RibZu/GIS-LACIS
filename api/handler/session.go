package handler

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const sessionCookieName = "session"

const duracionSesion = 3600

type estadoSesion int

const (
	sesionNinguna estadoSesion = iota
	sesionValida
	sesionVencida
	sesionInvalida
)

var claveSesion = claveAleatoria()

func claveAleatoria() []byte {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("no se pudo generar la clave de sesión: " + err.Error())
	}
	return b
}

func ConfigurarSesion(secreto string, logger *zap.Logger) {
	if secreto != "" {
		claveSesion = []byte(secreto)
		return
	}
	logger.Warn("SESSION_SECRET no definida: las sesiones se pierden al reiniciar")
}

func firmar(id int, vence int64) string {
	m := hmac.New(sha256.New, claveSesion)
	m.Write([]byte(fmt.Sprintf("%d.%d", id, vence)))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}

func leerSesion(c *gin.Context) (int, estadoSesion) {
	cookie, err := c.Cookie(sessionCookieName)
	if err != nil || cookie == "" {
		return 0, sesionNinguna
	}
	partes := strings.Split(cookie, ".")
	if len(partes) != 3 {
		return 0, sesionInvalida
	}
	id, err := strconv.Atoi(partes[0])
	if err != nil || id <= 0 {
		return 0, sesionInvalida
	}
	vence, err := strconv.ParseInt(partes[1], 10, 64)
	if err != nil {
		return 0, sesionInvalida
	}
	if !hmac.Equal([]byte(partes[2]), []byte(firmar(id, vence))) {
		return 0, sesionInvalida
	}
	if time.Now().Unix() > vence {
		return 0, sesionVencida
	}
	return id, sesionValida
}

func CurrentUserID(c *gin.Context) (int, bool) {
	id, estado := leerSesion(c)
	return id, estado == sesionValida
}

func esHTTPS(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}
	primero := strings.TrimSpace(strings.Split(c.GetHeader("X-Forwarded-Proto"), ",")[0])
	return strings.EqualFold(primero, "https")
}

func IniciarSesion(c *gin.Context, id int) {
	vence := time.Now().Unix() + duracionSesion
	valor := fmt.Sprintf("%d.%d.%s", id, vence, firmar(id, vence))
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(sessionCookieName, valor, duracionSesion, "/", "", esHTTPS(c), true)
}

func CerrarSesion(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(sessionCookieName, "", -1, "/", "", esHTTPS(c), true)
}

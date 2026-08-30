package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const sessionCookieName = "session"

// CurrentUserID lee el ID del usuario logueado desde la cookie de sesión.
// La cookie guarda el ID (no un valor fijo) desde que auth.go pasó a bcrypt.
func CurrentUserID(c *gin.Context) (int, bool) {
	cookie, err := c.Cookie(sessionCookieName)
	if err != nil {
		return 0, false
	}
	id, err := strconv.Atoi(cookie)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

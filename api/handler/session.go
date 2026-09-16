package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const sessionCookieName = "session"

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

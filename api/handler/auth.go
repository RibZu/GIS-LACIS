package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "Login.html", gin.H{})
}

func (h *AuthHandler) ProcessLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "admin" && password == "admin123" {
		c.SetCookie("session", "auth", 3600, "/", "", false, true)
		c.Redirect(http.StatusFound, "/admin/dashboard")
		return
	}

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

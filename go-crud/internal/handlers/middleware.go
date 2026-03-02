package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {
	s := sessions.Default(c)
	if s.Get("user_id") == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}
	c.Next()
}

func RequireAdmin(c *gin.Context) {
	// Validamos admin contra el usuario real (BD) para que coincida con navbar y permisos.
	u := currentUserFromSession(c)
	if u == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}
	if !u.IsAdmin {
		c.String(http.StatusForbidden, "solo admin")
		c.Abort()
		return
	}
	c.Next()
}
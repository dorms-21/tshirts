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
	s := sessions.Default(c)
	isAdmin, ok := s.Get("is_admin").(bool)
	if !ok || !isAdmin {
		c.String(http.StatusForbidden, "solo admin")
		c.Abort()
		return
	}
	c.Next()
}

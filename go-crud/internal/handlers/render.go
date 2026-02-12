package handlers

import (
	"go-crud/internal/db"
	"go-crud/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func currentUserFromSession(c *gin.Context) *models.User {
	sess := sessions.Default(c)
	uidRaw := sess.Get("user_id")
	if uidRaw == nil {
		return nil
	}

	var uid int64
	switch v := uidRaw.(type) {
	case int64:
		uid = v
	case int:
		uid = int64(v)
	case float64:
		uid = int64(v)
	default:
		return nil
	}

	var u models.User
	if err := db.DB.First(&u, uid).Error; err != nil {
		return nil
	}
	return &u
}

func Render(c *gin.Context, status int, tpl string, data gin.H) {
	if data == nil {
		data = gin.H{}
	}
	data["CurrentUser"] = currentUserFromSession(c)
	c.HTML(status, tpl, data)
}

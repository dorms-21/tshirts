package handlers

import (
	"fmt"
	"strconv"

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

	var uid uint64

	switch v := uidRaw.(type) {
	case uint:
		uid = uint64(v)
	case uint64:
		uid = v
	case int:
		if v < 0 {
			return nil
		}
		uid = uint64(v)
	case int64:
		if v < 0 {
			return nil
		}
		uid = uint64(v)
	case float64:
		if v < 0 {
			return nil
		}
		uid = uint64(v)
	case string:
		parsed, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil
		}
		uid = parsed
	default:
		// útil para debug: ver qué tipo guardó la sesión
		_ = fmt.Sprintf("unsupported user_id type: %T", uidRaw)
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
package handlers

import (
	"net/http"
	"strings"

	"go-crud/internal/db"
	"go-crud/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func getSessionUserID(c *gin.Context) (uint, bool) {
	s := sessions.Default(c)
	v := s.Get("user_id")
	if v == nil {
		return 0, false
	}
	switch t := v.(type) {
	case uint:
		return t, true
	case int:
		if t < 0 {
			return 0, false
		}
		return uint(t), true
	case int64:
		if t < 0 {
			return 0, false
		}
		return uint(t), true
	case float64:
		if t < 0 {
			return 0, false
		}
		return uint(t), true
	default:
		return 0, false
	}
}

func ProfileSurveyGet(c *gin.Context) {
	uid, _ := getSessionUserID(c)

	// Si ya existe perfil, lo mandamos para pre-rellenar (opcional)
	var prof models.UserProfile
	_ = db.DB.Where("user_id = ?", uid).First(&prof).Error

	c.HTML(http.StatusOK, "profile/survey.html", gin.H{
		"Title": "Perfil",
		"Breadcrumbs": []Crumb{
			{Label: "Tienda", Href: "/shop", Active: false},
			{Label: "Perfil", Href: "", Active: true},
		},
		"Profile": prof,
	})
}

func ProfileSurveyPost(c *gin.Context) {
	uid, ok := getSessionUserID(c)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	gender := strings.TrimSpace(strings.ToLower(c.PostForm("gender")))
	location := strings.TrimSpace(c.PostForm("location"))
	style := strings.TrimSpace(strings.ToLower(c.PostForm("style")))
	colors := strings.TrimSpace(strings.ToLower(c.PostForm("colors")))
	fit := strings.TrimSpace(strings.ToLower(c.PostForm("fit")))
	occasion := strings.TrimSpace(strings.ToLower(c.PostForm("occasion")))

	// Validaciones mínimas
	if gender != "hombre" && gender != "mujer" {
		c.HTML(http.StatusBadRequest, "profile/survey.html", gin.H{
			"Title": "Perfil",
			"Error": "Selecciona un género.",
		})
		return
	}
	if location == "" {
		c.HTML(http.StatusBadRequest, "profile/survey.html", gin.H{
			"Title": "Perfil",
			"Error": "Escribe tu localidad.",
		})
		return
	}
	if style == "" || colors == "" || fit == "" || occasion == "" {
		c.HTML(http.StatusBadRequest, "profile/survey.html", gin.H{
			"Title": "Perfil",
			"Error": "Selecciona una opción en cada fila de preferencias.",
		})
		return
	}

	// Upsert (si existe, actualiza; si no, crea)
	var prof models.UserProfile
	err := db.DB.Where("user_id = ?", uid).First(&prof).Error

	if err == nil {
		// update
		prof.Gender = gender
		prof.Location = location
		prof.Style = style
		prof.Colors = colors
		prof.Fit = fit
		prof.Occasion = occasion

		if err := db.DB.Save(&prof).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "errors/500.html", gin.H{
				"Title": "500",
			})
			return
		}
	} else {
		// create
		prof = models.UserProfile{
			UserID:   uid,
			Gender:   gender,
			Location: location,
			Style:    style,
			Colors:   colors,
			Fit:      fit,
			Occasion: occasion,
		}

		if err := db.DB.Create(&prof).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "errors/500.html", gin.H{
				"Title": "500",
			})
			return
		}
	}

	// Después de guardar, regresa a tienda (luego aquí haremos recomendaciones)
	c.Redirect(http.StatusSeeOther, "/shop")
}

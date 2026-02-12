package handlers

import (
	"net/http"
	"strings"

	"go-crud/internal/db"
	"go-crud/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ShowRegister(c *gin.Context) {
	q, t := NewCaptchaForForm(c)
	Render(c, http.StatusOK, "auth/register.html", gin.H{
		"Title": "Registro",
		"Breadcrumbs": []Crumb{
			{Label: "Tienda", Href: "/shop", Active: false},
			{Label: "Registro", Href: "", Active: true},
		},
		"CaptchaQ":     q,
		"CaptchaToken": t,
	})
}

func Register(c *gin.Context) {
	email := strings.TrimSpace(strings.ToLower(c.PostForm("email")))
	name := strings.TrimSpace(c.PostForm("name"))
	pass := c.PostForm("password")

	if email == "" || name == "" || len(pass) < 6 {
		q, t := NewCaptchaForForm(c)
		Render(c, http.StatusBadRequest, "auth/register.html", gin.H{
			"Title":        "Registro",
			"Error":        "Datos inválidos (password mínimo 6).",
			"CaptchaQ":     q,
			"CaptchaToken": t,
		})
		return
	}

	if !CheckCaptcha(c) {
		RenderCaptchaError(c, "auth/register.html", "Registro",
			[]Crumb{{Label: "Tienda", Href: "/shop", Active: false}, {Label: "Registro", Href: "", Active: true}},
			gin.H{},
		)
		return
	}

	var existing models.User
	if err := db.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		q, t := NewCaptchaForForm(c)
		Render(c, http.StatusBadRequest, "auth/register.html", gin.H{
			"Title":        "Registro",
			"Error":        "Ese correo ya existe.",
			"CaptchaQ":     q,
			"CaptchaToken": t,
		})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	u := models.User{
		Email:        email,
		Name:         name,
		PasswordHash: string(hash),
		IsAdmin:      false,
	}
	if err := db.DB.Create(&u).Error; err != nil {
		q, t := NewCaptchaForForm(c)
		Render(c, http.StatusInternalServerError, "auth/register.html", gin.H{
			"Title":        "Registro",
			"Error":        "No se pudo crear el usuario.",
			"CaptchaQ":     q,
			"CaptchaToken": t,
		})
		return
	}

	// login automático
	s := sessions.Default(c)
	s.Set("user_id", u.ID)
	s.Set("is_admin", u.IsAdmin)
	_ = s.Save()

	c.Redirect(http.StatusSeeOther, "/shop")
}

func ShowLogin(c *gin.Context) {
	Render(c, http.StatusOK, "auth/login.html", gin.H{
		"Title": "Login",
		"Breadcrumbs": []Crumb{
			{Label: "Tienda", Href: "/shop", Active: false},
			{Label: "Login", Href: "", Active: true},
		},
	})
}

func Login(c *gin.Context) {
	email := strings.TrimSpace(strings.ToLower(c.PostForm("email")))
	pass := c.PostForm("password")

	var u models.User
	if err := db.DB.Where("email = ?", email).First(&u).Error; err != nil {
		Render(c, http.StatusUnauthorized, "auth/login.html", gin.H{
			"Title": "Login",
			"Error": "Credenciales inválidas.",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pass)); err != nil {
		Render(c, http.StatusUnauthorized, "auth/login.html", gin.H{
			"Title": "Login",
			"Error": "Credenciales inválidas.",
		})
		return
	}

	s := sessions.Default(c)
	s.Set("user_id", u.ID)
	s.Set("is_admin", u.IsAdmin)
	_ = s.Save()

	c.Redirect(http.StatusSeeOther, "/shop")
}

func Logout(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	_ = s.Save()
	c.Redirect(http.StatusSeeOther, "/shop")
}

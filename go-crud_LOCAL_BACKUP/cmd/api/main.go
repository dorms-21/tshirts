package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"go-crud/internal/db"
	"go-crud/internal/handlers"
	"go-crud/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db.Connect()

	// -------------------------
	// Migraciones
	// -------------------------
	if err := db.DB.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.Product{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	r := gin.Default()

	// -------------------------
	// Templates (sin glob que truene)
	// -------------------------
	funcMap := template.FuncMap{
		"money": func(cents int64) string {
			return fmt.Sprintf("%.2f", float64(cents)/100.0)
		},
	}

	tmpl := template.New("").Funcs(funcMap)

	if files, _ := filepath.Glob("internal/web/templates/*.html"); len(files) > 0 {
		template.Must(tmpl.ParseFiles(files...))
	}
	if files, _ := filepath.Glob("internal/web/templates/*/*.html"); len(files) > 0 {
		template.Must(tmpl.ParseFiles(files...))
	}

	r.SetHTMLTemplate(tmpl)
	r.Static("/static", "internal/web/static")

	// -------------------------
	// Sesiones
	// -------------------------
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "dev-session-secret-change-me"
	}

	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7,
		HttpOnly: true,
		SameSite: 2, // Lax
	})
	r.Use(sessions.Sessions("go_app_session", store))

	// -------------------------
	// Captcha secret
	// -------------------------
	captchaSecret := os.Getenv("CAPTCHA_SECRET")
	if captchaSecret == "" {
		captchaSecret = "dev-captcha-secret-change-me-32-bytes-min"
	}
	r.Use(func(c *gin.Context) {
		c.Set("CAPTCHA_SECRET", []byte(captchaSecret))
		c.Next()
	})

	// -------------------------
	// Health
	// -------------------------
	r.GET("/health", func(c *gin.Context) {
		sqlDB, err := db.DB.DB()
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "error": err.Error()})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"ok": true})
	})

	// -------------------------
	// Auth
	// -------------------------
	r.GET("/register", handlers.ShowRegister)
	r.POST("/register", handlers.Register)
	r.GET("/login", handlers.ShowLogin)
	r.POST("/login", handlers.Login)
	r.POST("/logout", handlers.Logout)

	// -------------------------
	// Shop
	// -------------------------
	r.GET("/", handlers.ShopIndex)
	r.GET("/shop", handlers.ShopIndex)
	r.GET("/shop/:id", handlers.ShopDetail)

	// -------------------------
	// Profile Survey (protegido)
	// -------------------------
	profile := r.Group("/profile")
	profile.Use(handlers.RequireAuth)
	{
		profile.GET("/survey", handlers.ProfileSurveyGet)
		profile.POST("/survey", handlers.ProfileSurveyPost)
	}

	// -------------------------
	// Admin (protegido)
	// -------------------------
	admin := r.Group("/admin")
	admin.Use(handlers.RequireAuth, handlers.RequireAdmin)
	{
		admin.GET("/products", handlers.AdminProductsIndex)
		admin.GET("/products/new", handlers.AdminProductsNewForm)
		admin.POST("/products", handlers.AdminProductsCreate)
		admin.GET("/products/:id/edit", handlers.AdminProductsEditForm)
		admin.POST("/products/:id/edit", handlers.AdminProductsUpdate)
	}

	// -------------------------
	// Puerto
	// -------------------------
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	_ = r.Run(":" + port)
}

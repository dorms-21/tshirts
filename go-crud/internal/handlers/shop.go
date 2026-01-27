package handlers

import (
	"net/http"

	"go-crud/internal/db"
	"go-crud/internal/models"

	"github.com/gin-gonic/gin"
)

func ShopIndex(c *gin.Context) {
	var products []models.Product
	db.DB.Where("active = true").Order("id desc").Find(&products)

	c.HTML(http.StatusOK, "shop/index.html", gin.H{
		"Title": "Tienda",
		"Breadcrumbs": []Crumb{
			{Label: "Tienda", Href: "/shop", Active: true},
		},
		"Products": products,
	})
}

func ShopDetail(c *gin.Context) {
	var p models.Product
	if err := db.DB.First(&p, c.Param("id")).Error; err != nil {
		c.HTML(http.StatusNotFound, "errors/404.html", gin.H{
			"Title": "404",
			"Breadcrumbs": []Crumb{
				{Label: "Tienda", Href: "/shop", Active: false},
				{Label: "404", Href: "", Active: true},
			},
		})
		return
	}

	c.HTML(http.StatusOK, "shop/detail.html", gin.H{
		"Title": p.Name,
		"Breadcrumbs": []Crumb{
			{Label: "Tienda", Href: "/shop", Active: false},
			{Label: p.Name, Href: "", Active: true},
		},
		"Product": p,
	})
}

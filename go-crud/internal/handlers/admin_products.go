package handlers

import (
	"net/http"
	"strconv"

	"go-crud/internal/db"
	"go-crud/internal/models"

	"github.com/gin-gonic/gin"
)

type Crumb struct {
	Label  string
	Href   string
	Active bool
}

func AdminProductsIndex(c *gin.Context) {
	var products []models.Product
	db.DB.Order("id desc").Find(&products)

	c.HTML(http.StatusOK, "admin/products_index.html", gin.H{
		"Title": "Admin - Productos",
		"Breadcrumbs": []Crumb{
			{Label: "Admin", Href: "/admin/products", Active: false},
			{Label: "Productos", Href: "", Active: true},
		},
		"Products": products,
	})
}

func AdminProductsNewForm(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/products_new.html", gin.H{
		"Title": "Admin - Nuevo Producto",
		"Breadcrumbs": []Crumb{
			{Label: "Admin", Href: "/admin/products", Active: false},
			{Label: "Nuevo producto", Href: "", Active: true},
		},
	})
}

func AdminProductsCreate(c *gin.Context) {
	name := c.PostForm("name")
	desc := c.PostForm("description")

	priceCents, err := strconv.ParseInt(c.PostForm("price_cents"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "price_cents inválido")
		return
	}

	stock, err := strconv.Atoi(c.PostForm("stock"))
	if err != nil {
		c.String(http.StatusBadRequest, "stock inválido")
		return
	}

	p := models.Product{
		Name:        name,
		Description: desc,
		PriceCents:  priceCents,
		Stock:       stock,
		Active:      true,
	}

	if err := db.DB.Create(&p).Error; err != nil {
		c.String(http.StatusInternalServerError, "error creando producto")
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/products")
}

func AdminProductsEditForm(c *gin.Context) {
	var p models.Product
	if err := db.DB.First(&p, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "producto no encontrado")
		return
	}

	c.HTML(http.StatusOK, "admin/products_edit.html", gin.H{
		"Title": "Admin - Editar Producto",
		"Breadcrumbs": []Crumb{
			{Label: "Admin", Href: "/admin/products", Active: false},
			{Label: "Editar", Href: "", Active: true},
		},
		"Product": p,
	})
}

func AdminProductsUpdate(c *gin.Context) {
	var p models.Product
	if err := db.DB.First(&p, c.Param("id")).Error; err != nil {
		c.String(http.StatusNotFound, "producto no encontrado")
		return
	}

	p.Name = c.PostForm("name")
	p.Description = c.PostForm("description")

	priceCents, err := strconv.ParseInt(c.PostForm("price_cents"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "price_cents inválido")
		return
	}

	stock, err := strconv.Atoi(c.PostForm("stock"))
	if err != nil {
		c.String(http.StatusBadRequest, "stock inválido")
		return
	}

	p.PriceCents = priceCents
	p.Stock = stock

	if err := db.DB.Save(&p).Error; err != nil {
		c.String(http.StatusInternalServerError, "error guardando producto")
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/products")
}

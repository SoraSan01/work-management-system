package server

import (
	"html/template"
	"workms/internal/routes"

	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {
	r := gin.Default()

	r.Static("/static", "./static")

	tmpl := template.Must(template.ParseGlob("templates/**/*.html"))
	r.SetHTMLTemplate(tmpl)

	routes.SetupRoutes(r)
	return r
}

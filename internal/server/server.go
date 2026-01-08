package server

import (
	"html/template"
	"workms/internal/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {
	r := gin.Default()

	// Static files
	r.Static("/static", "./static")

	// Sessions
	store := cookie.NewStore([]byte("Sora@09089831215")) // choose a secure secret
	r.Use(sessions.Sessions("workms_session", store))

	// Load templates
	tmpl := template.Must(template.ParseGlob("templates/**/*.html"))
	r.SetHTMLTemplate(tmpl)

	// Routes
	routes.SetupRoutes(r)

	return r
}

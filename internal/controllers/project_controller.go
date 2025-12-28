package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListProjects(c *gin.Context) {
	c.HTML(http.StatusOK, "projects/list.html", gin.H{
		"title": "Projects List",
	})
}

func CreateProjectForm(c *gin.Context) {
	c.HTML(http.StatusOK, "projects/create.html", gin.H{
		"title": "Projects Create",
	})
}

func EditProjectForm(c *gin.Context) {
	c.HTML(http.StatusOK, "projects/edit.html", gin.H{
		"title": "Projects Edit",
	})
}

func DetailProject(c *gin.Context) {
	c.HTML(http.StatusOK, "projects/detail.html", gin.H{
		"title": "Projects Detail",
	})
}

func BoardProject(c *gin.Context) {
	c.HTML(http.StatusOK, "projects/board.html", gin.H{
		"title": "Projects Board",
	})
}

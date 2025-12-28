package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReportTask(c *gin.Context) {
	c.HTML(http.StatusOK, "reports/tasks.html", gin.H{
		"title": "Task Reports",
	})
}

func ReportEmployee(c *gin.Context) {
	c.HTML(http.StatusOK, "reports/employees.html", gin.H{
		"title": "Employee Reports",
	})
}

func ReportProject(c *gin.Context) {
	c.HTML(http.StatusOK, "reports/projects.html", gin.H{
		"title": "Project Reports",
	})
}

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard/index.html", gin.H{
		"title": "dashboard",
	})
}

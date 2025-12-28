package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListDocument(c *gin.Context) {
	c.HTML(http.StatusOK, "document/list.html", gin.H{
		"title": "Document List",
	})
}

func UploadDocument(c *gin.Context) {
	c.HTML(http.StatusOK, "document/upload.html", gin.H{
		"title": "Document Upload",
	})
}

func DetailDocument(c *gin.Context) {
	c.HTML(http.StatusOK, "document/detail.html", gin.H{
		"title": "Document Details",
	})
}

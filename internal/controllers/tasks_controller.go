package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListTask(c *gin.Context) {
	c.HTML(http.StatusOK, "task/list.html", gin.H{
		"title": "Task List",
	})
}

func CreateTaskForm(c *gin.Context) {
	c.HTML(http.StatusOK, "task/create.html", gin.H{
		"title": "Task Create",
	})
}

func EditTaskForm(c *gin.Context) {
	c.HTML(http.StatusOK, "task/edit.html", gin.H{
		"title": "Task Edit",
	})
}

func DetailTask(c *gin.Context) {
	c.HTML(http.StatusOK, "task/detail.html", gin.H{
		"title": "Task Detail",
	})
}

func BoardTask(c *gin.Context) {
	c.HTML(http.StatusOK, "task/board.html", gin.H{
		"title": "Task Board",
	})
}

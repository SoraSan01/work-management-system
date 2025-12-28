package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CalendarIndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "calendar/index.html", gin.H{
		"title": "Calendar",
	})
}

func CreateEventForm(c *gin.Context) {
	c.HTML(http.StatusOK, "calendar/event-create.html", gin.H{
		"title": "Create Event",
	})
}

func DetailEvent(c *gin.Context) {
	c.HTML(http.StatusOK, "calendar/event-detail.html", gin.H{
		"title": "Detail Event",
	})
}

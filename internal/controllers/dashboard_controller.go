package controllers

import (
	"net/http"
	"workms/internal/repositories"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	UserRepo *repositories.UserRepository
}

func NewDashboardController(userRepo *repositories.UserRepository) *DashboardController {
	return &DashboardController{UserRepo: userRepo}
}

func (dc *DashboardController) Index(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Use the repository instance (dc.UserRepo) instead of NewUserRepository
	user, err := dc.UserRepo.GetByID(userID.(uint64))
	if err != nil || user == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Example notifications
	notifications := []struct {
		Message string
		TimeAgo string
	}{
		{"New user registered", "2 minutes ago"},
		{"Task completed successfully", "1 hour ago"},
		{"System update required", "3 hours ago"},
	}

	c.HTML(http.StatusOK, "dashboard/index.html", gin.H{
		"pageTitle":     "Dashboard",
		"User":          user,
		"Notifications": notifications,
	})
}

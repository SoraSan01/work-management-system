package controllers

import (
	"net/http"
	"workms/internal/repositories"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	Repo *repositories.AuthenticationRepository
}

func NewAuthenticationController(repo *repositories.AuthenticationRepository) *AuthenticationController {
	return &AuthenticationController{Repo: repo}
}

// Login page (GET)
func (ac *AuthenticationController) Login(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID != nil {
		// Already logged in, redirect to dashboard
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}

	// Not logged in, show login page
	c.HTML(http.StatusOK, "authentication/login.html", gin.H{
		"Error": "",
	})
}

// Authenticate user (POST)
func (ac *AuthenticationController) Authenticate(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := ac.Repo.Authenticate(email, password)
	if err != nil || user == nil {
		c.HTML(http.StatusUnauthorized, "/login", gin.H{
			"Error": "Invalid email or password",
		})
		return
	}

	// Set session
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/dashboard")
}

// Logout user
func (ac *AuthenticationController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusSeeOther, "/login")
}

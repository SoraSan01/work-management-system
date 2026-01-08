package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthRequired checks if a user is logged in
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			// Not logged in
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// User is logged in, continue
		c.Next()
	}
}

// Logger logs the request method and path
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("%s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

// AuthAdminOrManager ensures the logged-in user is admin or manager
func AuthAdminOrManager() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		role, ok := session.Get("role").(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if role != "admin" && role != "manager" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: admin or manager only"})
			c.Abort()
			return
		}

		c.Next()
	}
}

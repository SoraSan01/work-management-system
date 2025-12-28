// internal/middleware/auth.go
package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Logger middleware example
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("Request: %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

// AuthRequired is a simple auth middleware example
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple example: check query param ?token=123
		token := c.Query("token")
		if token != "123" {
			c.String(http.StatusUnauthorized, "Unauthorized: missing or wrong token")
			c.Abort()
			return
		}
		c.Next()
	}
}

package middleware

import (
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	// Define allowed origins - replace with your actual frontend domains
	allowedOrigins := []string{
		"https://ccp15.dev.bri.co.id",
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if the origin is in the allowed list
		// Use the whitelisted value directly (not the user-supplied origin)
		var matchedOrigin string
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				matchedOrigin = allowedOrigin
				break
			}
		}

		// Set CORS headers only for allowed origins
		if matchedOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", matchedOrigin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, File-Name, Content-Type, Content-Length")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // Cache preflight for 24 hours

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

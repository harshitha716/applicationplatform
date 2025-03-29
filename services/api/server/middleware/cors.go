package middleware

import (
	"fmt"

	"github.com/Zampfi/application-platform/services/api/helper"
	"github.com/gin-gonic/gin"
)

// CORS middleware
func GetCORSMiddleware(allowedOrigins string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", fmt.Sprintf("Content-Type, Accept, Origin, Referer, %s", helper.PROXY_WORKSPACE_ID_HEADER))
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Location")
		c.Writer.Header().Set("Access-Control-Max-Age", "600")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

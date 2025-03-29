package health

import "github.com/gin-gonic/gin"

func RegisterHealthcheckRoute(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}

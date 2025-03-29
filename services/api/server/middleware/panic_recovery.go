package middleware

import (
	"github.com/gin-gonic/gin"
)

func GetPanicRecoveryMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}

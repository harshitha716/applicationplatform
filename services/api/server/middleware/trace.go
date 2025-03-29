package middleware

import (
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Middleware to inject a new trace ID in gin context
func GetTraceGenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		traceId := uuid.New().String()
		apicontext.AddTraceIdToGinContext(c, traceId)

		c.Header("X-Trace-id", traceId)
		c.Next()
	}
}

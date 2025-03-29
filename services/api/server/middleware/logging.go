package middleware

import (
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		Context: func(c *gin.Context) []zapcore.Field {
			traceId := apicontext.GetTraceIdFromContext(c)
			return []zapcore.Field{zap.String("trace_id", traceId), zap.String("label", "request")}
		},
	})
}

func GetContextLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apicontext.AddLoggerToGinContext(ctx, logger)
		ctx.Next()
	}
}

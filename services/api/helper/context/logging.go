package apicontext

import (
	"context"

	"github.com/Zampfi/application-platform/services/api/pkg/logging"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

const (
	contextKeyLogger string = "logger"
)

// List of sensitive context variables that should not be logged
var sensitiveContextKeys = []string{
	contextKeyUserEmail,
	contextKeyUserIPAddress,
	// Add other sensitive keys as needed
}

func AddLoggerToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

func AddLoggerToGinContext(ctx *gin.Context, logger *zap.Logger) {
	ctx.Set(contextKeyLogger, logger)
}

// GetFilteredCtxVariablesFromCtx returns context variables with sensitive information removed
func GetFilteredCtxVariablesFromCtx(ctx context.Context) map[string]interface{} {
	allVars := GetAllCtxVariablesFromCtx(ctx)
	filteredVars := make(map[string]interface{})
	
	// Copy all variables except sensitive ones to the filtered map
	for key, value := range allVars {
		isSensitive := false
		for _, sensitiveKey := range sensitiveContextKeys {
			if key == sensitiveKey {
				isSensitive = true
				break
			}
		}
		if !isSensitive {
			filteredVars[key] = value
		}
	}
	
	return filteredVars
}

func GetLoggerFromCtx(ctx context.Context) *zap.Logger {
	l := ctx.Value(contextKeyLogger)
	var logger *zap.Logger

	if l != nil {
		lgr, ok := l.(*zap.Logger)
		if !ok {
			logger = lgr
		}
	}

	if logger == nil {
		var err error
		logger, err = logging.GetLogger()
		if err != nil {
			logger = zap.NewNop()
		}
	}
	loggerWithCtx := logger.With(zap.Any("context_variables", GetFilteredCtxVariablesFromCtx(ctx)))
	return loggerWithCtx
}

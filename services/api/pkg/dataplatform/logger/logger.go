package logger

import (
	"context"

	"go.uber.org/zap"
)

const (
	contextKeyLogger       string = "logger"
	contextKeyCtxVariables string = "context_variables"
)

func GetLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func GetAllCtxVariablesFromCtx(ctx context.Context) map[string]interface{} {
	ctxVariables, ok := ctx.Value(contextKeyCtxVariables).(map[string]interface{})
	if !ok {
		return make(map[string]interface{})
	}
	return ctxVariables
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
		logger, err = GetLogger()
		if err != nil {
			logger = zap.NewNop()
		}
	}

	loggerWithCtx := logger.With(zap.Any("context_variables", GetAllCtxVariablesFromCtx(ctx)))
	return loggerWithCtx
}

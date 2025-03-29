package apicontext

import (
	"context"

	"github.com/gin-gonic/gin"
)

const contextKeyTraceId string = "trace_id"

func AddTraceIdToGinContext(ctx *gin.Context, traceId string) {
	AddContextVariableToGinContext(ctx, contextKeyTraceId, traceId)
}

func GetTraceIdFromContext(ctx context.Context) string {
	traceId := getCtxVariableFromCtx(ctx, contextKeyTraceId)
	if traceId == nil {
		return ""
	}
	traceId, ok := traceId.(string)
	if !ok {
		return ""
	}
	return traceId.(string)
}

package apicontext

import (
	"github.com/gin-gonic/gin"
)

func AddContextVariableToGinContext(ctx *gin.Context, key string, value interface{}) {
	ginCtxVarsRaw := ctx.Value(contextKeyCtxVariables)

	var ginCtxVars map[string]interface{}
	if ginCtxVarsRaw == nil {
		ginCtxVars = make(map[string]interface{})
	} else {
		var ok bool
		ginCtxVars, ok = ginCtxVarsRaw.(map[string]interface{})
		if !ok {
			ginCtxVars = make(map[string]interface{})
		}
	}

	ginCtxVars[key] = value

	ctx.Set(contextKeyCtxVariables, ginCtxVars)
}

package apicontext

import (
	"context"
)

const contextKeyCtxVariables string = "context_variables"

func AddCtxVariableToCtx(ctx context.Context, key string, value interface{}) context.Context {
	ctxVariables, ok := ctx.Value(contextKeyCtxVariables).(map[string]interface{})
	if !ok {
		ctxVariables = make(map[string]interface{})
	}

	ctxVariables[key] = value
	return context.WithValue(ctx, contextKeyCtxVariables, ctxVariables)
}

func RemoveCtxVariableFromCtx(ctx context.Context, key string) context.Context {
	ctxVariables, ok := ctx.Value(contextKeyCtxVariables).(map[string]interface{})
	if !ok {
		return ctx
	}
	delete(ctxVariables, key)
	return context.WithValue(ctx, contextKeyCtxVariables, ctxVariables)
}

func GetAllCtxVariablesFromCtx(ctx context.Context) map[string]interface{} {
	ctxVariables, ok := ctx.Value(contextKeyCtxVariables).(map[string]interface{})
	if !ok {
		return make(map[string]interface{})
	}
	return ctxVariables
}

func getCtxVariableFromCtx(ctx context.Context, key string) interface{} {

	ctxVarsRaw := ctx.Value(contextKeyCtxVariables)
	if ctxVarsRaw == nil {
		return nil
	}

	ctxVariables, ok := ctxVarsRaw.(map[string]interface{})
	if !ok {
		return nil
	}
	return ctxVariables[key]
}

package apicontext

import "go.temporal.io/sdk/workflow"

func AddContextVariableToWorkflowContext(ctx workflow.Context, key string, value interface{}) workflow.Context {
	workflowCtxVarsRaw := ctx.Value(contextKeyCtxVariables)

	var workflowCtxVars map[string]interface{}
	if workflowCtxVarsRaw == nil {
		workflowCtxVars = make(map[string]interface{})
	} else {
		var ok bool
		workflowCtxVars, ok = workflowCtxVarsRaw.(map[string]interface{})
		if !ok {
			workflowCtxVars = make(map[string]interface{})
		}
	}

	workflowCtxVars[key] = value

	return workflow.WithValue(ctx, contextKeyCtxVariables, workflowCtxVars)
}

func GetContextVariableFromWorkflowContext(ctx workflow.Context, key string) interface{} {
	return ctx.Value(contextKeyCtxVariables)
}

package workflowutil

import (
	"go.temporal.io/sdk/workflow"
)

func PanicRecoveryHook(wtx workflow.Context) {
	logger := workflow.GetLogger(wtx)

	panicErr := recover()
	if panicErr != nil {
		logger.Error("Recovered from panic", "panic", panicErr)
	}

}

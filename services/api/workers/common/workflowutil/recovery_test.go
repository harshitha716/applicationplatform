package workflowutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
)

// TestPanicRecoveryHook tests the PanicRecoveryHook function
func TestPanicRecoveryHook(t *testing.T) {
	// Create a test workflow environment
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	
	// Register a workflow that will panic
	env.RegisterWorkflow(func(ctx workflow.Context) error {
		// This will be recovered by the PanicRecoveryHook
		defer PanicRecoveryHook(ctx)
		
		// Cause a panic
		panic("test panic")
		
		// This line should not be reached
		return nil
	})
	
	// Execute the workflow - it should not fail due to the panic recovery
	env.ExecuteWorkflow(func(ctx workflow.Context) error {
		// This will be recovered by the PanicRecoveryHook
		defer PanicRecoveryHook(ctx)
		
		// Cause a panic
		panic("test panic")
		
		// This line should not be reached
		return nil
	})
	
	// Verify that the workflow completed without error
	assert.NoError(t, env.GetWorkflowError())
}

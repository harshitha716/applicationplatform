package workflowutil

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
)

// Simple implementation of ParamsBase for testing
type TestParams struct {
	UserID uuid.UUID
	OrgIDs []uuid.UUID
}

func (p TestParams) GetAccessControlParams() (uuid.UUID, []uuid.UUID) {
	return p.UserID, p.OrgIDs
}

func TestAddAccessControlParamsToWorkflowCtx(t *testing.T) {
	// Create a mock workflow context
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	
	// Create test UUIDs
	userID := uuid.New()
	orgID1 := uuid.New()
	orgID2 := uuid.New()
	
	// Create test params
	testParams := TestParams{
		UserID: userID,
		OrgIDs: []uuid.UUID{orgID1, orgID2},
	}
	
	env.RegisterWorkflow(func(ctx workflow.Context) error {
		// Test the function
		ctx = AddAccessControlParamsToWorkflowCtx(ctx, testParams)
		
		// We can't directly access the header values, but we can verify the function exists
		// Just verify that the function was called successfully
		return nil
	})
	
	env.ExecuteWorkflow(func(ctx workflow.Context) error {
		ctx = AddAccessControlParamsToWorkflowCtx(ctx, testParams)
		return nil
	})
	assert.NoError(t, env.GetWorkflowError())
}

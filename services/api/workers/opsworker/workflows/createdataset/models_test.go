package createdataset

import (
	"testing"

	"github.com/Zampfi/application-platform/services/api/core/datasets/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateDatasetWorkflowInitPayload_GetAccessControlParams(t *testing.T) {
	// Create test UUIDs
	orgId := uuid.New()
	userId := uuid.New()
	
	// Create test payload
	payload := CreateDatasetWorkflowInitPayload{
		RegisterDatasetPayload: models.DatasetCreationInfo{},
		OrganizationId:         orgId,
		UserId:                 userId,
	}
	
	// Call the function
	returnedOrgId, returnedUserIds := payload.GetAccessControlParams()
	
	// Verify the return values
	assert.Equal(t, &payload.OrganizationId, returnedOrgId)
	assert.Equal(t, []uuid.UUID{payload.UserId}, returnedUserIds)
}

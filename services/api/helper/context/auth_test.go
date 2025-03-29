package apicontext

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddAuthToCtx(t *testing.T) {
	ctx := context.Background()

	// get auth values from context
	role, userID, userOrganizations := GetAuthFromContext(ctx)
	assert.Equal(t, "anonymous", role)
	assert.Nil(t, userID)
	assert.Empty(t, userOrganizations)

	// add auth values to context
	r := "user"
	wId := uuid.New()
	uId := uuid.New()
	am := []uuid.UUID{wId}
	ctx = AddAuthToContext(ctx, r, uId, am)

	// get auth values from context
	userID = nil
	userOrganizations = []uuid.UUID{}
	role, userID, userOrganizations = GetAuthFromContext(ctx)
	assert.Equal(t, r, role)
	assert.NotNil(t, userID)
	assert.Equal(t, uId, *userID)
	assert.NotEmpty(t, userOrganizations)
	assert.Equal(t, wId, userOrganizations[0])

	// add user id but no allowed organizations
	ctx = AddAuthToContext(ctx, r, uId, []uuid.UUID{})
	userID = nil
	userOrganizations = []uuid.UUID{}
	role, userID, userOrganizations = GetAuthFromContext(ctx)
	assert.Equal(t, r, role)
	assert.NotNil(t, userID)
	assert.Equal(t, uId, *userID)
	assert.Empty(t, userOrganizations)

}

func TestAddAuthVariablesToGinContext(t *testing.T) {

	gin.SetMode(gin.TestMode)
	ginCtx, _ := gin.CreateTestContext(nil)

	// get auth values from context
	role := "user"
	userId := uuid.New()
	userOrganizations := []uuid.UUID{uuid.New()}

	// add auth values to gin context
	AddAuthToGinContext(ginCtx, role, userId, userOrganizations)

	// get auth values from gin context
	r, uId, wIds := GetAuthFromContext(ginCtx)
	assert.Equal(t, role, r)
	assert.NotNil(t, uId)
	assert.Equal(t, userId, *uId)
	assert.NotEmpty(t, wIds)
	assert.Equal(t, userOrganizations[0], wIds[0])

}

func TestIsZampEmail(t *testing.T) {
	assert.True(t, IsZampEmail("admin@zamp.finance"))
	assert.True(t, IsZampEmail("admin@zamp.ai"))
	assert.False(t, IsZampEmail("admin@zamp.com"))
}

func TestAddAuditInfoToContext(t *testing.T) {
	ctx := context.Background()
	email := "test@example.com"
	ipAddress := "192.168.1.1"
	userAgent := "TestUserAgent"
	
	// Add audit info to context
	enrichedCtx := AddAuditInfoToContext(ctx, email, ipAddress, userAgent)
	
	// Retrieve and verify audit info from context
	retrievedEmail, retrievedIP, retrievedAgent := GetAuditInfoFromContext(enrichedCtx)
	assert.Equal(t, email, retrievedEmail)
	assert.Equal(t, ipAddress, retrievedIP)
	assert.Equal(t, userAgent, retrievedAgent)
	
	// Test with empty values
	emptyCtx := AddAuditInfoToContext(ctx, "", "", "")
	emptyEmail, emptyIP, emptyAgent := GetAuditInfoFromContext(emptyCtx)
	assert.Equal(t, "", emptyEmail)
	assert.Equal(t, "", emptyIP)
	assert.Equal(t, "", emptyAgent)
	
	// Test with nil context
	nilCtx := context.Background()
	nilEmail, nilIP, nilAgent := GetAuditInfoFromContext(nilCtx)
	assert.Equal(t, "", nilEmail)
	assert.Equal(t, "", nilIP)
	assert.Equal(t, "", nilAgent)
}

func TestAddAuditInfoToGinContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ginCtx, _ := gin.CreateTestContext(nil)
	
	email := "test@example.com"
	ipAddress := "192.168.1.1"
	userAgent := "TestUserAgent"
	
	// Add audit info to gin context
	AddAuditInfoToGinContext(ginCtx, email, ipAddress, userAgent)
	
	// Retrieve and verify audit info from gin context
	retrievedEmail, retrievedIP, retrievedAgent := GetAuditInfoFromContext(ginCtx)
	assert.Equal(t, email, retrievedEmail)
	assert.Equal(t, ipAddress, retrievedIP)
	assert.Equal(t, userAgent, retrievedAgent)
}

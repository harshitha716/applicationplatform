package apicontext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockLogger is a mock implementation of zap.Logger for testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) With(fields ...zap.Field) *zap.Logger {
	args := m.Called(fields)
	return args.Get(0).(*zap.Logger)
}

func TestAddLoggerToContext(t *testing.T) {
	// Create a test logger
	testLogger := zap.NewNop()

	// Create a base context
	ctx := context.Background()

	// Add logger to context
	ctxWithLogger := AddLoggerToContext(ctx, testLogger)

	// Retrieve logger from context
	retrievedLogger := ctxWithLogger.Value(contextKeyLogger)
	assert.NotNil(t, retrievedLogger, "Logger should be added to context")
}

func TestGetLoggerFromCtx(t *testing.T) {
	// Scenario 1: Logger already in context

	ctx := context.Background()
	retrievedLogger := GetLoggerFromCtx(ctx)
	assert.NotNil(t, retrievedLogger, "Should return a logger when GetLogger succeeds")

	// Scenario 2: GetLogger fails

	ctx = context.Background()
	retrievedLogger = GetLoggerFromCtx(ctx)
	assert.NotNil(t, retrievedLogger, "Should return a no-op logger when GetLogger fails")

	// Scenario 3: Incorrect logger type in context
	ctx = context.WithValue(context.Background(), contextKeyLogger, "not a logger")
	retrievedLogger = GetLoggerFromCtx(ctx)
	assert.NotNil(t, retrievedLogger, "Should handle incorrect logger type gracefully")
}

func TestLoggerWithContextVariables(t *testing.T) {
	// Create a context with variables
	ctx := context.Background()
	ctx = AddCtxVariableToCtx(ctx, "key1", "value1")
	ctx = AddCtxVariableToCtx(ctx, "key2", 42)

	// Create a test logger
	testLogger := zap.NewNop()
	ctx = AddLoggerToContext(ctx, testLogger)

	// Retrieve logger
	retrievedLogger := GetLoggerFromCtx(ctx)
	assert.NotNil(t, retrievedLogger, "Logger should be retrieved with context variables")

	// Verify context variables are added to logger
	contextVars := GetAllCtxVariablesFromCtx(ctx)
	assert.Len(t, contextVars, 2, "Should have two context variables")
	assert.Equal(t, "value1", contextVars["key1"], "First context variable should match")
	assert.Equal(t, 42, contextVars["key2"], "Second context variable should match")
}

func TestGetFilteredCtxVariablesFromCtx(t *testing.T) {
	// Create a context with both sensitive and non-sensitive variables
	ctx := context.Background()
	ctx = AddCtxVariableToCtx(ctx, contextKeyUserID, "user-123")
	ctx = AddCtxVariableToCtx(ctx, contextKeyUserEmail, "test@example.com")
	ctx = AddCtxVariableToCtx(ctx, contextKeyUserIPAddress, "192.168.1.1")
	ctx = AddCtxVariableToCtx(ctx, contextKeyUserRole, "admin")
	
	// Get filtered variables
	filteredVars := GetFilteredCtxVariablesFromCtx(ctx)
	
	// Verify sensitive information is removed
	assert.NotContains(t, filteredVars, contextKeyUserEmail, "Email should be filtered out")
	assert.NotContains(t, filteredVars, contextKeyUserIPAddress, "IP address should be filtered out")
	
	// Verify non-sensitive information is preserved
	assert.Contains(t, filteredVars, contextKeyUserID, "User ID should be preserved")
	assert.Contains(t, filteredVars, contextKeyUserRole, "User role should be preserved")
}

func TestGetLoggerFromCtxWithFiltering(t *testing.T) {
	// Create a context with sensitive information
	ctx := context.Background()
	ctx = AddCtxVariableToCtx(ctx, contextKeyUserEmail, "test@example.com")
	ctx = AddCtxVariableToCtx(ctx, contextKeyUserRole, "admin")
	
	// Get logger with filtered context
	logger := GetLoggerFromCtx(ctx)
	
	// Cannot directly test the logger contents, but we're testing the integration
	// This test ensures the code runs without errors
	assert.NotNil(t, logger, "Logger should be created with filtered context")
}

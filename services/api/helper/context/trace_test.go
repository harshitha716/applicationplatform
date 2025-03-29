package apicontext

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAddTraceIdToGinContext(t *testing.T) {
	// Initialize a Gin context
	ginCtx := &gin.Context{
		Request: &http.Request{},
	}

	// Define a trace ID to add
	traceId := "test-trace-id"

	// Call the function to add the trace ID to the Gin context
	AddTraceIdToGinContext(ginCtx, traceId)

	// Retrieve the value directly from the Gin context
	ctxVars, exists := ginCtx.Get(contextKeyCtxVariables)
	assert.True(t, exists, "Context variables should exist in the Gin context")
	traceIdFromGinCtx := ctxVars.(map[string]interface{})[contextKeyTraceId]

	// Assertions
	assert.Equal(t, traceId, traceIdFromGinCtx, "Trace ID should match the added value")
}

func TestGetTraceIdFromContext(t *testing.T) {
	// Create a parent context
	parentCtx := context.Background()

	// Add the trace ID to the context using a helper function
	traceId := "test-trace-id"

	ctxVars := make(map[string]interface{})
	ctxVars[contextKeyTraceId] = traceId

	ctx := context.WithValue(parentCtx, contextKeyCtxVariables, ctxVars)

	// Call the function to retrieve the trace ID
	retrievedTraceId := GetTraceIdFromContext(ctx)

	// Assertions
	assert.Equal(t, traceId, retrievedTraceId, "Trace ID should match the added value")
}

func TestGetTraceIdFromContext_NoTraceId(t *testing.T) {
	// Create a parent context without a trace ID
	ctx := context.Background()

	// Call the function to retrieve the trace ID
	retrievedTraceId := GetTraceIdFromContext(ctx)

	// Assertions
	assert.Equal(t, "", retrievedTraceId, "Should return an empty string when trace ID is not present")
}

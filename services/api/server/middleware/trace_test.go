package middleware

import (
	"net/http/httptest"
	"testing"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetTraceGenMiddleware(t *testing.T) {
	// Test cases
	tests := []struct {
		name string
	}{
		{
			name: "should generate and inject trace ID",
		},
		{
			name: "should generate unique trace IDs for different requests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new gin context with a response writer
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Get the middleware function
			middleware := GetTraceGenMiddleware()

			// Execute the middleware
			middleware(c)

			// Get the trace ID from context
			traceID := apicontext.GetTraceIdFromContext(c)

			// Verify that a trace ID was set
			assert.NotEmpty(t, traceID)

			// Verify that it's a valid UUID
			_, err := uuid.Parse(traceID)
			assert.NoError(t, err)

			// Verify header was set correctly
			assert.Equal(t, traceID, w.Header().Get("X-Trace-Id"))

			if tt.name == "should generate unique trace IDs for different requests" {
				// Create a second context and verify uniqueness
				w2 := httptest.NewRecorder()
				c2, _ := gin.CreateTestContext(w2)
				middleware(c2)
				traceID2 := apicontext.GetTraceIdFromContext(c2)

				// Verify that the trace IDs are different
				assert.NotEqual(t, traceID, traceID2)
				assert.Equal(t, traceID2, w2.Header().Get("X-Trace-Id"))
			}
		})
	}
}

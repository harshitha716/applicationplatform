package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetCORSMiddleware(t *testing.T) {
	allowedOrigins := "http://localhost"
	middleware := GetCORSMiddleware(allowedOrigins)

	if middleware == nil {
		t.Errorf("Expected middleware to be a function, got nil")
	}

	// Test the middleware
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	// Create a test context
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/", nil)

	// Set the request method to OPTIONS
	c.Request.Method = "OPTIONS"

	// Call the middleware
	middleware(c)

	// Check if the response status is 204
	if c.Writer.Status() != 204 {
		t.Errorf("Expected status 204, got %d", c.Writer.Status())
	}

	// Set the request method to GET
	c.Request.Method = "GET"

	// Call the middleware
	middleware(c)

	// Check if the response headers are set
	assert.Equal(t, allowedOrigins, c.Writer.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, PATCH, DELETE, OPTIONS", c.Writer.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Accept, Origin, Referer, X-Zamp-Organization-Id", c.Writer.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "Location", c.Writer.Header().Get("Access-Control-Expose-Headers"))
	assert.Equal(t, "600", c.Writer.Header().Get("Access-Control-Max-Age"))
	assert.Equal(t, "true", c.Writer.Header().Get("Access-Control-Allow-Credentials"))

}

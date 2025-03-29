package helper

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockHTTPClient implements HTTPClient interface for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestForwardResponseHeaders(t *testing.T) {
	tests := []struct {
		name            string
		sourceHeaders   map[string][]string
		expectedHeaders map[string]string
	}{
		{
			name: "forward all non-excluded headers",
			sourceHeaders: map[string][]string{
				"Content-Type":    {"application/json"},
				"X-Custom-Header": {"value1", "value2"},
				"Content-Length":  {"100"},
				"Date":            {"Mon, 01 Jan 2024 00:00:00 GMT"},
			},
			expectedHeaders: map[string]string{
				"Content-Type":    "application/json",
				"X-Custom-Header": "value2", // Last value should be used
			},
		},
		{
			name:            "empty headers",
			sourceHeaders:   map[string][]string{},
			expectedHeaders: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create source response
			srcResp := &http.Response{
				Header: tt.sourceHeaders,
			}

			// Create Gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Call function
			ForwardResponseHeaders(srcResp, c)

			// Verify headers
			for key, expectedValue := range tt.expectedHeaders {
				assert.Equal(t, expectedValue, w.Header().Get(key))
			}

			// Verify excluded headers are not present
			assert.Empty(t, w.Header().Get("Content-Length"))
			assert.Empty(t, w.Header().Get("Date"))
		})
	}
}

func TestHttpGet(t *testing.T) {
	tests := []struct {
		name           string
		serverHandler  http.HandlerFunc
		mockClient     *MockHTTPClient
		headers        map[string]string
		expectedError  bool
		expectedStatus int
	}{
		{
			name: "successful get request",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString("success")),
					}
					return resp, nil
				},
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "client error",
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("client error")
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, resp, err := HttpGet(tt.mockClient, "http://example.com", tt.headers)

			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotEmpty(t, body)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}

	// Test invalid URL
	t.Run("invalid URL", func(t *testing.T) {
		mockClient := &MockHTTPClient{}
		_, _, err := HttpGet(mockClient, "!@#$%^&invalid-url", nil)
		assert.Error(t, err)
	})
}

func TestHttpPost(t *testing.T) {
	tests := []struct {
		name           string
		payload        interface{}
		mockClient     *MockHTTPClient
		headers        map[string]string
		expectedError  bool
		expectedStatus int
	}{
		{
			name:    "successful post request",
			payload: map[string]string{"key": "value"},
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewBufferString("success")),
					}
					return resp, nil
				},
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "client error",
			payload: map[string]string{"key": "value"},
			mockClient: &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("client error")
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, resp, err := HttpPost(tt.mockClient, "http://example.com", tt.payload, tt.headers)

			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotEmpty(t, body)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}

	// Test invalid payload
	t.Run("invalid payload", func(t *testing.T) {
		mockClient := &MockHTTPClient{}
		payload := make(chan int) // channels cannot be marshaled to JSON
		_, _, err := HttpPost(mockClient, "http://example.com", payload, nil)
		assert.Error(t, err)
	})

	// Test invalid URL
	t.Run("invalid URL", func(t *testing.T) {
		mockClient := &MockHTTPClient{}
		_, _, err := HttpPost(mockClient, "!@#$%^invalid-url", map[string]string{}, nil)
		assert.Error(t, err)
	})
}

func TestAddHeaders(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-Custom-Key": "custom-value",
	}

	addHeaders(req, headers)

	for key, value := range headers {
		assert.Equal(t, value, req.Header.Get(key))
	}
}

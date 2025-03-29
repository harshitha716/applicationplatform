package kratosclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func GetMockKratosServer(httpStatus int, responseBody string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// return a resposne with write status code
		w.Header().Set("Content-Type", "application/json")
		w.Header()["Content-Type"] = []string{"application/json"}
		w.WriteHeader(httpStatus)
		// write response body
		io.WriteString(w, responseBody)
	}))
}

func TestNewClient(t *testing.T) {

	publicUrl := "http://example.com/"
	client, err := NewClient(publicUrl)
	assert.Nil(t, err)
	assert.Equal(t, "http://example.com/", client.publicUrl.String())

	client, err = NewClient("badurl")
	assert.NotNil(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "invalid public url")

}

const sessionMock = `
 {
      "id": "fc59625f-8ad6-491a-b3e0-f5676a9232f3",
      "active": true,
      "expires_at": "2024-08-02T21:13:10.021766508Z",
      "authenticated_at": "2024-08-01T21:13:10.021766508Z",
      "authenticator_assurance_level": "aal1",
      "authentication_methods": [
          {
              "method": "password",
              "aal": "aal1",
              "completed_at": "2024-08-01T21:13:10.021757758Z"
          }
      ],
      "issued_at": "2024-08-01T21:13:10.021766508Z",
      "identity": {
          "id": "d93e9fb4-2451-4eb5-aa86-aaf017c74c39",
          "schema_id": "default",
          "schema_url": "http://auth:4433/schemas/ZGVmYXVsdA",
          "state": "active",
          "state_changed_at": "2024-08-01T19:59:14.935789Z",
          "traits": {
              "email": "rishichandra1@gmail.com"
          },
          "metadata_public": null,
          "created_at": "2024-08-01T19:59:14.937324Z",
          "updated_at": "2024-08-01T19:59:14.937324Z",
          "organization_id": null
      },
      "devices": [
          {
              "id": "3da58b8b-1f62-47aa-9a0e-e506470efc20",
              "ip_address": "192.168.65.1:58625",
              "user_agent": "PostmanRuntime/7.40.0",
              "location": ""
          }
      ]
  }
`

func TestGetSession(t *testing.T) {

	// logger := zap.NewNop()
	logger, _ := zap.NewProduction()

	// test internal server resposne from kratos
	server := GetMockKratosServer(http.StatusInternalServerError, `{ "error": {"message":"internal error"}}`)
	defer server.Close()

	client, err := NewClient(server.URL)
	assert.Nil(t, err)

	ctx := context.Background()

	_, httpResp, kerr := client.GetSessionInfo(ctx, logger, "cookie")
	assert.NotNil(t, kerr)
	assert.Equal(t, "500 Internal Server Error", kerr.Message)
	assert.Equal(t, http.StatusInternalServerError, httpResp.StatusCode)
	assert.Equal(t, http.StatusInternalServerError, kerr.Code)

	server.Close()

	// test a 400 response with custom message from kratos
	server = GetMockKratosServer(http.StatusBadRequest, `{ "error": {"message":"bad request"}}`)
	defer server.Close()

	client, err = NewClient(server.URL)
	assert.Nil(t, err)

	_, httpResp, kerr = client.GetSessionInfo(ctx, logger, "cookie")
	assert.NotNil(t, kerr)
	assert.Equal(t, "400 Bad Request", kerr.Message)
	assert.Equal(t, http.StatusBadRequest, httpResp.StatusCode)
	assert.Equal(t, http.StatusBadRequest, kerr.Code)
	server.Close()
	kerr = nil

	// test a successful response from kratos
	server = GetMockKratosServer(http.StatusOK, sessionMock)
	defer server.Close()

	client1, err := NewClient(server.URL)
	assert.Nil(t, err)

	session, httpResp, kerr := client1.GetSessionInfo(ctx, logger, "cookie")
	assert.Nil(t, kerr)
	assert.NotNil(t, session)
	assert.Equal(t, http.StatusOK, httpResp.StatusCode)
	assert.Equal(t, "fc59625f-8ad6-491a-b3e0-f5676a9232f3", session.Id)
	server.Close()
}

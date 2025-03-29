package pantheon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// CustomResponseRecorder extends httptest.ResponseRecorder with http.CloseNotifier
type CustomResponseRecorder struct {
	*httptest.ResponseRecorder
}

// CloseNotify implements http.CloseNotifier
func (r *CustomResponseRecorder) CloseNotify() <-chan bool {
	return make(<-chan bool)
}

func getMockServerConfig() *serverconfig.ServerConfig {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))

	}))

	config := serverconfig.GetEmptyServerConfig()
	config.Env.PantheonURL = server.URL

	return config
}

func TestPantheonReverseProxy(t *testing.T) {
	// Initialize Gin engine
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Get mock server config
	serverConfig := getMockServerConfig()

	// Initialize the reverse proxy
	RegisterPantheonRoutes(router.Group("/"), serverConfig)

	// Create a test request
	w := &CustomResponseRecorder{ResponseRecorder: httptest.NewRecorder()}
	req, err := http.NewRequest("POST", "/ai/*", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Serve http request
	router.ServeHTTP(w, req)

	// Check if the reverse proxy correctly modifies the request URL
	assert.Equal(t, 405, w.Code, "Expected HTTP status code 405")
}

package connectivity_catalog

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/db/models"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CustomResponseRecorder struct {
	*httptest.ResponseRecorder
}

// CloseNotify implements http.CloseNotifier
func (r *CustomResponseRecorder) CloseNotify() <-chan bool {
	return make(<-chan bool)
}

type testCase struct {
	name          string
	method        string
	path          string
	skip          bool
	outputPayload string
	statusCode    int
	initServerCfg func() *serverconfig.ServerConfig
}

func TestListConnectors(t *testing.T) {
	testCases := []testCase{
		{
			name:          "List connectors db error",
			method:        "GET",
			path:          "/connectors",
			outputPayload: `{"error":"something went wrong"}`,
			statusCode:    http.StatusInternalServerError,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				mockStore := mock_store.NewMockStore(t)
				mockStore.EXPECT().GetAllConnectors(mock.Anything).Return([]models.ConnectorWithActiveConnectionsCount{}, errors.New("something went wrong"))

				svCfg.Store = mockStore

				return svCfg
			},
		},
		{
			name:          "List connectors success",
			method:        "GET",
			path:          "/connectors",
			outputPayload: `[{"id":"00000000-0000-0000-0000-000000000000","name":"test","description":"test", "display_name":"test", "logo_url": "test", "category": "test", "status": "active"}]`,
			statusCode:    http.StatusOK,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				mockStore := mock_store.NewMockStore(t)
				mockStore.EXPECT().GetAllConnectors(mock.Anything).Return([]models.ConnectorWithActiveConnectionsCount{}, nil)

				svCfg.Store = mockStore

				return svCfg
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skip()
			}
			svCfg := tc.initServerCfg()

			e := gin.New()
			gin.SetMode(gin.TestMode)
			g := e.Group("/")
			RegisterCatalogRoutes(g, svCfg)

			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := &CustomResponseRecorder{ResponseRecorder: httptest.NewRecorder()}
			e.ServeHTTP(w, req)

			assert.Equal(t, tc.statusCode, w.Code)
		})
	}
}

func TestGetConnectorById(t *testing.T) {
	testCases := []testCase{
		{
			name:          "Get connector by id db error",
			method:        "GET",
			path:          "/connectors/00000000-0000-0000-0000-000000000000",
			outputPayload: `{"error":"something went wrong"}`,
			statusCode:    http.StatusInternalServerError,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				mockStore := mock_store.NewMockStore(t)
				mockStore.EXPECT().GetConnectorById(mock.Anything, mock.Anything).Return(nil, errors.New("something went wrong"))

				svCfg.Store = mockStore

				return svCfg
			},
		},
		{
			name:          "Get connector by id success",
			method:        "GET",
			path:          "/connectors/00000000-0000-0000-0000-000000000000",
			outputPayload: `{"id":"00000000-0000-0000-0000-000000000000","name":"test","description":"test", "display_name":"test", "logo_url": "test", "category": "test", "status": "active"}`,
			statusCode:    http.StatusOK,
			initServerCfg: func() *serverconfig.ServerConfig {
				svCfg := serverconfig.GetEmptyServerConfig()

				mockStore := mock_store.NewMockStore(t)
				mockStore.EXPECT().GetConnectorById(mock.Anything, mock.Anything).Return(&models.Connector{}, nil)

				svCfg.Store = mockStore

				return svCfg
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skip()
			}

			svCfg := tc.initServerCfg()

			e := gin.New()
			gin.SetMode(gin.TestMode)
			g := e.Group("/")
			RegisterCatalogRoutes(g, svCfg)

			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := &CustomResponseRecorder{ResponseRecorder: httptest.NewRecorder()}
			e.ServeHTTP(w, req)

			assert.Equal(t, tc.statusCode, w.Code)
		})
	}
}

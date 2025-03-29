package service

import (
	"context"
	"testing"

	databricksmockprovider "github.com/Zampfi/application-platform/services/api/mocks/pkg/dataplatform/providers/databricks"
	pinotmockprovider "github.com/Zampfi/application-platform/services/api/mocks/pkg/dataplatform/providers/pinot"
	postgresmockprovider "github.com/Zampfi/application-platform/services/api/mocks/pkg/dataplatform/providers/postgres"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/errors"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/providers/databricks"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/providers/pinot"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/providers/postgres"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	mockDatabricksService *databricksmockprovider.MockDatabricksService
	mockPinotService      *pinotmockprovider.MockPinotService
	mockPostgresService   *postgresmockprovider.MockPostgresService
	service               *providerService
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) SetupTest() {
	s.mockDatabricksService = new(databricksmockprovider.MockDatabricksService)
	s.mockPinotService = new(pinotmockprovider.MockPinotService)
	s.mockPostgresService = new(postgresmockprovider.MockPostgresService)
	s.service = &providerService{
		providerRegistry: &ProviderRegistry{
			databricksServices: map[string]databricks.DatabricksService{"databricks1": s.mockDatabricksService},
			pinotServices:      map[string]pinot.PinotService{"pinot1": s.mockPinotService},
			postgresServices:   map[string]postgres.PostgresService{"postgres1": s.mockPostgresService},
		},
	}
}

func (s *ServiceTestSuite) TestGetDatabricksService() {

	service, err := s.service.GetDatabricksService(context.Background(), "databricks1")
	s.Require().NoError(err)
	s.Require().NotNil(service)
	s.Equal(s.mockDatabricksService, service)

	_, err = s.service.GetDatabricksService(context.Background(), "nonexistent")
	s.Require().Error(err)
	s.EqualError(err, errors.ErrDatabricksServiceInitializationFailed.Error())
}

func (s *ServiceTestSuite) TestGetPinotService() {

	service, err := s.service.GetPinotService(context.Background(), "pinot1")
	s.Require().NoError(err)
	s.Require().NotNil(service)
	s.Equal(s.mockPinotService, service)

	_, err = s.service.GetPinotService(context.Background(), "nonexistent")
	s.Require().Error(err)
	s.EqualError(err, errors.ErrPinotServiceInitializationFailed.Error())
}

func (s *ServiceTestSuite) TestGetPostgresService() {

	service, err := s.service.GetPostgresService(context.Background(), "postgres1")
	s.Require().NoError(err)
	s.Require().NotNil(service)
	s.Equal(s.mockPostgresService, service)

	_, err = s.service.GetPostgresService(context.Background(), "nonexistent")
	s.Require().Error(err)
	s.EqualError(err, errors.ErrPostgresServiceInitializationFailed.Error())
}

func (s *ServiceTestSuite) TestGetService() {
	tests := []struct {
		providerType    constants.ProviderType
		dataProviderId  string
		expectedService interface{}
		expectedError   error
	}{
		{constants.ProviderTypeDatabricks, "databricks1", s.mockDatabricksService, nil},
		{constants.ProviderTypePinot, "pinot1", s.mockPinotService, nil},
		{constants.ProviderTypePostgres, "postgres1", s.mockPostgresService, nil},
		{constants.ProviderTypeDatabricks, "nonexistent", nil, errors.ErrProviderServiceNotFound},
		{constants.ProviderTypePinot, "nonexistent", nil, errors.ErrProviderServiceNotFound},
		{constants.ProviderTypePostgres, "nonexistent", nil, errors.ErrProviderServiceNotFound},
		{"unsupported", "anyId", nil, errors.ErrUnsupportedProviderType},
	}

	for _, tt := range tests {
		service, err := s.service.GetService(context.Background(), tt.providerType, tt.dataProviderId)

		if tt.expectedError != nil {
			s.Require().Error(err)
			s.EqualError(err, tt.expectedError.Error())
		} else {
			s.Require().NoError(err)
			s.Require().NotNil(service)
			s.Equal(tt.expectedService, service)
		}
	}
}

func (s *ServiceTestSuite) TestQuery() {

	ctx := context.Background()

	s.mockDatabricksService.On("Query", ctx, "table1", "SELECT * FROM table1", mock.Anything).Return(models.QueryResult{Rows: []map[string]interface{}{{"id": "result1", "name": "result2"}}}, nil).Once()

	result, err := s.service.Query(ctx, constants.ProviderTypeDatabricks, "databricks1", "table1", "SELECT * FROM table1")
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Equal(models.Rows{{"id": "result1", "name": "result2"}}, result.Rows)

	_, err = s.service.Query(ctx, constants.ProviderTypeDatabricks, "nonexistent", "table1", "SELECT * FROM table1")
	s.Require().Error(err)
	s.EqualError(err, errors.ErrProviderServiceNotFound.Error())

	s.mockDatabricksService.On("Query", ctx, "table1", "SELECT * FROM table1", mock.Anything).Return(models.QueryResult{}, errors.ErrQueryingDatabricks).Once()

	result, err = s.service.Query(ctx, constants.ProviderTypeDatabricks, "databricks1", "table1", "SELECT * FROM table1")

	s.Require().Error(err)
	s.EqualError(err, errors.ErrQueryingDatabricks.Error())
}

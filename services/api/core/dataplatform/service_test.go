package dataplatform

import (
	"context"
	"errors"
	"testing"

	datamodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
	mockdataservice "github.com/Zampfi/application-platform/services/api/mocks/core/dataplatform/data"
	"github.com/stretchr/testify/suite"
)

type DataPlatformServiceTestSuite struct {
	suite.Suite
	mockDataService *mockdataservice.MockDataService
	service         *dataPlatformService
}

func TestDataPlatformServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DataPlatformServiceTestSuite))
}

func (s *DataPlatformServiceTestSuite) SetupTest() {
	s.mockDataService = new(mockdataservice.MockDataService)
	s.service = &dataPlatformService{
		dataService: s.mockDataService,
	}
}

func (s *DataPlatformServiceTestSuite) TestGetDagsMultipleEdges() {
	// Test case for multiple edges
	s.mockDataService.EXPECT().
		GetDatasetEdgesByMerchant(context.Background(), "merchant1").
		Return([]datamodels.JobDatasetMapping{
			{
				SourceValue:      "dataset1",
				DestinationValue: "dataset2",
				JobId:            1,
				JobParams:        `{"param1": "value1"}`,
			},
			{
				SourceValue:      "dataset2",
				DestinationValue: "dataset3",
				JobId:            2,
				JobParams:        `{"param2": "value2"}`,
			},
		}, nil)

	got, err := s.service.GetDags(context.Background(), "merchant1")
	s.NoError(err)
	s.Equal(3, len(got)) // Three datasets in the DAG

	// Verify dataset1 node
	s.NotNil(got["dataset1"])
	s.Equal("dataset1", got["dataset1"].NodeId)
	s.Nil(got["dataset1"].Parents)

	// Verify dataset2 node
	s.NotNil(got["dataset2"])
	s.Equal("dataset2", got["dataset2"].NodeId)
	s.Len(got["dataset2"].Parents, 1)
	s.Equal("dataset1", got["dataset2"].Parents[0].NodeId)
	s.Equal(1, got["dataset2"].EdgeConfig["job_id"])

	// Verify dataset3 node
	s.NotNil(got["dataset3"])
	s.Equal("dataset3", got["dataset3"].NodeId)
	s.Len(got["dataset3"].Parents, 1)
	s.Equal("dataset2", got["dataset3"].Parents[0].NodeId)
	s.Equal(2, got["dataset3"].EdgeConfig["job_id"])
}

func (s *DataPlatformServiceTestSuite) TestGetDagsSingleEdge() {
	s.mockDataService.EXPECT().
		GetDatasetEdgesByMerchant(context.Background(), "merchant1").
		Return([]datamodels.JobDatasetMapping{
			{
				SourceValue:      "dataset1",
				DestinationValue: "dataset2",
				JobId:            1,
				JobParams:        `{"param1": "value1"}`,
			},
		}, nil)

	got, err := s.service.GetDags(context.Background(), "merchant1")
	s.NoError(err)
	s.Equal(2, len(got)) // Two datasets in the DAG

	// Verify dataset1 node
	s.NotNil(got["dataset1"])
	s.Equal("dataset1", got["dataset1"].NodeId)
	s.Nil(got["dataset1"].Parents)

	// Verify dataset2 node
	s.NotNil(got["dataset2"])
	s.Equal("dataset2", got["dataset2"].NodeId)
	s.Len(got["dataset2"].Parents, 1)
	s.Equal("dataset1", got["dataset2"].Parents[0].NodeId)
	s.Equal(1, got["dataset2"].EdgeConfig["job_id"])
}

func (s *DataPlatformServiceTestSuite) TestGetDagsError() {
	s.mockDataService.EXPECT().
		GetDatasetEdgesByMerchant(context.Background(), "merchant1").
		Return(nil, errors.New("database error"))

	got, err := s.service.GetDags(context.Background(), "merchant1")
	s.Error(err)
	s.Nil(got)
}

func (s *DataPlatformServiceTestSuite) TestGetDagsInvalidJSON() {
	s.mockDataService.EXPECT().
		GetDatasetEdgesByMerchant(context.Background(), "merchant1").
		Return([]datamodels.JobDatasetMapping{
			{
				SourceValue:      "dataset1",
				DestinationValue: "dataset2",
				JobId:            1,
				JobParams:        `invalid json`,
			},
		}, nil)

	got, err := s.service.GetDags(context.Background(), "merchant1")
	s.Error(err)
	s.Nil(got)
}

func (s *DataPlatformServiceTestSuite) TestGetDagsEmptyResult() {
	s.mockDataService.EXPECT().
		GetDatasetEdgesByMerchant(context.Background(), "merchant1").
		Return([]datamodels.JobDatasetMapping{}, nil)

	got, err := s.service.GetDags(context.Background(), "merchant1")
	s.NoError(err)
	s.Empty(got)
}

func (s *DataPlatformServiceTestSuite) TestGetDagsMultipleConvergingPaths() {
	s.mockDataService.EXPECT().
		GetDatasetEdgesByMerchant(context.Background(), "merchant1").
		Return([]datamodels.JobDatasetMapping{
			{
				SourceValue:      "dataset1",
				DestinationValue: "dataset3",
				JobId:            1,
				JobParams:        `{"param1": "value1"}`,
			},
			{
				SourceValue:      "dataset2",
				DestinationValue: "dataset3",
				JobId:            2,
				JobParams:        `{"param2": "value2"}`,
			},
			{
				SourceValue:      "dataset3",
				DestinationValue: "dataset4",
				JobId:            3,
				JobParams:        `{"param3": "value3"}`,
			},
			{
				SourceValue:      "dataset5",
				DestinationValue: "dataset4",
				JobId:            4,
				JobParams:        `{"param4": "value4"}`,
			},
		}, nil)

	got, err := s.service.GetDags(context.Background(), "merchant1")
	s.NoError(err)
	s.Equal(5, len(got)) // Five datasets in the DAG

	// Verify root nodes (dataset1, dataset2, dataset5)
	for _, NodeId := range []string{"dataset1", "dataset2", "dataset5"} {
		s.NotNil(got[NodeId])
		s.Equal(NodeId, got[NodeId].NodeId)
		s.Nil(got[NodeId].Parents)
	}

	// Verify dataset3 with two parents
	s.NotNil(got["dataset3"])
	s.Equal("dataset3", got["dataset3"].NodeId)
	s.Len(got["dataset3"].Parents, 2)
	parentIds := []string{
		got["dataset3"].Parents[0].NodeId,
		got["dataset3"].Parents[1].NodeId,
	}
	s.Contains(parentIds, "dataset1")
	s.Contains(parentIds, "dataset2")

	// Verify dataset4 with two parents
	s.NotNil(got["dataset4"])
	s.Equal("dataset4", got["dataset4"].NodeId)
	s.Len(got["dataset4"].Parents, 2)
	parentIds = []string{
		got["dataset4"].Parents[0].NodeId,
		got["dataset4"].Parents[1].NodeId,
	}
	s.Contains(parentIds, "dataset3")
	s.Contains(parentIds, "dataset5")
}

package dataplatform

import (
	"context"
	"encoding/json"
	"fmt"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	actions "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions"
	actionconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/constants"
	actionmodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/models"
	data "github.com/Zampfi/application-platform/services/api/core/dataplatform/data"
	dataconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	datamodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform/errors"
	servicemodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/constants"
	models "github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
	"go.uber.org/zap"
)

type DataPlatformService interface {
	QueryRealTime(ctx context.Context, merchantId string, query string, params map[string]string, args ...interface{}) (models.QueryResult, error)
	Query(ctx context.Context, merchantId string, query string, params map[string]string, args ...interface{}) (models.QueryResult, error)
	GetDatasetMetadata(ctx context.Context, merchantId string, datasetId string) (datamodels.DatasetMetadata, error)
	GetDatasetParents(ctx context.Context, merchantId string, datasetId string) (datamodels.DatasetParents, error)
	CreateMV(ctx context.Context, payload servicemodels.CreateMVPayload) (actionmodels.CreateActionResponse, error)
	GetActionById(ctx context.Context, merchantId string, actionId string) (actionmodels.Action, error)
	UpdateDatasetData(ctx context.Context, payload servicemodels.UpdateDatasetDataPayload) (actionmodels.CreateActionResponse, error)
	UpdateAction(ctx context.Context, jobStatusUpdate servicemodels.DatabricksJobStatusUpdatePayload) (actionmodels.Action, error)
	RegisterDataset(ctx context.Context, payload servicemodels.RegisterDatasetPayload) (actionmodels.CreateActionResponse, error)
	RegisterJob(ctx context.Context, payload servicemodels.RegisterJobPayload) (actionmodels.CreateActionResponse, error)
	UpsertTemplate(ctx context.Context, payload servicemodels.UpsertTemplatePayload) (actionmodels.CreateActionResponse, error)
	UpdateDataset(ctx context.Context, payload servicemodels.UpdateDatasetPayload) (actionmodels.CreateActionResponse, error)
	CopyDataset(ctx context.Context, payload servicemodels.CopyDatasetPayload) (actionmodels.CreateActionResponse, error)
	GetDags(ctx context.Context, merchantId string) (map[string]*servicemodels.DAGNode, error)
	DeleteDataset(ctx context.Context, payload servicemodels.DeleteDatasetPayload) (string, error)
}

type dataPlatformService struct {
	dataService   data.DataService
	actionService actions.ActionService
}

func InitDataPlatformService(dataPlatformConfig *serverconfig.DataPlatformConfig) (DataPlatformService, error) {
	dataService, err := data.InitDataService(dataPlatformConfig)
	if err != nil {
		return nil, err
	}
	actionService := actions.InitActionService(dataService)

	return &dataPlatformService{
		dataService:   dataService,
		actionService: actionService,
	}, nil
}

func (s *dataPlatformService) GetDataProviderIdForMerchant(merchantId string, providerType constants.ProviderType) (string, error) {
	return s.dataService.GetDataProviderIdForMerchant(merchantId, providerType)
}

func (s *dataPlatformService) QueryRealTime(ctx context.Context, merchantId string, query string, params map[string]string, args ...interface{}) (models.QueryResult, error) {
	return s.dataService.QueryRealTime(ctx, merchantId, query, params, args...)
}

func (s *dataPlatformService) Query(ctx context.Context, merchantId string, query string, params map[string]string, args ...interface{}) (models.QueryResult, error) {
	return s.dataService.Query(ctx, merchantId, query, params, args...)
}

func (s *dataPlatformService) GetDatasetMetadata(ctx context.Context, merchantId string, datasetId string) (datamodels.DatasetMetadata, error) {
	return s.dataService.GetDatasetMetadata(ctx, merchantId, datasetId)
}

func (s *dataPlatformService) GetDatasetParents(ctx context.Context, merchantId string, datasetId string) (datamodels.DatasetParents, error) {
	return s.dataService.GetDatasetParents(ctx, merchantId, datasetId)
}

func (s *dataPlatformService) CreateMV(ctx context.Context, payload servicemodels.CreateMVPayload) (actionmodels.CreateActionResponse, error) {
	createActionPayload := actionmodels.CreateActionPayload{
		MerchantID:            payload.MerchantID,
		ActionType:            actionconstants.ActionTypeCreateMV,
		ActionMetadataPayload: payload.ActionMetadataPayload,
		ActorId:               payload.ActorId,
	}
	return s.actionService.CreateAction(ctx, createActionPayload)
}

func (s *dataPlatformService) UpdateDatasetData(ctx context.Context, payload servicemodels.UpdateDatasetDataPayload) (actionmodels.CreateActionResponse, error) {
	createActionPayload := actionmodels.CreateActionPayload{
		MerchantID:            payload.MerchantID,
		ActionType:            actionconstants.ActionTypeUpdateDatasetData,
		ActionMetadataPayload: payload.ActionMetadataPayload,
		ActorId:               payload.ActorId,
	}
	return s.actionService.CreateAction(ctx, createActionPayload)
}

func (s *dataPlatformService) UpdateAction(ctx context.Context, jobStatusUpdate servicemodels.DatabricksJobStatusUpdatePayload) (actionmodels.Action, error) {
	return s.actionService.UpdateAction(ctx, jobStatusUpdate)
}

func (s *dataPlatformService) GetActionById(ctx context.Context, merchantId string, actionId string) (actionmodels.Action, error) {
	return s.actionService.GetActionById(ctx, merchantId, actionId)
}

func (s *dataPlatformService) RegisterDataset(ctx context.Context, payload servicemodels.RegisterDatasetPayload) (actionmodels.CreateActionResponse, error) {
	createActionPayload := actionmodels.CreateActionPayload{
		MerchantID:            payload.MerchantID,
		ActionType:            actionconstants.ActionTypeRegisterDataset,
		ActionMetadataPayload: payload.ActionMetadataPayload,
		ActorId:               payload.ActorId,
	}
	return s.actionService.CreateAction(ctx, createActionPayload)
}

func (s *dataPlatformService) RegisterJob(ctx context.Context, payload servicemodels.RegisterJobPayload) (actionmodels.CreateActionResponse, error) {
	createActionPayload := actionmodels.CreateActionPayload{
		MerchantID:            payload.MerchantID,
		ActionType:            actionconstants.ActionTypeRegisterJob,
		ActionMetadataPayload: payload.ActionMetadataPayload,
		ActorId:               payload.ActorId,
	}
	return s.actionService.CreateAction(ctx, createActionPayload)
}

func (s *dataPlatformService) UpsertTemplate(ctx context.Context, payload servicemodels.UpsertTemplatePayload) (actionmodels.CreateActionResponse, error) {
	createActionPayload := actionmodels.CreateActionPayload{
		MerchantID:            payload.MerchantID,
		ActionType:            actionconstants.ActionTypeUpsertTemplate,
		ActionMetadataPayload: payload.ActionMetadataPayload,
		ActorId:               payload.ActorId,
	}
	return s.actionService.CreateAction(ctx, createActionPayload)
}

func (s *dataPlatformService) UpdateDataset(ctx context.Context, payload servicemodels.UpdateDatasetPayload) (actionmodels.CreateActionResponse, error) {
	createActionPayload := actionmodels.CreateActionPayload{
		MerchantID:            payload.MerchantID,
		ActionType:            actionconstants.ActionTypeUpdateDataset,
		ActionMetadataPayload: payload.ActionMetadataPayload,
		ActorId:               payload.ActorId,
	}
	return s.actionService.CreateAction(ctx, createActionPayload)
}

func (s *dataPlatformService) CopyDataset(ctx context.Context, payload servicemodels.CopyDatasetPayload) (actionmodels.CreateActionResponse, error) {
	createActionPayload := actionmodels.CreateActionPayload{
		MerchantID:            payload.MerchantID,
		ActionType:            actionconstants.ActionTypeCopyDataset,
		ActionMetadataPayload: payload.ActionMetadataPayload,
		ActorId:               payload.ActorId,
	}
	return s.actionService.CreateAction(ctx, createActionPayload)
}

func (s *dataPlatformService) DeleteDataset(ctx context.Context, payload servicemodels.DeleteDatasetPayload) (string, error) {
	createActionPayload := actionmodels.CreateActionPayload{
		MerchantID:            payload.MerchantID,
		ActionType:            actionconstants.ActionTypeDeleteDataset,
		ActionMetadataPayload: payload,
		ActorId:               payload.ActorId,
	}
	response, err := s.actionService.CreateAction(ctx, createActionPayload)
	if err != nil {
		return "", err
	}
	return response.ActionID, nil
}

func (s *dataPlatformService) GetDags(ctx context.Context, merchantId string) (map[string]*servicemodels.DAGNode, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	jobsMappings, err := s.dataService.GetDatasetEdgesByMerchant(ctx, merchantId)
	if err != nil {
		logger.Error("failed to get dataset edges by merchant", zap.Error(err), zap.String("merchant_id", merchantId))
		return nil, err
	}

	dags := make(map[string]*servicemodels.DAGNode)
	for _, jobMapping := range jobsMappings {
		if _, ok := dags[jobMapping.DestinationValue]; !ok {
			dags[jobMapping.DestinationValue] = &servicemodels.DAGNode{
				NodeId:   jobMapping.DestinationValue,
				NodeType: servicemodels.NodeTypeDataset,
			}
		}

		if _, ok := dags[jobMapping.SourceValue]; !ok {
			if jobMapping.SourceType == string(dataconstants.JobMappingTypeFolder) {
				dags[jobMapping.SourceValue] = &servicemodels.DAGNode{
					NodeId:   jobMapping.SourceValue,
					NodeType: servicemodels.NodeTypeFolder,
				}
			} else {
				dags[jobMapping.SourceValue] = &servicemodels.DAGNode{
					NodeId:   jobMapping.SourceValue,
					NodeType: servicemodels.NodeTypeDataset,
				}
			}
		}

		dags[jobMapping.DestinationValue].Parents = append(dags[jobMapping.DestinationValue].Parents, dags[jobMapping.SourceValue])

		jobParams := make(map[string]interface{})
		err = json.Unmarshal([]byte(jobMapping.JobParams), &jobParams)
		if err != nil {
			logger.Error("failed to unmarshal job params", zap.Error(err), zap.String("merchant_id", merchantId), zap.String("job_params", jobMapping.JobParams))
			return nil, err
		}
		jobParams[dataconstants.JobMappingJobIdColumnName] = jobMapping.JobId

		dags[jobMapping.DestinationValue].EdgeConfig = jobParams
	}

	if err := s.detectCycle(dags); err != nil {
		return nil, err
	}

	return dags, nil
}

func (s *dataPlatformService) detectCycle(dags map[string]*servicemodels.DAGNode) error {
	visited := make(map[string]bool)
	stack := make(map[string]bool)

	for node := range dags {
		if !visited[node] {
			if err := s.dfs(node, visited, stack, dags); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *dataPlatformService) dfs(node string, visited map[string]bool, stack map[string]bool, dags map[string]*servicemodels.DAGNode) error {
	visited[node] = true
	stack[node] = true

	for _, parent := range dags[node].Parents {
		if !visited[parent.NodeId] {
			if err := s.dfs(parent.NodeId, visited, stack, dags); err != nil {
				return err
			}
		} else if stack[parent.NodeId] {
			return fmt.Errorf("cycle detected in node %s dags: %w", node, errors.ErrDAGCyclic)
		}
	}
	stack[node] = false
	return nil
}

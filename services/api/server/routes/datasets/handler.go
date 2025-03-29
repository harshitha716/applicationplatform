package datasets

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	actionmodels "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/models"
	dataplatformDataModels "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/models"
	"github.com/Zampfi/application-platform/services/api/core/datasets/models"
	datasetservice "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	"github.com/Zampfi/application-platform/services/api/core/fileimports"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	querybuildermodels "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apictx "github.com/Zampfi/application-platform/services/api/helper/context"

	"github.com/Zampfi/application-platform/services/api/server/middleware"
	"github.com/Zampfi/application-platform/services/api/server/routes/datasets/dtos"
)

func GetFilterConfig(c *gin.Context, svc datasetservice.DatasetService) {
	ctx := c.MustGet("datasetContext").(middleware.DatasetContext)
	filterConfig, datasetConfig, err := svc.GetFilterConfigByDatasetId(c, ctx.MerchantID, ctx.DatasetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dtos.FilterConfigResponse{
		Data:   make([]dtos.FilterConfig, len(filterConfig)),
		Config: datasetConfig,
	}
	for i, config := range filterConfig {
		response.Data[i].FromModel(config)
	}
	c.JSON(http.StatusOK, response)
}

func GetData(c *gin.Context, svc datasetservice.DatasetService) {
	ctx := c.MustGet("datasetContext").(middleware.DatasetContext)

	jsonStringConfig := c.Query("query_config")
	if jsonStringConfig == "" {
		jsonStringConfig = "{}"
	}

	var getDataRequest dtos.GetDataRequest
	if err := json.Unmarshal([]byte(jsonStringConfig), &getDataRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format: " + err.Error()})
		return
	}

	queryConfig := getDataRequest.ToModel()

	data, err := svc.GetDataByDatasetId(c, ctx.MerchantID, ctx.DatasetID, queryConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	description := ""
	if data.Description != nil {
		description = *data.Description
	}

	response := dtos.GetDataResponse{
		Title:       data.Title,
		Description: description,
		Data:        dtos.DatasetData{},
	}
	response.Data.FromModel(data)
	c.JSON(http.StatusOK, response)
}

func UpdateDatasetData(c *gin.Context, svc datasetservice.DatasetService) {
	ctx := c.MustGet("datasetContext").(middleware.DatasetContext)

	userId := ctx.UserID

	if userId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID not found"})
		return
	}

	var updateDatasetDataRequest dtos.UpdateDatasetDataRequest
	if err := c.ShouldBindJSON(&updateDatasetDataRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateDatasetParams := updateDatasetDataRequest.ToModel(*userId)

	datasetId, err := uuid.Parse(ctx.DatasetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dataset id"})
		return
	}

	datasetAction, err := svc.UpdateDatasetData(c, ctx.MerchantID, datasetId, updateDatasetParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dtos.DatasetAction{}
	response.FromModel(datasetAction)

	c.JSON(http.StatusOK, response)
}

func GetRowDetails(c *gin.Context, svc datasetservice.DatasetService) {
	ctx := c.MustGet("datasetContext").(middleware.DatasetContext)

	rowUUID := c.Param("rowUUID")

	parentDatasetInfo, err := svc.GetRowDetailsByUUID(c, ctx.MerchantID, ctx.DatasetID, rowUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dtos.ParentDatasetInfo{}
	response.FromModel(parentDatasetInfo)

	c.JSON(http.StatusOK, response)
}

func GetDatasetListing(c *gin.Context, svc datasetservice.DatasetService) {
	_, _, merchantIds := apictx.GetAuthFromContext(c)
	if len(merchantIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid merchant id"})
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pagesize"))
	if err != nil {
		pageSize = 100
	}

	var sortParams []models.DatasetListingSortParams
	sortParamsString := c.Query("sort")
	if sortParamsString != "" {
		err := json.Unmarshal([]byte(sortParamsString), &sortParams)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sort params"})
			return
		}
	}

	merchantId := merchantIds[0]

	datasetsModel, err := svc.GetDatasetListing(c, merchantId, models.DatsetListingParams{
		Pagination: &querybuildermodels.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
		SortParams: sortParams,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var totalCount int64
	if page == 1 {
		totalCount, err = svc.GetDatasetCount(c, merchantId, models.DatsetListingParams{
			Pagination: &querybuildermodels.Pagination{
				Page:     page,
				PageSize: pageSize,
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	var response dtos.GetDatasetListingResponse
	var datasets []dtos.DatasetListing
	for _, dataset := range datasetsModel {
		datasetListing := dtos.DatasetListing{}
		datasetListing.FromModel(dataset)
		datasets = append(datasets, datasetListing)
	}
	response.Datasets = datasets
	response.TotalCount = totalCount
	c.JSON(http.StatusOK, response)
}

func RegisterDataset(c *gin.Context, svc datasetservice.DatasetService) {
	_, userId, merchantIds := apictx.GetAuthFromContext(c)

	var registerDatasetRequest dtos.RegisterDatasetRequest
	if err := c.ShouldBindJSON(&registerDatasetRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(merchantIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid merchant id"})
		return
	}

	merchantId := merchantIds[0]

	datasetCreationInfo := registerDatasetRequest.ToModel()
	actionId, datasetId, err := svc.RegisterDataset(c, merchantId, *userId, datasetCreationInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dataset_id": datasetId, "action_id": actionId})
}

func CopyDataset(c *gin.Context, svc datasetservice.DatasetService) {
	_, userId, merchantIds := apictx.GetAuthFromContext(c)

	if len(merchantIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid merchant id"})
		return
	}

	if userId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	merchantId := merchantIds[0]

	var params models.CopyDatasetParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actionId, datasetId, err := svc.CopyDataset(c, merchantId, *userId, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dataset_id": datasetId, "action_id": actionId})
}

func validateDatasetConfig(config *dataplatformDataModels.DatasetConfig) (bool, map[string]bool) {
	if config == nil {
		return true, nil
	}

	hasColumns := config.Columns != nil
	hasCustomGroups := config.CustomColumnGroups != nil
	hasRules := config.Rules != nil

	validation := map[string]bool{
		"has_columns":       hasColumns,
		"has_custom_groups": hasCustomGroups,
		"has_rules":         hasRules,
	}

	return hasColumns && hasCustomGroups && hasRules, validation
}

func UpdateDataset(c *gin.Context, svc datasetservice.DatasetService) {
	ctx := c.MustGet("datasetContext").(middleware.DatasetContext)

	var request dtos.UpdateDatasetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if isValid, validation := validateDatasetConfig(request.DatasetConfig); !isValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "dataset config must include columns, custom column groups, and rules",
			"validation": validation,
		})
		return
	}

	actionId, err := svc.UpdateDataset(c, ctx.MerchantID, ctx.DatasetID, request.ToModel())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"action_id": actionId})
}

func RegisterDatasetJob(c *gin.Context, svc datasetservice.DatasetService) {
	_, _, merchantIds := apictx.GetAuthFromContext(c)

	if len(merchantIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid merchant id"})
		return
	}

	merchantId := merchantIds[0]

	var rawJSON map[string]interface{}
	if err := c.ShouldBindJSON(&rawJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rawJSON["merchant_id"] = merchantId

	modifiedJSON, err := json.Marshal(rawJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	var jobInfo actionmodels.RegisterJobActionPayload
	if err := json.Unmarshal(modifiedJSON, &jobInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actionId, err := svc.RegisterDatasetJob(c, merchantId, jobInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"action_id": actionId})
}

func UpsertTemplate(c *gin.Context, svc datasetservice.DatasetService) {
	_, _, merchantIds := apictx.GetAuthFromContext(c)

	if len(merchantIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid merchant id"})
		return
	}

	merchantId := merchantIds[0]

	var upsertTemplateConfig actionmodels.UpsertTemplateActionPayload
	if err := c.ShouldBindJSON(&upsertTemplateConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actionId, err := svc.UpsertTemplate(c, merchantId, upsertTemplateConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"action_id": actionId})
}

func GetDatasetAudiences(c *gin.Context, svc datasetservice.DatasetService) {

	datasetIdStr := c.Param("datasetId")

	datasetId, err := uuid.Parse(datasetIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dataset id"})
		return
	}

	audiences, err := svc.GetDatasetAudiences(c, datasetId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, audiences)
}

func GetDatasetActions(c *gin.Context, svc datasetservice.DatasetService) {
	_, _, merchantIds := apictx.GetAuthFromContext(c)

	if len(merchantIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid merchant id"})
		return
	}

	merchantId := merchantIds[0]

	datasetIdStr := c.Param("datasetId")

	datasetId, err := uuid.Parse(datasetIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dataset id"})
		return
	}

	params := parseActionQueryParams(c)

	actions, err := svc.GetDatasetActions(c, merchantId, params.ToModel(datasetId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dtos.DatasetAction
	for _, action := range actions {
		datasetAction := dtos.DatasetAction{}
		datasetAction.FromModel(action)
		response = append(response, datasetAction)
	}

	c.JSON(http.StatusOK, response)
}

func parseActionQueryParams(c *gin.Context) dtos.DatasetActionQueryParams {
	var params dtos.DatasetActionQueryParams

	actionIdsStr := c.Query("action_ids")
	actionTypeStr := c.Query("action_type")
	actionByStr := c.Query("action_by")
	statusStr := c.Query("status")

	if actionIdsStr != "" {
		params.ActionIds = strings.Split(actionIdsStr, ",")
	}
	if actionTypeStr != "" {
		params.ActionType = strings.Split(actionTypeStr, ",")
	}
	if actionByStr != "" {
		actionByStrings := strings.Split(actionByStr, ",")
		for _, id := range actionByStrings {
			if uid, err := uuid.Parse(id); err == nil {
				params.ActionBy = append(params.ActionBy, uid)
			}
		}
	}
	if statusStr != "" {
		params.Status = strings.Split(statusStr, ",")
	}

	return params
}

func addDatasetAudiences(c *gin.Context, svc datasetservice.DatasetService) {
	datasetId, err := uuid.Parse(c.Param("datasetId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dataset ID"})
		return
	}

	var payload dtos.BulkAddAudienceRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bulkAddAudiencePayload models.BulkAddDatasetAudiencePayload
	for _, audience := range payload.Audiences {
		bulkAddAudiencePayload.Audiences = append(bulkAddAudiencePayload.Audiences, models.AddDatasetAudiencePayload{
			AudienceId:   audience.AudienceId,
			AudienceType: audience.AudienceType,
			Privilege:    audience.Role,
		})
	}

	response, errs := svc.BulkAddAudienceToDataset(c, datasetId, bulkAddAudiencePayload)
	if errs.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errs.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"audiences": response, "audience_errors": errs.Audiences})
}

func updateDatasetAudience(c *gin.Context, svc datasetservice.DatasetService) {
	datasetId, err := uuid.Parse(c.Param("datasetId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dataset ID"})
		return
	}

	var payload dtos.UpdateAudienceRoleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := svc.UpdateDatasetAudiencePrivilege(c, datasetId, payload.AudiencId, dbmodels.ResourcePrivilege(payload.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func deleteDatasetAudience(c *gin.Context, svc datasetservice.DatasetService) {
	datasetId, err := uuid.Parse(c.Param("datasetId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dataset ID"})
		return
	}

	var payload dtos.DeleteAudienceRoleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = svc.RemoveAudienceFromDataset(c, datasetId, payload.AudiencId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func GetRulesByDatasetColumns(c *gin.Context, svc datasetservice.DatasetService) {
	_, _, organizationIds := apictx.GetAuthFromContext(c)

	if len(organizationIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	organizationId := organizationIds[0]

	datasetColumnsString := c.Query("dataset_columns")
	if datasetColumnsString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dataset_columns parameter is required"})
		return
	}

	var request dtos.GetRulesByDatasetColumnsRequest
	if err := json.Unmarshal([]byte(datasetColumnsString), &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format: " + err.Error()})
		return
	}

	datasetColumns := request.ToModel()
	rules, err := svc.GetRulesByDatasetColumns(c, organizationId, datasetColumns)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rules)
}

func GetRulesByIds(c *gin.Context, svc datasetservice.DatasetService) {
	_, _, organizationIds := apictx.GetAuthFromContext(c)

	if len(organizationIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	ruleIdsString := c.Query("rule_ids")
	if ruleIdsString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rule_ids parameter is required"})
		return
	}

	ruleIds := strings.Split(ruleIdsString, ",")

	rules, err := svc.GetRulesByIds(c, ruleIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rules)
}

func UpdateRulesPriority(c *gin.Context, svc datasetservice.DatasetService) {
	_, userId, organizationIds := apictx.GetAuthFromContext(c)

	if len(organizationIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	if userId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	organizationId := organizationIds[0]

	var request models.UpdateRulePriorityParams
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	datasetAction, err := svc.UpdateRulePriority(c, organizationId, *userId, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dtos.DatasetAction{}
	response.FromModel(datasetAction)

	c.JSON(http.StatusOK, response)
}

func CreateDatasetExportAction(c *gin.Context, svc datasetservice.DatasetService) {
	ctx := c.MustGet("datasetContext").(middleware.DatasetContext)
	_, userId, _ := apictx.GetAuthFromContext(c)

	jsonStringConfig := c.Query("query_config")
	if jsonStringConfig == "" {
		jsonStringConfig = "{}"
	}

	var getDataRequest dtos.GetDataRequest
	if err := json.Unmarshal([]byte(jsonStringConfig), &getDataRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format: " + err.Error()})
		return
	}

	queryConfig := getDataRequest.ToModel()

	workflowID, err := svc.CreateDatasetExportAction(c, ctx.MerchantID, ctx.DatasetID, queryConfig, *userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"workflow_id": workflowID})
}

func DeleteDataset(c *gin.Context, svc datasetservice.DatasetService) {
	ctx := c.MustGet("datasetContext").(middleware.DatasetContext)

	actionId, err := svc.DeleteDataset(c, ctx.MerchantID, ctx.DatasetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"action_id": actionId})
}

func GetDownloadableDataExportUrl(c *gin.Context, svc datasetservice.DatasetService) {
	workflowId := c.Param("workflowId")

	signedURL, err := svc.GetDownloadableDataExportUrl(c, workflowId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"signed_url": signedURL})

}

func InitDatasetFileImport(c *gin.Context, svc datasetservice.DatasetService, fileUploadService fileimports.FileImportService) {
	_, _, organizationIds := apictx.GetAuthFromContext(c)

	if len(organizationIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	var request dtos.InitiateDatasetFileImportRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := fileUploadService.InitiateFileImport(c, organizationIds[0], request.FileName, request.FileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func AckDatasetFileImport(c *gin.Context, svc datasetservice.DatasetService, fileUploadService fileimports.FileImportService) {
	_, _, organizationIds := apictx.GetAuthFromContext(c)

	if len(organizationIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	fileUploadId, err := uuid.Parse(c.Param("fileUploadId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file upload id"})
		return
	}

	var request dtos.AckDatasetFileImportRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	_, err = fileUploadService.AcknowledgeFileImportCompletion(c, fileUploadId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	datasetActionId, err := svc.InitiateFilePreparationForDatasetImport(c, request.DatasetId, fileUploadId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dataset_action_id": datasetActionId})

}

func GetDatasetImportPreview(c *gin.Context, svc datasetservice.DatasetService, fileUploadService fileimports.FileImportService) {
	_, _, organizationIds := apictx.GetAuthFromContext(c)

	if len(organizationIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	fileUploadId, err := uuid.Parse(c.Param("fileUploadId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file upload id"})
		return
	}

	preview, err := svc.GetFileUploadPreview(c, fileUploadId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data_preview": preview})
}

func ConfirmDatasetImport(c *gin.Context, svc datasetservice.DatasetService) {
	_, _, organizationIds := apictx.GetAuthFromContext(c)

	if len(organizationIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	fileUploadId, err := uuid.Parse(c.Param("fileUploadId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file upload id"})
		return
	}

	var request dtos.ConfirmDatasetImportRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err = svc.ImportDataFromFile(c, organizationIds[0], request.DatasetId, fileUploadId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func GetDatasetDisplayConfig(c *gin.Context, svc datasetservice.DatasetService) {
	ctx := c.MustGet("datasetContext").(middleware.DatasetContext)

	displayConfig, err := svc.GetDatasetDisplayConfig(c, ctx.MerchantID, ctx.DatasetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dtos.GetDatasetDisplayConfigResponse{
		DisplayConfig: displayConfig,
	}
	c.JSON(http.StatusOK, response)
}

func SetDatasetDisplayConfig(c *gin.Context, svc datasetservice.DatasetService) {
	ctx := c.MustGet("datasetContext").(middleware.DatasetContext)

	var request dtos.SetDatasetDisplayConfigRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Reuse the existing UpdateDataset method with only the display config
	actionId, err := svc.UpdateDataset(c, ctx.MerchantID, ctx.DatasetID, models.UpdateDatasetParams{
		DisplayConfig: &request.DisplayConfig,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"action_id": actionId})
}

func GetFileImportsHistory(c *gin.Context, svc datasetservice.DatasetService, fileUploadService fileimports.FileImportService) {
	ctx := c.MustGet("datasetContext").(middleware.DatasetContext)

	datasetId, err := uuid.Parse(ctx.DatasetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dataset id"})
		return
	}

	fileUploads, err := svc.GetDatasetFileUploads(c, datasetId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(fileUploads) == 0 {
		fileUploads = []models.DatasetFileUpload{}
	}

	c.JSON(http.StatusOK, gin.H{"file_uploads": fileUploads})

}

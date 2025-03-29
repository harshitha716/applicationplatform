package dataplatformwebhooks

import (
	"net/http"

	dataplatformservice "github.com/Zampfi/application-platform/services/api/core/dataplatform"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform/models"
	datasetservice "github.com/Zampfi/application-platform/services/api/core/datasets/service"
	apictx "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func HandleJobStatusUpdate(c *gin.Context, dataplatformService dataplatformservice.DataPlatformService, datasetService datasetservice.DatasetService) {
	logger := apictx.GetLoggerFromCtx(c)

	var jobStatusUpdate models.DatabricksJobStatusUpdatePayload
	logger.Info("JOB STATUS UPDATE RECEIVED")

	err := c.ShouldBindJSON(&jobStatusUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Info("JOB STATUS UPDATE BODY:", zap.Any("jobStatusUpdate", jobStatusUpdate))

	action, err := dataplatformService.UpdateAction(c.Request.Context(), jobStatusUpdate)
	if err != nil {
		logger.Error("ERROR UPDATING ACTION", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_id, err := uuid.Parse(action.ActorId)
	if err != nil {
		logger.Error("ERROR PARSING USER ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctxWithUserId := apictx.AddAuthToContext(c.Request.Context(), "user", user_id, []uuid.UUID{})

	err = datasetService.UpdateDatasetActionStatus(ctxWithUserId, action.ID, string(action.ActionStatus))
	if err != nil {
		logger.Error("ERROR CONSUMING ACTION IN DATASET SERVICE", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Info("ACTION UPDATED", zap.Any("action", action))
	// TODO: APPLICATION PLATFORM TO HANDLE DATASET CREATION
	c.JSON(http.StatusOK, gin.H{"message": "Job status updated"})
}

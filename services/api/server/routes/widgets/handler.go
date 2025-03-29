package widgets

import (
	"encoding/json"
	"net/http"

	widgetservice "github.com/Zampfi/application-platform/services/api/core/widgets/service"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/Zampfi/application-platform/services/api/server/routes/widgets/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GetWidgetInstanceData(c *gin.Context, widgetService widgetservice.WidgetsService) {
	_, _, orgIds := apicontext.GetAuthFromContext(c)
	if len(orgIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no organization ids found"})
		return
	}

	orgId := orgIds[0]

	widgetInstanceId, err := uuid.Parse(c.Param("widgetInstanceId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid widget instance ID"})
		return
	}

	// Get individual query parameters
	filtersStr := c.Query("filters")
	timeColumnsStr := c.Query("time_columns")
	periodicity := c.Query("periodicity")
	currency := c.Query("currency")

	queryParams := dtos.WidgetQueryParams{
		Filters: []dtos.WidgetFilters{},
	}

	if filtersStr != "" {
		if err := json.Unmarshal([]byte(filtersStr), &queryParams.Filters); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filters format"})
			return
		}
	}

	if timeColumnsStr != "" {
		if err := json.Unmarshal([]byte(timeColumnsStr), &queryParams.TimeColumns); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time columns format"})
			return
		}
	}

	if periodicity != "" {
		queryParams.Periodicity = &periodicity
	}

	if currency != "" {
		queryParams.Currency = &currency
	}

	queryParamModels := queryParams.ToModels()
	widgetInstanceData, err := widgetService.GetWidgetInstanceData(c, orgId, widgetInstanceId, queryParamModels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get widget instance data"})
		return
	}

	periodicityValue := "daily"
	if queryParams.Periodicity != nil {
		periodicityValue = *queryParams.Periodicity
	}

	resp := dtos.NewWidgetInstanceDataResponse(widgetInstanceData, periodicityValue, queryParams.Currency)

	c.JSON(http.StatusOK, resp)
}

func GetWidgetInstance(c *gin.Context, widgetService widgetservice.WidgetsService) {
	ctxLogger := apicontext.GetLoggerFromCtx(c)

	widgetInstanceId, err := uuid.Parse(c.Param("widgetInstanceId"))
	if err != nil {
		ctxLogger.Error("failed to parse widget instance ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid widget instance ID"})
		return
	}

	widgetInstance, err := widgetService.GetWidgetInstance(c, widgetInstanceId)
	if err != nil {
		ctxLogger.Error("failed to get widget instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get widget instance"})
		return
	}

	resp, err := dtos.NewWidgetInstanceResponse(&widgetInstance)
	if err != nil {
		ctxLogger.Error("failed to create widget instance response", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create widget instance response"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

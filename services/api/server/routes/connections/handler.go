package connections

import (
	"errors"
	"net/http"
	"time"

	"github.com/Zampfi/application-platform/services/api/core/connections/service"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	"github.com/Zampfi/application-platform/services/api/server/routes/connections/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go.uber.org/zap"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"

	manager "github.com/Zampfi/application-platform/services/api/core/connections/managers"
	organization "github.com/Zampfi/application-platform/services/api/core/organizations"
	schedules "github.com/Zampfi/application-platform/services/api/core/schedules/service"
)

func HandleCreateConnection(c *gin.Context, connectionService service.ConnectionService, txStore store.Store, connectionManagerRegistry manager.ConnectionManagerRegistry, scheduleService schedules.ScheduleService, organizationService organization.OrganizationService) {
	ctxLogger := apicontext.GetLoggerFromCtx(c)
	var params dto.CreateConnectionRequest

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	connectionManager, err := connectionManagerRegistry.GetManager(params.ConnectorName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	connectorID, err := uuid.Parse(params.ConnectorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid connector ID"})
		return
	}

	createConnectionParams := dbmodels.CreateConnectionParams{
		ConnectorID: connectorID,
		Name:        params.DisplayName,
		Status:      "active",
		Config:      params.ConnectionConfig,
	}

	txStore.WithTx(c, func(store store.Store) error {

		connectionId, err := connectionService.CreateConnection(c, createConnectionParams, store)

		if err != nil {
			ctxLogger.Error("failed to create connection", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return err
		}

		organizations, err := organizationService.GetOrganizations(c)

		if err != nil {
			ctxLogger.Error("failed to get organization", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return err
		}

		if len(organizations) != 1 {
			ctxLogger.Error("organization access forbidden")
			c.JSON(http.StatusForbidden, gin.H{"error": "organization access forbidden"})
			return errors.New("organization access forbidden")
		}

		createScheduleParams, err := connectionManager.ExtractAndValidateDefaultSchedules(connectionId, organizations[0].ID, createConnectionParams)

		if err != nil {
			ctxLogger.Error("failed to extract and validate default schedules", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return err
		}

		err = scheduleService.CreateSchedules(c, createScheduleParams, store)

		if err != nil {
			ctxLogger.Error("failed to create schedules", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return err
		}

		executeWorkflowParams, err := connectionManager.ExtractAndValidateTemporalParams(createScheduleParams)

		if err != nil {
			ctxLogger.Error("failed to extract and validate temporal params", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return err
		}

		err = scheduleService.CreateSchedulesInTemporal(c, executeWorkflowParams)

		if err != nil {
			ctxLogger.Error("failed to execute scheduled workflow in temporal", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return err
		}

		c.JSON(http.StatusCreated, gin.H{"connection_id": connectionId})

		return nil
	})

}

func handleGetConnections(c *gin.Context, connectionService service.ConnectionService, scheduleService schedules.ScheduleService, txStore store.Store) {
	ctxLogger := apicontext.GetLoggerFromCtx(c)
	connections, err := connectionService.GetConnections(c, txStore)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	connectionDetails := []dto.ConnectionDetails{}
	for _, connection := range connections {
		lastSyncedAt, err := scheduleService.GetLastSyncedAt(c, connection.ID)

		connectionDetails = append(connectionDetails, dto.ConnectionDetails{
			ID:               connection.ID.String(),
			Name:             connection.Name,
			ConnectorIconURL: connection.Connector.LogoURL,
			LastSyncedAt:     lastSyncedAt.UTC().Format(time.RFC3339),
			CreatedAt:        connection.CreatedAt.UTC().Format(time.RFC3339),
		})

		if err != nil {
			ctxLogger.Error("failed to get last synced at", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, connectionDetails)

}

func handleGetSchedules(c *gin.Context, connectionService service.ConnectionService, scheduleService schedules.ScheduleService) {
	var params dto.GetSchedulesRequest

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	connectionId, err := uuid.Parse(params.ConnectionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid connection ID"})
		return
	}

	connection, err := connectionService.GetConnectionByID(c, connectionId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	schedules, err := scheduleService.GetSchedules(c, connectionId)

	scheduleDetails := []dto.ScheduleDetails{}

	for _, schedule := range schedules {
		temporalSchedule, err := scheduleService.GetScheduleDetailsFromTemporal(c, schedule.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		scheduleDetails = append(scheduleDetails, dto.ScheduleDetails{
			ID:        schedule.ID.String(),
			Name:      schedule.Name,
			Status:    schedule.Status,
			NextRunAt: temporalSchedule.Info.NextActionTimes[0].Format(time.RFC3339),
			LastRunAt: temporalSchedule.Info.RecentActions[len(temporalSchedule.Info.RecentActions)-1].ScheduleTime.Format(time.RFC3339),
			CreatedAt: temporalSchedule.Info.CreatedAt.Format(time.RFC3339),
			UpdatedAt: temporalSchedule.Info.LastUpdateAt.Format(time.RFC3339),
			LogoURL:   connection.Connector.LogoURL,
		})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scheduleDetails)
}

package middleware

import (
	"net/http"

	"github.com/Zampfi/application-platform/services/api/db/models"
	dbmodel "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DatasetContext struct {
	DatasetID  string
	MerchantID uuid.UUID
	UserID     *uuid.UUID
}

func ValidateDatasetAccess(datasetStore store.DatasetStore) gin.HandlerFunc {
	return func(c *gin.Context) {

		logger := apicontext.GetLoggerFromCtx(c)

		datasetId := c.Param("datasetId")
		_, userId, orgIds := apicontext.GetAuthFromContext(c)

		if userId == nil {
			logger.Error("user ID not found", zap.String("error", "user ID not found"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID not found"})
			c.Abort()
			return
		}

		if len(orgIds) == 0 {
			logger.Error("organization ID not found", zap.String("error", "organization ID not found"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "organization ID not found"})
			c.Abort()
			return
		}

		// Validating dataset access for user
		_, err := datasetStore.GetDatasetById(c, datasetId)
		if err != nil {
			logger.Error("failed to authorize dataset access", zap.String("error", err.Error()))
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized dataset access"})
			c.Abort()
			return
		}

		merchantId := orgIds[0]

		// TMP:if env is local, set merchantId to d815eb60-9258-41c4-a2bd-e8b7ae4a05ae
		// if os.Getenv("ENVIRONMENT") == "local" {
		// 	merchantId = uuid.MustParse("0d0e3757-15e6-470d-be7d-bfcce90eeb03")
		// }

		c.Set("datasetContext", DatasetContext{
			DatasetID:  datasetId,
			MerchantID: merchantId,
			UserID:     userId,
		})

		c.Next()
	}
}

func ValidateDatasetAdminAccess(flattenedResourceAudiencePoliciesStore store.FlattenedResourceAudiencePoliciesStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := apicontext.GetLoggerFromCtx(c)

		datasetId := c.Param("datasetId")
		_, userId, orgIds := apicontext.GetAuthFromContext(c)

		if userId == nil {
			logger.Error("user ID not found", zap.String("error", "user ID not found"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID not found"})
			c.Abort()
			return
		}

		if len(orgIds) == 0 {
			logger.Error("organization ID not found", zap.String("error", "organization ID not found"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "organization ID not found"})
			c.Abort()
			return
		}

		datasetUUID, err := uuid.Parse(datasetId)
		if err != nil {
			logger.Error("invalid dataset ID", zap.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dataset ID"})
			c.Abort()
			return
		}

		flattenedResourceAudiencePolicies, err := flattenedResourceAudiencePoliciesStore.GetFlattenedResourceAudiencePolicies(c, dbmodel.FlattenedResourceAudiencePoliciesFilters{
			ResourceIds:   []uuid.UUID{datasetUUID},
			UserIds:       []uuid.UUID{*userId},
			ResourceTypes: []string{"dataset"},
			Privileges:    []models.ResourcePrivilege{models.PrivilegeDatasetAdmin},
		})
		if err != nil {
			logger.Error("failed to get flattened resource audience policies", zap.String("user_id", userId.String()), zap.String("dataset_id", datasetId), zap.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get flattened resource audience policies"})
			c.Abort()
			return
		}

		if len(flattenedResourceAudiencePolicies) == 0 {
			logger.Error("user does not have admin access to dataset", zap.String("user_id", userId.String()), zap.String("dataset_id", datasetId), zap.String("error", "user does not have admin access to dataset"))
			c.JSON(http.StatusForbidden, gin.H{"error": "user does not have admin access to dataset"})
			c.Abort()
			return
		}

		merchantId := orgIds[0]

		// TMP:if env is local, set merchantId to d815eb60-9258-41c4-a2bd-e8b7ae4a05ae
		// if os.Getenv("ENVIRONMENT") == constants.ENVLOCAL {
		// 	merchantId = uuid.MustParse("0d0e3757-15e6-470d-be7d-bfcce90eeb03")
		// }

		c.Set("datasetContext", DatasetContext{
			DatasetID:  datasetId,
			MerchantID: merchantId,
			UserID:     userId,
		})

		c.Next()

	}
}

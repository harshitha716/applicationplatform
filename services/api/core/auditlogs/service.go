package auditlogs

import (
	"context"
	"encoding/json"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/core/auditlogs/models"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuditLogService interface {
	WithResource(resourceType dbmodels.ResourceType) AuditLogServiceWithResource
	GetResourceAuditLogMiddleware(resource dbmodels.ResourceType, resourceIdRouteParam string) gin.HandlerFunc
	GetAuditLogsByOrganizationId(ctx context.Context, organizationId uuid.UUID, kind dbmodels.AuditLogKind) ([]models.AuditLog, error)
	EmitAuditLog(ctx context.Context, params models.AuditLogEmitParams) error
}

type AuditLogServiceWithResource interface {
	EmitAuditLog(ctx context.Context, resourceId uuid.UUID, auditLogKind dbmodels.AuditLogKind, eventName string, payload map[string]interface{}) error
}

type auditLogService struct {
	store store.Store
}

type auditLogServiceWithResource struct {
	service      *auditLogService
	resourceType dbmodels.ResourceType
}

func NewAuditLogService(serverConfig *serverconfig.ServerConfig) AuditLogService {
	return &auditLogService{
		store: serverConfig.Store,
	}
}

func (s *auditLogService) WithResource(resourceType dbmodels.ResourceType) AuditLogServiceWithResource {
	return &auditLogServiceWithResource{
		service:      s,
		resourceType: resourceType,
	}
}

func (s *auditLogService) GetAuditLogsByOrganizationId(ctx context.Context, organizationId uuid.UUID, kind dbmodels.AuditLogKind) ([]models.AuditLog, error) {
	logger := apicontext.GetLoggerFromCtx(ctx)

	dbAuditLogs, err := s.store.GetAuditLogsByOrganizationId(ctx, organizationId, kind)
	if err != nil {
		logger.Error("failed to get audit logs by organization id", zap.Error(err))
		return nil, err
	}

	var auditLogs []models.AuditLog
	for _, dbAuditLog := range dbAuditLogs {
		var auditLog models.AuditLog
		auditLog.FromSchema(dbAuditLog)
		auditLogs = append(auditLogs, auditLog)
	}

	return auditLogs, nil
}

func (s *auditLogService) EmitAuditLog(ctx context.Context, params models.AuditLogEmitParams) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	// Create payload JSON
	payloadBytes, err := json.Marshal(params.Payload)
	if err != nil {
		logger.Error("failed to marshal audit log payload", zap.Error(err))
		return err
	}

	// Get email, IP address, and user agent from context
	contextEmail, contextIPAddress, contextUserAgent := apicontext.GetAuditInfoFromContext(ctx)
	
	// Use context values if params values are empty
	if params.UserEmail == "" {
		params.UserEmail = contextEmail
	}
	
	if params.IPAddress == "" {
		params.IPAddress = contextIPAddress
	}
	
	if params.UserAgent == "" {
		params.UserAgent = contextUserAgent
	}

	// Create audit log
	auditLog := dbmodels.AuditLog{
		Kind:           params.Kind,
		OrganizationID: params.OrganizationID,
		IPAddress:      params.IPAddress,
		UserEmail:      params.UserEmail,
		UserAgent:      params.UserAgent,
		ResourceType:   dbmodels.ResourceType(params.ResourceType),
		ResourceID:     params.ResourceID,
		EventName:      params.EventName,
		Payload:        payloadBytes,
	}

	// Save audit log
	_, err = s.store.CreateAuditLog(ctx, auditLog)
	if err != nil {
		logger.Error("failed to create audit log", zap.Error(err))
		return err
	}

	return nil
}

func (s *auditLogService) GetResourceAuditLogMiddleware(resource dbmodels.ResourceType, resourceIdRouteParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This will be implemented in Phase 2
		c.Next()
	}
}

func (s *auditLogServiceWithResource) EmitAuditLog(ctx context.Context, resourceId uuid.UUID, auditLogKind dbmodels.AuditLogKind, eventName string, payload map[string]interface{}) error {
	logger := apicontext.GetLoggerFromCtx(ctx)

	// Get user info from context
	_, userId, orgIds := apicontext.GetAuthFromContext(ctx)
	if userId == nil || len(orgIds) == 0 {
		logger.Error("failed to get user or organization from context")
		return nil // Don't fail the request if audit logging fails
	}

	// Get email, IP address, and user agent from context
	contextEmail, contextIPAddress, contextUserAgent := apicontext.GetAuditInfoFromContext(ctx)

	// If email is not in context, try to get it from user record
	userEmail := contextEmail
	if userEmail == "" {
		// Get user information
		user, err := s.service.store.GetUserById(ctx, userId.String())
		if err != nil {
			logger.Error("failed to get user", zap.Error(err))
			return nil // Don't fail the request if audit logging fails
		}
		userEmail = user.Email
	}

	params := models.AuditLogEmitParams{
		Kind:           auditLogKind,
		ResourceType:   s.resourceType,
		ResourceID:     resourceId,
		EventName:      eventName,
		Payload:        payload,
		UserEmail:      userEmail,
		IPAddress:      contextIPAddress,
		UserAgent:      contextUserAgent,
		OrganizationID: orgIds[0], // Use the first organization from context
	}

	return s.service.EmitAuditLog(ctx, params)
}

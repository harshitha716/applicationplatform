package apicontext

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
)

const (
	contextKeyUserID            string = "user_id"
	contextKeyUserOrganizations string = "user_organizations"
	contextKeyUserRole          string = "user_role"
	contextKeyUserEmail         string = "user_email"
	contextKeyUserIPAddress     string = "user_ip_address"
	contextKeyUserAgent         string = "user_agent"
)

func AddAuthToContext(ctx context.Context, role string, userID uuid.UUID, userOrganizations []uuid.UUID) context.Context {
	enrichedCtx := AddCtxVariableToCtx(ctx, contextKeyUserID, userID)
	enrichedCtx = AddCtxVariableToCtx(enrichedCtx, contextKeyUserOrganizations, userOrganizations)
	enrichedCtx = AddCtxVariableToCtx(enrichedCtx, contextKeyUserRole, role)
	return enrichedCtx
}

func AddAuditInfoToContext(ctx context.Context, email string, ipAddress string, userAgent string) context.Context {
	enrichedCtx := AddCtxVariableToCtx(ctx, contextKeyUserEmail, email)
	enrichedCtx = AddCtxVariableToCtx(enrichedCtx, contextKeyUserIPAddress, ipAddress)
	enrichedCtx = AddCtxVariableToCtx(enrichedCtx, contextKeyUserAgent, userAgent)
	return enrichedCtx
}

func GetAuditInfoFromContext(ctx context.Context) (string, string, string) {
	emailRaw := getCtxVariableFromCtx(ctx, contextKeyUserEmail)
	ipAddressRaw := getCtxVariableFromCtx(ctx, contextKeyUserIPAddress)
	userAgentRaw := getCtxVariableFromCtx(ctx, contextKeyUserAgent)
	
	email := ""
	if emailRaw != nil {
		if emailStr, ok := emailRaw.(string); ok {
			email = emailStr
		}
	}
	
	ipAddress := ""
	if ipAddressRaw != nil {
		if ipAddressStr, ok := ipAddressRaw.(string); ok {
			ipAddress = ipAddressStr
		}
	}
	
	userAgent := ""
	if userAgentRaw != nil {
		if userAgentStr, ok := userAgentRaw.(string); ok {
			userAgent = userAgentStr
		}
	}
	
	return email, ipAddress, userAgent
}

func GetAuthFromContext(ctx context.Context) (string, *uuid.UUID, []uuid.UUID) {

	userRoleRaw := getCtxVariableFromCtx(ctx, contextKeyUserRole)

	userIDRaw := getCtxVariableFromCtx(ctx, contextKeyUserID)

	userOrganizationsRaw := getCtxVariableFromCtx(ctx, contextKeyUserOrganizations)

	return parseContextAuth(userRoleRaw, userIDRaw, userOrganizationsRaw)

}

func parseContextAuth(userRoleRaw interface{}, userIDRaw interface{}, userOrganizationsRaw interface{}) (string, *uuid.UUID, []uuid.UUID) {
	userRole := "anonymous"

	if userRoleRaw == nil {
		return "anonymous", nil, []uuid.UUID{}
	}

	userRole, ok := userRoleRaw.(string)
	if !ok {
		userRole = "anonymous"
	}

	if userIDRaw == nil {
		return userRole, nil, []uuid.UUID{}
	}

	userId, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return userRole, nil, []uuid.UUID{}
	}

	if userOrganizationsRaw == nil {
		return userRole, &userId, []uuid.UUID{}
	}

	userOrganizations, ok := userOrganizationsRaw.([]uuid.UUID)
	if !ok {
		return userRole, &userId, []uuid.UUID{}
	}

	return userRole, &userId, userOrganizations
}

func AddAuthToGinContext(ctx *gin.Context, role string, userID uuid.UUID, userOrganizations []uuid.UUID) {
	AddContextVariableToGinContext(ctx, contextKeyUserRole, role)
	AddContextVariableToGinContext(ctx, contextKeyUserID, userID)
	AddContextVariableToGinContext(ctx, contextKeyUserOrganizations, userOrganizations)
}

func AddAuditInfoToGinContext(ctx *gin.Context, email string, ipAddress string, userAgent string) {
	AddContextVariableToGinContext(ctx, contextKeyUserEmail, email)
	AddContextVariableToGinContext(ctx, contextKeyUserIPAddress, ipAddress)
	AddContextVariableToGinContext(ctx, contextKeyUserAgent, userAgent)
}

func IsZampEmail(email string) bool {
	return strings.HasSuffix(email, "@zamp.finance") || strings.HasSuffix(email, "@zamp.ai")
}

func AddAuthToWorkflowContext(ctx workflow.Context, role string, userID uuid.UUID, userOrganizations []uuid.UUID) workflow.Context {
	ctx = AddContextVariableToWorkflowContext(ctx, contextKeyUserRole, role)
	ctx = AddContextVariableToWorkflowContext(ctx, contextKeyUserID, userID)
	ctx = AddContextVariableToWorkflowContext(ctx, contextKeyUserOrganizations, userOrganizations)
	return ctx
}

func GetAuthFromWorkflowContext(ctx workflow.Context) (string, *uuid.UUID, []uuid.UUID) {
	role := GetContextVariableFromWorkflowContext(ctx, contextKeyUserRole)
	userID := GetContextVariableFromWorkflowContext(ctx, contextKeyUserID)
	userOrganizations := GetContextVariableFromWorkflowContext(ctx, contextKeyUserOrganizations)

	return parseContextAuth(role, userID, userOrganizations)
}

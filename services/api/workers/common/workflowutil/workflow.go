package workflowutil

import (
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	temporalWorkflow "go.temporal.io/sdk/workflow"
)

type ParamsBase interface {
	GetAccessControlParams() (userId uuid.UUID, orgIds []uuid.UUID)
}

func AddAccessControlParamsToWorkflowCtx(ctx temporalWorkflow.Context, params ParamsBase) temporalWorkflow.Context {
	userID, userOrganizations := params.GetAccessControlParams()

	ctx = apicontext.AddAuthToWorkflowContext(ctx, "user", userID, userOrganizations)

	return ctx
}

package helper

import (
	"net/http"
	"testing"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddAuthHeaders(t *testing.T) {
	// AddAuthHeaders(headers http.Header, ginctx *gin.Context)
	headers := http.Header{}

	ctx, _ := gin.CreateTestContext(nil)

	AddAuthHeaders(headers, ctx)

	assert.Equal(t, headers.Get(PROXY_USER_ID_HEADER), "")
	assert.Equal(t, headers.Get(PROXY_WORKSPACE_IDS_HEADER), "")

	userId := uuid.New()
	organizationIds := []uuid.UUID{uuid.New(), uuid.New()}

	apicontext.AddAuthToGinContext(ctx, "user", userId, organizationIds)

	AddAuthHeaders(headers, ctx)

	assert.Equal(t, headers.Get(PROXY_USER_ID_HEADER), userId.String())

	assert.Equal(t, headers.Values(PROXY_WORKSPACE_IDS_HEADER), []string{organizationIds[0].String(), organizationIds[1].String()})

}

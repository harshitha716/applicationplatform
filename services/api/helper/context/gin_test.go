package apicontext

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAddContextVariablesToGinCtx(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	value := "value"

	AddContextVariableToGinContext(ctx, contextKeyUserID, value)

	assert.NotEmpty(t, ctx.Value(contextKeyCtxVariables))
	assert.Equal(t, value, ctx.Value(contextKeyCtxVariables).(map[string]interface{})[contextKeyUserID])
}

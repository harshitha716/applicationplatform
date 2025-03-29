package tests

import (
	"testing"

	"github.com/Zampfi/application-platform/services/api/core/dataplatform/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGenerateUUIDWithUnderscores(t *testing.T) {
	result := helpers.GenerateUUIDWithUnderscores()
	assert.NotContains(t, result, "-")

}

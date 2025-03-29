package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnector_TableName(t *testing.T) {
	t.Parallel()
	connector := Connector{}
	assert.Equal(t, "connectors", connector.TableName())
}

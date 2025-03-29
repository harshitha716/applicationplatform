package provider

import (
	"context"

	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/models"
)

type ProviderService interface {
	Query(ctx context.Context, table string, query string, args ...interface{}) (models.QueryResult, error)
}

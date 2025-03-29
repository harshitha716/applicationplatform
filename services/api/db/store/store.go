package store

import (
	"context"

	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	"gorm.io/gorm"
)

type TransactionStore interface {
	WithTx(ctx context.Context, fn func(Store) error) error
}

type Store interface {
	UserStore
	OrganizationStore
	DatasetStore
	PageStore
	SheetStore
	TransactionStore
	WidgetStore
	DatasetActionStore
	RuleStore
	DatasetActionStore
	DatasetFileUploadStore
	RuleStore
	FlattenedResourceAudiencePoliciesStore
	ConnectorStore
	ConnectionStore
	ConnectionPoliciesStore
	ScheduleStore
	TeamStore
	FileUploadStore
	AuditLogStore
	PaymentsConfigStore
}

type appStore struct {
	client *pgclient.PostgresClient
}

func NewStore(client *pgclient.PostgresClient) (Store, func()) {

	cleanup := func() {
		sqlDb, err := client.DB.DB()
		if err != nil {
			panic(err)
		}
		sqlDb.Close()
	}

	return &appStore{client: client}, cleanup
}
func (s *appStore) WithTx(ctx context.Context, fn func(Store) error) error {
	return s.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txClient := &pgclient.PostgresClient{DB: tx}
		txStore, _ := NewStore(txClient)
		// TODO: this is a connection leak; figure out how to close db connection
		return fn(txStore)
	})
}

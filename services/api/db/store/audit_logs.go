package store

import (
	"context"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLogStore interface {
	CreateAuditLog(ctx context.Context, auditLog models.AuditLog) (*models.AuditLog, error)
	GetAuditLogsByOrganizationId(ctx context.Context, organizationId uuid.UUID, kind models.AuditLogKind) ([]models.AuditLog, error)
	WithAuditLogTransaction(ctx context.Context, fn func(AuditLogStore) error) error
}

func (s *appStore) CreateAuditLog(ctx context.Context, auditLog models.AuditLog) (*models.AuditLog, error) {
	err := s.client.WithContext(ctx).Create(&auditLog).Error
	if err != nil {
		return nil, err
	}
	return &auditLog, nil
}

func (s *appStore) GetAuditLogsByOrganizationId(ctx context.Context, organizationId uuid.UUID, kind models.AuditLogKind) ([]models.AuditLog, error) {
	var auditLogs []models.AuditLog
	query := s.client.WithContext(ctx).Model(&models.AuditLog{})

	if kind != "" {
		query = query.Where("kind = ?", kind)
	}

	err := query.Where("organization_id = ?", organizationId).
		Order("created_at DESC").
		Find(&auditLogs).Error

	if err != nil {
		return nil, err
	}

	return auditLogs, nil
}

func (s *appStore) WithAuditLogTransaction(ctx context.Context, fn func(AuditLogStore) error) error {
	return s.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txStore := &appStore{client: &pgclient.PostgresClient{DB: tx}}
		return fn(txStore)
	})
}

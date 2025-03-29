package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Zampfi/application-platform/services/api/db/models"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
)

type PaymentsConfigStore interface {
	PaymentsConfigReadStore
	PaymentsConfigWriteStore
}

type PaymentsConfigReadStore interface {
	GetPaymentsConfigsByOrganizationId(ctx context.Context, organizationId string) (models.PaymentsConfig, error)
}

type PaymentsConfigWriteStore interface {
	CreatePaymentsConfig(ctx context.Context, paymentsConfig models.PaymentsConfig) (models.PaymentsConfig, error)
	UpdatePaymentsConfig(ctx context.Context, paymentsConfigId uuid.UUID, config json.RawMessage) (models.PaymentsConfig, error)
	UpdatePaymentsConfigStatus(ctx context.Context, paymentsConfigId uuid.UUID, status models.PaymentsConfigStatus) (models.PaymentsConfig, error)
	DeletePaymentsConfigById(ctx context.Context, paymentsConfigId uuid.UUID) error
}

func (s *appStore) GetPaymentsConfigsByOrganizationId(ctx context.Context, organizationId string) (models.PaymentsConfig, error) {
	var paymentsConfig models.PaymentsConfig
	err := s.client.WithContext(ctx).Model(&paymentsConfig).Where("organization_id = ?", organizationId).First(&paymentsConfig).Error
	if err != nil {
		return models.PaymentsConfig{}, err
	}
	return paymentsConfig, nil
}

func (s *appStore) CreatePaymentsConfig(ctx context.Context, paymentsConfig models.PaymentsConfig) (models.PaymentsConfig, error) {
	_, userId, organizationIds := apicontext.GetAuthFromContext(ctx)
	if userId == nil {
		return models.PaymentsConfig{}, fmt.Errorf("no user ID found in context")
	}

	if len(organizationIds) != 1 {
		return models.PaymentsConfig{}, fmt.Errorf("organization access forbidden")
	}

	paymentsConfig.OrganizationID = organizationIds[0]

	err := s.client.WithContext(ctx).Create(&paymentsConfig).Error
	if err != nil {
		return models.PaymentsConfig{}, err
	}

	return paymentsConfig, nil
}

func (s *appStore) UpdatePaymentsConfig(ctx context.Context, paymentsConfigId uuid.UUID, config json.RawMessage) (models.PaymentsConfig, error) {
	_, userId, _ := apicontext.GetAuthFromContext(ctx)
	if userId == nil {
		return models.PaymentsConfig{}, fmt.Errorf("no user ID found in context")
	}

	paymentsConfig := models.PaymentsConfig{
		ID:            paymentsConfigId,
		MappingConfig: config,
	}

	err := s.client.WithContext(ctx).Model(&paymentsConfig).Where("id = ?", paymentsConfigId).Update("mapping_config", config).Error
	if err != nil {
		return models.PaymentsConfig{}, err
	}

	return paymentsConfig, nil
}

func (s *appStore) UpdatePaymentsConfigStatus(ctx context.Context, paymentsConfigId uuid.UUID, status models.PaymentsConfigStatus) (models.PaymentsConfig, error) {
	_, userId, _ := apicontext.GetAuthFromContext(ctx)
	if userId == nil {
		return models.PaymentsConfig{}, fmt.Errorf("no user ID found in context")
	}

	paymentsConfig := models.PaymentsConfig{
		ID:     paymentsConfigId,
		Status: status,
	}

	err := s.client.WithContext(ctx).Model(&paymentsConfig).Where("id = ?", paymentsConfigId).Update("status", status).Error
	if err != nil {
		return models.PaymentsConfig{}, err
	}

	return paymentsConfig, nil
}

func (s *appStore) DeletePaymentsConfigById(ctx context.Context, paymentsConfigId uuid.UUID) error {
	_, userId, _ := apicontext.GetAuthFromContext(ctx)
	if userId == nil {
		return fmt.Errorf("no user ID found in context")
	}

	err := s.client.WithContext(ctx).Where("id = ?", paymentsConfigId).Delete(&models.PaymentsConfig{}).Error
	if err != nil {
		return err
	}

	return nil
}

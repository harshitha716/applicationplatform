package mock_store

import (
	"context"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	"github.com/google/uuid"
)

// CreateAuditLog provides a mock function with given fields: ctx, auditLog
func (_m *MockStore) CreateAuditLog(ctx context.Context, auditLog models.AuditLog) (*models.AuditLog, error) {
	ret := _m.Called(ctx, auditLog)

	var r0 *models.AuditLog
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.AuditLog) (*models.AuditLog, error)); ok {
		return rf(ctx, auditLog)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.AuditLog) *models.AuditLog); ok {
		r0 = rf(ctx, auditLog)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.AuditLog)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.AuditLog) error); ok {
		r1 = rf(ctx, auditLog)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAuditLogsByOrganizationId provides a mock function with given fields: ctx, organizationId, kind
func (_m *MockStore) GetAuditLogsByOrganizationId(ctx context.Context, organizationId uuid.UUID, kind models.AuditLogKind) ([]models.AuditLog, error) {
	ret := _m.Called(ctx, organizationId, kind)

	var r0 []models.AuditLog
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.AuditLogKind) ([]models.AuditLog, error)); ok {
		return rf(ctx, organizationId, kind)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.AuditLogKind) []models.AuditLog); ok {
		r0 = rf(ctx, organizationId, kind)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.AuditLog)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, models.AuditLogKind) error); ok {
		r1 = rf(ctx, organizationId, kind)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAuditLogsByResource provides a mock function with given fields: ctx, resource, resourceId
func (_m *MockStore) GetAuditLogsByResource(ctx context.Context, resource string, resourceId string) ([]models.AuditLog, error) {
	ret := _m.Called(ctx, resource, resourceId)

	var r0 []models.AuditLog
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]models.AuditLog, error)); ok {
		return rf(ctx, resource, resourceId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []models.AuditLog); ok {
		r0 = rf(ctx, resource, resourceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.AuditLog)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, resource, resourceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithAuditLogTransaction provides a mock function with given fields: ctx, fn
func (_m *MockStore) WithAuditLogTransaction(ctx context.Context, fn func(store.AuditLogStore) error) error {
	ret := _m.Called(ctx, fn)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(store.AuditLogStore) error) error); ok {
		r0 = rf(ctx, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

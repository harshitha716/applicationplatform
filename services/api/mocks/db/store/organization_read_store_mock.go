// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_store

import (
	context "context"

	models "github.com/Zampfi/application-platform/services/api/db/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockOrganizationReadStore is an autogenerated mock type for the OrganizationReadStore type
type MockOrganizationReadStore struct {
	mock.Mock
}

type MockOrganizationReadStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockOrganizationReadStore) EXPECT() *MockOrganizationReadStore_Expecter {
	return &MockOrganizationReadStore_Expecter{mock: &_m.Mock}
}

// GetOrganizationById provides a mock function with given fields: ctx, organizationId
func (_m *MockOrganizationReadStore) GetOrganizationById(ctx context.Context, organizationId string) (*models.Organization, error) {
	ret := _m.Called(ctx, organizationId)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationById")
	}

	var r0 *models.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.Organization, error)); ok {
		return rf(ctx, organizationId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.Organization); ok {
		r0 = rf(ctx, organizationId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Organization)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, organizationId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationReadStore_GetOrganizationById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationById'
type MockOrganizationReadStore_GetOrganizationById_Call struct {
	*mock.Call
}

// GetOrganizationById is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId string
func (_e *MockOrganizationReadStore_Expecter) GetOrganizationById(ctx interface{}, organizationId interface{}) *MockOrganizationReadStore_GetOrganizationById_Call {
	return &MockOrganizationReadStore_GetOrganizationById_Call{Call: _e.mock.On("GetOrganizationById", ctx, organizationId)}
}

func (_c *MockOrganizationReadStore_GetOrganizationById_Call) Run(run func(ctx context.Context, organizationId string)) *MockOrganizationReadStore_GetOrganizationById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationById_Call) Return(_a0 *models.Organization, _a1 error) *MockOrganizationReadStore_GetOrganizationById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationById_Call) RunAndReturn(run func(context.Context, string) (*models.Organization, error)) *MockOrganizationReadStore_GetOrganizationById_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationInvitationById provides a mock function with given fields: ctx, invitationId
func (_m *MockOrganizationReadStore) GetOrganizationInvitationById(ctx context.Context, invitationId uuid.UUID) (*models.OrganizationInvitation, error) {
	ret := _m.Called(ctx, invitationId)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationInvitationById")
	}

	var r0 *models.OrganizationInvitation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*models.OrganizationInvitation, error)); ok {
		return rf(ctx, invitationId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *models.OrganizationInvitation); ok {
		r0 = rf(ctx, invitationId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OrganizationInvitation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, invitationId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationReadStore_GetOrganizationInvitationById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationInvitationById'
type MockOrganizationReadStore_GetOrganizationInvitationById_Call struct {
	*mock.Call
}

// GetOrganizationInvitationById is a helper method to define mock.On call
//   - ctx context.Context
//   - invitationId uuid.UUID
func (_e *MockOrganizationReadStore_Expecter) GetOrganizationInvitationById(ctx interface{}, invitationId interface{}) *MockOrganizationReadStore_GetOrganizationInvitationById_Call {
	return &MockOrganizationReadStore_GetOrganizationInvitationById_Call{Call: _e.mock.On("GetOrganizationInvitationById", ctx, invitationId)}
}

func (_c *MockOrganizationReadStore_GetOrganizationInvitationById_Call) Run(run func(ctx context.Context, invitationId uuid.UUID)) *MockOrganizationReadStore_GetOrganizationInvitationById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationInvitationById_Call) Return(_a0 *models.OrganizationInvitation, _a1 error) *MockOrganizationReadStore_GetOrganizationInvitationById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationInvitationById_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*models.OrganizationInvitation, error)) *MockOrganizationReadStore_GetOrganizationInvitationById_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationInvitationsAll provides a mock function with given fields: ctx
func (_m *MockOrganizationReadStore) GetOrganizationInvitationsAll(ctx context.Context) ([]models.OrganizationInvitation, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationInvitationsAll")
	}

	var r0 []models.OrganizationInvitation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.OrganizationInvitation, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.OrganizationInvitation); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.OrganizationInvitation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationReadStore_GetOrganizationInvitationsAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationInvitationsAll'
type MockOrganizationReadStore_GetOrganizationInvitationsAll_Call struct {
	*mock.Call
}

// GetOrganizationInvitationsAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockOrganizationReadStore_Expecter) GetOrganizationInvitationsAll(ctx interface{}) *MockOrganizationReadStore_GetOrganizationInvitationsAll_Call {
	return &MockOrganizationReadStore_GetOrganizationInvitationsAll_Call{Call: _e.mock.On("GetOrganizationInvitationsAll", ctx)}
}

func (_c *MockOrganizationReadStore_GetOrganizationInvitationsAll_Call) Run(run func(ctx context.Context)) *MockOrganizationReadStore_GetOrganizationInvitationsAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationInvitationsAll_Call) Return(_a0 []models.OrganizationInvitation, _a1 error) *MockOrganizationReadStore_GetOrganizationInvitationsAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationInvitationsAll_Call) RunAndReturn(run func(context.Context) ([]models.OrganizationInvitation, error)) *MockOrganizationReadStore_GetOrganizationInvitationsAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationInvitationsAndMembershipRequests provides a mock function with given fields: ctx, organizationId
func (_m *MockOrganizationReadStore) GetOrganizationInvitationsAndMembershipRequests(ctx context.Context, organizationId uuid.UUID) (*models.Organization, error) {
	ret := _m.Called(ctx, organizationId)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationInvitationsAndMembershipRequests")
	}

	var r0 *models.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*models.Organization, error)); ok {
		return rf(ctx, organizationId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *models.Organization); ok {
		r0 = rf(ctx, organizationId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Organization)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, organizationId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationReadStore_GetOrganizationInvitationsAndMembershipRequests_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationInvitationsAndMembershipRequests'
type MockOrganizationReadStore_GetOrganizationInvitationsAndMembershipRequests_Call struct {
	*mock.Call
}

// GetOrganizationInvitationsAndMembershipRequests is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
func (_e *MockOrganizationReadStore_Expecter) GetOrganizationInvitationsAndMembershipRequests(ctx interface{}, organizationId interface{}) *MockOrganizationReadStore_GetOrganizationInvitationsAndMembershipRequests_Call {
	return &MockOrganizationReadStore_GetOrganizationInvitationsAndMembershipRequests_Call{Call: _e.mock.On("GetOrganizationInvitationsAndMembershipRequests", ctx, organizationId)}
}

func (_c *MockOrganizationReadStore_GetOrganizationInvitationsAndMembershipRequests_Call) Run(run func(ctx context.Context, organizationId uuid.UUID)) *MockOrganizationReadStore_GetOrganizationInvitationsAndMembershipRequests_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationInvitationsAndMembershipRequests_Call) Return(_a0 *models.Organization, _a1 error) *MockOrganizationReadStore_GetOrganizationInvitationsAndMembershipRequests_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationInvitationsAndMembershipRequests_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*models.Organization, error)) *MockOrganizationReadStore_GetOrganizationInvitationsAndMembershipRequests_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationInvitationsByOrganizationId provides a mock function with given fields: ctx, organizationId
func (_m *MockOrganizationReadStore) GetOrganizationInvitationsByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.OrganizationInvitation, error) {
	ret := _m.Called(ctx, organizationId)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationInvitationsByOrganizationId")
	}

	var r0 []models.OrganizationInvitation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]models.OrganizationInvitation, error)); ok {
		return rf(ctx, organizationId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []models.OrganizationInvitation); ok {
		r0 = rf(ctx, organizationId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.OrganizationInvitation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, organizationId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationReadStore_GetOrganizationInvitationsByOrganizationId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationInvitationsByOrganizationId'
type MockOrganizationReadStore_GetOrganizationInvitationsByOrganizationId_Call struct {
	*mock.Call
}

// GetOrganizationInvitationsByOrganizationId is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
func (_e *MockOrganizationReadStore_Expecter) GetOrganizationInvitationsByOrganizationId(ctx interface{}, organizationId interface{}) *MockOrganizationReadStore_GetOrganizationInvitationsByOrganizationId_Call {
	return &MockOrganizationReadStore_GetOrganizationInvitationsByOrganizationId_Call{Call: _e.mock.On("GetOrganizationInvitationsByOrganizationId", ctx, organizationId)}
}

func (_c *MockOrganizationReadStore_GetOrganizationInvitationsByOrganizationId_Call) Run(run func(ctx context.Context, organizationId uuid.UUID)) *MockOrganizationReadStore_GetOrganizationInvitationsByOrganizationId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationInvitationsByOrganizationId_Call) Return(_a0 []models.OrganizationInvitation, _a1 error) *MockOrganizationReadStore_GetOrganizationInvitationsByOrganizationId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationInvitationsByOrganizationId_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]models.OrganizationInvitation, error)) *MockOrganizationReadStore_GetOrganizationInvitationsByOrganizationId_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationMembershipRequestsAll provides a mock function with given fields: ctx
func (_m *MockOrganizationReadStore) GetOrganizationMembershipRequestsAll(ctx context.Context) ([]models.OrganizationMembershipRequest, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationMembershipRequestsAll")
	}

	var r0 []models.OrganizationMembershipRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.OrganizationMembershipRequest, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.OrganizationMembershipRequest); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.OrganizationMembershipRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationReadStore_GetOrganizationMembershipRequestsAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationMembershipRequestsAll'
type MockOrganizationReadStore_GetOrganizationMembershipRequestsAll_Call struct {
	*mock.Call
}

// GetOrganizationMembershipRequestsAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockOrganizationReadStore_Expecter) GetOrganizationMembershipRequestsAll(ctx interface{}) *MockOrganizationReadStore_GetOrganizationMembershipRequestsAll_Call {
	return &MockOrganizationReadStore_GetOrganizationMembershipRequestsAll_Call{Call: _e.mock.On("GetOrganizationMembershipRequestsAll", ctx)}
}

func (_c *MockOrganizationReadStore_GetOrganizationMembershipRequestsAll_Call) Run(run func(ctx context.Context)) *MockOrganizationReadStore_GetOrganizationMembershipRequestsAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationMembershipRequestsAll_Call) Return(_a0 []models.OrganizationMembershipRequest, _a1 error) *MockOrganizationReadStore_GetOrganizationMembershipRequestsAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationMembershipRequestsAll_Call) RunAndReturn(run func(context.Context) ([]models.OrganizationMembershipRequest, error)) *MockOrganizationReadStore_GetOrganizationMembershipRequestsAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationMembershipRequestsByOrganizationId provides a mock function with given fields: ctx, organizationId
func (_m *MockOrganizationReadStore) GetOrganizationMembershipRequestsByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.OrganizationMembershipRequest, error) {
	ret := _m.Called(ctx, organizationId)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationMembershipRequestsByOrganizationId")
	}

	var r0 []models.OrganizationMembershipRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]models.OrganizationMembershipRequest, error)); ok {
		return rf(ctx, organizationId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []models.OrganizationMembershipRequest); ok {
		r0 = rf(ctx, organizationId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.OrganizationMembershipRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, organizationId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationReadStore_GetOrganizationMembershipRequestsByOrganizationId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationMembershipRequestsByOrganizationId'
type MockOrganizationReadStore_GetOrganizationMembershipRequestsByOrganizationId_Call struct {
	*mock.Call
}

// GetOrganizationMembershipRequestsByOrganizationId is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
func (_e *MockOrganizationReadStore_Expecter) GetOrganizationMembershipRequestsByOrganizationId(ctx interface{}, organizationId interface{}) *MockOrganizationReadStore_GetOrganizationMembershipRequestsByOrganizationId_Call {
	return &MockOrganizationReadStore_GetOrganizationMembershipRequestsByOrganizationId_Call{Call: _e.mock.On("GetOrganizationMembershipRequestsByOrganizationId", ctx, organizationId)}
}

func (_c *MockOrganizationReadStore_GetOrganizationMembershipRequestsByOrganizationId_Call) Run(run func(ctx context.Context, organizationId uuid.UUID)) *MockOrganizationReadStore_GetOrganizationMembershipRequestsByOrganizationId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationMembershipRequestsByOrganizationId_Call) Return(_a0 []models.OrganizationMembershipRequest, _a1 error) *MockOrganizationReadStore_GetOrganizationMembershipRequestsByOrganizationId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationMembershipRequestsByOrganizationId_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]models.OrganizationMembershipRequest, error)) *MockOrganizationReadStore_GetOrganizationMembershipRequestsByOrganizationId_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationPolicies provides a mock function with given fields: ctx, orgId
func (_m *MockOrganizationReadStore) GetOrganizationPolicies(ctx context.Context, orgId uuid.UUID) ([]models.ResourceAudiencePolicy, error) {
	ret := _m.Called(ctx, orgId)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationPolicies")
	}

	var r0 []models.ResourceAudiencePolicy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]models.ResourceAudiencePolicy, error)); ok {
		return rf(ctx, orgId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []models.ResourceAudiencePolicy); ok {
		r0 = rf(ctx, orgId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ResourceAudiencePolicy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, orgId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationReadStore_GetOrganizationPolicies_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationPolicies'
type MockOrganizationReadStore_GetOrganizationPolicies_Call struct {
	*mock.Call
}

// GetOrganizationPolicies is a helper method to define mock.On call
//   - ctx context.Context
//   - orgId uuid.UUID
func (_e *MockOrganizationReadStore_Expecter) GetOrganizationPolicies(ctx interface{}, orgId interface{}) *MockOrganizationReadStore_GetOrganizationPolicies_Call {
	return &MockOrganizationReadStore_GetOrganizationPolicies_Call{Call: _e.mock.On("GetOrganizationPolicies", ctx, orgId)}
}

func (_c *MockOrganizationReadStore_GetOrganizationPolicies_Call) Run(run func(ctx context.Context, orgId uuid.UUID)) *MockOrganizationReadStore_GetOrganizationPolicies_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationPolicies_Call) Return(_a0 []models.ResourceAudiencePolicy, _a1 error) *MockOrganizationReadStore_GetOrganizationPolicies_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationPolicies_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]models.ResourceAudiencePolicy, error)) *MockOrganizationReadStore_GetOrganizationPolicies_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationPoliciesByEmail provides a mock function with given fields: ctx, organizationId, email
func (_m *MockOrganizationReadStore) GetOrganizationPoliciesByEmail(ctx context.Context, organizationId uuid.UUID, email string) ([]models.ResourceAudiencePolicy, error) {
	ret := _m.Called(ctx, organizationId, email)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationPoliciesByEmail")
	}

	var r0 []models.ResourceAudiencePolicy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) ([]models.ResourceAudiencePolicy, error)); ok {
		return rf(ctx, organizationId, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) []models.ResourceAudiencePolicy); ok {
		r0 = rf(ctx, organizationId, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ResourceAudiencePolicy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, string) error); ok {
		r1 = rf(ctx, organizationId, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationReadStore_GetOrganizationPoliciesByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationPoliciesByEmail'
type MockOrganizationReadStore_GetOrganizationPoliciesByEmail_Call struct {
	*mock.Call
}

// GetOrganizationPoliciesByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
//   - email string
func (_e *MockOrganizationReadStore_Expecter) GetOrganizationPoliciesByEmail(ctx interface{}, organizationId interface{}, email interface{}) *MockOrganizationReadStore_GetOrganizationPoliciesByEmail_Call {
	return &MockOrganizationReadStore_GetOrganizationPoliciesByEmail_Call{Call: _e.mock.On("GetOrganizationPoliciesByEmail", ctx, organizationId, email)}
}

func (_c *MockOrganizationReadStore_GetOrganizationPoliciesByEmail_Call) Run(run func(ctx context.Context, organizationId uuid.UUID, email string)) *MockOrganizationReadStore_GetOrganizationPoliciesByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(string))
	})
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationPoliciesByEmail_Call) Return(_a0 []models.ResourceAudiencePolicy, _a1 error) *MockOrganizationReadStore_GetOrganizationPoliciesByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationPoliciesByEmail_Call) RunAndReturn(run func(context.Context, uuid.UUID, string) ([]models.ResourceAudiencePolicy, error)) *MockOrganizationReadStore_GetOrganizationPoliciesByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationPolicyByUser provides a mock function with given fields: ctx, organizationId, userId
func (_m *MockOrganizationReadStore) GetOrganizationPolicyByUser(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID) (*models.ResourceAudiencePolicy, error) {
	ret := _m.Called(ctx, organizationId, userId)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationPolicyByUser")
	}

	var r0 *models.ResourceAudiencePolicy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) (*models.ResourceAudiencePolicy, error)); ok {
		return rf(ctx, organizationId, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) *models.ResourceAudiencePolicy); ok {
		r0 = rf(ctx, organizationId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ResourceAudiencePolicy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, organizationId, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationReadStore_GetOrganizationPolicyByUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationPolicyByUser'
type MockOrganizationReadStore_GetOrganizationPolicyByUser_Call struct {
	*mock.Call
}

// GetOrganizationPolicyByUser is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
//   - userId uuid.UUID
func (_e *MockOrganizationReadStore_Expecter) GetOrganizationPolicyByUser(ctx interface{}, organizationId interface{}, userId interface{}) *MockOrganizationReadStore_GetOrganizationPolicyByUser_Call {
	return &MockOrganizationReadStore_GetOrganizationPolicyByUser_Call{Call: _e.mock.On("GetOrganizationPolicyByUser", ctx, organizationId, userId)}
}

func (_c *MockOrganizationReadStore_GetOrganizationPolicyByUser_Call) Run(run func(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID)) *MockOrganizationReadStore_GetOrganizationPolicyByUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationPolicyByUser_Call) Return(_a0 *models.ResourceAudiencePolicy, _a1 error) *MockOrganizationReadStore_GetOrganizationPolicyByUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationPolicyByUser_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID) (*models.ResourceAudiencePolicy, error)) *MockOrganizationReadStore_GetOrganizationPolicyByUser_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationSSOConfigsByOrganizationId provides a mock function with given fields: ctx, organizationId
func (_m *MockOrganizationReadStore) GetOrganizationSSOConfigsByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.OrganizationSSOConfig, error) {
	ret := _m.Called(ctx, organizationId)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationSSOConfigsByOrganizationId")
	}

	var r0 []models.OrganizationSSOConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]models.OrganizationSSOConfig, error)); ok {
		return rf(ctx, organizationId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []models.OrganizationSSOConfig); ok {
		r0 = rf(ctx, organizationId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.OrganizationSSOConfig)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, organizationId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationReadStore_GetOrganizationSSOConfigsByOrganizationId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationSSOConfigsByOrganizationId'
type MockOrganizationReadStore_GetOrganizationSSOConfigsByOrganizationId_Call struct {
	*mock.Call
}

// GetOrganizationSSOConfigsByOrganizationId is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
func (_e *MockOrganizationReadStore_Expecter) GetOrganizationSSOConfigsByOrganizationId(ctx interface{}, organizationId interface{}) *MockOrganizationReadStore_GetOrganizationSSOConfigsByOrganizationId_Call {
	return &MockOrganizationReadStore_GetOrganizationSSOConfigsByOrganizationId_Call{Call: _e.mock.On("GetOrganizationSSOConfigsByOrganizationId", ctx, organizationId)}
}

func (_c *MockOrganizationReadStore_GetOrganizationSSOConfigsByOrganizationId_Call) Run(run func(ctx context.Context, organizationId uuid.UUID)) *MockOrganizationReadStore_GetOrganizationSSOConfigsByOrganizationId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationSSOConfigsByOrganizationId_Call) Return(_a0 []models.OrganizationSSOConfig, _a1 error) *MockOrganizationReadStore_GetOrganizationSSOConfigsByOrganizationId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationSSOConfigsByOrganizationId_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]models.OrganizationSSOConfig, error)) *MockOrganizationReadStore_GetOrganizationSSOConfigsByOrganizationId_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationsAll provides a mock function with given fields: ctx
func (_m *MockOrganizationReadStore) GetOrganizationsAll(ctx context.Context) ([]models.Organization, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationsAll")
	}

	var r0 []models.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.Organization, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.Organization); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Organization)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationReadStore_GetOrganizationsAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationsAll'
type MockOrganizationReadStore_GetOrganizationsAll_Call struct {
	*mock.Call
}

// GetOrganizationsAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockOrganizationReadStore_Expecter) GetOrganizationsAll(ctx interface{}) *MockOrganizationReadStore_GetOrganizationsAll_Call {
	return &MockOrganizationReadStore_GetOrganizationsAll_Call{Call: _e.mock.On("GetOrganizationsAll", ctx)}
}

func (_c *MockOrganizationReadStore_GetOrganizationsAll_Call) Run(run func(ctx context.Context)) *MockOrganizationReadStore_GetOrganizationsAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationsAll_Call) Return(_a0 []models.Organization, _a1 error) *MockOrganizationReadStore_GetOrganizationsAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationsAll_Call) RunAndReturn(run func(context.Context) ([]models.Organization, error)) *MockOrganizationReadStore_GetOrganizationsAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationsByMemberId provides a mock function with given fields: ctx, memberId
func (_m *MockOrganizationReadStore) GetOrganizationsByMemberId(ctx context.Context, memberId uuid.UUID) ([]models.Organization, error) {
	ret := _m.Called(ctx, memberId)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationsByMemberId")
	}

	var r0 []models.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]models.Organization, error)); ok {
		return rf(ctx, memberId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []models.Organization); ok {
		r0 = rf(ctx, memberId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Organization)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, memberId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationReadStore_GetOrganizationsByMemberId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationsByMemberId'
type MockOrganizationReadStore_GetOrganizationsByMemberId_Call struct {
	*mock.Call
}

// GetOrganizationsByMemberId is a helper method to define mock.On call
//   - ctx context.Context
//   - memberId uuid.UUID
func (_e *MockOrganizationReadStore_Expecter) GetOrganizationsByMemberId(ctx interface{}, memberId interface{}) *MockOrganizationReadStore_GetOrganizationsByMemberId_Call {
	return &MockOrganizationReadStore_GetOrganizationsByMemberId_Call{Call: _e.mock.On("GetOrganizationsByMemberId", ctx, memberId)}
}

func (_c *MockOrganizationReadStore_GetOrganizationsByMemberId_Call) Run(run func(ctx context.Context, memberId uuid.UUID)) *MockOrganizationReadStore_GetOrganizationsByMemberId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationsByMemberId_Call) Return(_a0 []models.Organization, _a1 error) *MockOrganizationReadStore_GetOrganizationsByMemberId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationReadStore_GetOrganizationsByMemberId_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]models.Organization, error)) *MockOrganizationReadStore_GetOrganizationsByMemberId_Call {
	_c.Call.Return(run)
	return _c
}

// GetPrimarySSOConfigByDomain provides a mock function with given fields: ctx, domain
func (_m *MockOrganizationReadStore) GetPrimarySSOConfigByDomain(ctx context.Context, domain string) (*models.OrganizationSSOConfig, error) {
	ret := _m.Called(ctx, domain)

	if len(ret) == 0 {
		panic("no return value specified for GetPrimarySSOConfigByDomain")
	}

	var r0 *models.OrganizationSSOConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.OrganizationSSOConfig, error)); ok {
		return rf(ctx, domain)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.OrganizationSSOConfig); ok {
		r0 = rf(ctx, domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OrganizationSSOConfig)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationReadStore_GetPrimarySSOConfigByDomain_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPrimarySSOConfigByDomain'
type MockOrganizationReadStore_GetPrimarySSOConfigByDomain_Call struct {
	*mock.Call
}

// GetPrimarySSOConfigByDomain is a helper method to define mock.On call
//   - ctx context.Context
//   - domain string
func (_e *MockOrganizationReadStore_Expecter) GetPrimarySSOConfigByDomain(ctx interface{}, domain interface{}) *MockOrganizationReadStore_GetPrimarySSOConfigByDomain_Call {
	return &MockOrganizationReadStore_GetPrimarySSOConfigByDomain_Call{Call: _e.mock.On("GetPrimarySSOConfigByDomain", ctx, domain)}
}

func (_c *MockOrganizationReadStore_GetPrimarySSOConfigByDomain_Call) Run(run func(ctx context.Context, domain string)) *MockOrganizationReadStore_GetPrimarySSOConfigByDomain_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockOrganizationReadStore_GetPrimarySSOConfigByDomain_Call) Return(_a0 *models.OrganizationSSOConfig, _a1 error) *MockOrganizationReadStore_GetPrimarySSOConfigByDomain_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationReadStore_GetPrimarySSOConfigByDomain_Call) RunAndReturn(run func(context.Context, string) (*models.OrganizationSSOConfig, error)) *MockOrganizationReadStore_GetPrimarySSOConfigByDomain_Call {
	_c.Call.Return(run)
	return _c
}

// GetSSOConfigByDomain provides a mock function with given fields: ctx, domain
func (_m *MockOrganizationReadStore) GetSSOConfigByDomain(ctx context.Context, domain string) (*models.OrganizationSSOConfig, error) {
	ret := _m.Called(ctx, domain)

	if len(ret) == 0 {
		panic("no return value specified for GetSSOConfigByDomain")
	}

	var r0 *models.OrganizationSSOConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.OrganizationSSOConfig, error)); ok {
		return rf(ctx, domain)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.OrganizationSSOConfig); ok {
		r0 = rf(ctx, domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OrganizationSSOConfig)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationReadStore_GetSSOConfigByDomain_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSSOConfigByDomain'
type MockOrganizationReadStore_GetSSOConfigByDomain_Call struct {
	*mock.Call
}

// GetSSOConfigByDomain is a helper method to define mock.On call
//   - ctx context.Context
//   - domain string
func (_e *MockOrganizationReadStore_Expecter) GetSSOConfigByDomain(ctx interface{}, domain interface{}) *MockOrganizationReadStore_GetSSOConfigByDomain_Call {
	return &MockOrganizationReadStore_GetSSOConfigByDomain_Call{Call: _e.mock.On("GetSSOConfigByDomain", ctx, domain)}
}

func (_c *MockOrganizationReadStore_GetSSOConfigByDomain_Call) Run(run func(ctx context.Context, domain string)) *MockOrganizationReadStore_GetSSOConfigByDomain_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockOrganizationReadStore_GetSSOConfigByDomain_Call) Return(_a0 *models.OrganizationSSOConfig, _a1 error) *MockOrganizationReadStore_GetSSOConfigByDomain_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationReadStore_GetSSOConfigByDomain_Call) RunAndReturn(run func(context.Context, string) (*models.OrganizationSSOConfig, error)) *MockOrganizationReadStore_GetSSOConfigByDomain_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockOrganizationReadStore creates a new instance of MockOrganizationReadStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockOrganizationReadStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockOrganizationReadStore {
	mock := &MockOrganizationReadStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

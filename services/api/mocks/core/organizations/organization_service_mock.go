// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_organizations

import (
	context "context"

	models "github.com/Zampfi/application-platform/services/api/db/models"
	mock "github.com/stretchr/testify/mock"

	organizations "github.com/Zampfi/application-platform/services/api/core/organizations"

	teams "github.com/Zampfi/application-platform/services/api/core/organizations/teams"

	uuid "github.com/google/uuid"
)

// MockOrganizationService is an autogenerated mock type for the OrganizationService type
type MockOrganizationService struct {
	mock.Mock
}

type MockOrganizationService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockOrganizationService) EXPECT() *MockOrganizationService_Expecter {
	return &MockOrganizationService_Expecter{mock: &_m.Mock}
}

// ApprovePendingOrganizationMembershipRequest provides a mock function with given fields: ctx, organizationId, userId
func (_m *MockOrganizationService) ApprovePendingOrganizationMembershipRequest(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID) (*models.OrganizationMembershipRequest, error) {
	ret := _m.Called(ctx, organizationId, userId)

	if len(ret) == 0 {
		panic("no return value specified for ApprovePendingOrganizationMembershipRequest")
	}

	var r0 *models.OrganizationMembershipRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) (*models.OrganizationMembershipRequest, error)); ok {
		return rf(ctx, organizationId, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) *models.OrganizationMembershipRequest); ok {
		r0 = rf(ctx, organizationId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OrganizationMembershipRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, organizationId, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationService_ApprovePendingOrganizationMembershipRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ApprovePendingOrganizationMembershipRequest'
type MockOrganizationService_ApprovePendingOrganizationMembershipRequest_Call struct {
	*mock.Call
}

// ApprovePendingOrganizationMembershipRequest is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
//   - userId uuid.UUID
func (_e *MockOrganizationService_Expecter) ApprovePendingOrganizationMembershipRequest(ctx interface{}, organizationId interface{}, userId interface{}) *MockOrganizationService_ApprovePendingOrganizationMembershipRequest_Call {
	return &MockOrganizationService_ApprovePendingOrganizationMembershipRequest_Call{Call: _e.mock.On("ApprovePendingOrganizationMembershipRequest", ctx, organizationId, userId)}
}

func (_c *MockOrganizationService_ApprovePendingOrganizationMembershipRequest_Call) Run(run func(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID)) *MockOrganizationService_ApprovePendingOrganizationMembershipRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MockOrganizationService_ApprovePendingOrganizationMembershipRequest_Call) Return(_a0 *models.OrganizationMembershipRequest, _a1 error) *MockOrganizationService_ApprovePendingOrganizationMembershipRequest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationService_ApprovePendingOrganizationMembershipRequest_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID) (*models.OrganizationMembershipRequest, error)) *MockOrganizationService_ApprovePendingOrganizationMembershipRequest_Call {
	_c.Call.Return(run)
	return _c
}

// BulkInviteMembers provides a mock function with given fields: ctx, organizationId, payload
func (_m *MockOrganizationService) BulkInviteMembers(ctx context.Context, organizationId uuid.UUID, payload organizations.BulkInvitationPayload) ([]models.OrganizationInvitation, organizations.BulkInvitationError) {
	ret := _m.Called(ctx, organizationId, payload)

	if len(ret) == 0 {
		panic("no return value specified for BulkInviteMembers")
	}

	var r0 []models.OrganizationInvitation
	var r1 organizations.BulkInvitationError
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, organizations.BulkInvitationPayload) ([]models.OrganizationInvitation, organizations.BulkInvitationError)); ok {
		return rf(ctx, organizationId, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, organizations.BulkInvitationPayload) []models.OrganizationInvitation); ok {
		r0 = rf(ctx, organizationId, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.OrganizationInvitation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, organizations.BulkInvitationPayload) organizations.BulkInvitationError); ok {
		r1 = rf(ctx, organizationId, payload)
	} else {
		r1 = ret.Get(1).(organizations.BulkInvitationError)
	}

	return r0, r1
}

// MockOrganizationService_BulkInviteMembers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BulkInviteMembers'
type MockOrganizationService_BulkInviteMembers_Call struct {
	*mock.Call
}

// BulkInviteMembers is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
//   - payload organizations.BulkInvitationPayload
func (_e *MockOrganizationService_Expecter) BulkInviteMembers(ctx interface{}, organizationId interface{}, payload interface{}) *MockOrganizationService_BulkInviteMembers_Call {
	return &MockOrganizationService_BulkInviteMembers_Call{Call: _e.mock.On("BulkInviteMembers", ctx, organizationId, payload)}
}

func (_c *MockOrganizationService_BulkInviteMembers_Call) Run(run func(ctx context.Context, organizationId uuid.UUID, payload organizations.BulkInvitationPayload)) *MockOrganizationService_BulkInviteMembers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(organizations.BulkInvitationPayload))
	})
	return _c
}

func (_c *MockOrganizationService_BulkInviteMembers_Call) Return(_a0 []models.OrganizationInvitation, _a1 organizations.BulkInvitationError) *MockOrganizationService_BulkInviteMembers_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationService_BulkInviteMembers_Call) RunAndReturn(run func(context.Context, uuid.UUID, organizations.BulkInvitationPayload) ([]models.OrganizationInvitation, organizations.BulkInvitationError)) *MockOrganizationService_BulkInviteMembers_Call {
	_c.Call.Return(run)
	return _c
}

// CreateOrganization provides a mock function with given fields: ctx, name, description, ownerId
func (_m *MockOrganizationService) CreateOrganization(ctx context.Context, name string, description *string, ownerId uuid.UUID) (*models.Organization, error) {
	ret := _m.Called(ctx, name, description, ownerId)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrganization")
	}

	var r0 *models.Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *string, uuid.UUID) (*models.Organization, error)); ok {
		return rf(ctx, name, description, ownerId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *string, uuid.UUID) *models.Organization); ok {
		r0 = rf(ctx, name, description, ownerId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Organization)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *string, uuid.UUID) error); ok {
		r1 = rf(ctx, name, description, ownerId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationService_CreateOrganization_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateOrganization'
type MockOrganizationService_CreateOrganization_Call struct {
	*mock.Call
}

// CreateOrganization is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - description *string
//   - ownerId uuid.UUID
func (_e *MockOrganizationService_Expecter) CreateOrganization(ctx interface{}, name interface{}, description interface{}, ownerId interface{}) *MockOrganizationService_CreateOrganization_Call {
	return &MockOrganizationService_CreateOrganization_Call{Call: _e.mock.On("CreateOrganization", ctx, name, description, ownerId)}
}

func (_c *MockOrganizationService_CreateOrganization_Call) Run(run func(ctx context.Context, name string, description *string, ownerId uuid.UUID)) *MockOrganizationService_CreateOrganization_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*string), args[3].(uuid.UUID))
	})
	return _c
}

func (_c *MockOrganizationService_CreateOrganization_Call) Return(_a0 *models.Organization, _a1 error) *MockOrganizationService_CreateOrganization_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationService_CreateOrganization_Call) RunAndReturn(run func(context.Context, string, *string, uuid.UUID) (*models.Organization, error)) *MockOrganizationService_CreateOrganization_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllOrganizationInvitations provides a mock function with given fields: ctx, organizationId
func (_m *MockOrganizationService) GetAllOrganizationInvitations(ctx context.Context, organizationId uuid.UUID) ([]models.OrganizationInvitation, error) {
	ret := _m.Called(ctx, organizationId)

	if len(ret) == 0 {
		panic("no return value specified for GetAllOrganizationInvitations")
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

// MockOrganizationService_GetAllOrganizationInvitations_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllOrganizationInvitations'
type MockOrganizationService_GetAllOrganizationInvitations_Call struct {
	*mock.Call
}

// GetAllOrganizationInvitations is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
func (_e *MockOrganizationService_Expecter) GetAllOrganizationInvitations(ctx interface{}, organizationId interface{}) *MockOrganizationService_GetAllOrganizationInvitations_Call {
	return &MockOrganizationService_GetAllOrganizationInvitations_Call{Call: _e.mock.On("GetAllOrganizationInvitations", ctx, organizationId)}
}

func (_c *MockOrganizationService_GetAllOrganizationInvitations_Call) Run(run func(ctx context.Context, organizationId uuid.UUID)) *MockOrganizationService_GetAllOrganizationInvitations_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockOrganizationService_GetAllOrganizationInvitations_Call) Return(_a0 []models.OrganizationInvitation, _a1 error) *MockOrganizationService_GetAllOrganizationInvitations_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationService_GetAllOrganizationInvitations_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]models.OrganizationInvitation, error)) *MockOrganizationService_GetAllOrganizationInvitations_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationAudiences provides a mock function with given fields: ctx, organizationId
func (_m *MockOrganizationService) GetOrganizationAudiences(ctx context.Context, organizationId uuid.UUID) ([]models.ResourceAudiencePolicy, error) {
	ret := _m.Called(ctx, organizationId)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationAudiences")
	}

	var r0 []models.ResourceAudiencePolicy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]models.ResourceAudiencePolicy, error)); ok {
		return rf(ctx, organizationId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []models.ResourceAudiencePolicy); ok {
		r0 = rf(ctx, organizationId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ResourceAudiencePolicy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, organizationId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationService_GetOrganizationAudiences_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationAudiences'
type MockOrganizationService_GetOrganizationAudiences_Call struct {
	*mock.Call
}

// GetOrganizationAudiences is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
func (_e *MockOrganizationService_Expecter) GetOrganizationAudiences(ctx interface{}, organizationId interface{}) *MockOrganizationService_GetOrganizationAudiences_Call {
	return &MockOrganizationService_GetOrganizationAudiences_Call{Call: _e.mock.On("GetOrganizationAudiences", ctx, organizationId)}
}

func (_c *MockOrganizationService_GetOrganizationAudiences_Call) Run(run func(ctx context.Context, organizationId uuid.UUID)) *MockOrganizationService_GetOrganizationAudiences_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockOrganizationService_GetOrganizationAudiences_Call) Return(_a0 []models.ResourceAudiencePolicy, _a1 error) *MockOrganizationService_GetOrganizationAudiences_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationService_GetOrganizationAudiences_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]models.ResourceAudiencePolicy, error)) *MockOrganizationService_GetOrganizationAudiences_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationMembershipRequestsAll provides a mock function with given fields: ctx
func (_m *MockOrganizationService) GetOrganizationMembershipRequestsAll(ctx context.Context) ([]models.OrganizationMembershipRequest, error) {
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

// MockOrganizationService_GetOrganizationMembershipRequestsAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationMembershipRequestsAll'
type MockOrganizationService_GetOrganizationMembershipRequestsAll_Call struct {
	*mock.Call
}

// GetOrganizationMembershipRequestsAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockOrganizationService_Expecter) GetOrganizationMembershipRequestsAll(ctx interface{}) *MockOrganizationService_GetOrganizationMembershipRequestsAll_Call {
	return &MockOrganizationService_GetOrganizationMembershipRequestsAll_Call{Call: _e.mock.On("GetOrganizationMembershipRequestsAll", ctx)}
}

func (_c *MockOrganizationService_GetOrganizationMembershipRequestsAll_Call) Run(run func(ctx context.Context)) *MockOrganizationService_GetOrganizationMembershipRequestsAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockOrganizationService_GetOrganizationMembershipRequestsAll_Call) Return(_a0 []models.OrganizationMembershipRequest, _a1 error) *MockOrganizationService_GetOrganizationMembershipRequestsAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationService_GetOrganizationMembershipRequestsAll_Call) RunAndReturn(run func(context.Context) ([]models.OrganizationMembershipRequest, error)) *MockOrganizationService_GetOrganizationMembershipRequestsAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationMembershipRequestsByOrganizationId provides a mock function with given fields: ctx, organizationId
func (_m *MockOrganizationService) GetOrganizationMembershipRequestsByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]models.OrganizationMembershipRequest, error) {
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

// MockOrganizationService_GetOrganizationMembershipRequestsByOrganizationId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationMembershipRequestsByOrganizationId'
type MockOrganizationService_GetOrganizationMembershipRequestsByOrganizationId_Call struct {
	*mock.Call
}

// GetOrganizationMembershipRequestsByOrganizationId is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
func (_e *MockOrganizationService_Expecter) GetOrganizationMembershipRequestsByOrganizationId(ctx interface{}, organizationId interface{}) *MockOrganizationService_GetOrganizationMembershipRequestsByOrganizationId_Call {
	return &MockOrganizationService_GetOrganizationMembershipRequestsByOrganizationId_Call{Call: _e.mock.On("GetOrganizationMembershipRequestsByOrganizationId", ctx, organizationId)}
}

func (_c *MockOrganizationService_GetOrganizationMembershipRequestsByOrganizationId_Call) Run(run func(ctx context.Context, organizationId uuid.UUID)) *MockOrganizationService_GetOrganizationMembershipRequestsByOrganizationId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockOrganizationService_GetOrganizationMembershipRequestsByOrganizationId_Call) Return(_a0 []models.OrganizationMembershipRequest, _a1 error) *MockOrganizationService_GetOrganizationMembershipRequestsByOrganizationId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationService_GetOrganizationMembershipRequestsByOrganizationId_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]models.OrganizationMembershipRequest, error)) *MockOrganizationService_GetOrganizationMembershipRequestsByOrganizationId_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizations provides a mock function with given fields: ctx
func (_m *MockOrganizationService) GetOrganizations(ctx context.Context) ([]models.Organization, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizations")
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

// MockOrganizationService_GetOrganizations_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizations'
type MockOrganizationService_GetOrganizations_Call struct {
	*mock.Call
}

// GetOrganizations is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockOrganizationService_Expecter) GetOrganizations(ctx interface{}) *MockOrganizationService_GetOrganizations_Call {
	return &MockOrganizationService_GetOrganizations_Call{Call: _e.mock.On("GetOrganizations", ctx)}
}

func (_c *MockOrganizationService_GetOrganizations_Call) Run(run func(ctx context.Context)) *MockOrganizationService_GetOrganizations_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockOrganizationService_GetOrganizations_Call) Return(_a0 []models.Organization, _a1 error) *MockOrganizationService_GetOrganizations_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationService_GetOrganizations_Call) RunAndReturn(run func(context.Context) ([]models.Organization, error)) *MockOrganizationService_GetOrganizations_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveOrganizationMember provides a mock function with given fields: ctx, organizationId, userId
func (_m *MockOrganizationService) RemoveOrganizationMember(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID) error {
	ret := _m.Called(ctx, organizationId, userId)

	if len(ret) == 0 {
		panic("no return value specified for RemoveOrganizationMember")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(ctx, organizationId, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockOrganizationService_RemoveOrganizationMember_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveOrganizationMember'
type MockOrganizationService_RemoveOrganizationMember_Call struct {
	*mock.Call
}

// RemoveOrganizationMember is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
//   - userId uuid.UUID
func (_e *MockOrganizationService_Expecter) RemoveOrganizationMember(ctx interface{}, organizationId interface{}, userId interface{}) *MockOrganizationService_RemoveOrganizationMember_Call {
	return &MockOrganizationService_RemoveOrganizationMember_Call{Call: _e.mock.On("RemoveOrganizationMember", ctx, organizationId, userId)}
}

func (_c *MockOrganizationService_RemoveOrganizationMember_Call) Run(run func(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID)) *MockOrganizationService_RemoveOrganizationMember_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MockOrganizationService_RemoveOrganizationMember_Call) Return(_a0 error) *MockOrganizationService_RemoveOrganizationMember_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOrganizationService_RemoveOrganizationMember_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID) error) *MockOrganizationService_RemoveOrganizationMember_Call {
	_c.Call.Return(run)
	return _c
}

// TeamService provides a mock function with no fields
func (_m *MockOrganizationService) TeamService() teams.TeamService {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for TeamService")
	}

	var r0 teams.TeamService
	if rf, ok := ret.Get(0).(func() teams.TeamService); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(teams.TeamService)
		}
	}

	return r0
}

// MockOrganizationService_TeamService_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TeamService'
type MockOrganizationService_TeamService_Call struct {
	*mock.Call
}

// TeamService is a helper method to define mock.On call
func (_e *MockOrganizationService_Expecter) TeamService() *MockOrganizationService_TeamService_Call {
	return &MockOrganizationService_TeamService_Call{Call: _e.mock.On("TeamService")}
}

func (_c *MockOrganizationService_TeamService_Call) Run(run func()) *MockOrganizationService_TeamService_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockOrganizationService_TeamService_Call) Return(_a0 teams.TeamService) *MockOrganizationService_TeamService_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOrganizationService_TeamService_Call) RunAndReturn(run func() teams.TeamService) *MockOrganizationService_TeamService_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateMemberRole provides a mock function with given fields: ctx, organizationId, userId, privilege
func (_m *MockOrganizationService) UpdateMemberRole(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID, privilege models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error) {
	ret := _m.Called(ctx, organizationId, userId, privilege)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMemberRole")
	}

	var r0 *models.ResourceAudiencePolicy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)); ok {
		return rf(ctx, organizationId, userId, privilege)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, models.ResourcePrivilege) *models.ResourceAudiencePolicy); ok {
		r0 = rf(ctx, organizationId, userId, privilege)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ResourceAudiencePolicy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID, models.ResourcePrivilege) error); ok {
		r1 = rf(ctx, organizationId, userId, privilege)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrganizationService_UpdateMemberRole_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateMemberRole'
type MockOrganizationService_UpdateMemberRole_Call struct {
	*mock.Call
}

// UpdateMemberRole is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
//   - userId uuid.UUID
//   - privilege models.ResourcePrivilege
func (_e *MockOrganizationService_Expecter) UpdateMemberRole(ctx interface{}, organizationId interface{}, userId interface{}, privilege interface{}) *MockOrganizationService_UpdateMemberRole_Call {
	return &MockOrganizationService_UpdateMemberRole_Call{Call: _e.mock.On("UpdateMemberRole", ctx, organizationId, userId, privilege)}
}

func (_c *MockOrganizationService_UpdateMemberRole_Call) Run(run func(ctx context.Context, organizationId uuid.UUID, userId uuid.UUID, privilege models.ResourcePrivilege)) *MockOrganizationService_UpdateMemberRole_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID), args[3].(models.ResourcePrivilege))
	})
	return _c
}

func (_c *MockOrganizationService_UpdateMemberRole_Call) Return(_a0 *models.ResourceAudiencePolicy, _a1 error) *MockOrganizationService_UpdateMemberRole_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrganizationService_UpdateMemberRole_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID, models.ResourcePrivilege) (*models.ResourceAudiencePolicy, error)) *MockOrganizationService_UpdateMemberRole_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateAudienceInOrganization provides a mock function with given fields: ctx, organizationId, audienceType, audienceId
func (_m *MockOrganizationService) ValidateAudienceInOrganization(ctx context.Context, organizationId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID) error {
	ret := _m.Called(ctx, organizationId, audienceType, audienceId)

	if len(ret) == 0 {
		panic("no return value specified for ValidateAudienceInOrganization")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.AudienceType, uuid.UUID) error); ok {
		r0 = rf(ctx, organizationId, audienceType, audienceId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockOrganizationService_ValidateAudienceInOrganization_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateAudienceInOrganization'
type MockOrganizationService_ValidateAudienceInOrganization_Call struct {
	*mock.Call
}

// ValidateAudienceInOrganization is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
//   - audienceType models.AudienceType
//   - audienceId uuid.UUID
func (_e *MockOrganizationService_Expecter) ValidateAudienceInOrganization(ctx interface{}, organizationId interface{}, audienceType interface{}, audienceId interface{}) *MockOrganizationService_ValidateAudienceInOrganization_Call {
	return &MockOrganizationService_ValidateAudienceInOrganization_Call{Call: _e.mock.On("ValidateAudienceInOrganization", ctx, organizationId, audienceType, audienceId)}
}

func (_c *MockOrganizationService_ValidateAudienceInOrganization_Call) Run(run func(ctx context.Context, organizationId uuid.UUID, audienceType models.AudienceType, audienceId uuid.UUID)) *MockOrganizationService_ValidateAudienceInOrganization_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(models.AudienceType), args[3].(uuid.UUID))
	})
	return _c
}

func (_c *MockOrganizationService_ValidateAudienceInOrganization_Call) Return(_a0 error) *MockOrganizationService_ValidateAudienceInOrganization_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOrganizationService_ValidateAudienceInOrganization_Call) RunAndReturn(run func(context.Context, uuid.UUID, models.AudienceType, uuid.UUID) error) *MockOrganizationService_ValidateAudienceInOrganization_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockOrganizationService creates a new instance of MockOrganizationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockOrganizationService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockOrganizationService {
	mock := &MockOrganizationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

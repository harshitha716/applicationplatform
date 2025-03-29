// Code generated by mockery v2.50.0. DO NOT EDIT.

package mock_store

import (
	context "context"

	models "github.com/Zampfi/application-platform/services/api/db/models"
	mock "github.com/stretchr/testify/mock"

	store "github.com/Zampfi/application-platform/services/api/db/store"

	uuid "github.com/google/uuid"
)

// MockTeamStore is an autogenerated mock type for the TeamStore type
type MockTeamStore struct {
	mock.Mock
}

type MockTeamStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTeamStore) EXPECT() *MockTeamStore_Expecter {
	return &MockTeamStore_Expecter{mock: &_m.Mock}
}

// CreateOrganizationTeam provides a mock function with given fields: ctx, organizationId, team
func (_m *MockTeamStore) CreateOrganizationTeam(ctx context.Context, organizationId uuid.UUID, team models.Team) (*models.Team, error) {
	ret := _m.Called(ctx, organizationId, team)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrganizationTeam")
	}

	var r0 *models.Team
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.Team) (*models.Team, error)); ok {
		return rf(ctx, organizationId, team)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.Team) *models.Team); ok {
		r0 = rf(ctx, organizationId, team)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Team)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, models.Team) error); ok {
		r1 = rf(ctx, organizationId, team)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTeamStore_CreateOrganizationTeam_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateOrganizationTeam'
type MockTeamStore_CreateOrganizationTeam_Call struct {
	*mock.Call
}

// CreateOrganizationTeam is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
//   - team models.Team
func (_e *MockTeamStore_Expecter) CreateOrganizationTeam(ctx interface{}, organizationId interface{}, team interface{}) *MockTeamStore_CreateOrganizationTeam_Call {
	return &MockTeamStore_CreateOrganizationTeam_Call{Call: _e.mock.On("CreateOrganizationTeam", ctx, organizationId, team)}
}

func (_c *MockTeamStore_CreateOrganizationTeam_Call) Run(run func(ctx context.Context, organizationId uuid.UUID, team models.Team)) *MockTeamStore_CreateOrganizationTeam_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(models.Team))
	})
	return _c
}

func (_c *MockTeamStore_CreateOrganizationTeam_Call) Return(_a0 *models.Team, _a1 error) *MockTeamStore_CreateOrganizationTeam_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTeamStore_CreateOrganizationTeam_Call) RunAndReturn(run func(context.Context, uuid.UUID, models.Team) (*models.Team, error)) *MockTeamStore_CreateOrganizationTeam_Call {
	_c.Call.Return(run)
	return _c
}

// CreateTeamMembership provides a mock function with given fields: ctx, teamId, userId
func (_m *MockTeamStore) CreateTeamMembership(ctx context.Context, teamId uuid.UUID, userId uuid.UUID) (*models.TeamMembership, error) {
	ret := _m.Called(ctx, teamId, userId)

	if len(ret) == 0 {
		panic("no return value specified for CreateTeamMembership")
	}

	var r0 *models.TeamMembership
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) (*models.TeamMembership, error)); ok {
		return rf(ctx, teamId, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) *models.TeamMembership); ok {
		r0 = rf(ctx, teamId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.TeamMembership)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, teamId, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTeamStore_CreateTeamMembership_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateTeamMembership'
type MockTeamStore_CreateTeamMembership_Call struct {
	*mock.Call
}

// CreateTeamMembership is a helper method to define mock.On call
//   - ctx context.Context
//   - teamId uuid.UUID
//   - userId uuid.UUID
func (_e *MockTeamStore_Expecter) CreateTeamMembership(ctx interface{}, teamId interface{}, userId interface{}) *MockTeamStore_CreateTeamMembership_Call {
	return &MockTeamStore_CreateTeamMembership_Call{Call: _e.mock.On("CreateTeamMembership", ctx, teamId, userId)}
}

func (_c *MockTeamStore_CreateTeamMembership_Call) Run(run func(ctx context.Context, teamId uuid.UUID, userId uuid.UUID)) *MockTeamStore_CreateTeamMembership_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MockTeamStore_CreateTeamMembership_Call) Return(_a0 *models.TeamMembership, _a1 error) *MockTeamStore_CreateTeamMembership_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTeamStore_CreateTeamMembership_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID) (*models.TeamMembership, error)) *MockTeamStore_CreateTeamMembership_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteTeam provides a mock function with given fields: ctx, organizationId, teamId
func (_m *MockTeamStore) DeleteTeam(ctx context.Context, organizationId uuid.UUID, teamId uuid.UUID) error {
	ret := _m.Called(ctx, organizationId, teamId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteTeam")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(ctx, organizationId, teamId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTeamStore_DeleteTeam_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteTeam'
type MockTeamStore_DeleteTeam_Call struct {
	*mock.Call
}

// DeleteTeam is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
//   - teamId uuid.UUID
func (_e *MockTeamStore_Expecter) DeleteTeam(ctx interface{}, organizationId interface{}, teamId interface{}) *MockTeamStore_DeleteTeam_Call {
	return &MockTeamStore_DeleteTeam_Call{Call: _e.mock.On("DeleteTeam", ctx, organizationId, teamId)}
}

func (_c *MockTeamStore_DeleteTeam_Call) Run(run func(ctx context.Context, organizationId uuid.UUID, teamId uuid.UUID)) *MockTeamStore_DeleteTeam_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MockTeamStore_DeleteTeam_Call) Return(_a0 error) *MockTeamStore_DeleteTeam_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTeamStore_DeleteTeam_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID) error) *MockTeamStore_DeleteTeam_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteTeamMembership provides a mock function with given fields: ctx, teamId, teamMembershipId
func (_m *MockTeamStore) DeleteTeamMembership(ctx context.Context, teamId uuid.UUID, teamMembershipId uuid.UUID) error {
	ret := _m.Called(ctx, teamId, teamMembershipId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteTeamMembership")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(ctx, teamId, teamMembershipId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTeamStore_DeleteTeamMembership_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteTeamMembership'
type MockTeamStore_DeleteTeamMembership_Call struct {
	*mock.Call
}

// DeleteTeamMembership is a helper method to define mock.On call
//   - ctx context.Context
//   - teamId uuid.UUID
//   - teamMembershipId uuid.UUID
func (_e *MockTeamStore_Expecter) DeleteTeamMembership(ctx interface{}, teamId interface{}, teamMembershipId interface{}) *MockTeamStore_DeleteTeamMembership_Call {
	return &MockTeamStore_DeleteTeamMembership_Call{Call: _e.mock.On("DeleteTeamMembership", ctx, teamId, teamMembershipId)}
}

func (_c *MockTeamStore_DeleteTeamMembership_Call) Run(run func(ctx context.Context, teamId uuid.UUID, teamMembershipId uuid.UUID)) *MockTeamStore_DeleteTeamMembership_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MockTeamStore_DeleteTeamMembership_Call) Return(_a0 error) *MockTeamStore_DeleteTeamMembership_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTeamStore_DeleteTeamMembership_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID) error) *MockTeamStore_DeleteTeamMembership_Call {
	_c.Call.Return(run)
	return _c
}

// GetTeam provides a mock function with given fields: ctx, organizationId, teamId
func (_m *MockTeamStore) GetTeam(ctx context.Context, organizationId uuid.UUID, teamId uuid.UUID) (*models.Team, error) {
	ret := _m.Called(ctx, organizationId, teamId)

	if len(ret) == 0 {
		panic("no return value specified for GetTeam")
	}

	var r0 *models.Team
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) (*models.Team, error)); ok {
		return rf(ctx, organizationId, teamId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) *models.Team); ok {
		r0 = rf(ctx, organizationId, teamId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Team)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, organizationId, teamId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTeamStore_GetTeam_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTeam'
type MockTeamStore_GetTeam_Call struct {
	*mock.Call
}

// GetTeam is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
//   - teamId uuid.UUID
func (_e *MockTeamStore_Expecter) GetTeam(ctx interface{}, organizationId interface{}, teamId interface{}) *MockTeamStore_GetTeam_Call {
	return &MockTeamStore_GetTeam_Call{Call: _e.mock.On("GetTeam", ctx, organizationId, teamId)}
}

func (_c *MockTeamStore_GetTeam_Call) Run(run func(ctx context.Context, organizationId uuid.UUID, teamId uuid.UUID)) *MockTeamStore_GetTeam_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MockTeamStore_GetTeam_Call) Return(_a0 *models.Team, _a1 error) *MockTeamStore_GetTeam_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTeamStore_GetTeam_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID) (*models.Team, error)) *MockTeamStore_GetTeam_Call {
	_c.Call.Return(run)
	return _c
}

// GetTeamByName provides a mock function with given fields: ctx, organizationId, name
func (_m *MockTeamStore) GetTeamByName(ctx context.Context, organizationId uuid.UUID, name string) ([]models.Team, error) {
	ret := _m.Called(ctx, organizationId, name)

	if len(ret) == 0 {
		panic("no return value specified for GetTeamByName")
	}

	var r0 []models.Team
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) ([]models.Team, error)); ok {
		return rf(ctx, organizationId, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) []models.Team); ok {
		r0 = rf(ctx, organizationId, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Team)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, string) error); ok {
		r1 = rf(ctx, organizationId, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTeamStore_GetTeamByName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTeamByName'
type MockTeamStore_GetTeamByName_Call struct {
	*mock.Call
}

// GetTeamByName is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
//   - name string
func (_e *MockTeamStore_Expecter) GetTeamByName(ctx interface{}, organizationId interface{}, name interface{}) *MockTeamStore_GetTeamByName_Call {
	return &MockTeamStore_GetTeamByName_Call{Call: _e.mock.On("GetTeamByName", ctx, organizationId, name)}
}

func (_c *MockTeamStore_GetTeamByName_Call) Run(run func(ctx context.Context, organizationId uuid.UUID, name string)) *MockTeamStore_GetTeamByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(string))
	})
	return _c
}

func (_c *MockTeamStore_GetTeamByName_Call) Return(_a0 []models.Team, _a1 error) *MockTeamStore_GetTeamByName_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTeamStore_GetTeamByName_Call) RunAndReturn(run func(context.Context, uuid.UUID, string) ([]models.Team, error)) *MockTeamStore_GetTeamByName_Call {
	_c.Call.Return(run)
	return _c
}

// GetTeamMembershipById provides a mock function with given fields: ctx, teamMembershipId
func (_m *MockTeamStore) GetTeamMembershipById(ctx context.Context, teamMembershipId uuid.UUID) (*models.TeamMembership, error) {
	ret := _m.Called(ctx, teamMembershipId)

	if len(ret) == 0 {
		panic("no return value specified for GetTeamMembershipById")
	}

	var r0 *models.TeamMembership
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*models.TeamMembership, error)); ok {
		return rf(ctx, teamMembershipId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *models.TeamMembership); ok {
		r0 = rf(ctx, teamMembershipId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.TeamMembership)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, teamMembershipId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTeamStore_GetTeamMembershipById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTeamMembershipById'
type MockTeamStore_GetTeamMembershipById_Call struct {
	*mock.Call
}

// GetTeamMembershipById is a helper method to define mock.On call
//   - ctx context.Context
//   - teamMembershipId uuid.UUID
func (_e *MockTeamStore_Expecter) GetTeamMembershipById(ctx interface{}, teamMembershipId interface{}) *MockTeamStore_GetTeamMembershipById_Call {
	return &MockTeamStore_GetTeamMembershipById_Call{Call: _e.mock.On("GetTeamMembershipById", ctx, teamMembershipId)}
}

func (_c *MockTeamStore_GetTeamMembershipById_Call) Run(run func(ctx context.Context, teamMembershipId uuid.UUID)) *MockTeamStore_GetTeamMembershipById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockTeamStore_GetTeamMembershipById_Call) Return(_a0 *models.TeamMembership, _a1 error) *MockTeamStore_GetTeamMembershipById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTeamStore_GetTeamMembershipById_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*models.TeamMembership, error)) *MockTeamStore_GetTeamMembershipById_Call {
	_c.Call.Return(run)
	return _c
}

// GetTeamMembershipByUserIdTeamId provides a mock function with given fields: ctx, userId, teamId
func (_m *MockTeamStore) GetTeamMembershipByUserIdTeamId(ctx context.Context, userId uuid.UUID, teamId uuid.UUID) (*models.TeamMembership, error) {
	ret := _m.Called(ctx, userId, teamId)

	if len(ret) == 0 {
		panic("no return value specified for GetTeamMembershipByUserIdTeamId")
	}

	var r0 *models.TeamMembership
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) (*models.TeamMembership, error)); ok {
		return rf(ctx, userId, teamId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) *models.TeamMembership); ok {
		r0 = rf(ctx, userId, teamId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.TeamMembership)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, userId, teamId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTeamStore_GetTeamMembershipByUserIdTeamId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTeamMembershipByUserIdTeamId'
type MockTeamStore_GetTeamMembershipByUserIdTeamId_Call struct {
	*mock.Call
}

// GetTeamMembershipByUserIdTeamId is a helper method to define mock.On call
//   - ctx context.Context
//   - userId uuid.UUID
//   - teamId uuid.UUID
func (_e *MockTeamStore_Expecter) GetTeamMembershipByUserIdTeamId(ctx interface{}, userId interface{}, teamId interface{}) *MockTeamStore_GetTeamMembershipByUserIdTeamId_Call {
	return &MockTeamStore_GetTeamMembershipByUserIdTeamId_Call{Call: _e.mock.On("GetTeamMembershipByUserIdTeamId", ctx, userId, teamId)}
}

func (_c *MockTeamStore_GetTeamMembershipByUserIdTeamId_Call) Run(run func(ctx context.Context, userId uuid.UUID, teamId uuid.UUID)) *MockTeamStore_GetTeamMembershipByUserIdTeamId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MockTeamStore_GetTeamMembershipByUserIdTeamId_Call) Return(_a0 *models.TeamMembership, _a1 error) *MockTeamStore_GetTeamMembershipByUserIdTeamId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTeamStore_GetTeamMembershipByUserIdTeamId_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID) (*models.TeamMembership, error)) *MockTeamStore_GetTeamMembershipByUserIdTeamId_Call {
	_c.Call.Return(run)
	return _c
}

// GetTeamMemberships provides a mock function with given fields: ctx, teamId
func (_m *MockTeamStore) GetTeamMemberships(ctx context.Context, teamId uuid.UUID) ([]models.TeamMembership, error) {
	ret := _m.Called(ctx, teamId)

	if len(ret) == 0 {
		panic("no return value specified for GetTeamMemberships")
	}

	var r0 []models.TeamMembership
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]models.TeamMembership, error)); ok {
		return rf(ctx, teamId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []models.TeamMembership); ok {
		r0 = rf(ctx, teamId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.TeamMembership)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, teamId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTeamStore_GetTeamMemberships_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTeamMemberships'
type MockTeamStore_GetTeamMemberships_Call struct {
	*mock.Call
}

// GetTeamMemberships is a helper method to define mock.On call
//   - ctx context.Context
//   - teamId uuid.UUID
func (_e *MockTeamStore_Expecter) GetTeamMemberships(ctx interface{}, teamId interface{}) *MockTeamStore_GetTeamMemberships_Call {
	return &MockTeamStore_GetTeamMemberships_Call{Call: _e.mock.On("GetTeamMemberships", ctx, teamId)}
}

func (_c *MockTeamStore_GetTeamMemberships_Call) Run(run func(ctx context.Context, teamId uuid.UUID)) *MockTeamStore_GetTeamMemberships_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockTeamStore_GetTeamMemberships_Call) Return(_a0 []models.TeamMembership, _a1 error) *MockTeamStore_GetTeamMemberships_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTeamStore_GetTeamMemberships_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]models.TeamMembership, error)) *MockTeamStore_GetTeamMemberships_Call {
	_c.Call.Return(run)
	return _c
}

// GetTeams provides a mock function with given fields: ctx, organizationId
func (_m *MockTeamStore) GetTeams(ctx context.Context, organizationId uuid.UUID) ([]models.Team, error) {
	ret := _m.Called(ctx, organizationId)

	if len(ret) == 0 {
		panic("no return value specified for GetTeams")
	}

	var r0 []models.Team
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]models.Team, error)); ok {
		return rf(ctx, organizationId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []models.Team); ok {
		r0 = rf(ctx, organizationId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Team)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, organizationId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTeamStore_GetTeams_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTeams'
type MockTeamStore_GetTeams_Call struct {
	*mock.Call
}

// GetTeams is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
func (_e *MockTeamStore_Expecter) GetTeams(ctx interface{}, organizationId interface{}) *MockTeamStore_GetTeams_Call {
	return &MockTeamStore_GetTeams_Call{Call: _e.mock.On("GetTeams", ctx, organizationId)}
}

func (_c *MockTeamStore_GetTeams_Call) Run(run func(ctx context.Context, organizationId uuid.UUID)) *MockTeamStore_GetTeams_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockTeamStore_GetTeams_Call) Return(_a0 []models.Team, _a1 error) *MockTeamStore_GetTeams_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTeamStore_GetTeams_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]models.Team, error)) *MockTeamStore_GetTeams_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateTeam provides a mock function with given fields: ctx, organizationId, team
func (_m *MockTeamStore) UpdateTeam(ctx context.Context, organizationId uuid.UUID, team models.Team) (*models.Team, error) {
	ret := _m.Called(ctx, organizationId, team)

	if len(ret) == 0 {
		panic("no return value specified for UpdateTeam")
	}

	var r0 *models.Team
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.Team) (*models.Team, error)); ok {
		return rf(ctx, organizationId, team)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.Team) *models.Team); ok {
		r0 = rf(ctx, organizationId, team)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Team)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, models.Team) error); ok {
		r1 = rf(ctx, organizationId, team)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTeamStore_UpdateTeam_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateTeam'
type MockTeamStore_UpdateTeam_Call struct {
	*mock.Call
}

// UpdateTeam is a helper method to define mock.On call
//   - ctx context.Context
//   - organizationId uuid.UUID
//   - team models.Team
func (_e *MockTeamStore_Expecter) UpdateTeam(ctx interface{}, organizationId interface{}, team interface{}) *MockTeamStore_UpdateTeam_Call {
	return &MockTeamStore_UpdateTeam_Call{Call: _e.mock.On("UpdateTeam", ctx, organizationId, team)}
}

func (_c *MockTeamStore_UpdateTeam_Call) Run(run func(ctx context.Context, organizationId uuid.UUID, team models.Team)) *MockTeamStore_UpdateTeam_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(models.Team))
	})
	return _c
}

func (_c *MockTeamStore_UpdateTeam_Call) Return(_a0 *models.Team, _a1 error) *MockTeamStore_UpdateTeam_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTeamStore_UpdateTeam_Call) RunAndReturn(run func(context.Context, uuid.UUID, models.Team) (*models.Team, error)) *MockTeamStore_UpdateTeam_Call {
	_c.Call.Return(run)
	return _c
}

// WithTeamTransaction provides a mock function with given fields: ctx, fn
func (_m *MockTeamStore) WithTeamTransaction(ctx context.Context, fn func(store.TeamStore) error) error {
	ret := _m.Called(ctx, fn)

	if len(ret) == 0 {
		panic("no return value specified for WithTeamTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(store.TeamStore) error) error); ok {
		r0 = rf(ctx, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTeamStore_WithTeamTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTeamTransaction'
type MockTeamStore_WithTeamTransaction_Call struct {
	*mock.Call
}

// WithTeamTransaction is a helper method to define mock.On call
//   - ctx context.Context
//   - fn func(store.TeamStore) error
func (_e *MockTeamStore_Expecter) WithTeamTransaction(ctx interface{}, fn interface{}) *MockTeamStore_WithTeamTransaction_Call {
	return &MockTeamStore_WithTeamTransaction_Call{Call: _e.mock.On("WithTeamTransaction", ctx, fn)}
}

func (_c *MockTeamStore_WithTeamTransaction_Call) Run(run func(ctx context.Context, fn func(store.TeamStore) error)) *MockTeamStore_WithTeamTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(store.TeamStore) error))
	})
	return _c
}

func (_c *MockTeamStore_WithTeamTransaction_Call) Return(_a0 error) *MockTeamStore_WithTeamTransaction_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTeamStore_WithTeamTransaction_Call) RunAndReturn(run func(context.Context, func(store.TeamStore) error) error) *MockTeamStore_WithTeamTransaction_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTeamStore creates a new instance of MockTeamStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTeamStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTeamStore {
	mock := &MockTeamStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

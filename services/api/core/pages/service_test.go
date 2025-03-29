package pages

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupTest(t *testing.T) (PagesService, *mock_store.MockStore, context.Context) {
	mockStore := mock_store.NewMockStore(t)
	service := NewPagesService(mockStore)

	// Create a context with logger
	ctx := context.Background()

	return service, mockStore, ctx
}

func TestNewPagesService(t *testing.T) {
	mockStore := mock_store.NewMockStore(t)
	service := NewPagesService(mockStore)

	assert.NotNil(t, service)
	assert.Equal(t, mockStore, service.store)
}

func TestStructImplementsInterface(t *testing.T) {
	var _ PagesService = &pagesService{}
}

func TestGetPagesAll(t *testing.T) {

	pageId := uuid.New()
	tests := []struct {
		name          string
		setupMock     func(*mock_store.MockStore)
		expectedPages []models.Page
		expectedErr   error
	}{
		{
			name: "successful retrieval",
			setupMock: func(m *mock_store.MockStore) {
				m.EXPECT().GetPagesAll(mock.Anything, models.PageFilters{SortParams: []models.PageSortParams{{Column: "created_at", Desc: false}}}).
					Return([]models.Page{{ID: pageId}}, nil)
			},
			expectedPages: []models.Page{{ID: pageId}},
			expectedErr:   nil,
		},
		{
			name: "store error",
			setupMock: func(m *mock_store.MockStore) {
				m.On("GetPagesAll", mock.Anything, models.PageFilters{SortParams: []models.PageSortParams{{Column: "created_at", Desc: false}}}).
					Return(nil, errors.New("store error"))
			},
			expectedPages: nil,
			expectedErr:   errors.New("store error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockStore, ctx := setupTest(t)

			// Setup mock expectations
			tt.setupMock(mockStore)

			// Execute
			pages, err := service.GetPagesAll(ctx)

			// Verify
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, pages, len(tt.expectedPages))
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestGetPageByID(t *testing.T) {
	pageID := uuid.New()

	sheetId := uuid.New()

	tests := []struct {
		name         string
		setupMock    func(*mock_store.MockStore)
		expectedPage *models.Page
		expectedErr  error
	}{
		{
			name: "successful retrieval",
			setupMock: func(m *mock_store.MockStore) {
				m.EXPECT().GetPageById(mock.Anything, pageID).
					Return(&models.Page{ID: pageID, Sheets: []models.Sheet{{ID: sheetId}}}, nil)
			},
			expectedPage: &models.Page{ID: pageID, Sheets: []models.Sheet{{ID: sheetId}}},
			expectedErr:  nil,
		},
		{
			name: "store error",
			setupMock: func(m *mock_store.MockStore) {
				m.EXPECT().GetPageById(mock.Anything, pageID).
					Return(nil, errors.New("store error"))
			},
			expectedPage: nil,
			expectedErr:  errors.New("store error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockStore, ctx := setupTest(t)

			// Setup mock expectations
			tt.setupMock(mockStore)

			// Execute
			page, err := service.GetPageByID(ctx, pageID)

			// Verify
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPage.ID, page.ID)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestPagesAudiences(t *testing.T) {
	t.Parallel()

	pageId := uuid.New()
	audience1 := uuid.New()
	audience2 := uuid.New()
	testAudiences := []models.ResourceAudiencePolicy{{
		ID:                   uuid.New(),
		ResourceID:           pageId,
		ResourceType:         "page",
		ResourceAudienceType: "user",
		ResourceAudienceID:   audience1,
		User:                 &models.User{ID: audience1, Email: "audience1@user.com", Name: "audience1"},
	}, {
		ID:                   uuid.New(),
		ResourceID:           pageId,
		ResourceType:         "page",
		ResourceAudienceType: "user",
		ResourceAudienceID:   audience2,
		User:                 &models.User{ID: audience2, Email: "audience2@user.com", Name: "audience2"},
	}}

	tests := []struct {
		name      string
		want      []models.ResourceAudiencePolicy
		mockSetup func(*mock_store.MockStore)
		wantErr   bool
	}{
		{
			name: "success",
			want: testAudiences,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return(testAudiences, nil)
			},
			wantErr: false,
		},
		{
			name: "store error",
			want: nil,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return(nil, errors.New("test error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockStore := mock_store.NewMockStore(t)
			ctx := context.Background()
			tt.mockSetup(mockStore)

			// Execute
			service := NewPagesService(mockStore)
			got, err := service.GetPageAudiences(ctx, pageId)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAddAudienceToPage_TeamWithoutMembers(t *testing.T) {
	currentUserId := uuid.New()
	organizationId := uuid.New()
	pageId := uuid.New()
	teamId := uuid.New()

	// Setup
	mockStore := mock_store.NewMockStore(t)
	mockStore.EXPECT().WithPageTransaction(mock.Anything, mock.Anything).Return(nil).Run(func(ctx context.Context, fn func(store.PageStore) error) {
		mockStore.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
			{
				ResourceAudienceType: models.AudienceTypeUser,
				ResourceAudienceID:   currentUserId,
				Privilege:            models.PrivilegePageAdmin,
				UserPolicies: []models.FlattenedResourceAudiencePolicy{
					{
						UserId:    currentUserId,
						Privilege: models.PrivilegePageAdmin,
					},
				},
			},
		}, nil)
		// Create a team policy with empty UserPolicies to simulate a team without members
		mockStore.EXPECT().CreatePagePolicy(mock.Anything, pageId, models.AudienceTypeTeam, teamId, models.PrivilegePageAdmin).Return(&models.ResourceAudiencePolicy{
			ResourceAudienceType: models.AudienceTypeTeam,
			ResourceAudienceID:   teamId,
			Privilege:            models.PrivilegePageAdmin,
			UserPolicies:         []models.FlattenedResourceAudiencePolicy{},
		}, nil)
		fn(mockStore)
	})

	ctx := apicontext.AddAuthToContext(context.Background(), "user", currentUserId, []uuid.UUID{organizationId})
	service := NewPagesService(mockStore)

	// Execute
	got, err := service.AddAudienceToPage(ctx, pageId, models.AudienceTypeTeam, teamId, models.PrivilegePageAdmin)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, got)
}

func TestAddAudienceToPage(t *testing.T) {
	currentUserId := uuid.New()
	organizationId := uuid.New()
	pageId := uuid.New()
	policyId := uuid.New()
	teamId := uuid.New()
	tests := []struct {
		name          string
		audienceType  models.AudienceType
		audienceId    uuid.UUID
		privilege     models.ResourcePrivilege
		currentUserId *uuid.UUID
		mockSetup     func(*mock_store.MockStore)
		wantErr       bool
		expectedErr   string
	}{
		{
			name:          "success - add user audience",
			audienceType:  models.AudienceTypeUser,
			audienceId:    uuid.New(),
			privilege:     models.PrivilegePageAdmin,
			currentUserId: &currentUserId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithPageTransaction(mock.Anything, mock.Anything).Return(nil).Run(func(ctx context.Context, fn func(store.PageStore) error) {
					m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
						{
							ID:                   policyId,
							ResourceAudienceType: models.AudienceTypeUser,
							ResourceAudienceID:   currentUserId,
							Privilege:            models.PrivilegePageAdmin,
							UserPolicies: []models.FlattenedResourceAudiencePolicy{
								{
									ResourceAudiencePolicyId: policyId,
									UserId:                   currentUserId,
									ResourceAudienceType:     string(models.AudienceTypeUser),
									Privilege:                models.PrivilegePageAdmin,
								},
							},
						},
					}, nil)
					m.EXPECT().CreatePagePolicy(mock.Anything, pageId, models.AudienceTypeUser, mock.Anything, models.PrivilegePageAdmin).Return(&models.ResourceAudiencePolicy{}, nil)
					fn(m)
				})
			},
		},
		{
			name:          "success - add team audience",
			audienceType:  models.AudienceTypeTeam,
			audienceId:    teamId,
			privilege:     models.PrivilegePageAdmin,
			currentUserId: &currentUserId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithPageTransaction(mock.Anything, mock.Anything).Return(nil).Run(func(ctx context.Context, fn func(store.PageStore) error) {
					m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
						{
							ID:                   policyId,
							ResourceAudienceType: models.AudienceTypeUser,
							ResourceAudienceID:   currentUserId,
							Privilege:            models.PrivilegePageAdmin,
							UserPolicies: []models.FlattenedResourceAudiencePolicy{
								{
									ResourceAudiencePolicyId: policyId,
									UserId:                   currentUserId,
									ResourceAudienceType:     string(models.AudienceTypeUser),
									Privilege:                models.PrivilegePageAdmin,
								},
							},
						},
					}, nil)
					m.EXPECT().CreatePagePolicy(mock.Anything, pageId, models.AudienceTypeTeam, teamId, models.PrivilegePageAdmin).Return(&models.ResourceAudiencePolicy{}, nil)
					fn(m)
				})
			},
		},
		{
			name:          "success - add organization audience",
			audienceType:  models.AudienceTypeOrganization,
			audienceId:    organizationId,
			privilege:     models.PrivilegePageAdmin,
			currentUserId: &currentUserId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithPageTransaction(mock.Anything, mock.Anything).Return(nil).Run(func(ctx context.Context, fn func(store.PageStore) error) {
					m.EXPECT().GetPageById(mock.Anything, pageId).Return(&models.Page{OrganizationId: organizationId}, nil)
					m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
						{
							ID:                   policyId,
							ResourceAudienceType: models.AudienceTypeUser,
							ResourceAudienceID:   currentUserId,
							Privilege:            models.PrivilegePageAdmin,
							UserPolicies: []models.FlattenedResourceAudiencePolicy{
								{
									ResourceAudiencePolicyId: policyId,
									UserId:                   currentUserId,
									ResourceAudienceType:     string(models.AudienceTypeUser),
									Privilege:                models.PrivilegePageAdmin,
								},
							},
						},
					}, nil)
					m.EXPECT().CreatePagePolicy(mock.Anything, pageId, models.AudienceTypeOrganization, organizationId, models.PrivilegePageAdmin).Return(&models.ResourceAudiencePolicy{}, nil)
					fn(m)
				})
			},
		},
		{
			name:          "success - user is admin through team",
			audienceType:  models.AudienceTypeUser,
			audienceId:    uuid.New(),
			privilege:     models.PrivilegePageAdmin,
			currentUserId: &currentUserId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithPageTransaction(mock.Anything, mock.Anything).Return(nil).Run(func(ctx context.Context, fn func(store.PageStore) error) {
					m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
						{
							ID:                   policyId,
							ResourceAudienceType: models.AudienceTypeTeam,
							ResourceAudienceID:   teamId,
							Privilege:            models.PrivilegePageAdmin,
							UserPolicies: []models.FlattenedResourceAudiencePolicy{
								{
									ResourceAudiencePolicyId: policyId,
									UserId:                   currentUserId,
									ResourceAudienceType:     string(models.AudienceTypeTeam),
									Privilege:                models.PrivilegePageAdmin,
								},
							},
						},
					}, nil)
					fn(m)
				})
				m.EXPECT().CreatePagePolicy(mock.Anything, pageId, models.AudienceTypeUser, mock.Anything, models.PrivilegePageAdmin).Return(&models.ResourceAudiencePolicy{}, nil)
			},
		},
		{
			name:          "error - invalid privilege",
			audienceType:  models.AudienceTypeUser,
			audienceId:    uuid.New(),
			privilege:     "invalid",
			currentUserId: &currentUserId,
			wantErr:       true,
			expectedErr:   "invalid privilege",
		},
		{
			name:          "error - no user in context",
			audienceType:  models.AudienceTypeUser,
			audienceId:    uuid.New(),
			privilege:     models.PrivilegePageAdmin,
			currentUserId: nil,
			wantErr:       true,
			expectedErr:   "no user ID found in the context",
		},
		{
			name:          "error - user already exists",
			audienceType:  models.AudienceTypeUser,
			audienceId:    currentUserId,
			privilege:     models.PrivilegePageAdmin,
			currentUserId: &currentUserId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithPageTransaction(mock.Anything, mock.Anything).Return(fmt.Errorf("audience already exists on the page")).Run(func(ctx context.Context, fn func(store.PageStore) error) {
					m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
						{
							ResourceAudienceType: models.AudienceTypeUser,
							ResourceAudienceID:   currentUserId,
							Privilege:            models.PrivilegePageAdmin,
						},
					}, nil)
					fn(m)
				})
			},
			wantErr:     true,
			expectedErr: "audience already exists on the page",
		},
		{
			name:          "error - user does not have admin access",
			audienceType:  models.AudienceTypeUser,
			audienceId:    uuid.New(),
			privilege:     models.PrivilegePageAdmin,
			currentUserId: &currentUserId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithPageTransaction(mock.Anything, mock.Anything).Return(fmt.Errorf("current user does not have access to change permissions on the page")).Run(func(ctx context.Context, fn func(store.PageStore) error) {
					m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
						{
							ResourceAudienceType: models.AudienceTypeUser,
							ResourceAudienceID:   currentUserId,
							Privilege:            models.PrivilegePageRead,
						},
					}, nil)
					fn(m)
				})
			},
			wantErr:     true,
			expectedErr: "current user does not have access to change permissions on the page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockStore := mock_store.NewMockStore(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockStore)
			}

			ctx := context.Background()
			if tt.currentUserId != nil {
				ctx = apicontext.AddAuthToContext(ctx, "user", *tt.currentUserId, []uuid.UUID{organizationId})
			}
			ctx = apicontext.AddLoggerToContext(ctx, zap.NewNop())

			service := NewPagesService(mockStore)

			// Execute
			got, err := service.AddAudienceToPage(ctx, pageId, tt.audienceType, tt.audienceId, tt.privilege)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}

func TestRemoveAudienceFromPage(t *testing.T) {
	currentUserId := uuid.New()
	organizationId := uuid.New()
	pageId := uuid.New()
	policyId1 := uuid.New()
	policyId2 := uuid.New()
	otherUserId := uuid.New()
	tests := []struct {
		name          string
		currentUserId *uuid.UUID
		audienceId    uuid.UUID
		mockSetup     func(m *mock_store.MockStore)
		wantErr       bool
		expectedErr   string
	}{
		{
			name:          "Error - No user ID in context",
			currentUserId: nil,
			audienceId:    uuid.New(),
			wantErr:       true,
			expectedErr:   "no user ID found in the context",
		},
		{
			name:       "Error - User does not have admin access",
			audienceId: uuid.New(),
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
					{
						ResourceAudienceType: models.AudienceTypeUser,
						ResourceAudienceID:   currentUserId,
						Privilege:            models.PrivilegePageRead,
					},
				}, nil)
			},
			currentUserId: &currentUserId,
			wantErr:       true,
			expectedErr:   "current user does not have access to change permissions on the page",
		},
		{
			name:       "Error - Invalid audience ID",
			audienceId: uuid.New(),
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
					{
						ResourceAudienceType: models.AudienceTypeUser,
						ResourceAudienceID:   currentUserId,
						Privilege:            models.PrivilegePageAdmin,
						UserPolicies: []models.FlattenedResourceAudiencePolicy{
							{
								UserId:    currentUserId,
								Privilege: models.PrivilegePageAdmin,
							},
						},
					},
				}, nil)
			},
			currentUserId: &currentUserId,
			wantErr:       true,
			expectedErr:   "invalid audience id",
		},
		{
			name:       "Error - User trying to remove their own admin access",
			audienceId: currentUserId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
					{
						ResourceAudienceType: models.AudienceTypeUser,
						ResourceAudienceID:   currentUserId,
						Privilege:            models.PrivilegePageAdmin,
						UserPolicies: []models.FlattenedResourceAudiencePolicy{
							{
								UserId:    currentUserId,
								Privilege: models.PrivilegePageAdmin,
							},
						},
					},
				}, nil)
			},
			currentUserId: &currentUserId,
			wantErr:       true,
			expectedErr:   "you cannot change own permissions",
		},
		{
			name:       "Success - Remove other user's access",
			audienceId: otherUserId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
					{
						ID:                   policyId1,
						ResourceAudienceType: models.AudienceTypeUser,
						ResourceAudienceID:   currentUserId,
						Privilege:            models.PrivilegePageAdmin,
						UserPolicies: []models.FlattenedResourceAudiencePolicy{
							{
								ResourceAudiencePolicyId: policyId1,
								UserId:                   currentUserId,
								ResourceAudienceType:     string(models.AudienceTypeUser),
								Privilege:                models.PrivilegePageAdmin,
							},
						},
					},
					{
						ID:                   policyId2,
						ResourceAudienceType: models.AudienceTypeUser,
						ResourceAudienceID:   otherUserId,
						Privilege:            models.PrivilegePageRead,
						UserPolicies: []models.FlattenedResourceAudiencePolicy{
							{
								ResourceAudiencePolicyId: policyId2,
								UserId:                   otherUserId,
								ResourceAudienceType:     string(models.AudienceTypeUser),
								Privilege:                models.PrivilegePageRead,
							},
						},
					},
				}, nil)
				m.EXPECT().DeletePagePolicy(mock.Anything, pageId, models.AudienceTypeUser, otherUserId).Return(nil)
			},
			currentUserId: &currentUserId,
			wantErr:       false,
		},
		{
			name:       "Error - User trying to remove their own admin access",
			audienceId: currentUserId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
					{
						ResourceAudienceType: models.AudienceTypeUser,
						ResourceAudienceID:   currentUserId,
						Privilege:            models.PrivilegePageAdmin,
						UserPolicies: []models.FlattenedResourceAudiencePolicy{
							{
								UserId:                   currentUserId,
								ResourceAudienceType:     string(models.AudienceTypeUser),
								Privilege:                models.PrivilegePageAdmin,
								ResourceAudiencePolicyId: policyId1,
							},
						},
					},
				}, nil)
			},
			currentUserId: &currentUserId,
			wantErr:       true,
			expectedErr:   "you cannot change own permissions",
		},
		{
			name:       "Error - User trying to remove their own admin access when they have admin access through organization",
			audienceId: organizationId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
					{
						ID:                   policyId1,
						ResourceAudienceType: models.AudienceTypeOrganization,
						ResourceAudienceID:   organizationId,
						Privilege:            models.PrivilegePageAdmin,
						UserPolicies: []models.FlattenedResourceAudiencePolicy{
							{
								ResourceAudiencePolicyId: policyId1,
								UserId:                   currentUserId,
								ResourceAudienceType:     string(models.AudienceTypeOrganization),
								Privilege:                models.PrivilegePageAdmin,
							},
						},
					},
				}, nil)
			},
			currentUserId: &currentUserId,
			wantErr:       true,
			expectedErr:   "you cannot change own permissions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockStore := mock_store.NewMockStore(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockStore)
			}

			ctx := context.Background()
			if tt.currentUserId != nil {
				ctx = apicontext.AddAuthToContext(ctx, "user", *tt.currentUserId, []uuid.UUID{organizationId})
			}
			ctx = apicontext.AddLoggerToContext(ctx, zap.NewNop())

			service := NewPagesService(mockStore)

			// Execute
			err := service.RemoveAudienceFromPage(ctx, pageId, tt.audienceId)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestUpdatePageAudiencePrivilege(t *testing.T) {
	currentUserId := uuid.New()
	organizationId := uuid.New()
	pageId := uuid.New()
	policyId1 := uuid.New()
	policyId2 := uuid.New()
	otherUserId := uuid.New()
	teamId := uuid.New()
	tests := []struct {
		name          string
		currentUserId *uuid.UUID
		audienceId    uuid.UUID
		privilege     models.ResourcePrivilege
		mockSetup     func(m *mock_store.MockStore)
		wantErr       bool
		expectedErr   string
	}{
		{
			name:          "Error - No user ID in context",
			currentUserId: nil,
			audienceId:    otherUserId,
			privilege:     models.PrivilegePageAdmin,
			wantErr:       true,
			expectedErr:   "no user ID found in the context",
		},
		{
			name:          "Error - Invalid privilege",
			audienceId:    otherUserId,
			privilege:     "invalid",
			mockSetup:     nil,
			currentUserId: &currentUserId,
			wantErr:       true,
			expectedErr:   "invalid privilege",
		},
		{
			name:       "Error - User does not have admin access",
			audienceId: otherUserId,
			privilege:  models.PrivilegePageAdmin,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
					{
						ResourceAudienceType: models.AudienceTypeUser,
						ResourceAudienceID:   currentUserId,
						Privilege:            models.PrivilegePageRead,
						UserPolicies: []models.FlattenedResourceAudiencePolicy{
							{
								ResourceAudiencePolicyId: policyId1,
								UserId:                   currentUserId,
								ResourceAudienceType:     string(models.AudienceTypeUser),
								Privilege:                models.PrivilegePageRead,
							},
						},
					},
					{
						ResourceAudienceType: models.AudienceTypeUser,
						ResourceAudienceID:   otherUserId,
						Privilege:            models.PrivilegePageRead,
						UserPolicies: []models.FlattenedResourceAudiencePolicy{
							{
								ResourceAudiencePolicyId: policyId2,
								UserId:                   otherUserId,
								ResourceAudienceType:     string(models.AudienceTypeUser),
								Privilege:                models.PrivilegePageRead,
							},
						},
					},
				}, nil)
			},
			currentUserId: &currentUserId,
			wantErr:       true,
			expectedErr:   "current user does not have access to change permissions on the page",
		},
		{
			name:       "Error - Invalid audience ID",
			audienceId: uuid.New(),
			privilege:  models.PrivilegePageAdmin,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
					{
						ResourceAudienceType: models.AudienceTypeUser,
						ResourceAudienceID:   currentUserId,
						Privilege:            models.PrivilegePageAdmin,
						UserPolicies: []models.FlattenedResourceAudiencePolicy{
							{
								UserId:    currentUserId,
								Privilege: models.PrivilegePageAdmin,
							},
						},
					},
				}, nil)
			},
			currentUserId: &currentUserId,
			wantErr:       true,
			expectedErr:   "invalid audience id",
		},
		{
			name:       "Success - Update other user's privilege",
			audienceId: otherUserId,
			privilege:  models.PrivilegePageAdmin,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
					{
						ID:                   policyId1,
						ResourceAudienceType: models.AudienceTypeUser,
						ResourceAudienceID:   currentUserId,
						Privilege:            models.PrivilegePageAdmin,
						UserPolicies: []models.FlattenedResourceAudiencePolicy{
							{
								ResourceAudiencePolicyId: policyId1,
								UserId:                   currentUserId,
								ResourceAudienceType:     string(models.AudienceTypeUser),
								Privilege:                models.PrivilegePageAdmin,
							},
						},
					},
					{
						ID:                   policyId2,
						ResourceAudienceType: models.AudienceTypeUser,
						ResourceAudienceID:   otherUserId,
						Privilege:            models.PrivilegePageRead,
						UserPolicies: []models.FlattenedResourceAudiencePolicy{
							{
								ResourceAudiencePolicyId: policyId2,
								UserId:                   otherUserId,
								ResourceAudienceType:     string(models.AudienceTypeUser),
								Privilege:                models.PrivilegePageRead,
							},
						},
					},
				}, nil)
				m.EXPECT().UpdatePagePolicy(mock.Anything, pageId, otherUserId, models.PrivilegePageAdmin).Return(&models.ResourceAudiencePolicy{
					ResourceAudienceType: models.AudienceTypeUser,
					ResourceAudienceID:   otherUserId,
					Privilege:            models.PrivilegePageAdmin,
				}, nil)
			},
			currentUserId: &currentUserId,
			wantErr:       false,
		},
		{
			name:       "user tries to change thier own admin privilege when they have admin privilege through a team",
			audienceId: teamId,
			privilege:  models.PrivilegePageAdmin,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
					{
						ResourceAudienceType: models.AudienceTypeTeam,
						ResourceAudienceID:   teamId,
						Privilege:            models.PrivilegePageAdmin,
						UserPolicies: []models.FlattenedResourceAudiencePolicy{
							{
								ResourceAudiencePolicyId: policyId1,
								UserId:                   currentUserId,
								ResourceAudienceType:     string(models.AudienceTypeTeam),
								Privilege:                models.PrivilegePageAdmin,
							},
						},
					},
				}, nil)
				m.EXPECT().UpdatePagePolicy(mock.Anything, pageId, teamId, models.PrivilegePageAdmin).Return(nil, fmt.Errorf("you cannot change own permissions"))
			},
			currentUserId: &currentUserId,
			wantErr:       true,
			expectedErr:   "you cannot change own permissions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockStore := mock_store.NewMockStore(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockStore)
			}

			ctx := context.Background()
			if tt.currentUserId != nil {
				ctx = apicontext.AddAuthToContext(ctx, "user", *tt.currentUserId, []uuid.UUID{organizationId})
			}
			ctx = apicontext.AddLoggerToContext(ctx, zap.NewNop())

			service := NewPagesService(mockStore)

			// Execute
			got, err := service.UpdatePageAudiencePrivilege(ctx, pageId, tt.audienceId, tt.privilege)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
				assert.Nil(t, got)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}

func TestRemoveAudienceFromPage_TeamWithoutMembers(t *testing.T) {
	currentUserId := uuid.New()
	organizationId := uuid.New()
	pageId := uuid.New()
	teamId := uuid.New()

	// Setup
	mockStore := mock_store.NewMockStore(t)
	mockStore.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
		{
			ResourceAudienceType: models.AudienceTypeUser,
			ResourceAudienceID:   currentUserId,
			Privilege:            models.PrivilegePageAdmin,
			UserPolicies: []models.FlattenedResourceAudiencePolicy{
				{
					UserId:    currentUserId,
					Privilege: models.PrivilegePageAdmin,
				},
			},
		},
		{
			ID:                   uuid.New(),
			ResourceAudienceType: models.AudienceTypeTeam,
			ResourceAudienceID:   teamId,
			Privilege:            models.PrivilegePageAdmin,
			UserPolicies:         []models.FlattenedResourceAudiencePolicy{},
		},
	}, nil)
	mockStore.EXPECT().DeletePagePolicy(mock.Anything, pageId, models.AudienceTypeTeam, teamId).Return(nil)

	ctx := apicontext.AddAuthToContext(context.Background(), "user", currentUserId, []uuid.UUID{organizationId})
	service := NewPagesService(mockStore)

	// Execute
	err := service.RemoveAudienceFromPage(ctx, pageId, teamId)

	// Assert
	assert.NoError(t, err)
}

func TestUpdatePageAudiencePrivilege_TeamWithoutMembers(t *testing.T) {
	currentUserId := uuid.New()
	organizationId := uuid.New()
	pageId := uuid.New()
	teamId := uuid.New()
	policyId := uuid.New()

	// Setup
	mockStore := mock_store.NewMockStore(t)
	mockStore.EXPECT().GetPagesPolicies(mock.Anything, pageId).Return([]models.ResourceAudiencePolicy{
		{
			ResourceAudienceType: models.AudienceTypeUser,
			ResourceAudienceID:   currentUserId,
			Privilege:            models.PrivilegePageAdmin,
			UserPolicies: []models.FlattenedResourceAudiencePolicy{
				{
					UserId:    currentUserId,
					Privilege: models.PrivilegePageAdmin,
				},
			},
		},
		{
			ID:                   policyId,
			ResourceAudienceType: models.AudienceTypeTeam,
			ResourceAudienceID:   teamId,
			Privilege:            models.PrivilegePageRead,
			UserPolicies:         []models.FlattenedResourceAudiencePolicy{},
		},
	}, nil)
	mockStore.EXPECT().UpdatePagePolicy(mock.Anything, pageId, teamId, models.PrivilegePageAdmin).Return(&models.ResourceAudiencePolicy{
		ID:                   policyId,
		ResourceAudienceType: models.AudienceTypeTeam,
		ResourceAudienceID:   teamId,
		Privilege:            models.PrivilegePageAdmin,
		UserPolicies:         []models.FlattenedResourceAudiencePolicy{},
	}, nil)

	ctx := apicontext.AddAuthToContext(context.Background(), "user", currentUserId, []uuid.UUID{organizationId})
	service := NewPagesService(mockStore)

	// Execute
	got, err := service.UpdatePageAudiencePrivilege(ctx, pageId, teamId, models.PrivilegePageAdmin)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, models.PrivilegePageAdmin, got.Privilege)
}

func TestCreatePage(t *testing.T) {
	t.Parallel()
	currentUserId := uuid.New()
	organizationId := uuid.New()

	tests := []struct {
		name          string
		payload       CreatePagePayload
		currentUserId *uuid.UUID
		mockSetup     func(*mock_store.MockStore)
		wantErr       bool
		expectedErr   string
	}{
		{
			name: "success - creates page and policy",
			payload: CreatePagePayload{
				PageName:        "Test Page",
				PageDescription: "Test Description",
			},
			currentUserId: &currentUserId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithPageTransaction(mock.Anything, mock.Anything).Return(nil).Run(func(ctx context.Context, fn func(store.PageStore) error) {
					pageDescription := "Test Description"
					m.EXPECT().CreatePage(mock.Anything, "Test Page", "Test Description").Return(&models.Page{
						ID:          uuid.New(),
						Name:        "Test Page",
						Description: &pageDescription,
					}, nil)
					m.EXPECT().CreatePagePolicy(mock.Anything, mock.Anything, models.AudienceTypeUser, currentUserId, models.PrivilegePageAdmin).Return(&models.ResourceAudiencePolicy{}, nil)
					fn(m)
				})
			},
		},
		{
			name: "error - no user in context",
			payload: CreatePagePayload{
				PageName:        "Test Page",
				PageDescription: "Test Description",
			},
			currentUserId: nil,
			mockSetup:     func(m *mock_store.MockStore) {},
			wantErr:       true,
			expectedErr:   "no user ID found in the context",
		},
		{
			name: "error - create page fails",
			payload: CreatePagePayload{
				PageName:        "Test Page",
				PageDescription: "Test Description",
			},
			currentUserId: &currentUserId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithPageTransaction(mock.Anything, mock.Anything).Return(fmt.Errorf("failed to create page")).Run(func(ctx context.Context, fn func(store.PageStore) error) {
					m.EXPECT().CreatePage(mock.Anything, "Test Page", "Test Description").Return(nil, fmt.Errorf("failed to create page"))
					fn(m)
				})
			},
			wantErr:     true,
			expectedErr: "failed to create page",
		},
		{
			name: "error - create page policy fails",
			payload: CreatePagePayload{
				PageName:        "Test Page",
				PageDescription: "Test Description",
			},
			currentUserId: &currentUserId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithPageTransaction(mock.Anything, mock.Anything).Return(fmt.Errorf("failed to create page policy")).Run(func(ctx context.Context, fn func(store.PageStore) error) {
					m.EXPECT().CreatePage(mock.Anything, "Test Page", "Test Description").Return(&models.Page{}, nil)
					m.EXPECT().CreatePagePolicy(mock.Anything, mock.Anything, models.AudienceTypeUser, currentUserId, models.PrivilegePageAdmin).Return(nil, fmt.Errorf("failed to create page policy"))
					fn(m)
				})
			},
			wantErr:     true,
			expectedErr: "failed to create page policy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockStore := mock_store.NewMockStore(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockStore)
			}

			ctx := context.Background()
			if tt.currentUserId != nil {
				ctx = apicontext.AddAuthToContext(ctx, "user", *tt.currentUserId, []uuid.UUID{organizationId})
			}
			ctx = apicontext.AddLoggerToContext(ctx, zap.NewNop())

			service := NewPagesService(mockStore)

			// Execute
			page, err := service.CreatePage(ctx, tt.payload)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
				assert.Nil(t, page)
				return
			}

			assert.NoError(t, err)
		})
	}
}

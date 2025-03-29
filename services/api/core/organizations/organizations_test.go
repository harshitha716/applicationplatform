package organizations

import (
	"context"
	"errors"
	"fmt"
	"testing"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/core/mailer"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	mock_mailer "github.com/Zampfi/application-platform/services/api/mocks/core/mailer"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestNewOrganizationService(t *testing.T) {
	t.Parallel()

	ms := mock_store.NewMockStore(t)
	service := NewOrganizationService(&serverconfig.ServerConfig{Store: ms, Env: &serverconfig.ConfigVariables{}})

	assert.NotNil(t, service)
}

func TestGetOrganizations(t *testing.T) {
	t.Parallel()

	orgID := uuid.New()
	testOrgs := []models.Organization{{ID: orgID, Name: "Test Org"}}

	tests := []struct {
		name      string
		want      []models.Organization
		mockSetup func(*mock_store.MockStore)
		wantErr   bool
	}{
		{
			name: "success",
			want: testOrgs,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetOrganizationsAll(mock.Anything).Return(testOrgs, nil)
			},
			wantErr: false,
		},
		{
			name: "store error",
			want: nil,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetOrganizationsAll(mock.Anything).Return(nil, errors.New("test error"))
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
			logger := zap.NewNop()
			ctx := context.WithValue(context.Background(), "logger", logger)
			tt.mockSetup(mockStore)

			// Execute
			service := &organizationService{store: mockStore}
			got, err := service.GetOrganizations(ctx)

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

func TestGetOrganizationAudiences(t *testing.T) {
	t.Parallel()

	orgID := uuid.New()
	audience1 := uuid.New()
	audience2 := uuid.New()
	testAudiences := []models.ResourceAudiencePolicy{{
		ID:                   uuid.New(),
		ResourceID:           orgID,
		ResourceType:         "organization",
		ResourceAudienceType: "user",
		ResourceAudienceID:   audience1,
		User:                 &models.User{ID: audience1, Email: "audience1@user.com", Name: "audience1"},
	}, {
		ID:                   uuid.New(),
		ResourceID:           orgID,
		ResourceType:         "organization",
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
				m.EXPECT().GetOrganizationPolicies(mock.Anything, orgID).Return(testAudiences, nil)
			},
			wantErr: false,
		},
		{
			name: "store error",
			want: nil,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetOrganizationPolicies(mock.Anything, orgID).Return(nil, errors.New("test error"))
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
			service := &organizationService{store: mockStore}
			got, err := service.store.GetOrganizationPolicies(ctx, orgID)

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

func TestUpdateMemberRole(t *testing.T) {
	t.Parallel()
	orgID := uuid.New()
	actorId := uuid.New()
	userID := uuid.New()
	testPolicy := &models.ResourceAudiencePolicy{
		ID:                   uuid.New(),
		ResourceID:           orgID,
		ResourceType:         "organization",
		ResourceAudienceType: "user",
		ResourceAudienceID:   userID,
		Privilege:            models.PrivilegeOrganizationSystemAdmin,
	}

	tests := []struct {
		name      string
		ctx       func() context.Context
		orgID     uuid.UUID
		userID    uuid.UUID
		privilege models.ResourcePrivilege
		mockSetup func(*mock_store.MockStore)
		want      *models.ResourceAudiencePolicy
		wantErr   bool
	}{
		{
			name: "success - update existing policy",
			ctx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", actorId, []uuid.UUID{orgID})
			},
			orgID:     orgID,
			userID:    userID,
			privilege: models.PrivilegeOrganizationSystemAdmin,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetOrganizationPolicyByUser(mock.Anything, orgID, userID).Return(
					&models.ResourceAudiencePolicy{Privilege: models.PrivilegeOrganizationMember}, nil)
				m.EXPECT().UpdateOrganizationPolicy(mock.Anything, orgID, userID, models.PrivilegeOrganizationSystemAdmin).
					Return(testPolicy, nil)
			},
			want:    testPolicy,
			wantErr: false,
		},
		{
			name: "success - no update needed (same privilege)",
			ctx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", actorId, []uuid.UUID{orgID})
			},
			orgID:     orgID,
			userID:    userID,
			privilege: models.PrivilegeOrganizationSystemAdmin,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetOrganizationPolicyByUser(mock.Anything, orgID, userID).Return(
					&models.ResourceAudiencePolicy{Privilege: models.PrivilegeOrganizationSystemAdmin}, nil)
			},
			want:    &models.ResourceAudiencePolicy{Privilege: models.PrivilegeOrganizationSystemAdmin},
			wantErr: false,
		},
		{
			name: "error - unauthorized organization",
			ctx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", actorId, []uuid.UUID{uuid.New()})
			},
			orgID:     orgID,
			userID:    userID,
			privilege: models.PrivilegeOrganizationSystemAdmin,
			mockSetup: func(m *mock_store.MockStore) {},
			want:      nil,
			wantErr:   true,
		},
		{
			name: "error - invalid privilege",
			ctx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", actorId, []uuid.UUID{orgID})
			},
			orgID:     orgID,
			userID:    userID,
			privilege: "invalid_privilege",
			mockSetup: func(m *mock_store.MockStore) {},
			want:      nil,
			wantErr:   true,
		},
		{
			name: "error - restricting user from changing their own role",
			ctx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", actorId, []uuid.UUID{orgID})
			},
			orgID:     orgID,
			userID:    actorId,
			privilege: "invalid_privilege",
			mockSetup: func(m *mock_store.MockStore) {},
			want:      nil,
			wantErr:   true,
		},
		{
			name: "error - failed to get policy",
			ctx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", actorId, []uuid.UUID{orgID})
			},
			orgID:     orgID,
			userID:    userID,
			privilege: models.PrivilegeOrganizationSystemAdmin,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetOrganizationPolicyByUser(mock.Anything, orgID, userID).Return(nil, errors.New("test error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - failed to update policy",
			ctx: func() context.Context {
				return apicontext.AddAuthToContext(context.Background(), "user", actorId, []uuid.UUID{orgID})
			},
			orgID:     orgID,
			userID:    userID,
			privilege: models.PrivilegeOrganizationSystemAdmin,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetOrganizationPolicyByUser(mock.Anything, orgID, userID).Return(
					&models.ResourceAudiencePolicy{Privilege: models.PrivilegeOrganizationMember}, nil)
				m.EXPECT().UpdateOrganizationPolicy(mock.Anything, orgID, userID, models.PrivilegeOrganizationSystemAdmin).
					Return(nil, errors.New("test error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup
			mockStore := mock_store.NewMockStore(t)
			logger := zap.NewNop()
			ctx := context.WithValue(tt.ctx(), "logger", logger)
			tt.mockSetup(mockStore)

			// Execute
			service := organizationService{store: mockStore}
			got, err := service.UpdateMemberRole(ctx, tt.orgID, tt.userID, tt.privilege)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if tt.want != nil {
				assert.Equal(t, tt.want.Privilege, got.Privilege)
			}
		})
	}
}

func TestOrganizationService_inviteMember(t *testing.T) {

	actorId := uuid.New()

	tests := []struct {
		name           string
		organizationId uuid.UUID
		userEmail      string
		privilege      models.ResourcePrivilege
		orgIds         []uuid.UUID
		mockSetup      func(m *mock_store.MockStore, mc *mock_mailer.MockMailerService)
		wantErr        bool
		expectedError  string
	}{
		{
			name:           "Success - Valid invitation",
			organizationId: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			userEmail:      "test@example.com",
			privilege:      models.PrivilegeOrganizationMember,
			orgIds:         []uuid.UUID{uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")},
			mockSetup: func(m *mock_store.MockStore, mc *mock_mailer.MockMailerService) {
				m.On("WithOrganizationTransaction", mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(store.OrganizationStore) error)
						fn(m)
					}).Return(nil)

				m.On("GetOrganizationInvitationsByOrganizationId", mock.Anything,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")).
					Return([]models.OrganizationInvitation{}, nil)

				m.On("GetOrganizationPoliciesByEmail", mock.Anything,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					"test@example.com").
					Return([]models.ResourceAudiencePolicy{}, nil)

				invitation := &models.OrganizationInvitation{
					OrganizationInvitationID: uuid.New(),
					OrganizationID:           uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					TargetEmail:              "test@example.com",
					InvitedBy:                actorId,
					Privilege:                models.PrivilegeOrganizationMember,
				}

				m.On("CreateOrganizationInvitation", mock.Anything,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					"test@example.com",
					models.PrivilegeOrganizationMember).
					Return(invitation, nil)

				m.On("GetOrganizationInvitationsAndMembershipRequests", mock.Anything,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")).
					Return(&models.Organization{
						MembershipRequests: []models.OrganizationMembershipRequest{},
					}, nil)

				m.On("GetOrganizationById", mock.Anything, mock.Anything).Return(&models.Organization{
					Name: "Test Org",
				}, nil)

				m.On("GetUserById", mock.Anything, mock.Anything).Return(&models.User{
					Name: "Test User",
				}, nil)

				mc.On("SendInvitationEmail", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:           "Failure - User does not have access to organization",
			organizationId: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			userEmail:      "test@example.com",
			privilege:      models.PrivilegeOrganizationMember,
			orgIds:         []uuid.UUID{uuid.MustParse("660e8400-e29b-41d4-a716-446655440000")}, // different org ID
			mockSetup:      func(m *mock_store.MockStore, mc *mock_mailer.MockMailerService) {},
			wantErr:        true,
			expectedError:  "forbidden",
		},
		{
			name:           "Failure - Invalid email",
			organizationId: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			userEmail:      "invalid-email",
			privilege:      models.PrivilegeOrganizationMember,
			orgIds:         []uuid.UUID{uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")},
			mockSetup:      func(m *mock_store.MockStore, mc *mock_mailer.MockMailerService) {},
			wantErr:        true,
			expectedError:  "invalid email",
		},
		{
			name:           "Failure - User already member",
			organizationId: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			userEmail:      "existing@example.com",
			privilege:      models.PrivilegeOrganizationMember,
			orgIds:         []uuid.UUID{uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")},
			mockSetup: func(m *mock_store.MockStore, mc *mock_mailer.MockMailerService) {
				m.On("WithOrganizationTransaction", mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(store.OrganizationStore) error)
						fn(m)
					}).Return(fmt.Errorf("User is already a member of the organization"))

				m.On("GetOrganizationInvitationsByOrganizationId", mock.Anything,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")).
					Return([]models.OrganizationInvitation{}, nil)

				m.On("GetOrganizationPoliciesByEmail", mock.Anything,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					"existing@example.com").
					Return([]models.ResourceAudiencePolicy{{
						ID:                   uuid.New(),
						ResourceID:           uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
						ResourceType:         models.ResourceTypeOrganization,
						ResourceAudienceID:   uuid.New(),
						ResourceAudienceType: models.AudienceTypeUser,
						Privilege:            models.PrivilegeOrganizationMember,
					}}, nil)
			},
			wantErr:       true,
			expectedError: "User is already a member of the organization",
		},
		{
			name:           "Failure - User already invited",
			organizationId: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			userEmail:      "invited@example.com",
			privilege:      models.PrivilegeOrganizationMember,
			orgIds:         []uuid.UUID{uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")},
			mockSetup: func(m *mock_store.MockStore, mc *mock_mailer.MockMailerService) {
				m.On("WithOrganizationTransaction", mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(store.OrganizationStore) error)
						fn(m)
					}).Return(fmt.Errorf("User is already invited to the organization."))

				m.On("GetOrganizationInvitationsByOrganizationId", mock.Anything,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")).
					Return([]models.OrganizationInvitation{{
						OrganizationInvitationID: uuid.New(),
						OrganizationID:           uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
						TargetEmail:              "invited@example.com",
						InvitedBy:                actorId,
						Privilege:                models.PrivilegeOrganizationMember,
					}}, nil)
			},
			wantErr:       true,
			expectedError: "User is already invited to the organization.",
		},
		{
			name:           "Success - No pending request to approve",
			organizationId: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			userEmail:      "norequest@example.com",
			privilege:      models.PrivilegeOrganizationMember,
			orgIds:         []uuid.UUID{uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")},
			mockSetup: func(m *mock_store.MockStore, mc *mock_mailer.MockMailerService) {
				m.On("WithOrganizationTransaction", mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(store.OrganizationStore) error)
						fn(m)
					}).Return(nil)

				m.On("GetOrganizationInvitationsByOrganizationId", mock.Anything,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")).
					Return([]models.OrganizationInvitation{}, nil)

				m.On("GetOrganizationPoliciesByEmail", mock.Anything,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					"norequest@example.com").
					Return([]models.ResourceAudiencePolicy{}, nil)

				invitation := &models.OrganizationInvitation{
					OrganizationInvitationID: uuid.New(),
					OrganizationID:           uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					TargetEmail:              "norequest@example.com",
					InvitedBy:                actorId,
					Privilege:                models.PrivilegeOrganizationMember,
				}

				m.On("CreateOrganizationInvitation", mock.Anything,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					"norequest@example.com",
					models.PrivilegeOrganizationMember).
					Return(invitation, nil)

				m.On("GetOrganizationInvitationsAndMembershipRequests", mock.Anything,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")).
					Return(&models.Organization{
						MembershipRequests: []models.OrganizationMembershipRequest{},
					}, nil)

				m.On("GetOrganizationById", mock.Anything, mock.Anything).Return(&models.Organization{
					Name: "Test Org",
				}, nil)

				m.On("GetUserById", mock.Anything, mock.Anything).Return(&models.User{
					Name: "Test User",
				}, nil)

				mc.On("SendInvitationEmail", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock store and mailer
			mockStore := mock_store.NewMockStore(t)
			mockMailer := mock_mailer.NewMockMailerService(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockStore, mockMailer)
			}

			// Create service with mock store and mailer
			s := &organizationService{
				store:        mockStore,
				mailerClient: mockMailer,
			}

			// Create context with required values
			ctx := apicontext.AddLoggerToContext(context.Background(), zap.NewNop())
			ctx = apicontext.AddAuthToContext(ctx, "user", actorId, tt.orgIds)

			// Call the function
			invitation, err := s.inviteMember(ctx, tt.organizationId, tt.userEmail, tt.privilege)

			// Assert results
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, invitation)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, invitation)
				assert.Equal(t, tt.userEmail, invitation.TargetEmail)
				assert.Equal(t, tt.organizationId, invitation.OrganizationID)
				assert.Equal(t, tt.privilege, invitation.Privilege)
			}

			// Verify all mocked calls were made
			mockStore.AssertExpectations(t)
			mockMailer.AssertExpectations(t)
		})
	}
}

func TestGetAllOrganizationInvitations(t *testing.T) {
	t.Parallel()

	orgID := uuid.New()
	actorID := uuid.New()
	testInvitations := []models.OrganizationInvitation{
		{
			OrganizationInvitationID: uuid.New(),
			OrganizationID:           orgID,
			TargetEmail:              "test1@example.com",
			InvitedBy:                actorID,
			Privilege:                models.PrivilegeOrganizationMember,
		},
		{
			OrganizationInvitationID: uuid.New(),
			OrganizationID:           orgID,
			TargetEmail:              "test2@example.com",
			InvitedBy:                actorID,
			Privilege:                models.PrivilegeOrganizationSystemAdmin,
		},
	}

	tests := []struct {
		name           string
		organizationID uuid.UUID
		orgIDs         []uuid.UUID
		mockSetup      func(*mock_store.MockStore)
		want           []models.OrganizationInvitation
		wantErr        bool
		expectedError  string
	}{
		{
			name:           "success",
			organizationID: orgID,
			orgIDs:         []uuid.UUID{orgID},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationInvitationsByOrganizationId(mock.Anything, orgID).
					Return(testInvitations, nil)
			},
			want:    testInvitations,
			wantErr: false,
		},
		{
			name:           "error - unauthorized organization",
			organizationID: orgID,
			orgIDs:         []uuid.UUID{uuid.New()}, // different org ID
			mockSetup:      func(m *mock_store.MockStore) {},
			want:           nil,
			wantErr:        true,
			expectedError:  "forbidden",
		},
		{
			name:           "error - store error",
			organizationID: orgID,
			orgIDs:         []uuid.UUID{orgID},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationInvitationsByOrganizationId(mock.Anything, orgID).
					Return(nil, errors.New("test error"))
			},
			want:          nil,
			wantErr:       true,
			expectedError: "test error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockStore := mock_store.NewMockStore(t)
			logger := zap.NewNop()
			ctx := context.WithValue(context.Background(), "logger", logger)
			ctx = apicontext.AddAuthToContext(ctx, "user", actorID, tt.orgIDs)
			tt.mockSetup(mockStore)

			// Execute
			service := &organizationService{store: mockStore}
			got, err := service.GetAllOrganizationInvitations(ctx, tt.organizationID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, err.Error())
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOrganizationService_RemoveOrganizationMember(t *testing.T) {
	t.Parallel()

	orgID := uuid.New()
	actorID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name           string
		organizationID uuid.UUID
		userID         uuid.UUID
		orgIDs         []uuid.UUID
		currentUserID  *uuid.UUID
		mockSetup      func(*mock_store.MockStore)
		wantErr        bool
		expectedError  string
	}{
		{
			name:           "Success - Remove member",
			organizationID: orgID,
			userID:         userID,
			currentUserID:  &actorID,
			orgIDs:         []uuid.UUID{orgID},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationPolicyByUser(mock.Anything, orgID, actorID).
					Return(&models.ResourceAudiencePolicy{
						Privilege: models.PrivilegeOrganizationSystemAdmin,
					}, nil)

				m.EXPECT().
					DeleteOrganizationPolicy(mock.Anything, orgID, userID).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:           "Error - No user ID in context",
			organizationID: orgID,
			userID:         userID,
			currentUserID:  nil,
			orgIDs:         []uuid.UUID{orgID},
			mockSetup:      func(m *mock_store.MockStore) {},
			wantErr:        true,
			expectedError:  "no user id found in context",
		},
		{
			name:           "Error - User does not have access to organization",
			organizationID: orgID,
			userID:         userID,
			currentUserID:  &actorID,
			orgIDs:         []uuid.UUID{uuid.New()}, // Different org ID
			mockSetup:      func(m *mock_store.MockStore) {},
			wantErr:        true,
			expectedError:  "forbidden",
		},
		{
			name:           "Error - User trying to remove themselves",
			organizationID: orgID,
			userID:         actorID, // Same as current user
			currentUserID:  &actorID,
			orgIDs:         []uuid.UUID{orgID},
			mockSetup:      func(m *mock_store.MockStore) {},
			wantErr:        true,
			expectedError:  "forbidden",
		},
		{
			name:           "Error - Failed to get organization policy",
			organizationID: orgID,
			userID:         userID,
			currentUserID:  &actorID,
			orgIDs:         []uuid.UUID{orgID},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationPolicyByUser(mock.Anything, orgID, actorID).
					Return(nil, errors.New("test error"))
			},
			wantErr:       true,
			expectedError: "test error",
		},
		{
			name:           "Error - User not a member of organization",
			organizationID: orgID,
			userID:         userID,
			currentUserID:  &actorID,
			orgIDs:         []uuid.UUID{orgID},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationPolicyByUser(mock.Anything, orgID, actorID).
					Return(nil, nil)
			},
			wantErr:       true,
			expectedError: "user is not a member of the organization",
		},
		{
			name:           "Error - User does not have system admin privilege",
			organizationID: orgID,
			userID:         userID,
			currentUserID:  &actorID,
			orgIDs:         []uuid.UUID{orgID},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationPolicyByUser(mock.Anything, orgID, actorID).
					Return(&models.ResourceAudiencePolicy{
						Privilege: models.PrivilegeOrganizationMember,
					}, nil)
			},
			wantErr:       true,
			expectedError: "forbidden",
		},
		{
			name:           "Error - Failed to delete organization policy",
			organizationID: orgID,
			userID:         userID,
			currentUserID:  &actorID,
			orgIDs:         []uuid.UUID{orgID},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationPolicyByUser(mock.Anything, orgID, actorID).
					Return(&models.ResourceAudiencePolicy{
						Privilege: models.PrivilegeOrganizationSystemAdmin,
					}, nil)

				m.EXPECT().
					DeleteOrganizationPolicy(mock.Anything, orgID, userID).
					Return(errors.New("test error"))
			},
			wantErr:       true,
			expectedError: "test error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockStore := mock_store.NewMockStore(t)
			logger := zap.NewNop()
			ctx := apicontext.AddLoggerToContext(context.Background(), logger)
			if tt.currentUserID != nil {
				ctx = apicontext.AddAuthToContext(ctx, "user", *tt.currentUserID, tt.orgIDs)
			}
			tt.mockSetup(mockStore)

			// Execute
			service := &organizationService{store: mockStore}
			err := service.RemoveOrganizationMember(ctx, tt.organizationID, tt.userID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, err.Error())
				}
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestOrganizationService_CreateOrganization(t *testing.T) {
	t.Parallel()
	orgID := uuid.New()
	ownerID := uuid.New()
	description := "Test Description"
	testOrg := &models.Organization{
		ID:          orgID,
		Name:        "Test Org",
		Description: &description,
		OwnerId:     ownerID,
	}

	tests := []struct {
		name          string
		role          string
		orgName       string
		description   *string
		ownerID       uuid.UUID
		mockSetup     func(*mock_store.MockStore)
		want          *models.Organization
		wantErr       bool
		expectedError string
	}{
		{
			name:        "success",
			role:        "admin",
			orgName:     "Test Org",
			description: &description,
			ownerID:     ownerID,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					WithOrganizationTransaction(mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					Run(func(ctx context.Context, fn func(store.OrganizationStore) error) {
						fn(m)
					}).
					Return(nil)

				m.EXPECT().
					CreateOrganization(mock.Anything, "Test Org", &description, ownerID).
					Return(testOrg, nil)

				m.EXPECT().
					CreateOrganizationPolicy(mock.Anything, orgID, models.AudienceTypeUser, ownerID, models.PrivilegeOrganizationSystemAdmin).
					Return(&models.ResourceAudiencePolicy{}, nil)
			},
			want:    testOrg,
			wantErr: false,
		},
		{
			name:          "non-admin access",
			role:          "user",
			orgName:       "Test Org",
			description:   &description,
			ownerID:       ownerID,
			mockSetup:     func(m *mock_store.MockStore) {},
			want:          nil,
			wantErr:       true,
			expectedError: "forbidden",
		},
		{
			name:        "create organization fails",
			role:        "admin",
			orgName:     "Test Org",
			description: &description,
			ownerID:     ownerID,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					WithOrganizationTransaction(mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					Run(func(ctx context.Context, fn func(store.OrganizationStore) error) {
						fn(m)
					}).
					Return(fmt.Errorf("test error"))

				m.EXPECT().
					CreateOrganization(mock.Anything, "Test Org", &description, ownerID).
					Return(nil, errors.New("test error"))
			},
			want:          nil,
			wantErr:       true,
			expectedError: "test error",
		},
		{
			name:        "create policy fails",
			role:        "admin",
			orgName:     "Test Org",
			description: &description,
			ownerID:     ownerID,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					WithOrganizationTransaction(mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					Run(func(ctx context.Context, fn func(store.OrganizationStore) error) {
						fn(m)
					}).
					Return(fmt.Errorf("test error"))

				m.EXPECT().
					CreateOrganization(mock.Anything, "Test Org", &description, ownerID).
					Return(testOrg, nil)

				m.EXPECT().
					CreateOrganizationPolicy(mock.Anything, orgID, models.AudienceTypeUser, ownerID, models.PrivilegeOrganizationSystemAdmin).
					Return(nil, errors.New("test error"))
			},
			want:          nil,
			wantErr:       true,
			expectedError: "test error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockStore := mock_store.NewMockStore(t)
			logger := zap.NewNop()
			ctx := apicontext.AddLoggerToContext(context.Background(), logger)
			ctx = apicontext.AddAuthToContext(ctx, tt.role, uuid.New(), nil)
			tt.mockSetup(mockStore)

			// Execute
			service := &organizationService{store: mockStore}
			got, err := service.CreateOrganization(ctx, tt.orgName, tt.description, tt.ownerID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, err.Error())
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetOrganizationMembershipRequestsByOrganizationId(t *testing.T) {
	t.Parallel()

	orgID := uuid.New()
	testRequests := []models.OrganizationMembershipRequest{{ID: uuid.New(), OrganizationID: orgID}}

	tests := []struct {
		name          string
		orgID         uuid.UUID
		want          []models.OrganizationMembershipRequest
		mockSetup     func(*mock_store.MockStore)
		wantErr       bool
		expectedError string
	}{
		{
			name:  "success",
			orgID: orgID,
			want:  testRequests,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationMembershipRequestsByOrganizationId(mock.Anything, orgID).
					Return(testRequests, nil)
			},
			wantErr: false,
		},
		{
			name:  "store error",
			orgID: orgID,
			want:  nil,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationMembershipRequestsByOrganizationId(mock.Anything, orgID).
					Return(nil, errors.New("test error"))
			},
			wantErr:       true,
			expectedError: "test error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockStore := mock_store.NewMockStore(t)
			logger := zap.NewNop()
			ctx := apicontext.AddLoggerToContext(context.Background(), logger)
			tt.mockSetup(mockStore)

			// Execute
			service := &organizationService{store: mockStore}
			got, err := service.GetOrganizationMembershipRequestsByOrganizationId(ctx, tt.orgID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, err.Error())
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetOrganizationMembershipRequestsAll(t *testing.T) {
	t.Parallel()

	orgID := uuid.New()
	testRequests := []models.OrganizationMembershipRequest{{ID: uuid.New(), OrganizationID: orgID}}

	tests := []struct {
		name          string
		want          []models.OrganizationMembershipRequest
		mockSetup     func(*mock_store.MockStore)
		wantErr       bool
		expectedError string
	}{
		{
			name: "success",
			want: testRequests,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationMembershipRequestsAll(mock.Anything).
					Return(testRequests, nil)
			},
			wantErr: false,
		},
		{
			name: "store error",
			want: nil,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationMembershipRequestsAll(mock.Anything).
					Return(nil, errors.New("test error"))
			},
			wantErr:       true,
			expectedError: "test error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockStore := mock_store.NewMockStore(t)
			logger := zap.NewNop()
			ctx := apicontext.AddLoggerToContext(context.Background(), logger)
			tt.mockSetup(mockStore)

			// Execute
			service := &organizationService{store: mockStore}
			got, err := service.GetOrganizationMembershipRequestsAll(ctx)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, err.Error())
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOrganizationService_ApprovePendingOrganizationMembershipRequest(t *testing.T) {
	testOrgId := uuid.New()
	testUserId := uuid.New()
	testCurrentUserId := uuid.New()

	testPolicies := []models.ResourceAudiencePolicy{
		{
			ResourceAudienceType: models.AudienceTypeUser,
			ResourceAudienceID:   testCurrentUserId,
			Privilege:            models.PrivilegeOrganizationSystemAdmin,
		},
	}

	testRequest := &models.OrganizationMembershipRequest{
		OrganizationID: testOrgId,
		UserID:         testUserId,
		Status:         models.OrgMembershipStatusApproved,
	}

	tests := []struct {
		name          string
		orgId         uuid.UUID
		userId        uuid.UUID
		currentUserId *uuid.UUID
		want          *models.OrganizationMembershipRequest
		mockSetup     func(m *mock_store.MockStore)
		wantErr       bool
		expectedError string
	}{
		{
			name:          "success",
			orgId:         testOrgId,
			userId:        testUserId,
			currentUserId: &testCurrentUserId,
			want:          testRequest,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					WithOrganizationTransaction(mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					RunAndReturn(func(ctx context.Context, fn func(store.OrganizationStore) error) error {
						return fn(m)
					})
				m.EXPECT().
					GetOrganizationPolicies(mock.Anything, testOrgId).
					Return(testPolicies, nil)
				m.EXPECT().
					CreateOrganizationPolicy(mock.Anything, testOrgId, models.AudienceTypeUser, testUserId, models.PrivilegeOrganizationMember).
					Return(&models.ResourceAudiencePolicy{}, nil)
				m.EXPECT().
					UpdatePendingOrganizationMembershipRequest(mock.Anything, testOrgId, testUserId, models.OrgMembershipStatusApproved).
					Return(testRequest, nil)
			},
			wantErr: false,
		},
		{
			name:          "create policy error",
			orgId:         testOrgId,
			userId:        testUserId,
			currentUserId: &testCurrentUserId,
			want:          nil,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					WithOrganizationTransaction(mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					RunAndReturn(func(ctx context.Context, fn func(store.OrganizationStore) error) error {
						return fn(m)
					})
				m.EXPECT().
					GetOrganizationPolicies(mock.Anything, testOrgId).
					Return(testPolicies, nil)
				m.EXPECT().
					CreateOrganizationPolicy(mock.Anything, testOrgId, models.AudienceTypeUser, testUserId, models.PrivilegeOrganizationMember).
					Return(nil, errors.New("failed to create policy"))
			},
			wantErr:       true,
			expectedError: "failed to create policy",
		},
		{
			name:          "no user in context",
			orgId:         testOrgId,
			userId:        testUserId,
			currentUserId: nil,
			want:          nil,
			mockSetup:     func(m *mock_store.MockStore) {},
			wantErr:       true,
			expectedError: "no user id found in context",
		},
		{
			name:          "get policies error",
			orgId:         testOrgId,
			userId:        testUserId,
			currentUserId: &testCurrentUserId,
			want:          nil,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					WithOrganizationTransaction(mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					RunAndReturn(func(ctx context.Context, fn func(store.OrganizationStore) error) error {
						return fn(m)
					})
				m.EXPECT().
					GetOrganizationPolicies(mock.Anything, testOrgId).
					Return(nil, errors.New("test error"))
			},
			wantErr:       true,
			expectedError: "test error",
		},
		{
			name:          "user not member",
			orgId:         testOrgId,
			userId:        testUserId,
			currentUserId: &testCurrentUserId,
			want:          nil,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					WithOrganizationTransaction(mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					RunAndReturn(func(ctx context.Context, fn func(store.OrganizationStore) error) error {
						return fn(m)
					})
				m.EXPECT().
					GetOrganizationPolicies(mock.Anything, testOrgId).
					Return([]models.ResourceAudiencePolicy{}, nil)
			},
			wantErr:       true,
			expectedError: "current user is not a member of the organization",
		},
		{
			name:          "user not admin",
			orgId:         testOrgId,
			userId:        testUserId,
			currentUserId: &testCurrentUserId,
			want:          nil,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					WithOrganizationTransaction(mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					RunAndReturn(func(ctx context.Context, fn func(store.OrganizationStore) error) error {
						return fn(m)
					})
				m.EXPECT().
					GetOrganizationPolicies(mock.Anything, testOrgId).
					Return([]models.ResourceAudiencePolicy{{
						ResourceAudienceType: models.AudienceTypeUser,
						ResourceAudienceID:   testCurrentUserId,
						Privilege:            models.PrivilegeOrganizationMember,
					}}, nil)
			},
			wantErr:       true,
			expectedError: "You do not have the necessary privileges to approve membership requests",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockStore := mock_store.NewMockStore(t)
			logger := zap.NewNop()
			ctx := apicontext.AddLoggerToContext(context.Background(), logger)
			if tt.currentUserId != nil {
				ctx = apicontext.AddAuthToContext(ctx, "test", *tt.currentUserId, nil)
			}
			tt.mockSetup(mockStore)

			// Execute
			service := &organizationService{store: mockStore}
			got, err := service.ApprovePendingOrganizationMembershipRequest(ctx, tt.orgId, tt.userId)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, err.Error())
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOrganizationService_validateOrganizationAccessByEmail(t *testing.T) {
	t.Parallel()

	testOrgId := uuid.New()
	testEmail := "test@example.com"
	testUser := &models.User{
		Email: testEmail,
	}

	tests := []struct {
		name          string
		orgId         uuid.UUID
		email         string
		organization  *models.Organization
		mockSetup     func(*mock_store.MockStore)
		want          userMembershipState
		wantErr       bool
		expectedError string
	}{
		{
			name:  "success - no membership state",
			orgId: testOrgId,
			email: testEmail,
			organization: &models.Organization{
				Invitations:        []models.OrganizationInvitation{},
				MembershipRequests: []models.OrganizationMembershipRequest{},
			},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationInvitationsAndMembershipRequests(mock.Anything, testOrgId).
					Return(&models.Organization{
						Invitations:        []models.OrganizationInvitation{},
						MembershipRequests: []models.OrganizationMembershipRequest{},
					}, nil)
			},
			want:    userMembershipStateNone,
			wantErr: false,
		},
		{
			name:  "success - invited state",
			orgId: testOrgId,
			email: testEmail,
			organization: &models.Organization{
				Invitations: []models.OrganizationInvitation{
					{
						TargetEmail:        testEmail,
						InvitationStatuses: []models.OrganizationInvitationStatus{},
					},
				},
			},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationInvitationsAndMembershipRequests(mock.Anything, testOrgId).
					Return(&models.Organization{
						Invitations: []models.OrganizationInvitation{
							{
								TargetEmail:        testEmail,
								InvitationStatuses: []models.OrganizationInvitationStatus{},
							},
						},
					}, nil)
			},
			want:    userMembershipStateInvited,
			wantErr: false,
		},
		{
			name:  "success - under review state",
			orgId: testOrgId,
			email: testEmail,
			organization: &models.Organization{
				MembershipRequests: []models.OrganizationMembershipRequest{
					{
						User:   *testUser,
						Status: models.OrgMembershipStatusPending,
					},
				},
			},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationInvitationsAndMembershipRequests(mock.Anything, testOrgId).
					Return(&models.Organization{
						MembershipRequests: []models.OrganizationMembershipRequest{
							{
								User:   *testUser,
								Status: models.OrgMembershipStatusPending,
							},
						},
					}, nil)
			},
			want:    userMembershipStateUnderReview,
			wantErr: false,
		},
		{
			name:  "store error",
			orgId: testOrgId,
			email: testEmail,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationInvitationsAndMembershipRequests(mock.Anything, testOrgId).
					Return(nil, errors.New("test error"))
			},
			want:          userMembershipStateNone,
			wantErr:       true,
			expectedError: "test error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockStore := mock_store.NewMockStore(t)
			logger := zap.NewNop()
			ctx := apicontext.AddLoggerToContext(context.Background(), logger)
			tt.mockSetup(mockStore)

			// Execute
			service := &organizationService{store: mockStore}
			got, err := service.getOranizationAccessStateByEmail(ctx, tt.orgId, tt.email)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, err.Error())
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOrganizationService_sendInvitationEmail(t *testing.T) {
	t.Parallel()

	testOrgId := uuid.New()
	testUserId := uuid.New()
	testOrg := &models.Organization{
		ID:   testOrgId,
		Name: "Test Org",
	}
	testUser := &models.User{
		ID:    testUserId,
		Name:  "Test User",
		Email: "test@example.com",
	}
	testInvitation := &models.OrganizationInvitation{
		InvitedBy:      testUserId,
		TargetEmail:    "invited@example.com",
		OrganizationID: testOrgId,
	}

	tests := []struct {
		name          string
		orgId         uuid.UUID
		invitation    *models.OrganizationInvitation
		mockSetup     func(*mock_store.MockStore, *mock_mailer.MockMailerService)
		wantErr       bool
		expectedError string
	}{
		{
			name:       "success",
			orgId:      testOrgId,
			invitation: testInvitation,
			mockSetup: func(m *mock_store.MockStore, mc *mock_mailer.MockMailerService) {
				m.EXPECT().
					GetOrganizationById(mock.Anything, testOrgId.String()).
					Return(testOrg, nil)
				m.EXPECT().
					GetUserById(mock.Anything, testUserId.String()).
					Return(testUser, nil)
				mc.EXPECT().
					SendInvitationEmail(mock.Anything, mailer.InvitationEmailData{
						OrganizationName:   testOrg.Name,
						RecipientEmail:     testInvitation.TargetEmail,
						InvitedByFirstName: testUser.Name,
						InvitationLink:     "https://app.zamp.ai",
					}).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "get organization error",
			orgId:      testOrgId,
			invitation: testInvitation,
			mockSetup: func(m *mock_store.MockStore, mc *mock_mailer.MockMailerService) {
				m.EXPECT().
					GetOrganizationById(mock.Anything, testOrgId.String()).
					Return(nil, errors.New("failed to get organization by id"))
			},
			wantErr:       true,
			expectedError: "failed to get organization by id",
		},
		{
			name:       "get user error",
			orgId:      testOrgId,
			invitation: testInvitation,
			mockSetup: func(m *mock_store.MockStore, mc *mock_mailer.MockMailerService) {
				m.EXPECT().
					GetOrganizationById(mock.Anything, testOrgId.String()).
					Return(testOrg, nil)
				m.EXPECT().
					GetUserById(mock.Anything, testUserId.String()).
					Return(nil, errors.New("failed to get user by id"))
			},
			wantErr:       true,
			expectedError: "failed to get user by id",
		},
		{
			name:       "send email error",
			orgId:      testOrgId,
			invitation: testInvitation,
			mockSetup: func(m *mock_store.MockStore, mc *mock_mailer.MockMailerService) {
				m.EXPECT().
					GetOrganizationById(mock.Anything, testOrgId.String()).
					Return(testOrg, nil)
				m.EXPECT().
					GetUserById(mock.Anything, testUserId.String()).
					Return(testUser, nil)
				mc.EXPECT().
					SendInvitationEmail(mock.Anything, mock.Anything).
					Return(errors.New("failed to send email"))
			},
			wantErr:       true,
			expectedError: "failed to send email",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockStore := mock_store.NewMockStore(t)
			mockMailer := mock_mailer.NewMockMailerService(t)
			logger := zap.NewNop()
			ctx := apicontext.AddLoggerToContext(context.Background(), logger)
			tt.mockSetup(mockStore, mockMailer)

			// Execute
			service := &organizationService{
				store:        mockStore,
				mailerClient: mockMailer,
			}
			err := service.sendInvitationEmail(ctx, tt.invitation)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, err.Error())
				}
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestOrganizationService_ValidateAudienceInOrganization(t *testing.T) {
	t.Parallel()

	testOrgId := uuid.New()
	testUserId := uuid.New()
	testTeamId := uuid.New()
	testPolicies := []models.ResourceAudiencePolicy{
		{
			ResourceID:           testOrgId,
			ResourceAudienceType: models.AudienceTypeUser,
			ResourceAudienceID:   testUserId,
		},
		{
			ResourceID:           testOrgId,
			ResourceAudienceType: models.AudienceTypeTeam,
			ResourceAudienceID:   testTeamId,
		},
	}

	tests := []struct {
		name          string
		orgId         uuid.UUID
		audienceType  models.AudienceType
		audienceId    uuid.UUID
		mockSetup     func(*mock_store.MockStore)
		wantErr       bool
		expectedError string
	}{
		{
			name:         "success - user audience found",
			orgId:        testOrgId,
			audienceType: models.AudienceTypeUser,
			audienceId:   testUserId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationPolicies(mock.Anything, testOrgId).
					Return(testPolicies, nil)
			},
			wantErr: false,
		},
		{
			name:         "success - team audience found",
			orgId:        testOrgId,
			audienceType: models.AudienceTypeTeam,
			audienceId:   testTeamId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationPolicies(mock.Anything, testOrgId).
					Return(testPolicies, nil)
			},
			wantErr: false,
		},
		{
			name:         "error - audience not found",
			orgId:        testOrgId,
			audienceType: models.AudienceTypeUser,
			audienceId:   uuid.New(),
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationPolicies(mock.Anything, testOrgId).
					Return(testPolicies, nil)
			},
			wantErr:       true,
			expectedError: "audience not found in organization",
		},
		{
			name:         "error - failed to get policies",
			orgId:        testOrgId,
			audienceType: models.AudienceTypeUser,
			audienceId:   testUserId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().
					GetOrganizationPolicies(mock.Anything, testOrgId).
					Return(nil, errors.New("db error"))
			},
			wantErr:       true,
			expectedError: "failed to find organization information",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockStore := mock_store.NewMockStore(t)
			logger := zap.NewNop()
			ctx := apicontext.AddLoggerToContext(context.Background(), logger)
			tt.mockSetup(mockStore)

			// Execute
			service := &organizationService{store: mockStore}
			err := service.ValidateAudienceInOrganization(ctx, tt.orgId, tt.audienceType, tt.audienceId)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, err.Error())
				}
				return
			}

			assert.NoError(t, err)
		})
	}
}

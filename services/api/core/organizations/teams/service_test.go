package teams

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
)

func TestNewTeamService(t *testing.T) {
	t.Parallel()

	mockStore := mock_store.NewMockStore(t)
	service := NewTeamService(mockStore)

	assert.NotNil(t, service)
}

func TestValidateContext(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	orgId := uuid.New()

	tests := []struct {
		name    string
		ctx     context.Context
		wantErr bool
	}{
		{
			name:    "success",
			ctx:     apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			wantErr: false,
		},
		{
			name:    "no user ID",
			ctx:     context.Background(),
			wantErr: true,
		},
		{
			name:    "no org IDs",
			ctx:     apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{}),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStore := mock_store.NewMockStore(t)
			service := &teamService{store: mockStore}

			_, _, err := service.validateContext(tt.ctx)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGetTeamsByOrganizationID(t *testing.T) {
	t.Parallel()

	orgId := uuid.New()
	teams := []models.Team{
		{TeamID: uuid.New(), OrganizationID: orgId, Name: "Team 1"},
		{TeamID: uuid.New(), OrganizationID: orgId, Name: "Team 2"},
	}

	tests := []struct {
		name      string
		orgId     uuid.UUID
		mockSetup func(*mock_store.MockStore)
		want      []models.Team
		wantErr   bool
	}{
		{
			name:  "success",
			orgId: orgId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetTeams(mock.Anything, orgId).Return(teams, nil)
			},
			want:    teams,
			wantErr: false,
		},
		{
			name:  "store error",
			orgId: orgId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetTeams(mock.Anything, orgId).Return(nil, errors.New("test error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStore := mock_store.NewMockStore(t)
			tt.mockSetup(mockStore)

			service := &teamService{store: mockStore}
			got, err := service.GetTeamsByOrganizationID(context.Background(), tt.orgId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetTeamById(t *testing.T) {
	t.Parallel()

	orgId := uuid.New()
	teamId := uuid.New()
	team := models.Team{TeamID: teamId, OrganizationID: orgId, Name: "Test Team"}

	tests := []struct {
		name      string
		orgId     uuid.UUID
		teamId    uuid.UUID
		mockSetup func(*mock_store.MockStore)
		want      *models.Team
		wantErr   bool
	}{
		{
			name:   "success",
			orgId:  orgId,
			teamId: teamId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetTeam(mock.Anything, orgId, teamId).Return(&team, nil)
			},
			want:    &team,
			wantErr: false,
		},
		{
			name:   "store error",
			orgId:  orgId,
			teamId: teamId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetTeam(mock.Anything, orgId, teamId).Return(nil, errors.New("test error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStore := mock_store.NewMockStore(t)
			tt.mockSetup(mockStore)

			service := &teamService{store: mockStore}
			got, err := service.GetTeamById(context.Background(), tt.orgId, tt.teamId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateTeam(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	orgId := uuid.New()
	team := models.Team{
		OrganizationID: orgId,
		Name:           "Test Team",
		Description:    "Test Description",
		CreatedBy:      userId,
		Metadata:       []byte(`{"color_hex_code":"#000000"}`),
	}

	tests := []struct {
		name      string
		ctx       context.Context
		orgId     uuid.UUID
		payload   CreateTeamPayload
		mockSetup func(*mock_store.MockStore)
		want      *models.Team
		wantErr   bool
	}{
		{
			name:  "success",
			ctx:   apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId: orgId,
			payload: CreateTeamPayload{
				Name:         "Test Team",
				Description:  "Test Description",
				ColorHexCode: "#000000",
			},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithTeamTransaction(mock.Anything, mock.AnythingOfType("func(store.TeamStore) error")).
					Run(func(_ context.Context, fn func(store.TeamStore) error) {
						mockTxStore := mock_store.NewMockStore(t)
						mockTxStore.EXPECT().GetTeamByName(mock.Anything, orgId, "Test Team").Return([]models.Team{}, nil)
						mockTxStore.EXPECT().CreateOrganizationTeam(mock.Anything, orgId, team).Return(&team, nil)
						fn(mockTxStore)
					}).Return(nil)
			},
			want:    &team,
			wantErr: false,
		},
		{
			name:  "team name already exists",
			ctx:   apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId: orgId,
			payload: CreateTeamPayload{
				Name:         "Test Team",
				Description:  "Test Description",
				ColorHexCode: "#000000",
			},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithTeamTransaction(mock.Anything, mock.AnythingOfType("func(store.TeamStore) error")).
					Run(func(_ context.Context, fn func(store.TeamStore) error) {
						mockTxStore := mock_store.NewMockStore(t)
						mockTxStore.EXPECT().GetTeamByName(mock.Anything, orgId, "Test Team").Return([]models.Team{team}, nil)
						fn(mockTxStore)
					}).Return(fmt.Errorf("Team with name Test Team already exists"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "unauthorized org",
			ctx:   apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{uuid.New()}),
			orgId: orgId,
			payload: CreateTeamPayload{
				Name:         "Test Team",
				Description:  "Test Description",
				ColorHexCode: "#000000",
			},
			mockSetup: func(m *mock_store.MockStore) {},
			want:      nil,
			wantErr:   true,
		},
		{
			name:  "missing name",
			ctx:   apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId: orgId,
			payload: CreateTeamPayload{
				Description:  "Test Description",
				ColorHexCode: "#000000",
			},
			mockSetup: func(m *mock_store.MockStore) {},
			want:      nil,
			wantErr:   true,
		},
		{
			name:  "name too long",
			ctx:   apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId: orgId,
			payload: CreateTeamPayload{
				Name: "This name is way too long and should fail validation",
			},
			mockSetup: func(m *mock_store.MockStore) {},
			want:      nil,
			wantErr:   true,
		},
		{
			name:  "description too long",
			ctx:   apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId: orgId,
			payload: CreateTeamPayload{
				Name:        "Test Team",
				Description: "This description is way too long and should fail validation. It should be less than 64 characters.",
			},
			mockSetup: func(m *mock_store.MockStore) {},
			want:      nil,
			wantErr:   true,
		},
		{
			name:  "invalid color hex code",
			ctx:   apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId: orgId,
			payload: CreateTeamPayload{
				Name:         "Test Team",
				Description:  "Test Description",
				ColorHexCode: "This is not a valid color hex code",
			},
			mockSetup: func(m *mock_store.MockStore) {},
			want:      nil,
			wantErr:   true,
		},
		{
			name:  "missing color",
			ctx:   apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId: orgId,
			payload: CreateTeamPayload{
				Name:        "Test Team",
				Description: "Test Description",
			},
			mockSetup: func(m *mock_store.MockStore) {},
			want:      nil,
			wantErr:   true,
		},
		{
			name:  "store error",
			ctx:   apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId: orgId,
			payload: CreateTeamPayload{
				Name:         "Test Team",
				Description:  "Test Description",
				ColorHexCode: "#000000",
			},
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithTeamTransaction(mock.Anything, mock.AnythingOfType("func(store.TeamStore) error")).
					Run(func(_ context.Context, fn func(store.TeamStore) error) {
						mockTxStore := mock_store.NewMockStore(t)
						mockTxStore.EXPECT().GetTeamByName(mock.Anything, orgId, "Test Team").Return([]models.Team{}, nil)
						mockTxStore.EXPECT().CreateOrganizationTeam(mock.Anything, orgId, team).Return(nil, errors.New("test error"))
						fn(mockTxStore)
					}).Return(fmt.Errorf("error creating team; please try again later"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStore := mock_store.NewMockStore(t)
			tt.mockSetup(mockStore)

			service := &teamService{store: mockStore}
			got, err := service.CreateTeam(tt.ctx, tt.orgId, tt.payload)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeleteTeam(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	orgId := uuid.New()
	teamId := uuid.New()

	tests := []struct {
		name      string
		ctx       context.Context
		orgId     uuid.UUID
		teamId    uuid.UUID
		mockSetup func(*mock_store.MockStore)
		wantErr   bool
	}{
		{
			name:   "success",
			ctx:    apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId:  orgId,
			teamId: teamId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().DeleteTeam(mock.Anything, orgId, teamId).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "unauthorized org",
			ctx:       apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{uuid.New()}),
			orgId:     orgId,
			teamId:    teamId,
			mockSetup: func(m *mock_store.MockStore) {},
			wantErr:   true,
		},
		{
			name:      "nil team ID",
			ctx:       apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId:     orgId,
			teamId:    uuid.Nil,
			mockSetup: func(m *mock_store.MockStore) {},
			wantErr:   true,
		},
		{
			name:   "store error",
			ctx:    apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId:  orgId,
			teamId: teamId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().DeleteTeam(mock.Anything, orgId, teamId).Return(errors.New("test error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStore := mock_store.NewMockStore(t)
			tt.mockSetup(mockStore)

			service := &teamService{store: mockStore}
			err := service.DeleteTeam(tt.ctx, tt.orgId, tt.teamId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestRenameTeam(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	orgId := uuid.New()
	teamId := uuid.New()

	tests := []struct {
		name      string
		ctx       context.Context
		orgId     uuid.UUID
		teamId    uuid.UUID
		payload   RenameTeamPayload
		mockSetup func(*mock_store.MockStore)
		wantErr   bool
	}{
		{
			name: "success",
			ctx:  apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			payload: RenameTeamPayload{
				Name: "New Name",
			},
			teamId: teamId,
			orgId:  orgId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().UpdateTeam(mock.Anything, orgId, models.Team{TeamID: teamId, OrganizationID: orgId, Name: "New Name"}).Return(&models.Team{}, nil)
			},
			wantErr: false,
		},
		{
			name: "unauthorized org",
			ctx:  apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{uuid.New()}),
			payload: RenameTeamPayload{
				Name: "New Name",
			},
			orgId:     orgId,
			teamId:    teamId,
			mockSetup: func(m *mock_store.MockStore) {},
			wantErr:   true,
		},
		{
			name: "nil team ID",
			ctx:  apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			payload: RenameTeamPayload{
				Name: "New Name",
			},
			orgId:     orgId,
			teamId:    uuid.Nil,
			mockSetup: func(m *mock_store.MockStore) {},
			wantErr:   true,
		},
		{
			name: "empty name",
			ctx:  apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			payload: RenameTeamPayload{
				Name: "",
			},
			orgId:     orgId,
			teamId:    teamId,
			mockSetup: func(m *mock_store.MockStore) {},
			wantErr:   true,
		},
		{
			name: "name too long",
			ctx:  apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			payload: RenameTeamPayload{
				Name: "This name is way too long and should fail validation",
			},
			orgId:     orgId,
			teamId:    teamId,
			mockSetup: func(m *mock_store.MockStore) {},
			wantErr:   true,
		},
		{
			name: "store error",
			ctx:  apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			payload: RenameTeamPayload{
				Name: "New Name",
			},
			orgId:  orgId,
			teamId: teamId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().UpdateTeam(mock.Anything, orgId, models.Team{TeamID: teamId, OrganizationID: orgId, Name: "New Name"}).Return(nil, errors.New("test error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStore := mock_store.NewMockStore(t)
			tt.mockSetup(mockStore)

			service := &teamService{store: mockStore}
			err := service.RenameTeam(tt.ctx, tt.orgId, tt.teamId, tt.payload)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestAddUserToTeam(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	orgId := uuid.New()
	teamId := uuid.New()
	membership := &models.TeamMembership{
		TeamID: teamId,
		UserID: userId,
	}
	team := &models.Team{
		TeamID:          teamId,
		OrganizationID:  orgId,
		TeamMemberships: []models.TeamMembership{},
	}
	policy := &models.ResourceAudiencePolicy{}

	tests := []struct {
		name      string
		ctx       context.Context
		orgId     uuid.UUID
		teamId    uuid.UUID
		userId    uuid.UUID
		mockSetup func(*mock_store.MockStore)
		want      *models.TeamMembership
		wantErr   bool
	}{
		{
			name:   "success",
			ctx:    apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId:  orgId,
			teamId: teamId,
			userId: userId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetOrganizationPolicyByUser(mock.Anything, orgId, userId).Return(policy, nil)
				m.EXPECT().WithTeamTransaction(mock.Anything, mock.AnythingOfType("func(store.TeamStore) error")).
					Run(func(ctx context.Context, fn func(store.TeamStore) error) {
						fn(m)
					}).Return(nil)
				m.EXPECT().GetTeam(mock.Anything, orgId, teamId).Return(team, nil)
				m.EXPECT().CreateTeamMembership(mock.Anything, teamId, userId).Return(membership, nil)
			},
			want:    membership,
			wantErr: false,
		},
		{
			name:      "unauthorized org",
			ctx:       apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{uuid.New()}),
			orgId:     orgId,
			teamId:    teamId,
			userId:    userId,
			mockSetup: func(m *mock_store.MockStore) {},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "nil team ID",
			ctx:       apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId:     orgId,
			teamId:    uuid.Nil,
			userId:    userId,
			mockSetup: func(m *mock_store.MockStore) {},
			want:      nil,
			wantErr:   true,
		},
		{
			name:   "user not in org",
			ctx:    apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId:  orgId,
			teamId: teamId,
			userId: userId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetOrganizationPolicyByUser(mock.Anything, orgId, userId).Return(nil, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "user already in team",
			ctx:    apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId:  orgId,
			teamId: teamId,
			userId: userId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetOrganizationPolicyByUser(mock.Anything, orgId, userId).Return(policy, nil)
				m.EXPECT().WithTeamTransaction(mock.Anything, mock.AnythingOfType("func(store.TeamStore) error")).
					Run(func(ctx context.Context, fn func(store.TeamStore) error) {
						fn(m)
					}).Return(errors.New("user is already a member of this team"))
				m.EXPECT().GetTeam(mock.Anything, orgId, teamId).Return(&models.Team{
					TeamID:          teamId,
					OrganizationID:  orgId,
					TeamMemberships: []models.TeamMembership{{UserID: userId}},
				}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "store error",
			ctx:    apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId:  orgId,
			teamId: teamId,
			userId: userId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().GetOrganizationPolicyByUser(mock.Anything, orgId, userId).Return(nil, errors.New("test error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStore := mock_store.NewMockStore(t)
			tt.mockSetup(mockStore)

			service := &teamService{store: mockStore}
			got, err := service.AddUserToTeam(tt.ctx, tt.orgId, tt.teamId, tt.userId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRemoveUserFromTeam(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	orgId := uuid.New()
	teamId := uuid.New()
	membershipId := uuid.New()

	tests := []struct {
		name         string
		ctx          context.Context
		orgId        uuid.UUID
		teamId       uuid.UUID
		membershipId uuid.UUID
		mockSetup    func(*mock_store.MockStore)
		wantErr      bool
	}{
		{
			name:         "success",
			ctx:          apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId:        orgId,
			teamId:       teamId,
			membershipId: membershipId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithTeamTransaction(mock.Anything, mock.AnythingOfType("func(store.TeamStore) error")).
					Run(func(ctx context.Context, fn func(store.TeamStore) error) {
						fn(m)
					}).Return(nil)
				m.EXPECT().GetTeam(mock.Anything, orgId, teamId).Return(&models.Team{
					TeamID:          teamId,
					OrganizationID:  orgId,
					TeamMemberships: []models.TeamMembership{{TeamMembershipID: membershipId}},
				}, nil)
				m.EXPECT().DeleteTeamMembership(mock.Anything, teamId, membershipId).Return(nil)
			},
			wantErr: false,
		},
		{
			name:         "unauthorized org",
			ctx:          apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{uuid.New()}),
			orgId:        orgId,
			teamId:       teamId,
			membershipId: membershipId,
			mockSetup:    func(m *mock_store.MockStore) {},
			wantErr:      true,
		},
		{
			name:         "nil team ID",
			ctx:          apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId:        orgId,
			teamId:       uuid.Nil,
			membershipId: membershipId,
			mockSetup:    func(m *mock_store.MockStore) {},
			wantErr:      true,
		},
		{
			name:         "nil membership ID",
			ctx:          apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId:        orgId,
			teamId:       teamId,
			membershipId: uuid.Nil,
			mockSetup:    func(m *mock_store.MockStore) {},
			wantErr:      true,
		},
		{
			name:         "membership not found",
			ctx:          apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId:        orgId,
			teamId:       teamId,
			membershipId: membershipId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithTeamTransaction(mock.Anything, mock.AnythingOfType("func(store.TeamStore) error")).
					Run(func(ctx context.Context, fn func(store.TeamStore) error) {
						fn(m)
					}).Return(fmt.Errorf("team membership not found"))
				m.EXPECT().GetTeam(mock.Anything, orgId, teamId).Return(&models.Team{
					TeamID:          teamId,
					OrganizationID:  orgId,
					TeamMemberships: []models.TeamMembership{},
				}, nil)
			},
			wantErr: true,
		},
		{
			name:         "get team error",
			ctx:          apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId:        orgId,
			teamId:       teamId,
			membershipId: membershipId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithTeamTransaction(mock.Anything, mock.AnythingOfType("func(store.TeamStore) error")).
					Run(func(ctx context.Context, fn func(store.TeamStore) error) {
						fn(m)
					}).Return(errors.New("test error"))
				m.EXPECT().GetTeam(mock.Anything, orgId, teamId).Return(nil, errors.New("test error"))
			},
			wantErr: true,
		},
		{
			name:         "delete membership error",
			ctx:          apicontext.AddAuthToContext(context.Background(), "user", userId, []uuid.UUID{orgId}),
			orgId:        orgId,
			teamId:       teamId,
			membershipId: membershipId,
			mockSetup: func(m *mock_store.MockStore) {
				m.EXPECT().WithTeamTransaction(mock.Anything, mock.AnythingOfType("func(store.TeamStore) error")).
					Run(func(ctx context.Context, fn func(store.TeamStore) error) {
						fn(m)
					}).Return(errors.New("test error"))
				m.EXPECT().GetTeam(mock.Anything, orgId, teamId).Return(&models.Team{
					TeamID:          teamId,
					OrganizationID:  orgId,
					TeamMemberships: []models.TeamMembership{{TeamMembershipID: membershipId}},
				}, nil)
				m.EXPECT().DeleteTeamMembership(mock.Anything, teamId, membershipId).Return(errors.New("test error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStore := mock_store.NewMockStore(t)
			tt.mockSetup(mockStore)

			service := &teamService{store: mockStore}
			err := service.RemoveUserFromTeam(tt.ctx, tt.orgId, tt.teamId, tt.membershipId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/store"
	helper "github.com/Zampfi/application-platform/services/api/helper"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	mockAuth "github.com/Zampfi/application-platform/services/api/mocks/core/auth"
	mockStore "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	mock_store "github.com/Zampfi/application-platform/services/api/mocks/db/store"
	mockKratos "github.com/Zampfi/application-platform/services/api/mocks/pkg/kratosclient"
	mock_kratosclient "github.com/Zampfi/application-platform/services/api/mocks/pkg/kratosclient"
	"github.com/Zampfi/application-platform/services/api/pkg/kratosclient"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	kratos "github.com/ory/kratos-client-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewAuthService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		adminSecrets []string
		authURL      string
		wantErr      bool
	}{
		{
			name:         "success",
			adminSecrets: []string{"secret1", "secret2"},
			authURL:      "http://valid-url.com",
			wantErr:      false,
		},
		{
			name:         "invalid auth URL",
			adminSecrets: []string{"secret1"},
			authURL:      "invalid-url",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockStore := mockStore.NewMockStore(t)

			// Execute
			service, err := NewAuthService(tt.adminSecrets, tt.authURL, mockStore, "local")

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, service)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, service)
		})
	}
}

func TestResolveSessionCookie(t *testing.T) {
	t.Parallel()

	testSession := &kratos.Session{
		Id: "test-session",
	}

	tests := []struct {
		name       string
		cookie     string
		mockSetup  func(*mockKratos.MockKratosClient)
		wantResult *kratos.Session
		wantErr    bool
	}{
		{
			name:   "success",
			cookie: "valid-cookie",
			mockSetup: func(m *mockKratos.MockKratosClient) {
				m.EXPECT().GetSessionInfo(mock.Anything, mock.Anything, "valid-cookie").
					Return(testSession, &http.Response{StatusCode: http.StatusOK}, nil)
			},
			wantResult: testSession,
			wantErr:    false,
		},
		{
			name:   "invalid session",
			cookie: "invalid-cookie",
			mockSetup: func(m *mockKratos.MockKratosClient) {
				m.EXPECT().GetSessionInfo(mock.Anything, mock.Anything, "invalid-cookie").
					Return(nil, &http.Response{StatusCode: http.StatusUnauthorized}, &kratosclient.Error{
						Code:    http.StatusUnauthorized,
						Message: "invalid session",
					})
			},
			wantResult: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockKratosClient := mockKratos.NewMockKratosClient(t)
			tt.mockSetup(mockKratosClient)

			service := &authService{
				kratosClient: mockKratosClient,
				environment:  "local",
			}

			// Execute
			session, resp, err := service.ResolveSessionCookie(
				context.Background(),
				tt.cookie,
			)

			// Assert
			if tt.wantErr {
				assert.Nil(t, session)
				assert.NotNil(t, err)
				return
			}

			assert.Equal(t, tt.wantResult, session)
			assert.NotNil(t, resp)
		})
	}
}

func TestResolveAdminInfo(t *testing.T) {
	t.Parallel()

	validUserID := uuid.New()
	validOrgID := uuid.New()

	tests := []struct {
		name         string
		setupHeaders func() http.Header
		adminSecrets []string
		wantRole     string
		wantUserID   uuid.UUID
		wantOrgIDs   []uuid.UUID
	}{
		{
			name: "valid admin with user and org",
			setupHeaders: func() http.Header {
				h := http.Header{}
				h.Set(helper.ADMIN_SECRET_HEADER, "valid-secret")
				h.Set(helper.PROXY_USER_ID_HEADER, validUserID.String())
				h.Add(helper.PROXY_WORKSPACE_IDS_HEADER, validOrgID.String())
				return h
			},
			adminSecrets: []string{"valid-secret"},
			wantRole:     "admin",
			wantUserID:   validUserID,
			wantOrgIDs:   []uuid.UUID{validOrgID},
		},
		{
			name: "invalid admin secret",
			setupHeaders: func() http.Header {
				h := http.Header{}
				h.Set(helper.ADMIN_SECRET_HEADER, "invalid-secret")
				return h
			},
			adminSecrets: []string{"valid-secret"},
			wantRole:     "anonymous",
			wantUserID:   uuid.Nil,
			wantOrgIDs:   []uuid.UUID{},
		},
		{
			name: "missing user ID",
			setupHeaders: func() http.Header {
				h := http.Header{}
				h.Set(helper.ADMIN_SECRET_HEADER, "valid-secret")
				return h
			},
			adminSecrets: []string{"valid-secret"},
			wantRole:     "anonymous",
			wantUserID:   uuid.Nil,
			wantOrgIDs:   []uuid.UUID{},
		},
		{
			name: "invalid user ID format",
			setupHeaders: func() http.Header {
				h := http.Header{}
				h.Set(helper.ADMIN_SECRET_HEADER, "valid-secret")
				h.Set(helper.PROXY_USER_ID_HEADER, "invalid-uuid")
				return h
			},
			adminSecrets: []string{"valid-secret"},
			wantRole:     "anonymous",
			wantUserID:   uuid.Nil,
			wantOrgIDs:   []uuid.UUID{},
		},
		{
			name: "missing organization IDs",
			setupHeaders: func() http.Header {
				h := http.Header{}
				h.Set(helper.ADMIN_SECRET_HEADER, "valid-secret")
				h.Set(helper.PROXY_USER_ID_HEADER, validUserID.String())
				return h
			},
			adminSecrets: []string{"valid-secret"},
			wantRole:     "anonymous",
			wantUserID:   uuid.Nil,
			wantOrgIDs:   []uuid.UUID{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			service := &authService{
				adminSecrets: tt.adminSecrets,
				environment:  "local",
			}

			// Execute
			role, userID, orgIDs := service.ResolveAdminInfo(
				context.Background(),
				tt.setupHeaders(),
			)

			// Assert
			assert.Equal(t, tt.wantRole, role)
			assert.Equal(t, tt.wantUserID, userID)
			assert.Equal(t, tt.wantOrgIDs, orgIDs)
		})
	}
}

func TestGetOrganizations(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	orgID := uuid.New()
	testOrgs := []models.Organization{{ID: orgID, Name: "Test Org"}}

	tests := []struct {
		name      string
		want      []models.Organization
		mockSetup func(*mockAuth.MockAuthServiceStore)
		wantErr   bool
	}{
		{
			name: "success",
			want: testOrgs,
			mockSetup: func(m *mockAuth.MockAuthServiceStore) {
				m.EXPECT().GetOrganizationsByMemberId(mock.Anything, userId).Return(testOrgs, nil)
			},
			wantErr: false,
		},
		{
			name: "store error",
			want: nil,
			mockSetup: func(m *mockAuth.MockAuthServiceStore) {
				m.EXPECT().GetOrganizationsByMemberId(mock.Anything, userId).Return(nil, errors.New("store error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			mockStore := mockAuth.NewMockAuthServiceStore(t)
			tt.mockSetup(mockStore)

			service := &authService{
				authServiceStore: mockStore,
				environment:      "local",
			}

			// Execute
			got, err := service.GetUserOrganizations(context.Background(), userId)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetAuthFlowForUser(t *testing.T) {
	t.Parallel()
	testEmail := "test@example.com"
	testDomain := "example.com"
	testProviderId := "google"
	testError := errors.New("test error")

	testSSOConfig := &models.OrganizationSSOConfig{
		SSOProviderID: testProviderId,
	}

	tests := []struct {
		name      string
		email     string
		mockSetup func(*mockAuth.MockAuthServiceStore, *mockKratos.MockKratosClient)
		want      *kratos.LoginFlow
		env       string
		wantErr   bool
	}{
		{
			name:  "success - get only oidc node on production",
			env:   "production",
			email: testEmail,
			mockSetup: func(m *mockAuth.MockAuthServiceStore, k *mockKratos.MockKratosClient) {

				testLoginFlow := &kratos.LoginFlow{
					Ui: kratos.UiContainer{
						Nodes: []kratos.UiNode{
							{
								Group: "oidc",
								Attributes: kratos.UiNodeAttributes{
									UiNodeInputAttributes: &kratos.UiNodeInputAttributes{
										Name:  "provider",
										Value: testProviderId,
									},
								},
							},
							{
								Group: "password",
								Attributes: kratos.UiNodeAttributes{
									UiNodeInputAttributes: &kratos.UiNodeInputAttributes{
										Name:  "password",
										Value: "testpassword",
									},
								},
							},
						},
					},
				}

				m.EXPECT().GetSSOConfigByDomain(mock.Anything, testDomain).Return(testSSOConfig, nil)
				k.EXPECT().CreateLoginFlow(mock.Anything, mock.Anything, testEmail).Return(testLoginFlow, nil, nil)
			},
			want: &kratos.LoginFlow{
				Ui: kratos.UiContainer{
					Nodes: []kratos.UiNode{
						{
							Group: "oidc",
							Attributes: kratos.UiNodeAttributes{
								UiNodeInputAttributes: &kratos.UiNodeInputAttributes{
									Name:  "provider",
									Value: testProviderId,
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "success - get all nodes on local",
			env:   "local",
			email: testEmail,
			mockSetup: func(m *mockAuth.MockAuthServiceStore, k *mockKratos.MockKratosClient) {
				testLoginFlow := &kratos.LoginFlow{
					Ui: kratos.UiContainer{
						Nodes: []kratos.UiNode{
							{
								Group: "oidc",
								Attributes: kratos.UiNodeAttributes{
									UiNodeInputAttributes: &kratos.UiNodeInputAttributes{
										Name:  "provider",
										Value: testProviderId,
									},
								},
							},
							{
								Group: "password",
								Attributes: kratos.UiNodeAttributes{
									UiNodeInputAttributes: &kratos.UiNodeInputAttributes{
										Name:  "password",
										Value: "testpassword",
									},
								},
							},
						},
					},
				}
				m.EXPECT().GetSSOConfigByDomain(mock.Anything, testDomain).Return(testSSOConfig, nil)
				k.EXPECT().CreateLoginFlow(mock.Anything, mock.Anything, testEmail).Return(testLoginFlow, nil, nil)
			},
			want: &kratos.LoginFlow{
				Ui: kratos.UiContainer{
					Nodes: []kratos.UiNode{
						{
							Group: "oidc",
							Attributes: kratos.UiNodeAttributes{
								UiNodeInputAttributes: &kratos.UiNodeInputAttributes{
									Name:  "provider",
									Value: testProviderId,
								},
							},
						},
						{
							Group: "password",
							Attributes: kratos.UiNodeAttributes{
								UiNodeInputAttributes: &kratos.UiNodeInputAttributes{
									Name:  "password",
									Value: "testpassword",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "invalid email",
			email: "invalid-email",
			env:   "production",
			mockSetup: func(m *mockAuth.MockAuthServiceStore, k *mockKratos.MockKratosClient) {
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "no sso config",
			env:   "production",
			email: testEmail,
			mockSetup: func(m *mockAuth.MockAuthServiceStore, k *mockKratos.MockKratosClient) {
				m.EXPECT().GetSSOConfigByDomain(mock.Anything, testDomain).Return(nil, testError)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "kratos error",
			email: testEmail,
			env:   "production",
			mockSetup: func(m *mockAuth.MockAuthServiceStore, k *mockKratos.MockKratosClient) {
				m.EXPECT().GetSSOConfigByDomain(mock.Anything, testDomain).Return(testSSOConfig, nil)
				k.EXPECT().CreateLoginFlow(mock.Anything, mock.Anything, testEmail).Return(nil, nil, &kratosclient.Error{Message: "kratos error"})
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
			mockStore := mockAuth.NewMockAuthServiceStore(t)
			mockKratosClient := mockKratos.NewMockKratosClient(t)
			tt.mockSetup(mockStore, mockKratosClient)

			service := &authService{
				authServiceStore: mockStore,
				kratosClient:     mockKratosClient,
				environment:      tt.env,
			}

			gin.SetMode(gin.TestMode)
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

			// Execute
			got, _, err := service.GetAuthFlowForUser(ctx, tt.email)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.want.Ui.Nodes), len(got.Ui.Nodes))
			assert.Equal(t, tt.want.Ui.Nodes[0].Attributes.UiNodeInputAttributes.Name, got.Ui.Nodes[0].Attributes.UiNodeInputAttributes.Name)
			assert.Equal(t, tt.want.Ui.Nodes[0].Attributes.UiNodeInputAttributes.Value, got.Ui.Nodes[0].Attributes.UiNodeInputAttributes.Value)

		})
	}
}

func TestAuthService_IsUserExposedKratosPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "starts with /sessions",
			path: "/sessions/123",
			want: true,
		},
		{
			name: "starts with /self-service/methods/oidc",
			path: "/self-service/methods/oidc/callback",
			want: true,
		},
		{
			name: "exact match /self-service/login",
			path: "/self-service/login",
			want: true,
		},
		{
			name: "exact match /self-service/logout/browser",
			path: "/self-service/logout/browser",
			want: true,
		},
		{
			name: "exact match /self-service/logout",
			path: "/self-service/logout",
			want: true,
		},
		{
			name: "no match - different path",
			path: "/api/v1/users",
			want: false,
		},
		{
			name: "no match - partial match of exact path",
			path: "/self-service/login/something",
			want: false,
		},
		{
			name: "no match - empty path",
			path: "",
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			service := &authService{}

			// Execute
			got := service.IsUserExposedKratosPath(tt.path)

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHandleNewUserCreated(t *testing.T) {
	t.Parallel()

	userId := uuid.New()
	adminSecret := "test-admin-secret"

	tests := []struct {
		name    string
		setup   func(t *testing.T) *authService
		userId  uuid.UUID
		wantErr bool
	}{
		{
			name: "success - user has invitation",
			setup: func(t *testing.T) *authService {
				mockStore := mock_store.NewMockStore(t)

				user := &models.User{
					ID:    userId,
					Email: "test@example.com",
				}

				invitation := models.OrganizationInvitation{
					OrganizationInvitationID: uuid.New(),
					OrganizationID:           uuid.New(),
					TargetEmail:              "test@example.com",
					Privilege:                models.PrivilegeOrganizationSystemAdmin,
				}

				mockStore.EXPECT().GetUserById(mock.Anything, userId.String()).Return(user, nil)

				mockStore.On("WithOrganizationTransaction", mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(store.OrganizationStore) error)
						fn(mockStore)
					}).Return(nil)

				mockStore.EXPECT().GetOrganizationInvitationsAll(mock.Anything).Return([]models.OrganizationInvitation{invitation}, nil)
				mockStore.EXPECT().CreateOrganizationInvitationStatus(mock.Anything, invitation.OrganizationInvitationID, models.InvitationStatusAccepted).Return(&models.OrganizationInvitationStatus{}, nil)
				mockStore.EXPECT().CreateOrganizationPolicy(mock.Anything, invitation.OrganizationID, models.AudienceTypeUser, userId, invitation.Privilege).Return(&models.ResourceAudiencePolicy{}, nil)

				return &authService{
					authServiceStore: mockStore,
					adminSecrets:     []string{adminSecret},
				}
			},
			userId:  userId,
			wantErr: false,
		},
		{
			name:   "success - user has no invitation",
			userId: userId,
			setup: func(t *testing.T) *authService {
				mockStore := mock_store.NewMockStore(t)

				user := &models.User{
					ID:    userId,
					Email: "test@example.com",
				}

				ssoConfig := &models.OrganizationSSOConfig{
					OrganizationID: uuid.New(),
					EmailDomain:    "example.com",
				}

				mockStore.On("GetUserById", mock.Anything, userId.String()).Return(user, nil)
				mockStore.On("WithOrganizationTransaction", mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(store.OrganizationStore) error)
						fn(mockStore)
					}).Return(nil)
				mockStore.EXPECT().GetOrganizationInvitationsAll(mock.Anything).Return([]models.OrganizationInvitation{}, nil)
				mockStore.On("GetPrimarySSOConfigByDomain", mock.Anything, "example.com").Return(ssoConfig, nil)
				mockStore.On("CreateOrganizationMembershipRequest", mock.Anything, ssoConfig.OrganizationID, userId, models.OrgMembershipStatusPending).Return(&models.OrganizationMembershipRequest{ID: uuid.New()}, nil)

				return &authService{
					adminSecrets:     []string{adminSecret},
					authServiceStore: mockStore,
				}
			},
			wantErr: false,
		},
		{
			name: "error - failed to get user",
			setup: func(t *testing.T) *authService {
				mockStore := mock_store.NewMockStore(t)

				mockStore.On("GetUserById", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to get user"))

				return &authService{
					adminSecrets:     []string{adminSecret},
					authServiceStore: mockStore,
				}
			},
			userId:  userId,
			wantErr: true,
		},
		{
			name: "error - failed to get invitations",
			setup: func(t *testing.T) *authService {
				mockStore := mock_store.NewMockStore(t)

				user := &models.User{
					ID:    userId,
					Email: "test@example.com",
				}

				mockStore.EXPECT().GetUserById(mock.Anything, userId.String()).Return(user, nil)

				mockStore.On("WithOrganizationTransaction", mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(store.OrganizationStore) error)
						fn(mockStore)
					}).Return(fmt.Errorf("failed to get invitations"))

				mockStore.EXPECT().GetOrganizationInvitationsAll(mock.Anything).Return(nil, fmt.Errorf("failed to get invitations"))

				return &authService{
					adminSecrets:     []string{adminSecret},
					authServiceStore: mockStore,
				}
			},
			userId:  userId,
			wantErr: true,
		},
		{
			name: "error - failed to get sso config",
			setup: func(t *testing.T) *authService {

				user := &models.User{
					ID:    userId,
					Email: "test@example.com",
				}

				mockStore := mock_store.NewMockStore(t)

				mockStore.EXPECT().GetUserById(mock.Anything, userId.String()).Return(user, nil)

				mockStore.On("WithOrganizationTransaction", mock.Anything, mock.AnythingOfType("func(store.OrganizationStore) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(store.OrganizationStore) error)
						fn(mockStore)
					}).Return(fmt.Errorf("failed to get sso config"))

				mockStore.EXPECT().GetOrganizationInvitationsAll(mock.Anything).Return([]models.OrganizationInvitation{}, nil)
				mockStore.On("GetPrimarySSOConfigByDomain", mock.Anything, "example.com").Return(nil, fmt.Errorf("failed to get sso config"))

				return &authService{
					adminSecrets:     []string{adminSecret},
					authServiceStore: mockStore,
				}
			},
			userId:  userId,
			wantErr: true,
		},
		{
			name: "unauthorized - invalid admin secret",
			setup: func(t *testing.T) *authService {
				mockStore := mock_store.NewMockStore(t)

				return &authService{
					adminSecrets:     []string{"invalid-admin-secret"},
					authServiceStore: mockStore,
				}
			},
			userId:  userId,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			service := tt.setup(t)

			// Execute
			err := service.HandleNewUserCreated(context.Background(), adminSecret, tt.userId)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestSignupUserAsAdmin(t *testing.T) {
	t.Parallel()

	actorId := uuid.New()

	validUUID := uuid.New()
	validEmail := "test@example.com"
	validPassword := "test-password"

	tests := []struct {
		name    string
		setup   func(t *testing.T) *authService
		ctx     context.Context
		email   string
		pass    string
		want    *models.User
		wantErr bool
	}{
		{
			name: "success",
			setup: func(t *testing.T) *authService {
				mockKratos := mock_kratosclient.NewMockKratosClient(t)
				identity := &kratos.Identity{
					Id: validUUID.String(),
					Traits: map[string]interface{}{
						"email": validEmail,
						"name":  "Test User",
					},
				}

				mockKratos.EXPECT().SignupUserEmailPassword(mock.Anything, mock.Anything, validEmail, validPassword).
					Return(identity, nil, nil)

				return &authService{
					kratosClient: mockKratos,
				}
			},
			ctx:   apicontext.AddAuthToContext(context.Background(), "admin", actorId, []uuid.UUID{}),
			email: validEmail,
			pass:  validPassword,
			want: &models.User{
				ID:    validUUID,
				Email: validEmail,
				Name:  "Test User",
			},
			wantErr: false,
		},
		{
			name: "unauthorized - not admin",
			setup: func(t *testing.T) *authService {
				return &authService{}
			},
			ctx:     apicontext.AddAuthToContext(context.Background(), "user", actorId, []uuid.UUID{}),
			email:   validEmail,
			pass:    validPassword,
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - kratos signup fails",
			setup: func(t *testing.T) *authService {
				mockKratos := mock_kratosclient.NewMockKratosClient(t)
				mockKratos.EXPECT().SignupUserEmailPassword(mock.Anything, mock.Anything, validEmail, validPassword).
					Return(nil, &http.Response{StatusCode: http.StatusBadRequest}, &kratosclient.Error{
						Code:    http.StatusBadRequest,
						Message: "signup failed",
					})

				return &authService{
					kratosClient: mockKratos,
				}
			},
			ctx:     apicontext.AddAuthToContext(context.Background(), "admin", actorId, []uuid.UUID{}),
			email:   validEmail,
			pass:    validPassword,
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - invalid identity",
			setup: func(t *testing.T) *authService {
				mockKratos := mock_kratosclient.NewMockKratosClient(t)
				identity := &kratos.Identity{
					Id:     "invalid-uuid",
					Traits: nil,
				}

				mockKratos.EXPECT().SignupUserEmailPassword(mock.Anything, mock.Anything, validEmail, validPassword).
					Return(identity, nil, nil)

				return &authService{
					kratosClient: mockKratos,
				}
			},
			ctx:     apicontext.AddAuthToContext(context.Background(), "admin", actorId, []uuid.UUID{}),
			email:   validEmail,
			pass:    validPassword,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			service := tt.setup(t)

			// Execute
			got, err := service.SignupUserAsAdmin(tt.ctx, tt.email, tt.pass)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

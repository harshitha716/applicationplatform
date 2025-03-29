package mailer

import (
	"context"
	"fmt"
	"testing"

	mock_emailtemplates "github.com/Zampfi/application-platform/services/api/mocks/core/mailer/email_templates"
	mock_sparkpost "github.com/Zampfi/application-platform/services/api/mocks/pkg/sparkpost"
	"github.com/stretchr/testify/mock"
	"github.com/zeebo/assert"
)

func TestMailerService_SendInvitationEmail(t *testing.T) {
	// Test data
	testData := InvitationEmailData{
		RecipientEmail:     "test@example.com",
		InvitedByFirstName: "John",
		OrganizationName:   "Test Org",
		InvitationLink:     "https://test.com/invite",
	}

	// Expected template data
	expectedTemplateData := map[string]string{
		"invited_by_first_name": testData.InvitedByFirstName,
		"organization_name":     testData.OrganizationName,
		"invitation_link":       testData.InvitationLink,
	}

	tests := []struct {
		name          string
		templateErr   error
		templateResp  string
		mockSetup     func(templater *mock_emailtemplates.MockTemplater, sparkpostClient *mock_sparkpost.MockSparkPostClient)
		expectedError error
	}{
		{
			name:         "successful invitation email",
			templateResp: "<html>Test template</html>",
			mockSetup: func(templater *mock_emailtemplates.MockTemplater, sparkpostClient *mock_sparkpost.MockSparkPostClient) {
				templater.On("GetTemplate",
					"send_membership_invitation",
					"/templates",
					expectedTemplateData,
				).Return(mock.Anything, nil).Once()
				sparkpostClient.On("SendEmail",
					context.Background(),
					"from@test.com",
					fmt.Sprintf("%s invited you to join Zamp", testData.InvitedByFirstName),
					mock.Anything,
					[]string{testData.RecipientEmail},
				).Return(nil).Once()
			},
		},
		{
			name: "template error",
			mockSetup: func(templater *mock_emailtemplates.MockTemplater, sparkpostClient *mock_sparkpost.MockSparkPostClient) {
				templater.On("GetTemplate",
					"send_membership_invitation",
					"/templates",
					mock.Anything,
				).Return(mock.Anything, fmt.Errorf("template error")).Once()
			},
			templateErr:   fmt.Errorf("template error"),
			expectedError: fmt.Errorf("template error"),
		},
		{
			name: "sparkpost error",
			mockSetup: func(templater *mock_emailtemplates.MockTemplater, sparkpostClient *mock_sparkpost.MockSparkPostClient) {
				templater.On("GetTemplate",
					"send_membership_invitation",
					"/templates",
					mock.Anything,
				).Return(mock.Anything, nil).Once()

				sparkpostClient.On("SendEmail",
					context.Background(),
					"from@test.com",
					fmt.Sprintf("%s invited you to join Zamp", testData.InvitedByFirstName),
					mock.Anything,
					[]string{testData.RecipientEmail},
				).Return(fmt.Errorf("sparkpost error")).Once()
			},
			expectedError: fmt.Errorf("sparkpost error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockTemplater := mock_emailtemplates.NewMockTemplater(t)
			mockSparkpostClient := mock_sparkpost.NewMockSparkPostClient(t)

			// Setup mock expectations
			tt.mockSetup(mockTemplater, mockSparkpostClient)

			// Create service instance
			service := NewMailerService(mockSparkpostClient, "from@test.com", "/templates")
			// Override templater with mock
			service.(*mailerService).templater = mockTemplater

			// Call method
			err := service.SendInvitationEmail(context.Background(), testData)

			// Verify expectations
			mockTemplater.AssertExpectations(t)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

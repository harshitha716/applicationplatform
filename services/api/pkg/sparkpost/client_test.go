package sparkpost

import (
	"context"
	"fmt"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestClient_SendEmail tests the SendEmail method

func TestClient_SendEmail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		from            string
		subject         string
		htmlContent     string
		recipients      []string
		mockErr         error
		wantErr         bool
		wantErrType     error
		skipMock        bool // Skip mock setup for validation error cases
		wantErrContains string
		mockResp        *sp.Response
	}{
		{
			name:        "valid email",
			from:        "test@example.com",
			subject:     "Test Subject",
			htmlContent: "<p>Test content</p>",
			recipients:  []string{"recipient@example.com"},
			mockErr:     nil,
			wantErr:     false,
			mockResp: &sp.Response{
				Results: map[string]interface{}{
					"total_accepted_recipients": float64(1),
					"total_rejected_recipients": float64(0),
				},
			},
		},
		{
			name:            "empty recipients",
			from:            "test@example.com",
			subject:         "Test Subject",
			htmlContent:     "<p>Test content</p>",
			recipients:      []string{},
			wantErr:         true,
			wantErrType:     ErrInvalidRecipients,
			wantErrContains: "recipients list cannot be empty",
			skipMock:        true, // Skip mock since validation fails before API call
		},
		{
			name:        "empty sender",
			from:        "",
			subject:     "Test Subject",
			htmlContent: "<p>Test content</p>",
			recipients:  []string{"recipient@example.com"},
			wantErr:     true,
			wantErrType: ErrInvalidSender,
			skipMock:    true, // Skip mock since validation fails before API call
		},
		{
			name:            "empty subject",
			from:            "test@example.com",
			subject:         "",
			htmlContent:     "<p>Test content</p>",
			recipients:      []string{"recipient@example.com"},
			wantErr:         true,
			wantErrContains: "subject and content cannot be empty",
			skipMock:        true, // Skip mock since validation fails before API call
		},
		{
			name:            "empty content",
			from:            "test@example.com",
			subject:         "Test Subject",
			htmlContent:     "",
			recipients:      []string{"recipient@example.com"},
			wantErr:         true,
			wantErrContains: "subject and content cannot be empty",
			skipMock:        true, // Skip mock since validation fails before API call
		},
		{
			name:            "sparkpost client error",
			from:            "test@example.com",
			subject:         "Test Subject",
			htmlContent:     "<p>Test content</p>",
			recipients:      []string{"recipient@example.com"},
			mockErr:         fmt.Errorf("sparkpost error"),
			wantErr:         true,
			wantErrContains: "failed to send email: sparkpost error",
		},
		{
			name:        "multiple recipients",
			from:        "test@example.com",
			subject:     "Test Subject",
			htmlContent: "<p>Test content</p>",
			recipients:  []string{"recipient1@example.com", "recipient2@example.com"},
			mockErr:     nil,
			wantErr:     false,
			mockResp: &sp.Response{
				Results: map[string]interface{}{
					"total_accepted_recipients": float64(1),
					"total_rejected_recipients": float64(0),
				},
			},
		},
		{
			name:        "context cancelled",
			from:        "test@example.com",
			subject:     "Test Subject",
			htmlContent: "<p>Test content</p>",
			recipients:  []string{"recipient@example.com"},
			wantErr:     true,
			wantErrType: context.Canceled,
			skipMock:    true,
		},

		{
			name:            "nil recipients",
			from:            "test@example.com",
			subject:         "Test Subject",
			htmlContent:     "<p>Test content</p>",
			recipients:      nil,
			wantErr:         true,
			wantErrType:     ErrInvalidRecipients,
			wantErrContains: "recipients list cannot be empty",
			skipMock:        true,
		},
		{
			name:            "invalid config - empty API key",
			from:            "test@example.com",
			subject:         "Test Subject",
			htmlContent:     "<p>Test content</p>",
			recipients:      []string{"recipient@example.com"},
			wantErr:         true,
			wantErrType:     ErrInvalidConfig,
			wantErrContains: "invalid sparkpost configuration",
		},
		{
			name:            "invalid config - empty API URL",
			from:            "test@example.com",
			subject:         "Test Subject",
			htmlContent:     "<p>Test content</p>",
			recipients:      []string{"recipient@example.com"},
			wantErr:         true,
			wantErrType:     ErrInvalidConfig,
			wantErrContains: "invalid sparkpost configuration",
			skipMock:        true,
		},
		{
			name:        "successful send with multiple recipients",
			from:        "test@example.com",
			subject:     "Test Subject",
			htmlContent: "<p>Test content</p>",
			recipients:  []string{"recipient1@example.com", "recipient2@example.com", "recipient3@example.com"},
			mockErr:     nil,
			wantErr:     false,
			mockResp: &sp.Response{
				Results: map[string]interface{}{
					"total_accepted_recipients": float64(1),
					"total_rejected_recipients": float64(0),
				},
			},
		},
		{
			name:            "client init error",
			from:            "test@example.com",
			subject:         "Test Subject",
			htmlContent:     "<p>Test content</p>",
			recipients:      []string{"recipient@example.com"},
			mockErr:         fmt.Errorf("client init error"),
			wantErr:         true,
			wantErrContains: "failed to send email: client init error",
		},
		{
			name:        "successful send with html content",
			from:        "test@example.com",
			subject:     "Test Subject",
			htmlContent: "<p>Test content with <strong>HTML</strong></p>",
			recipients:  []string{"recipient@example.com"},
			mockErr:     nil,
			wantErr:     false,
			mockResp: &sp.Response{
				Results: map[string]interface{}{
					"total_accepted_recipients": float64(1),
					"total_rejected_recipients": float64(0),
				},
			},
		},
		{
			name:        "successful send with special characters",
			from:        "test@example.com",
			subject:     "Test Subject with ñ and é",
			htmlContent: "<p>Content with special chars: áéíóú</p>",
			recipients:  []string{"recipient@example.com"},
			mockErr:     nil,
			wantErr:     false,
			mockResp: &sp.Response{
				Results: map[string]interface{}{
					"total_accepted_recipients": float64(1),
					"total_rejected_recipients": float64(0),
				},
			},
		},
		{
			name:        "mock response error",
			from:        "test@example.com",
			subject:     "Test Subject",
			htmlContent: "<p>Test content</p>",
			recipients:  []string{"recipient@example.com"},
			mockResp: &sp.Response{
				Results: map[string]interface{}{
					"total_rejected_recipients": float64(1),
					"total_accepted_recipients": float64(0),
				},
			},
			mockErr:         fmt.Errorf("mock response error"),
			wantErr:         true,
			wantErrContains: "failed to send email: mock response error",
		},
		{
			name:        "http status error",
			from:        "test@example.com",
			subject:     "Test Subject",
			htmlContent: "<p>Test content</p>",
			recipients:  []string{"recipient@example.com"},
			mockResp: &sp.Response{
				Results: map[string]interface{}{
					"total_rejected_recipients": float64(1),
					"total_accepted_recipients": float64(0),
				},
			},
			wantErr:         true,
			wantErrType:     ErrNoRecipientsAccepted,
			wantErrContains: "no recipients accepted by SparkPost API",
		},
		{
			name:        "invalid response structure",
			from:        "test@example.com",
			subject:     "Test Subject",
			htmlContent: "<p>Test content</p>",
			recipients:  []string{"recipient@example.com"},
			mockResp: &sp.Response{
				Results: nil,
			},
			wantErr:         true,
			wantErrType:     ErrInvalidResponse,
			wantErrContains: "invalid response from SparkPost API",
		},
		{
			name:        "no recipients accepted",
			from:        "test@example.com",
			subject:     "Test Subject",
			htmlContent: "<p>Test content</p>",
			recipients:  []string{"recipient@example.com"},
			mockResp: &sp.Response{
				Results: map[string]interface{}{
					"total_accepted_recipients": float64(0),
					"total_rejected_recipients": float64(1),
				},
			},
			wantErr:         true,
			wantErrType:     ErrNoRecipientsAccepted,
			wantErrContains: "no recipients accepted by SparkPost API",
		},
		{
			name:            "sparkpost client error",
			from:            "test@example.com",
			subject:         "Test Subject",
			htmlContent:     "<p>Test content</p>",
			recipients:      []string{"recipient@example.com"},
			mockErr:         fmt.Errorf("sparkpost error"),
			wantErr:         true,
			wantErrContains: "failed to send email: sparkpost error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx context.Context
			switch tt.name {
			case "context cancelled":
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(context.Background())
				cancel()
			case "context timeout":
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(context.Background(), 1)
				cancel()
			default:
				ctx = context.Background()
			}

			var client SparkPostClient
			var err error

			switch tt.name {
			case "invalid config - empty API key", "invalid config - empty API URL":
				client, err = NewClient(Config{})
				require.Error(t, err)
				require.ErrorIs(t, err, ErrInvalidConfig)
				require.ErrorContains(t, err, tt.wantErrContains)
				return
			default:
				// Create a mock SparkPost client for non-config tests
				mockClient := &MockClient{}
				client = &Client{client: mockClient}

				// Set up expectations for the Send method
				if !tt.skipMock {
					mockClient.On("Send", mock.AnythingOfType("*gosparkpost.Transmission")).Return("test-id", tt.mockResp, tt.mockErr)
					defer mockClient.AssertExpectations(t)
				}
			}

			err = client.SendEmail(ctx, tt.from, tt.subject, tt.htmlContent, tt.recipients)
			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrType != nil {
					require.ErrorIs(t, err, tt.wantErrType)
				}
				if tt.wantErrContains != "" {
					require.ErrorContains(t, err, tt.wantErrContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_buildRecipients(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		emails   []string
		expected []sp.Recipient
	}{
		{
			name:   "single recipient",
			emails: []string{"test@example.com"},
			expected: []sp.Recipient{
				{Address: sp.Address{Email: "test@example.com"}},
			},
		},
		{
			name:   "multiple recipients",
			emails: []string{"test1@example.com", "test2@example.com"},
			expected: []sp.Recipient{
				{Address: sp.Address{Email: "test1@example.com"}},
				{Address: sp.Address{Email: "test2@example.com"}},
			},
		},
		{
			name:     "empty recipients",
			emails:   []string{},
			expected: []sp.Recipient{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildRecipients(tt.emails)
			assert.Equal(t, tt.expected, result)
		})
	}
}

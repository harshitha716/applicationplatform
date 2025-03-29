package sparkpost

import (
	"context"
	"fmt"

	sp "github.com/SparkPost/gosparkpost"
)

// mockSparkPostClient defines the interface for mocking SparkPost client
type mockSparkPostClient interface {
	Send(transmission *sp.Transmission) (string, *sp.Response, error)
}

type SparkPostClient interface {
	SendEmail(ctx context.Context, from, subject, htmlContent string, recipients []string) error
}

// Client represents a SparkPost client wrapper
type Client struct {
	client mockSparkPostClient
}

// Config holds SparkPost configuration
type Config struct {
	APIKey string
	APIUrl string
}

// NewClient creates a new SparkPost client
func NewClient(cfg Config) (SparkPostClient, error) {
	if cfg.APIKey == "" || cfg.APIUrl == "" {
		return nil, ErrInvalidConfig
	}
	client := &sp.Client{}
	err := client.Init(&sp.Config{
		BaseUrl:    cfg.APIUrl,
		ApiKey:     cfg.APIKey,
		ApiVersion: 1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sparkpost client: %w", err)
	}
	return &Client{client: client}, nil
}

// SendEmail sends an email using SparkPost
func (c *Client) SendEmail(ctx context.Context, from, subject, htmlContent string, recipients []string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if from == "" {
		return ErrInvalidSender
	}
	if subject == "" || htmlContent == "" {
		return fmt.Errorf("subject and content cannot be empty")
	}
	if len(recipients) == 0 {
		return ErrInvalidRecipients
	}

	tx := &sp.Transmission{
		Recipients: buildRecipients(recipients),
		Content: sp.Content{
			HTML:    htmlContent,
			From:    from,
			Subject: subject,
		},
	}

	_, resp, err := c.client.Send(tx)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Validate response
	if resp == nil || resp.Results == nil {
		return ErrInvalidResponse
	}

	results := resp.Results.(map[string]interface{})
	accepted, hasAccepted := results["total_accepted_recipients"].(float64)
	rejected, hasRejected := results["total_rejected_recipients"].(float64)

	if !hasAccepted || accepted == 0 {
		if hasRejected && rejected > 0 {
			return ErrNoRecipientsAccepted
		}
		return ErrInvalidResponse
	}

	if hasRejected && rejected > 0 {
		return fmt.Errorf("sparkpost API error: rejected=%v", rejected)
	}

	return nil
}

func buildRecipients(emails []string) []sp.Recipient {
	recipients := make([]sp.Recipient, len(emails))
	for i, email := range emails {
		recipients[i] = sp.Recipient{Address: sp.Address{Email: email}}
	}
	return recipients
}

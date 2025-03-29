package sparkpost

import "errors"

var (
	// ErrInvalidConfig is returned when the SparkPost configuration is invalid
	ErrInvalidConfig = errors.New("invalid sparkpost configuration")

	// ErrInvalidRecipients is returned when the recipients list is empty
	ErrInvalidRecipients = errors.New("recipients list cannot be empty")

	// ErrInvalidSender is returned when the sender email is invalid
	ErrInvalidSender = errors.New("invalid sender email")

	// ErrInvalidResponse is returned when the SparkPost API response is invalid
	ErrInvalidResponse = errors.New("invalid response from SparkPost API")

	// ErrNoRecipientsAccepted is returned when no recipients were accepted by SparkPost
	ErrNoRecipientsAccepted = errors.New("no recipients accepted by SparkPost API")
)

package errors

import "errors"

const (
	ErrCodeInvalidProviderName = "INVALID_PROVIDER_NAME"
	ErrMsgInvalidProviderName  = "Invalid provider name"
)

var (
	ErrInvalidProviderName = errors.New(ErrCodeInvalidProviderName)
)

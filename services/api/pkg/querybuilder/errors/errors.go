package errors

import "errors"

const (
	ErrInvalidDataTypeMessage             = "ERR_INVALID_DATA_TYPE"
	ErrInvalidOperatorMessage             = "ERR_INVALID_OPERATOR"
	ErrNoConditionsMessage                = "ERR_NO_CONDITIONS"
	ErrInvalidCustomDataTypeMessage       = "ERR_INVALID_CUSTOM_DATA_TYPE"
	ErrInvalidCustomDataTypeConfigMessage = "ERR_INVALID_CUSTOM_DATA_TYPE_CONFIG"
)

var (
	ErrInvalidDataType             = errors.New(ErrInvalidDataTypeMessage)
	ErrInvalidOperator             = errors.New(ErrInvalidOperatorMessage)
	ErrNoConditions                = errors.New(ErrNoConditionsMessage)
	ErrInvalidCustomDataType       = errors.New(ErrInvalidCustomDataTypeMessage)
	ErrInvalidCustomDataTypeConfig = errors.New(ErrInvalidCustomDataTypeConfigMessage)
)

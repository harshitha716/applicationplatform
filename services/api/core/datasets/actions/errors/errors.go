package errors

import "errors"

const (
	ErrInvalidDatasetActionStatusMessage = "ERR_INVALID_DATASET_ACTION_STATUS"
	ErrInvalidDatasetActionTypeMessage   = "ERR_INVALID_DATASET_ACTION_TYPE"
)

var (
	ErrInvalidDatasetActionStatus = errors.New(ErrInvalidDatasetActionStatusMessage)
	ErrInvalidDatasetActionType   = errors.New(ErrInvalidDatasetActionTypeMessage)
)

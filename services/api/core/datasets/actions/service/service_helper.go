package service

import (
	"slices"

	"github.com/Zampfi/application-platform/services/api/core/datasets/actions/constants"
)

// TODO to add more validations later on
func (s *datasetActionService) validateDatasetActionStatus(status string) bool {
	return slices.Contains(constants.ValidDatasetActionStatuses, status)
}

func (s *datasetActionService) validateDatasetActionType(actionType string) bool {
	return slices.Contains(constants.ValidDatasetActionTypes, actionType)
}

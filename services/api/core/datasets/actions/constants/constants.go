package constants

import dataplatfromactionconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/constants"

type ActionType string

const (
	ActionTypeDatasetExport     ActionType = "DATASET_EXPORT"
	ActionTypeDatasetFileImport ActionType = "DATASET_FILE_IMPORT"
)

var ValidDatasetActionStatuses = []string{
	string(dataplatfromactionconstants.ActionStatusInitiated),
	string(dataplatfromactionconstants.ActionStatusSuccessful),
	string(dataplatfromactionconstants.ActionStatusFailed),
}

var ValidDatasetActionTypes = []string{
	string(dataplatfromactionconstants.ActionTypeUpdateDatasetData),
	string(dataplatfromactionconstants.ActionTypeRegisterDataset),
	string(dataplatfromactionconstants.ActionTypeRegisterJob),
	string(dataplatfromactionconstants.ActionTypeUpsertTemplate),
	string(dataplatfromactionconstants.ActionTypeUpdateDataset),
	string(ActionTypeDatasetExport),
}

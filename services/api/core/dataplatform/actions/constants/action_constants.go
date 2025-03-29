package constants

type ActionType string

const (
	ActionTypeCreateMV          ActionType = "CREATE_MV"
	ActionTypeDeleteMV          ActionType = "DELETE_MV"
	ActionTypeUpdateDatasetData ActionType = "UPDATE_DATASET_DATA"
	ActionTypeRegisterDataset   ActionType = "REGISTER_DATASET"
	ActionTypeRegisterJob       ActionType = "REGISTER_JOB"
	ActionTypeUpsertTemplate    ActionType = "UPSERT_TEMPLATE"
	ActionTypeUpdateDataset     ActionType = "UPDATE_DATASET"
	ActionTypeCopyDataset       ActionType = "COPY_DATASET"
	ActionTypeDeleteDataset     ActionType = "DELETE_DATASET"
)

type ActionStatus string

const (
	ActionStatusInitiated  ActionStatus = "INITIATED"
	ActionStatusSuccessful ActionStatus = "SUCCESSFUL"
	ActionStatusFailed     ActionStatus = "FAILED"
)

var ActionTerminationStatuses = []ActionStatus{
	ActionStatusSuccessful,
	ActionStatusFailed,
}

type ActionActor string

const (
	ActionActorUser   ActionActor = "USER"
	ActionActorSystem ActionActor = "SYSTEM"
)

const (
	MerchantIdNotebookParam  = "merchant_id"
	QueryNotebookParam       = "query"
	UpdatedColumnParam       = "updated_column"
	UpdatedValuesParam       = "updated_values"
	UpdateDatasetDataParams  = "update_dataset_data_params"
	DatasetIdParam           = "dataset_id"
	DataPlatformModulesSrc   = "dataplatform_modules_src"
	RegisterDatasetParams    = "register_dataset_params"
	RegisterJobParams        = "register_job_params"
	UpsertTemplateParams     = "upsert_template_params"
	UpdateDatasetParams      = "update_dataset_params"
	UpdateDatasetEventParams = "update_dataset_event"
	CreateMVParams           = "create_mv_params"
	CopyDatasetParams        = "copy_dataset_params"
)

var SubmitOneTimeJobActions = []ActionType{
	ActionTypeCreateMV,
	ActionTypeDeleteMV,
	ActionTypeRegisterDataset,
	ActionTypeRegisterJob,
	ActionTypeUpsertTemplate,
	ActionTypeUpdateDataset,
	ActionTypeCopyDataset,
	ActionTypeDeleteDataset,
}

var SubmitJobActions = []ActionType{
	ActionTypeUpdateDatasetData,
}

type UpdateDatasetEventType string

const (
	UpdateDatasetEventTypeUpdateCustomColumn UpdateDatasetEventType = "update_custom_column"
	UpdateDatasetEventTypeUpsertRules        UpdateDatasetEventType = "upsert_rules"
)

type UpsertRuleOperation string

const (
	UpsertRuleOperationCreate  UpsertRuleOperation = "create"
	UpsertRuleOperationReorder UpsertRuleOperation = "reorder"
)

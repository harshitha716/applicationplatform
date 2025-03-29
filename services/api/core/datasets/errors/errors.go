package errors

import "errors"

const (
	ErrFailedToBuildQueryMessage                 = "ERR_FAILED_TO_BUILD_QUERY"
	ErrFailedToGetDataMessage                    = "ERR_FAILED_TO_GET_DATA"
	ErrNoRowMessage                              = "ERR_NO_ROW"
	ErrMoreThanOneRowMessage                     = "ERR_MORE_THAN_ONE_ROW"
	ErrInvalidMetadataFormatMessage              = "ERR_INVALID_METADATA_FORMAT"
	ErrFailedToUnmarshalMetadataMessage          = "ERR_FAILED_TO_UNMARSHAL_METADATA"
	ErrFailedToGetDatasetMessage                 = "ERR_FAILED_TO_GET_DATASET_METADATA"
	ErrFailedToGetDatasetByIdMessage             = "ERR_FAILED_TO_GET_DATASET_BY_ID"
	ErrFailedToPopulateMultiSelectOptionsMessage = "ERR_FAILED_TO_POPULATE_MULTI_SELECT_OPTIONS"
	ErrFailedToRegisterDatasetMessage            = "ERR_FAILED_TO_REGISTER_DATASET"
	ErrFailedToRegisterDatasetJobMessage         = "ERR_FAILED_TO_REGISTER_DATASET_JOB"
	ErrFailedToMarshalMetadataMessage            = "ERR_FAILED_TO_MARSHAL_METADATA"
	ErrFailedToUpdateDatasetMessage              = "ERR_FAILED_TO_UPDATE_DATASET"
	ErrFailedToUpsertTemplateMessage             = "ERR_FAILED_TO_UPSERT_TEMPLATE"
	ErrFailedToGetDatasetColumnDatatypeMessage   = "ERR_FAILED_TO_GET_DATASET_COLUMN_DATATYPE"
	ErrColumnMetadataTypeIsEmptyMessage          = "ERR_COLUMN_METADATA_TYPE_IS_EMPTY"
	ErrFailedToParseDatasetIdMessage             = "ERR_FAILED_TO_PARSE_DATASET_ID"
	ErrFailedToGetDatasetMetadataMessage         = "ERR_FAILED_TO_GET_DATASET_METADATA"
	ErrInvalidDataplatformProviderMessage        = "ERR_INVALID_DATAPLATFORM_PROVIDER"
	ErrInvalidDatalistinSortColumnMessage        = "ERR_INVALID_DATALISTIN_SORT_COLUMN"
	ErrFailedToUpdateRulePriorityMessage         = "ERR_FAILED_TO_UPDATE_RULE_PRIORITY"
	ErrNoRulesPresentMessage                     = "ERR_NO_RULES_PRESENT"
	ErrFailedToGetRuleMessage                    = "ERR_FAILED_TO_GET_RULES"
	ErrFailedToUpdateDatasetActionMessage        = "ERR_FAILED_TO_UPDATE_DATASET_ACTION"
	ErrInvalidDatasetTypeMessage                 = "ERR_INVALID_DATASET_TYPE"
	ErrFailedToGetDatasetDagsMessage             = "ERR_FAILED_TO_GET_DATASET_DAGS"
)

var (
	ErrFailedToBuildQuery                 = errors.New(ErrFailedToBuildQueryMessage)
	ErrFailedToGetData                    = errors.New(ErrFailedToGetDataMessage)
	ErrFailedToGetDatasetColumnDatatype   = errors.New(ErrFailedToGetDatasetColumnDatatypeMessage)
	ErrNoRow                              = errors.New(ErrNoRowMessage)
	ErrMoreThanOneRow                     = errors.New(ErrMoreThanOneRowMessage)
	ErrInvalidMetadataFormat              = errors.New(ErrInvalidMetadataFormatMessage)
	ErrFailedToUnmarshalMetadata          = errors.New(ErrFailedToUnmarshalMetadataMessage)
	ErrFailedToGetDataset                 = errors.New(ErrFailedToGetDatasetMessage)
	ErrFailedToGetDatasetById             = errors.New(ErrFailedToGetDatasetByIdMessage)
	ErrFailedToPopulateMultiSelectOptions = errors.New(ErrFailedToPopulateMultiSelectOptionsMessage)
	ErrFailedToRegisterDataset            = errors.New(ErrFailedToRegisterDatasetMessage)
	ErrFailedToRegisterDatasetJob         = errors.New(ErrFailedToRegisterDatasetJobMessage)
	ErrFailedToMarshalMetadata            = errors.New(ErrFailedToMarshalMetadataMessage)
	ErrFailedToUpdateDataset              = errors.New(ErrFailedToUpdateDatasetMessage)
	ErrFailedToUpdateDatasetAction        = errors.New(ErrFailedToUpdateDatasetActionMessage)
	ErrFailedToUpdateRulePriority         = errors.New(ErrFailedToUpdateRulePriorityMessage)
	ErrNoRulesPresent                     = errors.New(ErrNoRulesPresentMessage)
	ErrFailedToGetRule                    = errors.New(ErrFailedToGetRuleMessage)
	ErrFailedToUpsertTemplate             = errors.New(ErrFailedToUpsertTemplateMessage)
	ErrColumnMetadataTypeIsEmpty          = errors.New(ErrColumnMetadataTypeIsEmptyMessage)
	ErrFailedToParseDatasetId             = errors.New(ErrFailedToParseDatasetIdMessage)
	ErrFailedToGetDatasetMetadata         = errors.New(ErrFailedToGetDatasetMetadataMessage)
	ErrInvalidDataplatformProvider        = errors.New(ErrInvalidDataplatformProviderMessage)
	ErrInvalidDatalistinSortColumn        = errors.New(ErrInvalidDatalistinSortColumnMessage)
	ErrInvalidDatasetType                 = errors.New(ErrInvalidDatasetTypeMessage)
	ErrFailedToGetDatasetDags             = errors.New(ErrFailedToGetDatasetDagsMessage)
)

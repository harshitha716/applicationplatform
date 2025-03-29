package constants

import "fmt"

const JobMappingsTableName = "job_mappings"
const JobMappingsTableNameQueryParam = "job_mappings_table_name"

type JobMappingType string

const (
	JobMappingTypeDataset JobMappingType = "dataset"
	JobMappingTypeFolder  JobMappingType = "folder"
)

const (
	JobMappingIdColumnName                = "id"
	JobMappingDestinationTypeColumnName   = "destination_type"
	JobMappingDestinationValueColumnName  = "destination_value"
	JobMappingIsDeletedColumnName         = "is_deleted"
	JobMappingSourceTypeFolderColumnName  = "source_type_folder"
	JobMappingSourceTypeDatasetColumnName = "source_type_dataset"
	JobMappingSourceTypeColumnName        = "source_type"
	JobMappingSourceValueColumnName       = "source_value"
	JobMappingJobIdColumnName             = "job_id"
	JobMappingJobParamsColumnName         = "job_params"
	JobMappingDeletedAtColumnName         = "deleted_at"
	JobMappingDeletedByColumnName         = "deleted_by"
	JobMappingMerchantIdColumnName        = "merchant_id"
)

var SelectJobMappingColumnNames string = fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s", JobMappingIdColumnName, JobMappingDestinationTypeColumnName, JobMappingDestinationValueColumnName, JobMappingIsDeletedColumnName, JobMappingSourceTypeColumnName, JobMappingSourceValueColumnName, JobMappingJobIdColumnName, JobMappingJobParamsColumnName, JobMappingDeletedAtColumnName, JobMappingDeletedByColumnName, JobMappingMerchantIdColumnName)

var QueryGetDatasetParents = fmt.Sprintf("SELECT %s FROM {{.%s}} WHERE destination_type = '{{.%s}}' AND destination_value = '{{.%s}}' AND is_deleted = false", SelectJobMappingColumnNames, JobMappingsTableNameQueryParam, JobMappingDestinationTypeColumnName, JobMappingDestinationValueColumnName)

var QueryGetDatasetEdgesByMerchant = fmt.Sprintf("SELECT %s FROM {{.%s}} WHERE merchant_id = '{{.%s}}' AND source_type IN ('{{.%s}}', '{{.%s}}') AND destination_type = '{{.%s}}' AND is_deleted = false", SelectJobMappingColumnNames, JobMappingsTableNameQueryParam, JobMappingMerchantIdColumnName, JobMappingSourceTypeFolderColumnName, JobMappingSourceTypeDatasetColumnName, JobMappingDestinationTypeColumnName)

package constants

const JobsTableNameQueryParam = "jobs_table_name"
const JobsTableName = "jobs"

const (
	JobIdColumnName        = "id"
	JobTypeColumnName      = "type"
	JobIsDeletedColumnName = "is_deleted"
)

type JobType string

const (
	UpdateJobType JobType = "update_dataset_data"
)

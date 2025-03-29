package constants

const (
	OrganizationID                = "organization_id"
	TemporalWorkflowID            = "temporal_workflow_id"
	ConnectivityTemporalTaskQueue = "connectivity"
	CronExpression                = "cron_expression"
	DefaultScheduleGroup          = "default"
	ActiveScheduleStatus          = "active"

	// Object Storage
	BucketName = "bucket_name"

	// GCS
	GCS = "gcs"

	// Snowflake
	Snowflake                 = "snowflake"
	SnowflakeDatabase         = "snowflake_database"
	SnowflakeSchema           = "snowflake_schema"
	SnowflakeTable            = "snowflake_table"
	SnowflakeFilterColumnName = "snowflake_filter_column_name"
	S3Bucket                  = "s3_bucket"
	S3DestinationPath         = "s3_destination_path"
	S3StorageIntegration      = "s3_storage_integration"
	SnowflakeWorkflow         = "snowflake_workflow"

	// General
	StartOffsetDays = "start_offset_days"
	EndOffsetDays   = "end_offset_days"
)

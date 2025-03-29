package constants

type DatabricksJobEventType string

const (
	DatabricksJobEventTypeOnSuccess DatabricksJobEventType = "jobs.on_success"
	DatabricksJobEventTypeOnFailure DatabricksJobEventType = "jobs.on_failure"
	DatabricksJobEventTypeOnStart   DatabricksJobEventType = "jobs.on_start"
)

type DatabricksJobType string

const (
	DatabricksJobTypeIngestion      DatabricksJobType = "ingestion"
	DatabricksJobTypeTransformation DatabricksJobType = "transformation"
	DatabricksJobTypeDag            DatabricksJobType = "dag"
	DatabricksJobTypeCron           DatabricksJobType = "cron"
)

type DatabricksJobSourceType string

const (
	DatabricksJobSourceTypeFolder  DatabricksJobSourceType = "folder"
	DatabricksJobSourceTypeDataset DatabricksJobSourceType = "dataset"
)

type DatabricksJobDestinationType string

const (
	DatabricksJobDestinationTypeDataset DatabricksJobDestinationType = "dataset"
)

type DatabricksColumnCustomType string

const (
	DatabricksColumnCustomTypeCurrency DatabricksColumnCustomType = "currency"
	DatabricksColumnCustomTypeAmount   DatabricksColumnCustomType = "amount"
	DatabricksColumnCustomTypeDateTime DatabricksColumnCustomType = "date_time"
	DatabricksColumnCustomTypeCountry  DatabricksColumnCustomType = "country"
	DatabricksColumnCustomTypeTags     DatabricksColumnCustomType = "tags"
	DatabricksColumnCustomTypeBank     DatabricksColumnCustomType = "bank"
)

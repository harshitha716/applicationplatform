package models

import "github.com/Zampfi/application-platform/services/api/core/dataplatform/constants"

type RegisterJobActionPayload struct {
	MerchantId           string                                 `json:"merchant_id"`
	JobType              constants.DatabricksJobType            `json:"job_type"`
	SourceType           constants.DatabricksJobSourceType      `json:"source_type"`
	SourceValue          string                                 `json:"source_value"`
	DestinationType      constants.DatabricksJobDestinationType `json:"destination_type"`
	DestinationValue     string                                 `json:"destination_value"`
	TemplateId           string                                 `json:"template_id"`
	QuartzCronExpression string                                 `json:"quartz_cron_expression"`
}

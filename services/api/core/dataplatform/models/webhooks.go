package models

import "github.com/Zampfi/application-platform/services/api/core/dataplatform/constants"

type Run struct {
	RunId int64 `json:"run_id"`
}

type Job struct {
	JobId int64  `json:"job_id"`
	Name  string `json:"name"`
}

type DatabricksJobStatusUpdatePayload struct {
	EventType   constants.DatabricksJobEventType `json:"event_type"`
	WorkspaceId int64                            `json:"workspace_id"`
	Run         Run                              `json:"run"`
	Job         Job                              `json:"job"`
}

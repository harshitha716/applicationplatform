package constants

import (
	"fmt"

	dataconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
)

const ActionsTableNameQueryParam = "actions_table_name"
const ActionsTableName = "actions"

const (
	ActionIdColumnName          = "id"
	ActionWorkspaceIdColumnName = "workspace_id"
	ActionTypeColumnName        = "action_type"
	ActionMetadataColumnName    = "action_metadata"
	ActionStatusColumnName      = "status"
	ActionCreatedAtColumnName   = "created_at"
	ActionUpdatedAtColumnName   = "updated_at"
	ActionActorIdColumnName     = "actor_id"
)
const (
	ActionRunIdColumnName = "run_id"
)

var InsertActionColumnNames string = fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s, %s", ActionIdColumnName, ActionWorkspaceIdColumnName, ActionTypeColumnName, ActionMetadataColumnName, ActionStatusColumnName, ActionCreatedAtColumnName, ActionUpdatedAtColumnName, ActionActorIdColumnName)

var SelectActionColumnNames string = fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s, %s, %s", ActionIdColumnName, ActionWorkspaceIdColumnName, ActionTypeColumnName, ActionMetadataColumnName, ActionStatusColumnName, ActionCreatedAtColumnName, ActionUpdatedAtColumnName, ActionRunIdColumnName, ActionActorIdColumnName)

var QueryCreateAction string = fmt.Sprintf("INSERT INTO {{.%s}} (%s) VALUES ('{{.%s}}', '{{.%s}}', '{{.%s}}', '{{.%s}}', '{{.%s}}', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '{{.%s}}')", ActionsTableNameQueryParam, InsertActionColumnNames, ActionIdColumnName, ActionWorkspaceIdColumnName, ActionTypeColumnName, ActionMetadataColumnName, ActionStatusColumnName, ActionActorIdColumnName)

var QueryGetActionByRunId string = fmt.Sprintf("SELECT %s FROM {{.%s}} WHERE %s = '{{.%s}}' and %s = '{{.%s}}'", SelectActionColumnNames, ActionsTableNameQueryParam, ActionRunIdColumnName, ActionRunIdColumnName, ActionWorkspaceIdColumnName, ActionWorkspaceIdColumnName)

var QueryUpdateActionRunId string = fmt.Sprintf("UPDATE {{.%s}} SET %s = '{{.%s}}', %s = CURRENT_TIMESTAMP WHERE id = '{{.%s}}'", ActionsTableNameQueryParam, ActionRunIdColumnName, ActionRunIdColumnName, ActionUpdatedAtColumnName, ActionIdColumnName)

const (
	DatasetId = "dataset_id"
	Query     = "query"
)

var QueryCreateMV string = fmt.Sprintf("CREATE OR REPLACE MATERIALIZED VIEW {{.%s}} AS {{.%s}}", DatasetId, Query)

var QueryUpdateActionStatus string = fmt.Sprintf("UPDATE {{.%s}} SET %s = '{{.%s}}', %s = CURRENT_TIMESTAMP WHERE %s = '{{.%s}}' and %s = '{{.%s}}'", ActionsTableNameQueryParam, ActionStatusColumnName, ActionStatusColumnName, ActionUpdatedAtColumnName, ActionRunIdColumnName, ActionRunIdColumnName, ActionWorkspaceIdColumnName, ActionWorkspaceIdColumnName)

var QueryGetJobIdForDatasetForJobType string = fmt.Sprintf(
	"SELECT {{.%s}}.%s FROM {{.%s}} JOIN {{.%s}} ON {{.%s}}.%s = {{.%s}}.%s WHERE {{.%s}}.%s = '{{.%s}}' AND {{.%s}}.%s = '{{.%s}}' AND {{.%s}}.%s = false AND {{.%s}}.%s = false AND {{.%s}}.%s = '%s';",
	dataconstants.JobMappingsTableNameQueryParam,
	dataconstants.JobMappingJobIdColumnName,
	dataconstants.JobMappingsTableNameQueryParam,
	dataconstants.JobsTableNameQueryParam,
	dataconstants.JobMappingsTableNameQueryParam,
	dataconstants.JobMappingJobIdColumnName,
	dataconstants.JobsTableNameQueryParam,
	dataconstants.JobIdColumnName,
	dataconstants.JobMappingsTableNameQueryParam,
	dataconstants.JobMappingSourceValueColumnName,
	dataconstants.JobMappingSourceValueColumnName,
	dataconstants.JobMappingsTableNameQueryParam,
	dataconstants.JobMappingSourceTypeColumnName,
	dataconstants.JobMappingSourceTypeColumnName,
	dataconstants.JobMappingsTableNameQueryParam,
	dataconstants.JobMappingIsDeletedColumnName,
	dataconstants.JobsTableNameQueryParam,
	dataconstants.JobIsDeletedColumnName,
	dataconstants.JobsTableNameQueryParam,
	dataconstants.JobTypeColumnName,
	dataconstants.UpdateJobType,
)

var QueryGetActionById string = fmt.Sprintf("SELECT %s FROM {{.%s}} WHERE %s = '{{.%s}}'", SelectActionColumnNames, ActionsTableNameQueryParam, ActionIdColumnName, ActionIdColumnName)

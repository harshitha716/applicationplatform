package models

import (
	"encoding/json"
	"time"

	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	querybuildermodels "github.com/Zampfi/application-platform/services/api/pkg/querybuilder/models"
	"github.com/google/uuid"
)

type Rule struct {
	ID             uuid.UUID    `json:"rule_id"`
	OrganizationId uuid.UUID    `json:"organization_id"`
	DatasetId      uuid.UUID    `json:"dataset_id"`
	Column         string       `json:"column"`
	Value          string       `json:"value"`
	FilterConfig   FilterConfig `json:"filter_config"`
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	Priority       int          `json:"priority"`
	CreatedAt      time.Time    `json:"created_at"`
	CreatedBy      uuid.UUID    `json:"created_by"`
	UpdatedAt      time.Time    `json:"updated_at"`
	UpdatedBy      uuid.UUID    `json:"updated_by"`
	DeletedAt      *time.Time   `json:"deleted_at"`
	DeletedBy      *uuid.UUID   `json:"deleted_by"`
}

func (r *Rule) FromSchema(schema *dbmodels.Rule) error {
	filterConfig := FilterConfig{}

	err := json.Unmarshal(schema.FilterConfig, &filterConfig)
	if err != nil {
		return err
	}

	r.ID = schema.ID
	r.OrganizationId = schema.OrganizationId
	r.DatasetId = schema.DatasetId
	r.Column = schema.Column
	r.Value = schema.Value
	r.FilterConfig = filterConfig
	r.Title = schema.Title
	r.Description = schema.Description
	r.Priority = schema.Priority
	r.CreatedAt = schema.CreatedAt
	r.CreatedBy = schema.CreatedBy
	r.UpdatedAt = schema.UpdatedAt
	r.UpdatedBy = schema.UpdatedBy
	r.DeletedAt = schema.DeletedAt
	r.DeletedBy = schema.DeletedBy

	return nil
}

type FilterConfig struct {
	QueryConfig querybuildermodels.QueryConfig `json:"query_config"`
	Sql         string                         `json:"sql"`
	Args        map[string]interface{}         `json:"args"`
}

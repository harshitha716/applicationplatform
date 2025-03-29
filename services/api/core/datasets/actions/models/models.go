package models

import (
	"encoding/json"
	"time"

	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"

	"github.com/google/uuid"
)

type DatasetAction struct {
	ID             uuid.UUID
	ActionId       string
	ActionType     string
	DatasetId      uuid.UUID
	OrganizationId uuid.UUID
	Status         string
	Config         interface{}
	ActionBy       uuid.UUID
	StartedAt      time.Time
	CompletedAt    *time.Time
}

func (d *DatasetAction) FromSchema(schema dbmodels.DatasetAction) {
	d.ID = schema.ID
	d.ActionId = schema.ActionId
	d.DatasetId = schema.DatasetId
	d.OrganizationId = schema.OrganizationId
	d.ActionType = schema.ActionType
	d.Status = schema.Status
	var config interface{}
	_ = json.Unmarshal(schema.Config, &config)
	d.Config = config
	d.ActionBy = schema.ActionBy
	d.StartedAt = schema.StartedAt
	d.CompletedAt = schema.CompletedAt
}

func (d *DatasetAction) ToSchema() dbmodels.DatasetAction {
	return dbmodels.DatasetAction{
		ID:             d.ID,
		ActionId:       d.ActionId,
		ActionType:     d.ActionType,
		DatasetId:      d.DatasetId,
		OrganizationId: d.OrganizationId,
		Status:         d.Status,
		Config:         json.RawMessage(d.Config.([]byte)),
		ActionBy:       d.ActionBy,
		StartedAt:      d.StartedAt,
		CompletedAt:    d.CompletedAt,
	}
}

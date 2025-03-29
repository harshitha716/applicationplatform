package models

import (
	"encoding/json"
	"fmt"
	"time"

	widgetmodels "github.com/Zampfi/application-platform/services/api/core/widgets/models"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/google/uuid"
)

type SheetConfigVersion string

type Sheet struct {
	ID              uuid.UUID                     `json:"sheet_id" gorm:"column:sheet_id"`
	Name            string                        `json:"name"`
	Description     *string                       `json:"description"`
	WidgetInstances []widgetmodels.WidgetInstance `json:"widget_instances,omitempty"`
	CreatedAt       time.Time                     `json:"created_at"`
	UpdatedAt       time.Time                     `json:"updated_at"`
	DeletedAt       *time.Time                    `json:"deleted_at,omitempty"`
	FractionalIndex float64                       `json:"fractional_index"`
	PageId          uuid.UUID                     `json:"page_id"`
	SheetConfig     SheetConfig                   `json:"sheet_config"`
}

type SheetConfig struct {
	Version            string               `json:"version"`
	NativeFilterConfig []NativeFilterConfig `json:"native_filter_config"`
	SheetLayout        []WidgetGroupLayout  `json:"sheet_layout"`
	Currency           *SheetCurrencyConfig `json:"currency,omitempty"`
}

type SheetCurrencyConfig struct {
	HideCurrencyFilter bool   `json:"hide_currency_filter"`
	DefaultCurrency    string `json:"default_currency,omitempty"`
}

type WidgetGroupLayout struct {
	Name          string      `json:"name,omitempty"`
	Layout        *Layout     `json:"layout,omitempty"`
	DefaultWidget *uuid.UUID  `json:"default_widget,omitempty"`
	WidgetGroup   []uuid.UUID `json:"widget_group,omitempty"`
}

type Layout struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	W float64 `json:"w"`
	H float64 `json:"h"`
}

type NativeFilterConfig struct {
	Name           string              `json:"name"`
	Id             string              `json:"id"`
	FilterType     string              `json:"filter_type"`
	DataType       string              `json:"data_type"`
	WidgetsInScope []string            `json:"widgets_in_scope"`
	Targets        []FilterTarget      `json:"targets"`
	DefaultValue   *DefaultFilterValue `json:"default_value,omitempty"`
}

type DefaultFilterValue struct {
	Operator string   `json:"operator,omitempty"`
	Value    []string `json:"value,omitempty"`
}

type FilterTarget struct {
	DatasetId uuid.UUID `json:"dataset_id"`
	Column    string    `json:"column"`
}

type RangeDefaultValue struct {
	Operator string      `json:"operator"`
	From     interface{} `json:"from"`
	To       interface{} `json:"to,omitempty"`
}

// TODO implement JSON versioning for sheet config
func (s *Sheet) FromDB(schema *dbmodels.Sheet) error {
	s.ID = schema.ID
	s.Name = schema.Name
	s.Description = schema.Description
	s.WidgetInstances = make([]widgetmodels.WidgetInstance, len(schema.WidgetInstances))
	for i, widgetInstance := range schema.WidgetInstances {
		s.WidgetInstances[i] = widgetmodels.WidgetInstance{}
		if err := s.WidgetInstances[i].FromDB(&widgetInstance); err != nil {
			return err
		}
	}
	s.CreatedAt = schema.CreatedAt
	s.UpdatedAt = schema.UpdatedAt
	s.DeletedAt = schema.DeletedAt
	s.FractionalIndex = schema.FractionalIndex
	s.PageId = schema.PageId

	sheetConfig := SheetConfig{}
	if err := json.Unmarshal(schema.SheetConfig, &sheetConfig); err != nil {
		return err
	}
	s.SheetConfig = sheetConfig

	return nil
}

func (s *Sheet) ToDB() (*dbmodels.Sheet, error) {
	sheetConfig, err := json.Marshal(s.SheetConfig)
	if err != nil {
		return nil, err
	}

	return &dbmodels.Sheet{
		ID:              s.ID,
		Name:            s.Name,
		Description:     s.Description,
		FractionalIndex: s.FractionalIndex,
		PageId:          s.PageId,
		SheetConfig:     sheetConfig,
	}, nil
}

type CreateSheetPayload struct {
	Name        string  `form:"name"`
	Description *string `form:"description"`
	PageId      string  `form:"page_id"`
	SheetConfig string  `form:"sheet_config"`
}

func (c *CreateSheetPayload) ToModel() (*Sheet, error) {
	var sheetConfig SheetConfig
	if err := json.Unmarshal([]byte(c.SheetConfig), &sheetConfig); err != nil {
		return nil, fmt.Errorf("invalid sheet config JSON: %w", err)
	}

	pageUUID, err := uuid.Parse(c.PageId)
	if err != nil {
		return nil, fmt.Errorf("invalid page id: %w", err)
	}

	return &Sheet{
		Name:        c.Name,
		Description: c.Description,
		PageId:      pageUUID,
		SheetConfig: sheetConfig,
	}, nil
}

// Helper function to get a string pointer, returning nil for empty strings
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

type UpdateSheetPayload struct {
	SheetId     string `form:"sheet_id"`
	Name        string `form:"name"`
	Description string `form:"description"`
	PageId      string `form:"page_id"`
	SheetConfig string `form:"sheet_config"`
}

func (u *UpdateSheetPayload) ToModel() (*Sheet, error) {
	sheetUUID, err := uuid.Parse(u.SheetId)
	if err != nil {
		return nil, fmt.Errorf("invalid sheet id: %w", err)
	}

	model := &Sheet{
		ID: sheetUUID,
	}

	// Handle Name
	if name := stringPtr(u.Name); name != nil {
		model.Name = *name
	}

	// Handle Description
	model.Description = stringPtr(u.Description)

	// Handle PageId
	if pageId := stringPtr(u.PageId); pageId != nil {
		pageUUID, err := uuid.Parse(*pageId)
		if err != nil {
			return nil, fmt.Errorf("invalid page id: %w", err)
		}
		model.PageId = pageUUID
	}

	// Handle SheetConfig
	if config := stringPtr(u.SheetConfig); config != nil {
		var sheetConfig SheetConfig
		if err := json.Unmarshal([]byte(*config), &sheetConfig); err != nil {
			return nil, fmt.Errorf("invalid sheet config JSON: %w", err)
		}
		model.SheetConfig = sheetConfig
	}

	return model, nil
}

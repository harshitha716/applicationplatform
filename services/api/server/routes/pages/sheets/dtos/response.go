package dtos

import (
	"encoding/json"
	"time"

	"github.com/Zampfi/application-platform/services/api/core/sheets/models"
	widgetmodels "github.com/Zampfi/application-platform/services/api/core/widgets/models"
	dbmodels "github.com/Zampfi/application-platform/services/api/db/models"
	widgetInstanceDtos "github.com/Zampfi/application-platform/services/api/server/routes/widgets/dtos"
	"github.com/google/uuid"
)

type SheetResponse struct {
	SheetId         uuid.UUID                                   `json:"sheet_id"`
	Name            string                                      `json:"name"`
	Description     *string                                     `json:"description,omitempty"`
	CreatedAt       time.Time                                   `json:"created_at"`
	UpdatedAt       time.Time                                   `json:"updated_at"`
	WidgetInstances []widgetInstanceDtos.WidgetInstanceResponse `json:"widget_instances"`
	SheetConfig     SheetConfig                                 `json:"sheet_config"`
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

func (s *SheetResponse) NewSheetResponse(sheet *dbmodels.Sheet) error {

	widgetInstances := make([]widgetInstanceDtos.WidgetInstanceResponse, 0)
	for _, widgetInstance := range sheet.WidgetInstances {
		widgetInstanceModel := widgetmodels.WidgetInstance{}
		widgetInstanceModel.FromDB(&widgetInstance)
		widgetInstanceResponse, err := widgetInstanceDtos.NewWidgetInstanceResponse(&widgetInstanceModel)
		if err != nil {
			return err
		}
		widgetInstances = append(widgetInstances, *widgetInstanceResponse)
	}

	s.SheetId = sheet.ID
	s.Name = sheet.Name
	s.Description = sheet.Description
	s.CreatedAt = sheet.CreatedAt
	s.UpdatedAt = sheet.UpdatedAt
	s.WidgetInstances = widgetInstances
	s.SheetConfig = SheetConfig{}
	err := json.Unmarshal(sheet.SheetConfig, &s.SheetConfig)
	if err != nil {
		return err
	}
	return nil
}

type SheetFilterConfig struct {
	NativeFilterConfig []NativeFilterConfig `json:"native_filter_config"`
}

type NativeFilterConfig struct {
	Name           string                     `json:"name"`
	Id             string                     `json:"id"`
	FilterType     string                     `json:"filter_type"`
	DataType       string                     `json:"data_type"`
	WidgetsInScope []string                   `json:"widgets_in_scope"`
	Targets        []FilterTarget             `json:"targets"`
	Options        []interface{}              `json:"options,omitempty"`
	DefaultValue   *models.DefaultFilterValue `json:"default_value,omitempty"`
}

type FilterTarget struct {
	DatasetID string `json:"dataset_id"`
	Column    string `json:"column"`
}

func (s *SheetFilterConfig) NewSheetFilterConfig(filterConfig *models.FilterOptionsConfig) {
	for _, filter := range filterConfig.NativeFilterConfig {
		targets := make([]FilterTarget, 0)
		for _, target := range filter.Targets {
			targets = append(targets, FilterTarget{
				DatasetID: target.DatasetId.String(),
				Column:    target.Column,
			})
		}
		s.NativeFilterConfig = append(s.NativeFilterConfig, NativeFilterConfig{
			Id:             filter.Id,
			Name:           filter.Name,
			FilterType:     filter.FilterType,
			DataType:       filter.DataType,
			WidgetsInScope: filter.WidgetsInScope,
			Targets:        targets,
			Options:        filter.Options,
			DefaultValue:   filter.DefaultValue,
		})
	}
}

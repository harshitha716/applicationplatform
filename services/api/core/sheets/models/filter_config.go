package models

type FilterOptionsConfig struct {
	NativeFilterConfig []FilterOptionsModel `json:"native_filter_config"`
}

type FilterOptionsModel struct {
	Id             string              `json:"id"`
	Name           string              `json:"name"`
	FilterType     string              `json:"filter_type"`
	DataType       string              `json:"data_type"`
	WidgetsInScope []string            `json:"widgets_in_scope"`
	Targets        []FilterTarget      `json:"targets"`
	Options        []interface{}       `json:"options,omitempty"`
	DefaultValue   *DefaultFilterValue `json:"default_value,omitempty"`
}

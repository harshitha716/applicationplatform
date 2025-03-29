package models

type UpsertTemplateActionPayload struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	TemplateType  string `json:"template_type"`
	Configuration string `json:"configuration"`
}

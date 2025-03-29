package models

type TranslateQueryRequest struct {
	Query        string `json:"query"`
	OutputFormat string `json:"output_format"`
}

type TranslateQueryResponse struct {
	Query string `json:"query"`
}

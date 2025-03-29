package models

type DatabricksConfig struct {
	ServerHostname string `json:"serverHostname"`
	HttpPath       string `json:"httpPath"`
	Port           int    `json:"port"`
	AccessToken    string `json:"accessToken"`
	WarehouseId    string `json:"warehouseId"`
}

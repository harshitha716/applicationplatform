package models

type PinotConfig struct {
	BrokerList  []string `json:"brokerList"`
	AccessToken string   `json:"accessToken"`
}

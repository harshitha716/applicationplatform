package constants

type RosettaApiPath string

const RosettaTranslateSqlApiPath RosettaApiPath = "/translate/sql"

var RosettaApiHeaders = map[string]string{
	"Content-Type": "application/json",
}

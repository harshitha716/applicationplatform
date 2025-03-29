package helpers

import (
	"fmt"
	"time"

	"github.com/lithammer/shortuuid"
)

func GenerateUUIDWithUnderscores() string {
	delimiter := "_"
	randomId := shortuuid.New()
	_, month, day := time.Now().Date()
	uuid := fmt.Sprintf("%s%s%02d%s%02d", randomId, delimiter, int(month), delimiter, day)
	return uuid
}

package helpers

import (
	"encoding/json"
	"strings"
)

func ConvertToJSONString(data interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func ConvertToJSONStringWithReplacements(data interface{}, replacements map[string]string) (string, error) {
	jsonString, err := ConvertToJSONString(data)
	if err != nil {
		return "", err
	}
	for key, value := range replacements {
		jsonString = strings.ReplaceAll(jsonString, key, value)
	}
	return jsonString, nil
}

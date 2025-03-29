package helper

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Zampfi/application-platform/services/api/pkg/querybuilder/errors"
)

// ConvertInterfaceSliceToStrings converts a slice of interfaces to a slice of strings.
func ConvertInterfaceSliceToStrings(input interface{}) ([]string, error) {
	interfaces, ok := input.([]interface{})
	if !ok {
		return nil, errors.ErrInvalidDataType
	}

	stringsSlice := []string{}
	for _, inter := range interfaces {
		str, ok := inter.(string)
		if !ok {
			return nil, errors.ErrInvalidDataType
		}
		stringsSlice = append(stringsSlice, str)
	}

	return stringsSlice, nil
}

// removes the first and last character of a string if they are single quotes.
func RemoveFirstAndLastQuote(s string) string {
	if len(s) < 2 {
		return s
	}

	if strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'") {
		return s[1 : len(s)-1]
	}
	return s
}

func Contains[T comparable](slice []T, item T) bool {
	for _, eachItem := range slice {
		if eachItem == item {
			return true
		}
	}
	return false
}

func ToSqlValue(v interface{}) (interface{}, error) {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		var sqlValues []interface{}
		for i := 0; i < val.Len(); i++ {
			element := val.Index(i)
			elemStr, err := ToSqlValue(element.Interface())
			if err != nil {
				return "", err
			}
			sqlValues = append(sqlValues, elemStr)
		}
		return sqlValues, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatInt(val.Int(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'f', -1, 64), nil
	case reflect.String:
		return "'" + val.String() + "'", nil
	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), nil
	default:
		if t, ok := v.(time.Time); ok {
			return "'" + t.Format("2006-01-02 15:04:05") + "'::timestamp", nil
		}
		return "", errors.ErrInvalidDataType
	}
}

package models

import (
	"encoding/json"

	"github.com/Zampfi/application-platform/services/api/pkg/dataplatform/constants"

	"github.com/jmoiron/sqlx"
)

type Rows []map[string]interface{}

func (r *Rows) FromSqlRows(rows *sqlx.Rows, columnMetadata []ColumnMetadata) error {
	result := []map[string]interface{}{}

	columnMetadataMap := make(map[string]ColumnMetadata)
	for _, cm := range columnMetadata {
		columnMetadataMap[cm.Name] = cm
	}

	for rows.Next() {
		mapElement := map[string]interface{}{}
		if err := rows.MapScan(mapElement); err != nil {
			continue
		}

		result = append(result, processMap(mapElement, columnMetadataMap))
	}

	if err := rows.Err(); err != nil {
		return err
	}

	*r = result
	return nil
}

func processMap(m map[string]interface{}, columnMetadata map[string]ColumnMetadata) map[string]interface{} {
	for key, value := range m {
		colType := columnMetadata[key].DatabaseType

		switch colType {
		case constants.DatabaseTypeArray:
			if strValue, ok := value.(string); ok {
				var array []interface{}
				if err := json.Unmarshal([]byte(strValue), &array); err == nil {
					m[key] = array
				}
			}
		default:
			// handle other types as needed, for now just keep the value as it is
			m[key] = value
		}
	}
	return m
}

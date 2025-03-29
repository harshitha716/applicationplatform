package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type ColumnMetadata struct {
	Name         string `json:"name"`
	DatabaseType string `json:"database_type"`
}

func (cm *ColumnMetadata) FromSqlColumnType(colType *sql.ColumnType) {
	cm.Name = colType.Name()
	cm.DatabaseType = colType.DatabaseTypeName()
}

func GetColumnMetadataFromSqlRows(rows *sqlx.Rows) ([]ColumnMetadata, error) {
	var metadata []ColumnMetadata
	sqlColumnTypes, err := rows.ColumnTypes()
	if err != nil {
		return metadata, err
	}

	for _, colType := range sqlColumnTypes {
		var cm ColumnMetadata
		cm.FromSqlColumnType(colType)
		metadata = append(metadata, cm)
	}
	return metadata, nil
}

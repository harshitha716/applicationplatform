package models

import "github.com/jmoiron/sqlx"

type QueryResult struct {
	Rows    Rows             `json:"rows"`
	Columns []ColumnMetadata `json:"columns"`
}

func (qr *QueryResult) FromSqlRows(rows *sqlx.Rows) error {
	columnMetadata, err := GetColumnMetadataFromSqlRows(rows)
	if err != nil {
		return err
	}

	rowsData := Rows{}
	if err := rowsData.FromSqlRows(rows, columnMetadata); err != nil {
		return err
	}

	qr.Columns = columnMetadata
	qr.Rows = rowsData
	return nil
}

package dataset

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgresClient(dbURL string, maxIdleConns, maxOpenConns int) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	defer db.Close()

	return db, nil
}

func RunDatasetSeeds(pgClient *sql.DB) {

	var orgID, userId string
	err := pgClient.QueryRowContext(context.Background(), `
        SELECT organization_id, owner_id 
        FROM organizations 
        WHERE deleted_at IS NULL 
        LIMIT 1
    `).Scan(&orgID, &userId)
	if err != nil {
		panic(fmt.Errorf("error fetching organization data: %v", err))
	}

	datasetID := "10d8e092-ea1c-4e20-a1b4-a364201f9c99"
	createdBy := userId
	title := "test dataset"
	description := "test dataset description"
	dataType := "source"
	metadata := `{"columns": {"COUNTRY": {"custom_type": "country"}, "CURRENCY": {"custom_type": "currency"}, "PG_CHARGE_AMOUNT": {"custom_type": "amount"}, "PG_CHARGE_CREATED_AT_DATE": {"custom_type": "date_time", "custom_type_config": {"format": "MM/dd/yyyy"}}, "PGP_INTERNAL_TRANSACTION_STATUS": {"custom_type": "tags"}, "PGP_INTERNAL_TRANSACTIONS_AMOUNT": {"custom_type": "amount"}}, "display_config": [{"type": "currency_amount", "column": "PG_CHARGE_AMOUNT", "config": {"amount_column": "PG_CHARGE_AMOUNT", "currency_column": "CURRENCY"}, "is_hidden": false, "is_editable": false}, {"column": "CARD_BRAND", "is_hidden": false, "is_editable": false}, {"column": "PG_CHARGE_CREATED_AT_DATE", "is_hidden": false, "is_editable": false}, {"column": "PGP_INTERNAL_TRANSACTIONS_AMOUNT", "is_hidden": false, "is_editable": false}, {"column": "PGP_INTERNAL_TRANSACTIONS_CARD_TOKENIZATION_METHOD", "is_hidden": false, "is_editable": false}, {"column": "PGP_INTERNAL_TRANSACTIONS_CHARGE_STATUS", "is_hidden": false, "is_editable": false}, {"column": "PGP_INTERNAL_TRANSACTIONS_CODE", "is_hidden": false, "is_editable": false}, {"column": "PGP_INTERNAL_TRANSACTIONS_CREATED_AT", "is_hidden": false, "is_editable": false}, {"column": "PGP_INTERNAL_TRANSACTIONS_PAYMENT_INTENT_ID", "is_hidden": false, "is_editable": false}, {"column": "PGP_INTERNAL_TRANSACTIONS_SOURCE_ID", "is_hidden": false, "is_editable": false}, {"column": "PGP_INTERNAL_TRANSACTIONS_SUBMITTED_FOR_SETTLEMENT", "is_hidden": false, "is_editable": false}, {"column": "PGP_INTERNAL_TRANSACTION_STATUS", "is_hidden": false, "is_editable": false}, {"column": "PG_CHARGE_STATUS", "is_hidden": false, "is_editable": false}, {"column": "PG_EVENT_TYPE", "is_hidden": false, "is_editable": false}, {"column": "PG_ISSUING_BANK", "is_hidden": false, "is_editable": false}, {"column": "PG_PAYMENT_METHOD_TYPE", "is_hidden": false, "is_editable": false}, {"column": "PG_RECON_STATUS", "is_hidden": false, "is_editable": false}, {"column": "RECORD_DATE", "is_hidden": false, "is_editable": false}, {"column": "TRANSACTION_TYPE", "is_hidden": false, "is_editable": false}, {"column": "amount", "is_hidden": false, "is_editable": false}, {"column": "PG_CHARGE_TAGS", "is_hidden": false, "is_editable": true}, {"column": "PG_CHARGE_TAGS", "is_hidden": false, "is_editable": true}, {"column": "_zamp_id", "is_hidden": true, "is_editable": false}], "databricks_config": {"dedup_columns": ["PGP_INTERNAL_TRANSACTIONS_CREATED_AT", "PGP_INTERNAL_TRANSACTIONS_PAYMENT_INTENT_ID"], "cluster_columns": [], "order_by_column": "PGP_INTERNAL_TRANSACTIONS_CREATED_AT", "partition_columns": []}, "custom_column_groups": [{"type": "fx_handling", "config": {"amount_column": "PG_CHARGE_AMOUNT", "currency_column": "CURRENCY", "datetime_column": "PG_CHARGE_CREATED_AT_DATE"}}, {"type": "fx_handling", "config": {"amount_column": "PGP_INTERNAL_TRANSACTIONS_AMOUNT", "currency_column": "CURRENCY", "datetime_column": "PG_CHARGE_CREATED_AT_DATE"}}, {"type": "currency_amount", "config": {"alias": "currency_amount", "amount_column": "PG_CHARGE_AMOUNT", "currency_column": "CURRENCY"}}]}`

	_, err = pgClient.ExecContext(context.Background(), `
		INSERT INTO datasets (dataset_id, organization_id, title, description, type, created_by, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		datasetID, orgID, title, description, dataType, createdBy, metadata)
	if err != nil {
		fmt.Println("Error inserting into datasets:", err)
		return
	} else {
		fmt.Println("Inserted into datasets successfully")
	}

	_, err = pgClient.ExecContext(context.Background(), `
		INSERT INTO app.resource_audience_policies (
			resource_audience_type, resource_audience_id, privilege, resource_type, resource_id
		)
		VALUES ($1, $2, $3, $4, $5);
		`,
		"organization", orgID, "admin", "dataset", datasetID)
	if err != nil {
		fmt.Println("Error inserting into resource_access:", err)
		return
	} else {
		fmt.Println("Inserted into resource_access successfully")
	}
}

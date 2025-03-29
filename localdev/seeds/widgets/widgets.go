package widgets

import (
	"context"
	"database/sql"
	"fmt"
)

// Predefined UUIDs for widgets
var (
	// Widget Instance IDs
	BarChartInstanceId   = "44444444-4444-4444-4444-444444444444" // Fixed UUID for bar instance
	LineChartInstanceId  = "55555555-5555-5555-5555-555555555555" // Fixed UUID for line instance
	PieChartInstanceId   = "66666666-6666-6666-6666-666666666666" // Fixed UUID for pie instance
	PivotTableInstanceId = "88888888-8888-8888-8888-888888888888" // Fixed UUID for pivot instance
	KPIInstanceId        = "99999999-9999-9999-9999-999999999999" // Fixed UUID for KPI instance

	BarChartInstanceId2   = "44444444-4444-4444-4444-444444444445" // Fixed UUID for bar instance
	PivotTableInstanceId2 = "88888888-8888-8888-8888-888888888889" // Fixed UUID for pivot instance
	DonutChartInstanceId  = "66666666-6666-6666-6666-666666666667" // Fixed UUID for donut instance
	// Sheet ID
	SheetId = "d81fe5bf-8f4c-4b12-ba5d-20c2d8de8220"

	// Dataset ID
	DatasetId = "10d8e092-ea1c-4e20-a1b4-a364201f9c99"
)

func AddWidgetInstances(pgClient *sql.DB) {
	// Bar Chart Instance
	_, err := pgClient.ExecContext(context.Background(), `
			INSERT INTO widget_instances (widget_instance_id, widget_type, sheet_id, title, data_mappings, created_at, updated_at, deleted_at, default_filters)
			VALUES ($1, $2, $3, 'Sales by Category', 
			'{
				"version": "1",
				"mappings": [
					{
						"dataset_id": "10d8e092-ea1c-4e20-a1b4-a364201f9c99",	
						"fields": {
							"x_axis": [
								{
									"column": "CARD_BRAND",
									"type": "string",
									"field_type": "dimension"
								}
							],
							"y_axis": [
								{
									"column": "PGP_INTERNAL_TRANSACTIONS_AMOUNT",
									"type": "number",
									"field_type": "measure",
									"aggregation": "sum"
								}
							]
						}
					}
				]
			}', NOW(), NOW(), NULL,
			'{
				"logical_operator": "and",
				"conditions": [
					{
						"column": "CARD_BRAND",
						"operator": "eq",
						"value": "Visa"
					}
				]
			}'
			)
		`, BarChartInstanceId, "bar_chart", SheetId)

	if err != nil {
		fmt.Println("Error inserting bar chart widget instance: ", err)
		err = nil
	} else {
		fmt.Println("Bar chart widget instance inserted successfully")
	}

	// Bar Chart Instance 2
	_, err = pgClient.ExecContext(context.Background(), `
		INSERT INTO widget_instances (widget_instance_id, widget_type, sheet_id, title, data_mappings, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, 'Sales by Category 2', 
		'{
			"version": "1",
			"mappings": [
				{
					"dataset_id": "10d8e092-ea1c-4e20-a1b4-a364201f9c99",
					"fields": {
						"x_axis": [
							{
								"column": "PGP_INTERNAL_TRANSACTIONS_CREATED_AT",
								"type": "timestamp",
								"field_type": "dimension",
								"cardinality": "single",
								"expression": "PGP_INTERNAL_TRANSACTIONS_CREATED_AT"
							}
						],
						"y_axis": [
							{
								"column": "PGP_INTERNAL_TRANSACTIONS_AMOUNT",
								"type": "number",
								"field_type": "measure",
								"aggregation": "sum"
							}
						],
						"group_by": [
							{
								"column": "CARD_BRAND",
								"type": "string",
								"field_type": "dimension"
							}
						]
					}
				}
			]
		}', NOW(), NOW(), NULL)
		`, BarChartInstanceId2, "bar_chart", SheetId)

	if err != nil {
		fmt.Println("Error inserting bar chart widget instance 2: ", err)
		err = nil
	} else {
		fmt.Println("Bar chart widget instance 2 inserted successfully")
	}

	// Line Chart Instance
	_, err = pgClient.ExecContext(context.Background(), `
			INSERT INTO widget_instances (widget_instance_id, widget_type, sheet_id, title, data_mappings, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, 'Monthly Trend', 
			'{
				"version": "1",
				"mappings": [
					{
						"dataset_id": "10d8e092-ea1c-4e20-a1b4-a364201f9c99",
						"fields": {
							"x_axis": [
								{
									"column": "TRANSACTION_TYPE",
									"type": "string",
									"field_type": "dimension"
								}
							],
							"y_axis": [
								{
									"column": "PG_AMOUNT_REFUNDED",
									"type": "number",
									"field_type": "measure",
									"aggregation": "sum"
								}
							]
						}
					}
				]
			}', NOW(), NOW(), NULL)
		`, LineChartInstanceId, "line_chart", SheetId)

	if err != nil {
		fmt.Println("Error inserting line chart widget instance: ", err)
		err = nil
	} else {
		fmt.Println("Line chart widget instance inserted successfully")
	}

	// Pie Chart Instance
	_, err = pgClient.ExecContext(context.Background(), `
			INSERT INTO widget_instances (widget_instance_id, widget_id, sheet_id, title, data_mappings, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, 'Revenue Distribution', 
			'{
				"version": "1",
				"mappings": [
					{
						"dataset_id": "10d8e092-ea1c-4e20-a1b4-a364201f9c99",
						"fields": {
							"slices": [
								{
									"column": "TRANSACTION_TYPE",
									"type": "string",
									"field_type": "dimension",
									"cardinality": "single"
								}
							],
							"values": [
								{
									"column": "PG_AMOUNT_REFUNDED",
									"type": "number",
									"field_type": "measure",
									"aggregation": "sum",
									"cardinality": "single"
								}
							]
						}
					}
				]
			}', NOW(), NOW(), NULL)
		`, PieChartInstanceId, "pie_chart", SheetId)

	// Donut Chart Instance
	_, err = pgClient.ExecContext(context.Background(), `
			INSERT INTO widget_instances (widget_instance_id, widget_type, sheet_id, title, data_mappings, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, 'Donut Chart', 
			'{
				"version": "1",
				"mappings": [
					{
						"dataset_id": "10d8e092-ea1c-4e20-a1b4-a364201f9c99",
						"fields": {
							"slices": [
								{
									"column": "CARD_BRAND",
									"type": "string",
									"field_type": "dimension"
								}
							],
							"values": [
								{
									"column": "PGP_INTERNAL_TRANSACTIONS_AMOUNT",
									"type": "number",
									"field_type": "measure",
									"aggregation": "sum"
								}
							]
						}
					}
				]
			}', NOW(), NOW(), NULL)
		`, DonutChartInstanceId, "donut_chart", SheetId)
	if err != nil {
		fmt.Println("Error inserting pie chart widget instance: ", err)
		err = nil
	} else {
		fmt.Println("Pie chart widget instance inserted successfully")
	}

	// Pivot Table Instance
	_, err = pgClient.ExecContext(context.Background(), `
			INSERT INTO widget_instances (widget_instance_id, widget_type, sheet_id, title, data_mappings, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, 'Transaction Analysis', 
			'{
				"version": "1",
				"mappings": [
					{
						"dataset_id": "10d8e092-ea1c-4e20-a1b4-a364201f9c99",
						"ref": "Card Category",
						"fields": {
							"rows": [
								{
									"column": "CARD_BRAND",
									"type": "string",
									"field_type": "dimension"
								},
								{
									"column": "PG_RECON_STATUS",
									"type": "string",
									"field_type": "dimension"
								}
							],
							"columns": [
								{
									"column": "PGP_INTERNAL_TRANSACTIONS_CREATED_AT",
									"type": "date",
									"field_type": "dimension",
									"expression": "date(PGP_INTERNAL_TRANSACTIONS_CREATED_AT)"
								}
							],
							"values": [
								{
									"column": "PGP_INTERNAL_TRANSACTIONS_AMOUNT",
									"type": "number",
									"field_type": "measure",
									"aggregation": "sum"
								},
								{
									"column": "PG_CHARGE_AMOUNT",
									"type": "number",
									"field_type": "measure",
									"aggregation": "sum"
								}
							]
						}
					}
				]
			}', NOW(), NOW(), NULL)
		`, PivotTableInstanceId, "pivot_table", SheetId)

	if err != nil {
		fmt.Println("Error inserting pivot table widget instance: ", err)
		err = nil
	} else {
		fmt.Println("Pivot table widget instance inserted successfully")
	}

	// KPI Instance
	_, err = pgClient.ExecContext(context.Background(), `
		INSERT INTO widget_instances (widget_instance_id, widget_type, sheet_id, title, data_mappings, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, 'Revenue', 
		'{
			"version": "1",
			"mappings": [
				{
					"dataset_id": "10d8e092-ea1c-4e20-a1b4-a364201f9c99",
					"fields": {
						"primary_value": [
							{
								"column": "PG_AMOUNT_REFUNDED",
								"type": "number",
								"field_type": "measure",
								"aggregation": "sum"
							}
						]
					}
				}
			]
		}', NOW(), NOW(), NULL)
		`, KPIInstanceId, "kpi", SheetId)

	if err != nil {
		fmt.Println("Error inserting KPI widget instance: ", err)
		err = nil
	} else {
		fmt.Println("KPI widget instance inserted successfully")
	}
}

func RunWidgetsSeeds(pgClient *sql.DB) {
	// AddWidgetTypes(pgClient)
	AddWidgetInstances(pgClient)
}

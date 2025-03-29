package pages

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

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

	pages := []struct {
		pageId      string
		name        string
		description string
		orgId       string
	}{
		{
			pageId:      "c71fe5bf-8f4c-4b12-ba5d-20c2d8de8124",
			name:        "Test page 1",
			description: "Test page 1 description",
			orgId:       orgID,
		},
		{
			pageId:      "d81fe5bf-8f4c-4b12-ba5d-20c2d8de8225",
			name:        "Test page 1",
			description: "Test page 1 description",
			orgId:       orgID,
		},
	}

	for index, page := range pages {

		_, err = pgClient.ExecContext(context.Background(), `
			INSERT INTO pages (page_id, name, description, organization_id) 
			VALUES ($1, $2, $3, $4)`,
			page.pageId, page.name, page.description, page.orgId,
		)
		if err != nil {
			fmt.Println("Error inserting into pages:", err)
			return
		}
		fmt.Println("Inserted into pages successfully. PageId: ", page.pageId)

		// insert into resource audience policies
		_, err = pgClient.ExecContext(context.Background(), `
		INSERT INTO app.resource_audience_policies (
			resource_audience_type, resource_audience_id, privilege, resource_type, resource_id
		)
		VALUES ($1, $2, $3, $4, $5);
		`, "organization", orgID, "admin", "page", page.pageId)
		if err != nil {
			fmt.Println("Error inserting into resource audience policies:", err)
			return
		}
		fmt.Println("Inserted into resource audience policies successfully")

		sheets := []struct {
			sheetId     string
			name        string
			description string
		}{
			{
				sheetId:     "d81fe5bf-8f4c-4b12-ba5d-20c2d8de822" + strconv.Itoa(index),
				name:        "Test sheet 1 page " + strconv.Itoa(index),
				description: "Test description sheet 1 page " + strconv.Itoa(index),
			},
			{
				sheetId:     "c71fe5bf-8f4c-4b12-ba5d-20c2d8de812" + strconv.Itoa(index),
				name:        "Test sheet 1 page " + strconv.Itoa(index),
				description: "Test description sheet 1 page " + strconv.Itoa(index),
			},
		}

		for _, sheet := range sheets {
			var query string
			var args []interface{}

			if index == 0 && sheet.name == "Test sheet 1 page 0" {
				query = `
					INSERT INTO sheets (sheet_id, name, description, page_id, sheet_config) 
					VALUES ($1, $2, $3, $4, $5)`
				args = []interface{}{
					sheet.sheetId,
					sheet.name,
					sheet.description,
					page.pageId,
					`{"version":"1","sheet_layout":{"44444444-4444-4444-4444-444444444444":{"h":5.5,"w":16,"x":0,"y":0},"55555555-5555-5555-5555-555555555555":{"h":5.5,"w":8,"x":0,"y":0},"66666666-6666-6666-6666-666666666666":{"h":5.5,"w":8,"x":8,"y":0}},"native_filter_config":[{"id":"11111111-1111-1111-1111-111111111111","name":"Card Brands","targets":[{"dataset_id":"10d8e092-ea1c-4e20-a1b4-a364201f9c99","column":"CARD_BRAND"}],"data_type":"string","filter_type":"search","default_value":{"value":["Master"],"operator":"contains"},"widgets_in_scope":["44444444-4444-4444-4444-444444444444","55555555-5555-5555-5555-555555555555","66666666-6666-6666-6666-666666666666"]},{"id":"22222222-2222-2222-2222-222222222222","name":"Charge Status","targets":[{"dataset_id":"10d8e092-ea1c-4e20-a1b4-a364201f9c99","column":"PGP_INTERNAL_TRANSACTIONS_CHARGE_STATUS"}],"data_type":"string","filter_type":"multi-select","default_value":{"value":["MasterCard","Discover","Apple Pay - MasterCard","Cashapp"],"operator":"in"},"widgets_in_scope":["44444444-4444-4444-4444-444444444444","55555555-5555-5555-5555-555555555555","66666666-6666-6666-6666-666666666666"]},{"id":"33333333-3333-3333-3333-333333333333","name":"Date","targets":[{"dataset_id":"10d8e092-ea1c-4e20-a1b4-a364201f9c99","column":"RECORD_DATE"}],"data_type":"date","filter_type":"date-range","default_value":{"value":["2025-01-01T18:30:00.000Z","2025-01-07T18:30:00.000Z"],"operator":"inbetween"},"widgets_in_scope":["44444444-4444-4444-4444-444444444444","55555555-5555-5555-5555-555555555555","66666666-6666-6666-6666-666666666666"]},{"id":"44444444-4444-4444-4444-444444444444","name":"Transaction Amount","targets":[{"dataset_id":"10d8e092-ea1c-4e20-a1b4-a364201f9c99","column":"PGP_INTERNAL_TRANSACTIONS_AMOUNT"}],"data_type":"number","filter_type":"amount-range","default_value":{"value":["1000"],"operator":"gte"},"widgets_in_scope":["44444444-4444-4444-4444-444444444444","55555555-5555-5555-5555-555555555555","66666666-6666-6666-6666-666666666666"]}]}`,
				}
			} else {
				query = `
					INSERT INTO sheets (sheet_id, name, description, page_id, sheet_config) 
					VALUES ($1, $2, $3, $4, $5)`
				args = []interface{}{
					sheet.sheetId,
					sheet.name,
					sheet.description,
					page.pageId,
					`{"version":"1","native_filter_config":[],"sheet_layout":{}}`,
				}
			}

			_, err = pgClient.ExecContext(context.Background(), query, args...)
			if err != nil {
				fmt.Println("Error inserting into sheets:", err)
				return
			}
			fmt.Println("Inserted into sheets successfully. SheetId: ", sheet.sheetId)
		}
	}
}

package tests

import (
	"context"
	"testing"

	actionconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/actions/constants"
	dataconstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	"github.com/Zampfi/application-platform/services/api/core/dataplatform/helpers"
	"github.com/stretchr/testify/assert"
)

func TestFillQueryTemplate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		query       string
		params      map[string]string
		expected    string
		expectError bool
	}{
		{
			name:        "Valid template",
			query:       "SELECT * FROM users WHERE id = {{.UserID}}",
			params:      map[string]string{"UserID": "123"},
			expected:    "SELECT * FROM users WHERE id = 123",
			expectError: false,
		},
		{
			name:        "Empty template",
			query:       "",
			params:      map[string]string{"UserID": "123"},
			expected:    "",
			expectError: false,
		},
		{
			name:        "Invalid template",
			query:       "SELECT * FROM users WHERE id = {{.UserID",
			params:      map[string]string{"UserID": "123"},
			expected:    "",
			expectError: true,
		},
		{
			name:        "Missing parameter",
			query:       "SELECT * FROM users WHERE id = {{.UserID}}",
			params:      map[string]string{}, // No UserID provided
			expected:    "SELECT * FROM users WHERE id = <no value>",
			expectError: false,
		},
		{
			name:        "Valid template with multiple parameters",
			query:       "SELECT * FROM users WHERE id = {{.UserID}} AND name = {{.UserName}}",
			params:      map[string]string{"UserID": "123", "UserName": "John"},
			expected:    "SELECT * FROM users WHERE id = 123 AND name = John",
			expectError: false,
		},
		{
			name:  "Valid Update query",
			query: actionconstants.QueryGetJobIdForDatasetForJobType,
			params: map[string]string{
				dataconstants.JobMappingsTableNameQueryParam:  "`zamp`.`platform`.`job_mappings`",
				dataconstants.JobsTableNameQueryParam:         "`zamp`.`platform`.`jobs`",
				dataconstants.JobMappingSourceValueColumnName: "123",
				dataconstants.JobMappingSourceTypeColumnName:  "dataset",
			},
			expected:    "SELECT `zamp`.`platform`.`job_mappings`.job_id FROM `zamp`.`platform`.`job_mappings` JOIN `zamp`.`platform`.`jobs` ON `zamp`.`platform`.`job_mappings`.job_id = `zamp`.`platform`.`jobs`.id WHERE `zamp`.`platform`.`job_mappings`.source_value = '123' AND `zamp`.`platform`.`job_mappings`.source_type = 'dataset' AND `zamp`.`platform`.`job_mappings`.is_deleted = false AND `zamp`.`platform`.`jobs`.is_deleted = false AND `zamp`.`platform`.`jobs`.type = 'update_dataset_data';",
			expectError: false,
		},
		{
			name:  "Valid Update query with multiple parameters",
			query: "SELECT DISTINCT PG_RECON_STATUS FROM {{.zamp_10d8e092-ea1c-4e20-a1b4-a364201f9c99}} where PG_RECON_STATUS = '{{.PG_RECON_STATUS}}'",
			params: map[string]string{
				"zamp_10d8e092-ea1c-4e20-a1b4-a364201f9c99": "10d8e092-ea1c-4e20-a1b4-a364201f9c99",
				"PG_RECON_STATUS": "SUCCESS",
			},
			expected:    "SELECT DISTINCT PG_RECON_STATUS FROM 10d8e092-ea1c-4e20-a1b4-a364201f9c99 where PG_RECON_STATUS = 'SUCCESS'",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := helpers.FillQueryTemplate(ctx, tt.query, tt.params)
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestAddCommentsToQuery(t *testing.T) {

	tests := []struct {
		query            string
		metadata         map[string]string
		expectedContains []string
	}{
		{
			query:            "SELECT * FROM users",
			metadata:         map[string]string{"org_id": "123", "trace_id": "456"},
			expectedContains: []string{"-- org_id='123'", "-- trace_id='456'"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			result := helpers.AddCommentsToQuery(tt.query, tt.metadata)
			for _, expectedComment := range tt.expectedContains {
				assert.Contains(t, result, expectedComment)
			}
		})
	}
}

func TestBuildDatabricksTableName(t *testing.T) {
	tests := []struct {
		catalog  string
		schema   string
		table    string
		expected string
	}{
		{
			catalog:  "a-b",
			schema:   "c-d",
			table:    "e-f",
			expected: "`a-b`.`c-d`.`e-f`",
		},
		{
			catalog:  "a_b",
			schema:   "c_d",
			table:    "e_f",
			expected: "`a_b`.`c_d`.`e_f`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.catalog, func(t *testing.T) {
			result := helpers.BuildDatabricksTableName(tt.catalog, tt.schema, tt.table)
			assert.Equal(t, tt.expected, result)
		})
	}
}

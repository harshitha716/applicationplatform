package store

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zampfi/application-platform/services/api/db/models"
	"github.com/Zampfi/application-platform/services/api/db/pgclient"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateRule(t *testing.T) {
	t.Parallel()

	orgID := uuid.New()
	datasetID := uuid.New()
	userID := uuid.New()
	column := "column"
	value := "value"

	tests := []struct {
		name      string
		params    models.CreateRuleParams
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			params: models.CreateRuleParams{
				Title:          "test",
				Description:    "test",
				OrganizationId: orgID,
				DatasetId:      datasetID,
				Column:         column,
				Value:          value,
				FilterConfig:   map[string]interface{}{"condition": "test"},
				CreatedBy:      userID,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("dataset", sqlmock.AnyArg(), userID, "admin", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), "dataset", uuid.New(), userID, "admin"))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "rules" SET "priority"=priority + 1,"updated_at"=$1,"updated_by"=$2 WHERE organization_id = $3 AND dataset_id = $4 AND "column" = $5 AND deleted_at IS NULL`)).
					WithArgs(sqlmock.AnyArg(), userID, orgID, datasetID, column).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("dataset", sqlmock.AnyArg(), userID, "admin", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), "dataset", uuid.New(), userID, "admin"))
				mock.ExpectExec(`INSERT INTO "rules"`).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			params: models.CreateRuleParams{
				Title:          "test",
				Description:    "test",
				OrganizationId: orgID,
				DatasetId:      datasetID,
				Column:         column,
				Value:          value,
				FilterConfig:   map[string]interface{}{"condition": "test"},
				CreatedBy:      userID,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("dataset", sqlmock.AnyArg(), userID, "admin", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), "dataset", uuid.New(), userID, "admin"))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "rules" SET "priority"=priority + 1,"updated_at"=$1,"updated_by"=$2 WHERE organization_id = $3 AND dataset_id = $4 AND "column" = $5 AND deleted_at IS NULL`)).
					WithArgs(sqlmock.AnyArg(), userID, orgID, datasetID, column).
					WillReturnError(gorm.ErrInvalidData)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			// add auth to ctx
			ctx := apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{orgID})

			// Execute
			err := store.CreateRule(ctx, tt.params)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetRuleById(t *testing.T) {
	t.Parallel()

	ruleID := uuid.New()
	orgID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		ruleID    uuid.UUID
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:   "success",
			ruleID: ruleID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"rule_id",
					"organization_id",
					"filter_config",
					"action_config",
					"priority",
					"created_at",
				}).AddRow(
					ruleID,
					orgID,
					[]byte(`{"condition":"test"}`),
					[]byte(`{"action":"test"}`),
					1,
					now,
				)

				mock.ExpectQuery(`SELECT \* FROM "rules" WHERE rule_id = \$1 ORDER BY "rules"."rule_id" LIMIT \$2`).
					WithArgs(ruleID, 1).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:   "not found",
			ruleID: ruleID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "rules" WHERE rule_id = \$1 ORDER BY "rules"."rule_id" LIMIT \$2`).
					WithArgs(ruleID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			rule, err := store.GetRuleById(context.Background(), tt.ruleID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, rule)
				assert.Equal(t, tt.ruleID, rule.ID)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateRulePriority(t *testing.T) {
	t.Parallel()

	ruleID1 := uuid.New()
	orgID := uuid.New()
	ruleID2 := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name      string
		params    models.UpdateRulePriorityParams
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success - update multiple rule priorities",
			params: models.UpdateRulePriorityParams{
				RulePriority: []models.RulePriority{
					{RuleId: ruleID1, Priority: 2},
					{RuleId: ruleID2, Priority: 1},
				},
				UpdatedBy: userID,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("dataset", sqlmock.AnyArg(), userID, "admin", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), "dataset", uuid.New(), userID, "admin"))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "rules" SET "priority"=CASE WHEN rule_id = '`+ruleID1.String()+`' THEN 2 WHEN rule_id = '`+ruleID2.String()+`' THEN 1 END,"updated_at"=$1,"updated_by"=$2 WHERE rule_id IN ($3,$4)`)).
					WithArgs(sqlmock.AnyArg(), userID, ruleID1, ruleID2).
					WillReturnResult(sqlmock.NewResult(1, 2))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "success - update single rule priority",
			params: models.UpdateRulePriorityParams{
				RulePriority: []models.RulePriority{
					{RuleId: ruleID1, Priority: 2},
				},
				UpdatedBy: userID,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("dataset", sqlmock.AnyArg(), userID, "admin", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), "dataset", uuid.New(), userID, "admin"))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "rules" SET "priority"=CASE WHEN rule_id = '`+ruleID1.String()+`' THEN 2 END,"updated_at"=$1,"updated_by"=$2 WHERE rule_id IN ($3)`)).
					WithArgs(sqlmock.AnyArg(), userID, ruleID1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			params: models.UpdateRulePriorityParams{
				RulePriority: []models.RulePriority{
					{RuleId: ruleID1, Priority: 1},
				},
				UpdatedBy: userID,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("dataset", sqlmock.AnyArg(), userID, "admin", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), "dataset", uuid.New(), userID, "admin"))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "rules" SET "priority"=CASE WHEN rule_id = '`+ruleID1.String()+`' THEN 1 END,"updated_at"=$1,"updated_by"=$2 WHERE rule_id IN ($3)`)).
					WithArgs(sqlmock.AnyArg(), userID, ruleID1).
					WillReturnError(gorm.ErrInvalidData)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "error - empty rule priority list",
			params: models.UpdateRulePriorityParams{
				RulePriority: []models.RulePriority{},
				UpdatedBy:    userID,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// No database calls expected
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			ctx := apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{orgID})

			err := store.UpdateRulePriority(ctx, tt.params)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteRule(t *testing.T) {
	t.Parallel()

	ruleID := uuid.New()
	orgID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		params    models.DeleteRuleParams
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "success",
			params: models.DeleteRuleParams{
				RuleId:    ruleID,
				DeletedBy: userID,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Get current rule
				rows := sqlmock.NewRows([]string{
					"rule_id", "organization_id", "priority", "created_at", "deleted_at",
				}).AddRow(ruleID, orgID, 2, now, nil)

				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "rules" WHERE rule_id = \$1 ORDER BY "rules"."rule_id" LIMIT \$2`).
					WithArgs(ruleID, 1).
					WillReturnRows(rows)

				// Update other rules' priorities
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("dataset", sqlmock.AnyArg(), userID, "admin", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), "dataset", uuid.New(), userID, "admin"))
				mock.ExpectExec(`UPDATE "rules" SET`).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mark rule as deleted
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "flattened_resource_audience_policies" WHERE resource_type = $1 AND resource_id = $2 AND user_id = $3 AND privilege = $4 AND deleted_at IS NULL LIMIT $5`)).
					WithArgs("dataset", sqlmock.AnyArg(), userID, "admin", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "resource_type", "resource_id", "user_id", "privilege"}).
						AddRow(uuid.New(), "dataset", uuid.New(), userID, "admin"))
				mock.ExpectExec(`UPDATE "rules" SET`).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "rule already deleted",
			params: models.DeleteRuleParams{
				RuleId:    ruleID,
				DeletedBy: userID,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"rule_id", "organization_id", "priority", "created_at", "deleted_at",
				}).AddRow(ruleID, orgID, 2, now, now)
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "rules" WHERE rule_id = \$1 ORDER BY "rules"."rule_id" LIMIT \$2`).
					WithArgs(ruleID, 1).
					WillReturnRows(rows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "rule not found",
			params: models.DeleteRuleParams{
				RuleId:    ruleID,
				DeletedBy: userID,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT \* FROM "rules" WHERE rule_id = \$1 ORDER BY "rules"."rule_id" LIMIT \$2`).
					WithArgs(ruleID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			ctx := apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{orgID})

			err := store.DeleteRule(ctx, tt.params)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetRules(t *testing.T) {
	t.Parallel()

	orgID := uuid.New()
	dataset1ID := uuid.New()
	dataset2ID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		params    models.FilterRuleParams
		mockSetup func(sqlmock.Sqlmock)
		want      map[string]map[string][]models.Rule
		wantErr   bool
	}{
		{
			name: "success - filter by organization",
			params: models.FilterRuleParams{
				OrganizationId: orgID,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"rule_id",
					"organization_id",
					"dataset_id",
					"column",
					"value",
					"filter_config",
					"title",
					"description",
					"priority",
					"created_at",
					"created_by",
				}).AddRow(
					uuid.New(),
					orgID,
					dataset1ID,
					"column1",
					"value1",
					[]byte(`{"condition":"test"}`),
					"title1",
					"desc1",
					1,
					now,
					userID,
				)

				mock.ExpectQuery(`SELECT \* FROM "rules" WHERE organization_id = \$1 AND deleted_at IS NULL ORDER BY priority asc`).
					WithArgs(orgID).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "success - filter by datasets and columns",
			params: models.FilterRuleParams{
				OrganizationId: orgID,
				DatasetColumns: []models.DatasetColumn{
					{
						DatasetId: dataset1ID,
						Columns:   []string{"column1", "column2"},
					},
					{
						DatasetId: dataset2ID,
						Columns:   []string{"column3"},
					},
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"rule_id",
					"organization_id",
					"dataset_id",
					"column",
					"value",
					"filter_config",
					"title",
					"description",
					"priority",
					"created_at",
					"created_by",
				}).AddRows(
					[]driver.Value{
						uuid.New(),
						orgID,
						dataset1ID,
						"column1",
						"value1",
						[]byte(`{"condition":"test"}`),
						"title1",
						"desc1",
						1,
						now,
						userID,
					},
					[]driver.Value{
						uuid.New(),
						orgID,
						dataset1ID,
						"column2",
						"value1",
						[]byte(`{"condition":"test"}`),
						"title1",
						"desc1",
						1,
						now,
						userID,
					},
					[]driver.Value{
						uuid.New(),
						orgID,
						dataset2ID,
						"column3",
						"value3",
						[]byte(`{"condition":"test"}`),
						"title3",
						"desc3",
						3,
						now,
						userID,
					},
				)

				mock.ExpectQuery(`SELECT \* FROM "rules" WHERE organization_id = \$1 AND dataset_id IN \(\$2,\$3\) AND "column" IN \(\$4,\$5,\$6\) AND deleted_at IS NULL ORDER BY priority asc`).
					WithArgs(orgID, dataset1ID, dataset2ID, "column1", "column2", "column3").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "database error",
			params: models.FilterRuleParams{
				OrganizationId: orgID,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "rules"`).
					WillReturnError(gorm.ErrInvalidData)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			ctx := apicontext.AddAuthToContext(context.Background(), "role", userID, []uuid.UUID{orgID})

			got, err := store.GetRules(ctx, tt.params)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				// For successful cases, verify the map structure
				if len(tt.params.DatasetColumns) > 0 {
					// Verify that we have entries for the specified datasets
					for _, dc := range tt.params.DatasetColumns {
						datasetRules, exists := got[dc.DatasetId.String()]
						assert.True(t, exists)
						// Verify that we have entries for the specified columns
						for _, col := range dc.Columns {
							_, exists := datasetRules[col]
							assert.True(t, exists)
						}
					}
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetRuleByIds(t *testing.T) {
	t.Parallel()

	ruleID1 := uuid.New()
	ruleID2 := uuid.New()
	orgID := uuid.New()
	now := time.Now()

	tests := []struct {
		name      string
		ruleIds   []uuid.UUID
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:    "success",
			ruleIds: []uuid.UUID{ruleID1, ruleID2},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"rule_id",
					"organization_id",
					"filter_config",
					"action_config",
					"priority",
					"created_at",
				}).AddRow(
					ruleID1,
					orgID,
					[]byte(`{"condition":"test1"}`),
					[]byte(`{"action":"test1"}`),
					1,
					now,
				).AddRow(
					ruleID2,
					orgID,
					[]byte(`{"condition":"test2"}`),
					[]byte(`{"action":"test2"}`),
					2,
					now,
				)

				mock.ExpectQuery(`SELECT \* FROM "rules" WHERE rule_id IN`).
					WithArgs(ruleID1, ruleID2).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:    "database error",
			ruleIds: []uuid.UUID{ruleID1, ruleID2},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "rules" WHERE rule_id IN`).
					WithArgs(ruleID1, ruleID2).
					WillReturnError(gorm.ErrInvalidData)
			},
			wantErr: true,
		},
		{
			name:    "empty ids list",
			ruleIds: []uuid.UUID{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"rule_id",
					"organization_id",
					"filter_config",
					"action_config",
					"priority",
					"created_at",
				})

				mock.ExpectQuery(`SELECT \* FROM "rules" WHERE rule_id IN`).
					WithArgs().
					WillReturnRows(rows)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gormDB, mock := getMockDB(t)
			store := &appStore{
				client: &pgclient.PostgresClient{DB: gormDB},
			}
			tt.mockSetup(mock)

			rules, err := store.GetRuleByIds(context.Background(), tt.ruleIds)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, rules)
			} else {
				assert.NoError(t, err)
				if len(tt.ruleIds) > 0 {
					assert.Len(t, rules, len(tt.ruleIds))
					assert.Equal(t, tt.ruleIds[0], rules[0].ID)
					if len(tt.ruleIds) > 1 {
						assert.Equal(t, tt.ruleIds[1], rules[1].ID)
					}
				} else {
					assert.Empty(t, rules)
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

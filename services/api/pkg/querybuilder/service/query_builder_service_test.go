package service

import (
	"context"
	"testing"
	"time"

	dataplatformConstants "github.com/Zampfi/application-platform/services/api/core/dataplatform/data/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/querybuilder/constants"
	"github.com/Zampfi/application-platform/services/api/pkg/querybuilder/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	service QueryBuilder
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) SetupTest() {
	s.service = NewQueryBuilder()
}

func (s *ServiceTestSuite) TestSelectToSQL() {
	ctx := context.Background()

	dataTypeString := dataplatformConstants.StringDataType
	dataTypeDecimal := dataplatformConstants.DecimalDataType
	dataTypeDate := dataplatformConstants.DateDataType
	dataTypeArray := dataplatformConstants.ArrayOfStringDataType
	dataTypeBoolean := dataplatformConstants.BooleanDataType
	alias := "alias"

	testCases := []struct {
		name           string
		queryConfig    models.QueryConfig
		expectedSQL    string
		expectedParams map[string]interface{}
		expectError    bool
	}{
		{
			name: "Basic Select with Pagination",
			queryConfig: models.QueryConfig{
				TableConfig: models.TableConfig{
					DatasetId: "products",
					Columns: []models.ColumnConfig{
						{
							Column: "id",
						},
						{
							Column: "name",
						},
						{
							Column: "price",
						},
					},
				},
				Pagination: &models.Pagination{
					Page:     1,
					PageSize: 10,
				},
			},
			expectedSQL:    "SELECT id, name, price FROM {{.zamp_products}} LIMIT 10 OFFSET 0",
			expectedParams: map[string]interface{}{"zamp_products": "products"},
			expectError:    false,
		},
		{
			name: "Select with Nested Filters",
			queryConfig: models.QueryConfig{
				TableConfig: models.TableConfig{
					DatasetId: "orders",
					Columns: []models.ColumnConfig{
						{
							Column: "order_id",
						},
						{
							Column: "customer_id",
						},
						{
							Column: "total_amount",
						},
					},
				},
				Filters: models.FilterModel{
					LogicalOperator: constants.LogicalOperatorAnd,
					Conditions: []models.Filter{
						{
							Column: models.ColumnConfig{
								Column:   "status",
								Datatype: &dataTypeString,
							},
							Operator: constants.EqualOperator,
							Value:    "shipped",
						},
						{
							LogicalOperator: &constants.LogicalOperatorOr,
							Conditions: []models.Filter{
								{
									Column: models.ColumnConfig{
										Column:   "total_amount",
										Datatype: &dataTypeDecimal,
									},
									Operator: constants.GreaterThanOperator,
									Value:    100,
								},
								{
									Column: models.ColumnConfig{
										Column:   "created_at",
										Datatype: &dataTypeDate,
									},
									Operator: constants.LessThanOperator,
									Value:    time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
								},
							},
						},
					},
				},
			},
			expectedSQL:    "SELECT order_id, customer_id, total_amount FROM {{.zamp_orders}} WHERE ( status = 'shipped' ) AND (( total_amount > 100 ) OR ( created_at < '2022-01-01 00:00:00'::timestamp ))",
			expectedParams: map[string]interface{}{"zamp_orders": "orders"},
			expectError:    false,
		},
		{
			name: "Complex Nested Filters with Multiple Levels",
			queryConfig: models.QueryConfig{
				TableConfig: models.TableConfig{
					DatasetId: "users",
					Columns: []models.ColumnConfig{
						{
							Column: "id",
						},
						{
							Column: "name",
						},
						{
							Column: "email",
						},
						{
							Column: "status",
						},
					},
				},
				Filters: models.FilterModel{
					LogicalOperator: constants.LogicalOperatorAnd,
					Conditions: []models.Filter{
						{
							Column: models.ColumnConfig{
								Column:   "status",
								Datatype: &dataTypeString,
							},
							Operator: constants.EqualOperator,
							Value:    "active",
						},
						{
							LogicalOperator: &constants.LogicalOperatorOr,
							Conditions: []models.Filter{
								{
									Column: models.ColumnConfig{
										Column:   "last_login",
										Datatype: &dataTypeDate,
									},
									Operator: constants.GreaterThanOperator,
									Value:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
								},
								{
									LogicalOperator: &constants.LogicalOperatorAnd,
									Conditions: []models.Filter{
										{
											Column: models.ColumnConfig{
												Column:   "email_verified",
												Datatype: &dataTypeBoolean,
											},
											Operator: constants.EqualOperator,
											Value:    true,
										},
										{
											Column: models.ColumnConfig{
												Column:   "subscription_tier",
												Datatype: &dataTypeString,
											},
											Operator: constants.EqualOperator,
											Value:    "premium",
										},
									},
								},
							},
						},
					},
				},
			},
			expectedSQL:    "SELECT id, name, email, status FROM {{.zamp_users}} WHERE ( status = 'active' ) AND (( last_login > '2024-01-01 00:00:00'::timestamp ) OR (( email_verified = true ) AND ( subscription_tier = 'premium' )))",
			expectedParams: map[string]interface{}{"zamp_users": "users"},
			expectError:    false,
		},
		{
			name: "Complex Group By with Aggregation",
			queryConfig: models.QueryConfig{
				TableConfig: models.TableConfig{
					DatasetId: "sales",
					Columns: []models.ColumnConfig{
						{
							Column: "region",
						},
					},
				},
				Aggregations: []models.Aggregation{
					{
						Column: models.ColumnConfig{
							Column: "revenue",
						},
						Function: constants.AggregationFunctionSum,
						Alias:    "total_revenue",
					},
					{
						Column: models.ColumnConfig{
							Column: "profit",
						},
						Function: constants.AggregationFunctionAvg,
						Alias:    "average_profit",
					},
				},
				GroupBy: []models.GroupBy{
					{
						Column: models.ColumnConfig{
							Column:   "region",
							Datatype: &dataTypeString,
						},
					},
				},
			},
			expectedSQL: "SELECT region, SUM(revenue) AS \"total_revenue\", AVG(profit) AS \"average_profit\" FROM {{.zamp_sales}} GROUP BY region",
			expectedParams: map[string]interface{}{
				"zamp_sales": "sales",
			},
			expectError: false,
		},
		{
			name: "Simple Group By with Aggregation",
			queryConfig: models.QueryConfig{
				TableConfig: models.TableConfig{
					DatasetId: "sales",
					Columns: []models.ColumnConfig{
						{
							Column: "country",
						},
					},
				},
				Aggregations: []models.Aggregation{
					{
						Column: models.ColumnConfig{
							Column: "revenue",
						},
						Function: constants.AggregationFunctionSum,
						Alias:    "total_revenue",
					},
					{
						Column: models.ColumnConfig{
							Column: "profit",
						},
						Function: constants.AggregationFunctionAvg,
						Alias:    "average_profit",
					},
				},
				GroupBy: []models.GroupBy{
					{
						Column: models.ColumnConfig{
							Column:   "region",
							Datatype: &dataTypeString,
						},
					},
				},
			},
			expectedSQL: "SELECT region, SUM(revenue) AS \"total_revenue\", AVG(profit) AS \"average_profit\" FROM {{.zamp_sales}} GROUP BY region",
			expectedParams: map[string]interface{}{
				"zamp_sales": "sales",
			},
			expectError: false,
		},
		{
			name: "Simple Array Group By with Aggregation",
			queryConfig: models.QueryConfig{
				TableConfig: models.TableConfig{
					DatasetId: "sales",
					Columns: []models.ColumnConfig{
						{
							Column: "country",
						},
					},
				},
				Aggregations: []models.Aggregation{
					{
						Column: models.ColumnConfig{
							Column: "revenue",
						},
						Function: constants.AggregationFunctionSum,
						Alias:    "total_revenue",
					},
					{
						Column: models.ColumnConfig{
							Column: "profit",
						},
						Function: constants.AggregationFunctionAvg,
						Alias:    "average_profit",
					},
				},
				GroupBy: []models.GroupBy{
					{
						Column: models.ColumnConfig{
							Column:   "region",
							Datatype: &dataTypeArray,
						},
					},
				},
			},
			expectedSQL: "SELECT unnest(region), SUM(revenue) AS \"total_revenue\", AVG(profit) AS \"average_profit\" FROM {{.zamp_sales}} GROUP BY unnest(region)",
			expectedParams: map[string]interface{}{
				"zamp_sales": "sales",
			},
			expectError: false,
		},
		{
			name: "Simple Array Group By with Alias",
			queryConfig: models.QueryConfig{
				TableConfig: models.TableConfig{
					DatasetId: "sales",
					Columns: []models.ColumnConfig{
						{
							Column: "country",
						},
					},
				},
				GroupBy: []models.GroupBy{
					{
						Column: models.ColumnConfig{
							Column:   "region",
							Datatype: &dataTypeArray,
							Alias:    &alias,
						},
					},
				},
			},
			expectedSQL: "SELECT unnest(region) AS \"alias\" FROM {{.zamp_sales}} GROUP BY \"alias\"",
			expectedParams: map[string]interface{}{
				"zamp_sales": "sales",
			},
			expectError: false,
		},
		{
			name: "Simple Array Group By with Aggregation and Filters",
			queryConfig: models.QueryConfig{
				TableConfig: models.TableConfig{
					DatasetId: "sales",
					Columns: []models.ColumnConfig{
						{
							Column: "country",
						},
					},
				},
				Filters: models.FilterModel{
					LogicalOperator: constants.LogicalOperatorAnd,
					Conditions: []models.Filter{
						{
							Column: models.ColumnConfig{
								Column:   "region",
								Datatype: &dataTypeArray,
							},
							Operator: constants.ArrayInOperator,
							Value:    []string{"Europe", "Asia"},
						},
					},
				},
				Aggregations: []models.Aggregation{
					{
						Column: models.ColumnConfig{
							Column: "revenue",
						},
						Function: constants.AggregationFunctionSum,
						Alias:    "total_revenue",
					},
					{
						Column: models.ColumnConfig{
							Column: "profit",
						},
						Function: constants.AggregationFunctionAvg,
						Alias:    "average_profit",
					},
				},
				GroupBy: []models.GroupBy{
					{
						Column: models.ColumnConfig{
							Column:   "region",
							Datatype: &dataTypeArray,
						},
					},
				},
			},
			expectedSQL: "SELECT unnest(region), SUM(revenue) AS \"total_revenue\", AVG(profit) AS \"average_profit\" FROM {{.zamp_sales}} WHERE ( ( LOWER(ARRAY_TO_STRING(region, ',')) = 'europe' ) OR ( LOWER(ARRAY_TO_STRING(region, ',')) = 'asia' ) ) GROUP BY unnest(region)",
			expectedParams: map[string]interface{}{
				"zamp_sales": "sales",
			},
			expectError: false,
		},
		{
			name: "Simple Array Group Error",
			queryConfig: models.QueryConfig{
				TableConfig: models.TableConfig{
					DatasetId: "sales",
					Columns: []models.ColumnConfig{
						{
							Column: "country",
						},
					},
				},
				GroupBy: []models.GroupBy{
					{
						Column: models.ColumnConfig{
							Column: "region",
						},
					},
				},
			},
			expectedSQL:    "",
			expectedParams: map[string]interface{}{},
			expectError:    true,
		},
		{
			name: "Simple Aggregation",
			queryConfig: models.QueryConfig{
				TableConfig: models.TableConfig{
					DatasetId: "sales",
					Columns: []models.ColumnConfig{
						{
							Column: "region",
						},
					},
				},
				Aggregations: []models.Aggregation{
					{
						Column: models.ColumnConfig{
							Column: "revenue",
						},
						Function: constants.AggregationFunctionSum,
						Alias:    "total_revenue",
					},
					{
						Column: models.ColumnConfig{
							Column: "profit",
						},
						Function: constants.AggregationFunctionAvg,
						Alias:    "average_profit",
					},
				},
			},
			expectedSQL: "SELECT SUM(revenue) AS \"total_revenue\", AVG(profit) AS \"average_profit\" FROM {{.zamp_sales}}",
			expectedParams: map[string]interface{}{
				"zamp_sales": "sales",
			},
			expectError: false,
		},
		{
			name: "Error Case - Unsupported Operator",
			queryConfig: models.QueryConfig{
				TableConfig: models.TableConfig{
					DatasetId: "logs",
					Columns: []models.ColumnConfig{
						{
							Column: "timestamp",
						},
						{
							Column: "level",
						},
					},
				},
				Filters: models.FilterModel{
					LogicalOperator: constants.LogicalOperatorAnd,
					Conditions: []models.Filter{
						{
							Column: models.ColumnConfig{
								Column:   "level",
								Datatype: &dataTypeString,
							},
							Operator: "INVALID_OPERATOR",
							Value:    "error",
						},
					},
				},
			},
			expectedSQL:    "",
			expectedParams: map[string]interface{}{},
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			sql, params, err := s.service.ToSQL(ctx, tc.queryConfig)

			if tc.expectError {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedSQL, sql)
				assert.Equal(s.T(), tc.expectedParams, params)
			}
		})
	}
}

func (s *ServiceTestSuite) TestToFilterSQL() {
	ctx := context.Background()

	dataTypeString := dataplatformConstants.StringDataType
	dataTypeDecimal := dataplatformConstants.DecimalDataType
	dataTypeDate := dataplatformConstants.DateDataType
	dataTypeBoolean := dataplatformConstants.BooleanDataType
	dataTypeArray := dataplatformConstants.ArrayOfStringDataType

	testCases := []struct {
		name           string
		filterConfig   models.FilterModel
		expectedSQL    string
		expectedParams map[string]interface{}
		expectError    bool
	}{
		{
			name: "Simple AND condition",
			filterConfig: models.FilterModel{
				LogicalOperator: constants.LogicalOperatorAnd,
				Conditions: []models.Filter{
					{
						Column: models.ColumnConfig{
							Column:   "status",
							Datatype: &dataTypeString,
						},
						Operator: constants.EqualOperator,
						Value:    "active",
					},
					{
						Column: models.ColumnConfig{
							Column:   "age",
							Datatype: &dataTypeDecimal,
						},
						Operator: constants.GreaterThanOperator,
						Value:    18,
					},
				},
			},
			expectedSQL:    "( status = 'active' ) AND ( age > 18 )",
			expectedParams: map[string]interface{}{},
			expectError:    false,
		},
		{
			name: "Nested OR conditions",
			filterConfig: models.FilterModel{
				LogicalOperator: constants.LogicalOperatorOr,
				Conditions: []models.Filter{
					{
						Column: models.ColumnConfig{
							Column:   "created_at",
							Datatype: &dataTypeDate,
						},
						Operator: constants.GreaterThanOperator,
						Value:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						LogicalOperator: &constants.LogicalOperatorAnd,
						Conditions: []models.Filter{
							{
								Column: models.ColumnConfig{
									Column:   "is_verified",
									Datatype: &dataTypeBoolean,
								},
								Operator: constants.EqualOperator,
								Value:    true,
							},
							{
								Column: models.ColumnConfig{
									Column:   "role",
									Datatype: &dataTypeString,
								},
								Operator: constants.EqualOperator,
								Value:    "admin",
							},
						},
					},
				},
			},
			expectedSQL:    "( created_at > '2024-01-01 00:00:00'::timestamp ) OR (( is_verified = true ) AND ( role = 'admin' ))",
			expectedParams: map[string]interface{}{},
			expectError:    false,
		},
		{
			name: "Array operations",
			filterConfig: models.FilterModel{
				LogicalOperator: constants.LogicalOperatorAnd,
				Conditions: []models.Filter{
					{
						Column: models.ColumnConfig{
							Column:   "tags",
							Datatype: &dataTypeArray,
						},
						Operator: constants.ArrayInOperator,
						Value:    []string{"featured", "new"},
					},
				},
			},
			expectedSQL:    "( ( LOWER(ARRAY_TO_STRING(tags, ',')) = 'featured' ) OR ( LOWER(ARRAY_TO_STRING(tags, ',')) = 'new' ) )",
			expectedParams: map[string]interface{}{},
			expectError:    false,
		},
		{
			name: "Invalid operator",
			filterConfig: models.FilterModel{
				LogicalOperator: constants.LogicalOperatorAnd,
				Conditions: []models.Filter{
					{
						Column: models.ColumnConfig{
							Column:   "status",
							Datatype: &dataTypeString,
						},
						Operator: "INVALID_OPERATOR",
						Value:    "active",
					},
				},
			},
			expectedSQL:    "",
			expectedParams: nil,
			expectError:    true,
		},
		{
			name: "Invalid logical operator",
			filterConfig: models.FilterModel{
				LogicalOperator: "INVALID_LOGICAL_OPERATOR",
				Conditions: []models.Filter{
					{
						Column: models.ColumnConfig{
							Column:   "status",
							Datatype: &dataTypeString,
						},
						Operator: constants.EqualOperator,
						Value:    "active",
					},
					{
						Column: models.ColumnConfig{
							Column:   "age",
							Datatype: &dataTypeDecimal,
						},
						Operator: constants.GreaterThanOperator,
						Value:    18,
					},
				},
			},
			expectedSQL:    "",
			expectedParams: nil,
			expectError:    true,
		},
		{
			name: "Single value in array",
			filterConfig: models.FilterModel{
				LogicalOperator: constants.LogicalOperatorAnd,
				Conditions: []models.Filter{
					{
						Column: models.ColumnConfig{
							Column:   "tags",
							Datatype: &dataTypeArray,
						},
						Operator: constants.ArrayContainsOperator,
						Value:    []string{"test"},
					},
				},
			},
			expectedSQL:    "( ( LOWER(ARRAY_TO_STRING(tags, ',')) LIKE '%test%' ) )",
			expectedParams: map[string]interface{}{},
			expectError:    false,
		},
		{
			name: "Multiple values in array",
			filterConfig: models.FilterModel{
				LogicalOperator: constants.LogicalOperatorAnd,
				Conditions: []models.Filter{
					{
						Column: models.ColumnConfig{
							Column:   "categories",
							Datatype: &dataTypeArray,
						},
						Operator: constants.ArrayContainsOperator,
						Value:    []string{"food", "drink"},
					},
				},
			},
			expectedSQL:    "( ( LOWER(ARRAY_TO_STRING(categories, ',')) LIKE '%food%' ) OR ( LOWER(ARRAY_TO_STRING(categories, ',')) LIKE '%drink%' ) )",
			expectedParams: map[string]interface{}{},
			expectError:    false,
		},
		{
			name: "Invalid value type",
			filterConfig: models.FilterModel{
				LogicalOperator: constants.LogicalOperatorAnd,
				Conditions: []models.Filter{
					{
						Column: models.ColumnConfig{
							Column:   "tags",
							Datatype: &dataTypeArray,
						},
						Operator: constants.ArrayContainsOperator,
						Value:    "not_an_array",
					},
				},
			},
			expectedSQL:    "",
			expectedParams: map[string]interface{}{},
			expectError:    true,
		},
		{
			name: "Same column multiple values",
			filterConfig: models.FilterModel{
				LogicalOperator: constants.LogicalOperatorAnd,
				Conditions: []models.Filter{
					{
						Column: models.ColumnConfig{
							Column:   "tags",
							Datatype: &dataTypeArray,
						},
						Operator: constants.ArrayContainsOperator,
						Value:    []string{"urgent", "important", "follow-up"},
					},
				},
			},
			expectedSQL:    "( ( LOWER(ARRAY_TO_STRING(tags, ',')) LIKE '%urgent%' ) OR ( LOWER(ARRAY_TO_STRING(tags, ',')) LIKE '%important%' ) OR ( LOWER(ARRAY_TO_STRING(tags, ',')) LIKE '%follow-up%' ) )",
			expectedParams: map[string]interface{}{},
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			sql, params, err := s.service.ToFilterSQL(ctx, tc.filterConfig)

			if tc.expectError {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedSQL, sql)
				assert.Equal(s.T(), tc.expectedParams, params)
			}
		})
	}
}

func (s *ServiceTestSuite) TestBuildEqualClause() {
	testCases := []struct {
		name        string
		column      string
		value       interface{}
		expectedSQL string
		expectError bool
	}{
		{
			name:        "Valid string value",
			column:      "status",
			value:       "active",
			expectedSQL: "( status = active )",
			expectError: false,
		},
		{
			name:        "Non-string value",
			column:      "age",
			value:       42,
			expectedSQL: "",
			expectError: true,
		},
		{
			name:        "Empty string value",
			column:      "name",
			value:       "",
			expectedSQL: "( name =  )",
			expectError: false,
		},
		{
			name:        "Column with special characters",
			column:      "user_status",
			value:       "pending",
			expectedSQL: "( user_status = pending )",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			sql, err := s.service.(*queryBuilder).buildEqualClause(tc.column, tc.value)

			if tc.expectError {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), "", sql)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedSQL, sql)
			}
		})
	}
}

func (s *ServiceTestSuite) TestBuildNotEqualClause() {
	testCases := []struct {
		name        string
		column      string
		value       interface{}
		expectedSQL string
		expectError bool
	}{
		{
			name:        "Valid string value",
			column:      "status",
			value:       "inactive",
			expectedSQL: "( status != inactive )",
			expectError: false,
		},
		{
			name:        "Non-string value",
			column:      "age",
			value:       42,
			expectedSQL: "",
			expectError: true,
		},
		{
			name:        "Empty string value",
			column:      "name",
			value:       "",
			expectedSQL: "( name !=  )",
			expectError: false,
		},
		{
			name:        "Column with special characters",
			column:      "user_status",
			value:       "pending",
			expectedSQL: "( user_status != pending )",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			sql, err := s.service.(*queryBuilder).buildNotEqualClause(tc.column, tc.value)

			if tc.expectError {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), "", sql)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedSQL, sql)
			}
		})
	}
}

func (s *ServiceTestSuite) TestWindowFunctions() {
	ctx := context.Background()
	dataTypeString := dataplatformConstants.StringDataType

	testCases := []struct {
		name           string
		queryConfig    models.QueryConfig
		expectedSQL    string
		expectedParams map[string]interface{}
		expectError    bool
	}{
		{
			name: "Basic Window Function",
			queryConfig: models.QueryConfig{
				TableConfig: models.TableConfig{
					DatasetId: "transactions",
					Columns: []models.ColumnConfig{
						{Column: "account_number"},
						{Column: "currency_code"},
						{Column: "_time_stamp_utc"},
						{Column: "balance_type"},
					},
				},
				Windows: []models.WindowConfig{
					{
						Function: models.WindowFunctionRowNumber,
						PartitionBy: []models.ColumnConfig{
							{Column: "account_number"},
							{Column: "currency_code"},
						},
						OrderBy: []models.OrderBy{
							{
								Column: models.ColumnConfig{Column: "_time_stamp_utc"},
								Order:  constants.OrderAsc,
							},
						},
						Alias: "rn",
					},
				},
			},
			expectedSQL: "SELECT account_number, currency_code, _time_stamp_utc, balance_type, ROW_NUMBER() OVER (  PARTITION BY account_number, currency_code ORDER BY _time_stamp_utc ASC ) AS \"rn\" FROM {{.zamp_transactions}}",
			expectedParams: map[string]interface{}{
				"zamp_transactions": "transactions",
			},
			expectError: false,
		},
		{
			name: "Window Function with Subquery and Complex Filters",
			queryConfig: models.QueryConfig{
				Subquery: &models.QueryConfig{
					TableConfig: models.TableConfig{
						DatasetId: "transactions",
						Columns: []models.ColumnConfig{
							{Column: "account_number"},
							{Column: "currency_code"},
							{Column: "DATE_TRUNC('month', _time_stamp_utc) AS \"month\""},
							{Column: "balance_type"},
						},
					},
					Windows: []models.WindowConfig{
						{
							Function: models.WindowFunctionRowNumber,
							PartitionBy: []models.ColumnConfig{
								{Column: "account_number"},
								{Column: "currency_code"},
								{Column: "DATE_TRUNC('month', _time_stamp_utc)"},
							},
							OrderBy: []models.OrderBy{
								{
									Column: models.ColumnConfig{Column: "_time_stamp_utc"},
									Order:  constants.OrderAsc,
								},
							},
							Alias: "rn",
						},
					},
					Filters: models.FilterModel{
						LogicalOperator: constants.LogicalOperatorAnd,
						Conditions: []models.Filter{
							{
								Column:   models.ColumnConfig{Column: "balance_type", Datatype: &dataTypeString},
								Operator: constants.EqualOperator,
								Value:    "opening",
							},
							{
								Column:   models.ColumnConfig{Column: "account_number", Datatype: &dataTypeString},
								Operator: constants.EqualOperator,
								Value:    "8912672444",
							},
						},
					},
				},
				TableConfig: models.TableConfig{
					Columns: []models.ColumnConfig{
						{Column: "account_number"},
						{Column: "currency_code"},
						{Column: "month"},
						{Column: "balance_type"},
					},
				},
				Filters: models.FilterModel{
					LogicalOperator: constants.LogicalOperatorAnd,
					Conditions: []models.Filter{
						{
							Column:   models.ColumnConfig{Column: "rn", Datatype: &dataTypeString},
							Operator: constants.EqualOperator,
							Value:    "1",
						},
					},
				},
				GroupBy: []models.GroupBy{
					{Column: models.ColumnConfig{Column: "account_number", Datatype: &dataTypeString}},
					{Column: models.ColumnConfig{Column: "currency_code", Datatype: &dataTypeString}},
					{Column: models.ColumnConfig{Column: "month", Datatype: &dataTypeString}},
					{Column: models.ColumnConfig{Column: "balance_type", Datatype: &dataTypeString}},
				},
			},
			expectedSQL: "SELECT account_number, currency_code, month, balance_type FROM ( SELECT account_number, currency_code, DATE_TRUNC('month', _time_stamp_utc) AS \"month\", balance_type, ROW_NUMBER() OVER (  PARTITION BY account_number, currency_code, DATE_TRUNC('month', _time_stamp_utc) ORDER BY _time_stamp_utc ASC ) AS \"rn\" FROM {{.zamp_transactions}} WHERE ( balance_type = 'opening' ) AND ( account_number = '8912672444' ) ) subquery WHERE ( rn = '1' ) GROUP BY account_number, currency_code, month, balance_type",
			expectedParams: map[string]interface{}{
				"zamp_transactions": "transactions",
			},
			expectError: false,
		},
		{
			name: "Multiple Window Functions",
			queryConfig: models.QueryConfig{
				TableConfig: models.TableConfig{
					DatasetId: "transactions",
					Columns: []models.ColumnConfig{
						{Column: "account_number"},
						{Column: "amount"},
						{Column: "_time_stamp_utc"},
					},
				},
				Windows: []models.WindowConfig{
					{
						Function: models.WindowFunctionRowNumber,
						PartitionBy: []models.ColumnConfig{
							{Column: "account_number"},
						},
						OrderBy: []models.OrderBy{
							{
								Column: models.ColumnConfig{Column: "_time_stamp_utc"},
								Order:  constants.OrderAsc,
							},
						},
						Alias: "transaction_order",
					},
					{
						Function: models.WindowFunctionRowNumber,
						PartitionBy: []models.ColumnConfig{
							{Column: "account_number"},
						},
						OrderBy: []models.OrderBy{
							{
								Column: models.ColumnConfig{Column: "amount"},
								Order:  constants.OrderDesc,
							},
						},
						Alias: "amount_rank",
					},
				},
			},
			expectedSQL: "SELECT account_number, amount, _time_stamp_utc, ROW_NUMBER() OVER (  PARTITION BY account_number ORDER BY _time_stamp_utc ASC ) AS \"transaction_order\", ROW_NUMBER() OVER (  PARTITION BY account_number ORDER BY amount DESC ) AS \"amount_rank\" FROM {{.zamp_transactions}}",
			expectedParams: map[string]interface{}{
				"zamp_transactions": "transactions",
			},
			expectError: false,
		},
		{
			name: "Window Function with Order By and Pagination",
			queryConfig: models.QueryConfig{
				TableConfig: models.TableConfig{
					DatasetId: "transactions",
					Columns: []models.ColumnConfig{
						{Column: "account_number"},
						{Column: "amount"},
					},
				},
				Windows: []models.WindowConfig{
					{
						Function: models.WindowFunctionRowNumber,
						PartitionBy: []models.ColumnConfig{
							{Column: "account_number"},
						},
						OrderBy: []models.OrderBy{
							{
								Column: models.ColumnConfig{Column: "amount"},
								Order:  constants.OrderDesc,
							},
						},
						Alias: "amount_rank",
					},
				},
				OrderBy: []models.OrderBy{
					{
						Column: models.ColumnConfig{Column: "amount_rank"},
						Order:  constants.OrderAsc,
					},
				},
				Pagination: &models.Pagination{
					Page:     1,
					PageSize: 10,
				},
			},
			expectedSQL: "SELECT account_number, amount, ROW_NUMBER() OVER (  PARTITION BY account_number ORDER BY amount DESC ) AS \"amount_rank\" FROM {{.zamp_transactions}} ORDER BY amount_rank ASC LIMIT 10 OFFSET 0",
			expectedParams: map[string]interface{}{
				"zamp_transactions": "transactions",
			},
			expectError: false,
		},
		{
			name: "Nested Subqueries with Window Functions",
			queryConfig: models.QueryConfig{
				Subquery: &models.QueryConfig{
					TableConfig: models.TableConfig{
						DatasetId: "transactions",
						Columns: []models.ColumnConfig{
							{Column: "account_number"},
							{Column: "amount"},
							{Column: "_time_stamp_utc"},
						},
					},
					Windows: []models.WindowConfig{
						{
							Function: models.WindowFunctionRowNumber,
							PartitionBy: []models.ColumnConfig{
								{Column: "account_number"},
							},
							OrderBy: []models.OrderBy{
								{
									Column: models.ColumnConfig{Column: "_time_stamp_utc"},
									Order:  constants.OrderDesc,
								},
							},
							Alias: "latest_rank",
						},
					},
				},
				TableConfig: models.TableConfig{
					Columns: []models.ColumnConfig{
						{Column: "account_number"},
						{Column: "amount"},
					},
				},
				Filters: models.FilterModel{
					LogicalOperator: constants.LogicalOperatorAnd,
					Conditions: []models.Filter{
						{
							Column:   models.ColumnConfig{Column: "latest_rank", Datatype: &dataTypeString},
							Operator: constants.LessThanOrEqualOperator,
							Value:    "3",
						},
					},
				},
				OrderBy: []models.OrderBy{
					{
						Column: models.ColumnConfig{Column: "latest_rank"},
						Order:  constants.OrderAsc,
					},
				},
			},
			expectedSQL: "SELECT account_number, amount FROM ( SELECT account_number, amount, _time_stamp_utc, ROW_NUMBER() OVER (  PARTITION BY account_number ORDER BY _time_stamp_utc DESC ) AS \"latest_rank\" FROM {{.zamp_transactions}} ) subquery WHERE ( latest_rank <= '3' ) ORDER BY latest_rank ASC",
			expectedParams: map[string]interface{}{
				"zamp_transactions": "transactions",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			sql, params, err := s.service.ToSQL(ctx, tc.queryConfig)

			if tc.expectError {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedSQL, sql)
				assert.Equal(s.T(), tc.expectedParams, params)
			}
		})
	}
}

package models

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/Zampfi/application-platform/services/api/db/pgclient"
// 	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func TestWidgetInstance_TableName(t *testing.T) {
// 	t.Parallel()
// 	widgetInstance := WidgetInstance{}
// 	assert.Equal(t, "widget_instances", widgetInstance.TableName())
// }

// func TestWidgetInstance_GetQueryFilters(t *testing.T) {
// 	t.Parallel()

// 	userId := uuid.New()
// 	org1ID := uuid.New()
// 	org2ID := uuid.New()
// 	orgIDs := []uuid.UUID{org1ID, org2ID}
// 	instanceID := uuid.New()

// 	tests := []struct {
// 		name            string
// 		userId          uuid.UUID
// 		organizationIDs []uuid.UUID
// 		setupMock       func(mock sqlmock.Sqlmock)
// 		wantErr         bool
// 	}{
// 		{
// 			name:            "user and organization access",
// 			userId:          userId,
// 			organizationIDs: orgIDs,
// 			setupMock: func(mock sqlmock.Sqlmock) {
// 				rows := sqlmock.NewRows([]string{"widget_instance_id"}).AddRow(instanceID)
// 				mock.ExpectQuery(`SELECT \* FROM "widget_instances" WHERE EXISTS \(.*\)`).
// 					WithArgs(userId).
// 					WillReturnRows(rows)
// 			},
// 		},
// 		{
// 			name:            "empty organization list",
// 			userId:          userId,
// 			organizationIDs: []uuid.UUID{},
// 			setupMock: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(`SELECT \* FROM "widget_instances" WHERE EXISTS \(.*\)`).
// 					WithArgs(userId).
// 					WillReturnRows(sqlmock.NewRows([]string{"widget_instance_id"}))
// 			},
// 		},
// 		{
// 			name:            "only user access",
// 			userId:          userId,
// 			organizationIDs: nil,
// 			setupMock: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectQuery(`SELECT \* FROM "widget_instances" WHERE EXISTS \(.*\)`).
// 					WithArgs(userId).
// 					WillReturnRows(sqlmock.NewRows([]string{"widget_instance_id"}))
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			db, mock := setupTestDB(t)
// 			widgetInstance := &WidgetInstance{ID: instanceID}
// 			tt.setupMock(mock)

// 			query := widgetInstance.GetQueryFilters(db.Model(&WidgetInstance{}), tt.userId, tt.organizationIDs)

// 			// Verify query executes without error
// 			var results []WidgetInstance
// 			err := query.Find(&results).Error
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			assert.NoError(t, mock.ExpectationsWereMet())
// 		})
// 	}
// }

// func TestStructImplementsBaseModel_WidgetInstance(t *testing.T) {
// 	var _ pgclient.BaseModel = &WidgetInstance{}
// }

// func TestWidgetInstance_BeforeCreate(t *testing.T) {
// 	t.Parallel()

// 	tests := []struct {
// 		name      string
// 		setupMock func(mock sqlmock.Sqlmock, userId uuid.UUID)
// 		setupCtx  func() (context.Context, uuid.UUID)
// 		wantErr   bool
// 		errMsg    string
// 	}{
// 		{
// 			name: "successful creation with admin privilege",
// 			setupCtx: func() (context.Context, uuid.UUID) {
// 				userId := uuid.New()
// 				ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{})
// 				return ctx, userId
// 			},
// 			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID) {
// 				mock.ExpectQuery(`SELECT \* FROM "sheets" WHERE EXISTS \(.*\)`).
// 					WithArgs(sqlmock.AnyArg(), userId, 1).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "sheet_id", "page_id"}).
// 						AddRow(uuid.New(), uuid.New(), uuid.New()))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "failure - no user ID in context",
// 			setupCtx: func() (context.Context, uuid.UUID) {
// 				return context.Background(), uuid.Nil
// 			},
// 			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID) {
// 				// No mock expectations needed as it should fail before DB query
// 			},
// 			wantErr: true,
// 			errMsg:  "no user id found in context",
// 		},
// 		{
// 			name: "failure - no admin privilege found",
// 			setupCtx: func() (context.Context, uuid.UUID) {
// 				userId := uuid.New()
// 				ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{})
// 				return ctx, userId
// 			},
// 			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID) {
// 				mock.ExpectQuery(`SELECT \* FROM "sheets" WHERE EXISTS \(.*\)`).
// 					WithArgs(sqlmock.AnyArg(), userId, 1).
// 					WillReturnRows(sqlmock.NewRows([]string{"id", "sheet_id", "page_id"}))
// 			},
// 			wantErr: true,
// 			errMsg:  "page access forbidden",
// 		},
// 		{
// 			name: "failure - database error",
// 			setupCtx: func() (context.Context, uuid.UUID) {
// 				userId := uuid.New()
// 				ctx := apicontext.AddAuthToContext(context.Background(), "role", userId, []uuid.UUID{})
// 				return ctx, userId
// 			},
// 			setupMock: func(mock sqlmock.Sqlmock, userId uuid.UUID) {
// 				mock.ExpectQuery(`SELECT \* FROM "sheets" WHERE EXISTS \(.*\)`).
// 					WithArgs(sqlmock.AnyArg(), userId, 1).
// 					WillReturnError(fmt.Errorf("database error"))
// 			},
// 			wantErr: true,
// 			errMsg:  "database error",
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			db, mock := setupTestDB(t)

// 			ctx, userId := tt.setupCtx()
// 			db = db.WithContext(ctx)

// 			widgetInstance := &WidgetInstance{
// 				ID:        uuid.New(),
// 				SheetID:   uuid.New(),
// 				CreatedAt: time.Now(),
// 				UpdatedAt: time.Now(),
// 			}

// 			tt.setupMock(mock, userId)

// 			err := widgetInstance.BeforeCreate(db)

// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				if tt.errMsg != "" {
// 					assert.Contains(t, err.Error(), tt.errMsg)
// 				}
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 			assert.NoError(t, mock.ExpectationsWereMet())
// 		})
// 	}
// }

// func TestWidgetInstance_FromDB(t *testing.T) {
// 	sheetID := uuid.New()
// 	instanceID := uuid.New()

// 	dataMappings := DataMappings{
// 		Version: DataMappingVersion1,
// 		Mappings: []DataMappingFields{
// 			{
// 				DatasetID: "dataset1",
// 				Fields: map[string][]Field{
// 					"metrics": {{
// 						Column:      "revenue",
// 						Aggregation: "sum",
// 						Type:        "number",
// 						FieldType:   "metric",
// 					}},
// 				},
// 			},
// 		},
// 	}

// 	mappingsJSON, err := json.Marshal(dataMappings)
// 	require.NoError(t, err)

// 	dbModel := &dbmodels.WidgetInstance{
// 		ID:           instanceID,
// 		SheetID:      sheetID,
// 		Title:        "Test Widget",
// 		DataMappings: mappingsJSON,
// 		WidgetType:   "chart",
// 	}

// 	wi := &WidgetInstance{}
// 	err = wi.FromDB(dbModel)
// 	require.NoError(t, err)

// 	assert.Equal(t, instanceID, wi.ID)
// 	assert.Equal(t, "chart", wi.WidgetType)
// 	assert.Equal(t, sheetID, wi.SheetID)
// 	assert.Equal(t, "Test Widget", wi.Title)
// 	assert.Equal(t, DataMappingVersion1, wi.DataMappings.Version)
// 	assert.Len(t, wi.DataMappings.Mappings, 1)
// }

// func TestWidgetInstance_ToDB(t *testing.T) {
// 	sheetID := uuid.New()
// 	instanceID := uuid.New()

// 	wi := &WidgetInstance{
// 		ID:         instanceID,
// 		WidgetType: "chart",
// 		SheetID:    sheetID,
// 		Title:      "Test Widget",
// 		DataMappings: DataMappings{
// 			Version: DataMappingVersion1,
// 			Mappings: []DataMappingFields{
// 				{
// 					DatasetID: "dataset1",
// 					Fields: map[string][]Field{
// 						"metrics": {{
// 							Column:      "revenue",
// 							Aggregation: "sum",
// 							Type:        "number",
// 							FieldType:   "metric",
// 						}},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	dbModel, err := wi.ToDB()
// 	require.NoError(t, err)

// 	assert.Equal(t, instanceID, dbModel.ID)
// 	assert.Equal(t, "chart", dbModel.WidgetType)
// 	assert.Equal(t, sheetID, dbModel.SheetID)
// 	assert.Equal(t, "Test Widget", dbModel.Title)

// 	var dataMappings DataMappings
// 	err = json.Unmarshal(dbModel.DataMappings, &dataMappings)
// 	require.NoError(t, err)

// 	assert.Equal(t, DataMappingVersion1, dataMappings.Version)
// 	assert.Len(t, dataMappings.Mappings, 1)
// }

// func TestCreateWidgetInstancePayload_ToModel(t *testing.T) {
// 	sheetID := uuid.New()

// 	dataMappings := map[string]interface{}{
// 		"data_mappings": map[string]interface{}{
// 			"version": "1",
// 			"mappings": []map[string]interface{}{
// 				{
// 					"dataset_id": "dataset1",
// 					"fields": map[string]interface{}{
// 						"metrics": []map[string]interface{}{
// 							{
// 								"column":      "revenue",
// 								"aggregation": "sum",
// 								"type":        "number",
// 								"field_type":  "metric",
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	dataMappingsJSON, err := json.Marshal(dataMappings)
// 	require.NoError(t, err)

// 	payload := &CreateWidgetInstancePayload{
// 		WidgetType:   "chart",
// 		SheetId:      sheetID.String(),
// 		Title:        "Test Widget",
// 		DataMappings: string(dataMappingsJSON),
// 	}

// 	model, err := payload.ToModel()
// 	require.NoError(t, err)

// 	assert.Equal(t, "chart", model.WidgetType)
// 	assert.Equal(t, sheetID, model.SheetID)
// 	assert.Equal(t, "Test Widget", model.Title)
// 	assert.Equal(t, DataMappingVersion1, model.DataMappings.Version)
// 	assert.Len(t, model.DataMappings.Mappings, 1)
// }

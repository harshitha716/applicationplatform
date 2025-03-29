package pgclient

import (
	"context"
	"fmt"
	"testing"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type userBaseModel struct {
	ID uuid.UUID `json:"user_id" gorm:"column:user_id"`
}

func (u *userBaseModel) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(`"test_user_id" = ?`, userId)
}
func (u *userBaseModel) GetUpdateFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(`"test_user_id" = ?`, userId)
}
func (u *userBaseModel) GetDeleteFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(`"test_user_id" = ?`, userId)
}
func (u *userBaseModel) BeforeCreate(db *gorm.DB) error {
	return nil
}

func (u userBaseModel) TableName() string {
	return "user_base_models"
}

type userBaseModelPtr struct {
	ID uuid.UUID `json:"user_id" gorm:"column:user_id"`
}

func (u *userBaseModelPtr) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(`"test_user_id" = ?`, userId)
}
func (u *userBaseModelPtr) GetUpdateFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(`"test_user_id" = ?`, userId)
}
func (u *userBaseModelPtr) GetDeleteFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
	return db.Where(`"test_user_id" = ?`, userId)
}
func (u *userBaseModelPtr) BeforeCreate(db *gorm.DB) error {
	return nil
}

type userNotBaseModel struct {
	ID uuid.UUID `json:"user_id" gorm:"column:user_id"`
}

type userNotBaseModelPtr struct {
	ID uuid.UUID `json:"user_id" gorm:"column:user_id"`
}

func GetMockPgClient() (*PostgresClient, sqlmock.Sqlmock) {

	client := &PostgresClient{}

	sqldb, mock, _ := sqlmock.New()

	gormdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	client.DB = gormdb

	return client, mock
}

func GetMockPGClientWithHooks() (*PostgresClient, sqlmock.Sqlmock) {
	client, mock := GetMockPgClient()

	userId := uuid.New()
	orgId := uuid.New()

	ctx := context.Background()
	ctx = apicontext.AddAuthToContext(ctx, "user", userId, []uuid.UUID{orgId})
	client.DB = client.DB.WithContext(ctx)

	registerAccessControlHooks(client.DB)

	return client, mock
}

func TestRegisterAccessControlHooks_NoUser(t *testing.T) {
	client, _ := GetMockPgClient()

	context := context.Background()

	client.DB = client.DB.WithContext(context)

	registerAccessControlHooks(client.DB)

	_, err := client.DB.Model(&userBaseModel{}).Select("*").Rows()

	assert.NotNil(t, err)
}

func TestRegisterAccessControlHooks_NoOrg(t *testing.T) {
	client, _ := GetMockPgClient()

	registerAccessControlHooks(client.DB)

	userID := uuid.New()
	ctx := context.Background()
	ctx = apicontext.AddAuthToContext(ctx, "user", userID, []uuid.UUID{})
	client.DB = client.DB.WithContext(ctx)

	queryDb := client.DB.Model([]userBaseModel{})

	queryDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
		result := tx.Select("*").Find(&[]userBaseModel{})
		assert.Nil(t, result.Error)
		assert.Equal(t, `SELECT * FROM "user_base_models" WHERE "test_user_id" = $1`, result.Statement.SQL.String())
		return result
	})

}

func TestInjectFilterFromContext_NotBaseModel(t *testing.T) {
	client, _ := GetMockPgClient()

	queryDb := client.DB.Model(userNotBaseModel{})
	injectFilterFromContext(queryDb)

	sql := queryDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
		result := tx.Select("*")
		assert.NotNil(t, result.Error)
		assert.Equal(t, "model must implement BaseModel or be a slice of BaseModels", result.Error.Error())
		return result
	})
	assert.Equal(t, "", sql)
}

func TestInjectFilterFromContext_NotBaseModelPtr(t *testing.T) {
	client, _ := GetMockPgClient()

	queryDb := client.DB.Model(&userNotBaseModelPtr{})
	injectFilterFromContext(queryDb)

	sql := queryDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
		result := tx.Select("*")
		assert.NotNil(t, result.Error)
		assert.Equal(t, "model must implement BaseModel or be a slice of BaseModels", result.Error.Error())
		return result
	})
	assert.Equal(t, "", sql)
}

func TestInjectFilterFromContext_NonBaseModel_Slice(t *testing.T) {
	client, _ := GetMockPgClient()

	ctx := context.Background()
	client.DB = client.DB.WithContext(ctx)

	queryDb := client.DB.Model([]userNotBaseModel{})
	injectFilterFromContext(queryDb)

	_ = queryDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
		result := tx.Select("*").Find([]userNotBaseModel{})
		assert.Equal(t, "model must implement BaseModel or be a slice of BaseModels", result.Error.Error())
		return result
	})
}

func TestInjectFilterFromContext_BaseModelPtr(t *testing.T) {
	client, _ := GetMockPgClient()

	userID := uuid.New()
	orgId := uuid.New()
	ctx := context.Background()
	ctx = apicontext.AddAuthToContext(ctx, "user", userID, []uuid.UUID{orgId})
	client.DB = client.DB.WithContext(ctx)

	queryDb := client.DB.Model(&userBaseModelPtr{})
	injectFilterFromContext(queryDb)

	_ = queryDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
		result := tx.Select("*").Take(&userBaseModel{})
		assert.Nil(t, result.Error)
		assert.Equal(t, `SELECT * FROM "user_base_model_ptrs" WHERE "test_user_id" = $1 LIMIT $2`, result.Statement.SQL.String())
		return result
	})
}

func TestInjectFilterFromContext_BaseModel(t *testing.T) {
	client, _ := GetMockPgClient()

	userID := uuid.New()
	orgId := uuid.New()
	ctx := context.Background()
	ctx = apicontext.AddAuthToContext(ctx, "user", userID, []uuid.UUID{orgId})
	client.DB = client.DB.WithContext(ctx)

	queryDb := client.DB.Model(userBaseModel{})
	injectFilterFromContext(queryDb)

	_ = queryDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
		result := tx.Select("*").Take(&userBaseModel{})
		assert.Nil(t, result.Error)
		assert.Equal(t, `SELECT * FROM "user_base_models" WHERE "test_user_id" = $1 LIMIT $2`, result.Statement.SQL.String())
		return result
	})
}

func TestInjectFilterFromContext_BaseModel_Slice(t *testing.T) {
	client, _ := GetMockPgClient()

	userID := uuid.New()
	orgId := uuid.New()
	ctx := context.Background()
	ctx = apicontext.AddAuthToContext(ctx, "user", userID, []uuid.UUID{orgId})
	client.DB = client.DB.WithContext(ctx)

	queryDb := client.DB.Model([]userBaseModel{})
	startTime := time.Now()
	injectFilterFromContext(queryDb)
	duration := time.Since(startTime)
	fmt.Println("Duration for complete reflection", duration)

	_ = queryDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
		result := tx.Select("*").Find([]userBaseModel{})
		assert.Nil(t, result.Error)
		assert.Equal(t, `SELECT * FROM "user_base_models" WHERE "test_user_id" = $1`, result.Statement.SQL.String())
		return result
	})

}

func TestInjectFilterFromContext_NoUserInCtx_Slice(t *testing.T) {
	client, _ := GetMockPgClient()

	ctx := context.Background()
	client.DB = client.DB.WithContext(ctx)

	queryDb := client.DB.Model([]userBaseModel{})
	injectFilterFromContext(queryDb)

	_ = queryDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
		result := tx.Select("*").Find([]userBaseModel{})
		assert.Nil(t, result.Error)
		return result
	})

}

func TestApplyFilters_Success(t *testing.T) {
	client, _ := GetMockPgClient()

	userID := uuid.New()
	orgId := uuid.New()
	ctx := context.Background()
	ctx = apicontext.AddAuthToContext(ctx, "user", userID, []uuid.UUID{orgId})
	client.DB = client.DB.WithContext(ctx)

	applyFilters(client.DB, &userBaseModel{ID: userID})

	assert.Nil(t, client.DB.Error)
}

func TestValidateContext(t *testing.T) {
	t.Run("no user in context", func(t *testing.T) {
		ctx := context.Background()
		userId, orgIds, err := validateContext(ctx)

		assert.Equal(t, uuid.Nil, userId)
		assert.Nil(t, orgIds)
		assert.EqualError(t, err, "no user ID in context")
	})

	t.Run("with user in context", func(t *testing.T) {
		userID := uuid.New()
		orgId := uuid.New()
		ctx := context.Background()
		ctx = apicontext.AddAuthToContext(ctx, "user", userID, []uuid.UUID{orgId})

		userId, orgIds, err := validateContext(ctx)

		assert.Equal(t, userID, userId)
		assert.Equal(t, []uuid.UUID{orgId}, orgIds)
		assert.Nil(t, err)
	})
}

package pgclient

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/Zampfi/application-platform/services/api/pkg/errorreporting"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresClient struct {
	*gorm.DB
}

func NewPostgresClient(dsn string, maxIdleConnections int, maxOpenConnections int) (*PostgresClient, error) {

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	pgConnection := postgres.Open(dsn)

	db, err := gorm.Open(pgConnection, &gorm.Config{Logger: dbLogger})
	if err != nil {
		errorreporting.CaptureException(fmt.Errorf("error opening database connection: %w", err), context.Background())
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {

		return nil, err
	}

	sqlDB.SetMaxIdleConns(maxIdleConnections)
	sqlDB.SetMaxOpenConns(maxOpenConnections)

	client := &PostgresClient{db}

	registerAccessControlHooks(client.DB)

	return client, nil
}

// register access control hooks
func registerAccessControlHooks(db *gorm.DB) {
	beforeQuery(db)
	beforeUpdate(db)
	beforeDelete(db)
	beforePreload(db)
}

// ensure the access control filters are applied to all select queries
func beforeQuery(db *gorm.DB) {
	db.Callback().Query().Before("gorm:query").Register("inject_access_control_filters", func(d *gorm.DB) {
		injectFilterFromContext(d)
	})
}

// ensure the access control filters are applied to all update queries
func beforeUpdate(db *gorm.DB) {
	db.Callback().Update().Before("gorm:update").Register("inject_access_control_filters", func(d *gorm.DB) {
		injectFilterFromContext(d)
	})
}

// ensure the access control filters are applied to all delete queries
func beforeDelete(db *gorm.DB) {
	db.Callback().Delete().Before("gorm:delete").Register("inject_access_control_filters", func(d *gorm.DB) {
		injectFilterFromContext(d)
	})
}

// ensure the access control filters are applied to all preload queries
func beforePreload(db *gorm.DB) {
	db.Callback().Query().Before("gorm:preload").Register("inject_access_control_filters_preload", func(d *gorm.DB) {
		injectFilterFromContext(d)
	})
}

func coerceBaseModel(model interface{}) (BaseModel, error) {

	if model == nil {
		return nil, fmt.Errorf("query model is nil")
	}

	modelValue := reflect.ValueOf(model)
	modelType := modelValue.Type()

	// If it's a pointer, get the underlying value
	if modelType.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
		modelType = modelValue.Type()
	}

	// If it's a slice, get the element type
	if modelType.Kind() == reflect.Slice {
		elementType := modelType.Elem()

		// If element is a pointer, get its underlying type
		if elementType.Kind() == reflect.Ptr {
			elementType = elementType.Elem()
		}

		// Create a new instance of the element type
		newInstance := reflect.New(elementType).Interface()
		if baseModel, ok := newInstance.(BaseModel); ok {
			return baseModel, nil
		}
	} else {
		// For non-slice types
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}
		newInstance := reflect.New(modelType).Interface()
		if baseModel, ok := newInstance.(BaseModel); ok {
			return baseModel, nil
		}
	}

	return nil, fmt.Errorf("model must implement BaseModel or be a slice of BaseModels")
}

func injectFilterFromContext(db *gorm.DB) {
	model := db.Statement.Model

	baseModel, err := coerceBaseModel(model)
	if err != nil {
		errorreporting.CaptureException(fmt.Errorf("error coercing base model: %w", err), db.Statement.Context)
		db.AddError(err)
		return
	}

	if baseModel == nil {
		errorreporting.CaptureException(fmt.Errorf("error coercing base model: %w", err), db.Statement.Context)
		db.AddError(errors.New("model is nil"))
		return
	}

	applyFilters(db, baseModel)

}

func validateContext(ctx context.Context) (uuid.UUID, []uuid.UUID, error) {
	_, userId, orgIds := apicontext.GetAuthFromContext(ctx)

	if userId == nil {
		return uuid.Nil, nil, errors.New("no user ID in context")
	}

	return *userId, orgIds, nil
}

func applyFilters(db *gorm.DB, model BaseModel) {
	if db.Statement.Context == nil {
		return
	}

	userID, orgIDs, _ := validateContext(db.Statement.Context)

	filters := model.GetQueryFilters(db, userID, orgIDs)

	if filters.Statement.SQL.String() != "" {
		db.Where(filters.Statement.SQL.String(), filters.Statement.Vars...)
	}
}

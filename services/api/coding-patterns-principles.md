# API Service Coding Patterns and Principles

This document serves as a comprehensive guide for developers working on the Zampfi Application Platform API service. It outlines the coding patterns, architectural decisions, and best practices that have evolved through the development of the codebase. By following these patterns, you'll be able to contribute effectively while maintaining consistency and quality.

## Table of Contents
- [Overview](#overview)
- [Code Organization and Architecture](#code-organization-and-architecture)
  - [Directory Structure](#directory-structure)
  - [Core vs Pkg vs Helper](#core-vs-pkg-vs-helper)
  - [Domain-Driven Design](#domain-driven-design)
- [Service Pattern](#service-pattern)
  - [Interface-Based Design](#interface-based-design)
  - [Service Implementation](#service-implementation)
  - [Service Registration](#service-registration)
- [Models and Database](#models-and-database)
  - [Model Definition](#model-definition)
  - [GORM Hooks](#gorm-hooks)
  - [Query Filters](#query-filters)
  - [Store Interfaces](#store-interfaces)
  - [Store Implementation](#store-implementation)
- [Dependency Injection](#dependency-injection)
  - [Constructor Injection](#constructor-injection)
  - [Service Dependencies](#service-dependencies)
  - [Configuration Injection](#configuration-injection)
- [Error Handling](#error-handling)
  - [Domain-Specific Errors](#domain-specific-errors)
  - [Error Wrapping](#error-wrapping)
  - [Error Response Formatting](#error-response-formatting)
  - [Centralized Error Handling](#centralized-error-handling)
- [Middleware Pattern](#middleware-pattern)
  - [Authentication Middleware](#authentication-middleware)
  - [Logging Middleware](#logging-middleware)
  - [Middleware Chaining](#middleware-chaining)
- [Strategy Pattern](#strategy-pattern)
  - [Strategy Interface](#strategy-interface)
  - [Strategy Implementation](#strategy-implementation)
  - [Strategy Selection](#strategy-selection)
- [Testing Patterns](#testing-patterns)
  - [Table-Driven Tests](#table-driven-tests)
  - [Mocking Dependencies](#mocking-dependencies)
  - [Test Coverage Requirements](#test-coverage-requirements)
  - [Integration Tests](#integration-tests)
- [Request Handling](#request-handling)
  - [Handler Structure](#handler-structure)
  - [Request Validation](#request-validation)
  - [Response Formatting](#response-formatting)
- [Authentication and Authorization](#authentication-and-authorization)
  - [Auth Flow](#auth-flow)
  - [Permission Checking](#permission-checking)
  - [Role-Based Access Control](#role-based-access-control)
- [Transaction Handling](#transaction-handling)
  - [Transaction Patterns](#transaction-patterns)
  - [Rollback Handling](#rollback-handling)
- [Logging](#logging)
  - [Structured Logging](#structured-logging)
  - [Log Levels](#log-levels)
  - [Context-Based Logging](#context-based-logging)
- [Configuration Management](#configuration-management)
  - [Server Configuration](#server-configuration)
  - [Environment Variables](#environment-variables)
  - [Feature Flags](#feature-flags)
- [API Versioning](#api-versioning)
  - [URL-Based Versioning](#url-based-versioning)
  - [Version Compatibility](#version-compatibility)
- [Performance Considerations](#performance-considerations)
  - [Query Optimization](#query-optimization)
  - [Caching Strategies](#caching-strategies)
  - [Pagination](#pagination)
- [Common Pitfalls and Best Practices](#common-pitfalls-and-best-practices)
  - [Avoiding Circular Dependencies](#avoiding-circular-dependencies)
  - [Proper Error Handling](#proper-error-handling)
  - [Context Propagation](#context-propagation)

## Overview
This document outlines the coding patterns and principles used throughout the API service of the application platform. It serves as a comprehensive guide for new developers joining the team, ensuring they can contribute effectively while maintaining the codebase's quality and consistency. The patterns described here have evolved through practical experience and are designed to make the codebase maintainable, testable, and scalable.

## Code Organization and Architecture

### Directory Structure

The API service follows a clean architecture with distinct layers:

```
services/api/
├── core/           # Business logic organized by domain
│   ├── datasets/
│   ├── dataplatform/
│   ├── organizations/
│   ├── widgets/
│   └── auth/
├── server/         # HTTP server, routes, and middleware
│   ├── middleware/
│   ├── routes/
│   └── server.go
├── db/             # Database models and store implementations
│   ├── models/
│   ├── store/
│   └── pgclient/
├── pkg/            # Shared utilities and libraries
├── helper/         # Helper functions and utilities
├── config/         # Configuration management
└── main.go         # Application entry point
```

### Core vs Pkg vs Helper

Understanding the distinction between these directories is crucial for proper code organization:

1. **Core**: Contains domain-specific business logic organized by feature area. Each domain has its own directory with services, models, and errors.

   Example: `core/datasets/service/service.go` implements the dataset service with methods like `CreateDataset`, `GetDatasetById`, etc.

   ```go
   // From core/datasets/service/service.go
   func (s *datasetService) GetDatasetById(ctx context.Context, datasetId uuid.UUID) (*models.Dataset, error) {
       dataset, err := s.store.GetDatasetById(ctx, datasetId)
       if err != nil {
           return nil, errors.NewDatasetNotFoundError(datasetId)
       }
       return dataset, nil
   }
   ```

2. **Pkg**: Contains reusable packages that are not tied to specific domains. These packages can be used across different parts of the application or even in other services.

   Example: `pkg/apierrors/error.go` defines a standard error type used throughout the application:

   ```go
   // From pkg/apierrors/error.go
   type APIError struct {
       code    int
       message string
       detail  string
   }

   func (e *APIError) Error() string {
       return e.message
   }

   func (e *APIError) StatusCode() int {
       return e.code
   }
   ```

3. **Helper**: Contains utility functions that assist with common tasks but don't fit into a specific domain or package.

   Example: `helper/context/auth.go` provides functions for working with authentication in the context:

   ```go
   // From helper/context/auth.go
   func AddAuthToContext(ctx context.Context, role string, userID uuid.UUID, userOrganizations []uuid.UUID) context.Context {
       enrichedCtx := AddCtxVariableToCtx(ctx, contextKeyUserID, userID)
       enrichedCtx = AddCtxVariableToCtx(enrichedCtx, contextKeyUserOrganizations, userOrganizations)
       enrichedCtx = AddCtxVariableToCtx(enrichedCtx, contextKeyUserRole, role)
       return enrichedCtx
   }
   ```

### Domain-Driven Design

Each domain (datasets, organizations, widgets, etc.) has its own directory in `core/` with a consistent structure:

```
core/datasets/
├── service/       # Service interface and implementation
├── models/        # Domain-specific models
├── errors/        # Domain-specific errors
└── util/          # Domain-specific utilities
```

This structure makes it easy to understand the domain boundaries and responsibilities. When adding a new feature, you should identify which domain it belongs to and place it in the appropriate directory.

Example of domain-specific error definitions:

```go
// From core/datasets/errors/errors.go
func NewDatasetNotFoundError(datasetId uuid.UUID) *apierrors.APIError {
    return apierrors.NewNotFoundError(
        fmt.Sprintf("dataset with ID %s not found", datasetId),
        "The requested dataset could not be found",
    )
}
```

## Service Pattern

### Interface-Based Design

Services are defined as interfaces, promoting loose coupling and testability:

```go
// From core/datasets/service/service.go
type DatasetService interface {
    GetDatasetById(ctx context.Context, datasetId uuid.UUID) (*models.Dataset, error)
    CreateDataset(ctx context.Context, dataset models.CreateDatasetRequest) (*models.Dataset, error)
    UpdateDataset(ctx context.Context, datasetId uuid.UUID, dataset models.UpdateDatasetRequest) (*models.Dataset, error)
    DeleteDataset(ctx context.Context, datasetId uuid.UUID) error
    // ... other methods
}
```

This approach allows for:
- Easy mocking in tests
- Multiple implementations (e.g., for different data sources)
- Clear contract definition

### Service Implementation

Implementations use private structs with dependencies injected via constructors:

```go
// From core/datasets/service/service.go
type datasetService struct {
    store store.DatasetStore
    dataPlatformService dataplatform.DataPlatformService
    // Other dependencies
}

func NewDatasetService(
    store store.DatasetStore,
    dataPlatformService dataplatform.DataPlatformService,
    // Other dependencies
) (DatasetService, error) {
    return &datasetService{
        store: store,
        dataPlatformService: dataPlatformService,
        // Initialize other dependencies
    }, nil
}
```

Services implement the interface methods:

```go
// From core/datasets/service/service.go
func (s *datasetService) GetDatasetById(ctx context.Context, datasetId uuid.UUID) (*models.Dataset, error) {
    // Implementation with error handling
    dataset, err := s.store.GetDatasetById(ctx, datasetId)
    if err != nil {
        return nil, errors.NewDatasetNotFoundError(datasetId)
    }
    return dataset, nil
}
```

### Service Registration

Services are registered in `main.go` and accessed via the `ServerConfig`:

```go
// From config/config.go
type ServerConfig struct {
    AuthClient         *kratosclient.Client
    Env                *ConfigVariables
    DataPlatformConfig *DataPlatformConfig
    Store              store.Store
    TemporalSdk        temporalsdk.TemporalService
    // Other configurations
}

// In main.go
dataPlatformService, err := dataplatform.NewDataPlatformService(serverCfg)
if err != nil {
    // Handle error
}

datasetService, err := datasets.NewDatasetService(
    serverCfg.Store,
    dataPlatformService,
    // Other dependencies
)
if err != nil {
    // Handle error
}
```

## Models and Database

### Model Definition

Models are defined in `db/models/` and use GORM tags for mapping to database tables:

```go
// From db/models/datasets.go
type Dataset struct {
    ID             uuid.UUID       `json:"dataset_id" gorm:"column:dataset_id"`
    Title          string          `json:"title"`
    Description    *string         `json:"description"`
    OrganizationId uuid.UUID       `json:"organization_id" gorm:"column:organization_id"`
    Type           DatasetType     `json:"type"`
    CreatedBy      uuid.UUID       `json:"created_by" gorm:"column:created_by"`
    Metadata       json.RawMessage `json:"metadata"`
    CreatedAt      time.Time       `json:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at"`
    DeletedAt      *time.Time      `json:"deleted_at,omitempty"`
}
```

Always implement the `TableName()` method to explicitly set the database table:

```go
// From db/models/datasets.go
func (d *Dataset) TableName() string {
    return "datasets"
}
```

### GORM Hooks

Models implement GORM hooks for validation and permission checks:

```go
// From db/models/datasets.go
func (a *Dataset) BeforeCreate(db *gorm.DB) error {
    // Get user ID from context for permission check
    _, userId, _ := apicontext.GetAuthFromContext(db.Statement.Context)
    if userId == nil {
        return fmt.Errorf("no user id found in context")
    }

    // Check organization access
    fraps := []FlattenedResourceAudiencePolicy{}
    err := db.Where("resource_type = ? AND resource_id = ? AND user_id = ? AND deleted_at IS NULL", 
                   "organization", a.OrganizationId, userId).Limit(1).Find(&fraps).Error

    if err != nil {
        return err
    }

    if len(fraps) == 0 {
        return fmt.Errorf("organization access forbidden")
    }

    return nil
}
```

Common hooks include:
- `BeforeCreate`: Validation and permission checks before record creation
- `BeforeUpdate`: Validation and permission checks before record update
- `BeforeDelete`: Permission checks and cleanup before record deletion
- `AfterFind`: Post-processing after retrieving records

When implementing a new model, you should consider which hooks are necessary for validation and access control. For example, the `Widget` model prevents insertions with a simple hook:

```go
// From db/models/widgets.go
func (w *Widget) BeforeCreate(db *gorm.DB) error {
    return fmt.Errorf("insert forbidden")
}
```

### Query Filters

Models can implement query filter methods for access control:

```go
// From db/models/datasets.go
func (d *Dataset) GetQueryFilters(db *gorm.DB, userId uuid.UUID, orgIds []uuid.UUID) *gorm.DB {
    return db.Where(
        `EXISTS (
            SELECT 1 FROM "app"."flattened_resource_audience_policies" frap
            WHERE frap.resource_type = 'dataset'
            AND frap.resource_id = datasets.dataset_id
            AND frap.user_id = ?
            AND frap.deleted_at IS NULL
        )`, userId,
    )
}
```

These filters are used in store methods to enforce access control:

```go
// From db/store/datasets.go
func (s *appStore) GetDatasetById(ctx context.Context, datasetId uuid.UUID) (*models.Dataset, error) {
    var dataset models.Dataset
    err := s.client.WithContext(ctx).
        Scopes(dataset.GetQueryFilters).  // Apply query filters
        Where("dataset_id = ?", datasetId).
        First(&dataset).Error
    if err != nil {
        return nil, err
    }
    return &dataset, nil
}
```

### Store Interfaces

Store interfaces are defined in `db/store/store.go` and provide data access methods:

```go
// From db/store/datasets.go
type DatasetStore interface {
    GetDatasetById(ctx context.Context, datasetId uuid.UUID) (*models.Dataset, error)
    CreateDataset(ctx context.Context, dataset models.Dataset) (*models.Dataset, error)
    UpdateDataset(ctx context.Context, dataset models.Dataset) (*models.Dataset, error)
    DeleteDataset(ctx context.Context, datasetId uuid.UUID) error
    // Other methods
}
```

The main `Store` interface aggregates all domain-specific stores:

```go
// From db/store/store.go
type Store interface {
    DatasetStore
    TeamStore
    OrganizationStore
    // Other store interfaces
    
    WithTx(ctx context.Context, fn func(Store) error) error
}
```

### Store Implementation

Store methods implement the interfaces and handle database operations:

```go
// From db/store/datasets.go
func (s *appStore) GetDatasetById(ctx context.Context, datasetId uuid.UUID) (*models.Dataset, error) {
    var dataset models.Dataset
    err := s.client.WithContext(ctx).
        Scopes(dataset.GetQueryFilters).
        Where("dataset_id = ?", datasetId).
        First(&dataset).Error
    if err != nil {
        return nil, err
    }
    return &dataset, nil
}
```

When implementing store methods, follow these patterns:
1. Always use `WithContext(ctx)` to propagate context
2. Use model query filters for access control
3. Return raw errors from database operations
4. Use transactions for multi-step operations

Example of a store method with pagination:

```go
// From db/store/datasets.go
func (s *appStore) GetDatasetsByOrganizationId(ctx context.Context, organizationId uuid.UUID, page int, pageSize int, searchTerm string) ([]models.Dataset, int64, error) {
    var datasets []models.Dataset
    var count int64

    query := s.client.WithContext(ctx).
        Scopes(models.Dataset{}.GetQueryFilters).
        Where("organization_id = ?", organizationId)

    if searchTerm != "" {
        query = query.Where("title ILIKE ?", fmt.Sprintf("%%%s%%", searchTerm))
    }

    err := query.Count(&count).Error
    if err != nil {
        return nil, 0, err
    }

    if page > 0 && pageSize > 0 {
        offset := (page - 1) * pageSize
        query = query.Offset(offset).Limit(pageSize)
    }

    err = query.Order("created_at DESC").Find(&datasets).Error
    if err != nil {
        return nil, 0, err
    }

    return datasets, count, nil
}
```

## Dependency Injection

### Constructor Injection

Dependencies are provided through constructor parameters, making services testable and modular:

```go
// From core/widgets/service/service.go
func NewWidgetsService(
    store store.Store,
    datasetService datasetservice.DatasetService,
    dataPlatformService dataplatform.DataPlatformService,
) (*widgetsService, error) {
    return &widgetsService{
        store:               store,
        datasetService:      datasetService,
        dataPlatformService: dataPlatformService,
    }, nil
}
```

This approach:
- Makes dependencies explicit
- Facilitates testing with mocks
- Prevents hidden dependencies
- Allows for dependency lifecycle management

### Service Dependencies

Services often depend on other services. These dependencies should be clearly defined in the constructor:

```go
// From core/dataplatform/service.go
func NewDataPlatformService(
    serverCfg *serverconfig.ServerConfig,
) (DataPlatformService, error) {
    // Initialize dependencies
    dataService, err := data.NewDataService(serverCfg)
    if err != nil {
        return nil, err
    }
    
    actionService, err := actions.NewActionService(serverCfg)
    if err != nil {
        return nil, err
    }
    
    return &dataPlatformService{
        dataService:    dataService,
        actionService:  actionService,
        store:          serverCfg.Store,
        temporalClient: serverCfg.TemporalSdk,
    }, nil
}
```

### Configuration Injection

The `ServerConfig` struct in `config/config.go` is used to centralize configuration and dependencies:

```go
// From config/config.go
type ServerConfig struct {
    AuthClient         *kratosclient.Client
    Env                *ConfigVariables
    DataPlatformConfig *DataPlatformConfig
    Store              store.Store
    TemporalSdk        temporalsdk.TemporalService
    // Other configurations
}
```

This configuration is initialized in `main.go` and passed to services:

```go
// From main.go
serverCfg := &serverconfig.ServerConfig{
    AuthClient:         authClient,
    Env:                env,
    DataPlatformConfig: dataPlatformConfig,
    Store:              store,
    TemporalSdk:        temporalSdk,
}

// Pass to services
datasetService, err := datasets.NewDatasetService(serverCfg)
if err != nil {
    log.Fatal(err)
}
```

## Error Handling

### Domain-Specific Errors

Domain-specific errors are defined in `core/*/errors/` packages:

```go
// From core/datasets/errors/errors.go
func NewDatasetNotFoundError(datasetId uuid.UUID) *apierrors.APIError {
    return apierrors.NewNotFoundError(
        fmt.Sprintf("dataset with ID %s not found", datasetId),
        "The requested dataset could not be found",
    )
}
```

The `APIError` struct in `pkg/apierrors/` provides a standard error format:

```go
// From pkg/apierrors/error.go
type APIError struct {
    code    int
    message string
    detail  string
}

func (e *APIError) Error() string {
    return e.message
}

func (e *APIError) StatusCode() int {
    return e.code
}

func (e *APIError) Detail() string {
    return e.detail
}
```

### Error Wrapping

Errors are wrapped to preserve context:

```go
// From core/datasets/service/service.go
func (s *datasetService) GetDatasetById(ctx context.Context, datasetId uuid.UUID) (*models.Dataset, error) {
    dataset, err := s.store.GetDatasetById(ctx, datasetId)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.NewDatasetNotFoundError(datasetId)
        }
        return nil, fmt.Errorf("failed to get dataset: %w", err)
    }
    return dataset, nil
}
```

### Error Response Formatting

HTTP handlers format errors into appropriate responses:

```go
// From server/routes/datasets/handler.go
func (h *DatasetHandler) GetDatasetById(c *gin.Context) {
    datasetId := c.Param("datasetId")
    
    parsedId, err := uuid.Parse(datasetId)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dataset ID"})
        return
    }
    
    dataset, err := h.datasetService.GetDatasetById(c, parsedId)
    if err != nil {
        apiErr, ok := err.(*apierrors.APIError)
        if ok {
            c.JSON(apiErr.StatusCode(), gin.H{"error": apiErr.Error(), "details": apiErr.Detail()})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        }
        return
    }
    
    c.JSON(http.StatusOK, dataset)
}
```

### Centralized Error Handling

The API service uses middleware for centralized error handling:

```go
// From server/middleware/panic_recovery.go
func PanicRecoveryMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                logger := apicontext.GetLoggerFromCtx(c)
                logger.Error("panic recovered", zap.Any("error", err), zap.String("stack", string(debug.Stack())))
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
            }
        }()
        c.Next()
    }
}
```

## Middleware Pattern

### Authentication Middleware

Authentication middleware validates user credentials and adds authentication information to the context:

```go
// From server/middleware/auth.go
func GetAuthMiddleware(serverCfg *serverconfig.ServerConfig) (gin.HandlerFunc, error) {
    authSvc, err := authservice.NewAuthService(
        serverCfg.Env.ControlPlaneAdminSecrets, 
        serverCfg.Env.AuthBaseUrl, 
        serverCfg.Store, 
        serverCfg.Env.Environment,
    )
    if err != nil {
        return nil, err
    }

    return func(c *gin.Context) {
        role, userId, organizationIds := authSvc.ResolveAdminInfo(c, c.Request.Header)

        if role == "admin" {
            apictx.AddAuthToGinContext(c, role, userId, organizationIds)
        } else {
            userMiddleWare(authSvc)(c)
        }

        c.Next()
    }, nil
}
```

### Logging Middleware

Logging middleware adds request-specific information to the logger:

```go
// From server/middleware/logging.go
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        logger := zap.L().With(
            zap.String("request_id", uuid.New().String()),
            zap.String("path", c.Request.URL.Path),
            zap.String("method", c.Request.Method),
        )
        
        ctx := apicontext.AddLoggerToContext(c.Request.Context(), logger)
        c.Request = c.Request.WithContext(ctx)
        
        start := time.Now()
        c.Next()
        duration := time.Since(start)
        
        logger.Info("request completed",
            zap.Int("status", c.Writer.Status()),
            zap.Duration("duration", duration),
        )
    }
}
```

### Middleware Chaining

Middleware is chained in the router setup:

```go
// From server/server.go
func SetupRouter(serverCfg *serverconfig.ServerConfig) (*gin.Engine, error) {
    router := gin.New()
    
    // Add middleware
    router.Use(gin.Recovery())
    router.Use(middleware.PanicRecoveryMiddleware())
    router.Use(middleware.LoggingMiddleware())
    
    corsMiddleware := middleware.GetCORSMiddleware(serverCfg.Env.AllowedOrigins)
    router.Use(corsMiddleware)
    
    authMiddleware, err := middleware.GetAuthMiddleware(serverCfg)
    if err != nil {
        return nil, err
    }
    
    // Set up routes with middleware
    api := router.Group("/api")
    api.Use(authMiddleware)
    
    // Add route groups
    datasets.AddRoutes(api, serverCfg)
    organizations.AddRoutes(api, serverCfg)
    
    return router, nil
}
```

## Common Pitfalls and Best Practices

### Avoiding Circular Dependencies

Circular dependencies can cause compilation errors and make the codebase harder to understand. To avoid them:

1. Use interfaces to break dependency cycles
2. Move shared code to a common package
3. Restructure code to follow a clear dependency direction

Example of breaking a dependency cycle with an interface:

```go
// Before: Circular dependency between datasets and dataplatform
// datasets -> dataplatform -> datasets

// After: Break the cycle with an interface
// In core/datasets/service/service.go
type DataPlatformService interface {
    RegisterDataset(ctx context.Context, datasetId uuid.UUID) error
}

// In core/dataplatform/service.go
func (s *dataPlatformService) RegisterDataset(ctx context.Context, datasetId uuid.UUID) error {
    // Implementation
}
```

### Proper Error Handling

Always handle errors appropriately:

1. Return domain-specific errors when possible
2. Wrap errors to preserve context
3. Log errors with appropriate level and context
4. Don't expose internal errors to clients

### Context Propagation

Always propagate context through the call stack:

1. Pass context as the first parameter to functions
2. Use context to carry request-specific data (auth, logging, etc.)
3. Respect context cancellation and deadlines

```go
// Good context propagation
func (s *datasetService) GetDatasetById(ctx context.Context, datasetId uuid.UUID) (*models.Dataset, error) {
    // Pass context to store
    dataset, err := s.store.GetDatasetById(ctx, datasetId)
    if err != nil {
        return nil, err
    }
    
    // Pass context to other service
    err = s.dataPlatformService.ValidateDataset(ctx, dataset.ID)
    if err != nil {
        return nil, err
    }
    
    return dataset, nil
}
```

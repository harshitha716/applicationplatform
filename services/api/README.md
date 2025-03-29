# API service

API service is the backend for the application platform.

## Localdev

- Setup dependencies of this service by running `make setup-dependencies` from the `api` directory (one time setup)
- Launch the server from the `Run and debug` panel of vscode/cursor by running the `Launch Go API Server` process
- Optinally, if you want to run the server from terminal through go run:

    ```
    export AUTH_HOST=http://localhost:4433
    export PG_DATABASE_URL_APP="postgres://postgres:postgrespassword@localhost:5432/postgres?sslmode=disable&search_path=app"
    go run main.go`
    ```

## Code structure

This document explains the directory structure of the API service codebase, with a focus on clarity, maintainability, and scalability. The goal is to ensure that different parts of the codebase are organized logically to support easy navigation, independent development, and long-term growth.

### Root Directory
- **Dockerfile**: Instructions for building the Docker image of the service.
- **go.mod** and **go.sum**: Go module files that manage dependencies.

### Core Modules (`core/`)
The `core` directory contains the primary business logic of the application. It is responsible for encapsulating the key functional areas of the service, each focused on a specific part of the system’s operations. By centralizing this functionality, we can ensure that the core components of the service are self-contained and reusable across different parts and entrypoints (server, worker etc) of the application. This organization also makes it easier to extend or update a specific part of the service without affecting others.

### Database Layer (`db/`)
The `db` directory houses all database-related logic, including models, database clients, and CRUD operations. It centralizes data persistence functionality, ensuring that database interactions are isolated from the rest of the application’s logic. This separation makes it easier to manage, test, and optimize database-related tasks independently of business logic.

- **models**: Contains the database models and schemas that represent entities like users, organizations, and datasets.
- **pgclient**: Provides the database client and related utility functions for interacting with PostgreSQL.
- **store**: Contains the implementation of CRUD operations for the various database models, abstracting the database logic. The store is designed so that services can define what models they want/need and only have access to those.

### Mocks (`mocks/`)
The `mocks` directory holds mock implementations of interfaces, primarily for use in testing. These mocks simulate interactions with external systems or services, enabling more controlled and isolated unit tests. This organization helps keep test dependencies and external integrations separate from actual business logic, allowing tests to be faster and more reliable.

### Package Libraries (`pkg/`)
The `pkg` directory contains reusable libraries, clients, and SDKs that provide functionality shared across multiple parts of the service. It serves as a central location for commonly used code, such as third-party integrations (e.g., Kratos for authentication, Databricks for data processing), error handling utilities, and logging services. The `pkg` directory promotes code reuse, ensuring that shared functionality is maintained in one place and can be easily extended or updated. This organization also makes it easier to eventually share these libraries across different services if needed.

### Server (`server/`)
The `server` directory is the api server entrypoint and is responsible for managing the API's routing and the configuration of the HTTP server. This separation allows the server setup to be independent of the business logic, making it easier to modify the server’s behavior (e.g., adding middleware or new routes) without touching the core application functionality. By isolating server-related code, it ensures that the routing logic is clear and easy to maintain. The routes in the `server` directory are organized to have exactly the same structure as the actual routes.

### FAQ

- What's the difference between `pkg` and `core`?
    - `pkg` defines libraries, clients and SDKs -- essentially independent reusable packages that can be reused across the core modules. `core` contains the independent modules of business logic in the application that can be reused by multiple entrypoints of the application.
- What can be other entrypoints of the application?
    - `server` is an existing entrypoint. More examples of entrypoints that we might need in the future:
        - Temporal Workers: They will use core modules in their own desired way
        - CLI binaries or WASM bundles: Eventually if non-golang services want to use some code from the core modules, we can wrap them in a CLI binary and embed them in other runtimes
        - GRPC endpoints: If some services need to use business logic through GRPC, that would be another entrypoint

- Is this designed with microservices in mind?
    - No
- Why are the server routes separate from core business logic modules?
    - A server route or any other entrypoint can compose logic from different modules, so it technically does not belong in the same module. Routes and modules must not be coupled because logic can change regardless or the route and routes can change regardless of the logic.
- Why is `db` store outside?
    - It is a centralized abstraction over the database which will have all the models and the right way to query them with the right access control. Just like we have a centralized ORM, this is another layer above the ORM that makes the codebase abstracted from the database details and dialect.
 - How do I seggregate specific models for a core business logic module and ensure they don't access other modules?
    - Each model has a writer store and a reader store. You just have to define an interface that specifies what store it needs:

    ```go
    type OrganizationServiceStore interface {
        store.OrganizationStore
        GetUsersAll(ctx context.Context) ([]models.User, error)
    }

    type OrganizationService struct {
        store OrganizationServiceStore
    }

    func NewOrganizationService(appStore store.Store) *OrganizationService {
        return &OrganizationService{store: appStore}
    }

    func (s *OrganizationService) GetOrganizations(ctx context.Context) ([]models.Organization, error) {
        ctxLogger := apicontext.GetLoggerFromCtx(ctx)

        orgs, err := s.store.GetOrganizationsAll(ctx)
        if err != nil {
            ctxLogger.Error("failed to get organizations", zap.Error(err))
            return nil, err
        }
        return orgs, nil
    }
    ```

    In the above code, although `appStore` has methods like `GetDatasets`, `GetWidgets`, the organization service is only restricted to `GetOrganizationById`, `GetOrganizationsAll` and `GetUsersAll` methods. This way we can ensure seggregation of concerns in store.

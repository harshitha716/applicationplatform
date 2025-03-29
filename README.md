# ApplicationPlatform Services

This repository contains a collection of services that make up the application platform.

## Table of Contents
1. [Services Overview](#services-overview)  
   1.1 [Auth Service](#auth-service)  
   1.2 [API Service](#api-service)  
   1.3 [Dashboard Service](#dashboard-service)  
   1.4 [Database](#database)  
   1.5 [Migration Services](#migration-services)  
   1.6 [Development Tools](#development-tools)  
2. [Repository Structure](#repository-structure)
3. [Getting Started](#getting-started)  
   3.1 [Local Development](#local-development)  
   3.2 [Cloud Development](#cloud-development)  
4. [Development Workflow](#development-workflow)  
   4.1 [Code Verification](#code-verification)  
   4.2 [Pull Request Guidelines](#pull-request-guidelines)  
5. [Testing](#testing)  
6. [API Testing](#api-testing)  
   6.1 [Logging in for Manual API Testing Using Postman](#logging-in-for-manual-api-testing-using-postman)  
   6.2 [Emulating a User Request in Admin Mode](#emulating-a-user-request-in-admin-mode)  
7. [Observability](#observability)


---

## Services Overview

### Auth Service
- **Authentication service** powered by Ory Kratos
- Handles user identity, authentication flows, and SSO integration
- Exposed on ports:
  - `4433` (public)
  - `4434` (admin)
- Configurable via environment variables for:
  - URLs
  - SMTP
  - OIDC providers
- Located in `services/auth/`

### API Service
- **Backend API service** written in Go
- Runs on port `8080`
- Integrates with the auth service and database
- Configurable via environment variables for:
  - Database connection
  - CORS
- Located in `services/api/`
- Core modules include:
  - Authentication
  - Organizations
  - Datasets
  - Rules
  - Sheets
  - Widgets

### Dashboard Service
- **Frontend service** built with Next.js and TypeScript
- User interface for interacting with the application
- Located in `services/dashboard/`
- Key features:
  - Authentication UI
  - Data visualization
  - Configuration interfaces

### Database
- **PostgreSQL database**
- Used by both Auth and API services
- Separate schemas for:
  - Kratos
  - Application data
- Migrations handled by dedicated services

### Migration Services
- **Platform migrations service** for database setup
  - Located in `services/platform-migrations/`
- **App migrations service** for application schema changes
  - Located in `services/app-migrations/`
- Automatically run on startup

### Development Tools
- **Mailhog** (port `8025`) - Email testing interface
- **Adminer** (port `4040`) - Database management UI
- **Seeds service** for loading initial data
- **Temporal UI** (port `8888`) - Workflow management interface

---

## Repository Structure

```
application-platform/
├── .env.example             # Example environment variables
├── .env.cloud.example       # Example environment variables for cloud dev
├── Makefile                 # Development commands
├── docker-compose.yaml      # Local development services
├── docker-compose.cloud.yaml # Cloud development configuration
├── services/                # All application services
│   ├── api/                 # Backend API service (Go)
│   │   ├── core/            # Core business logic
│   │   ├── server/          # HTTP server implementation
│   │   ├── db/              # Database access layer
│   │   ├── pkg/             # Shared packages
│   │   └── workers/         # Background workers
│   ├── dashboard/           # Frontend service (Next.js)
│   │   ├── src/             # Source code
│   │   │   ├── apis/        # API client code
│   │   │   ├── components/  # Reusable UI components
│   │   │   ├── modules/     # Feature modules
│   │   │   ├── pages/       # Next.js pages
│   │   │   └── utils/       # Utility functions
│   ├── auth/                # Authentication service (Ory Kratos)
│   ├── app-migrations/      # Application database migrations
│   ├── platform-migrations/ # Platform database migrations
│   └── pinot-proxy/         # Proxy for Pinot database
└── localdev/                # Local development utilities
    ├── seeds/               # Seed data for development
    └── temporal/            # Temporal workflow configuration
```

---

## Getting Started

### Local Development

1. Set up environment variables:
   ```bash
   cp .env.example .env
   ```

2. Update the `DATAPLATFORM_CONFIG` environment variable with the config values from GCP Secret Manager.

3. Start all services:
   ```bash
   make localdev
   # OR docker compose up -d --build
   ```

4. Rebuild any service:
   ```bash
   docker compose up -d --build <service-name>
   ```
   Example:
   ```bash
   docker compose up -d --build api
   ```

5. Stop all services:
   ```bash
   make down
   # OR docker compose down
   ```

6. Clean the environment completely (removes volumes):
   ```bash
   make clean
   # OR docker compose --profile "*" down -v
   ```

7. Default login credentials: 
   - Email: `admin@zamp.ai` 
   - Password: `MeZamp@123`

### Cloud Development

You can also start the dev environment with services pointing to the database of the dev environment in GCP:

1. Copy the cloud environment example file:
   ```bash
   cp .env.cloud.example .env.cloud
   ```

2. Update the values in `.env.cloud` with the PostgreSQL database credentials of the dev environment.

3. Start the cloud development environment:
   ```bash
   make clouddev
   ```

4. Clean the environment:
   ```bash
   make clean
   ```

---

## Development Workflow

### Code Verification

When verifying your code:

- Run lint for the dashboard service:
  ```bash
  cd services/dashboard
  npm run lint
  ```

- Auto-fix linting issues:
  ```bash
  cd services/dashboard
  npm run lint -- --fix
  ```

- For Go code in the API service, use Go's built-in tools:
  ```bash
  cd services/api
  go fmt ./...
  go vet ./...
  ```

### Pull Request Guidelines

- PRs should target the `main` branch
- Changes are automatically deployed to staging on merge to main
- Production deployments require manual approval

#### CI Checks

The following checks run automatically on PRs:
- API service tests (for Go code changes)
- Dashboard tests (for dashboard service changes)
- Migration clean run (for SQL changes)
- Docker build verification
- Security scans (Semgrep for general security scanning, Gosec for Go-specific security scanning)

---

## Testing

- Ensure every line of logic in your codebase is covered with unit tests.
- Generate mocks using [Mockery](https://vektra.github.io/mockery/latest/):
  1. Install Mockery:
     ```bash
     brew install mockery
     ```
  2. Generate mocks:
     ```bash
     mockery
     ```
     Example: Run the `mockery` command in the `services/api` folder.

---

## API Testing

### Logging in for Manual API Testing Using Postman

To manually test APIs in Postman, follow these steps to authenticate and create a session:

#### Step 1: Retrieve Login Flow Information
Make a `GET` request to:
```
http://localhost:8080/auth/relay/self-service/login/browser
```

**Headers**:  
- `Accept: application/json`

**Sample Response**:
```json
{
  "id": "7d3786e9-5a2c-4fce-80b2-cd5ec05bdba0",
  "ui": {
    "action": "http://localhost:8080/auth/relay/self-service/login?flow=7d3786e9-5a2c-4fce-80b2-cd5ec05bdba0",
    "nodes": [
      {
        "attributes": {
          "name": "csrf_token",
          "value": "csrf_token_value"
        }
      }
    ]
  }
}
```

#### Step 2: Submit Login Credentials
Extract the `ui.action` URL and make a `POST` request to it.  
**Headers**:  
- `Accept: application/json`  
- `Content-Type: application/json`

**Body**:
```json
{
  "csrf_token": "csrf_token_value",
  "method": "password",
  "identifier": "admin@zamp.ai",
  "password": "Zamp@123"
}
```

#### Step 3: Verify Successful Login
On a successful request, a session is created. Use Postman to test other APIs as an authenticated user.

---

### Emulating a User Request in Admin Mode

To make queries as a user, pass these headers:  

```plaintext
X-Zamp-Admin-Secret: secret1
X-Zamp-User-Id: <user_id>
X-Zamp-Organization-Ids: <comma-separated-org-ids>
```

---

## Observability

View logs for all services on Google Cloud Console. Use the following query:  

```plaintext
resource.labels.namespace_name="hcp"
resource.labels.container_name="zamp-hcp-<svc-name>-<svc-name>"
```

Example for the API service:  

```plaintext
resource.labels.namespace_name="hcp"
resource.labels.container_name="zamp-hcp-api-api"
```

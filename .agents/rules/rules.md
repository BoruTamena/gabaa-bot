---
trigger: always_on
---

You are an expert Go backend engineer working with Gin and a clean modular architecture using GORM with PostgreSQL.

Agent Workflow Principles:
- Always follow clean architecture (handler → module → storage → platform)
- Keep code simple, scalable, and production-ready
- Never mix responsibilities across layers
- Think before coding, design first

Understanding the Task:
- Identify the feature or bug clearly
- Determine inputs, outputs, and side effects
- Map the domain (e.g., user, payment, ticket)

Architecture Mapping:
- handler → HTTP request/response
- module → business logic and use cases
- storage → database and cache
- platform → external services
- glue → routing and dependency injection

Implementation Flow:
- Define DTOs (request/response)
- Define interfaces (repository/platform)
- Implement storage layer (GORM)
- Implement module (business logic)
- Implement handler (HTTP layer)
- Wire routes in glue layer

Handler Rules:
- Parse and validate request
- Call module layer
- Return consistent JSON response
- Never write business logic
- Never access database directly
- MUST include detailed Swagger documentation for every endpoint

Swagger Documentation Rules (MANDATORY):
- Use Swagger annotations for all handlers
- Include:
  - Summary
  - Description
  - Tags (group by module/feature)
  - Accept / Produce (json)
  - Params (body, query, path)
  - Success response
  - Error responses
  - Route definition

Example:
```go
// CreateUser godoc
// @Summary Create user
// @Description Create a new user account
// @Tags User
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User data"
// @Success 201 {object} Response{data=UserResponse}
// @Failure 400 {object} Response{error=string}
// @Failure 500 {object} Response{error=string}
// @Router /api/v1/users [post]

Module Rules:
- Handle all business logic
- Manage transactions
- Orchestrate storage and platform
- Use interfaces for dependencies

Storage Rules (GORM + PostgreSQL):
- Use GORM for all DB operations
- Define models in `internal/constant/model/`
- Use repository pattern in `storage/persistence/`
- No business logic
- Always pass `context.Context`
- Use transactions when required
- Handle errors properly (e.g., record not found)

Platform Rules:
- Wrap external APIs/services
- Handle timeouts and retries
- Keep isolated and replaceable

Routing:
- Group routes by feature and version
- Follow RESTful conventions
- Handle 404 and 405 errors

Middleware:
- Logging
- Recovery
- Authentication (JWT/API Key)
- CORS
- Rate limiting

Request Handling:
- Bind JSON, query, and params
- Validate and sanitize input
- Handle file uploads if needed

Response Handling:
- Return consistent JSON structure:
  {
    "success": true,
    "data": {},
    "error": null
  }
- Use proper HTTP status codes
- Never expose internal errors

Database Integration:
- Use PostgreSQL with GORM
- Manage connection pool properly
- Use AutoMigrate carefully (dev only preferred)
- Handle DB errors gracefully
- Use transactions in module layer

Error Handling:
- Use centralized error definitions
- Return safe, structured errors
- Avoid leaking internal details 

Testing:
- Use mocks for repositories
- Write unit tests for modules
- Write integration tests for handlers

Security:
- Validate all inputs
- Sanitize user data
- Protect sensitive operations
- Use secure headers and HTTPS

Best Practices:
- Use dependency injection
- Keep modules small and focused
- Avoid tight coupling
- Log important events and errors
- Monitor performance

Anti-Patterns:
- Fat handlers
- Business logic in middleware
- Direct DB access from handler
- Global mutable state
- Tight coupling between modules

Final Rule:
- If logic grows → move to module
- If module grows → split by domain
- If structure breaks → refactor before continuing
# Todo Project – Architecture Overview

## 1. Monorepo Structure

```
.
├── docker-compose.yml       # Define and run multi-service containers
├── gateways
│   └── todo-bff             # Backend-for-Frontend (API gateway for clients)
├── go.work                  # Workspace to link multiple Go modules
├── go.work.sum              # Checksum for workspace dependencies
├── pkg
│   ├── auth                 # Shared authentication utilities
│   ├── config               # Shared configuration utilities
│   ├── go.mod               # Module definition for shared packages
│   └── logger               # Shared logging utilities
├── proto
│   ├── go.mod               # Module for protobuf contracts
│   ├── todo                 # Proto definitions for todo domain
│   └── user                 # Proto definitions for user domain
├── README.md                # Project documentation
└── services
    ├── todo-service         # Todo microservice
    └── user-service         # User microservice
```

---

## 2. Architecture Overview

This project follows a microservices architecture within a monorepo setup:

* Each service (`todo-service`, `user-service`) is an independent Go module
* Communication between services is done via gRPC using protobuf contracts (`proto/`)
* Shared logic is extracted into `pkg/`
* `gateways/todo-bff` acts as a BFF layer, aggregating and exposing APIs to clients

---

## 3. Todo Service – Internal Architecture

```
.
├── cmd
│   └── server
│       └── main.go                     # Application entry point
├── di
│   ├── wire_gen.go                     # Generated dependency injection code
│   └── wire.go                         # Dependency injection configuration (Google Wire)
├── go.mod                              # Module definition
├── go.sum                              # Dependency checksum
└── internal
    ├── config                          # Service-specific configuration
    │   └── config.go              
    ├── domain                          # Domain Layer
    │   ├── entity                      # Domain models with business rules
    │   │   ├── todo_enums.go
    │   │   ├── todo.go
    │   │   └── types.go
    │   └── gateway                     # Contracts for data access
    │       ├── todo_commands.go
    │       └── todo_queries.go
    handler                             # Handler Layer
    └── grpc
        ├── mapper
        │   └── todo.go                 # Map between proto ↔ domain
        └── todo
            ├── handler.go              # DI + register gRPC server
            ├── helper.go               # resource parsing, helpers
            ├── create_todo.go
            ├── delete_todo.go
            ├── get_todo.go
            ├── list_todos.go
            └── update_todo.go
    ├── infra
    │   └── persistence
            ├── database.go             # Initialize DB connection (from config)
    │       ├── mapper
    │       │   └── todo.go             # Map between DB model ↔ domain
    │       ├── model
    │       │   └── todo.go             # Database models
    │       ├── todo_commands.go        # Implementation of command operations
    │       ├── todo_queries.go         # Implementation of query operations
    │       └── wire.go                 # DI wiring for persistence layer
    └── usecase                         # Implements business logic orchestration
        └── todos
            ├── input
            │   └── todo.go             # Input DTOs (from handler → usecase)
            ├── output
            │   └── todo.go             # Output DTOs (from usecase → handler)
            ├── todo_creator.go         # Create todo use case (UseCase Implementation)
            ├── todo_deleter.go         # Delete todo use case (UseCase Implementation)
            ├── todo_getter.go          # Get single todo (UseCase Implementation)
            ├── todo_lister.go          # List todos (UseCase Implementation)
            ├── todo_updater.go         # Update todo (UseCase Implementation)
            └── wire.go                 # DI wiring for usecases
```

---

## 4. User Service – Internal Architecture

```
.
├── cmd
│   └── server
│       └── main.go                     # Application entry point
├── di
│   ├── wire_gen.go                     # Generated dependency injection code
│   └── wire.go                         # Dependency injection configuration (Google Wire)
├── go.mod                              # Module definition
├── go.sum                              # Dependency checksum
└── internal
    ├── config                          # Service-specific configuration
    │   └── config.go              
    ├── domain                          # Domain Layer
    │   ├── entity                      # Domain models (pure business objects)
    │   │   └── user.go
    │   └── gateway                     # Contracts (interfaces) for external systems
    │       ├── token_generator.go      # Token abstraction (e.g. JWT)
    │       ├── user_commands.go        # Write operations (create user, etc.)
    │       └── user_queries.go         # Read operations (find user, etc.)
    ├── handler                         # Handler Layer (transport layer)
    │   └── grpc
    │       ├── mapper
    │       │   └── user.go             # Map between proto ↔ domain
    │       └── user
    │           ├── handler.go          # DI + register gRPC server
    │           ├── login.go            # Login RPC handler
    │           └── register.go         # Register RPC handler
    ├── infra                           # Infrastructure Layer
    │   ├── persistence                 # Database implementation
    │   │   ├── database.go             # Initialize DB connection (from config)
    │   │   ├── mapper
    │   │   │   └── user.go             # Map between DB model ↔ domain
    │   │   ├── model
    │   │   │   └── user.go             # Database models (GORM)
    │   │   ├── user_repository.go      # Implements user gateway interfaces
    │   │   └── wire.go                 # DI wiring for persistence layer
    │   └── security                    # Security-related implementations
    │       ├── jwt_manager.go          # JWT implementation (TokenGenerator)
    │       └── wire.go                 # DI wiring for security layer
    └── usecase                         # Implements business logic orchestration
        └── user
            ├── input
            │   └── user.go             # Input DTOs (from handler → usecase)
            ├── output
            │   └── user.go             # Output DTOs (from usecase → handler)
            ├── user_creator.go         # Register user (UseCase Implementation)
            ├── user_authenticator.go   # Login / authentication logic (UseCase Implementation)
            ├── user.go                 # Shared interfaces
            └── wire.go                 # DI wiring for usecases
```

---

## 5. Data Flow (Simplified)

```
Client → BFF → gRPC → Handler → Usecase → Domain → Infrastructure → Database
```

---

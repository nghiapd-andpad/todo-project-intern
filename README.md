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
│       └── main.go                # Application entry point
├── di
│   ├── wire_gen.go                # Generated dependency injection code
│   └── wire.go                    # Dependency injection configuration (Google Wire)
├── go.mod                         # Module definition
├── go.sum                         # Dependency checksum
└── internal
    ├── config                     # Service-specific configuration
    │   └── config.go              
    ├── domain                     # Domain Layer
    │   ├── entity                 # Domain models with business rules
    │   │   ├── todo_enums.go
    │   │   ├── todo.go
    │   │   └── types.go
    │   └── gateway                # Contracts for data access
    │       ├── todo_commands.go
    │       └── todo_queries.go
    handler                    # Handler Layer
    └── grpc
        ├── mapper
        │   └── todo.go        # Map between proto ↔ domain
        └── todo
            ├── handler.go     # DI + register gRPC server
            ├── helper.go      # resource parsing, helpers
            ├── create_todo.go
            ├── delete_todo.go
            ├── get_todo.go
            ├── list_todos.go
            └── update_todo.go
    ├── infra
    │   └── persistence
    │       ├── mapper
    │       │   └── todo.go        # Map between DB model ↔ domain
    │       ├── model
    │       │   └── todo.go        # Database models
    │       ├── todo_commands.go   # Implementation of command operations
    │       ├── todo_queries.go    # Implementation of query operations
    │       └── wire.go            # DI wiring for persistence layer
    └── usecase                    # Implements business logic orchestration
        └── todos
            ├── input
            │   └── todo.go        # Input DTOs (from handler → usecase)
            ├── output
            │   └── todo.go        # Output DTOs (from usecase → handler)
            ├── todo_creator.go    # Create todo use case (UseCase Implementation)
            ├── todo_deleter.go    # Delete todo use case (UseCase Implementation)
            ├── todo_getter.go     # Get single todo (UseCase Implementation)
            ├── todo_lister.go     # List todos (UseCase Implementation)
            ├── todo_updater.go    # Update todo (UseCase Implementation)
            └── wire.go            # DI wiring for usecases
```

---

## 4. Data Flow (Simplified)

```
Client → BFF → gRPC → Handler → Usecase → Domain → Infrastructure → Database
```

---

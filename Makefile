.PHONY: proto

CORE_TODO_DIR := services/core-todo
CORE_TODO_COMPOSE := cd $(CORE_TODO_DIR) && docker compose --profile db

proto-buf:
	buf generate

proto-protoc:
	protoc \
	  --go_out=. \
	  --go_opt=paths=source_relative \
	  --go-grpc_out=. \
	  --go-grpc_opt=paths=source_relative \
	  proto/todo/v1/todo.proto \
	  proto/user/v1/user.proto

wire:
	cd services/core-todo && wire ./di/...
	cd services/core-user && wire ./di/...
	cd services/bff-web && wire ./di/...
	
generate:
	cd services/bff-web && go run github.com/99designs/gqlgen@v0.17.89 generate

run-core-todo:
	cd services/core-todo && go run ./cmd/server/...

run-core-user:
	cd services/core-user && go run ./cmd/server/...

run-bff-web:
	cd services/bff-web && go run ./cmd/server/...

run-core-worker:
	cd services/core-todo && go run ./cmd/worker/...
	
mock-core-todo:
	cd services/core-todo && go generate ./internal/domain/gateway/...

test-usecase-core-todo:
	cd services/core-todo && go test ./internal/usecase/... -v -count=1 -race

test-handler-core-todo:
	cd services/core-todo && go test ./internal/handler/... -v -count=1 -race

test-integration-core-todo:
	cd services/core-todo && \
	go test -v -tags integration ./test/integration/...
	  
test-usecase-core-user:
	cd services/core-user && go test ./internal/usecase/... -v -count=1 -race

test-unit-core-todo-coverage:
	cd services/core-todo && rm -f coverage.out coverage.html && \
    go test ./internal/... -coverprofile=coverage.out -covermode=atomic -count=1 && \
    go tool cover -html=coverage.out -o coverage.html && \
	go tool cover -func=coverage.out

core-todo-db-up:
	$(CORE_TODO_COMPOSE) up -d core-todo-database core-todo-redis

core-todo-db-down:
	$(CORE_TODO_COMPOSE) down

core-todo-db-reset:
	$(CORE_TODO_COMPOSE) down -v
	$(CORE_TODO_COMPOSE) up -d core-todo-database core-todo-redis

core-todo-dbmigrator-build:
	$(CORE_TODO_COMPOSE) build core-todo-dbmigrator

core-todo-dbmigrator-shell:
	$(CORE_TODO_COMPOSE) run --rm core-todo-dbmigrator bash

core-todo-ridgepole-version:
	$(CORE_TODO_COMPOSE) run --rm core-todo-dbmigrator ridgepole --version

core-todo-ridgepole-export:
	$(CORE_TODO_COMPOSE) run --rm core-todo-dbmigrator \
		ridgepole --config database/config/ridgepole.yaml -E local --export -o /tmp/schema.rb

core-todo-ridgepole-diff:
	$(CORE_TODO_COMPOSE) run --rm core-todo-dbmigrator \
		ridgepole --config database/config/ridgepole.yaml -E local --apply --dry-run -f database/schemas/Schemafile

core-todo-ridgepole-apply:
	$(CORE_TODO_COMPOSE) run --rm core-todo-dbmigrator \
		ridgepole --config database/config/ridgepole.yaml -E local --apply -f database/schemas/Schemafile
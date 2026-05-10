.PHONY: proto
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

mock-core-todo:
	cd services/core-todo && go generate ./internal/domain/gateway/...

test-usecase-core-todo:
	cd services/core-todo && go test ./internal/usecase/... -v -count=1 -race

test-handler-core-todo:
	cd services/core-todo && go test ./internal/handler/... -v -count=1 -race

test-integration-core-todo:
	cd services/core-todo && \
	  TEST_DB_USER=root \
	  TEST_DB_PASSWORD=root \
	  TEST_DB_HOST=localhost \
	  TEST_DB_PORT=3306 \
	  go test ./internal/handler/... -v -count=1 -race -timeout 120s -tags integration
	  
test-usecase-core-user:
	cd services/core-user && go test ./internal/usecase/... -v -count=1 -race

test-unit-core-todo-coverage:
	cd services/core-todo && rm -f coverage.out coverage.html && \
    go test ./internal/... -coverprofile=coverage.out -covermode=atomic -count=1 && \
    go tool cover -html=coverage.out -o coverage.html && \
	go tool cover -func=coverage.out

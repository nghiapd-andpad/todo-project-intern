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

test-unit-core-todo:
	cd services/core-todo && go test ./internal/usecase/... -v -count=1 -race

test-unit-core-user:
	cd services/core-user && go test ./internal/usecase/... -v -count=1 -race

test-unit-core-todo-coverage:
	cd services/core-todo && go test ./internal/usecase/... -coverprofile=coverage.out -covermode=atomic -count=1
	cd services/core-todo && gocov convert coverage.out | gocov-html > coverage.html

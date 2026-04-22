.PHONY: proto
proto:
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
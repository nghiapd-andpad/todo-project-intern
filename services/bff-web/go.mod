module github.com/nghiapd-andpad/todo-project-intern/services/bff-web

go 1.26.1

replace github.com/nghiapd-andpad/todo-project-intern/pkg => ../../pkg

replace github.com/nghiapd-andpad/todo-project-intern/proto => ../../proto

require (
	github.com/99designs/gqlgen v0.17.89
	github.com/google/wire v0.7.0
	github.com/nghiapd-andpad/todo-project-intern/proto v0.0.0-00010101000000-000000000000
	github.com/vektah/gqlparser/v2 v2.5.32
	google.golang.org/grpc v1.80.0
)

require (
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/sosodev/duration v1.4.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

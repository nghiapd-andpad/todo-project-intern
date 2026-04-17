module github.com/nghiapd-andpad/todo-project-intern/services/core-user

go 1.26.1

require (
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/wire v0.7.0
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/nghiapd-andpad/todo-project-intern/proto v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.50.0
	google.golang.org/grpc v1.80.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.31.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/nghiapd-andpad/todo-project-intern/proto => ../../proto

replace github.com/nghiapd-andpad/todo-project-intern/pkg => ../../pkg

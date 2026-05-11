# REVIEW 2026-05-04
## PROTO
- Cach dat ten request va response phai dong nhat voi nhau [todo.proto](proto/todo/v1/todo.proto)
  - RPC la GetTodoList thi request la GetTodoListRequest va response la GetTodoListResponse
- Comment phai dat o tren message hoac filed 
- Chua validate required field trong proto file


## CODE 
- Qua nhieu params trong func, vi du [auth_registerer.go](services/bff-web/internal/usecase/auth/auth_registerer.go)
- Mot so validate nen dat trong proto thay vi trong code, vi du
  - validate user name, email, password khac empty 
- File wire gen phai bao gom toan bo dependency cua service, khong tach rieng nhieu file, khong chi la usecase, vi du [wire.go](services/core-todo/internal/infra/persistence/wire.go)
  - Phai chua binding giua interface va implementation, khong chi la func new
- Cac func new chi nen tao instance cua struct khong nen tra ra interface, vi du [todo_commands.go](services/core-todo/internal/infra/persistence/todo_commands.go)
  - Co the them dong lenh de dam bao struct da implement interface, vi du `var _ gateway.TodoCommandsGateway = (*todoCommandsGateway)(nil)`
- Trong unit test khong nen mock anything, nen test mock object thuc te duoc goi va tra ve ket qua thuc te, vi du [todo_creator_test.go](services/core-todo/internal/usecase/todos/todo_creator_test.go)


# REVIEW 2026-05-08
## CODE
- Dang ignore mot so loi  tiem nang khi bypass error/validate, vi du [todo.resolvers.go](services/bff-web/internal/handler/graph/todo.resolvers.go)
- van con validate trong business code thay vi co the validate trong proto [auth_loginer.go](services/bff-web/internal/usecase/auth/auth_loginer.go)
- Input params va output khong dong nhat ve kieu du lieu, output la pointer nhung input la value, can dong nhat kieu du lieu, vi du [auth_loginer.go](services/bff-web/internal/usecase/auth/auth_loginer.go)
- van con nhieu file wire trong mot service, can gom chung vao 1 file, vi du [wire.go](services/bff-web/internal/usecase/todo/wire.go), [wire.go](services/bff-web/internal/usecase/auth/wire.go)
- dat ten package chua dung chuan Go vi du [auth_client.go](services/bff-web/internal/infra/grpc_client/auth_client.go)
- Interface va implement cua no dang dat trong cung 1 file, can tach ra, vi du  [todo_queries.go](services/core-todo/internal/domain/gateway/todo_queries.go)
- Unit test: follow Remote service
- 

# REVIEW 2026-05-11
- Trong resolver nen dat ten doi utong/field muon lay gia tri luon, vi du [resolver.go](services/bff-web/internal/handler/graph/resolver.go)
- Co the dua mot so gia tri vao trong config thay vi hardcode, vi du [user_authenticator.go](services/core-user/internal/usecase/user/user_authenticator.go)
- Su dung mot so method da bi deprecated, vi du: [auth_client.go](services/bff-web/internal/infra/grpcclient/auth_client.go)
- Van con mot so cho dang ignore error, vi du: [todo_client.go](services/bff-web/internal/infra/grpcclient/todo_client.go)
- Follow cach dat ten cho mot so thuat ngu trong Go, vi du: [jwt.go](services/bff-web/internal/infra/jwt/jwt.go)
- Van con co business logic trong gateway interface: [jwt.go](services/bff-web/internal/infra/jwt/jwt.go)
- Bo sung feature flag/blacklist trong API create user/todo va update unit test --> muc dich de xem cach implement va cach xu ly unit test khi on/off flag
  - Trong unit test, neu can truyen config vao trong tung bo test 

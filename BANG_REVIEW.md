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
- 

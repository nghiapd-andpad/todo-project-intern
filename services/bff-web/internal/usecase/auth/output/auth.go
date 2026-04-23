package output

import "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"

type LoginOutput struct {
	AccessToken string
	User        *entity.User
}

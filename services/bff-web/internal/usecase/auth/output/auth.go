package output

import "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain"

type LoginOutput struct {
	AccessToken string
	User        *domain.User
}

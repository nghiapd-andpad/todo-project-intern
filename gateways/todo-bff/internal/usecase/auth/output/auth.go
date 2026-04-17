package output

import "github.com/nghiapd-andpad/todo-project-intern/gateways/todo-bff/internal/domain"

type LoginOutput struct {
	AccessToken string
	User        *domain.User
}

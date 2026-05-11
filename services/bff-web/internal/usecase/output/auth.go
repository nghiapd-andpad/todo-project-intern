// Package output defines the data structures for output results returned by use cases.
package output

import "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"

type LoginOutput struct {
	AccessToken string
	User        *entity.User
}

type RegisterOutput struct {
	Username string
	Email    string
}

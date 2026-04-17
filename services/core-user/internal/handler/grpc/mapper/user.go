package mapper

import (
	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user/output"
)

func UserToPb(u *output.UserDTO) *userv1.User {
	if u == nil {
		return nil
	}
	return &userv1.User{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}

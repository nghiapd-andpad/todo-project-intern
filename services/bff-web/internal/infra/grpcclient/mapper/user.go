package mapper

import (
	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
)

func UserFromPb(pb *userv1.User) *entity.User {
	if pb == nil {
		return nil
	}
	return &entity.User{
		ID:       pb.Id,
		Username: pb.Username,
		Email:    pb.Email,
	}
}

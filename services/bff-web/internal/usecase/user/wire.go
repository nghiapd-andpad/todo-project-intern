package user

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewUserUsecase,
	ProvideUserGetter,
)

func ProvideUserGetter(u *userUsecase) UserGetter { return u }

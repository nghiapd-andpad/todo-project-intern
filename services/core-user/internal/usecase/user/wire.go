package user

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewUserCreator,
	NewUserAuthenticator,
)

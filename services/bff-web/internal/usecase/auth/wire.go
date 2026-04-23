package auth

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewRegisterer,
	NewLoginer,
)

package jwt

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(NewJwtManager)

package infra

import (
	"github.com/google/wire"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/security"
)

var WireSet = wire.NewSet(
	persistence.WireSet,
	security.WireSet,
)

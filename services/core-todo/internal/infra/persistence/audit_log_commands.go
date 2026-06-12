package persistence

import (
	"context"
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/model"
)

type AuditLogCommandsGateway struct {
	db *gorm.DB
}

func NewAuditLogCommandsGateway(db *gorm.DB) *AuditLogCommandsGateway {
	return &AuditLogCommandsGateway{db: db}
}

var _ gateway.AuditLogCommandsGateway = (*AuditLogCommandsGateway)(nil)

func (g *AuditLogCommandsGateway) Create(ctx context.Context, in *gatewayinput.CreateAuditLog) error {
	conn := connFromContext(ctx, g.db)

	m := &model.AuditLog{
		EventName:  in.EventName,
		EntityType: in.EntityType,
		EntityID:   in.EntityID,
		ActorID:    int64(in.ActorID),
		Payload:    datatypes.JSON(in.Payload),
	}

	if err := conn.Create(m).Error; err != nil {
		return fmt.Errorf("db create audit log: %w", err)
	}

	return nil
}

package model

import (
	"time"

	"gorm.io/datatypes"
)

type AuditLog struct {
	ID         int64          `gorm:"primaryKey;autoIncrement"`
	EventName  string         `gorm:"column:event_name;size:100;not null;index"`
	EntityType string         `gorm:"column:entity_type;size:50;not null;index:idx_audit_logs_entity"`
	EntityID   int64          `gorm:"column:entity_id;not null;index:idx_audit_logs_entity"`
	ActorID    int64          `gorm:"column:actor_id;not null;index"`
	Payload    datatypes.JSON `gorm:"column:payload;type:json;not null"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

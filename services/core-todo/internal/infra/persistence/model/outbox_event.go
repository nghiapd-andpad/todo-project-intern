package model

import (
	"time"

	"gorm.io/datatypes"
)

type OutboxEvent struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	EventName   string         `gorm:"column:event_name;size:100;not null"`
	RoutingKey  string         `gorm:"column:routing_key;size:100;not null;index"`
	Payload     datatypes.JSON `gorm:"column:payload;type:json;not null"`
	Status      string         `gorm:"column:status;size:20;not null;index:idx_outbox_events_status_created_at"`
	RetryCount  int            `gorm:"column:retry_count;not null;default:0"`
	LastError   *string        `gorm:"column:last_error;type:text"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime;index:idx_outbox_events_status_created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	PublishedAt *time.Time     `gorm:"column:published_at"`
}

func (OutboxEvent) TableName() string {
	return "outbox_events"
}

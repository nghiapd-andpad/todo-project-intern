package model

import "time"

type Notification struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`

	ReceiverID   int64     `gorm:"column:receiver_id;not null;index:idx_notifications_receiver_is_read;uniqueIndex:uq_notifications_idempotency"`
	ResourceType string    `gorm:"column:resource_type;size:50;not null;uniqueIndex:uq_notifications_idempotency"`
	ResourceID   int64     `gorm:"column:resource_id;not null;uniqueIndex:uq_notifications_idempotency"`
	EventName    string    `gorm:"column:event_name;size:100;not null;uniqueIndex:uq_notifications_idempotency"`
	OccurredAt   time.Time `gorm:"column:occurred_at;not null;uniqueIndex:uq_notifications_idempotency"`

	Title   string `gorm:"column:title;size:255;not null"`
	Content string `gorm:"column:content;type:text;not null"`

	IsRead bool `gorm:"column:is_read;not null;default:false;index:idx_notifications_receiver_is_read"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Notification) TableName() string {
	return "notifications"
}

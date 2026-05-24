package model

import "time"

type IdempotencyKey struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	UserID         int64     `gorm:"column:user_id;not null;uniqueIndex:uk_idempotency_scope"`
	Operation      string    `gorm:"column:operation;size:100;not null;uniqueIndex:uk_idempotency_scope"`
	IdempotencyKey string    `gorm:"column:idempotency_key;size:255;not null;uniqueIndex:uk_idempotency_scope"`
	RequestHash    string    `gorm:"column:request_hash;size:64;not null"`
	Status         string    `gorm:"column:status;size:20;not null"`
	ResourceID     *int64    `gorm:"column:resource_id"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
	ExpiresAt      time.Time `gorm:"column:expires_at;not null;index"`
}

func (IdempotencyKey) TableName() string {
	return "idempotency_keys"
}

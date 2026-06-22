package model

import "time"

type ProcessedEvent struct {
	EventHash   string    `gorm:"primaryKey;column:event_hash;size:64"`
	ConsumerKey string    `gorm:"primaryKey;column:consumer_key;size:64"`
	ProcessedAt time.Time `gorm:"column:processed_at;autoCreateTime"`
}

func (ProcessedEvent) TableName() string {
	return "processed_events"
}

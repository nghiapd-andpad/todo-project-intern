package model

import (
	"time"

	"gorm.io/gorm"
)

type TodoList struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"column:name;size:255;not null"`
	OwnerID   int64          `gorm:"column:owner_id;not null;index"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (TodoList) TableName() string {
	return "todo_lists"
}

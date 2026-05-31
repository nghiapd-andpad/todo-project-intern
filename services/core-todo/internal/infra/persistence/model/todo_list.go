package model

import (
	"time"

	"gorm.io/gorm"
)

type TodoList struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"column:name;size:255;not null;uniqueIndex:uk_owner_todo_list_name"`
	OwnerID   int64          `gorm:"column:owner_id;not null;index;uniqueIndex:uk_owner_todo_list_name"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Version   int64          `gorm:"column:version;not null;default:1"`
}

func (TodoList) TableName() string {
	return "todo_lists"
}

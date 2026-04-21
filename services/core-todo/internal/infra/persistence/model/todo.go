package model

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	TodoListID  int64          `gorm:"column:todo_list_id;index"`
	Title       string         `gorm:"column:title;size:255;not null"`
	Description *string        `gorm:"column:description;type:text"`
	Status      string         `gorm:"column:status;size:50;not null;index"`
	Priority    string         `gorm:"column:priority;size:50;not null;index"`
	DueDate     *time.Time     `gorm:"column:due_date"`
	CreatorID   int64          `gorm:"column:creator_id;not null;index"`
	AssigneeID  *int64         `gorm:"column:assignee_id;index"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Todo) TableName() string {
	return "todos"
}

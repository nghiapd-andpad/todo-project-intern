// Package model contains GORM models for Todo domain.
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
	Priority    string         `gorm:"column:priority;size:50;not null;index"`
	AssigneeID  *int64         `gorm:"column:assignee_id;index"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	Status      string         `gorm:"column:status;size:50;not null;index;index:idx_todos_overdue_marker,priority:2"`
	DueDate     *time.Time     `gorm:"column:due_date;index:idx_todos_overdue_marker,priority:3"`
	DeletedAt   gorm.DeletedAt `gorm:"index;index:idx_todos_overdue_marker,priority:1"`
}

func (Todo) TableName() string {
	return "todos"
}

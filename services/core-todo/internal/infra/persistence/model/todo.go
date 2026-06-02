// Package model contains GORM models for Todo domain.
package model

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	TodoListID  int64          `gorm:"column:todo_list_id;not null;index;uniqueIndex:uk_todos_business_duplicate"`
	Title       string         `gorm:"column:title;size:255;not null;uniqueIndex:uk_todos_business_duplicate"`
	Description *string        `gorm:"column:description;type:text"`
	Priority    string         `gorm:"column:priority;size:50;not null;index"`
	AssigneeID  *int64         `gorm:"column:assignee_id;index;uniqueIndex:uk_todos_business_duplicate"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	Status      string         `gorm:"column:status;size:50;not null;index;index:idx_todos_overdue_marker,priority:2;uniqueIndex:uk_todos_business_duplicate"`
	DueDate     *time.Time     `gorm:"column:due_date;index:idx_todos_overdue_marker,priority:3"`
	DeletedAt   gorm.DeletedAt `gorm:"index;index:idx_todos_overdue_marker,priority:1"`
}

func (Todo) TableName() string {
	return "todos"
}

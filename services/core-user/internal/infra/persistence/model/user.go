package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             int64  `gorm:"primaryKey;autoIncrement"`
	Username       string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Email          string `gorm:"type:varchar(191);uniqueIndex;not null"`
	HashedPassword string `gorm:"type:varchar(255);not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

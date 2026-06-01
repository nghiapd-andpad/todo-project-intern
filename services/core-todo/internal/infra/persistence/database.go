// Package persistence provides database connection and data access layer for the core-todo service
package persistence

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
)

func NewDatabase(cfg *config.Config) (*gorm.DB, func(), error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect database: %w", err)
	}

	cleanup := func() {
		sqlDB, err := db.DB()
		if err == nil {
			fmt.Println("Closing database connection...")
			sqlDB.Close()
		}
	}

	return db, cleanup, nil
}

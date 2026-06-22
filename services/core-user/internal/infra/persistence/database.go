// Package persistence provides the implementation for database interactions using GORM for the core-user service.
package persistence

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence/model"
)

func NewDatabase(cfg *config.Config) (*gorm.DB, func(), error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Auto Migrate
	if err := db.AutoMigrate(&model.User{}, &model.Notification{}, &model.OutboxEvent{}, &model.ProcessedEvent{}); err != nil {
		return nil, nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	cleanup := func() {
		sqlDB, err := db.DB()
		if err == nil {
			fmt.Println("Closing User Database connection...")
			sqlDB.Close()
		}
	}

	return db, cleanup, nil
}

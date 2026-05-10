// Package testutil provides test helpers for persistence layer testing.
package testutil

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/model"
)

// NewTestDB creates an isolated database for each test: creates a new database, migrates schema, and drops it after test.
func NewTestDB(t *testing.T, cfg *config.Config) *gorm.DB {
	t.Helper()

	// create 4 random bytes
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		t.Fatalf("failed to generate random string: %v", err)
	}
	randomSuffix := hex.EncodeToString(randomBytes)

	dbName := fmt.Sprintf("todo_test_%d_%s", time.Now().UnixNano(), randomSuffix)

	adminDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
	)

	adminDB, err := gorm.Open(mysql.Open(adminDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("connect admin db: %v", err)
	}

	// create test database
	if err := adminDB.Exec(fmt.Sprintf(
		"CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		dbName,
	)).Error; err != nil {
		t.Fatalf("create test db: %v", err)
	}

	sqlDB, _ := adminDB.DB()
	sqlDB.Close()

	// connect to test database
	testDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		dbName,
	)

	db, err := gorm.Open(mysql.Open(testDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("connect test db: %v", err)
	}

	// migrate schema
	if err := db.AutoMigrate(&model.Todo{}, &model.TodoList{}); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	// cleanup after test
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()

		dropDB, err := gorm.Open(mysql.Open(adminDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			t.Logf("connect admin for cleanup: %v", err)
			return
		}

		dropDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName))

		sqlDB2, _ := dropDB.DB()
		sqlDB2.Close()
	})

	return db
}

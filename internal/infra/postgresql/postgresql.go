package postgresql

import (
	"fmt"

	"github.com/jevvonn/ws-cicd-bcc2026/pkg/helper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		helper.Env("DB_HOST", "localhost"),
		helper.Env("DB_PORT", "5432"),
		helper.Env("DB_USER", "postgres"),
		helper.Env("DB_PASSWORD", "postgres"),
		helper.Env("DB_NAME", "todo_db"),
		helper.Env("DB_SSLMODE", "disable"),
	)

	logLevel := logger.Silent
	if helper.Env("APP_ENV", "development") == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db, nil
}

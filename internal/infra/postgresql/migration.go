package postgresql

import (
	"github.com/jevvonn/ws-cicd-bcc2026/internal/app/todo/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&entity.Todo{})
}

func Rollback(db *gorm.DB) error {
	return db.Migrator().DropTable(&entity.Todo{})
}

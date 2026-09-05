package bootstrap

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jevvonn/ws-cicd-bcc2026/internal/infra/postgresql"
	"github.com/jevvonn/ws-cicd-bcc2026/pkg/helper"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      helper.Env("APP_NAME", "ws-cicd-bcc2026"),
		ErrorHandler: errorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New())

	registerRoutes(app, db)

	return app
}

func Run() error {
	LoadEnv()

	db, err := postgresql.New()
	if err != nil {
		return err
	}

	app := NewApp(db)

	return app.Listen(":" + helper.Env("APP_PORT", "3000"))
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
	}

	return helper.Error(c, code, err.Error())
}

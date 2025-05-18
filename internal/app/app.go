package app

import (
	"context"
	"profkom/config"
	"profkom/internal/binder"
	"profkom/internal/repository"
	"profkom/internal/service"

	"profkom/pkg/postgres"
	"profkom/pkg/s3"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	txmanager "github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/gofiber/fiber/v2"
)

func Run(ctx context.Context, cfg *config.Config) (err error) {
	postgres, err := postgres.NewDB(cfg.Postgres)
	if err != nil {
		return err
	}

	repo := repository.New(postgres, trmsqlx.DefaultCtxGetter)

	txManager, err := txmanager.New(trmsqlx.NewDefaultFactory(postgres))
	if err != nil {
		return err
	}

	storage, err := s3.New(cfg.S3)
	if err != nil {
		return err
	}

	service := service.New(cfg.Services, repo, txManager, storage)

	handler := binder.NewHandler(service)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		if origin != "" {
			c.Set("Access-Control-Allow-Origin", origin)
			c.Set("Access-Control-Allow-Credentials", "true")
			c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie, X-TradeUnion-ID")
		}
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	mw := binder.New("asdqwe2131241eqeqw", service.Auth)

	binder := binder.NewBinder(app, handler, mw)
	binder.BindRoutes()

	return app.Listen(":8080")
}

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
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	binder := binder.NewBinder(app, handler)
	binder.BindRoutes()

	return app.Listen(":8080")
}

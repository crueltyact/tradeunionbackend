package config

import (
	"profkom/internal/service"
	"profkom/pkg/postgres"
	"profkom/pkg/s3"
)

type Config struct {
	Postgres postgres.Config
	S3       s3.Config
	Services service.Config
}

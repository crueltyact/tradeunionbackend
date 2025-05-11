package auth

import (
	"profkom/internal/repository/auth"
	"profkom/pkg/s3"

	txmanager "github.com/avito-tech/go-transaction-manager/trm/manager"
)


type Service struct {
	repo      *auth.Repository
	txManager *txmanager.Manager
	s3        *s3.Client
	cfg       Config
}

func New(cfg Config, repo *auth.Repository, txManager *txmanager.Manager, s3 *s3.Client) *Service {
	return &Service{
		repo:      repo,
		txManager: txManager,
		s3:        s3,
		cfg:       cfg,
	}
}

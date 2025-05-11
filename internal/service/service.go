package service

import (
	txmanager "github.com/avito-tech/go-transaction-manager/trm/manager"

	"profkom/internal/repository"
	"profkom/internal/service/auth"
	"profkom/internal/service/chat"
	"profkom/internal/service/documents"
	"profkom/internal/service/guide"
	"profkom/internal/service/news"
	"profkom/internal/service/projects"
	"profkom/pkg/s3"
)

type Service struct {
	Guide     *guide.Service
	Project   *projects.Service
	News      *news.Service
	Auth      *auth.Service
	Documents *documents.Service
	Chat      *chat.Service
}

func New(cfg Config, repository *repository.Repository, txManager *txmanager.Manager, s3storage *s3.Client) *Service {
	return &Service{
		Guide:     guide.New(repository.Guide),
		Project:   projects.New(repository.Project, s3storage),
		News:      news.New(repository.News, s3storage),
		Auth:      auth.New(cfg.Auth, repository.Auth, txManager, s3storage),
		Documents: documents.New(repository.Documents, s3storage),
		Chat:      chat.New(repository.Chat, txManager, s3storage),
	}
}

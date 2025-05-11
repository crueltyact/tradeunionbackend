package documents

import (
	"profkom/internal/repository/documents"
	"profkom/pkg/s3"
)

type Service struct {
	repo *documents.Repository
	s3   *s3.Client
}

func New(repository *documents.Repository, s3 *s3.Client) *Service {
	return &Service{
		repo: repository,
		s3:   s3,
	}
}

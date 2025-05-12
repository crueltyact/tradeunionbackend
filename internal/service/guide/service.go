package guide

import "profkom/internal/repository/guide"

type Service struct {
	repo *guide.Repository
}

func New(repo *guide.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

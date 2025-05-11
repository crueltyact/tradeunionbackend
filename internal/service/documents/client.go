package documents

import (
	"context"
	"profkom/internal/entities"
)

func (s *Service) GetDocuemnts(ctx context.Context) ([]entities.Document, error) {
	return s.repo.SelectDocuments(ctx)
}

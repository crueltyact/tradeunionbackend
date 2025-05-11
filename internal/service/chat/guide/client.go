package guide

import (
	"context"

	"profkom/internal/models"
)

func (s *Service) GetGuide(ctx context.Context) (models.AllGuides, error) {
	return s.repo.SelectGuide(ctx)
}

func (s *Service) InsertGuide(ctx context.Context, guidesType string, guides []models.Guide) error {
	return s.repo.InsertGuide(ctx, guidesType, guides)
}

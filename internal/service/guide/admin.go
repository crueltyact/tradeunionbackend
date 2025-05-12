package guide

import (
	"context"
	"profkom/internal/models"
)

func (s *Service) DeleteGuide(ctx context.Context, id int) (err error) {
	return s.repo.DeleteGuide(ctx, id)
}

func (s *Service) UpdateGuide(ctx context.Context) (err error) {
	return err
}

func (s *Service) DeleteTheme(ctx context.Context, id int) (err error) {
	return s.repo.DeleteTheme(ctx, id)
}

func (s *Service) CreateTheme(ctx context.Context, req models.PostThemeRequest) (err error) {
	err = s.repo.InsertTheme(ctx, req)
	if err != nil {
		return err
	}

	return err
}

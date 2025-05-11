package news

import (
	"context"
	"profkom/internal/models"

	"github.com/google/uuid"
)

func (s *Service) GetNews(ctx context.Context) (models.News, error) {
	return s.repository.SelectNews(ctx)
}

func (s *Service) GetNew(ctx context.Context, id string) (new models.New, err error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return new, err
	}

	new, err = s.repository.SelectNew(ctx, uuid)
	if err != nil {
		return new, err
	}

	return new, err
}

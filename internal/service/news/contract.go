package news

import (
	"context"

	"github.com/google/uuid"

	"profkom/internal/entities"
	"profkom/internal/models"
)

type (
	repository interface {
		SelectNews(ctx context.Context) (news models.News, err error)
		InsertNew(ctx context.Context, new *entities.New) (err error)
		SelectNew(ctx context.Context, id uuid.UUID) (new models.New, err error)
		DeleteNew(ctx context.Context, id uuid.UUID) (err error)
	}
)

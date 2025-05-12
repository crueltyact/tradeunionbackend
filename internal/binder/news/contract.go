package news

import (
	"context"
	"profkom/internal/models"

	"github.com/google/uuid"
)

type (
	service interface {
		GetNews(ctx context.Context) (models.News, error)
		UploadNews(ctx context.Context, request models.PostNewRequest) (err error)
		GetNew(ctx context.Context, id string) (new models.New, err error)
		DeleteNew(ctx context.Context, newID uuid.UUID) (err error)
	}
)

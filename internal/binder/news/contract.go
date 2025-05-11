package news

import (
	"context"
	"profkom/internal/models"
)

type (
	service interface {
		GetNews(ctx context.Context) (models.News, error)
		UploadNews(ctx context.Context, request models.PostNewRequest) (err error)
		GetNew(ctx context.Context, id string) (new models.New, err error)
	}
)

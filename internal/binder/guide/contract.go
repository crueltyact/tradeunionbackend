package guide

import (
	"context"
	"profkom/internal/models"
)

type (
	service interface {
		GetGuide(ctx context.Context) (models.AllGuides, error)
		InsertGuide(ctx context.Context, guidesType string, guides []models.Guide) error
		DeleteGuide(ctx context.Context, id int) (err error)
		DeleteTheme(ctx context.Context, id int) (err error)
	}
)

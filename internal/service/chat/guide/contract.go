package guide

import (
	"context"
	"profkom/internal/models"
)

type repository interface {
	SelectGuide(ctx context.Context) (result models.AllGuides, err error)
	InsertGuide(ctx context.Context, guideType string, guides []models.Guide) (err error)
	DeleteGuide(ctx context.Context, id int) (err error)
	DeleteTheme(ctx context.Context, id int) (err error)
}

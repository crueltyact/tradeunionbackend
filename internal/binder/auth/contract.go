package auth

import (
	"context"
	"profkom/internal/models"
)

type (
	service interface {
		AdminSingUp(ctx context.Context, req models.SignUpRequest) (response models.SignUpResponse, err error)
	}
)

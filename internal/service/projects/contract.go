package projects

import (
	"context"
	"profkom/internal/entities"
	"profkom/internal/models"

	"github.com/google/uuid"
)

type (
	repository interface {
		SelectProjects(ctx context.Context) (result models.GetProjectsResponse, err error)
		SelectProject(ctx context.Context, id uuid.UUID) (result models.Project, err error)
		InsertProject(ctx context.Context, project entities.Project) (err error)
	}
)

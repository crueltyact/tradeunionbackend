package projects

import (
	"context"
	"profkom/internal/models"
)

type (
	service interface {
		GetProjects(ctx context.Context) (porjects models.GetProjectsResponse, err error)
		GetProject(ctx context.Context, id string) (project models.Project, err error)
		UplaodProject(ctx context.Context, req models.PostProjectRequest) (err error)
	}
)

package projects

import (
	"context"
	"profkom/internal/models"
	"profkom/internal/repository/projects"
	"profkom/pkg/s3"

	"github.com/google/uuid"
)

type Service struct {
	repository *projects.Repository
	s3         *s3.Client
}

func New(repository *projects.Repository, s3 *s3.Client) *Service {
	return &Service{
		repository: repository,
		s3:         s3,
	}
}

func (s *Service) GetProjects(ctx context.Context) (porjects models.GetProjectsResponse, err error) {
	return s.repository.SelectProjects(ctx)
}

func (s *Service) GetProject(ctx context.Context, id string) (project models.Project, err error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return project, err
	}

	project, err = s.repository.SelectProject(ctx, uuid)
	if err != nil {
		return project, err
	}

	return project, err
}

func (s *Service) RemoveProject(ctx context.Context, id string) (err error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = s.repository.DeleteProject(ctx, uuid)
	if err != nil {
		return err
	}

	return err
}

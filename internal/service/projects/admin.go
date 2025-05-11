package projects

import (
	"bytes"
	"context"
	"fmt"
	"profkom/internal/entities"
	"profkom/internal/models"
	"time"

	"github.com/google/uuid"
)

const (
	filePattern = "%s-%s"
	fileURL     = "%s/%s"
	url         = "https://storage.yandexcloud.net/profkom-dev"
)

func (s *Service) UplaodProject(ctx context.Context, req models.PostProjectRequest) (err error) {
	filename := fmt.Sprintf(filePattern, time.Now(), req.Image.Filename)

	reader := bytes.NewReader(req.Image.Bytes)

	err = s.s3.UploadFile(ctx, &filename, reader)
	if err != nil {
		return err
	}

	new := entities.Project{
		ID:          uuid.New(),
		Intro:       req.Intro,
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		ImageUrl:    fmt.Sprintf(fileURL, url, filename),
	}

	err = s.repository.InsertProject(ctx, new)
	if err != nil {
		return err
	}

	return err
}

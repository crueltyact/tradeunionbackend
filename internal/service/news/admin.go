package news

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

func (s *Service) UploadNews(ctx context.Context, request models.PostNewRequest) (err error) {
	filename := fmt.Sprintf(filePattern, time.Now(), request.Image.Filename)

	reader := bytes.NewReader(request.Image.Bytes)

	err = s.s3.UploadFile(ctx, &filename, reader)
	if err != nil {
		return err
	}

	new := entities.New{
		ID:       uuid.New(),
		Title:    request.Title,
		Content:  request.Content,
		ImageUrl: fmt.Sprintf(fileURL, url, filename),
	}

	err = s.repository.InsertNew(ctx, &new)
	if err != nil {
		return err
	}

	return err
}

func (s *Service) DeleteNew(ctx context.Context, newID uuid.UUID) (err error) {
	return s.repository.DeleteNew(ctx, newID)
}

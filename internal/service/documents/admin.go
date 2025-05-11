package documents

import (
	"bytes"
	"context"
	"fmt"
	"profkom/internal/entities"
	"profkom/internal/models"
	"strings"
	"sync"

	"github.com/google/uuid"
)

const (
	fileURL = "%s/%s"
	url     = "https://storage.yandexcloud.net/profkom-dev"
)

func (s *Service) UploadDocuments(ctx context.Context, req models.PostDocumentRequest) (err error) {
	chanErr := make(chan error)

	documents := entities.DocumentBatch{
		ID:    make([]uuid.UUID, 0, len(req.Documents)),
		URL:   make([]string, 0, len(req.Documents)),
		Title: make([]string, 0, len(req.Documents)),
		Type:  req.Type,
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for _, document := range req.Documents {
		wg.Add(1)

		go func() {
			defer wg.Done()
			reader := bytes.NewReader(document.Bytes)

			format := strings.Split(document.Filename, ".")
			if len(format) < 2 {
				chanErr <- fmt.Errorf("no file format")
				return
			}

			uuid := uuid.New()

			idStr := uuid.String() + "." + format[1]

			err = s.s3.UploadFile(ctx, &idStr, reader)
			if err != nil {
				chanErr <- err
				return
			}

			mu.Lock()
			defer mu.Unlock()
			documents.ID = append(documents.ID, uuid)
			documents.URL = append(documents.URL, fmt.Sprintf(fileURL, url, idStr))
			documents.Title = append(documents.Title, document.Filename)
		}()
	}

	go func() {
		wg.Wait()
		close(chanErr)
	}()

	for err := range chanErr {
		return err
	}

	err = s.repo.InsertDocumnets(ctx, documents)
	if err != nil {
		return err
	}

	return err
}

func (s *Service) DeleteDocument(ctx context.Context, id string) (err error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = s.repo.DeleteDocument(ctx, uuid)
	if err != nil {
		return err
	}

	return err
}

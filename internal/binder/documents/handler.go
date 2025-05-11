package documents

import (
	"fmt"
	"io"
	"log"
	"profkom/internal/models"
	"profkom/internal/service/documents"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *documents.Service
}

func New(service *documents.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) PostDocument(c *fiber.Ctx) error {
	request := models.PostDocumentRequest{}

	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid multipart form")
	}

	files := form.File["documents"]
	log.Println(files)
	if len(files) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "No files uploaded")
	}

	request.Type = c.FormValue("type")
	switch request.Type {
	case "worker":
	case "student":
	default:
		return fmt.Errorf("fail type")
	}

	chanErr := make(chan error)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for _, fileHeader := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()

			file, err := fileHeader.Open()
			if err != nil {
				chanErr <- err
				return
			}
			defer file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				chanErr <- err
				return
			}

			mu.Lock()
			defer mu.Unlock()
			request.Documents = append(request.Documents, models.File{
				Filename: fileHeader.Filename,
				Bytes:    data,
			})
		}()
	}

	go func() {
		wg.Wait()
		close(chanErr)
	}()

	for err := range chanErr {
		return err
	}

	err = h.service.UploadDocuments(c.Context(), request)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) GetDocuments(c *fiber.Ctx) error {
	resp, err := h.service.GetDocuemnts(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

func (h *Handler) DeleteDocument(c *fiber.Ctx) error {
	id := c.Params("document_id")

	err := h.service.DeleteDocument(c.Context(), id)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

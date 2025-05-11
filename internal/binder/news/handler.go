package news

import (
	"io"
	"profkom/internal/models"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service service
}

func New(service service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetNews(c *fiber.Ctx) error {
	resp, err := h.service.GetNews(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

func (h *Handler) PostNews(c *fiber.Ctx) error {
	request := models.PostNewRequest{
		Title:   c.FormValue("title"),
		Content: c.FormValue("content"),
	}

	file, err := c.FormFile("image")
	if err != nil {
		return err
	}

	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	request.Image = models.File{
		Filename: file.Filename,
		Bytes:    fileBytes,
	}

	err = h.service.UploadNews(c.Context(), request)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) GetNew(c *fiber.Ctx) error {
	id := c.Params("new_id")

	resp, err := h.service.GetNew(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

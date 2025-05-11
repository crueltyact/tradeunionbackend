package projects

import (
	"io"
	"profkom/internal/models"
	"profkom/internal/service/projects"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *projects.Service
}

func New(service *projects.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetProjects(c *fiber.Ctx) error {
	resp, err := h.service.GetProjects(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

func (h *Handler) GetProject(c *fiber.Ctx) error {
	id := c.Params("project_id")

	resp, err := h.service.GetProject(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

func (h *Handler) PostProject(c *fiber.Ctx) error {
	request := models.PostProjectRequest{
		Title:       c.FormValue("title"),
		Intro:       c.FormValue("intro"),
		Description: c.FormValue("description"),
		Type:        c.FormValue("type"),
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

	err = h.service.UplaodProject(c.Context(), request)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) DeleteProject(c *fiber.Ctx) error {
	id := c.Params("project_id")

	err := h.service.RemoveProject(c.Context(), id)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

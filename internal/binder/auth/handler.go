package auth

import (
	"errors"
	"io"
	"profkom/internal/models"
	"profkom/internal/service/auth"
	"profkom/pkg/consts"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	service *auth.Service
}

func New(service *auth.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) SignUp(c *fiber.Ctx) error {
	var request models.SignUpRequest

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	resp, err := h.service.AdminSingUp(c.Context(), request)
	if err != nil {
		return err
	}

	cookie := &fiber.Cookie{
		Name:     "token",
		Value:    resp.JwtToken, // предположим, что в ответе есть поле Token
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,   // нельзя читать из JS
		Secure:   false,  // true в проде по HTTPS
		SameSite: "None", // обязательно для кросс-доменных запросов
	}

	c.Cookie(cookie)

	return c.JSON(resp)
}

func (h *Handler) PostInviteToken(c *fiber.Ctx) error {
	var request models.PostInviteTokenRequest

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	resp, err := h.service.CreateInviteToken(c.Context(), request)
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

func (h *Handler) SignIn(c *fiber.Ctx) error {
	var request models.AdminSignInRequest

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	resp, err := h.service.AdminSignIn(c.Context(), request)
	if err != nil {
		return err
	}

	cookie := &fiber.Cookie{
		Name:     "token",
		Value:    resp.Token, // предположим, что в ответе есть поле Token
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,   // нельзя читать из JS
		Secure:   false,  // true в проде по HTTPS
		SameSite: "None", // обязательно для кросс-доменных запросов
	}

	c.Cookie(cookie)

	return c.JSON(resp)
}

func (h *Handler) EnrichProfile(c *fiber.Ctx) error {
	user, ok := c.Locals(consts.UserContextKey).(*models.ClaimsJwt)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	request := models.EnrichProfileRequest{
		UserID:     user.UserID,
		FirstName:  c.FormValue("first_name"),
		Secondname: c.FormValue("second_name"),
		Patronymic: c.FormValue("patronymic"),
	}

	file, err := c.FormFile("image")
	if err == nil {
		f, err := file.Open()
		if err != nil {
			return err
		}
		defer f.Close()

		fileBytes, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		request.Image = &models.File{
			Filename: file.Filename,
			Bytes:    fileBytes,
		}
	} else {
		if !errors.Is(err, fasthttp.ErrMissingFile) {
			return err
		}
	}

	err = h.service.EnrichUserProfile(c.Context(), request)
	if err != nil {
		return err
	}

	return err
}

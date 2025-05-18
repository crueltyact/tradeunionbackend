package binder

import (
	"fmt"
	"profkom/internal/models"
	"profkom/internal/service/auth"
	"profkom/pkg/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

const (
	claimsKey    = "user"
	UserID       = "userID"
	TradeUnionID = "trade-union-id"
)

type Middleware struct {
	service *auth.Service
	secret  string
}

func New(secret string, service *auth.Service) *Middleware {
	return &Middleware{
		secret:  secret,
		service: service,
	}
}

func (m *Middleware) Auth(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")

	claims, err := m.parseJwt(token)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	exist, err := m.service.CheckUserInfoExist(ctx.Context(), claims.UserID)
	if err != nil {
		return err
	}

	if !exist {
		return ctx.Status(fiber.StatusForbidden).SendString("enrich profile")
	}

	ctx.Locals(claimsKey, claims)

	return ctx.Next()
}

func (m *Middleware) parseJwt(jwtToken string) (*models.ClaimsJwt, error) {
	claims := &models.ClaimsJwt{}

	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return []byte(m.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (m *Middleware) CheckTradeUnionID(ctx *fiber.Ctx) error {
	tradeUnionID := ctx.Query("tradeUnionID")

	ctx.Locals(consts.TradeUnionIDKey, tradeUnionID)

	return ctx.Next()
}

func (m *Middleware) AuthWebsocket(ctx *fiber.Ctx) error {
	tradeUnionID := ctx.Query("jwtToken")

	claims, err := m.parseJwt(tradeUnionID)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	ctx.Locals(claimsKey, claims)

	return ctx.Next()
}

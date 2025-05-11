package models

import (
	"profkom/internal/entities"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type (
	ClaimsJwt struct {
		UserID uuid.UUID
		Login  string
		Role   string
		jwt.StandardClaims
	}
	PostInviteTokenRequest struct {
		Role string `json:"role"`
	}
	PostInviteTokenResponse struct {
		IviteToken string `json:"invite_token"`
	}
	AdminSignInRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	AdminSignInResponse struct {
		Token string `json:"token"`
	}
	SignUpRequest struct {
		Login       string `json:"login"`
		Password    string `json:"password"`
		InviteToken string `json:"invite_token"`
	}
	SignUpResponse struct {
		JwtToken string `json:"token"`
		NextStep bool   `json:"next_step"`
	}
	EnrichProfileRequest struct {
		UserID     uuid.UUID `json:"id,omitempty"`
		FirstName  string    `json:"first_name"`
		Secondname string    `json:"second_name"`
		Patronymic string    `json:"patronymic,omitempty"`
		Image      *File     `json:"file,omitempty"`
	}
)

func (r *SignUpRequest) ToEntity() entities.User {
	return entities.User{
		Login:    r.Login,
		Password: r.Password,
	}
}

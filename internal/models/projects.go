package models

import "github.com/google/uuid"

type (
	GetProjectsResponse struct {
		Worker  Projects `json:"worker"`
		Student Projects `json:"student"`
	}
	Projects []Project
	Project  struct {
		ID          uuid.UUID `json:"id" db:"id"`
		Type        string    `json:"type" db:"type"`
		Intro       string    `json:"intro" db:"intro"`
		Title       string    `json:"title" db:"title"`
		Description string    `json:"description" db:"description"`
		ImageUrl    string    `json:"image_url" db:"image_url"`
		CreatedAt   int       `json:"created_at,omitempty" db:"created_at"`
	}
	PostProjectRequest struct {
		Title       string `json:"title"`
		Intro       string `json:"intro"`
		Description string `json:"description"`
		Type        string `json:"type"`
		Image       File   `json:"image"`
	}
)

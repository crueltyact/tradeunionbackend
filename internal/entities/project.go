package entities

import (
	"time"

	"github.com/google/uuid"
)

type (
	Project struct {
		ID          uuid.UUID `db:"id"`
		Intro       string    `db:"intro"`
		Title       string    `db:"title"`
		Description string    `db:"description"`
		Type        string    `db:"type"`
		ImageUrl    string    `db:"image_url"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
	}
)

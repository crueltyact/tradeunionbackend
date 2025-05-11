package entities

import (
	"time"

	"github.com/google/uuid"
)

type (
	New struct {
		ID        uuid.UUID `db:"id"`
		Title     string    `db:"title"`
		Content   string    `db:"content"`
		ImageUrl  string    `db:"image_url"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
)

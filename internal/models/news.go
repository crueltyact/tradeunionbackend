package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	News []New
	New  struct {
		ID        uuid.UUID `json:"id" db:"id"`
		Title     string    `json:"title" db:"title"`
		Content   string    `json:"content" db:"content"`
		ImageURL  string    `json:"image_url,omitempty" db:"image_url"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
	}
	PostNewRequest struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Image   File   `json:"image"`
	}
	File struct {
		Filename string `json:"filename"`
		Bytes    []byte `json:"content"`
	}
)

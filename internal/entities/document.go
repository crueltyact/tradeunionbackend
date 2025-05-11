package entities

import "github.com/google/uuid"

type (
	Document struct {
		ID    uuid.UUID `db:"id"`
		URL   string    `db:"url"`
		Title string    `db:"title"`
		Type  string    `db:"type"`
	}
	DocumentBatch struct {
		ID    []uuid.UUID `db:"id"`
		URL   []string    `db:"url"`
		Title []string    `db:"title"`
		Type  string      `db:"type"`
	}
)

package entities

import (
	"time"

	"github.com/google/uuid"
)

type (
	Chat struct {
		ID    uuid.UUID `db:"id"`
		Title string    `db:"title"`
	}
	ChatUserBatch struct {
		UserID []uuid.UUID `db:"user_id"`
		ChatID uuid.UUID   `db:"chat_id"`
	}
	ChatUser struct {
		UserID uuid.UUID `db:"user_id"`
		ChatID uuid.UUID `db:"chat_id"`
	}
	Message struct {
		ID        uuid.UUID `db:"id"`
		ChatID    uuid.UUID `db:"chat_id"`
		Content   string    `db:"content"`
		UserID    uuid.UUID `db:"user_id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
		Role      string    `db:"role"`
	}
)

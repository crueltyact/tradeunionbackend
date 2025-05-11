package entities

import (
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		ID       uuid.UUID `db:"id"`
		Role     string    `db:"role"`
		Login    string    `db:"login"`
		Password string    `db:"password"`
		CreateAt time.Time `db:"created_at"`
	}
)

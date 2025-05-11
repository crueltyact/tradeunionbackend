package entities

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v2"
)

type (
	User struct {
		ID       uuid.UUID   `db:"id"`
		Role     string      `db:"role"`
		Name     null.String `db:"name"`
		Login    string      `db:"login"`
		Password string      `db:"password"`
		CreateAt time.Time   `db:"created_at"`
	}
)

package entities

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v2"
)

type (
	InviteToken struct {
		ID      int       `db:"id"`
		Content uuid.UUID `db:"content"`
		Used    bool      `db:"used"`
		Role    string    `db:"role"`
	}
	UserInfo struct {
		UserID     uuid.UUID   `db:"user_id"`
		FirstName  string      `db:"first_name"`
		SecondName string      `db:"second_name"`
		Patronymic string      `db:"patronymic"`
		ImageUrl   null.String `db:"iamge_url"`
	}
)

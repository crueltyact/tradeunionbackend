package auth

import (
	"context"
	"profkom/internal/entities"
)

func (r *Repository) InsertUserInfo(ctx context.Context, userInfo entities.UserInfo) (err error) {
	query := `
		insert into auth.user_info(
			user_id,
			first_name,
			second_name,
			patronymic,
			image_url
		) values (
			$1,
			$2,
			$3,
			$4,
			$5
		)
	`

	_, err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).ExecContext(
		ctx,
		query,
		userInfo.UserID,
		userInfo.FirstName,
		userInfo.SecondName,
		userInfo.Patronymic,
		userInfo.ImageUrl,
	)
	if err != nil {
		return err
	}

	return err
}

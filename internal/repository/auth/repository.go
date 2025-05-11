package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"profkom/internal/entities"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db        *sqlx.DB
	ctxGetter *trmsqlx.CtxGetter
}

func New(db *sqlx.DB, ctxGetter *trmsqlx.CtxGetter) *Repository {
	return &Repository{
		db:        db,
		ctxGetter: ctxGetter,
	}
}

func (r *Repository) InsertUser(ctx context.Context, user *entities.User) (err error) {
	query := `
		INSERT INTO auth."user" (
			role,
			login,
			password
		) VALUES (
			$1,
			$2,
			$3
		) RETURNING *
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		user,
		query,
		user.Role,
		user.Login,
		user.Password,
	)
	if err != nil {
		return err
	}

	return err
}

func (r *Repository) CheckUserExist(ctx context.Context, login string) (exist bool, err error) {
	query := `
		select exists(
				select 1
				from auth."user"
				where 
					login = $1 
			) as result
	`
	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		&exist,
		query,
		login,
	)
	if err != nil {
		return exist, err
	}

	return exist, err
}

func (r *Repository) CheckInviteToken(ctx context.Context, inviteToken string) (role string, err error) {
	query := `
		SELECT
		 	role
		FROM auth.invite_token
		WHERE
			content = $1 AND used = false
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		&role,
		query,
		inviteToken,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return role, fmt.Errorf("invalid user token")
		}
		return role, err
	}

	query = `
		UPDATE auth.invite_token SET used = true
	`

	_, err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).ExecContext(
		ctx,
		query,
	)
	if err != nil {
		return role, err
	}

	return role, err
}

func (r *Repository) InsertInviteToken(ctx context.Context, token *entities.InviteToken) (err error) {
	query := `
		INSERT INTO auth.invite_token(
			content,
			role
		) VALUES (
			$1, 
			$2 
		) RETURNING *
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		token,
		query,
		token.Content,
		token.Role,
	)
	if err != nil {
		return err
	}

	return err
}

func (r *Repository) SelectUserByLogin(ctx context.Context, login string) (user entities.User, err error) {
	query := `
		SELECT
		*
		FROM auth."user"
		WHERE login = $1;
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		&user,
		query,
		login,
	)
	if err != nil {
		return user, err
	}

	return user, err
}


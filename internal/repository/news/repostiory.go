package news

import (
	"context"
	"profkom/internal/entities"
	"profkom/internal/models"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/google/uuid"
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

func (r Repository) SelectNew(ctx context.Context, id uuid.UUID) (new models.New, err error) {
	query := `
		SELECT
			id,
			title,
			content,
			image_url
		FROM content.news
		WHERE id = $1
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		&new,
		query,
		id,
	)
	if err != nil {
		return new, err
	}

	return new, err
}

func (r *Repository) SelectNews(ctx context.Context) (news models.News, err error) {
	query := `
		SELECT
			id,
			title,
			content,
			image_url,
			created_at
		FROM content.news
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).SelectContext(
		ctx,
		&news,
		query,
	)
	if err != nil {
		return news, err
	}

	return news, err
}

func (r *Repository) InsertNew(ctx context.Context, new *entities.New) (err error) {
	query := `
		INSERT INTO content.news(
			id,
			title,
			content,
			image_url
		) VALUES (
			$1,
			$2,
			$3,
			$4
		) RETURNING *
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		new,
		query,
		new.ID,
		new.Title,
		new.Content,
		new.ImageUrl,
	)
	if err != nil {
		return err
	}

	return err
}

package documents

import (
	"context"
	"profkom/internal/entities"

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

func (r *Repository) InsertDocumnets(ctx context.Context, documents entities.DocumentBatch) (err error) {
	query := `
		INSERT INTO content.documents(
			id,
			url,
			title,
			type
		) VALUES (
			unnest($1::UUID[]),
			unnest($2::text[]),
			unnest($3::text[]),
			$4
		)
	`

	_, err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).ExecContext(
		ctx,
		query,
		documents.ID,
		documents.URL,
		documents.Title,
		documents.Type,
	)

	if err != nil {
		return err
	}

	return err
}

func (r *Repository) SelectDocuments(ctx context.Context) (documents []entities.Document, err error) {
	query := `
		SELECT 
			*
		FROM 
			content.documents
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).SelectContext(
		ctx,
		&documents,
		query,
	)
	if err != nil {
		return documents, err
	}

	return documents, err
}

func (r *Repository) DeleteDocument(ctx context.Context, id uuid.UUID) (err error) {
	query := `
		DELETE FROM content.documents WHERE id = $1
	`
	_, err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).ExecContext(
		ctx,
		query,
		id,
	)

	if err != nil {
		return err
	}

	return err
}

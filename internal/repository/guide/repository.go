package guide

import (
	"context"
	"profkom/internal/models"

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

func (r *Repository) SelectGuide(ctx context.Context) (result models.AllGuides, err error) {
	query := `
		SELECT
			id,
			title
		FROM
			guides.guides
		WHERE 
			type = 'worker' 
	`
	var guides []models.Guide

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).SelectContext(
		ctx,
		&guides,
		query,
	)
	if err != nil {
		return result, err
	}

	query = `
		SELECT
			id,
			title,
			content
		FROM
			guides.themes
		WHERE 
			guide_id = $1 
	`

	for idx, guide := range guides {
		var subGuides []models.SubGuides

		err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).SelectContext(
			ctx,
			&subGuides,
			query,
			guide.ID,
		)
		if err != nil {
			return result, err
		}

		guides[idx].SubGuides = subGuides
	}

	result.Worker = guides

	query = `
		SELECT
			id,
			title
		FROM
			guides.guides
		WHERE 
			type = 'student' 
	`
	var studentGuides []models.Guide

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).SelectContext(
		ctx,
		&studentGuides,
		query,
	)
	if err != nil {
		return result, err
	}

	query = `
		SELECT
			id,
			title,
			content
		FROM
			guides.themes
		WHERE 
			guide_id = $1 
	`

	for idx, guide := range studentGuides {
		var subGuides []models.SubGuides

		err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).SelectContext(
			ctx,
			&subGuides,
			query,
			guide.ID,
		)
		if err != nil {
			return result, err
		}

		studentGuides[idx].SubGuides = subGuides
	}

	result.Student = studentGuides

	return result, err
}

func (r *Repository) InsertGuide(ctx context.Context, guideType string, guides []models.Guide) (err error) {
	query := `
		INSERT INTO guides.guides  (
			title,
			type
		) VALUES(
			$1,
			$2
		) RETURNING ID
	`

	themesQuery := `
		INSERT INTO guides.themes  (
			title,
			content,
			guide_id
		) VALUES(
			$1,
			$2,
			$3
		)
	`

	for _, guide := range guides {
		var id int

		err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
			ctx,
			&id,
			query,
			guide.Label,
			guideType,
		)
		if err != nil {
			return err
		}

		for _, subGuide := range guide.SubGuides {
			_, err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).ExecContext(
				ctx,
				themesQuery,
				subGuide.Label,
				subGuide.Content,
				id,
			)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// func (r *Repository) UpdateGuide(ctx context.Context) (err error) {
// 	query := `
// 	`

// 	return err
// }

func (r *Repository) DeleteGuide(ctx context.Context, id int) (err error) {
	query := `
		DELETE FROM guides.guides WHERE id = $1
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

func (r *Repository) DeleteTheme(ctx context.Context, id int) (err error) {
	query := `
		DELETE FROM guides.themes WHERE id = $1
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

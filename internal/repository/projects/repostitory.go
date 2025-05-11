package projects

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

func (r *Repository) SelectProjects(ctx context.Context) (result models.GetProjectsResponse, err error) {
	query := `
		SELECT
			id,
			type,
			intro,
			title,
			description,
			image_url
		FROM
			content.projects
	`

	var projects models.Projects

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).SelectContext(
		ctx,
		&projects,
		query,
	)
	if err != nil {
		return result, err
	}

	for _, project := range projects {
		switch project.Type {
		case "worker":
			result.Worker = append(result.Worker, project)
		case "student":
			result.Student = append(result.Student, project)
		}
	}

	return result, err
}

func (r *Repository) SelectProject(ctx context.Context, id uuid.UUID) (result models.Project, err error) {
	query := `
		SELECT
			id,
			type,
			intro,
			title,
			description,
			image_url
		FROM
			content.projects
		WHERE id = $1
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		&result,
		query,
		id,
	)
	if err != nil {
		return result, err
	}

	return result, err
}

func (r *Repository) InsertProject(ctx context.Context, project entities.Project) (err error) {
	query := `
		INSERT INTO content.projects (
			id,
			intro,
			title,
			description,
			type,
			image_url
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6 
		)
	`

	_, err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).ExecContext(
		ctx,
		query,
		project.ID,
		project.Intro,
		project.Title,
		project.Description,
		project.Type,
		project.ImageUrl,
	)
	if err != nil {
		return err
	}

	return err
}

func (r *Repository) DeleteProject(ctx context.Context, id uuid.UUID) (err error) {
	query := `
		DELETE FROM content.projects WHERE id = $1
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

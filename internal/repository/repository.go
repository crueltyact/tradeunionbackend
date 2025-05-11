package repository

import (
	"profkom/internal/repository/auth"
	"profkom/internal/repository/chat"
	"profkom/internal/repository/documents"
	"profkom/internal/repository/guide"
	"profkom/internal/repository/news"
	"profkom/internal/repository/projects"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Guide     *guide.Repository
	Project   *projects.Repository
	News      *news.Repository
	Auth      *auth.Repository
	Documents *documents.Repository
	Chat      *chat.Repository
}

func New(db *sqlx.DB, ctxGetter *trmsqlx.CtxGetter) *Repository {
	return &Repository{
		Guide:     guide.New(db, ctxGetter),
		Project:   projects.New(db, ctxGetter),
		News:      news.New(db, ctxGetter),
		Auth:      auth.New(db, ctxGetter),
		Documents: documents.New(db, ctxGetter),
		Chat:      chat.New(db, ctxGetter),
	}
}

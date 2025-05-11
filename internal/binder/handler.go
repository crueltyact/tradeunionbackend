package binder

import (
	"profkom/internal/binder/auth"
	"profkom/internal/binder/chat"
	"profkom/internal/binder/documents"
	"profkom/internal/binder/guide"
	"profkom/internal/binder/news"
	"profkom/internal/binder/projects"
	"profkom/internal/service"
)

type Handler struct {
	Guide     *guide.Handler
	Project   *projects.Handler
	News      *news.Handler
	Auth      *auth.Handler
	Documents *documents.Handler
	Chat      *chat.Handler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		Guide:     guide.New(service.Guide),
		Project:   projects.New(service.Project),
		News:      news.New(service.News),
		Auth:      auth.New(service.Auth),
		Documents: documents.New(service.Documents),
		Chat:      chat.New(service.Chat),
	}
}

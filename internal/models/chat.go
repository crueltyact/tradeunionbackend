package models

import "github.com/google/uuid"

type (
	PostChatRequest struct {
		Title string      `json:"title"`
		Users []uuid.UUID `json:"users"`
	}
	PostChatResponse struct {
		ChatID uuid.UUID `json:"chat_id"`
	}
	PostMessageRequest struct {
		Content string    `json:"content"`
		ChatID  uuid.UUID `json:"omitempty"`
		UserID  uuid.UUID `json:"omitempty"`
	}
	PostMessageResponse struct {
		ID        uuid.UUID `json:"id"`
		Content   string    `json:"content"`
		ChatID    string    `json:"chat_id,omitempty"`
		UserID    uuid.UUID `json:"user_id,omitempty"`
		CreatedAt int64     `json:"created_at"`
		UpdatedAt int64     `json:"updated_at"`
	}
	CheckAccessToChat struct {
		UserID uuid.UUID
		ChatID uuid.UUID
	}
	Chat struct {
		ID       uuid.UUID             `json:"chat_id"`
		Title    string                `json:"title"`
		Messages []PostMessageResponse `json:"messages"`
	}
	GetChatsRequest struct {
		UserID uuid.UUID
	}
	GetChatsResponse struct {
		Chats []Chat `json:"chats"`
	}
	PostClientChatRequest struct {
		TradeUnionID string
	}
	PostClientChatResponse struct {
		ChatID   string                `json:"chat_id"`
		Worker   Worker                `json:"worker"`
		Messages []PostMessageResponse `json:"messages"`
	}
	Worker struct {
		ImageURL   string `json:"iamge_url"`
		FirstName  string `json:"first_name"`
		SecondName string `json:"second_name"`
		Patronymic string `json:"patronymic"`
	}
)

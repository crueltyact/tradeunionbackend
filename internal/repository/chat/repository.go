package chat

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

func (r *Repository) InsertChat(ctx context.Context, chat *entities.Chat) (err error) {
	query := `
		insert into chat.chat(
			id,
			title
		) values(
			$1,
			$2 
		) RETURNING *
	`
	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		chat,
		query,
		chat.ID,
		chat.Title,
	)
	if err != nil {
		return err
	}

	return err
}

func (r *Repository) InsertChatUser(ctx context.Context, users entities.ChatUserBatch) (err error) {
	query := `
		insert into chat.chat_users (
			chat_id,
			user_id
		) values (
			$1,
			unnest($2::UUID[])
		)
	`

	_, err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).ExecContext(
		ctx,
		query,
		users.ChatID,
		users.UserID,
	)
	if err != nil {
		return err
	}

	return err
}

func (r Repository) InsertMessage(ctx context.Context, message *entities.Message) (err error) {
	query := `
		insert into chat.messages (
			id,
			content,
			user_id,
			chat_id
		) values(
			$1,
			$2,
			$3,
			$4	 
		) returning *
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		message,
		query,
		message.ID,
		message.Content,
		message.UserID,
		message.ChatID,
	)
	if err != nil {
		return err
	}

	return err
}

func (r *Repository) DeleteMessage(ctx context.Context, messageID uuid.UUID) (err error) {
	return err
}

func (r *Repository) UpdateMessage(ctx context.Context, message *entities.Message) (err error) {
	// query := `
	// 	update chat.messages
	// `
	return err
}

func (r *Repository) SelectExistChatUser(ctx context.Context, user entities.ChatUser) (exist bool, err error) {
	query := `
		select exists(
				select 1
				from chat.chat_users
				where 
					user_id = $1 and chat_id = $2
			) as result
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		&exist,
		query,
		user.UserID,
		user.ChatID,
	)
	if err != nil {
		return exist, err
	}

	return exist, err
}

func (r *Repository) SelectChats(ctx context.Context, userID uuid.UUID) (chats []entities.Chat, err error) {
	query := `
		select
			c.id,
			c.title
		from chat.chat c
		join chat.chat_users cu on cu.chat_id = c.id
		where cu.user_id = $1 
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).SelectContext(
		ctx,
		&chats,
		query,
		userID,
	)
	if err != nil {
		return chats, err
	}

	return chats, err
}

func (r *Repository) SelectMessages(ctx context.Context, userID, chatID uuid.UUID) (messages []entities.Message, err error) {
	query := `
		select
			*
		from chat.messages
		where chat_id = $1 and user_id = $2
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).SelectContext(
		ctx,
		&messages,
		query,
		chatID,
		userID,
	)
	if err != nil {
		return messages, err
	}

	return messages, err
}

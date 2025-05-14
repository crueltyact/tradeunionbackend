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

func (r *Repository) SelectMessages(ctx context.Context, chatID uuid.UUID) (messages []entities.Message, err error) {
	query := `
		select
			*
		from chat.messages
		where chat_id = $1
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).SelectContext(
		ctx,
		&messages,
		query,
		chatID,
	)
	if err != nil {
		return messages, err
	}

	return messages, err
}

func (r *Repository) SelectClientExist(ctx context.Context, unionTradeID string) (exist bool, err error) {
	query := `
		select exists(
				select 1
				from auth.client
				where 
					trade_union_id = $1 
			) as result
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		&exist,
		query,
		unionTradeID,
	)
	if err != nil {
		return exist, err
	}

	return exist, err
}

func (r *Repository) InsertClient(ctx context.Context, client *entities.Client) (err error) {
	query := `
		insert into auth.client(
			id,
			trade_union_id
		) values(
			$1,
			$2 
		)
	`

	_, err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).ExecContext(
		ctx,
		query,
		client.UserID,
		client.TradeUnionID,
	)
	if err != nil {
		return err
	}

	return err
}

func (r *Repository) SelectClient(ctx context.Context, client *entities.Client) (err error) {
	query := `
		select 
			id,
			trade_union_id
		from auth.client
		where 
			trade_union_id = $1 
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		client,
		query,
		client.TradeUnionID,
	)
	if err != nil {
		return err
	}

	return err
}

func (r *Repository) SelectTradeUnionWorkerIDForHelp(ctx context.Context) (worker entities.User, err error) {
	query := `
		SELECT * FROM auth."user"
		where role = 'worker'
		ORDER BY random()
		LIMIT 1
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		&worker,
		query,
	)
	if err != nil {
		return worker, err
	}

	return worker, err
}

func (r *Repository) SelectChatExist(ctx context.Context, clientID uuid.UUID) (exist bool, err error) {
	query := `
		select exists(
			select 1
			from chat.chat_users
			where user_id = $1
		) as result
	`

	err = r.ctxGetter.DefaultTrOrDB(
		ctx,
		r.db,
	).GetContext(
		ctx,
		&exist,
		query,
		clientID,
	)
	if err != nil {
		return exist, err
	}

	return exist, err
}

func (r *Repository) SelectChat(ctx context.Context, userID uuid.UUID) (chat entities.Chat, err error) {
	query := `
		select 
			chat.chat.id, 
			chat.chat.title
		from chat.chat_users
		join chat.chat on chat.chat.id = chat.chat_users.chat_id
		where chat.chat_users.user_id = $1
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		&chat,
		query,
		userID,
	)
	if err != nil {
		return chat, err
	}

	return chat, err
}

func (r *Repository) SelectWorkerInfo(ctx context.Context, userID uuid.UUID) (userInfo entities.UserInfo, err error) {
	query := `
		select 
			*
		from auth.user_info
		where user_id = $1
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		&userInfo,
		query,
		userID,
	)
	if err != nil {
		return userInfo, err
	}

	return userInfo, err
}

func (r *Repository) SelectUserIDByTradeUnionID(ctx context.Context, tradeUnionID string) (id uuid.UUID, err error) {
	query := `
		select
			id
		from auth.client
		where trade_union_id = $1
	`

	err = r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		&id,
		query,
		tradeUnionID,
	)
	if err != nil {
		return id, err
	}

	return id, err
}

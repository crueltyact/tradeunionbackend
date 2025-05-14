package chat

import (
	"context"
	"fmt"
	"profkom/internal/entities"
	"profkom/internal/models"
	"profkom/internal/repository/chat"
	"profkom/pkg/s3"

	txmanager "github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type Service struct {
	repo      *chat.Repository
	txManager *txmanager.Manager
	s3        *s3.Client
}

func New(repo *chat.Repository, txManager *txmanager.Manager, s3 *s3.Client) *Service {
	return &Service{
		repo:      repo,
		txManager: txManager,
		s3:        s3,
	}
}

func (s *Service) CreateChat(ctx context.Context, req models.PostChatRequest) (err error) {
	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		id := uuid.New()

		chat := &entities.Chat{
			ID:    id,
			Title: req.Title,
		}

		err = s.repo.InsertChat(ctx, chat)
		if err != nil {
			return err
		}

		chatUsers := entities.ChatUserBatch{
			ChatID: id,
			UserID: req.Users,
		}

		err = s.repo.InsertChatUser(ctx, chatUsers)
		if err != nil {
			return err
		}

		return err
	})

	return err
}

func (s *Service) SendMessage(ctx context.Context, req models.PostMessageRequest) (resp models.PostMessageResponse, err error) {
	var message entities.Message

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		id := uuid.New()

		message = entities.Message{
			ID:      id,
			ChatID:  req.ChatID,
			Content: req.Content,
			UserID:  req.UserID,
		}

		err = s.repo.InsertMessage(ctx, &message)
		if err != nil {
			return err
		}

		return err
	})

	return models.PostMessageResponse{
		ID:        message.ID,
		Content:   message.Content,
		ChatID:    message.ChatID.String(),
		UserID:    message.UserID,
		CreatedAt: message.CreatedAt.Unix(),
		UpdatedAt: message.UpdatedAt.Unix(),
	}, err
}

func (s *Service) ChechAccessToChat(ctx context.Context, req models.CheckAccessToChat) (exist bool, err error) {
	chatUser := entities.ChatUser{
		UserID: req.UserID,
		ChatID: req.ChatID,
	}

	exist, err = s.repo.SelectExistChatUser(ctx, chatUser)
	if err != nil {
		return exist, err
	}

	return exist, err
}

func (s *Service) GetChats(ctx context.Context, req models.GetChatsRequest) (resp models.GetChatsResponse, err error) {
	chats, err := s.repo.SelectChats(ctx, req.UserID)
	if err != nil {
		return resp, err
	}

	for _, chat := range chats {
		messages, err := s.repo.SelectMessages(ctx, chat.ID)
		if err != nil {
			return resp, err
		}

		var msgs []models.PostMessageResponse
		for _, msg := range messages {
			msgs = append(msgs, models.PostMessageResponse{
				ID:        msg.ID,
				Content:   msg.Content,
				ChatID:    msg.ChatID.String(),
				UserID:    msg.UserID,
				CreatedAt: msg.CreatedAt.Unix(),
				UpdatedAt: msg.UpdatedAt.Unix(),
			})
		}

		resp.Chats = append(resp.Chats, models.Chat{
			ID:       chat.ID,
			Title:    chat.Title,
			Messages: msgs,
		})
	}

	return resp, err
}

func (s *Service) CreateClientChat(ctx context.Context, req models.PostClientChatRequest) (resp models.PostClientChatResponse, err error) {
	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		client := entities.Client{
			TradeUnionID: req.TradeUnionID,
		}

		exist, err := s.repo.SelectClientExist(ctx, req.TradeUnionID)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("user %s - exist is %s", req.TradeUnionID, exist))

		if !exist {
			client.UserID = uuid.New()

			err := s.repo.InsertClient(ctx, &client)
			if err != nil {
				return err
			}
		} else {
			err = s.repo.SelectClient(ctx, &client)
			if err != nil {
				log.Error("fail get user", err)
				return err
			}
		}

		worker, err := s.repo.SelectTradeUnionWorkerIDForHelp(ctx)
		if err != nil {
			return err
		}

		workerInfo, err := s.repo.SelectWorkerInfo(ctx, worker.ID)
		if err != nil {
			return err
		}

		resp.Worker = models.Worker{
			FirstName:  workerInfo.FirstName,
			SecondName: workerInfo.SecondName,
			Patronymic: workerInfo.Patronymic,
			ImageURL:   workerInfo.ImageUrl.String,
		}

		chatExist, err := s.repo.SelectChatExist(ctx, client.UserID)
		if err != nil {
			return err
		}

		if chatExist {
			chat, err := s.repo.SelectChat(ctx, client.UserID)
			if err != nil {
				return err
			}

			messages, err := s.repo.SelectMessages(ctx, chat.ID)
			if err != nil {
				return err
			}

			var msgs []models.PostMessageResponse
			for _, msg := range messages {
				msgs = append(msgs, models.PostMessageResponse{
					ID:        msg.ID,
					Content:   msg.Content,
					ChatID:    msg.ChatID.String(),
					UserID:    msg.UserID,
					CreatedAt: msg.CreatedAt.Unix(),
					UpdatedAt: msg.UpdatedAt.Unix(),
				})
			}

			resp.ChatID = chat.ID.String()
			resp.Messages = msgs
		} else {
			chat := &entities.Chat{
				ID:    uuid.New(),
				Title: req.TradeUnionID,
			}

			err = s.repo.InsertChat(ctx, chat)
			if err != nil {
				return err
			}

			users := []uuid.UUID{worker.ID, client.UserID}

			chatUsers := entities.ChatUserBatch{
				ChatID: chat.ID,
				UserID: users,
			}

			err = s.repo.InsertChatUser(ctx, chatUsers)
			if err != nil {
				return err
			}

			resp.ChatID = chat.ID.String()
		}

		return err
	})
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (s *Service) GetUserIDByTradeUnionID(ctx context.Context, tradeUinonID string) (uuid.UUID, error) {
	return s.repo.SelectUserIDByTradeUnionID(ctx, tradeUinonID)
}

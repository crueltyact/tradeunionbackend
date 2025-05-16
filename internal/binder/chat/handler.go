package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"profkom/internal/models"
	"profkom/internal/service/chat"
	"profkom/pkg/consts"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

var (
	chats   map[uuid.UUID]map[uuid.UUID]*websocket.Conn = make(map[uuid.UUID]map[uuid.UUID]*websocket.Conn)
	rwmutex                                             = sync.RWMutex{}
)

type Handler struct {
	service *chat.Service
}

func New(service *chat.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleConnection(c *websocket.Conn) {
	var (
		code int
		err  error
	)

	ctx := context.Background()
	defer func() {
		if err != nil {
			c.WriteJSON(map[string]string{"error": err.Error()})
			c.CloseHandler()(code, err.Error())

			return
		}
	}()

	user, ok := c.Locals(consts.UserContextKey).(*models.ClaimsJwt)
	if !ok {
		code = fiber.StatusUnauthorized
		err = fmt.Errorf("Unathorized")

		return
	}

	chatID := c.Params("chat_id")
	chatUUID, err := uuid.Parse(chatID)
	if err != nil {
		code = fiber.StatusNotFound

		return
	}

	hasAccess, err := h.service.ChechAccessToChat(ctx, models.CheckAccessToChat{
		UserID: user.UserID,
		ChatID: chatUUID,
	})
	if err != nil {
		code = fiber.StatusInternalServerError

		return
	}

	if !hasAccess {
		err = fmt.Errorf("no access")
		code = fiber.StatusForbidden

		return
	}

	rwmutex.Lock()
	conns, ok := chats[chatUUID]
	if !ok {
		conns = map[uuid.UUID]*websocket.Conn{}

		chats[chatUUID] = conns
	}

	conns[user.UserID] = c
	rwmutex.Unlock()

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			break
		}

		var request models.PostMessageRequest
		err = json.Unmarshal(msg, &request)
		if err != nil {
			break
		}

		request.ChatID = chatUUID
		request.UserID = user.UserID

		resp, err := h.service.SendMessage(ctx, request)
		if err != nil {
			code = fiber.StatusInternalServerError

			return
		}

		message, err := json.Marshal(resp)
		if err != nil {
			code = fiber.StatusInternalServerError

			return
		}

		for _, conn := range conns {
			if err := conn.WriteMessage(mt, message); err != nil {
				break
			}
		}
	}
}

func (h *Handler) GetChats(c *fiber.Ctx) error {
	user, ok := c.Locals(consts.UserContextKey).(*models.ClaimsJwt)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	request := models.GetChatsRequest{
		UserID: user.UserID,
	}

	resp, err := h.service.GetChats(c.Context(), request)
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

func (h *Handler) PostClientChat(c *fiber.Ctx) error {
	tradeUnionID := c.Get(consts.TradeUnionIDKey)
	if tradeUnionID == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	log.Info(tradeUnionID)

	request := models.PostClientChatRequest{
		TradeUnionID: tradeUnionID,
	}

	resp, err := h.service.CreateClientChat(c.Context(), request)
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

func (h *Handler) HandleClientConnection(c *websocket.Conn) {
	var (
		code int
		err  error
	)

	ctx := context.Background()
	defer func() {
		if err != nil {
			c.WriteJSON(map[string]string{"error": err.Error()})
			c.CloseHandler()(code, err.Error())

			return
		}
	}()

	tradeUnionID := c.Locals(consts.TradeUnionIDKey).(string)
	if tradeUnionID == "" {
		log.Info(tradeUnionID)
		code = fiber.StatusUnauthorized
		err = fmt.Errorf("Unathorized")

		return
	}

	chatID := c.Params("chat_id")
	chatUUID, err := uuid.Parse(chatID)
	if err != nil {
		code = fiber.StatusNotFound

		return
	}

	userID, err := h.service.GetUserIDByTradeUnionID(ctx, tradeUnionID)
	if err != nil {
		code = fiber.StatusUnauthorized

		return
	}

	rwmutex.Lock()
	conns, ok := chats[chatUUID]
	if !ok {
		conns = map[uuid.UUID]*websocket.Conn{}

		chats[chatUUID] = conns
	}

	conns[userID] = c
	rwmutex.Unlock()

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			break
		}

		var request models.PostMessageRequest
		err = json.Unmarshal(msg, &request)
		if err != nil {
			break
		}

		request.ChatID = chatUUID
		request.UserID = userID

		resp, err := h.service.SendMessage(ctx, request)
		if err != nil {
			code = fiber.StatusInternalServerError

			return
		}

		message, err := json.Marshal(resp)
		if err != nil {
			code = fiber.StatusInternalServerError

			return
		}

		for _, conn := range conns {
			if err := conn.WriteMessage(mt, message); err != nil {
				break
			}
		}
	}
}

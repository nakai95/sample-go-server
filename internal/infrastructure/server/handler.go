package server

import (
	"encoding/json"
	"net/http"
	"sample-go-server/api"
	"sample-go-server/internal/adapter/controller"
	"sample-go-server/internal/adapter/presenter"
	"sample-go-server/internal/adapter/repository"
	"sample-go-server/internal/infrastructure/datastore"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type handler struct {
	sync.RWMutex
	event controller.EventController
	chat  controller.ChatController
}

func NewHandler() (*handler, error) {

	manager, err := datastore.NewPostgresManager()
	if err != nil {
		return nil, err
	}

	eventRepo := repository.NewEventRepository(manager.EventDataStore)
	eventPres := presenter.NewEventPresenter()
	eventCtrl := controller.NewEventController(eventRepo, eventPres)

	chatRepo := repository.NewChatRepository(manager.ChatDataStore)
	chatPres := presenter.NewChatPresenter()
	chatCtrl := controller.NewChatController(chatRepo, chatPres)

	return &handler{
		event: eventCtrl,
		chat:  chatCtrl,
	}, nil
}

func sendServerError(ctx echo.Context, code int, message string) error {
	errResponse := api.Error{
		Code:    int32(code),
		Message: message,
	}
	return ctx.JSON(code, errResponse)
}

// Ensure that we implement the server interface
var _ api.ServerInterface = (*handler)(nil)

func (h *handler) HealthCheck(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func (h *handler) GetToken(ctx echo.Context) error {
	var tokenRequest api.GetTokenFormdataRequestBody
	err := ctx.Bind(&tokenRequest)
	if err != nil {
		return sendServerError(ctx, http.StatusBadRequest, "could not bind request body")
	}

	// Check the username and password
	if (tokenRequest.Username != "demo1@example.com" || tokenRequest.Password != "#demo1") &&
		(tokenRequest.Username != "demo2@example.com" && tokenRequest.Username != "#demo2") {
		return sendServerError(ctx, http.StatusUnauthorized, "invalid username or password")
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": tokenRequest.Username,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	// Sign the token
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return sendServerError(ctx, http.StatusInternalServerError, "could not sign token")
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"token": tokenString,
	})
}

func (h *handler) ListEvents(ctx echo.Context) error {
	h.RLock()
	defer h.RUnlock()

	events, err := h.event.ListEvents()
	if err != nil {
		return sendServerError(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, events)
}

func (h *handler) ListChatRooms(ctx echo.Context) error {
	h.RLock()
	defer h.RUnlock()

	chatRooms, err := h.chat.ListChatRooms()
	if err != nil {
		return sendServerError(ctx, http.StatusInternalServerError, "could not list chat rooms")
	}

	return ctx.JSON(http.StatusOK, chatRooms)
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
)

func handleMessages(client *Client) {
	for message := range client.Send {
		client.Conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (h *handler) ChatWebSocket(ctx echo.Context, id string) error {
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return sendServerError(ctx, http.StatusInternalServerError, "could not upgrade connection")
	}
	defer ws.Close()

	client := &Client{
		ID:   id,
		Conn: ws,
		Send: make(chan []byte),
	}

	room := getOrCreateChatRoom(id)
	room.Register <- client
	defer func() { room.Unregister <- client }()

	go handleMessages(client)

	for {
		var msg api.ChatMessage

		// Read message from client
		_, readMsg, err := ws.ReadMessage()
		if err != nil {
			return sendServerError(ctx, http.StatusInternalServerError, "could not read message")
		}

		// Unmarshal the message into the WebSocketMessage struct
		err = json.Unmarshal(readMsg, &msg)
		if err != nil {
			return sendServerError(ctx, http.StatusInternalServerError, "could not unmarshal message")
		}
		if err := h.chat.SaveMessage(msg); err != nil {
			return sendServerError(ctx, http.StatusInternalServerError, "could not save message")
		}

		// Marshal the message back to JSON
		writeMsg, err := json.Marshal(msg)
		if err != nil {
			return sendServerError(ctx, http.StatusInternalServerError, "could not marshal message")
		}

		// Broadcast the message to all clients
		room.Broadcast <- writeMsg
	}
}

func (h *handler) ListChatMessages(ctx echo.Context, roomId string) error {
	h.RLock()
	defer h.RUnlock()

	messages, err := h.chat.GetMessages(roomId)
	if err != nil {
		return sendServerError(ctx, http.StatusInternalServerError, "could not get messages")
	}

	return ctx.JSON(http.StatusOK, messages)
}

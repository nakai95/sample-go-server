package server

import (
	"net/http"
	"sample-go-server/api"
	"sample-go-server/internal/adapter/controller"
	"sample-go-server/internal/adapter/presenter"
	"sample-go-server/internal/adapter/repository"
	"sample-go-server/internal/infrastructure/datastore"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type handler struct {
	sync.RWMutex
	event controller.EventController
}

func NewHandler() *handler {

	ds := datastore.NewFakeDataStore()

	eventRepo := repository.NewEventRepository(ds)
	eventPres := presenter.NewEventPresenter()
	eventCtrl := controller.NewEventController(eventRepo, eventPres)

	return &handler{
		event: eventCtrl,
	}
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
	if tokenRequest.Username != "demo@example.com" || tokenRequest.Password != "#demo" {
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
		return sendServerError(ctx, http.StatusInternalServerError, "could not list events")
	}

	return ctx.JSON(http.StatusOK, events)
}

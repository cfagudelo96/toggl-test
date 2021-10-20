// Package handler contain the different handlers for application inputs.
package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/cfagudelo96/toggle-test/deck/domain"
	"github.com/cfagudelo96/toggle-test/deck/service"
	"github.com/labstack/echo/v4"
)

const (
	uuidParam          = "uuid"
	shuffledQueryParam = "shuffled"
	cardsQueryParam    = "cards"
)

// DeckService represents the interface required to handle the decks use cases.
type DeckService interface {
	CreateDeck(ctx context.Context, opts ...service.DeckCreationOption) (service.CreateDeckOutput, error)
	OpenDeck(ctx context.Context, uuid string) (service.OpenDeckOutput, error)
	DrawCards(ctx context.Context, uuid string, amount int) (service.DrawCardsOutput, error)
}

// DeckEchoHandler handles the echo HTTP requests.
type DeckEchoHandler struct {
	deckService DeckService
}

// NewDeckEchoHandler returns a new deck handler for handling echo HTTP requests.
func NewDeckEchoHandler(s DeckService) *DeckEchoHandler {
	return &DeckEchoHandler{
		deckService: s,
	}
}

// HandleCreateDeck handles the endpoint to create a new deck.
func (h *DeckEchoHandler) HandleCreateDeck(c echo.Context) error {
	var opts []service.DeckCreationOption

	if c.QueryParam(shuffledQueryParam) == "n" {
		opts = append(opts, service.Shuffled(false))
	}

	if cardsStr := c.QueryParam(cardsQueryParam); cardsStr != "" {
		cardsSplit := strings.Split(cardsStr, ",")
		cards := make([]domain.Card, len(cardsSplit))

		for i, cardCode := range cardsSplit {
			cards[i] = domain.FromCode(cardCode)
		}

		opts = append(opts, service.WithCards(cards))
	}

	res, err := h.deckService.CreateDeck(c.Request().Context(), opts...)

	if err != nil {
		return mapError(c, err)
	}

	return c.JSON(http.StatusCreated, res)
}

// HandleOpenDeck handles the endpoint for opening a deck.
func (h *DeckEchoHandler) HandleOpenDeck(c echo.Context) error {
	uuid := c.Param(uuidParam)

	res, err := h.deckService.OpenDeck(c.Request().Context(), uuid)

	if err != nil {
		return mapError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

type drawCardsRequest struct {
	Amount int `json:"amount"`
}

// HandleDrawCars handles the endpoint for drawing cards from a deck.
func (h *DeckEchoHandler) HandleDrawCars(c echo.Context) error {
	uuid := c.Param(uuidParam)

	req := drawCardsRequest{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, buildErrorMap("Invalid body"))
	}

	if req.Amount < 0 {
		return c.JSON(http.StatusBadRequest, buildErrorMap("Invalid amount, must be greater or equal to 0"))
	}

	res, err := h.deckService.DrawCards(c.Request().Context(), uuid, req.Amount)

	if err != nil {
		return mapError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

func mapError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, domain.ErrDeckNotFound):
		return c.JSON(http.StatusBadRequest, buildErrorMap("The deck given wasn't found"))
	default:
		return err
	}
}

func buildErrorMap(message string) map[string]string {
	return map[string]string{
		"message": message,
	}
}

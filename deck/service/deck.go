// Package service contains the implementations for the use cases relating to decks.
package service

import (
	"context"
	"fmt"

	"github.com/cfagudelo96/toggle-test/deck/domain"
)

// DeckRepository represents the interface required for storing and retrieving decks.
type DeckRepository interface {
	Save(ctx context.Context, d *domain.Deck) error
	Get(ctx context.Context, uuid string) (*domain.Deck, error)
}

// DeckService handles the deck related use cases.
type DeckService struct {
	deckRepository DeckRepository
}

// NewDeckService returns a new DeckService.
func NewDeckService(r DeckRepository) *DeckService {
	return &DeckService{
		deckRepository: r,
	}
}

type deckCreationOptions struct {
	shuffled bool
	cards    []domain.Card
}

// DeckCreationOption is the interface implemented to allow options while creating a new Deck.
type DeckCreationOption interface {
	apply(*deckCreationOptions)
}

type shuffledOption bool

func (c shuffledOption) apply(o *deckCreationOptions) {
	o.shuffled = bool(c)
}

// Shuffled option to determine if a new deck should be shuffled or unshuffled.
func Shuffled(s bool) DeckCreationOption {
	return shuffledOption(s)
}

type cardsOption struct {
	cards []domain.Card
}

func (c cardsOption) apply(o *deckCreationOptions) {
	o.cards = c.cards
}

// WithCards allows to specify the cards in the new deck being created.
func WithCards(cards []domain.Card) DeckCreationOption {
	return cardsOption{cards: cards}
}

// CreateDeckOutput is the result of creating a new deck.
type CreateDeckOutput struct {
	DeckID    string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}

func createDeckOutputFromDeck(d *domain.Deck) CreateDeckOutput {
	return CreateDeckOutput{
		DeckID:    d.UUID,
		Shuffled:  d.Shuffled,
		Remaining: len(d.Cards),
	}
}

// CreateDeck creates a new deck. By default creates a shuffled complete deck, unless the options say otherwise.
// Returns an error if the new deck couldn't be saved.
func (s *DeckService) CreateDeck(ctx context.Context, opts ...DeckCreationOption) (CreateDeckOutput, error) {
	options := deckCreationOptions{
		shuffled: true,
		cards:    domain.CompleteDeckCards(),
	}

	for _, o := range opts {
		o.apply(&options)
	}

	d := domain.NewDeck(options.shuffled, options.cards)

	if err := s.deckRepository.Save(ctx, d); err != nil {
		return CreateDeckOutput{}, fmt.Errorf("saving the deck failed: %w", err)
	}

	return createDeckOutputFromDeck(d), nil
}

// OpenDeckOutput is the result of opening a deck.
type OpenDeckOutput struct {
	DeckID    string        `json:"deck_id"`
	Shuffled  bool          `json:"shuffled"`
	Remaining int           `json:"remaining"`
	Cards     []domain.Card `json:"cards"`
}

func openDeckOutputFromDeck(d *domain.Deck) OpenDeckOutput {
	return OpenDeckOutput{
		DeckID:    d.UUID,
		Shuffled:  d.Shuffled,
		Remaining: len(d.Cards),
		Cards:     d.Cards,
	}
}

// OpenDeck opens the deck with the given UUID.
// Returns an error if there is no deck with the given UUID.
func (s *DeckService) OpenDeck(ctx context.Context, uuid string) (OpenDeckOutput, error) {
	d, err := s.deckRepository.Get(ctx, uuid)

	if err != nil {
		return OpenDeckOutput{}, fmt.Errorf("getting the deck failed: %w", err)
	}

	return openDeckOutputFromDeck(d), nil
}

// DrawCardsOutput is the result of drawing cards from a deck.
type DrawCardsOutput struct {
	Cards []domain.Card `json:"cards"`
}

// DrawCards draws the given amount of cards from the deck with the given UUID.
// Returns an error if there is no deck with the given UUID or if saving the modified deck failed.
func (s *DeckService) DrawCards(ctx context.Context, uuid string, amount int) (DrawCardsOutput, error) {
	d, err := s.deckRepository.Get(ctx, uuid)

	if err != nil {
		return DrawCardsOutput{}, fmt.Errorf("getting the deck failed: %w", err)
	}

	drawnCards := d.Draw(amount)

	if err := s.deckRepository.Save(ctx, d); err != nil {
		return DrawCardsOutput{}, fmt.Errorf("saving the deck failed: %w", err)
	}

	return DrawCardsOutput{Cards: drawnCards}, nil
}

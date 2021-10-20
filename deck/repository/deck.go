// Package repository contains the implementations for deck repositories.
package repository

import (
	"context"

	"github.com/cfagudelo96/toggle-test/deck/domain"
)

// InMemoryDeckRepository represents a repository of decks implemented using memory.
type InMemoryDeckRepository struct {
	decks map[string]*domain.Deck
}

// NewInMemoryDeckRepository returns a new InMemoryDeckRepository.
func NewInMemoryDeckRepository() *InMemoryDeckRepository {
	return &InMemoryDeckRepository{
		decks: make(map[string]*domain.Deck),
	}
}

// Save saves the given deck in memory.
// Returns an error thinking about possible future implementations using some database.
func (r *InMemoryDeckRepository) Save(_ context.Context, d *domain.Deck) error {
	r.decks[d.UUID] = d

	return nil
}

// Get gets the deck with the given UUID. Returns an error if the deck is not found.
func (r *InMemoryDeckRepository) Get(_ context.Context, uuid string) (*domain.Deck, error) {
	d, ok := r.decks[uuid]

	if !ok {
		return nil, domain.ErrDeckNotFound
	}

	return d, nil
}

// Package service contains the implementations for the use cases relating to decks.
package service_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/cfagudelo96/toggle-test/deck/domain"
	"github.com/cfagudelo96/toggle-test/deck/service"
	"github.com/cfagudelo96/toggle-test/deck/service/mocks"
	"github.com/stretchr/testify/mock"
)

//go:generate mockery --name DeckRepository

func TestDeckService_CreateDeck(t *testing.T) {
	ctx := context.Background()
	deckRepositoryMock := &mocks.DeckRepository{}
	deckRepositoryMock.On("Save", ctx, mock.Anything).Return(nil)
	tests := []struct {
		name           string
		deckRepository service.DeckRepository
		opts           []service.DeckCreationOption
		want           service.CreateDeckOutput
		wantErr        bool
	}{
		{
			name: "returns an error if the repository fails to save",
			deckRepository: func() service.DeckRepository {
				m := &mocks.DeckRepository{}
				m.On("Save", ctx, mock.Anything).Return(errors.New("test"))
				return m
			}(),
			wantErr: true,
		},
		{
			name:           "works correctly with the shuffled option",
			deckRepository: deckRepositoryMock,
			opts:           []service.DeckCreationOption{service.Shuffled(false)},
			want: service.CreateDeckOutput{
				DeckID:    "some-deck-id",
				Shuffled:  false,
				Remaining: 52,
			},
		},
		{
			name:           "works correctly with the cards option",
			deckRepository: deckRepositoryMock,
			opts: []service.DeckCreationOption{
				service.WithCards([]domain.Card{{Value: "4", Suit: "HEARTS", Code: "4H"}, {Value: "5", Suit: "SPADES", Code: "5S"}}),
			},
			want: service.CreateDeckOutput{
				DeckID:    "some-deck-id",
				Shuffled:  true,
				Remaining: 2,
			},
		},
		{
			name:           "works correctly without options",
			deckRepository: deckRepositoryMock,
			want: service.CreateDeckOutput{
				DeckID:    "some-deck-id",
				Shuffled:  true,
				Remaining: 52,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewDeckService(tt.deckRepository)
			got, err := s.CreateDeck(ctx, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeckService.CreateDeck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got.DeckID = tt.want.DeckID
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeckService.CreateDeck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeckService_OpenDeck(t *testing.T) {
	ctx := context.Background()
	uuid := "some-deck-uuid"
	tests := []struct {
		name           string
		deckRepository service.DeckRepository
		want           service.OpenDeckOutput
		wantErr        bool
	}{
		{
			name: "returns an error if the repository fails to get the deck",
			deckRepository: func() service.DeckRepository {
				m := &mocks.DeckRepository{}
				m.On("Get", ctx, uuid).Return(nil, errors.New("test"))
				return m
			}(),
			wantErr: true,
		},
		{
			name: "works correctly if the repository finds the deck",
			deckRepository: func() service.DeckRepository {
				m := &mocks.DeckRepository{}
				m.On("Get", ctx, uuid).Return(&domain.Deck{
					UUID:     uuid,
					Shuffled: true,
					Cards:    []domain.Card{{Value: "4", Suit: "HEARTS", Code: "4H"}, {Value: "5", Suit: "SPADES", Code: "5S"}},
				}, nil)
				return m
			}(),
			want: service.OpenDeckOutput{
				DeckID:    uuid,
				Shuffled:  true,
				Remaining: 2,
				Cards:     []domain.Card{{Value: "4", Suit: "HEARTS", Code: "4H"}, {Value: "5", Suit: "SPADES", Code: "5S"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewDeckService(tt.deckRepository)
			got, err := s.OpenDeck(ctx, uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeckService.OpenDeck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeckService.OpenDeck() = %v, want %v", got, tt.want)
			}
		})
	}
}

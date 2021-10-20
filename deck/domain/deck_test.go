package domain_test

import (
	"reflect"
	"testing"

	"github.com/cfagudelo96/toggle-test/deck/domain"
)

func TestNewDeck(t *testing.T) {
	t.Run("works correctly with shuffling", func(t *testing.T) {
		c := []domain.Card{{Value: "3", Suit: "HEARTS", Code: "3H"}, {Value: "4", Suit: "SPADES", Code: "4S"}}
		got := domain.NewDeck(true, c)
		if len(got.Cards) != len(c) {
			t.Errorf("New deck length = %d, want %d", len(got.Cards), len(c))
		}
		if !got.Shuffled {
			t.Errorf("New deck should be shuffled")
		}
		if got.UUID == "" {
			t.Errorf("New deck should have an UUID")
		}
	})
	t.Run("works correctly without shuffling", func(t *testing.T) {
		c := []domain.Card{{Value: "4", Suit: "HEARTS", Code: "4H"}, {Value: "5", Suit: "SPADES", Code: "5S"}}
		want := &domain.Deck{
			UUID:     "some-uuid",
			Shuffled: false,
			Cards:    c,
		}
		got := domain.NewDeck(false, c)
		if got.UUID == "" {
			t.Errorf("New deck should have an UUID")
		}
		got.UUID = want.UUID
		if !reflect.DeepEqual(got, want) {
			t.Errorf("New deck = %v, want %v", got, want)
		}
	})
}

func TestCompleteDeckCards(t *testing.T) {
	t.Run("generates a complete deck correctly sorted", func(t *testing.T) {
		got := domain.CompleteDeckCards()
		if len(got) != 52 {
			t.Errorf("Complete deck len = %d, want 52", len(got))
		}

		if got[13*0].Code != "AC" && got[13*1].Code != "AD" && got[13*2].Code != "AH" && got[13*3].Code != "AS" {
			t.Error("The deck wasn't sorted")
		}
	})
}

func TestFromCode(t *testing.T) {
	tests := []struct {
		name string
		code string
		want domain.Card
	}{
		{
			name: "works with numbered cards",
			code: "4D",
			want: domain.Card{
				Value: "4",
				Suit:  "DIAMONDS",
				Code:  "4D",
			},
		},
		{
			name: "works for Ace cards",
			code: "AS",
			want: domain.Card{
				Value: "ACE",
				Suit:  "SPADES",
				Code:  "AS",
			},
		},
		{
			name: "works for Jack cards",
			code: "JH",
			want: domain.Card{
				Value: "JACK",
				Suit:  "HEARTS",
				Code:  "JH",
			},
		},
		{
			name: "works for Queen cards",
			code: "QC",
			want: domain.Card{
				Value: "QUEEN",
				Suit:  "CLUBS",
				Code:  "QC",
			},
		},
		{
			name: "works for King cards",
			code: "KS",
			want: domain.Card{
				Value: "KING",
				Suit:  "SPADES",
				Code:  "KS",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := domain.FromCode(tt.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testDeck() *domain.Deck {
	return &domain.Deck{
		UUID:     "test-uuid",
		Shuffled: true,
		Cards:    []domain.Card{{Value: "4", Suit: "HEARTS", Code: "4H"}, {Value: "5", Suit: "SPADES", Code: "5S"}},
	}
}

func TestDeck_Draw(t *testing.T) {
	tests := []struct {
		name     string
		deck     *domain.Deck
		amount   int
		want     []domain.Card
		wantDeck *domain.Deck
	}{
		{
			name:   "draws 1 card correctly",
			deck:   testDeck(),
			amount: 1,
			want:   []domain.Card{{Value: "4", Suit: "HEARTS", Code: "4H"}},
			wantDeck: &domain.Deck{
				UUID:     "test-uuid",
				Shuffled: true,
				Cards:    []domain.Card{{Value: "5", Suit: "SPADES", Code: "5S"}},
			},
		},
		{
			name:   "draws all cards correctly",
			deck:   testDeck(),
			amount: 2,
			want:   []domain.Card{{Value: "4", Suit: "HEARTS", Code: "4H"}, {Value: "5", Suit: "SPADES", Code: "5S"}},
			wantDeck: &domain.Deck{
				UUID:     "test-uuid",
				Shuffled: true,
				Cards:    []domain.Card{},
			},
		},
		{
			name:   "drawing more than all cards",
			deck:   testDeck(),
			amount: 3,
			want:   []domain.Card{{Value: "4", Suit: "HEARTS", Code: "4H"}, {Value: "5", Suit: "SPADES", Code: "5S"}},
			wantDeck: &domain.Deck{
				UUID:     "test-uuid",
				Shuffled: true,
				Cards:    []domain.Card{},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.deck.Draw(tt.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Deck.Draw() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.deck, tt.wantDeck) {
				t.Errorf("Deck after drawing = %v, want %v", tt.deck, tt.want)
			}
		})
	}
}

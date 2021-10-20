// Package domain contains the implementation of the domain of french playing cards
package domain

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	hearts   string = "HEARTS"
	diamonds string = "DIAMONDS"
	clubs    string = "CLUBS"
	spades   string = "SPADES"
)

const (
	ranksNumber = 13
	aceNumber   = 1
	jackNumber  = 11
	queenNumber = 12
	kingNumber  = 13
)

var (
	// ErrDeckNotFound error returned when a deck is not found in the system.
	ErrDeckNotFound = errors.New("deck_not_found")
)

func suitCode(s string) string {
	switch s {
	case hearts:
		return "H"
	case diamonds:
		return "D"
	case clubs:
		return "C"
	default:
		return "S"
	}
}

// Deck represents a french deck.
type Deck struct {
	UUID     string
	Shuffled bool
	Cards    []Card
}

// NewDeck creates a new deck with the cards given. If the shuffled flag is true, the deck gets shuffled.
func NewDeck(shuffled bool, cards []Card) *Deck {
	uuid := uuid.NewString()

	if shuffled {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	}

	return &Deck{
		UUID:     uuid,
		Shuffled: shuffled,
		Cards:    cards,
	}
}

// Draw draws the amount of cards given as parameter from the top of the deck.
// If the amount given is more than the number of cards in the deck, draws all the cards available.
func (d *Deck) Draw(amount int) []Card {
	if amount > len(d.Cards) {
		amount = len(d.Cards)
	}

	drawnCards := d.Cards[:amount]
	d.Cards = d.Cards[amount:]

	return drawnCards
}

// Card represents a card inside a french deck.
type Card struct {
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

// CompleteDeckCards returns a complete set of french deck cards sorted.
func CompleteDeckCards() []Card {
	suits := []string{clubs, diamonds, hearts, spades}
	cards := make([]Card, ranksNumber*len(suits))

	for i, s := range suits {
		for j := 1; j <= ranksNumber; j++ {
			cards[i*ranksNumber+j-1] = Card{
				Value: numberToValue(j),
				Suit:  s,
				Code:  fmt.Sprintf("%s%s", numberToCodePart(j), suitCode(s)),
			}
		}
	}

	return cards
}

// FromCode returns the card represented by the code given as a parameter.
func FromCode(code string) Card {
	codeSlice := strings.Split(code, "")

	var (
		value string
		suit  string
	)

	switch codeSlice[0] {
	case "A":
		value = "ACE"
	case "J":
		value = "JACK"
	case "Q":
		value = "QUEEN"
	case "K":
		value = "KING"
	default:
		value = codeSlice[0]
	}

	switch codeSlice[1] {
	case "H":
		suit = hearts
	case "D":
		suit = diamonds
	case "C":
		suit = clubs
	default:
		suit = spades
	}

	return Card{
		Value: value,
		Suit:  suit,
		Code:  code,
	}
}

func numberToValue(n int) string {
	switch n {
	case aceNumber:
		return "ACE"
	case jackNumber:
		return "JACK"
	case queenNumber:
		return "QUEEN"
	case kingNumber:
		return "KING"
	default:
		return strconv.Itoa(n)
	}
}

func numberToCodePart(n int) string {
	switch n {
	case aceNumber:
		return "A"
	case jackNumber:
		return "J"
	case queenNumber:
		return "Q"
	case kingNumber:
		return "K"
	default:
		return strconv.Itoa(n)
	}
}

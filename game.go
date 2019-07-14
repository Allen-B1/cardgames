package cardgames

import (
	"github.com/allen-b1/cards"
)

type Game interface {
	Players() []string // Returns a list of the player names
	Turn() int         // Whose turn it is; -1 if N/A
	Ended() bool

	Hands(player int) []cards.Deck // Returns the hands of all players from the given player's point of view
	Pile() cards.Deck              // Returns the pile in the middle

	Play(player int, cards []cards.Card) error // Make the given player play the given cards; passing is represented by playing no cards
}

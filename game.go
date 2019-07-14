/*
Package cardgames provides a set of utilities for working with card games.
*/

package cardgames

import (
	"github.com/allen-b1/cards"
)

// Type Game represents the state of a game.
//
// Each player has an index which is used for identification. Every method that takes in `player int` as an argument assumes that that argument is the index (in the slice returned by the method Players) of the player.
type Game interface {
	// Method Players returns a list of player names. The returned slice cannot be modified.
	Players() []string

	// Method Turn returns the current turn.
	Turn() int

	// Method Ended returns whether the game has ended (e.g. one person one).
	Ended() bool

	// Method Hands returns the hands from the player's point of view. The returned slice cannot be modified.
	Hands(player int) []cards.Deck

	// Method Pile returns the pile in the middle. The returned Deck cannot be modified.
	Pile() cards.Deck

	// Method Play makes the given player play the given cards.
	// Passing is represented by playing no cards. If this function returns an error,
	// it can't have also changed the game's state.
	Play(player int, cards []cards.Card) error
}

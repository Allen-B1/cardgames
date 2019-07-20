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

	// Method Hands returns the hands from the player's point of view. The returned slice cannot be modified.
	Hands(player int) []cards.Deck

	// Method Piles returns the piles. The returned slice cannot be modified.
	Piles() []cards.Deck

	// Method Winners returns an array showing what places the players recieved.
	// If a place is not yet decided, -1 should fill it.
	// The returned slice cannot be modified.
	Winners() []int

	// Method Play makes the given player play the given cards.
	// Passing is represented by playing no cards. If this function returns an error,
	// it can't have also changed the game's state.
	Play(player int, cards []cards.Card) error
}

// Function Type returns a string representation of the name of the card game.
func Type(g Game) string {
	switch g.(type) {
	case *President:
		return "president"
	default:
		return ""
	}
}

// Function Ended returns true if g.Winners() doesn't contain any negative numbers.
func Ended(g Game) bool {
	for _, winner := range g.Winners() {
		if winner < 0 {
			return false
		}
	}

	return true
}

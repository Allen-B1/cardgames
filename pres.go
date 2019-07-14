package cardgames

import (
	"fmt"
	"github.com/allen-b1/cards"
)

type PresGame struct {
	players  []string
	turn     int
	lastplay int  // the last person that played a card
	mode     uint // 0 for no mode

	hands []cards.Deck
	pile  cards.Deck
}

func NewPresGame(players []string) *PresGame {
	return &PresGame{
		players: players,
		turn:    0,                                // TODO: This should be whoever has 3S
		hands:   make([]cards.Deck, len(players)), // TODO: actually deal....
	}
}

func (g *PresGame) Players() []string {
	return g.players
}

func (g *PresGame) Turn() int {
	return int(g.turn)
}

func (g *PresGame) Ended() bool {
	return false
}

func (g *PresGame) Hands(player int) []cards.Deck {
	hands := make([]cards.Deck, len(g.players))
	for i := 0; i < len(g.players); i++ {
		if i == player {
			hands[i] = g.hands[i]
		} else {
			hands[i] = make(cards.Deck, len(g.hands[i]))
		}
	}
	return hands
}

func (g *PresGame) Pile() cards.Deck {
	return g.pile
}

func (g *PresGame) Play(player int, cards []cards.Card) error {
	if len(cards) == 0 {
		g.turn = (g.turn + 1) % len(g.players)

		// Clear if nobody can play
		if g.lastplay == g.turn {
			g.pile = nil
		}

		return nil
	}

	// DONT PLAY OUT OF TURN PLEASE
	if g.turn != player {
		return fmt.Errorf("playing out of turn: it's %s's turn", g.players[g.turn])
	}

	for _, card := range cards {
		has := false
		for _, handcard := range g.hands[player] {
			if card == handcard {
				has = true
				break
			}
		}
		if !has {
			return fmt.Errorf("can't play card: %v is not in your hand", card)
		}
	}

	if g.mode != 0 && int(g.mode) != len(cards) && len(cards) != 4 && cards[0].Value() != 2 {
		return fmt.Errorf("can't play cards: %d given but %d required", len(cards), g.mode)
	}

	// Check if 1) all cards have = face value and 2) not playing multiple bombs
	if len(cards) > 1 {
		for i := 0; i < len(cards)-1; i++ {
			if cards[i].Value() != cards[i+1].Value() {
				return fmt.Errorf("can't play card: %v and %v have different face values", cards[i], cards[i+1])
			}
		}

		if cards[0].Value() == 2 {
			return fmt.Errorf("can't play multiple 2s")
		}
	}

	// Check that it's greater than...
	if len(g.pile) != 0 && presLess(cards[0], g.pile[len(g.pile)-1]) {
		return fmt.Errorf("can't play card: %v is less than %v", cards[0], g.pile[len(g.pile)-1])
	}

	// Set mode
	g.mode = uint(len(cards))

	// Remove cards from hand
	for _, card := range cards {
		for i, handcard := range g.hands[player] {
			if card == handcard {
				g.hands[player] = append(g.hands[player][:i], g.hands[player][i+1:]...)
				i--
				break
			}
		}
	}

	if cards[0].Value() == 2 || len(cards) == 4 {
		// Bomb
		g.pile = nil
	} else {
		g.pile = append(g.pile, cards...)
		g.turn = (g.turn + 1) % len(g.players)
	}
	g.lastplay = player

	return nil
}

// returns true if c1 < c2
func presLess(c1 cards.Card, c2 cards.Card) bool {
	val1 := c1.Value()
	val2 := c2.Value()
	// Make A, 2 the highest cards
	if val1 <= 2 {
		val1 += 13
	}
	if val2 <= 2 {
		val2 += 13
	}

	return val1 < val2
}

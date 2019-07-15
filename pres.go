package cardgames

import (
	"fmt"
	"github.com/allen-b1/cards"
	"sort"
)

// Type PresGame implements Game and represents a game of President.
type PresGame struct {
	players  []string
	turn     int
	lastplay int  // the last person that played a card
	mode     uint // 0 for no mode

	hands   []cards.Deck
	pile    cards.Deck
	discard cards.Deck

	winners []int
	losers  []int // TODO: when somebody leaves
}

// Function NewPresGame creates a new game of President.
func NewPresGame(players []string) *PresGame {
	hands := make([]cards.Deck, len(players))
	fulldeck := cards.FullDeck()
	fulldeck.Shuffle()
	start := 0
	for index, card := range fulldeck {
		playerIndex := index % len(players)
		hands[playerIndex] = append(hands[playerIndex], card)
		if card == cards.New(3, cards.Spades) {
			start = playerIndex
		}
	}

	for _, hand := range hands {
		sort.Slice(hand, func(i, j int) bool {
			vali := hand[i].Value()
			valj := hand[j].Value()
			if vali == 1 {
				vali = 14
			}
			if valj == 1 {
				valj = 14
			}
			return vali < valj
		})
	}

	return &PresGame{
		players: players,
		turn:    start,
		hands:   hands,
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

func (g *PresGame) Piles() []cards.Deck {
	return []cards.Deck{g.pile, g.discard}
}

func (g *PresGame) bomb() {
	g.pile = nil
	g.mode = 0
	g.discard = append(g.discard, g.pile...)
}

func (g *PresGame) Play(player int, cards []cards.Card) error {
	defer func() {
		if len(g.players)-len(g.winners)-len(g.losers) <= 1 {
		outer:
			for player, _ := range g.players {
				combined := append(append([]int(nil), g.winners...), g.losers...)
				for _, winner := range combined {
					if winner == player {
						continue outer
					}
				}

				g.winners = append(g.winners, player)
				return
			}
		}
	}()

	// DONT PLAY OUT OF TURN PLEASE
	if g.turn != player {
		return fmt.Errorf("playing out of turn: it's %s's turn", g.players[g.turn])
	}

	if len(cards) == 0 {
		g.increaseTurnBy(1)

		// Clear if nobody can play
		if g.lastplay == g.turn {
			g.bomb()
		}

		return nil
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
		g.bomb()
	} else {
		g.pile = append(g.pile, cards...)
		// TODO: Skipping
		g.increaseTurnBy(1)
	}
	g.lastplay = player

	if len(g.hands[player]) == 0 {
		g.winners = append(g.winners, player)
		g.bomb()
	}

	return nil
}

func (g *PresGame) increaseTurnBy(n int) {
	targetTurn := g.turn
	for i := 0; i < n; i++ {
		targetTurn = (targetTurn + 1) % len(g.players)
		for _, winner := range g.winners {
			if targetTurn == winner {
				i--
				continue
			}
		}
	}
	g.turn = targetTurn
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

func (g *PresGame) Winners() []int {
	return g.winners
}

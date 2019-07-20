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
	g.discard = append(g.discard, g.pile...)
	g.pile = nil
	g.mode = 0
}

func (g *PresGame) Play(player int, cardlist []cards.Card) error {
	defer func() {
		if len(g.players)-len(g.winners) <= 1 {
		outer:
			for player, _ := range g.players {
				for _, winner := range g.winners {
					if winner == player {
						continue outer
					}
				}

				g.winners = append(g.winners, player)
				return
			}
		}
	}()

	// Check whether or not this will complete the set
	isCompletion := false
	combined := append(append(cards.Deck(nil), g.pile...), cardlist...)
	if len(combined) >= 4 {
		isCompletion = true
		for i := len(combined) - 4; i < len(combined)-1; i++ {
			if combined[i].Value() != combined[i+1].Value() {
				isCompletion = false
			}
		}
	}

	// DONT PLAY OUT OF TURN PLEASE
	if g.turn != player && !isCompletion {
		return fmt.Errorf("playing out of turn: it's %s's turn", g.players[g.turn])
	}

	// Pass
	if len(cardlist) == 0 {
		g.turn = g.after(g.turn, 1)

		// Clear if nobody can play
		if g.lastplay == g.turn {
			g.bomb()
		}

		return nil
	}

	// Check that all cardlist are in the hand
	for _, card := range cardlist {
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

	// Check that required # are played
	if g.mode != 0 && int(g.mode) != len(cardlist) && len(cardlist) != 4 && cardlist[0].Value() != 2 && !isCompletion {
		return fmt.Errorf("can't play cards: %d given but %d required", len(cardlist), g.mode)
	}

	// Check if 1) all cards have = face value and 2) not playing multiple bombs
	if len(cardlist) > 1 {
		for i := 0; i < len(cardlist)-1; i++ {
			if cardlist[i].Value() != cardlist[i+1].Value() {
				return fmt.Errorf("can't play card: %v and %v have different face values", cardlist[i], cardlist[i+1])
			}
		}

		if cardlist[0].Value() == 2 {
			return fmt.Errorf("can't play multiple 2s")
		}
	}

	// Check that it's greater than...
	if len(g.pile) != 0 && presLess(cardlist[0], g.pile[len(g.pile)-1]) {
		return fmt.Errorf("can't play card: %v is less than %v", cardlist[0], g.pile[len(g.pile)-1])
	}

	// Make sure that the player isn't bombing nothing
	if cardlist[0].Value() == 2 && len(g.pile) == 0 {
		return fmt.Errorf("can't bomb nothing")
	}

	// Set mode
	g.mode = uint(len(cardlist))

	// Remove cards from hand
	for _, card := range cardlist {
		for i, handcard := range g.hands[player] {
			if card == handcard {
				g.hands[player] = append(g.hands[player][:i], g.hands[player][i+1:]...)
				i--
				break
			}
		}
	}

	// Add cards to pile
	g.pile = append(g.pile, cardlist...)

	if cardlist[0].Value() == 2 || len(cardlist) == 4 || isCompletion { // Bomb
		g.bomb()
		g.turn = player
	} else {
		if g.mode == 1 && len(g.pile) != 0 && cardlist[0].Value() == g.pile[len(g.pile)-1].Value() {
			g.turn = g.after(g.turn, 2)
		} else {
			g.turn = g.after(g.turn, 1)
		}
	}

	g.lastplay = player

	if len(g.hands[player]) == 0 {
		g.winners = append(g.winners, player)
		g.lastplay = g.after(player, 1) // Prevent infinite loop if nobody can play
		g.bomb()
	}

	return nil
}

func (g *PresGame) after(turn int, n int) int {
	targetTurn := turn
outer:
	for i := 0; i < n; i++ {
		targetTurn = (targetTurn + 1) % len(g.players)
		for _, winner := range g.winners {
			if targetTurn == winner {
				i--
				continue outer
			}
		}
	}
	return targetTurn
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

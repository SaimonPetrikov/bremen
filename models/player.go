package models

import "github.com/SaimonPetrikov/bremen/models/cards"

type Player struct {
	*Combatant
}

func InitDefaultPlayerDeck() []PlayingCard {
	deck := []PlayingCard{}
	for i := 0; i < 15; i++ {
		deck = append(deck, cards.NewRat())
	}
	for i := 0; i < 5; i++ {
		deck = append(deck, cards.NewFatRat())
	}

	return deck
}

func NewPlayer(collection []PlayingCard) *Player {
	player := &Player{}
	player.Name = "Hero"
	player.HP = 10
	player.MaxHP = 10
	player.AP = 3
	player.MaxAP = 3
	for _, card := range collection {
		if card.InDeck() {
			player.Deck = append(player.Deck, card)
		}
	}

	return player
}

func (w *Player) Init() {

}

func (w *Player) Attack() {

}

func (w *Player) Trun() {

}

func (p *Player) AddCubeProgress() {
	if len(p.Cubes) == 0 {
		p.Cubes = append(p.Cubes, 1)
		return
	}
	lastIdx := len(p.Cubes) - 1
	if p.Cubes[lastIdx] < 3 {
		p.Cubes[lastIdx]++
	} else {
		p.Cubes = append(p.Cubes, 1)
	}
}
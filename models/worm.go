package models

import "github.com/SaimonPetrikov/bremen/models/cards"

type Worm struct {
	*Combatant
}

func NewWorm() *Worm {
	worm := &Worm{}
	worm.Name = "Giant Worm"
	worm.HP = 12
	worm.MaxHP = 12
	worm.AP = 2
	worm.MaxAP = 2
	worm.AIPattern = "swarm"
	for i := 0; i < 15; i++ {
		worm.Deck = append(worm.Deck, cards.NewWorm())
	}

	return worm
}

func (w *Worm) Init() {

}

func (w *Worm) Attack() {

}

func (w *Worm) Trun() {

}
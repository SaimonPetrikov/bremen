package models

import (
	"time"

	"math/rand"
)

type PlayingCard interface {
	Death() int
	Sacrifice() int
	UpdateStateAfterDeath()
	UpdateStateAfterChangeTime()
	InDeck() bool
}

type Combatant struct {
	Name      string
	HP        int
	MaxHP     int
	Blood     int
	Cubes     []int // НОВОЕ: Массив кубиков вместо костей
	AP        int
	MaxAP     int
	Lanes     [4]PlayingCard
	Buildings [4]PlayingCard
	Hand      []PlayingCard
	Deck      []PlayingCard

	// НОВОЕ: Тактика поведения ИИ
	// "swarm" — Червь (забивает линии и прет массой)
	// "hunter" — Волк (бьет в самые уязвимые точки / охотится на слабого)
	// "default" — Стандартный ИИ
	AIPattern string
}

func (c *Combatant) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(c.Deck), func(i, j int) { c.Deck[i], c.Deck[j] = c.Deck[j], c.Deck[i] })
}
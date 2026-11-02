package cards

type Worm struct {
	*Card
}

func NewWorm() *Worm {
	card := &Worm{}
	card.Name = "Worm"
	card.Attack = 1
	card.HP = 2
	card.Cost = 0
	return card
}

func (w *Worm) Death() int {
	return 1
}

func (w *Worm) Sacrifice() int {
	return 1
}

func (w *Worm) UpdateStateAfterDeath() {

}

func (w *Worm) UpdateStateAfterChangeTime() {

}
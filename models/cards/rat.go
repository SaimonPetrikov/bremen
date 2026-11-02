package cards

type Rat struct {
	*Card
}

func NewRat() *Rat {
	card := &Rat{}
	card.Name = "Rat"
	card.Attack = 1
	card.HP = 2
	card.Cost = 0
	card.IsActive = true
	return card
}

func (r *Rat) Death() int {
	return 1
}

func (r *Rat) Sacrifice() int {
	return 1
}

func (r *Rat) UpdateStateAfterDeath() {

}

func (r *Rat) UpdateStateAfterChangeTime() {

}
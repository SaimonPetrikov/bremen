package cards

type FatRat struct {
	*Card
}

func NewFatRat() *FatRat {
	card := &FatRat{}
	card.Name = "Fat rat"
	card.Attack = 2
	card.HP = 4
	card.Cost = 1
	card.IsActive = true
	return card
}

func (r *FatRat) Death() int {
	return 1
}

func (r *FatRat) Sacrifice() int {
	return 1
}

func (r *FatRat) UpdateStateAfterDeath() {

}

func (r *FatRat) UpdateStateAfterChangeTime() {

}
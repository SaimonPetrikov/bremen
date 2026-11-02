package cards

type Rooster struct {
	*Card
}

func NewRooster() *Rooster {
	return &Rooster{}
}

func (r *Rooster) Death() int {
	return 1
}

func (r *Rooster) Sacrifice() int {
	return 1
}

func (r *Rooster) UpdateStateAfterDeath() {

}

func (r *Rooster) UpdateStateAfterChangeTime() {

}
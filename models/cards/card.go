package cards

type Card struct {
	Name     string
	Attack   int
	HP       int
	Cost     int
	IsSpell  bool
	Effect   string
	IsActive bool  // Взята ли карта в боевую колоду (макс 20)
	Original *Card // Ссылка на оригинальную карту в коллекции (для эволюции)
	IsBuilding bool  // НОВОЕ: Является ли карта Зданием/Аурой
	TransformOnDeath   string // НОВОЕ: Имя карты при смерти (например, "Headless Cockroach")
	IsTempCopy         bool   // Временная боевая копия
	// Две формы Дровосека
	CanTransform     bool   // Может ли менять форму
	AltName          string // Имя второй формы (например, "Woodcutter (Beast)")
	AltAttack        int    // Атака второй формы
	AltHP            int    // HP второй формы
	IsBeastForm      bool   // Текущая форма (false = человек, true = зверь)

	// НОВОЕ: Свойство вампира
	IsVampire        bool
}

func (c *Card) InDeck() bool {
	return c.IsActive
} 
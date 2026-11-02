package scenes

import (
	"fmt"
	"image/color"

	"github.com/SaimonPetrikov/bremen/core"
	"github.com/SaimonPetrikov/bremen/models"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MillScene struct {
	State      *core.GameState
	Stage      int // 0: Разговор с мельником, 1: Ночь/Спать, 2: Утро/Крик петуха, 3: Разговор с петухом
}

func (s *MillScene) Update(g *core.Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if s.Stage == 0 {
			s.Stage = 1 // Идем спать
		} else if s.Stage == 1 {
			s.Stage = 2 // Утро, крик петуха
		} else if s.Stage == 2 {
			s.Stage = 3 // Диалог с петухом
		} else if s.Stage == 3 {
			s.State.PubQuestStep = 4
			
			// Добавляем Петуха в коллекцию игрока с механикой «Без головы»
			s.State.PlayerDeck = append(s.State.PlayerDeck, &models.Card{
				Name:             "Rooster",
				Attack:           2, // 2 атаки
				HP:               2, // 2 здоровья
				Cost:             0,
				TransformOnDeath: "Headless Rooster",
				IsActive:         false,
			})
			
			g.ChangeScene(&ForestEdgeScene{State: s.State}) // Или на опушку леса
		}
	}
	return nil
}

func (s *MillScene) Draw(screen *ebiten.Image) {
	if s.Stage >= 2 {
		screen.Fill(color.RGBA{100, 150, 100, 255}) // Утро
	} else {
		screen.Fill(color.RGBA{25, 30, 40, 255})    // Ночь
	}

	ebitenutil.DebugPrintAt(screen, "=== THE OLD MILL ===", 20, 20)

	if s.Stage == 0 {
		ebitenutil.DebugPrintAt(screen, "Miller: Ah, a traveler. The grain hasn't arrived from the city yet.", 20, 60)
		ebitenutil.DebugPrintAt(screen, "Miller: It's getting dark. You can sleep on the floor if you'd like.", 20, 80)
		ebitenutil.DebugPrintAt(screen, "[SPACE] Agree and go to sleep", 20, 130)
	} else if s.Stage == 1 {
		ebitenutil.DebugPrintAt(screen, "You fall asleep on the hard floor. Outside, the wind howls...", 20, 60)
		ebitenutil.DebugPrintAt(screen, "[SPACE] Wait for morning", 20, 100)
	} else if s.Stage == 2 {
		ebitenutil.DebugPrintAt(screen, "COCK-A-DOODLE-DOO!", 20, 60)
		ebitenutil.DebugPrintAt(screen, "You wake up suddenly from a piercing, desperate cry!", 20, 80)
		ebitenutil.DebugPrintAt(screen, "[SPACE] Look around", 20, 130)
	} else if s.Stage == 3 {
		ebitenutil.DebugPrintAt(screen, "Rooster: Listen to me! The miller is planning to put me in the soup pot today!", 20, 60)
		ebitenutil.DebugPrintAt(screen, "Rooster: Take me with you on your journey, I beg you!", 20, 80)
		ebitenutil.DebugPrintAt(screen, "[SPACE] Take the Rooster into your deck", 20, 130)
	}
}

type InventoryScene struct {
	State         *core.GameState
	PreviousScene core.Scene
	CursorIdx     int
}

func (s *InventoryScene) Update(g *core.Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyI) || inpututil.IsKeyJustPressed(ebiten.KeyX) {
		g.ChangeScene(s.PreviousScene)
		return nil
	}

	total := len(s.State.PlayerDeck)
	cols := 7

	if total > 0 {
		// Управление курсором
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) { s.CursorIdx = (s.CursorIdx + 1) % total }
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) { s.CursorIdx--; if s.CursorIdx < 0 { s.CursorIdx = total - 1 } }
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) { s.CursorIdx += cols; if s.CursorIdx >= total { s.CursorIdx = total - 1 } }
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) { s.CursorIdx -= cols; if s.CursorIdx < 0 { s.CursorIdx = 0 } }

		// Включение/выключение карты в колоду
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			card := s.State.PlayerDeck[s.CursorIdx]
			if card.IsActive {
				card.IsActive = false // Убрать из колоды можно всегда
			} else {
				// Считаем активные карты
				activeCount := 0
				for _, c := range s.State.PlayerDeck {
					if c.IsActive { activeCount++ }
				}
				if activeCount < 20 {
					card.IsActive = true // Добавляем, только если меньше 20
				}
			}
		}
	}
	return nil
}

func (s *InventoryScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 40, 255}) 

	activeCount := 0
	for _, c := range s.State.PlayerDeck {
		if c.IsActive { activeCount++ }
	}

	ebitenutil.DebugPrintAt(screen, "=== DECK BUILDING ===", 20, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Active Deck: %d / 20", activeCount), 200, 20)
	ebitenutil.DebugPrintAt(screen, "[SPACE] Toggle Card  |  [I] Close", 400, 20)

	// Статистика Героя
	ebitenutil.DrawRect(screen, 20, 60, 80, 80, color.RGBA{50, 150, 255, 255}) 
	ebitenutil.DebugPrintAt(screen, "HERO", 20, 150)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP: %d / 10", s.State.PlayerHP), 20, 170)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Collection: %d", len(s.State.PlayerDeck)), 20, 190)

	if len(s.State.PlayerDeck) == 0 { return }

	startX, startY := 140, 60
	cardW, cardH, gapX, gapY := 60, 80, 10, 10
	cols := 7 

	for i, card := range s.State.PlayerDeck {
		row := i / cols
		col := i % cols
		cx := startX + col*(cardW+gapX)
		cy := startY + row*(cardH+gapY)

		bgColor := color.RGBA{200, 200, 200, 255}
		if card.IsSpell { 
			bgColor = color.RGBA{150, 100, 200, 255} // Заклинание (фиолетовая)
		} else if card.IsBuilding {
			bgColor = color.RGBA{210, 105, 30, 255} // Здание (оранжевая)
		} else if card.Name == "Fat Rat" || card.Name == "Rat King" { 
			bgColor = color.RGBA{100, 200, 100, 255} 
		} else if card.Name == "Bone Wall" { 
			bgColor = color.RGBA{120, 120, 120, 255} 
		}

		// Если карта под курсором - рисуем желтую рамку
		if i == s.CursorIdx {
			ebitenutil.DrawRect(screen, float64(cx-3), float64(cy-3), float64(cardW+6), float64(cardH+6), color.RGBA{255, 255, 0, 255})
		}

		ebitenutil.DrawRect(screen, float64(cx), float64(cy), float64(cardW), float64(cardH), bgColor)
		ebitenutil.DebugPrintAt(screen, card.Name, cx+5, cy+10)
		if card.Cost > 0 { ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Bld:%d", card.Cost), cx+5, cy+30) }
		if !card.IsSpell { ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d/%d", card.Attack, card.HP), cx+5, cy+65) }

		// Если карта активна (в колоде) - рисуем зеленую метку [V]
		if card.IsActive {
			ebitenutil.DrawRect(screen, float64(cx+cardW-15), float64(cy+5), 10, 10, color.RGBA{0, 255, 0, 255})
		}
	}
}

// --- DEBUG MENU СЦЕНА ---
type DebugMenuScene struct {
	State *core.GameState
}

// Вспомогательная функция для генерации полной колоды игрока для дебаг-боев
func buildDebugPlayerDeck() []*models.Card {
	var deck []*models.Card
	for i := 0; i < 15; i++ { deck = append(deck, &models.Card{Name: "Rat", Attack: 1, HP: 2, Cost: 0}) }
	for i := 0; i < 5; i++ { deck = append(deck, &models.Card{Name: "Fat Rat", Attack: 2, HP: 4, Cost: 1}) }
	return deck
}

func (s *DebugMenuScene) Update(g *core.Game) error {
	var enemy *models.Combatant
	var player *models.Combatant

	// Создаем базового игрока для теста (с полным здоровьем и полной колодой)
	player = &models.Combatant{Name: "Hero (Debug)", HP: 10, MaxHP: 10, AP: 3, MaxAP: 3, Deck: buildDebugPlayerDeck()}

	// [1] БОЙ: ЧЕРВЬ
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		enemy = &models.Combatant{Name: "Giant Worm", HP: 10, MaxHP: 10, AP: 2, MaxAP: 2}
		for i := 0; i < 15; i++ { enemy.Deck = append(enemy.Deck, &models.Card{Name: "Worm", Attack: 1, HP: 1, Cost: 0}) }
	}
	
	// [2] БОЙ: БАНДИТ
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		enemy = &models.Combatant{Name: "Bandit Boss", HP: 15, MaxHP: 15, AP: 3, MaxAP: 3}
		for i := 0; i < 12; i++ { enemy.Deck = append(enemy.Deck, &models.Card{Name: "Thief", Attack: 1, HP: 2}) }
		for i := 0; i < 8; i++ { enemy.Deck = append(enemy.Deck, &models.Card{Name: "Brute", Attack: 3, HP: 4, Cost: 1}) }
	}

	// [3] БОЙ: СКЕЛЕТ
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		enemy = &models.Combatant{Name: "Skeleton Lord", HP: 25, MaxHP: 25, AP: 3, MaxAP: 3}
		for i := 0; i < 10; i++ { enemy.Deck = append(enemy.Deck, &models.Card{Name: "Skeleton", Attack: 1, HP: 2, Cost: 0}) }
		for i := 0; i < 4; i++ { enemy.Deck = append(enemy.Deck, &models.Card{Name: "Skel Archer", Attack: 2, HP: 1, Cost: 1}) }
		for i := 0; i < 4; i++ { enemy.Deck = append(enemy.Deck, &models.Card{Name: "Skel Wheel", Attack: 3, HP: 2, Cost: 2}) }
		for i := 0; i < 2; i++ { enemy.Deck = append(enemy.Deck, &models.Card{Name: "Skel Giant", Attack: 4, HP: 5, Cost: 3}) }
	}

	// Если выбран враг - запускаем Арену в режиме дебага
	if enemy != nil {
		g.ChangeScene(&CombatScene{
			State:       s.State,
			Player:      player,
			Enemy:       enemy,
			IsDebugMode: true, // ВАЖНО: Включаем дебаг!
		})
	}

	// [X] Выход
	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		g.ChangeScene(&HubScene{State: s.State})
	}
	return nil
}

func (s *DebugMenuScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "=== DEVELOPER DEBUG MODE ===", 20, 20)
	ebitenutil.DebugPrintAt(screen, "Choose an encounter to test:", 20, 60)
	
	ebitenutil.DebugPrintAt(screen, "[1] Test Fight: Giant Worm (Tutorial 1)", 20, 100)
	ebitenutil.DebugPrintAt(screen, "[2] Test Fight: Bandit Boss (Tutorial 2)", 20, 120)
	ebitenutil.DebugPrintAt(screen, "[3] Test Fight: Skeleton Lord (Final)", 20, 140)
	
	ebitenutil.DebugPrintAt(screen, "[X] Return to Hub", 20, 400)
}


// --- ЛЕС (Tutorial 2: Bones & Final Boss) ---
type ForestScene struct {
	State *core.GameState

	HeroX float64
	HeroY float64

	isInitialized bool
	HeroImg       *ebiten.Image
	MapImg        *ebiten.Image
}

func (s *ForestScene) Init() {
	if s.isInitialized { return }

	s.HeroImg = ebiten.NewImage(32, 32)
	s.HeroImg.Fill(color.RGBA{50, 150, 255, 255})

	s.MapImg = ebiten.NewImage(640, 480)
	s.MapImg.Fill(color.RGBA{20, 60, 20, 255}) // Темно-зеленый лес

	// Тропа обратно в Хаб (Сверху экрана)
	ebitenutil.DrawRect(s.MapImg, 280, 0, 80, 40, color.RGBA{100, 100, 100, 255}) 

	// Герой появляется у верхней тропы
	s.HeroX = 300
	s.HeroY = 60

	s.isInitialized = true
}

func (s *ForestScene) Update(g *core.Game) error {
	s.Init()

	speed := 4.0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) { s.HeroX -= speed }
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) { s.HeroX += speed }
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) { s.HeroY -= speed }
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) { s.HeroY += speed }

	if inpututil.IsKeyJustPressed(ebiten.KeyI) { g.ChangeScene(&InventoryScene{State: s.State, PreviousScene: s}); return nil }

	// Ограничения экрана (весь экран)
	if s.HeroX < 0 { s.HeroX = 0 }
	if s.HeroX > 640-32 { s.HeroX = 640 - 32 }
	if s.HeroY < 0 { s.HeroY = 0 }
	if s.HeroY > 480-32 { s.HeroY = 480 - 32 }

	isNearExit := s.HeroX > 260 && s.HeroX < 360 && s.HeroY < 60
	isNearEnemy := s.HeroX > 240 && s.HeroX < 400 && s.HeroY > 200 && s.HeroY < 340

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if isNearExit {
			g.ChangeScene(&HubScene{State: s.State})
		} else if isNearEnemy {
			if s.State.PubQuestStep == 8 {
				// Бой с Бандитом
				bandit := &models.Combatant{Name: "Bandit Boss", HP: 15, MaxHP: 15, AP: 3, MaxAP: 3}
				for i := 0; i < 12; i++ { bandit.Deck = append(bandit.Deck, &models.Card{Name: "Thief", Attack: 1, HP: 2}) }
				for i := 0; i < 8; i++ { bandit.Deck = append(bandit.Deck, &models.Card{Name: "Brute", Attack: 3, HP: 4, Cost: 1}) }

				player := &models.Combatant{Name: "Hero", HP: s.State.PlayerHP, MaxHP: 10, AP: 3, MaxAP: 3}
				for _, c := range s.State.PlayerDeck { cardCopy := *c; player.Deck = append(player.Deck, &cardCopy) }
				g.ChangeScene(&CombatScene{State: s.State, Player: player, Enemy: bandit})
			
			} else if s.State.PubQuestStep == 9 {
				for i := 0; i < 2; i++ {
					s.State.PlayerDeck = append(s.State.PlayerDeck, &models.Card{
						Name: "Caravan", IsBuilding: true, Cost: 1, IsActive: false,
					})
				}
				// Бой со Скелетом
				skeleton := &models.Combatant{Name: "Skeleton Lord", HP: 25, MaxHP: 25, AP: 3, MaxAP: 3}
				for i := 0; i < 10; i++ { skeleton.Deck = append(skeleton.Deck, &models.Card{Name: "Skeleton", Attack: 1, HP: 2, Cost: 0}) }
				for i := 0; i < 4; i++ { skeleton.Deck = append(skeleton.Deck, &models.Card{Name: "Skel Archer", Attack: 2, HP: 1, Cost: 1}) }
				for i := 0; i < 4; i++ { skeleton.Deck = append(skeleton.Deck, &models.Card{Name: "Skel Wheel", Attack: 3, HP: 2, Cost: 2}) }
				for i := 0; i < 2; i++ { skeleton.Deck = append(skeleton.Deck, &models.Card{Name: "Skel Giant", Attack: 4, HP: 5, Cost: 3}) }

				player := &models.Combatant{Name: "Hero", HP: s.State.PlayerHP, MaxHP: 10, AP: 3, MaxAP: 3}
				for _, c := range s.State.PlayerDeck {
					if c.IsActive {
						cardCopy := *c
						cardCopy.Original = c // Связываем с оригиналом
						player.Deck = append(player.Deck, &cardCopy)
					}
				}
				g.ChangeScene(&CombatScene{State: s.State, Player: player, Enemy: skeleton})
			}
		}
	}
	return nil
}

func (s *ForestScene) Draw(screen *ebiten.Image) {
	s.Init()

	screen.DrawImage(s.MapImg, &ebiten.DrawImageOptions{})

	// Динамическая отрисовка врагов по центру Леса (X:300, Y:260)
	if s.State.PubQuestStep == 8 {
		ebitenutil.DrawRect(screen, 290, 250, 60, 60, color.RGBA{139, 69, 19, 255}) // Коричневый Бандит
		ebitenutil.DrawRect(screen, 360, 260, 80, 40, color.RGBA{200, 180, 150, 255}) // Караван с пивом
		ebitenutil.DebugPrintAt(screen, "Bandit Boss", 285, 230)
	} else if s.State.PubQuestStep == 9 {
		ebitenutil.DrawRect(screen, 290, 250, 60, 60, color.RGBA{220, 220, 220, 255}) // Белый Скелет
		ebitenutil.DebugPrintAt(screen, "Skeleton Lord", 280, 230)
	}
	
	heroOp := &ebiten.DrawImageOptions{}
	heroOp.GeoM.Translate(s.HeroX, s.HeroY)
	screen.DrawImage(s.HeroImg, heroOp)

	// UI Фон
	ebitenutil.DrawRect(screen, 0, 0, 640, 75, color.RGBA{0, 0, 0, 180})
	ebitenutil.DebugPrintAt(screen, "=== THE DARK FOREST ===", 10, 10)

	isNearExit := s.HeroX > 260 && s.HeroX < 360 && s.HeroY < 60
	isNearEnemy := s.HeroX > 240 && s.HeroX < 400 && s.HeroY > 200 && s.HeroY < 340

	if isNearExit {
		ebitenutil.DebugPrintAt(screen, "[SPACE] Return to Hub", 10, 30)
	} else if isNearEnemy {
		if s.State.PubQuestStep == 8 {
			ebitenutil.DebugPrintAt(screen, "TUTORIAL: Use Bones to cast Rat Poison and build Bone Walls!", 10, 30)
			ebitenutil.DebugPrintAt(screen, "Deal 10 damage to Evolve your cards!", 10, 45)
			ebitenutil.DebugPrintAt(screen, "[SPACE] FIGHT THE BANDIT!", 10, 60)
		} else if s.State.PubQuestStep == 9 {
			ebitenutil.DebugPrintAt(screen, "The Bandit's corpse rises as a SKELETON LORD!", 10, 30)
			ebitenutil.DebugPrintAt(screen, "FINAL BOSS: No more tutorials. Survive.", 10, 45)
			ebitenutil.DebugPrintAt(screen, "[SPACE] FIGHT THE SKELETON LORD!", 10, 60)
		}
	} else if s.State.PubQuestStep == 10 {
		ebitenutil.DebugPrintAt(screen, "The forest is safe. The caravan is saved.", 10, 30)
	}
}

// --- ОПУШКА ЛЕСА (Forest Edge) ---
type ForestEdgeScene struct {
	State *core.GameState
}

func (s *ForestEdgeScene) Update(g *core.Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyI) { g.ChangeScene(&InventoryScene{State: s.State, PreviousScene: s}); return nil }
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.State.PubQuestStep = 5 // Переходим к следующему этапу квеста
		
		// Добавляем Дровосека и Собаку в колоду игрока
		// Дровосек сильный (3 ATK, 4 HP), но потом трагически погибнет в бою с волком
		s.State.PlayerDeck = append(s.State.PlayerDeck, &models.Card{
			Name:   "Woodcutter",
			Attack: 3,
			HP:     4,
			Cost:   1,
		})
		// Собака (2 ATK, 2 HP)
		s.State.PlayerDeck = append(s.State.PlayerDeck, &models.Card{
			Name:   "Dog",
			Attack: 2,
			HP:     2,
			Cost:   1,
		})

		g.ChangeScene(&CrossroadScene{State: s.State})
	}
	return nil
}

func (s *ForestEdgeScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 40, 20, 255})
	ebitenutil.DebugPrintAt(screen, "=== FOREST EDGE ===", 20, 20)
	ebitenutil.DebugPrintAt(screen, "Leaving the mill, you reach the edge of the dark woods.", 20, 60)
	ebitenutil.DebugPrintAt(screen, "A desperate Woodcutter and his loyal Dog block your path.", 20, 90)
	ebitenutil.DebugPrintAt(screen, "Woodcutter: 'My daughter went into the forest to visit Grandma and hasn't returned!'", 20, 120)
	ebitenutil.DebugPrintAt(screen, "Woodcutter: 'Please, help me search for her! We will join your journey.'", 20, 150)
	ebitenutil.DebugPrintAt(screen, "[SPACE] Agree to help. Woodcutter & Dog join your deck.", 20, 200)
}


// --- РАЗВИЛКА В ЛЕСУ (Crossroad Scene) ---
type CrossroadScene struct {
	State *core.GameState
}

func (s *CrossroadScene) Update(g *core.Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyI) { g.ChangeScene(&InventoryScene{State: s.State, PreviousScene: s}); return nil }
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		// Пес чует след и срывается в сторону домика Бабушки!
		s.State.PubQuestStep = 6
		g.ChangeScene(&GrandmaHouseScene{State: s.State})
	}
	return nil
}

func (s *CrossroadScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{15, 25, 15, 255})
	ebitenutil.DebugPrintAt(screen, "=== THE FOREST CROSSROAD ===", 20, 20)
	ebitenutil.DebugPrintAt(screen, "You walk deeper into the woods until you reach a split path.", 20, 60)
	ebitenutil.DebugPrintAt(screen, "One path leads toward the old trade caravan route. The other goes deep into the thicket.", 20, 90)
	ebitenutil.DebugPrintAt(screen, "Suddenly, the Woodcutter's Dog lets out a sharp bark...", 20, 130)
	ebitenutil.DebugPrintAt(screen, "The Dog catches a scent and darts headfirst down the path toward Grandma's house!", 20, 160)
	ebitenutil.DebugPrintAt(screen, "Woodcutter: 'Wait! Follow the dog! My daughter is that way!'", 20, 190)
	ebitenutil.DebugPrintAt(screen, "[SPACE] Chase after the dog towards Grandma's House", 20, 240)
}

// --- ДОМ БАБУШКИ (Grandma's House) ---
type GrandmaHouseScene struct {
	State *core.GameState
}

func (s *GrandmaHouseScene) Update(g *core.Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyI) { g.ChangeScene(&InventoryScene{State: s.State, PreviousScene: s}); return nil }
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.State.PubQuestStep = 7
		
		// 1. Пес погибает в сюжетной сцене — удаляем карту обычной Собаки из колоды навсегда
		var newDeck []*models.Card
		for _, c := range s.State.PlayerDeck {
			if c.Name == "Dog" { continue } // Пес погиб защищая дом
			
			// Настраиваем Дровосека: у него теперь две формы (Человек и Зверь-Пес)
			if c.Name == "Woodcutter" {
				c.CanTransform = true
				c.AltName = "Woodcutter (Feral Hound)"
				c.AltAttack = 5
				c.AltHP = 4
			}
			newDeck = append(newDeck, c)
		}
		s.State.PlayerDeck = newDeck

		// 2. Создаем босса: Гигантский Волк
		giantWolf := &models.Combatant{
			Name:      "Giant Wolf",
			HP:        28,
			MaxHP:     28,
			AP:        3,
			MaxAP:     3,
			AIPattern: "hunter", // Охотится на слабейших!
		}
		for i := 0; i < 12; i++ {
			giantWolf.Deck = append(giantWolf.Deck, &models.Card{Name: "Wolf Fang", Attack: 3, HP: 2, Cost: 1})
		}
		for i := 0; i < 12; i++ {
			giantWolf.Deck = append(giantWolf.Deck, &models.Card{Name: "Wolf test", Attack: 1, HP: 2, Cost: 0})
		}

		// 3. Собираем бойца-игрока
		player := &models.Combatant{
			Name:  "Hero (Mule)",
			HP:    s.State.PlayerHP,
			MaxHP: 10,
			AP:    3,
			MaxAP: 3,
		}
		for _, c := range s.State.PlayerDeck {
			if c.IsActive {
				copyCard := *c
				copyCard.Original = c
				player.Deck = append(player.Deck, &copyCard)
			}
		}

		// Переходим в бой против Гигантского Волка
		g.ChangeScene(&CombatScene{
			State:  s.State,
			Player: player,
			Enemy:  giantWolf,
		})
	}
	return nil
}

func (s *GrandmaHouseScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{40, 10, 10, 255})
	
	ebitenutil.DebugPrintAt(screen, "=== GRANDMA'S HOUSE: THE TRAGEDY ===", 20, 20)
	ebitenutil.DebugPrintAt(screen, "You burst into Grandma's cottage, following the tracks.", 20, 60)
	ebitenutil.DebugPrintAt(screen, "A massive, monstrous Giant Wolf stands over the bloody ruins.", 20, 90)
	ebitenutil.DebugPrintAt(screen, "Before anyone can react, the beast tears into the loyal Dog, killing it instantly.", 20, 120)
	ebitenutil.DebugPrintAt(screen, "The Woodcutter freezes... and then, a terrifying rage takes over.", 20, 150)
	ebitenutil.DebugPrintAt(screen, "His bones snap and shift. He takes his true, bestial form — a giant hound of vengeance.", 20, 180)
	ebitenutil.DebugPrintAt(screen, "[SPACE] Avenge the fallen dog. Fight the Giant Wolf!", 20, 230)
}
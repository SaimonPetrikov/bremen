package scenes

import (
	"fmt"
	"image/color"
	"math/rand"
	"sort"
	"time"

	"github.com/SaimonPetrikov/bremen/core"
	"github.com/SaimonPetrikov/bremen/models"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Combatant interface {
	Init()
	Attack()
	Trun()
	Shuffle()
}

type CombatScene struct {
	State  *core.GameState
	Player Combatant 
	Enemy  Combatant 


	SelectedHandIdx int
	isInitialized   bool
	Message         string
	MessageTimer    int

	ActionPhase     int 
	DelayTimer      int 
	TurnNumber      int

	DamageProgress  int
	DamageThreshold int
	UpgradeOptions  []*models.Card

	TutorialStep int
	IsDebugMode  bool

	RatImg *ebiten.Image
	IsNight bool
}

func getUpgradeInfo(name string) (bool, string, int, int, int) {
	if name == "Rat" { return true, "Fat Rat", 2, 4, 1 }
	if name == "Fat Rat" { return true, "Rat King", 3, 6, 2 }
	return false, "", 0, 0, 0
}

// НОВОЕ: Логика Эволюции Куба Смерти при смерти существа
func (s *CombatScene) AddCubeProgress() {
	if len(s.Player.Cubes) == 0 {
		s.Player.Cubes = append(s.Player.Cubes, 1)
		return
	}
	lastIdx := len(s.Player.Cubes) - 1
	if s.Player.Cubes[lastIdx] < 3 {
		s.Player.Cubes[lastIdx]++
	} else {
		s.Player.Cubes = append(s.Player.Cubes, 1)
	}
}

// Поиск и сжигание нужного куба для крафта
func (s *CombatScene) consumeCube(minValue int) bool {
	bestIdx := -1
	bestVal := 999
	for i, val := range s.Player.Cubes {
		if val >= minValue && val < bestVal {
			bestVal = val
			bestIdx = i
		}
	}
	if bestIdx != -1 {
		s.Player.Cubes = append(s.Player.Cubes[:bestIdx], s.Player.Cubes[bestIdx+1:]...)
		return true
	}
	return false
}

func (s *CombatScene) DrawCard(c *models.Combatant, showMsg bool) {
	if len(c.Deck) > 0 && len(c.Hand) < 5 {
		c.Hand = append(c.Hand, c.Deck[0])
		c.Deck = c.Deck[1:]
		if showMsg { s.Message = "Drew a card!"; s.MessageTimer = 45 }
	} else if len(c.Hand) >= 5 && showMsg {
		s.Message = "Hand is FULL!"; s.MessageTimer = 60
	} else if showMsg {
		s.Message = "Deck is EMPTY!"; s.MessageTimer = 60
	}
}

func (s *CombatScene) Init(g *core.Game) {
	if s.isInitialized { return }

	var err error
	s.RatImg, _, err = ebitenutil.NewImageFromFile("./assets/rat.png")
	if err != nil {
		fmt.Println("Warning: Could not load rat.png -", err)
	}

	s.Player.Shuffle()
	s.Enemy.Shuffle()

	for i := 0; i < 4; i++ { s.DrawCard(s.Player, false) }
	for i := 0; i < 3; i++ { s.DrawCard(s.Enemy, false) }

	s.TurnNumber = 1 
	s.DamageThreshold = 10 
	s.ActionPhase = -1
	s.IsNight = false
	// НОВОЕ: Если это бой с Гигантским Волком, Дровосек уже ждет на первой линии!
	if s.Enemy.Name == "Giant Wolf" {
		// Ищем Дровосека в колоде или руке, либо создаем его прямо на поле
		woodcutterCard := &models.Card{
			Name:         "Woodcutter",
			Attack:       2,
			HP:           4,
			Cost:         0,
			CanTransform: true,
			AltName:      "Woodcutter (Feral Hound)",
			AltAttack:    5,
			AltHP:        4,
		}
		s.Player.Lanes[0] = woodcutterCard // Ставим на 1-ю линию
	}
	s.isInitialized = true
	s.TutorialStep = g.State.TutorialStep
}

func (s *CombatScene) Update(g *core.Game) error {
	s.Init(g)

	if s.MessageTimer > 0 {
		s.MessageTimer--
		if s.MessageTimer <= 0 { s.Message = "" }
	}

	// ==========================================
	// ФАЗА 11: МЕНЮ КУБОВ СМЕРТИ (Вызывается по C)
	// ==========================================
	if s.ActionPhase == 11 {
		if inpututil.IsKeyJustPressed(ebiten.KeyX) || inpututil.IsKeyJustPressed(ebiten.KeyC) {
			s.ActionPhase = -1 // Закрыть меню
			return nil
		}
		if len(s.Player.Hand) < 5 {
			if inpututil.IsKeyJustPressed(ebiten.Key1) {
				if s.consumeCube(1) {
					s.Player.Hand = append(s.Player.Hand, &models.Card{Name: "Bone Dart", IsSpell: true, Effect: "Damage2", Cost: 0})
					s.Message = "Crafted: Bone Dart!"; s.MessageTimer = 60; s.ActionPhase = -1
				} else { s.Message = "Need Cube lvl 1+"; s.MessageTimer = 45 }
			} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
				if s.consumeCube(2) {
					s.Player.Hand = append(s.Player.Hand, &models.Card{Name: "Plague Rat", Attack: 1, HP: 1, Cost: 0})
					s.Message = "Crafted: Plague Rat!"; s.MessageTimer = 60; s.ActionPhase = -1
				} else { s.Message = "Need Cube lvl 2+"; s.MessageTimer = 45 }
			} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
				if s.consumeCube(3) {
					s.Player.Hand = append(s.Player.Hand, &models.Card{Name: "Bone Wall", Attack: 0, HP: 4, Cost: 0})
					s.Message = "Crafted: Bone Wall!"; s.MessageTimer = 60; s.ActionPhase = -1
				} else { s.Message = "Need Cube lvl 3"; s.MessageTimer = 45 }
			}
		} else {
			if inpututil.IsKeyJustPressed(ebiten.Key1) || inpututil.IsKeyJustPressed(ebiten.Key2) || inpututil.IsKeyJustPressed(ebiten.Key3) {
				s.Message = "Hand is FULL!"; s.MessageTimer = 60
			}
		}
		return nil
	}

	// ФАЗА ЭВОЛЮЦИИ
	if s.ActionPhase == 10 {
		selectedIdx := -1
		if inpututil.IsKeyJustPressed(ebiten.Key1) && len(s.UpgradeOptions) >= 1 { selectedIdx = 0 }
		if inpututil.IsKeyJustPressed(ebiten.Key2) && len(s.UpgradeOptions) >= 2 { selectedIdx = 1 }
		if inpututil.IsKeyJustPressed(ebiten.Key3) && len(s.UpgradeOptions) >= 3 { selectedIdx = 2 }

		if selectedIdx != -1 {
			card := s.UpgradeOptions[selectedIdx]
			_, newName, newAtk, newHP, newCost := getUpgradeInfo(card.Name)
			card.Name, card.Attack, card.HP, card.Cost = newName, newAtk, newHP, newCost

			s.DamageProgress = 0
			s.DamageThreshold += 5 
			rand.Shuffle(len(s.Player.Deck), func(i, j int) { s.Player.Deck[i], s.Player.Deck[j] = s.Player.Deck[j], s.Player.Deck[i] })

			s.Message = "Card Evolved!"
			s.MessageTimer = 60
			s.ActionPhase = -1 
		}
		return nil 
	}

	// ЛОГИКА АНИМАЦИИ И РАЗРЕШЕНИЯ БОЯ
	if s.ActionPhase != -1 && s.ActionPhase != 10 && s.ActionPhase != 11 {
		if s.DelayTimer > 0 { s.DelayTimer--; return nil }

		// Атака игрока
		if s.ActionPhase >= 0 && s.ActionPhase < 4 {
			lane := s.ActionPhase
			if s.Player.Lanes[lane] != nil && s.Player.Lanes[lane].Attack > 0 {
				dmg := s.Player.Lanes[lane].Attack
				s.DamageProgress += dmg 

				if s.Enemy.Lanes[lane] != nil {
					s.Enemy.Lanes[lane].HP -= dmg
					if s.Enemy.Lanes[lane].HP <= 0 { 
						s.Enemy.Lanes[lane] = nil 
						s.AddCubeProgress() // Враг умер -> Качаем куб!
					}
				} else {
					s.Enemy.HP -= dmg 
				}

				// НОВОЕ: Если это Безглавый Петух, он умирает сразу после атаки!
				if s.Player.Lanes[lane].Name == "Headless Rooster" {
					s.Player.Lanes[lane] = nil
					s.Message = "Headless Rooster falls after attack!"
					s.MessageTimer = 45
					// Опционально: можно дать Кровь или Кость за его окончательную смерть
					// s.Player.Blood++ 
				}

				s.DelayTimer = 25 
			} else { s.DelayTimer = 5 }
			s.ActionPhase++
			return nil
		}

		// Атака врага
		if s.ActionPhase >= 4 && s.ActionPhase < 8 {
			lane := s.ActionPhase - 4
			if s.Enemy.Lanes[lane] != nil && s.Enemy.Lanes[lane].Attack > 0 {
				dmg := s.Enemy.Lanes[lane].Attack
				s.DamageProgress += dmg 

				if s.Player.Lanes[lane] != nil {
					s.Player.Lanes[lane].HP -= dmg 
					if s.Player.Lanes[lane].HP <= 0 {
						card := s.Player.Lanes[lane]
						if card.TransformOnDeath != "" && !card.IsTempCopy {
							// НОВОЕ: Превращается в версию «без головы» прямо в руку!
							s.Player.Hand = append(s.Player.Hand, &models.Card{
								Name:       card.TransformOnDeath, // "Headless Rooster"
								Attack:     3, // Увеличена атака до 3!
								HP:         card.HP,
								Cost:       card.Cost,
								IsTempCopy: true,
							})
							s.Message = card.Name + " loses its head! (3 ATK)"
							s.MessageTimer = 60
						} else {
							s.AddCubeProgress()
						}
						s.Player.Lanes[lane] = nil
					}
				} else {
					s.Player.HP -= dmg 
				}
				s.DelayTimer = 25
			} else { s.DelayTimer = 5 }
			s.ActionPhase++
			return nil
		}

		// ПРОДВИНУТЫЙ ИИ С УЧЕТОМ ПАТТЕРНОВ БОССОВ
		if s.ActionPhase == 8 {
			if s.Enemy.HP > 0 {
				// Добираем карты
				for s.Enemy.AP > 0 && len(s.Enemy.Hand) < 5 && len(s.Enemy.Deck) > 0 {
					s.DrawCard(s.Enemy, false)
					s.Enemy.AP--
				}
				
				// Выбираем паттерн поведения
				if s.Enemy.AIPattern == "hunter" {
					// --- ТАКТИКА ВОЛКА: ОХОТНИК ---
					// Волк ищет самую слабую цель у игрока (или пустую линию для прямого удара)
					sort.Slice(s.Enemy.Hand, func(i, j int) bool { return s.Enemy.Hand[i].Attack > s.Enemy.Hand[j].Attack }) // Сначала самые сильные атаки
					
					for i := 0; i < len(s.Enemy.Hand); i++ {
						card := s.Enemy.Hand[i]
						if s.Enemy.Blood >= card.Cost {
							// Ищем идеальную линию для убийства
							targetLane := -1
							lowestHP := 999
							
							// Сначала ищем линии с живыми существами игрока, у которых мало HP (легкая добыча)
							for j := 0; j < 4; j++ {
								if s.Enemy.Lanes[j] == nil { // У Волка линия свободна
									playerCard := s.Player.Lanes[j]
									if playerCard != nil && playerCard.HP < lowestHP {
										lowestHP = playerCard.HP
										targetLane = j
									}
								}
							}
							
							// Если легких целей нет, занимаем любую свободную линию
							if targetLane == -1 {
								for j := 0; j < 4; j++ {
									if s.Enemy.Lanes[j] == nil {
										targetLane = j
										break
									}
								}
							}

							if targetLane != -1 {
								s.Enemy.Lanes[targetLane] = card
								s.Enemy.Blood -= card.Cost
								s.Enemy.Hand = append(s.Enemy.Hand[:i], s.Enemy.Hand[i+1:]...)
								s.Message = "The Wolf stalks its prey!"
								s.MessageTimer = 60
								break
							}
						}
					}
				} else {
					// --- ТАКТИКА ЧЕРВЯ / СТАНДАРТ (Роиться / Заполнять) ---
					sort.Slice(s.Enemy.Hand, func(i, j int) bool { return s.Enemy.Hand[i].Cost > s.Enemy.Hand[j].Cost })
					for i := 0; i < len(s.Enemy.Hand); i++ {
						if s.Enemy.Hand[i].Cost <= s.Enemy.Blood {
							for j := 0; j < 4; j++ {
								if s.Enemy.Lanes[j] == nil {
									s.Enemy.Lanes[j] = s.Enemy.Hand[i]
									s.Enemy.Blood -= s.Enemy.Hand[i].Cost
									s.Enemy.Hand = append(s.Enemy.Hand[:i], s.Enemy.Hand[i+1:]...)
									break
								}
							}
						}
					}
				}
			}
			s.DelayTimer = 45
			s.ActionPhase++
			return nil
		}

		// Завершение хода
		if s.ActionPhase == 9 {
			if s.Player.HP <= 0 {
				if s.IsDebugMode {
					g.ChangeScene(&DebugMenuScene{State: g.State})
					return nil
				}
				if g.State.PubQuestStep == 1 {
					g.ChangeScene(&BasementScene{State: g.State})
				} else {
					g.ChangeScene(&ForestScene{State: g.State})
				}
				return nil
			} else if s.Enemy.HP <= 0 {
				if s.IsDebugMode {
					g.ChangeScene(&DebugMenuScene{State: g.State})
					return nil
				}
				// НОВОЕ: Если мы победили Оборотня
				if s.Enemy.Name == "Werewolf (The Woodcutter)" {
					s.State.PlayerHP = s.Player.HP
					s.State.PubQuestStep = 8 // Идем дальше к каравану
					g.ChangeScene(&ForestScene{State: s.State}) // Или на сцену Каравана
					return nil
				}
				g.State.PubQuestStep++ 
				if g.State.PubQuestStep == 2 {
					g.ChangeScene(&BasementScene{State: g.State}) 
				} else {
					g.ChangeScene(&ForestScene{State: g.State}) 
				}
				return nil 
			}

			// ==========================================
			// МЕХАНИКА ДНЯ И НОЧИ (d20 каждый ход)
			// ==========================================
			if s.TurnNumber == 1 && s.Enemy.Name == "Giant Wolf" {
				// Скрипт для первого хода с Волком: гарантированная первая ночь
				s.IsNight = true
				s.Message = "Night falls... The Full Moon rises!"
				s.MessageTimer = 90
			} else {
				// Случайный бросок d20 для остальных ходов
				d20 := rand.Intn(20) + 1
				if d20 > 10 {
					s.IsNight = !s.IsNight // Меняем время суток
					if s.IsNight {
						s.Message = fmt.Sprintf("Night falls! (d20: %d)", d20)
					} else {
						s.Message = fmt.Sprintf("Dawn breaks! (d20: %d)", d20)
					}
					s.MessageTimer = 75
				} else {
					// Время не изменилось
					// s.Message = fmt.Sprintf("Time holds steady... (d20: %d)", d20)
					// s.MessageTimer = 45
				}
			}

			// Если наступила ночь, все Дровосеки на поле превращаются в зверей
			if s.IsNight {
				for i := 0; i < 4; i++ {
					c := s.Player.Lanes[i]
					if c != nil && c.CanTransform && !c.IsBeastForm {
						c.IsBeastForm = true
						c.Name = c.AltName
						c.Attack = c.AltAttack
						c.HP = c.AltHP
					}
				}
			}
			// ==========================================
			
			s.DrawCard(s.Player, false); s.Player.AP = s.Player.MaxAP
			s.DrawCard(s.Enemy, false); s.Enemy.AP = s.Enemy.MaxAP
			s.TurnNumber++ 

			if s.DamageProgress >= s.DamageThreshold && s.TutorialStep > 0 {
				s.UpgradeOptions = nil
				var eligibleCards []*models.Card
				for _, c := range s.Player.Deck {
					canUpgrade, _, _, _, _ := getUpgradeInfo(c.Name)
					if canUpgrade { eligibleCards = append(eligibleCards, c) }
				}
				if len(eligibleCards) > 0 {
					rand.Shuffle(len(eligibleCards), func(i, j int) { eligibleCards[i], eligibleCards[j] = eligibleCards[j], eligibleCards[i] })
					limit := 3; if len(eligibleCards) < 3 { limit = len(eligibleCards) }
					for i := 0; i < limit; i++ { s.UpgradeOptions = append(s.UpgradeOptions, eligibleCards[i]) }
					s.ActionPhase = 10 
					return nil
				} else { s.DamageProgress = 0 }
			}
			s.ActionPhase = -1 
			return nil
		}
	}

	// ==========================================
	// УПРАВЛЕНИЕ ИГРОКА 
	// ==========================================
	
	// Открытие меню кубов по клавише C
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		s.ActionPhase = 11
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if s.TurnNumber == 1 {
			s.ActionPhase = 8; s.DelayTimer = 20
		} else {
			s.ActionPhase = 0; s.DelayTimer = 10
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if s.Player.AP > 0 {
			if len(s.Player.Hand) < 5 && len(s.Player.Deck) > 0 {
				s.DrawCard(s.Player, true); s.Player.AP--
			} else { s.DrawCard(s.Player, true) }
		} else {
			s.Message = "0 Actions Left!"; s.MessageTimer = 60
		}
	}

	sacTarget := -1
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) { sacTarget = 0 }
	if inpututil.IsKeyJustPressed(ebiten.KeyW) { sacTarget = 1 }
	if inpututil.IsKeyJustPressed(ebiten.KeyE) { sacTarget = 2 }
	if inpututil.IsKeyJustPressed(ebiten.KeyR) { sacTarget = 3 }

	if sacTarget != -1 {
		if s.Player.Lanes[sacTarget] != nil {
			card := s.Player.Lanes[sacTarget]

			if card.TransformOnDeath != "" && !card.IsTempCopy {
				s.Player.Hand = append(s.Player.Hand, &models.Card{
					Name:       card.TransformOnDeath,
					Attack:     3, // 3 атаки
					HP:         card.HP,
					Cost:       card.Cost,
					IsTempCopy: true,
				})
				s.Message = "Rooster: 'I thought you'd save me! Better the soup pot!'"
				s.MessageTimer = 120
			} else {
				s.AddCubeProgress() // Обычные существа качают куб
			}

			s.Player.Lanes[sacTarget] = nil
			s.Player.Blood++
			s.Message = "+1 Blood, Cube Upgraded!"
			s.MessageTimer = 60
		}
	}

	// ПЕРЕМЕЩЕНИЕ В КАРАВАН (Используем клавишу B для третьей линии, чтобы не конфликтовать с C)
	moveTarget := -1
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) { moveTarget = 0 }
	if inpututil.IsKeyJustPressed(ebiten.KeyX) { moveTarget = 1 }
	if inpututil.IsKeyJustPressed(ebiten.KeyB) { moveTarget = 2 } // Изменено с C на B
	if inpututil.IsKeyJustPressed(ebiten.KeyV) { moveTarget = 3 }

	if moveTarget != -1 {
		if s.Player.Buildings[moveTarget] != nil && s.Player.Buildings[moveTarget].Name == "Caravan" {
			if s.Player.Lanes[moveTarget] == nil {
				moved := false
				if moveTarget > 0 && s.Player.Lanes[moveTarget-1] != nil {
					s.Player.Lanes[moveTarget] = s.Player.Lanes[moveTarget-1]
					s.Player.Lanes[moveTarget-1] = nil
					moved = true
				} else if moveTarget < 3 && s.Player.Lanes[moveTarget+1] != nil {
					s.Player.Lanes[moveTarget] = s.Player.Lanes[moveTarget+1]
					s.Player.Lanes[moveTarget+1] = nil
					moved = true
				}

				if moved {
					s.Message = "Moved to Caravan!"
					s.MessageTimer = 45
				} else {
					s.Message = "No adjacent creature!"
					s.MessageTimer = 45
				}
			} else {
				s.Message = "Caravan is full!"
				s.MessageTimer = 45
			}
		} else {
			s.Message = "No Caravan here!"
			s.MessageTimer = 45
		}
	}

	if len(s.Player.Hand) > 0 {
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) { s.SelectedHandIdx = (s.SelectedHandIdx + 1) % len(s.Player.Hand) }
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			s.SelectedHandIdx--; if s.SelectedHandIdx < 0 { s.SelectedHandIdx = len(s.Player.Hand) - 1 }
		}

		laneTarget := -1
		if inpututil.IsKeyJustPressed(ebiten.Key1) { laneTarget = 0 }
		if inpututil.IsKeyJustPressed(ebiten.Key2) { laneTarget = 1 }
		if inpututil.IsKeyJustPressed(ebiten.Key3) { laneTarget = 2 }
		if inpututil.IsKeyJustPressed(ebiten.Key4) { laneTarget = 3 }

		if laneTarget != -1 {
			card := s.Player.Hand[s.SelectedHandIdx]

			if card.IsSpell {
				spellCast := false
				if card.Effect == "Damage2" {
					if s.Enemy.Lanes[laneTarget] != nil {
						s.Enemy.Lanes[laneTarget].HP -= 2 
						if s.Enemy.Lanes[laneTarget].HP <= 0 {
							s.Enemy.Lanes[laneTarget] = nil
							s.AddCubeProgress() 
						}
						spellCast = true
					} else {
						s.Message = "No enemy in that lane!"; s.MessageTimer = 45
					}
				}
				
				if spellCast {
					s.Player.Hand = append(s.Player.Hand[:s.SelectedHandIdx], s.Player.Hand[s.SelectedHandIdx+1:]...)
					if len(s.Player.Hand) > 0 { if s.SelectedHandIdx >= len(s.Player.Hand) { s.SelectedHandIdx = len(s.Player.Hand) - 1 } } else { s.SelectedHandIdx = 0 }
					s.Message = fmt.Sprintf("Cast %s!", card.Name); s.MessageTimer = 60
				}
			} else if card.IsBuilding {
				if s.Player.Buildings[laneTarget] == nil {
					if s.Player.Blood >= card.Cost {
						s.Player.Buildings[laneTarget] = card
						s.Player.Blood -= card.Cost
						
						s.Player.Hand = append(s.Player.Hand[:s.SelectedHandIdx], s.Player.Hand[s.SelectedHandIdx+1:]...)
						if len(s.Player.Hand) > 0 { if s.SelectedHandIdx >= len(s.Player.Hand) { s.SelectedHandIdx = len(s.Player.Hand) - 1 } } else { s.SelectedHandIdx = 0 }
					} else { s.Message = "Not enough Blood!"; s.MessageTimer = 60 }
				} else { s.Message = "Building already here!"; s.MessageTimer = 45 }
			} else {
				if s.Player.Lanes[laneTarget] == nil {
					if s.Player.Blood >= card.Cost {
						s.Player.Lanes[laneTarget] = card
						s.Player.Blood -= card.Cost
						
						s.Player.Hand = append(s.Player.Hand[:s.SelectedHandIdx], s.Player.Hand[s.SelectedHandIdx+1:]...)
						if len(s.Player.Hand) > 0 { if s.SelectedHandIdx >= len(s.Player.Hand) { s.SelectedHandIdx = len(s.Player.Hand) - 1 } } else { s.SelectedHandIdx = 0 }
					} else { s.Message = "Not enough Blood!"; s.MessageTimer = 60 }
				} else { s.Message = "Lane is full!"; s.MessageTimer = 45 }
			}
		}
	}
	return nil
}

func (s *CombatScene) Draw(screen *ebiten.Image) {
	if s.IsNight {
		screen.Fill(color.RGBA{15, 15, 35, 255}) // Ночь (Темно-синий лунный свет)
	} else {
		screen.Fill(color.RGBA{50, 40, 30, 255}) // День (Теплый закатный свет у домика)
	}
	timeStatus := "☀️ DAY"
	if s.IsNight { timeStatus = "🌙 NIGHT (Beast Form Active)" }
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("=== TURN %d | %s ===", s.TurnNumber, timeStatus), 20, 20)
	
	// ФАЗА КРАФТА КУБОВ (11)
	if s.ActionPhase == 11 {
		screen.Fill(color.RGBA{0, 0, 0, 220})
		ebitenutil.DebugPrintAt(screen, "=== CUBE CRAFTING (Press [C] or [X] to Close) ===", 150, 100)
		ebitenutil.DebugPrintAt(screen, "Spend a Death Cube to craft a card instantly:", 150, 140)
		
		ebitenutil.DebugPrintAt(screen, "[1] Bone Dart (2 DMG Spell)  -> Needs Cube lvl 1+", 180, 180)
		ebitenutil.DebugPrintAt(screen, "[2] Plague Rat (1/1 Rat)     -> Needs Cube lvl 2+", 180, 220)
		ebitenutil.DebugPrintAt(screen, "[3] Bone Wall (0/4 Shield)   -> Needs Cube lvl 3", 180, 260)
		return 
	}

	if s.ActionPhase == 10 {
		screen.Fill(color.RGBA{0, 0, 0, 200})
		ebitenutil.DebugPrintAt(screen, "=== EVOLUTION READY! ===", 240, 100)
		startX := 150
		for i, card := range s.UpgradeOptions {
			cx := startX + i*120
			_, newName, newAtk, newHP, newCost := getUpgradeInfo(card.Name)
			ebitenutil.DrawRect(screen, float64(cx-2), 198, 84, 124, color.RGBA{255, 255, 0, 255})
			ebitenutil.DrawRect(screen, float64(cx), 200, 80, 120, color.RGBA{100, 200, 100, 255})
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[%d]", i+1), cx+35, 180)
			ebitenutil.DebugPrintAt(screen, card.Name, cx+5, 210)
			ebitenutil.DebugPrintAt(screen, "   ->", cx+5, 230)
			ebitenutil.DebugPrintAt(screen, newName, cx+5, 250)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("C:%d -> %d", card.Cost, newCost), cx+5, 280)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d/%d -> %d/%d", card.Attack, card.HP, newAtk, newHP), cx+5, 300)
		}
		return 
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("=== TURN %d ===", s.TurnNumber), 20, 20)
	ebitenutil.DebugPrintAt(screen, "[<, >] Select | [1-4] Play | [Q,W,E,R] Sac | [D] Draw | [C] Cubes", 20, 40)
	
	// Вражеский UI
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s HP: %d / %d", s.Enemy.Name, s.Enemy.HP, s.Enemy.MaxHP), 450, 40)
	
	// Интерфейс Игрока
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s HP: %d / %d", s.Player.Name, s.Player.HP, s.Player.MaxHP), 450, 340)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BLOOD: %d", s.Player.Blood), 450, 360)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ACTIONS: %d", s.Player.AP), 450, 380)

	// ОТРИСОВКА КУБОВ СМЕРТИ СПРАВА ВНИЗУ
	ebitenutil.DebugPrintAt(screen, "CUBES:", 450, 410)
	for i, val := range s.Player.Cubes {
		ebitenutil.DrawRect(screen, float64(450+i*25), 430, 20, 20, color.RGBA{220, 220, 220, 255})
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", val), 450+i*25+6, 432)
	}
	
	if s.TutorialStep > 0 {
		progBarX, progBarY := 220, 10
		fillWidth := int((float64(s.DamageProgress) / float64(s.DamageThreshold)) * 200)
		if fillWidth > 200 { fillWidth = 200 }
		ebitenutil.DrawRect(screen, float64(progBarX), float64(progBarY), 200, 15, color.RGBA{80, 80, 80, 255})
		ebitenutil.DrawRect(screen, float64(progBarX), float64(progBarY), float64(fillWidth), 15, color.RGBA{0, 200, 255, 255}) 
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("EVOLUTION: %d / %d DMG", s.DamageProgress, s.DamageThreshold), progBarX+30, progBarY+20)
	}

	if s.Message != "" { ebitenutil.DebugPrintAt(screen, ">>> "+s.Message+" <<<", 200, 250) }

	enemyHandX, enemyHandY := 20, 60
	for i := 0; i < len(s.Enemy.Hand); i++ {
		ebitenutil.DrawRect(screen, float64(enemyHandX+i*40), float64(enemyHandY), 30, 40, color.RGBA{100, 20, 20, 255})
	}

	startX, startY := 100, 140
	laneWidth, laneHeight, gap := 80, 100, 20

	for i := 0; i < 4; i++ {
		x := startX + i*(laneWidth+gap)
		
		enemyY := startY
		if s.ActionPhase >= 4 && s.ActionPhase < 8 && i == (s.ActionPhase-4) { enemyY += 15 }
		ebitenutil.DrawRect(screen, float64(x), float64(enemyY), float64(laneWidth), float64(laneHeight), color.RGBA{100, 50, 50, 255})
		if s.Enemy.Lanes[i] != nil {
			if s.Enemy.Lanes[i].Name == "Brute" {
				ebitenutil.DrawRect(screen, float64(x+5), float64(enemyY+5), float64(laneWidth-10), float64(laneHeight-10), color.RGBA{150, 0, 150, 255})
			} else {
				ebitenutil.DrawRect(screen, float64(x+5), float64(enemyY+5), float64(laneWidth-10), float64(laneHeight-10), color.RGBA{200, 50, 50, 255})
			}
			ebitenutil.DebugPrintAt(screen, s.Enemy.Lanes[i].Name, x+10, enemyY+20)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ATK:%d HP:%d", s.Enemy.Lanes[i].Attack, s.Enemy.Lanes[i].HP), x+10, enemyY+70)
		}

		playerY := startY + laneHeight + gap
		if s.ActionPhase >= 0 && s.ActionPhase < 4 && i == s.ActionPhase { playerY -= 15 }
		
		// 1. Отрисовка Здания (Фон линии)
		if s.Player.Buildings[i] != nil {
			ebitenutil.DrawRect(screen, float64(x), float64(playerY), float64(laneWidth), float64(laneHeight), color.RGBA{139, 69, 19, 255})
			ebitenutil.DebugPrintAt(screen, s.Player.Buildings[i].Name, x+5, playerY+5)
			
			// Подсказка на перемещение (клавиша B для 3-й линии)
			moveKey := []string{"Z", "X", "B", "V"}[i]
			ebitenutil.DebugPrintAt(screen, "["+moveKey+"] Move In", x+5, playerY+85)
		} else {
			ebitenutil.DrawRect(screen, float64(x), float64(playerY), float64(laneWidth), float64(laneHeight), color.RGBA{50, 100, 50, 255})
		}
		
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Sac:[%s]", []string{"Q", "W", "E", "R"}[i]), x+15, playerY+laneHeight+5)

		if s.Player.Lanes[i] != nil {
			cardH := float64(laneHeight-30)
			
			ebitenutil.DrawRect(screen, float64(x+5), float64(playerY+20), float64(laneWidth-10), cardH, color.RGBA{50, 200, 50, 255})
			
			if s.Player.Lanes[i].Name == "Rat" && s.RatImg != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x+5), float64(playerY+20))
				screen.DrawImage(s.RatImg, op)
			} else if s.Player.Lanes[i].Name == "Bone Wall" {
				ebitenutil.DrawRect(screen, float64(x+5), float64(playerY+5), float64(laneWidth-10), float64(laneHeight-10), color.RGBA{150, 150, 150, 255})
			} else {
				ebitenutil.DrawRect(screen, float64(x+5), float64(playerY+5), float64(laneWidth-10), float64(laneHeight-10), color.RGBA{50, 200, 50, 255})
			}
			ebitenutil.DebugPrintAt(screen, s.Player.Lanes[i].Name, x+10, playerY+20)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ATK:%d HP:%d", s.Player.Lanes[i].Attack, s.Player.Lanes[i].HP), x+10, playerY+70)
		}
	}

	ebitenutil.DebugPrintAt(screen, "DECK:", 20, 380)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d Cards", len(s.Player.Deck)), 20, 400)
	ebitenutil.DrawRect(screen, 20, 420, 40, 50, color.RGBA{150, 150, 150, 255})

	if len(s.Player.Hand) > 0 && s.ActionPhase == -1 {
		handStartX, handStartY := 120, 380
		for i, card := range s.Player.Hand {
			cx := handStartX + i*(60+gap)
			if i == s.SelectedHandIdx {
				ebitenutil.DrawRect(screen, float64(cx-2), float64(handStartY-2), 64, 64, color.RGBA{255, 255, 0, 255})
			}
			if card.IsSpell {
				ebitenutil.DrawRect(screen, float64(cx), float64(handStartY), 60, 60, color.RGBA{150, 100, 200, 255})
			} else if card.IsBuilding {
				ebitenutil.DrawRect(screen, float64(cx), float64(handStartY), 60, 60, color.RGBA{210, 105, 30, 255})
			} else {
				ebitenutil.DrawRect(screen, float64(cx), float64(handStartY), 60, 60, color.RGBA{200, 200, 200, 255})
			}
			ebitenutil.DebugPrintAt(screen, card.Name, cx+5, handStartY+10)
			if card.Cost > 0 { ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Bld:%d", card.Cost), cx+5, handStartY+25) }
		}
	}
}
package scenes

import (
	"image/color"

	"github.com/SaimonPetrikov/bremen/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// --- ПАБ ---
type PubScene struct {
	State *core.GameState

	HeroX float64
	HeroY float64

	isInitialized bool
	HeroImg       *ebiten.Image
	MapImg        *ebiten.Image
}

func (s *PubScene) Init() {
	if s.isInitialized { return }

	s.HeroImg = ebiten.NewImage(32, 32)
	s.HeroImg.Fill(color.RGBA{50, 150, 255, 255})

	s.MapImg = ebiten.NewImage(640, 480)
	// 1. Заливаем весь экран темным "ничто"
	s.MapImg.Fill(color.RGBA{15, 15, 20, 255}) 

	// 2. Рисуем сам Паб по центру (Координаты: X=120, Y=90, Ширина=400, Высота=300)
	ebitenutil.DrawRect(s.MapImg, 120, 90, 400, 300, color.RGBA{80, 40, 20, 255}) // Деревянный пол

	// Барная стойка
	ebitenutil.DrawRect(s.MapImg, 220, 110, 200, 40, color.RGBA{139, 69, 19, 255}) 
	// Хозяин за стойкой
	ebitenutil.DrawRect(s.MapImg, 300, 80, 40, 30, color.RGBA{200, 150, 100, 255}) 

	// Люк в Подвал (Левый верхний угол паба)
	ebitenutil.DrawRect(s.MapImg, 140, 110, 50, 50, color.RGBA{30, 30, 30, 255}) 

	// Выход наружу (Внизу по центру паба)
	ebitenutil.DrawRect(s.MapImg, 280, 370, 80, 20, color.RGBA{150, 150, 150, 255}) 

	// Герой появляется у выхода
	s.HeroX = 300
	s.HeroY = 330

	s.isInitialized = true
}

func (s *PubScene) Update(g *core.Game) error {
	s.Init()

	speed := 4.0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) { s.HeroX -= speed }
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) { s.HeroX += speed }
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) { s.HeroY -= speed }
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) { s.HeroY += speed }

	if inpututil.IsKeyJustPressed(ebiten.KeyI) { g.ChangeScene(&InventoryScene{State: s.State, PreviousScene: s}); return nil }

	// Ограничения стен (теперь герой не может выйти в темную зону!)
	if s.HeroX < 120 { s.HeroX = 120 }
	if s.HeroX > 520-32 { s.HeroX = 520 - 32 }
	if s.HeroY < 90 { s.HeroY = 90 }
	if s.HeroY > 390-32 { s.HeroY = 390 - 32 }

	// Зоны взаимодействия сдвинуты с учетом новых координат
	isNearOwner := s.HeroX > 200 && s.HeroX < 440 && s.HeroY < 160
	isNearBasement := s.HeroX < 200 && s.HeroY < 180
	isNearExit := s.HeroX > 260 && s.HeroX < 340 && s.HeroY > 320

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if isNearExit {
			g.ChangeScene(&HubScene{State: s.State})
		} else if isNearBasement && s.State.PubQuestStep <= 1 {
			g.ChangeScene(&BasementScene{State: s.State})
		} else if isNearOwner {
			if s.State.PubQuestStep == 2 {
				s.State.PubQuestStep = 3
				g.ChangeScene(&MillScene{State: s.State})
			} else if s.State.PubQuestStep == 5 {
				s.State.PubQuestStep = 6
			}
		}
	}

	return nil
}

func (s *PubScene) Draw(screen *ebiten.Image) {
	s.Init()

	screen.DrawImage(s.MapImg, &ebiten.DrawImageOptions{})
	
	heroOp := &ebiten.DrawImageOptions{}
	heroOp.GeoM.Translate(s.HeroX, s.HeroY)
	screen.DrawImage(s.HeroImg, heroOp)

	// UI Фон
	ebitenutil.DrawRect(screen, 0, 0, 640, 60, color.RGBA{0, 0, 0, 180})
	ebitenutil.DebugPrintAt(screen, "=== THE PUB ===", 10, 10)

	isNearOwner := s.HeroX > 200 && s.HeroX < 440 && s.HeroY < 160
	isNearBasement := s.HeroX < 200 && s.HeroY < 180
	isNearExit := s.HeroX > 260 && s.HeroX < 340 && s.HeroY > 320

	if isNearOwner {
		switch s.State.PubQuestStep {
		case 0, 1:
			ebitenutil.DebugPrintAt(screen, "Owner: Someone is stealing my food! Check the basement!", 10, 30)
		case 2:
			ebitenutil.DebugPrintAt(screen, "Owner: You killed the giant worm? Amazing!", 10, 30)
			ebitenutil.DebugPrintAt(screen, "Tavern Owner: Now go to the old Mill and ask if the grain was delivered!", 10, 40)
		ebitenutil.DebugPrintAt(screen, "[SPACE] Head out to the Mill", 10, 80)
		case 3, 4:
			ebitenutil.DebugPrintAt(screen, "Owner: Please find my beer caravan in the Forest!", 10, 30)
		case 5:
			ebitenutil.DebugPrintAt(screen, "Owner: You saved the beer AND killed a Skeleton Lord?!", 10, 30)
			ebitenutil.DebugPrintAt(screen, "[SPACE] Take Final Award", 10, 45)
		case 6:
			ebitenutil.DebugPrintAt(screen, "Owner: Have a drink, true hero!", 10, 30)
		}
	} else if isNearBasement {
		if s.State.PubQuestStep <= 1 {
			ebitenutil.DebugPrintAt(screen, "[SPACE] Go down to Basement", 10, 30)
		} else {
			ebitenutil.DebugPrintAt(screen, "The basement is locked. You already cleared it.", 10, 30)
		}
	} else if isNearExit {
		ebitenutil.DebugPrintAt(screen, "[SPACE] Leave to Hub", 10, 30)
	}
}
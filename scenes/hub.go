package scenes

import (
	"image/color"

	"github.com/SaimonPetrikov/bremen/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// --- ХАБ ---
type HubScene struct {
	State *core.GameState

	// Координаты героя
	HeroX float64
	HeroY float64

	// Графика
	isInitialized bool
	HeroImg       *ebiten.Image
	MapImg        *ebiten.Image
}

func (s *HubScene) Init() {
	if s.isInitialized {
		return
	}

	// 1. Создаем пиксельного героя (Синий квадрат 32x32)
	// В будущем замените на: s.HeroImg, _, _ = ebitenutil.NewImageFromFile("hero.png")
	s.HeroImg, _, _ = ebitenutil.NewImageFromFile("./assets/hero.png")

	// 2. Создаем карту (Зеленая трава и серые здания)
	s.MapImg = ebiten.NewImage(640, 480)
	s.MapImg.Fill(color.RGBA{34, 139, 34, 255}) // Трава

	// Рисуем Таверну (Паб) на карте
	ebitenutil.DrawRect(s.MapImg, 400, 100, 120, 100, color.RGBA{139, 69, 19, 255}) // Коричневое здание
	ebitenutil.DrawRect(s.MapImg, 440, 160, 40, 40, color.RGBA{100, 50, 0, 255})   // Дверь
	
	// Рисуем выход в Лес (Дорога вниз)
	ebitenutil.DrawRect(s.MapImg, 280, 400, 80, 80, color.RGBA{100, 100, 100, 255}) // Каменная тропа

	// Стартовая позиция героя (в центре)
	s.HeroX = 300
	s.HeroY = 240

	s.isInitialized = true
}

func (s *HubScene) Update(g *core.Game) error {
	s.Init()

	if s.State.PlayerHP <= 0 { s.State.PlayerHP = 10 }

	// ==========================================
	// ПЕРЕМЕЩЕНИЕ ГЕРОЯ (WASD или Стрелки)
	// ==========================================
	speed := 4.0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) { s.HeroX -= speed }
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) { s.HeroX += speed }
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) { s.HeroY -= speed }
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) { s.HeroY += speed }

	// Ограничение выхода за экран
	if s.HeroX < 0 { s.HeroX = 0 }
	if s.HeroX > 640-32 { s.HeroX = 640 - 32 }
	if s.HeroY < 0 { s.HeroY = 0 }
	if s.HeroY > 480-32 { s.HeroY = 480 - 32 }

	// ==========================================
	// ВЗАИМОДЕЙСТВИЕ (Заходим в здания)
	// ==========================================
	// Старые кнопки (на всякий случай)
	if inpututil.IsKeyJustPressed(ebiten.Key3) { g.ChangeScene(&PubScene{State: s.State}) }
	if inpututil.IsKeyJustPressed(ebiten.Key4) && s.State.PubQuestStep >= 3 { g.ChangeScene(&ForestScene{State: s.State}) }
	if inpututil.IsKeyJustPressed(ebiten.KeyD) { g.ChangeScene(&DebugMenuScene{State: s.State}) }

	// Новое взаимодействие: Нажать SPACE рядом со зданием
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		// Проверяем, стоит ли герой возле дверей Паба (координаты X:400-520, Y:100-200)
		if s.HeroX > 380 && s.HeroX < 540 && s.HeroY > 80 && s.HeroY < 220 {
			g.ChangeScene(&PubScene{State: s.State})
		}
		
		// Проверяем, стоит ли герой на тропе в Лес (X:280-360, Y:400-480)
		if s.State.PubQuestStep >= 3 && s.HeroX > 260 && s.HeroX < 380 && s.HeroY > 380 {
			g.ChangeScene(&ForestScene{State: s.State})
		}
	}

	return nil
}

func (s *HubScene) Draw(screen *ebiten.Image) {
	s.Init()
	// 1. Отрисовываем фон (Карту)
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(s.MapImg, op)

	// 2. Отрисовываем героя
	heroOp := &ebiten.DrawImageOptions{}
	heroOp.GeoM.Translate(s.HeroX, s.HeroY)
	screen.DrawImage(s.HeroImg, heroOp)

	// 3. Отрисовываем UI поверх всего
	ebitenutil.DrawRect(screen, 0, 0, 640, 40, color.RGBA{0, 0, 0, 150}) // Полупрозрачная плашка для текста
	ebitenutil.DebugPrintAt(screen, "=== VILLAGE HUB ===", 10, 10)
	
	// Подсказки меняются в зависимости от того, где стоит игрок
	isNearPub := s.HeroX > 380 && s.HeroX < 540 && s.HeroY > 80 && s.HeroY < 220
	isNearForest := s.State.PubQuestStep >= 3 && s.HeroX > 260 && s.HeroX < 380 && s.HeroY > 380

	if isNearPub {
		ebitenutil.DebugPrintAt(screen, "[SPACE] Enter Pub", 420, 80)
	} else if isNearForest {
		ebitenutil.DebugPrintAt(screen, "[SPACE] Enter Dark Forest", 260, 380)
	} else {
		ebitenutil.DebugPrintAt(screen, "Walk with WASD/Arrows. Press SPACE to interact.", 200, 10)
	}

	if s.State.PlayerHP < 10 {
		ebitenutil.DebugPrintAt(screen, "You feel rested. HP restored to 10.", 400, 10)
	}
}
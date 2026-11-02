package scenes

import (
	"image/color"

	"github.com/SaimonPetrikov/bremen/core"
	"github.com/SaimonPetrikov/bremen/models"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// --- ПОДВАЛ (Tutorial 1: Blood Mechanic) ---
type BasementScene struct {
	State *core.GameState

	HeroX float64
	HeroY float64

	isInitialized bool
	HeroImg       *ebiten.Image
	MapImg        *ebiten.Image
}

func (s *BasementScene) Init() {
	if s.isInitialized { return }

	s.HeroImg = ebiten.NewImage(32, 32)
	s.HeroImg.Fill(color.RGBA{50, 150, 255, 255})

	s.MapImg = ebiten.NewImage(640, 480)
	s.MapImg.Fill(color.RGBA{10, 10, 15, 255}) // Очень темная пустота вокруг

	// Пол подвала (по центру)
	ebitenutil.DrawRect(s.MapImg, 150, 100, 340, 280, color.RGBA{60, 60, 60, 255}) // Темно-серый камень

	// Лестница наверх (в Паб) - в левом верхнем углу подвала
	ebitenutil.DrawRect(s.MapImg, 170, 120, 50, 50, color.RGBA{100, 80, 60, 255}) 

	// Герой появляется у лестницы
	s.HeroX = 180
	s.HeroY = 180

	s.isInitialized = true
}

func (s *BasementScene) Update(g *core.Game) error {
	s.Init()

	speed := 4.0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		s.HeroX -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		s.HeroX += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		s.HeroY -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		s.HeroY += speed
	}

	// Ограничения стен подвала
	if s.HeroX < 150 {
		s.HeroX = 150
	}
	if s.HeroX > 490-32 {
		s.HeroX = 490 - 32
	}
	if s.HeroY < 100 {
		s.HeroY = 100
	}
	if s.HeroY > 380-32 {
		s.HeroY = 380 - 32
	}

	isNearStairs := s.HeroX < 240 && s.HeroY < 190
	isNearWorm := s.HeroX > 320 && s.HeroY > 160 && s.HeroY < 280

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if isNearStairs {
			g.ChangeScene(&PubScene{State: s.State})
		} else if isNearWorm && s.State.PubQuestStep == 0 {
			s.State.PubQuestStep = 1
			s.State.PlayerHP = 10

			s.State.PlayerDeck = models.InitDefaultPlayerDeck()
			player := models.NewPlayer(s.State.PlayerDeck)
			
			worm := models.NewWorm()

			g.ChangeScene(&CombatScene{State: s.State, Player: player, Enemy: worm})
		}
	}
	return nil
}

func (s *BasementScene) Draw(screen *ebiten.Image) {
	s.Init()

	screen.DrawImage(s.MapImg, &ebiten.DrawImageOptions{})

	// Динамическая отрисовка Червя (только если квест еще не сдан)
	if s.State.PubQuestStep == 0 {
		ebitenutil.DrawRect(screen, 380, 180, 60, 60, color.RGBA{180, 50, 100, 255}) // Розово-красный червь
		ebitenutil.DebugPrintAt(screen, "Giant Worm", 375, 160)
	}

	heroOp := &ebiten.DrawImageOptions{}
	heroOp.GeoM.Translate(s.HeroX, s.HeroY)
	screen.DrawImage(s.HeroImg, heroOp)

	// UI Фон
	ebitenutil.DrawRect(screen, 0, 0, 640, 60, color.RGBA{0, 0, 0, 180})
	ebitenutil.DebugPrintAt(screen, "=== PUB BASEMENT ===", 10, 10)

	isNearStairs := s.HeroX < 240 && s.HeroY < 190
	isNearWorm := s.HeroX > 320 && s.HeroY > 160 && s.HeroY < 280

	if isNearStairs {
		ebitenutil.DebugPrintAt(screen, "[SPACE] Go up to Pub", 10, 30)
	} else if isNearWorm && s.State.PubQuestStep == 0 {
		ebitenutil.DebugPrintAt(screen, "TUTORIAL: Sacrifice basic Rats to play Fat Rats!", 10, 30)
		ebitenutil.DebugPrintAt(screen, "[SPACE] Fight the Worm with the Rat Deck!", 10, 45)
	} else if s.State.PubQuestStep >= 2 {
		ebitenutil.DebugPrintAt(screen, "The basement is clear. The rats are now your loyal army.", 10, 30)
	}
}
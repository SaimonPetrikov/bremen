package core

import (
	"github.com/SaimonPetrikov/bremen/models"
	"github.com/hajimehoshi/ebiten/v2"
)

// GameState хранит прогресс квестов и инвентарь (колоду)
type GameState struct {
	TutorialStep int
	HasRatDeck   bool
	PubQuestStep int
	
	// НОВЫЕ ПОЛЯ: Постоянный прогресс игрока
	PlayerHP   int
	PlayerDeck []models.PlayingCard
}

// Scene — базовый интерфейс для всех сцен
type Scene interface {
	Update(g *Game) error
	Draw(screen *ebiten.Image)
}

// Game — главная структура игры
type Game struct {
	CurrentScene Scene
	State        *GameState
}

func (g *Game) Update() error {
	return g.CurrentScene.Update(g)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.CurrentScene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func (g *Game) ChangeScene(s Scene) {
	g.CurrentScene = s
}
package main

import (
	"log"

	"github.com/SaimonPetrikov/bremen/core"
	"github.com/SaimonPetrikov/bremen/scenes"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	state := &core.GameState{}
	game := &core.Game{
		CurrentScene: &scenes.MainMenuScene{State: state},
		State:        state,
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("2D Pixel RPG Card Battler Demo")
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
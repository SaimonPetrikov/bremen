package scenes

import (
	"github.com/SaimonPetrikov/bremen/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MainMenuScene struct {
	State *core.GameState
}

func (s *MainMenuScene) Update(g *core.Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.ChangeScene(&HubScene{State: s.State})
	}
	return nil
}

func (s *MainMenuScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "=== PIXEL RPG CARD BATTLER ===", 220, 200)
	ebitenutil.DebugPrintAt(screen, "Press ENTER to Start", 250, 240)
}
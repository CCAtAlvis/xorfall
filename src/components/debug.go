package components

import (
	"fmt"

	"github.com/CCAtAlvis/xorfall/src/configs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type DebugComponent struct {
	BaseComponent
}

func NewDebugComponent() *DebugComponent {
	return &DebugComponent{
		BaseComponent: NewBaseComponent("debug", 720, 480),
	}
}

func (d *DebugComponent) Update(gameTime *configs.GameTimeManager) {}

func (d *DebugComponent) Render() {
	d.Begin()

	rl.DrawRectangleLinesEx(rl.Rectangle{X: 0, Y: 0, Width: 720, Height: 480}, 2, rl.Green)
	// rl.DrawText("Hello Virtual World!", 100, 100, 40, rl.White)
	rl.DrawText(fmt.Sprintf("Screen Width: %d", rl.GetScreenWidth()), 100, 150, 40, rl.White)
	rl.DrawText(fmt.Sprintf("Screen Height: %d", rl.GetScreenHeight()), 100, 200, 40, rl.White)
	rl.DrawText(fmt.Sprintf("Render Width: %d", rl.GetRenderWidth()), 100, 250, 40, rl.White)
	rl.DrawText(fmt.Sprintf("Render Height: %d", rl.GetRenderHeight()), 100, 300, 40, rl.White)

	// rl.DrawRectangle(300, 200, 200, 100, rl.Red)
	// rl.DrawRectangle(700, 200, 100, 100, rl.Red)

	d.End()
}

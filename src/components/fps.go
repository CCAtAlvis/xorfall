package components

import (
	"fmt"

	"github.com/CCAtAlvis/xorfall/src/configs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type FPSComponent struct {
	BaseComponent
	fpsHistory []int32
	averageFPS int32
}

var bottom int32 = 120
var maxBarHeight int32 = 50

func NewFPSComponent() *FPSComponent {
	return &FPSComponent{
		BaseComponent: NewBaseComponent("fps", 120, 120),
		fpsHistory:    make([]int32, 116),
		averageFPS:    0,
	}
}

func (f *FPSComponent) Update(gameTime *configs.GameTimeManager) {
	f.fpsHistory = append(f.fpsHistory, rl.GetFPS())
	if len(f.fpsHistory) > 116 {
		f.fpsHistory = f.fpsHistory[1:]
	}

	f.averageFPS = 0
	for _, fps := range f.fpsHistory {
		f.averageFPS += fps
	}
	f.averageFPS /= int32(len(f.fpsHistory))
}

func (f *FPSComponent) Render() {
	f.Begin()

	rl.DrawRectangleLinesEx(rl.Rectangle{X: 0, Y: 0, Width: 120, Height: 120}, 2, rl.Red)
	rl.DrawText(
		fmt.Sprintf("FPS: %d", rl.GetFPS()),
		12, 12, 20, rl.Green,
	)
	rl.DrawText(
		fmt.Sprintf("Avg: %d", f.averageFPS),
		12, 40, 20, rl.Green,
	)

	for i, fps := range f.fpsHistory {
		targetFPS := configs.TargetFPS
		percentOfTargetFPS := float32(fps) / float32(targetFPS)
		if percentOfTargetFPS > 1.0 {
			percentOfTargetFPS = 1.0
		}
		barHeight := int32(percentOfTargetFPS * float32(maxBarHeight))
		top := bottom - 2 - barHeight

		rl.DrawRectangle(int32(i+2), int32(top), 1, barHeight, rl.Green)
	}

	f.End()
}

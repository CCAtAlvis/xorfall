package components

import (
	"github.com/CCAtAlvis/xorfall/src/configs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GridComponent struct {
	BaseComponent
	rows        int32
	cols        int32
	cellSize    int32
	cellPadding int32
	cellColor   rl.Color
}

func (g *GridComponent) Update(gameTime *configs.GameTimeManager) {
}

func (g *GridComponent) Render() {
	g.Begin()

	for i := int32(0); i < g.rows; i++ {
		for j := int32(0); j < g.cols; j++ {
			if j%2 == 0 {
				rl.DrawRectangle(j*(g.cellSize+g.cellPadding), i*(g.cellSize+g.cellPadding), g.cellSize, g.cellSize, g.cellColor)
			} else {
				rl.DrawRectangleLines(j*(g.cellSize+g.cellPadding), i*(g.cellSize+g.cellPadding), g.cellSize, g.cellSize, g.cellColor)
			}
		}
	}

	g.End()
}

func NewGridComponent() *GridComponent {
	return &GridComponent{
		BaseComponent: NewBaseComponent("grid", 1280, 720),
		rows:          20,
		cols:          8,
		cellSize:      32,
		cellPadding:   4,
		cellColor:     rl.Green,
	}
}

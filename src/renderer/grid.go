package renderer

import rl "github.com/gen2brain/raylib-go/raylib"

type Grid struct {
	rows        int32
	cols        int32
	cellSize    int32
	cellPadding int32
	cellColor   rl.Color
}

func (g *Grid) Draw() {
	for i := int32(0); i < g.rows; i++ {
		for j := int32(0); j < g.cols; j++ {
			if j%2 == 0 {
				rl.DrawRectangle(j*(g.cellSize+g.cellPadding), i*(g.cellSize+g.cellPadding), g.cellSize, g.cellSize, g.cellColor)
			} else {
				rl.DrawRectangleLines(j*(g.cellSize+g.cellPadding), i*(g.cellSize+g.cellPadding), g.cellSize, g.cellSize, g.cellColor)
			}
		}
	}
}

func NewGrid(rows int32, cols int32, cellSize int32, cellPadding int32, cellColor rl.Color) *Grid {
	return &Grid{
		rows:        rows,
		cols:        cols,
		cellSize:    cellSize,
		cellPadding: cellPadding,
		cellColor:   cellColor,
	}
}

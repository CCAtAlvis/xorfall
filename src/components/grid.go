package components

import (
	"fmt"
	"slices"

	"github.com/CCAtAlvis/xorfall/src/configs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const ColumnCount = 8
const RowCount = 20
const InitialRowCount = 3
const OffsetX = int32(10)
const OffsetY = int32(32)

type Row struct {
	bits uint8
}

func (r *Row) IsBitSet(col int) bool {
	return r.bits&(1<<col) != 0
}

func (r *Row) SetBit(col int) {
	r.bits |= (1 << col)
}

func (r *Row) ClearBit(col int) {
	r.bits &^= (1 << col)
}

func (r *Row) FlipBit(col int) {
	r.bits ^= (1 << col)
}

func (r *Row) IsAllOn() bool {
	return r.bits == 0b11111111
}

func (r *Row) IsAllOff() bool {
	return r.bits == 0
}

func (r *Row) Randomize() {
	r.bits = uint8(rl.GetRandomValue(0, 255))
}

func GenerateRandomRow() Row {
	return Row{bits: uint8(rl.GetRandomValue(0, 255))}
}

type Board struct {
	rows            []Row
	currentRowIndex int

	validRowStates          []uint8
	lastRowSpawnTime        float32
	currentRowSpawnInterval float32
}

type GridComponent struct {
	BaseComponent
	board Board

	cellSize    int32
	cellPadding int32
	cellColor   rl.Color

	gridBorderThickness int32
	gridBorderColor     rl.Color
}

func (g *GridComponent) Update(gameTime *configs.GameTimeManager) {
	if rl.IsKeyPressed(rl.KeyR) {
		// randomizeRow := int(rl.GetRandomValue(0, int32(g.board.currentRowIndex)))
		g.board.rows[0].Randomize()
	}

	row := g.board.rows[0]
	if slices.Contains(g.board.validRowStates, row.bits) {
		g.board.rows = g.board.rows[1 : g.board.currentRowIndex+1]
		g.board.currentRowIndex -= 1
		if g.board.currentRowIndex < 0 {
			// TODO: Handle this case
			g.board.currentRowIndex = 0
		}
		fmt.Println("row removed", g.board.currentRowIndex, len(g.board.rows))
	}

	g.board.lastRowSpawnTime += configs.GameTime.Delta
	if g.board.lastRowSpawnTime >= g.board.currentRowSpawnInterval {
		g.board.rows = append(g.board.rows, GenerateRandomRow())
		g.board.currentRowIndex++
		if g.board.currentRowIndex >= RowCount {
			// TODO: Handle this case
			g.board.currentRowIndex = RowCount - 1
		}
		g.board.lastRowSpawnTime = 0
		fmt.Println("row spawned", g.board.currentRowIndex, len(g.board.rows))
	}

	if rl.IsKeyPressed(rl.KeySpace) {
		g.board.rows[0].bits = 0b11111111
	}
}

func (g *GridComponent) Render() {
	g.Begin()
	rl.ClearBackground(rl.Black)

	// draw grid border
	// draw grid horizontal lines
	for i := 0; i <= RowCount; i++ {
		startX := OffsetX
		startY := OffsetY + int32(i)*(g.cellSize+g.cellPadding)
		startPos := rl.Vector2{X: float32(startX), Y: float32(startY)}

		endX := float32(startX + (ColumnCount * (g.cellSize + g.cellPadding)))
		endY := startY
		endPos := rl.Vector2{X: float32(endX), Y: float32(endY)}

		rl.DrawLineEx(startPos, endPos, float32(g.gridBorderThickness), g.gridBorderColor)
	}

	// draw grid vertical lines
	for i := 0; i <= ColumnCount; i++ {

		startX := OffsetX + int32(i)*(g.cellSize+g.cellPadding)
		startY := OffsetY
		startPos := rl.Vector2{X: float32(startX), Y: float32(startY)}

		endX := startX
		endY := startY + (RowCount * (g.cellSize + g.cellPadding))
		endPos := rl.Vector2{X: float32(endX), Y: float32(endY)}

		rl.DrawLineEx(startPos, endPos, float32(g.gridBorderThickness), g.gridBorderColor)
	}

	// draw cells
	// rows are inverted
	for rid := 0; rid <= g.board.currentRowIndex; rid++ {
		row := g.board.rows[rid]
		i := RowCount - g.board.currentRowIndex + rid - 1

		for j := range 8 {
			cellX := OffsetX + int32(j)*(g.cellSize+g.cellPadding) + (g.gridBorderThickness + g.gridBorderThickness/2)
			cellY := OffsetY + int32(i)*(g.cellSize+g.cellPadding) + (g.gridBorderThickness + g.gridBorderThickness/2)

			if row.IsBitSet(j) {
				rl.DrawRectangle(cellX, cellY, g.cellSize, g.cellSize, g.cellColor)
			} else {
				rl.DrawRectangleLinesEx(rl.Rectangle{X: float32(cellX), Y: float32(cellY), Width: float32(g.cellSize), Height: float32(g.cellSize)}, 4, g.cellColor)
			}
		}
	}

	g.End()
}

func NewGridComponent() *GridComponent {
	cellSize := int32(28)
	cellPadding := int32(8)
	cellColor := rl.Color{R: 0, G: 255, B: 0, A: 255}

	gridBorderThickness := int32(2)
	gridBorderColor := rl.Color{R: 0, G: 255, B: 0, A: 150}

	gridWidth := (ColumnCount * cellSize) + ((ColumnCount - 1) * cellPadding) + (2 * (gridBorderThickness + gridBorderThickness/2)) + (2 * OffsetX)
	gridHeight := (RowCount * cellSize) + ((RowCount - 1) * cellPadding) + (2 * (gridBorderThickness + gridBorderThickness/2)) + (2 * OffsetY)

	board := Board{
		rows:                    make([]Row, 0),
		currentRowIndex:         InitialRowCount - 1,
		validRowStates:          make([]uint8, 0),
		currentRowSpawnInterval: 1.0,
	}
	for range InitialRowCount {
		board.rows = append(board.rows, GenerateRandomRow())
	}
	board.validRowStates = append(board.validRowStates, 0b11111111)
	board.validRowStates = append(board.validRowStates, 0b0)

	return &GridComponent{
		BaseComponent:       NewBaseComponent("grid", gridWidth, gridHeight),
		cellSize:            cellSize,
		cellPadding:         cellPadding,
		cellColor:           cellColor,
		gridBorderThickness: gridBorderThickness,
		gridBorderColor:     gridBorderColor,
		board:               board,
	}
}

package grid

import (
	"fmt"
	"slices"

	"github.com/CCAtAlvis/xorfall/src/components"
	"github.com/CCAtAlvis/xorfall/src/configs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const ColumnCount = 8
const RowCount = 20
const InitialRowCount = 3
const OffsetX = int32(10)
const OffsetY = int32(32)

type GridComponent struct {
	components.BaseComponent
	board       *Board
	maskManager *MaskManager

	cellSize    int32
	cellPadding int32
	cellColor   rl.Color

	gridBorderThickness int32
	gridBorderColor     rl.Color
}

func (g *GridComponent) Update(gameTime *configs.GameTimeManager) {
	if rl.IsKeyPressed(rl.KeyQ) {
		g.Reset()
		configs.GameState().SetGameState(configs.GameStatePlaying)
		return
	}

	if configs.GameState().GetGameState() != configs.GameStatePlaying {
		return
	}

	g.maskManager.UpdateCurrentMask(gameTime)
	currentMaskRowIndex := g.maskManager.currentMaskRowIndex
	topRow := RowCount - g.board.currentRowIndex - 1

	if currentMaskRowIndex >= topRow {
		mask := g.maskManager.currentMask
		g.board.rows[0].ApplyMask(*mask, mask.StartCol)
		g.maskManager.DestroyCurrentMask()
		nextMask := g.maskManager.GetNextMask()
		g.maskManager.SetCurrentMask(nextMask)
		fmt.Println("mask applied", nextMask.MaskType, nextMask.MaskShape, nextMask.StartCol, nextMask.Length)
	}

	row := g.board.rows[0]
	if slices.Contains(g.board.validRowStates, row.bits) {
		g.board.rows = g.board.rows[1 : g.board.currentRowIndex+1]
		g.board.currentRowIndex--
		if g.board.currentRowIndex < 0 {
			// TODO: Handle this case
			g.board.currentRowIndex = 0
		}
		configs.GameState().IncrementRowsCleared()
		fmt.Println("row removed", g.board.currentRowIndex, len(g.board.rows))
	}

	// Row interval: max(1.5, 6.0 - (elapsed/10)*0.1)
	g.board.rowSpawnInterval = float32(configs.RowInterval(configs.GameState().SurvivalTime))
	g.board.lastRowSpawnTime += configs.GameTime.Delta
	if g.board.lastRowSpawnTime >= g.board.rowSpawnInterval {
		g.board.currentRowIndex++
		if g.board.currentRowIndex >= RowCount-1 {
			configs.GameState().SetGameState(configs.GameStateGameOver)
		}
		g.board.lastRowSpawnTime = 0
		g.board.rows = append(g.board.rows, GenerateRow(configs.GameState().GetPhase()))
		fmt.Println("row spawned", g.board.currentRowIndex, len(g.board.rows))
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

	// draw current mask
	g.maskManager.RenderCurrentMask(g)

	g.End()
}

func (g *GridComponent) Reset() {
	board := &Board{
		rows:             make([]Row, 0),
		currentRowIndex:  InitialRowCount - 1,
		validRowStates:   make([]uint8, 0),
		rowSpawnInterval: float32(configs.InitialRowInterval),
	}
	for range InitialRowCount {
		board.rows = append(board.rows, GenerateRow(configs.GameState().GetPhase()))
	}
	board.validRowStates = append(board.validRowStates, 0b11111111)
	board.validRowStates = append(board.validRowStates, 0b0)

	g.board = board
	configs.GameState().ResetScore()

	maskManager := GetGlobalMaskManager()
	maskManager.SetCurrentMask(GenerateNewMask())
	maskManager.QueueMask(GenerateNewMask())
	g.maskManager = maskManager
}

func NewGridComponent() *GridComponent {
	cellSize := int32(28)
	cellPadding := int32(8)
	// cellColor := rl.Color{R: 255, G: 255, B: 255, A: 200} // white
	cellColor := rl.Color{R: 0, G: 255, B: 0, A: 200} // green

	gridBorderThickness := int32(2)
	// gridBorderColor := rl.Color{R: 255, G: 255, B: 255, A: 150} // white
	gridBorderColor := rl.Color{R: 0, G: 255, B: 0, A: 150} // green

	gridWidth := (ColumnCount * cellSize) + ((ColumnCount - 1) * cellPadding) + (2 * (gridBorderThickness + gridBorderThickness/2)) + (2 * OffsetX)
	gridHeight := (RowCount * cellSize) + ((RowCount - 1) * cellPadding) + (2 * (gridBorderThickness + gridBorderThickness/2)) + (2 * OffsetY)

	gridComponent := &GridComponent{
		BaseComponent:       components.NewBaseComponent("grid", gridWidth, gridHeight),
		cellSize:            cellSize,
		cellPadding:         cellPadding,
		cellColor:           cellColor,
		gridBorderThickness: gridBorderThickness,
		gridBorderColor:     gridBorderColor,
	}
	gridComponent.Reset()

	return gridComponent
}

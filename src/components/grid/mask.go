package grid

import (
	"math"

	"github.com/CCAtAlvis/xorfall/src/configs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MaskType int

const (
	MaskTypeOR MaskType = iota
	MaskTypeXOR
	MaskTypeNOT
	// MaskTypeAND
)

var MaskColorMap = map[MaskType]rl.Color{
	MaskTypeOR:  rl.Yellow,
	MaskTypeXOR: rl.Blue,
	MaskTypeNOT: rl.Red,
	// MaskTypeAND: rl.Yellow,
}

type Mask struct {
	MaskType   MaskType
	MaskShape  uint8
	Speed      int32   // in rows per second
	SpeedFloat float32 // in rows per second
	StartCol   int
	Length     int
}

func NewMask(maskType MaskType, maskShape uint8, speed int32, speedFloat float32, startCol int, length int) *Mask {
	mask := &Mask{
		MaskType:   maskType,
		MaskShape:  maskShape,
		Speed:      speed,
		SpeedFloat: speedFloat,
		StartCol:   startCol,
		Length:     length,
	}
	return mask
}

func GenerateNewMask() *Mask {
	// TODO: manage these probabilities better
	maskType := MaskType(rl.GetRandomValue(0, 2))
	length := rl.GetRandomValue(1, 4)
	maskShape := uint8(rl.GetRandomValue(0, int32(math.Pow(2, float64(length))-1)))
	// speed := int32(rl.GetRandomValue(1, 10))
	// speedFloat := float32(rl.GetRandomValue(1, 10))
	speed := int32(1)
	speedFloat := float32(1)
	startCol := rl.GetRandomValue(0, ColumnCount-length)

	mask := NewMask(maskType, maskShape, speed, speedFloat, int(startCol), int(length))
	return mask
}

type MaskManager struct {
	currentMask               *Mask
	currentMaskLastUpdateTime float32
	currentMaskRowIndex       int
	currentMaskKeepFalling    bool

	nextMasks []*Mask
}

func (m *MaskManager) QueueMask(mask *Mask) {
	m.nextMasks = append(m.nextMasks, mask)
}

func (m *MaskManager) PopMask() *Mask {
	if len(m.nextMasks) == 0 {
		return GenerateNewMask()
	}

	mask := m.nextMasks[0]
	m.nextMasks = m.nextMasks[1:]
	return mask
}

func (m *MaskManager) DestroyCurrentMask() {
	m.currentMask = nil
	m.currentMaskRowIndex = -1
	m.currentMaskLastUpdateTime = 0
	m.currentMaskKeepFalling = false
}

func (m *MaskManager) SetCurrentMask(mask *Mask) {
	m.currentMask = mask
}

func (m *MaskManager) UpdateCurrentMask(gameTime *configs.GameTimeManager) {
	if rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressedRepeat(rl.KeyRight) {
		newCol := m.currentMask.StartCol + 1
		if newCol+m.currentMask.Length <= ColumnCount {
			m.currentMask.StartCol = newCol
		}
	}

	if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressedRepeat(rl.KeyLeft) {
		newCol := m.currentMask.StartCol - 1
		if newCol >= 0 {
			m.currentMask.StartCol = newCol
		}
	}

	if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressedRepeat(rl.KeyDown) {
		newRow := m.currentMaskRowIndex + 1
		if newRow < RowCount {
			m.currentMaskRowIndex = newRow
			m.currentMaskLastUpdateTime = 0
		} else {
			// TODO: handle this..
			// the timer remaining for new row to add gets added to survival time
		}
	}

	if rl.IsKeyPressed(rl.KeySpace) {
		m.currentMaskKeepFalling = true
		m.currentMask.Speed = 20
	}

	m.currentMaskLastUpdateTime += configs.GameTime.Delta
	if m.currentMaskLastUpdateTime >= 1.0/float32(m.currentMask.Speed) {
		m.currentMaskLastUpdateTime = 0
		m.currentMaskRowIndex++
		if m.currentMaskRowIndex >= RowCount {
			// TODO: handle this..
			// the timer remaining for new row to add gets added to survival time
			// Trigger some event here

		}
	}
}

func (m *MaskManager) RenderCurrentMask(g *GridComponent) {
	// start position of current mask is one cell above the grid border

	mask := m.currentMask
	for i := range mask.Length {
		maskBit := mask.MaskShape&(1<<i) != 0

		maskStartX := OffsetX + int32(mask.StartCol+int(i))*(g.cellSize+g.cellPadding) + (g.gridBorderThickness + g.gridBorderThickness/2)
		maskStartY := OffsetY + int32(m.currentMaskRowIndex-1)*(g.cellSize+g.cellPadding) + (g.gridBorderThickness + g.gridBorderThickness/2)

		maskColor := MaskColorMap[mask.MaskType]
		if maskBit {
			rl.DrawRectangle(maskStartX, maskStartY, g.cellSize, g.cellSize, maskColor)
		} else {
			rl.DrawRectangleLinesEx(rl.Rectangle{X: float32(maskStartX), Y: float32(maskStartY), Width: float32(g.cellSize), Height: float32(g.cellSize)}, 4, maskColor)
		}
	}

}

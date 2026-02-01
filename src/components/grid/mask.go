package grid

import (
	"github.com/CCAtAlvis/xorfall/src/configs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MaskType int

const (
	MaskTypeOR MaskType = iota
	MaskTypeXOR
	MaskTypeNOT
	MaskTypeAND
	MaskTypeXNOR
)

// Color hierarchy: OR = best (gold), XOR/NOT = ok (blue/orange), AND = worst (muted), XNOR = red.
var MaskColorMap = map[MaskType]rl.Color{
	MaskTypeOR:   rl.Yellow,
	MaskTypeXOR:  rl.Blue,
	MaskTypeNOT:  rl.Orange,
	MaskTypeAND:  rl.Color{R: 130, G: 130, B: 140, A: 255},
	MaskTypeXNOR: rl.Red,
}

type Mask struct {
	MaskType  MaskType
	MaskShape uint8
	Speed     int32 // in rows per second
	StartCol  int
	Length    int
}

func NewMask(maskType MaskType, maskShape uint8, speed int32, startCol int, length int) *Mask {
	mask := &Mask{
		MaskType:  maskType,
		MaskShape: maskShape,
		Speed:     speed,
		StartCol:  startCol,
		Length:    length,
	}
	return mask
}

func GenerateNewMask() *Mask {
	phase := configs.GameState().GetPhase()
	maskType := rollOperator(phase)
	length := rollLength(phase)

	var maskShape uint8
	switch {
	case maskType == MaskTypeNOT:
		// NOT: all bits of mask should be set to 1.
		maskShape = uint8(1<<length) - 1
	case maskType == MaskTypeOR && (phase == configs.Phase1Learning || phase == configs.Phase2SkillBuilding):
		// OR with all 0s does nothing; in Phase 1 and 2, OR mask is all 1s only.
		maskShape = uint8(1<<length) - 1
	case maskType == MaskTypeAND:
		// AND with all 1s does nothing; random shape otherwise (0..all 1s).
		maskShape = uint8(rl.GetRandomValue(0, int32(1<<length)-1))
	default:
		maskShape = uint8(rl.GetRandomValue(0, int32(1<<length)-1))
	}

	speed := int32(3)
	startCol := rl.GetRandomValue(0, int32(ColumnCount-length))

	return NewMask(maskType, maskShape, speed, int(startCol), int(length))
}

func rollOperator(phase configs.Phase) MaskType {
	roll := rl.GetRandomValue(0, 99)
	switch phase {
	case configs.Phase1Learning:
		// OR 53%, XOR 30%, NOT 10%, AND 2%, XNOR 5%
		if roll < 53 {
			return MaskTypeOR
		}
		if roll < 83 {
			return MaskTypeXOR
		}
		if roll < 93 {
			return MaskTypeNOT
		}
		if roll < 95 {
			return MaskTypeAND
		}
		return MaskTypeXNOR
	case configs.Phase2SkillBuilding:
		// OR 32%, XOR 37%, NOT 20%, AND 5%, XNOR 6%
		if roll < 32 {
			return MaskTypeOR
		}
		if roll < 69 {
			return MaskTypeXOR
		}
		if roll < 89 {
			return MaskTypeNOT
		}
		if roll < 94 {
			return MaskTypeAND
		}
		return MaskTypeXNOR
	default: // Phase3Mastery
		// OR 12%, XOR 45%, NOT 25%, AND 8%, XNOR 10%
		if roll < 12 {
			return MaskTypeOR
		}
		if roll < 57 {
			return MaskTypeXOR
		}
		if roll < 82 {
			return MaskTypeNOT
		}
		if roll < 90 {
			return MaskTypeAND
		}
		return MaskTypeXNOR
	}
}

func rollLength(phase configs.Phase) int {
	roll := rl.GetRandomValue(0, 99)
	switch phase {
	case configs.Phase1Learning:
		if roll < 25 {
			return 1
		}
		if roll < 75 {
			return 2
		}
		if roll < 95 {
			return 3
		}
		return 4
	case configs.Phase2SkillBuilding:
		if roll < 15 {
			return 1
		}
		if roll < 50 {
			return 2
		}
		if roll < 90 {
			return 3
		}
		return 4
	default: // Phase3Mastery
		if roll < 5 {
			return 1
		}
		if roll < 20 {
			return 2
		}
		if roll < 70 {
			return 3
		}
		return 4
	}
}

type MaskManager struct {
	currentMask               *Mask
	currentMaskLastUpdateTime float32
	currentMaskRowIndex       int
	currentMaskKeepFalling    bool

	nextMasks []*Mask
}

var globalMaskManager = &MaskManager{
	nextMasks: make([]*Mask, 0),
}

func GetGlobalMaskManager() *MaskManager {
	return globalMaskManager
}

func (m *MaskManager) QueueMask(mask *Mask) {
	m.nextMasks = append(m.nextMasks, mask)
}

func (m *MaskManager) GetNextMask() *Mask {
	mask := m.nextMasks[0]
	m.nextMasks = m.nextMasks[1:]

	if len(m.nextMasks) == 0 {
		m.QueueMask(GenerateNewMask())
	}

	return mask
}

func (m *MaskManager) DestroyCurrentMask() {
	m.currentMask = nil
	m.currentMaskRowIndex = 0
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

	// Fall speed: min(2.5, 0.5 + (elapsed/30)*0.15); hard drop = 20 rows/sec
	effectiveSpeed := configs.FallSpeed(configs.GameState().SurvivalTime)
	if m.currentMaskKeepFalling {
		effectiveSpeed = 20
	}
	m.currentMaskLastUpdateTime += configs.GameTime.Delta
	if m.currentMaskLastUpdateTime >= 1.0/effectiveSpeed {
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

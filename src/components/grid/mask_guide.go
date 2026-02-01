package grid

import (
	"github.com/CCAtAlvis/xorfall/src/components"
	"github.com/CCAtAlvis/xorfall/src/configs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	guideExampleBits    = 4
	guideCellSize       = int32(24)
	guideCellPadding    = int32(4)
	guideLabelFontSize  = int32(32)
	guideTitleFontSize  = int32(32)
	guideSectionSpacing = int32(6)
	guideRowSpacing     = int32(4)
	guidePadding        = int32(8)
	guideLabelWidth     = int32(52)
)

var (
	guideMaskBits = [guideExampleBits]bool{true, true, false, false}
	guideRowBits  = [guideExampleBits]bool{true, false, true, false}
)

var (
	guideCellColor   = rl.Color{R: 0, G: 255, B: 0, A: 200}
	guideCellOutline = rl.Color{R: 0, G: 255, B: 0, A: 150}
	guideLabelColor  = rl.White
)

type MaskGuideComponent struct {
	components.BaseComponent
}

func NewMaskGuideComponent() *MaskGuideComponent {
	cellRowWidth := int32(guideExampleBits)*(guideCellSize+guideCellPadding) - guideCellPadding
	blockWidth := guideLabelWidth + guideCellPadding + cellRowWidth

	oneRowHeight := guideLabelFontSize + guideRowSpacing
	sectionHeight := guideTitleFontSize + guideSectionSpacing + 3*oneRowHeight
	totalHeight := int32(guidePadding*2 + 5*sectionHeight + 4*guideSectionSpacing)

	w := guidePadding*2 + blockWidth
	w = w * 2
	h := totalHeight + 3*oneRowHeight

	return &MaskGuideComponent{
		BaseComponent: components.NewBaseComponent("mask_guide", w, h),
	}
}

func (m *MaskGuideComponent) Update(gameTime *configs.GameTimeManager) {}

func applyExample(maskType MaskType, maskBit, rowBit bool) bool {
	switch maskType {
	case MaskTypeOR:
		return maskBit || rowBit
	case MaskTypeXOR:
		return maskBit != rowBit
	case MaskTypeNOT:
		return !maskBit
	case MaskTypeAND:
		return maskBit && rowBit
	case MaskTypeXNOR:
		return maskBit == rowBit
	default:
		return false
	}
}

func (m *MaskGuideComponent) Render() {
	m.Begin()
	rl.ClearBackground(rl.Blank)

	y := int32(guidePadding)
	maskTypes := []MaskType{MaskTypeOR, MaskTypeXOR, MaskTypeNOT, MaskTypeAND, MaskTypeXNOR}
	titles := []string{"gold (OR)", "blue (XOR)", "orange (NOT)", "gray (AND)", "red (XNOR)"}

	for si, maskType := range maskTypes {
		color := MaskColorMap[maskType]
		rl.DrawText(titles[si], guidePadding, y, guideTitleFontSize, color)
		y += guideTitleFontSize + guideSectionSpacing

		cellX := int32(guidePadding + guideLabelWidth + guideCellPadding)

		// mask row
		rl.DrawText("mask:", guidePadding, y+2, guideLabelFontSize, guideLabelColor)
		for i := range guideExampleBits {
			cx := cellX + int32(i)*(guideCellSize+guideCellPadding) + 70
			cy := y + 8
			if guideMaskBits[i] {
				rl.DrawRectangle(cx, cy, guideCellSize, guideCellSize, color)
			} else {
				rl.DrawRectangleLinesEx(rl.Rectangle{X: float32(cx), Y: float32(cy), Width: float32(guideCellSize), Height: float32(guideCellSize)}, 2, color)
			}
		}
		y += guideLabelFontSize + guideRowSpacing

		// row row
		rl.DrawText("row:", guidePadding, y+2, guideLabelFontSize, guideLabelColor)
		for i := range guideExampleBits {
			cx := cellX + int32(i)*(guideCellSize+guideCellPadding) + 70
			cy := y + 8
			if guideRowBits[i] {
				rl.DrawRectangle(cx, cy, guideCellSize, guideCellSize, guideCellColor)
			} else {
				rl.DrawRectangleLinesEx(rl.Rectangle{X: float32(cx), Y: float32(cy), Width: float32(guideCellSize), Height: float32(guideCellSize)}, 2, guideCellOutline)
			}
		}
		y += guideLabelFontSize + guideRowSpacing

		// result row
		rl.DrawText("result:", guidePadding, y+2, guideLabelFontSize, guideLabelColor)
		for i := range guideExampleBits {
			res := applyExample(maskType, guideMaskBits[i], guideRowBits[i])
			cx := cellX + int32(i)*(guideCellSize+guideCellPadding) + 70
			cy := y + 8
			if res {
				rl.DrawRectangle(cx, cy, guideCellSize, guideCellSize, guideCellColor)
			} else {
				rl.DrawRectangleLinesEx(rl.Rectangle{X: float32(cx), Y: float32(cy), Width: float32(guideCellSize), Height: float32(guideCellSize)}, 2, guideCellOutline)
			}
		}
		y += guideLabelFontSize + guideRowSpacing + guideCellSize + guideSectionSpacing
	}

	//debuf bounding box
	width, height := m.GetSize()
	rl.DrawRectangleLinesEx(rl.Rectangle{X: 0, Y: 0, Width: float32(width), Height: float32(height)}, 2, rl.Red)

	m.End()
}

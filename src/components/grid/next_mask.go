package grid

import (
	"github.com/CCAtAlvis/xorfall/src/components"
	"github.com/CCAtAlvis/xorfall/src/configs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	nextMaskCellSize      = int32(24)
	nextMaskCellPadding   = int32(4)
	nextMaskLabelFontSize = int32(32)
	nextMaskItemSpacing   = int32(8)
	nextMaskPadding       = int32(8)
	nextMaskMaxLength     = 8
)

type NextMaskComponent struct {
	components.BaseComponent
	maskManager *MaskManager
}

func NewNextMaskComponent() *NextMaskComponent {
	labelW := rl.MeasureText("next mask:", nextMaskLabelFontSize)
	cellRowW := int32(nextMaskMaxLength) * (nextMaskCellSize + nextMaskCellPadding)
	w := nextMaskPadding*2 + labelW + nextMaskCellPadding + cellRowW

	rowH := nextMaskLabelFontSize + nextMaskItemSpacing + nextMaskCellSize
	h := nextMaskPadding*2 + rowH - nextMaskItemSpacing
	return &NextMaskComponent{
		BaseComponent: components.NewBaseComponent("next_mask", w, h),
		maskManager:   GetGlobalMaskManager(),
	}
}

func (n *NextMaskComponent) Update(gameTime *configs.GameTimeManager) {}

func (n *NextMaskComponent) Render() {
	n.Begin()

	masks := n.maskManager.nextMasks
	mask := masks[0]

	rl.DrawText("next mask:", nextMaskPadding, nextMaskPadding, nextMaskLabelFontSize, rl.White)
	textWidth := rl.MeasureText("next mask:", nextMaskLabelFontSize)
	cellXBase := nextMaskPadding + textWidth + nextMaskItemSpacing + nextMaskCellPadding
	y := int32(12)

	color := MaskColorMap[mask.MaskType]
	for j := 0; j < mask.Length; j++ {
		maskBit := mask.MaskShape&(1<<j) != 0
		cx := nextMaskItemSpacing + cellXBase + int32(j)*(nextMaskCellSize+nextMaskCellPadding)
		if maskBit {
			rl.DrawRectangle(cx, y, nextMaskCellSize, nextMaskCellSize, color)
		} else {
			rl.DrawRectangleLinesEx(rl.Rectangle{X: float32(cx), Y: float32(y), Width: float32(nextMaskCellSize), Height: float32(nextMaskCellSize)}, 2, color)
		}
	}

	n.End()
}

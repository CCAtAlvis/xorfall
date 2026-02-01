package components

import (
	"github.com/CCAtAlvis/xorfall/src/configs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const IntroScreenCount = 4

type IntroComponent struct {
	BaseComponent
	screenIndex    int
	startRequested bool
}

func NewIntroComponent() *IntroComponent {
	return &IntroComponent{
		BaseComponent: NewBaseComponent("intro", configs.VirtualWidth, configs.VirtualHeight),
		screenIndex:   0,
	}
}

func (c *IntroComponent) Update(gameTime *configs.GameTimeManager) {
	if rl.IsKeyPressed(rl.KeySpace) {
		if c.screenIndex < IntroScreenCount-1 {
			c.screenIndex++
		} else {
			c.startRequested = true
		}
	}
}

func (c *IntroComponent) StartRequested() bool {
	r := c.startRequested
	c.startRequested = false
	return r
}

func (c *IntroComponent) Render() {
	c.Begin()
	rl.ClearBackground(rl.Black)
	c.drawScreen(c.screenIndex)
	c.End()
}

func (c *IntroComponent) drawScreen(idx int) {
	vw := configs.VirtualWidth
	vh := configs.VirtualHeight

	switch idx {
	case 0:
		// Screen 1: Title XORFALL (large, each letter different color) + setup + "Press SPACE to continue"
		title := "XORFALL"
		letterColors := []rl.Color{
			rl.Red, rl.Orange, rl.Yellow, rl.Green, rl.SkyBlue, rl.Blue, rl.Violet,
		}
		fontSize := int32(140)
		letterW := fontSize * 5 / 6
		totalTitleW := int32(len(title)) * letterW
		titleX := (vw - totalTitleW) / 2
		titleY := vh/2 - 180
		for i, r := range title {
			col := letterColors[i%len(letterColors)]
			rl.DrawText(string(r), titleX+int32(i)*letterW, titleY, fontSize, col)
		}
		setup1 := "The signal is unstable."
		setup2 := "Apply masks to control it's flow."
		fs := int32(36)
		w1 := rl.MeasureText(setup1, fs)
		w2 := rl.MeasureText(setup2, fs)
		rl.DrawText(setup1, (vw-w1)/2, titleY+160, fs, rl.White)
		rl.DrawText(setup2, (vw-w2)/2, titleY+160+50, fs, rl.White)
		inst := "Press SPACE to continue"
		wi := rl.MeasureText(inst, 28)
		rl.DrawText(inst, (vw-wi)/2, vh-80, 28, rl.LightGray)
	case 1:
		// Screen 2: Core idea + graphic (row of bits, mask above) + "Press SPACE to continue"
		line1 := "Masks fall from the top."
		line2 := "Apply them to rows to change bits."
		fs := int32(40)
		w1 := rl.MeasureText(line1, fs)
		w2 := rl.MeasureText(line2, fs)
		rl.DrawText(line1, (vw-w1)/2, 200, fs, rl.White)
		rl.DrawText(line2, (vw-w2)/2, 260, fs, rl.White)
		cellSize := int32(40)
		cellPad := int32(6)
		rowY := int32(450)
		startX := vw/2 - (8*(cellSize+cellPad)-cellPad)/2
		for i := 0; i < 8; i++ {
			x := startX + int32(i)*(cellSize+cellPad)
			rl.DrawRectangle(x, rowY, cellSize, cellSize, rl.DarkGray)
			rl.DrawRectangleLinesEx(rl.Rectangle{X: float32(x), Y: float32(rowY), Width: float32(cellSize), Height: float32(cellSize)}, 2, rl.White)
		}
		maskW := (cellSize + cellPad) * 3
		maskX := startX + (cellSize+cellPad)*2
		maskY := rowY - 50
		rl.DrawRectangle(maskX, maskY, maskW, 30, rl.SkyBlue)
		rl.DrawRectangleLinesEx(rl.Rectangle{X: float32(maskX), Y: float32(maskY), Width: float32(maskW), Height: 30}, 2, rl.White)
		inst := "Press SPACE to continue"
		wi := rl.MeasureText(inst, 28)
		rl.DrawText(inst, (vw-wi)/2, vh-80, 28, rl.LightGray)
	case 2:
		c.drawOperatorsScreen(vw, vh)
	case 3:
		// Screen 4: Controls + "Clear rows before they overflow." + "Press SPACE to start"
		fs := int32(36)
		ctrl1 := "Move mask [left]/[right]/[down]"
		ctrl2 := "[space] when you are confident"
		w1 := rl.MeasureText(ctrl1, fs)
		w2 := rl.MeasureText(ctrl2, fs)
		rl.DrawText(ctrl1, (vw-w1)/2, 280, fs, rl.White)
		rl.DrawText(ctrl2, (vw-w2)/2, 330, fs, rl.White)
		goal := "Clear rows before they overflow."
		wg := rl.MeasureText(goal, fs)
		rl.DrawText(goal, (vw-wg)/2, 420, fs, rl.LightGray)
		inst := "Press SPACE to start"
		wi := rl.MeasureText(inst, 28)
		rl.DrawText(inst, (vw-wi)/2, vh-80, 28, rl.LightGray)
	}
}

// intro operators screen layout (aligned with mask guide: same order and colors)
const (
	introExampleBits   = 4
	introCellSize      = int32(22)
	introCellPadding   = int32(4)
	introLabelFontSize = int32(24)
	introTitleFontSize = int32(26)
	introSectionGap    = int32(6)
	introRowGap        = int32(4)
	introLabelWidth    = int32(56)
)

var (
	introMaskBits = [introExampleBits]bool{true, true, false, false}
	introRowBits  = [introExampleBits]bool{true, false, true, false}
)

// introOpIndex matches grid.MaskType order: OR, XOR, NOT, AND, XNOR
const (
	introOpOR introOpIndex = iota
	introOpXOR
	introOpNOT
	introOpAND
	introOpXNOR
)

type introOpIndex int

func introApplyExample(op introOpIndex, maskBit, rowBit bool) bool {
	switch op {
	case introOpOR:
		return maskBit || rowBit
	case introOpXOR:
		return maskBit != rowBit
	case introOpNOT:
		return !maskBit
	case introOpAND:
		return maskBit && rowBit
	case introOpXNOR:
		return maskBit == rowBit
	default:
		return false
	}
}

func (c *IntroComponent) drawOperatorsScreen(vw, vh int32) {
	ops := []struct {
		title string
		col   rl.Color
	}{
		{"OR", rl.Yellow},
		{"XOR", rl.Blue},
		{"NOT", rl.Orange},
		{"AND", rl.Color{R: 130, G: 130, B: 140, A: 255}},
		{"XNOR", rl.Red},
	}

	cellRowWidth := int32(introExampleBits)*(introCellSize+introCellPadding) - introCellPadding
	blockWidth := introLabelWidth + introCellPadding + cellRowWidth
	leftX := (vw - blockWidth) / 2

	y := int32(60)
	cellColor := rl.Color{R: 0, G: 255, B: 0, A: 200}
	cellOutline := rl.Color{R: 0, G: 255, B: 0, A: 150}

	for si, op := range ops {
		rl.DrawText(op.title, leftX, y, introTitleFontSize, op.col)
		y += introTitleFontSize + introSectionGap

		cellX := leftX + introLabelWidth + introCellPadding + 40

		// mask row
		rl.DrawText("mask:", leftX, y+2, introLabelFontSize, rl.White)
		for i := range introExampleBits {
			cx := cellX + int32(i)*(introCellSize+introCellPadding)
			cy := y
			if introMaskBits[i] {
				rl.DrawRectangle(cx, cy, introCellSize, introCellSize, op.col)
			} else {
				rl.DrawRectangleLinesEx(rl.Rectangle{X: float32(cx), Y: float32(cy), Width: float32(introCellSize), Height: float32(introCellSize)}, 2, op.col)
			}
		}
		y += introLabelFontSize + introRowGap

		// row row
		rl.DrawText("row:", leftX, y+2, introLabelFontSize, rl.White)
		for i := range introExampleBits {
			cx := cellX + int32(i)*(introCellSize+introCellPadding)
			cy := y
			if introRowBits[i] {
				rl.DrawRectangle(cx, cy, introCellSize, introCellSize, cellColor)
			} else {
				rl.DrawRectangleLinesEx(rl.Rectangle{X: float32(cx), Y: float32(cy), Width: float32(introCellSize), Height: float32(introCellSize)}, 2, cellOutline)
			}
		}
		y += introLabelFontSize + introRowGap

		// result row
		rl.DrawText("result:", leftX, y+2, introLabelFontSize, rl.White)
		for i := range introExampleBits {
			res := introApplyExample(introOpIndex(si), introMaskBits[i], introRowBits[i])
			cx := cellX + int32(i)*(introCellSize+introCellPadding)
			cy := y
			if res {
				rl.DrawRectangle(cx, cy, introCellSize, introCellSize, cellColor)
			} else {
				rl.DrawRectangleLinesEx(rl.Rectangle{X: float32(cx), Y: float32(cy), Width: float32(introCellSize), Height: float32(introCellSize)}, 2, cellOutline)
			}
		}
		y += introLabelFontSize + introRowGap + introCellSize + introSectionGap
	}

	inst := "Press SPACE to continue"
	wi := rl.MeasureText(inst, 28)
	rl.DrawText(inst, (vw-wi)/2, vh-80, 28, rl.LightGray)
}

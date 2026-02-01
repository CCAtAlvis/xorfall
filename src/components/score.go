package components

import (
	"fmt"

	"github.com/CCAtAlvis/xorfall/src/configs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	scoreLabelFontSize = 50
	scoreValueFontSize = 64 // larger = "bold" emphasis (default font has no bold)
	scorePadding       = 8
	scoreLineSpacing   = 4
)

type ScoreComponent struct {
	BaseComponent
}

func NewScoreComponent() *ScoreComponent {
	// Width/height to fit "time:" / "ssss.mm" and "rows cleared:" / "<number>"
	// Approximate: "rows cleared:" is longest label; value font 24 for "9999" etc.

	textWidth := rl.MeasureText("rows cleared:", scoreLabelFontSize)

	w := int32(textWidth + 50)
	h := int32(scorePadding + scoreLabelFontSize + scoreLineSpacing + scoreValueFontSize +
		scorePadding + scoreLabelFontSize + scoreLineSpacing + scoreValueFontSize + scorePadding)
	return &ScoreComponent{
		BaseComponent: NewBaseComponent("score", w, h),
	}
}

func formatSurvivalTime(seconds float64) string {
	if seconds >= 1000 {
		return fmt.Sprintf("%d", int(seconds))
	}
	sec := int(seconds)
	ms := int((seconds - float64(sec)) * 100)
	return fmt.Sprintf("%d:%02d", sec, ms)
}

func (s *ScoreComponent) Update(gameTime *configs.GameTimeManager) {
	configs.GameState().AddSurvivalTime(float64(gameTime.Delta))
}

func (s *ScoreComponent) Render() {
	s.Begin()

	gs := configs.GameState()
	timeStr := formatSurvivalTime(gs.SurvivalTime)
	rowsStr := fmt.Sprintf("%d", gs.RowsCleared)

	x := int32(scorePadding)
	y := int32(scorePadding)

	// time:
	rl.DrawText("time:", x, y, scoreLabelFontSize, rl.White)
	y += scoreLabelFontSize + scoreLineSpacing
	rl.DrawText(timeStr, x, y, scoreValueFontSize, rl.White)
	y += scoreValueFontSize + scorePadding

	// rows cleared:
	rl.DrawText("rows cleared:", x, y, scoreLabelFontSize, rl.White)
	y += scoreLabelFontSize + scoreLineSpacing
	rl.DrawText(rowsStr, x, y, scoreValueFontSize, rl.White)

	s.End()
}

package components

import (
	"strconv"

	"github.com/CCAtAlvis/xorfall/src/configs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameOverComponent struct {
	BaseComponent
	restartRequested bool
}

func NewGameOverComponent() *GameOverComponent {
	return &GameOverComponent{
		BaseComponent: NewBaseComponent("gameover", configs.VirtualWidth, configs.VirtualHeight),
	}
}

func (c *GameOverComponent) Update(gameTime *configs.GameTimeManager) {
	if rl.IsKeyPressed(rl.KeySpace) {
		c.restartRequested = true
	}
}

func (c *GameOverComponent) RestartRequested() bool {
	r := c.restartRequested
	c.restartRequested = false
	return r
}

func (c *GameOverComponent) Render() {
	c.Begin()
	rl.ClearBackground(rl.Black)
	c.drawContent()
	c.End()
}

func (c *GameOverComponent) drawContent() {
	vw := configs.VirtualWidth
	vh := configs.VirtualHeight
	gs := configs.GameState()

	title := "GAME OVER"
	fsTitle := int32(72)
	wTitle := rl.MeasureText(title, fsTitle)
	rl.DrawText(title, (vw-wTitle)/2, vh/2-120, fsTitle, rl.Red)

	timeSurvived := int(gs.SurvivalTime)
	rowsCleared := gs.RowsCleared
	fs := int32(36)
	line1 := "Time survived: " + strconv.Itoa(timeSurvived) + " seconds"
	line2 := "Rows cleared: " + strconv.Itoa(rowsCleared)
	w1 := rl.MeasureText(line1, fs)
	w2 := rl.MeasureText(line2, fs)
	rl.DrawText(line1, (vw-w1)/2, vh/2-20, fs, rl.White)
	rl.DrawText(line2, (vw-w2)/2, vh/2+30, fs, rl.White)

	inst := "Press SPACE to restart"
	wi := rl.MeasureText(inst, 28)
	rl.DrawText(inst, (vw-wi)/2, vh-80, 28, rl.LightGray)
}

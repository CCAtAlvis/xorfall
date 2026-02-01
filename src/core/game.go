package core

import (
	"math"

	"github.com/CCAtAlvis/xorfall/src/components"
	"github.com/CCAtAlvis/xorfall/src/components/grid"
	"github.com/CCAtAlvis/xorfall/src/configs"
	"github.com/CCAtAlvis/xorfall/src/render"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	containers []*render.Container

	scale   float32
	offsetX float32
	offsetY float32

	isFullscreen      bool
	introComponent    *components.IntroComponent
	gameOverComponent *components.GameOverComponent
	gridComponent     *grid.GridComponent
}

func NewGame() *Game {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(1280, 720, "My Raylib Game")
	isFullscreen := false
	if !rl.IsWindowFullscreen() {
		isFullscreen = true
		rl.ToggleFullscreen()
	}
	rl.SetTargetFPS(configs.TargetFPS)
	configs.Init()

	// debugComponent := components.NewDebugComponent()
	// debugWidth, debugHeight := debugComponent.GetSize()

	gridComponent := grid.NewGridComponent()
	gridWidth, gridHeight := gridComponent.GetSize()
	gridOffsetX := int32(float32(configs.VirtualWidth)*0.5 - float32(gridWidth)*0.5)
	gridOffsetY := int32(float32(configs.VirtualHeight)*0.5 - float32(gridHeight)*0.5)

	scoreComponent := components.NewScoreComponent()
	scoreWidth, _ := scoreComponent.GetSize()
	scoreOffsetX := configs.VirtualWidth - scoreWidth - 10
	scoreOffsetY := gridOffsetY

	// Left column: FPS, then mask guide, then next mask
	maskGuideComponent := grid.NewMaskGuideComponent()
	_, maskGuideHeight := maskGuideComponent.GetSize()

	nextMaskComponent := grid.NewNextMaskComponent()
	offsetX := int32(10)
	maskGuideY := int32(20)
	nextMaskY := maskGuideY + maskGuideHeight + 50

	introComponent := components.NewIntroComponent()
	gameOverComponent := components.NewGameOverComponent()

	containers := []*render.Container{
		// {
		// 	Component: debugComponent,
		// 	Tint:      rl.Color{R: 255, G: 255, B: 255, A: 200},
		// 	OffsetX:   int32(float32(configs.VirtualWidth)*0.5 - float32(debugWidth)*0.5),
		// 	OffsetY:   int32(float32(configs.VirtualHeight)*0.5 - float32(debugHeight)*0.5),
		// 	Visible:   true,
		// },
		&render.Container{
			Component: introComponent,
			Tint:      rl.White,
			OffsetX:   0,
			OffsetY:   0,
			Enabled:   true,
			Visible:   true,
		},
		&render.Container{
			Component: gameOverComponent,
			Tint:      rl.White,
			OffsetX:   0,
			OffsetY:   0,
			Enabled:   true,
			Visible:   true,
		},
		&render.Container{
			Component: gridComponent,
			Tint:      rl.White,
			OffsetX:   gridOffsetX,
			OffsetY:   gridOffsetY,
			Enabled:   true,
			Visible:   true,
		},
		&render.Container{
			Component: maskGuideComponent,
			Tint:      rl.White,
			OffsetX:   offsetX,
			OffsetY:   maskGuideY,
			Enabled:   true,
			Visible:   true,
		},
		&render.Container{
			Component: nextMaskComponent,
			Tint:      rl.White,
			OffsetX:   offsetX,
			OffsetY:   nextMaskY,
			Enabled:   true,
			Visible:   true,
		},
		&render.Container{
			Component: scoreComponent,
			Tint:      rl.White,
			OffsetX:   scoreOffsetX,
			OffsetY:   scoreOffsetY,
			Enabled:   true,
			Visible:   true,
		},
		// {
		// 	Component: components.NewFPSComponent(),
		// 	Tint:      rl.White,
		// 	OffsetX:   offsetX,
		// 	OffsetY:   offsetY,
		// },
	}

	game := &Game{
		containers:        containers,
		isFullscreen:      isFullscreen,
		introComponent:    introComponent,
		gameOverComponent: gameOverComponent,
		gridComponent:     gridComponent,
	}

	game.recalculateScale()

	return game
}

func (g *Game) Update() {
	if rl.IsWindowResized() {
		g.recalculateScale()
	}

	gs := configs.GameState()
	state := gs.GetGameState()

	// Set Enabled and Visible per container from state. Enabled => update and draw; Visible => draw only.
	for _, c := range g.containers {
		switch c.Component.(type) {
		case *components.IntroComponent:
			c.Enabled = state == configs.GameStateIntroScreen
			c.Visible = state == configs.GameStateIntroScreen
		case *components.GameOverComponent:
			c.Enabled = state == configs.GameStateGameOver
			c.Visible = state == configs.GameStateGameOver
		default:
			c.Enabled = state == configs.GameStatePlaying
			c.Visible = state == configs.GameStatePlaying
		}
	}

	if state == configs.GameStatePlaying {
		configs.GameTime.Update()
	}

	for _, c := range g.containers {
		if !c.Enabled {
			continue
		}
		c.Component.Update(configs.GameTime)
		c.Render()
	}

	if state == configs.GameStateIntroScreen && g.introComponent.StartRequested() {
		gs.SetGameState(configs.GameStatePlaying)
	}
	if state == configs.GameStateGameOver && g.gameOverComponent.RestartRequested() {
		g.gridComponent.Reset()
		gs.SetGameState(configs.GameStatePlaying)
	}
}

func (g *Game) Render() {
	rl.BeginTextureMode(configs.VirtualScreen)
	rl.ClearBackground(rl.Black)

	for _, c := range g.containers {
		if !c.Enabled || !c.Visible {
			continue
		}
		c.Draw()
	}

	rl.EndTextureMode()
}

func (g *Game) Draw() {
	rl.BeginDrawing()

	source := rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(configs.VirtualWidth),
		Height: -float32(configs.VirtualHeight), // Flip Y (raylib quirk)
	}

	dest := rl.Rectangle{
		X:      g.offsetX,
		Y:      g.offsetY,
		Width:  float32(configs.VirtualWidth) * g.scale,
		Height: float32(configs.VirtualHeight) * g.scale,
	}

	rl.DrawTexturePro(
		configs.VirtualScreen.Texture,
		source,
		dest,
		rl.Vector2{},
		0,
		rl.White,
	)

	// debug with a white rectangle border for virtual screen
	// rl.DrawRectangleLinesEx(rl.Rectangle{X: g.offsetX, Y: g.offsetY, Width: float32(configs.VirtualWidth) * g.scale, Height: float32(configs.VirtualHeight) * g.scale}, 10, rl.White)

	rl.EndDrawing()
}

func (g *Game) Close() {
	rl.CloseWindow()
}

func (g *Game) Start() {
	var update = func() {
		g.Update()
		g.Render()
		g.Draw()
	}

	// rl.SetMainLoop(update)
	for !rl.WindowShouldClose() {
		update()
	}
}

func (g *Game) recalculateScale() {
	var screenW, screenH float32
	if g.isFullscreen {
		screenW = float32(rl.GetRenderWidth())
		screenH = float32(rl.GetRenderHeight())
	} else {
		screenW = float32(rl.GetScreenWidth())
		screenH = float32(rl.GetScreenHeight())
	}

	scaleX := screenW / float32(configs.VirtualWidth)
	scaleY := screenH / float32(configs.VirtualHeight)

	// Keep aspect ratio
	g.scale = float32(math.Min(float64(scaleX), float64(scaleY)))

	// Center the game
	g.offsetX = (screenW - float32(configs.VirtualWidth)*g.scale) * 0.5
	g.offsetY = (screenH - float32(configs.VirtualHeight)*g.scale) * 0.5
}

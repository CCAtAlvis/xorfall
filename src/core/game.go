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
	containers []render.Container

	scale   float32
	offsetX float32
	offsetY float32

	isFullscreen bool
}

func NewGame() *Game {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(1280, 720, "My Raylib Game")
	isFullscreen := false
	if !rl.IsWindowFullscreen() {
		// isFullscreen = true
		// rl.ToggleFullscreen()
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
	_, mgH := maskGuideComponent.GetSize()

	nextMaskComponent := grid.NewNextMaskComponent()
	leftX := int32(10)
	fpsY := int32(10)
	maskGuideY := fpsY + 120 + 10
	nextMaskY := maskGuideY + mgH + 50

	containers := []render.Container{
		// {
		// 	Component: debugComponent,
		// 	Tint:      rl.Color{R: 255, G: 255, B: 255, A: 200},
		// 	OffsetX:   int32(float32(configs.VirtualWidth)*0.5 - float32(debugWidth)*0.5),
		// 	OffsetY:   int32(float32(configs.VirtualHeight)*0.5 - float32(debugHeight)*0.5),
		// 	Visible:   true,
		// },
		{
			Component: gridComponent,
			Tint:      rl.White,
			OffsetX:   gridOffsetX,
			OffsetY:   gridOffsetY,
		},
		{
			Component: maskGuideComponent,
			Tint:      rl.White,
			OffsetX:   leftX,
			OffsetY:   maskGuideY,
		},
		{
			Component: nextMaskComponent,
			Tint:      rl.White,
			OffsetX:   leftX,
			OffsetY:   nextMaskY,
		},
		// Score: right side, top-aligned (mirrors FPS on left). Alternatives: center vertically with OffsetY = (VirtualHeight - scoreHeight)/2, or bottom with OffsetY = VirtualHeight - scoreHeight - 10.
		{
			Component: components.NewScoreComponent(),
			Tint:      rl.White,
			OffsetX:   scoreOffsetX,
			OffsetY:   scoreOffsetY,
		},
		{
			Component: components.NewFPSComponent(),
			Tint:      rl.White,
			OffsetX:   leftX,
			OffsetY:   fpsY,
		},
	}

	game := &Game{
		containers:   containers,
		isFullscreen: isFullscreen,
	}

	game.recalculateScale()

	return game
}

func (g *Game) Update() {
	if rl.IsWindowResized() {
		g.recalculateScale()
	}

	configs.GameTime.Update()

	for _, container := range g.containers {
		component := container.Component
		component.Update(configs.GameTime)
		container.Render()
	}
}

func (g *Game) Render() {
	rl.BeginTextureMode(configs.VirtualScreen)
	rl.ClearBackground(rl.Blank)
	for _, container := range g.containers {
		container.Draw()
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
	rl.DrawRectangleLinesEx(rl.Rectangle{X: g.offsetX, Y: g.offsetY, Width: float32(configs.VirtualWidth) * g.scale, Height: float32(configs.VirtualHeight) * g.scale}, 10, rl.White)

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

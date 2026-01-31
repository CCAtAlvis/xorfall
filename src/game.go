package src

import (
	"math"

	"github.com/CCAtAlvis/xorfall/src/renderer"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	renderTarget rl.RenderTexture2D

	virtualWidth  int32
	virtualHeight int32

	scale   float32
	offsetX float32
	offsetY float32

	grid *renderer.Grid
}

func (g *Game) Update() {
	if rl.IsWindowResized() {
		g.recalculateScale()
	}
}

// func (g *Game) Draw() {
// 	rl.BeginDrawing()
// 	rl.ClearBackground(rl.Black)
// 	rl.DrawText("Congrats! You created your first window!", 190, 1180, 20, rl.White)
// 	rl.DrawText("Congrats! You created your first window!", 190, 1112, 20, rl.White)
// 	rl.DrawText(fmt.Sprintf("Screen Width: %d", rl.GetScreenWidth()), 190, 240, 20, rl.Red)
// 	rl.DrawText(fmt.Sprintf("Screen Height: %d", rl.GetScreenHeight()), 190, 280, 20, rl.Red)
// 	rl.DrawText(fmt.Sprintf("Render Width: %d", rl.GetRenderWidth()), 190, 320, 20, rl.Red)
// 	rl.DrawText(fmt.Sprintf("Render Height: %d", rl.GetRenderHeight()), 190, 360, 20, rl.Red)
// 	rl.EndDrawing()
// }

func (g *Game) Draw() {
	// 1️⃣ Draw everything to virtual screen
	rl.BeginTextureMode(g.renderTarget)
	rl.ClearBackground(rl.Black)

	rl.DrawText("Hello Virtual World!", 100, 100, 40, rl.White)
	rl.DrawRectangle(300, 200, 200, 100, rl.Red)
	rl.DrawRectangle(700, 200, 100, 100, rl.Red)

	g.grid.Draw()

	rl.EndTextureMode()

	// 2️⃣ Draw scaled result to real screen
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	source := rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(g.virtualWidth),
		Height: -float32(g.virtualHeight), // Flip Y (raylib quirk)
	}

	dest := rl.Rectangle{
		X:      g.offsetX,
		Y:      g.offsetY,
		Width:  float32(g.virtualWidth) * g.scale,
		Height: float32(g.virtualHeight) * g.scale,
	}

	rl.DrawTexturePro(
		g.renderTarget.Texture,
		source,
		dest,
		rl.Vector2{},
		0,
		rl.White,
	)

	rl.EndDrawing()
}

func NewGame() *Game {
	// screenWidth := rl.GetScreenWidth()
	// screenHeight := rl.GetScreenHeight()
	// rl.InitWindow(int32(screenWidth), int32(screenHeight), "raylib [core] example - basic window")
	// // rl.InitWindow(800, 450, "raylib [core] example - basic window")
	// rl.ToggleFullscreen()
	// rl.SetTargetFPS(60)

	// return &Game{
	// 	renderTarget:  rl.LoadRenderTexture(int32(screenWidth), int32(screenHeight)),
	// 	virtualWidth:  int32(screenWidth),
	// 	virtualHeight: int32(screenHeight),
	// 	scale:         1.0,
	// 	offsetX:       0.0,
	// 	offsetY:       0.0,
	// }

	const virtualWidth = 1280
	const virtualHeight = 720

	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(1280, 720, "My Raylib Game")
	rl.SetTargetFPS(60)

	renderTarget := rl.LoadRenderTexture(virtualWidth, virtualHeight)

	grid := renderer.NewGrid(20, 8, 32, 4, rl.Red)
	game := &Game{
		renderTarget:  renderTarget,
		virtualWidth:  virtualWidth,
		virtualHeight: virtualHeight,
		grid:          grid,
	}

	game.recalculateScale()

	return game
}

func (g *Game) Close() {
	rl.CloseWindow()
}

func (g *Game) Start() {
	var update = func() {
		g.Update()
		g.Draw()
	}

	// rl.SetMainLoop(update)
	for !rl.WindowShouldClose() {
		update()
	}
}

func (g *Game) recalculateScale() {
	screenW := float32(rl.GetScreenWidth())
	screenH := float32(rl.GetScreenHeight())

	scaleX := screenW / float32(g.virtualWidth)
	scaleY := screenH / float32(g.virtualHeight)

	// Keep aspect ratio
	g.scale = float32(math.Min(float64(scaleX), float64(scaleY)))

	// Center the game
	g.offsetX = (screenW - float32(g.virtualWidth)*g.scale) * 0.5
	g.offsetY = (screenH - float32(g.virtualHeight)*g.scale) * 0.5
}

package configs

import rl "github.com/gen2brain/raylib-go/raylib"

var TargetFPS int32 = 60
var VirtualWidth int32 = 1920
var VirtualHeight int32 = 1080
var VirtualScreen rl.RenderTexture2D
var EnableGhostPreview = false

func Init() {
	VirtualScreen = rl.LoadRenderTexture(VirtualWidth, VirtualHeight)
}

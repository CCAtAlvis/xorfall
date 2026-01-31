package configs

import rl "github.com/gen2brain/raylib-go/raylib"

var TargetFPS int32 = 60
var VirtualWidth int32 = 1280
var VirtualHeight int32 = 720
var VirtualScreen rl.RenderTexture2D

func Init() {
	VirtualScreen = rl.LoadRenderTexture(VirtualWidth, VirtualHeight)
}

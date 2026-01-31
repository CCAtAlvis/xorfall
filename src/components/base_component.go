package components

import rl "github.com/gen2brain/raylib-go/raylib"

type BaseComponent struct {
	rt     rl.RenderTexture2D
	width  int32
	height int32
}

func NewBaseComponent(w, h int32) BaseComponent {
	return BaseComponent{
		rt:     rl.LoadRenderTexture(w, h),
		width:  w,
		height: h,
	}
}

func (b *BaseComponent) GetTexture() rl.Texture2D {
	return b.rt.Texture
}

func (b *BaseComponent) GetSize() (int32, int32) {
	return b.width, b.height
}

func (b *BaseComponent) Begin() {
	rl.BeginTextureMode(b.rt)
	rl.ClearBackground(rl.Blank)
}

func (b *BaseComponent) End() {
	rl.EndTextureMode()
}

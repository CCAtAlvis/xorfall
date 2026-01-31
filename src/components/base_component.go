package components

import (
	"github.com/CCAtAlvis/xorfall/src/configs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Component interface {
	GetId() string
	GetTexture() rl.Texture2D
	GetSize() (int32, int32)
	Update(gameTime *configs.GameTimeManager)
	Render()
}

type BaseComponent struct {
	id     string
	rt     rl.RenderTexture2D
	width  int32
	height int32
}

func NewBaseComponent(id string, w, h int32) BaseComponent {
	return BaseComponent{
		rt:     rl.LoadRenderTexture(w, h),
		width:  w,
		height: h,
	}
}

func (b *BaseComponent) GetId() string {
	return b.id
}

func (b *BaseComponent) GetTexture() rl.Texture2D {
	return b.rt.Texture
}

func (b *BaseComponent) GetSize() (int32, int32) {
	return b.width, b.height
}

func (b *BaseComponent) Begin() {
	rl.BeginTextureMode(b.rt)
	rl.ClearBackground(rl.Black)
}

func (b *BaseComponent) End() {
	rl.EndTextureMode()
}

package render

import (
	"github.com/CCAtAlvis/xorfall/src/components"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Container struct {
	Component components.Component
	OffsetX   int32
	OffsetY   int32
	Width     int32
	Height    int32
	Tint      rl.Color
	Visible   bool
	Layer     int32
}

func (c *Container) Render() {
	c.Component.Render()
}

func (c *Container) Draw() {
	tex := c.Component.GetTexture()
	w, h := c.Component.GetSize()

	src := rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(w),
		Height: -float32(h), // flip Y
	}

	dst := rl.Rectangle{
		X:      float32(c.OffsetX),
		Y:      float32(c.OffsetY),
		Width:  float32(w),
		Height: float32(h),
	}

	rl.DrawTexturePro(tex, src, dst, rl.Vector2{}, 0, c.Tint)
}

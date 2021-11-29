package missile

import (
    "github.com/veandco/go-sdl2/sdl"
    . "TowerDefense/common"
    "math"
)

type IMissile interface {
    GetXPos() int32
    GetYPos() int32
    GetTexRotationDeg() float64

    Update()

    Render(renderer *sdl.Renderer)
}

type CannonBall struct {
    XPos int32
    YPos int32
    RotationRad float64
    Speed int32
}
var _ IMissile = (*CannonBall)(nil)

func (c *CannonBall) GetXPos() int32 { return c.XPos }
func (c *CannonBall) GetYPos() int32 { return c.YPos }
func (c *CannonBall) GetTexRotationDeg() float64 { return 0 }

func (c *CannonBall) Update() {
    ASSERT_TRUE(c.Speed > 0)
    c.XPos += int32(math.Cos(c.RotationRad) * float64(c.Speed))
    c.YPos += int32(math.Sin(c.RotationRad) * float64(c.Speed))
}

func (c *CannonBall) Render(renderer *sdl.Renderer) {
    tex := TEXTURES[TEXTURE_FILENAME_CANNONBALL]
    rect := sdl.Rect{
        X: c.GetXPos(), Y: c.GetYPos(),
        W: int32(FIELD_SIZE_PX)/2, H: int32(FIELD_SIZE_PX)/2}
    renderer.Copy(tex.Texture, nil, &rect)
}

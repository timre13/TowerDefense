package missile

import (
    "github.com/veandco/go-sdl2/sdl"
    . "TowerDefense/common"
    "TowerDefense/enemy"
    "math"
)

//------------------------------------------------------------------------------

func missileDoUpdate(m IMissile, enemies []enemy.IEnemy) {
    ASSERT_TRUE(m.getSpeed() > 0)
    m.setCol(m.getCol() + math.Cos(m.getRotRad()) * float64(m.getSpeed()))
    m.setRow(m.getRow() + math.Sin(m.getRotRad()) * float64(m.getSpeed()))

    for _, e := range enemies {
        // If hit an enemy
        if CalcDistance(Vec2DF{X: float64(m.GetXPos()), Y: float64(m.GetYPos())},
                        Vec2DF{X: float64(e.GetXPos()), Y: float64(e.GetYPos())}) < FIELD_SIZE_PX {
               e.Damage(1)
               m.setHit(true)
        }
    }
}

func missileRender(renderer *sdl.Renderer, m IMissile, texName string) {
    tex := TEXTURES[texName]
    rect := sdl.Rect{
        X: m.GetXPos(), Y: m.GetYPos(),
        W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
    renderer.CopyEx(tex.Texture, nil, &rect, RadToDeg(m.getRotRad()), nil, sdl.FLIP_NONE)
}

//------------------------------------------------------------------------------

type IMissile interface {
    GetXPos() int32
    GetYPos() int32
    getCol() float64
    getRow() float64
    getSpeed() float64
    getRotRad() float64
    HasHit() bool

    setCol(col float64)
    setRow(row float64)
    setHit(val bool)

    Update(enemies []enemy.IEnemy)

    Render(renderer *sdl.Renderer)
}

type CannonBall struct {
    Col float64
    Row float64
    RotationRad float64
    Speed float64
    hasHit bool
}
var _ IMissile = (*CannonBall)(nil)

func (m *CannonBall) GetXPos() int32 { return int32(m.Col*FIELD_SIZE_PX) }
func (m *CannonBall) GetYPos() int32 { return int32(m.Row*FIELD_SIZE_PX) }
func (m *CannonBall) getCol() float64 { return m.Col }
func (m *CannonBall) getRow() float64 { return m.Row }
func (m *CannonBall) getSpeed() float64 { return m.Speed }
func (m *CannonBall) getRotRad() float64 { return m.RotationRad }
func (m *CannonBall) HasHit() bool { return m.hasHit }

func (m *CannonBall) setCol(col float64) { m.Col = col }
func (m *CannonBall) setRow(row float64) { m.Row = row }
func (m *CannonBall) setHit(val bool) { m.hasHit = val }

func (m *CannonBall) Update(enemies []enemy.IEnemy) {
    missileDoUpdate(m, enemies)
}

func (m *CannonBall) Render(renderer *sdl.Renderer) {
    missileRender(renderer, m, TEXTURE_FILENAME_CANNONBALL)
}


type Rocket struct {
    Col float64
    Row float64
    RotationRad float64
    Speed float64
    hasHit bool
}
var _ IMissile = (*Rocket)(nil)

func (m *Rocket) GetXPos() int32 { return int32(m.Col*FIELD_SIZE_PX) }
func (m *Rocket) GetYPos() int32 { return int32(m.Row*FIELD_SIZE_PX) }
func (m *Rocket) getCol() float64 { return m.Col }
func (m *Rocket) getRow() float64 { return m.Row }
func (m *Rocket) getSpeed() float64 { return m.Speed }
func (m *Rocket) getRotRad() float64 { return m.RotationRad }
func (m *Rocket) GetTexRotationDeg() float64 { return 0 }
func (m *Rocket) HasHit() bool { return m.hasHit }

func (m *Rocket) setCol(col float64) { m.Col = col }
func (m *Rocket) setRow(row float64) { m.Row = row }
func (m *Rocket) setHit(val bool) { m.hasHit = val }

func (m *Rocket) Update(enemies []enemy.IEnemy) {
    missileDoUpdate(m, enemies)
}

func (m *Rocket) Render(renderer *sdl.Renderer) {
    missileRender(renderer, m, TEXTURE_FILENAME_ROCKET)
}

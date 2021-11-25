package tower

import (
    "github.com/veandco/go-sdl2/sdl"
    "math"
    . "TowerDefense/common"
    "TowerDefense/enemy"
)

//------------------------------------------------------------------------------

type TowerType int
const (
    TOWER_TYPE_NONE TowerType = iota
    TOWER_TYPE_CANNON
    TOWER_TYPE_ROCKETTOWER
)

func (t TowerType) GetPrice() int {
    switch (t) {
    default: panic(t)
    case TOWER_TYPE_NONE:               panic(t)
    case TOWER_TYPE_CANNON:             return 10
    case TOWER_TYPE_ROCKETTOWER:        return 50
    }
}

func (t TowerType) GetInitialHP() int {
    switch (t) {
    default: panic(t)
    case TOWER_TYPE_NONE:               panic(t)
    case TOWER_TYPE_CANNON:             return 100
    case TOWER_TYPE_ROCKETTOWER:        return 300
    }
}

//------------------------------------------------------------------------------

type ITower interface {
    GetFieldCol() int32
    GetFieldRow() int32
    GetHP() int
    /*
     * If the tower is phisically on the map, this is true.
     * If the tower is used on the bottom bar or otherwise an indicator, this is false.
    */
    IsReal() bool
    GetRotationDeg() float64

    SetReal(val bool)
    SetRotationDeg(val float64)

    Update(enemies []enemy.IEnemy)

    Render(renderer *sdl.Renderer)
}

//------------------------------------------------------------------------------

type Cannon struct {
    FieldCol int32;
    FieldRow int32;

    IsReal_ bool;

    RotationDeg float64

    Hp int
}
var _ ITower = (*Cannon)(nil)

func (c *Cannon) GetFieldCol() int32 { return c.FieldCol }
func (c *Cannon) GetFieldRow() int32 { return c.FieldRow }
func (c *Cannon) GetHP() int { return c.Hp; }
func (c *Cannon) IsReal() bool { return c.IsReal_ }
func (c *Cannon) GetRotationDeg() float64 { return c.RotationDeg }

func (c *Cannon) SetReal(val bool) { c.IsReal_ = val }
func (c *Cannon) SetRotationDeg(val float64) { c.RotationDeg = val }

func (c *Cannon) Update(enemies []enemy.IEnemy) {
    cloCol, cloRow := enemy.GetClosestEnemyPos(enemies, float64(c.FieldCol), float64(c.FieldRow))
    rotRad := math.Atan2(float64(cloRow - c.FieldRow), float64(cloCol - c.FieldCol))
    rotDeg := rotRad * (180/math.Pi) + 90
    c.RotationDeg += (rotDeg - c.RotationDeg) / 10
}

func (c *Cannon) Render(renderer *sdl.Renderer) {
    // Render body
    tex := TEXTURES[TEXTURE_FILENAME_CANNON_BASE]
    rect := sdl.Rect{
        X: int32(float64(c.GetFieldCol())*FIELD_SIZE_PX), Y: int32(float64(c.GetFieldRow())*FIELD_SIZE_PX),
        W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
    renderer.Copy(tex.Texture, nil, &rect)

    // Render head
    tex = TEXTURES[TEXTURE_FILENAME_CANNON_HEAD]
    rect = sdl.Rect{
        X: int32(float64(c.GetFieldCol())*FIELD_SIZE_PX), Y: int32(float64(c.GetFieldRow())*FIELD_SIZE_PX),
        W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
    renderer.CopyEx(tex.Texture, nil, &rect, c.GetRotationDeg(), nil, 0)
}

//------------------------------------------------------------------------------

type RocketTower struct {
    FieldCol int32;
    FieldRow int32;

    IsReal_ bool;

    RotationDeg float64

    Hp int
}
var _ ITower = (*RocketTower)(nil)

func (c *RocketTower) GetFieldCol() int32 { return c.FieldCol }
func (c *RocketTower) GetFieldRow() int32 { return c.FieldRow }
func (c *RocketTower) GetHP() int { return c.Hp; }
func (c *RocketTower) IsReal() bool { return c.IsReal_ }
func (c *RocketTower) GetRotationDeg() float64 { return c.RotationDeg }

func (c *RocketTower) SetReal(val bool) { c.IsReal_ = val }
func (c *RocketTower) SetRotationDeg(val float64) { c.RotationDeg = val }

func (c *RocketTower) Update(enemies []enemy.IEnemy) {
    cloCol, cloRow := enemy.GetClosestEnemyPos(enemies, float64(c.FieldCol), float64(c.FieldRow))
    rotRad := math.Atan2(float64(cloRow - c.FieldRow), float64(cloCol - c.FieldCol))
    rotDeg := rotRad * (180/math.Pi) + 90
    c.RotationDeg += (rotDeg - c.RotationDeg) / 10
}

func (c *RocketTower) Render(renderer *sdl.Renderer) {
    // Render body
    tex := TEXTURES[TEXTURE_FILENAME_ROCKETTOWER_BASE]
    rect := sdl.Rect{
        X: int32(float64(c.GetFieldCol())*FIELD_SIZE_PX), Y: int32(float64(c.GetFieldRow())*FIELD_SIZE_PX),
        W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
    renderer.Copy(tex.Texture, nil, &rect)

    // Render head
    tex = TEXTURES[TEXTURE_FILENAME_ROCKETTOWER_HEAD]
    rect = sdl.Rect{
        X: int32(float64(c.GetFieldCol())*FIELD_SIZE_PX), Y: int32(float64(c.GetFieldRow())*FIELD_SIZE_PX),
        W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
    renderer.CopyEx(tex.Texture, nil, &rect, c.GetRotationDeg(), nil, 0)
}


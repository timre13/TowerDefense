package tower

import (
    "github.com/veandco/go-sdl2/sdl"
    "math"
    "fmt"
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

func towerAsRect(t ITower) sdl.Rect {
    return sdl.Rect{
        X: int32(float64(t.GetFieldCol())*FIELD_SIZE_PX),
        Y: int32(float64(t.GetFieldRow())*FIELD_SIZE_PX),
        W: int32(FIELD_SIZE_PX),
        H: int32(FIELD_SIZE_PX)}
}

func renderTowerInfo(renderer *sdl.Renderer, t ITower, x int32, y int32) {
    charX := x+20
    for _, char := range fmt.Sprint(t.GetHP()) {
        tex := CHAR_TEXTURES[char-'!']
        rect := sdl.Rect{
            X: charX, Y: y,
            W: tex.Width, H: tex.Height}
        renderer.Copy(tex.Texture, nil, &rect)
        charX += tex.Width
    }
}

func renderTower(renderer *sdl.Renderer, t ITower, bodyTexName string, headTexName string) {
    // Render body
    tex := TEXTURES[bodyTexName]
    rect := sdl.Rect{
        X: int32(float64(t.GetFieldCol())*FIELD_SIZE_PX), Y: int32(float64(t.GetFieldRow())*FIELD_SIZE_PX),
        W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
    renderer.Copy(tex.Texture, nil, &rect)

    // Render head
    tex = TEXTURES[headTexName]
    rect = sdl.Rect{
        X: int32(float64(t.GetFieldCol())*FIELD_SIZE_PX), Y: int32(float64(t.GetFieldRow())*FIELD_SIZE_PX),
        W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
    renderer.CopyEx(tex.Texture, nil, &rect, t.GetRotationDeg(), nil, 0)
}

func towerUpdateRotation(t ITower, enemies []enemy.IEnemy) {
    cloCol, cloRow := enemy.GetClosestEnemyPos(enemies, float64(t.GetFieldCol()), float64(t.GetFieldRow()))
    rotRad := math.Atan2(float64(cloRow - t.GetFieldRow()), float64(cloCol - t.GetFieldCol()))
    rotDeg := rotRad * (180/math.Pi) + 90
    t.SetRotationDeg(t.GetRotationDeg() + (rotDeg - t.GetRotationDeg()) / 10)
}

func towerCheckCursorHover(t ITower, renderer *sdl.Renderer, x int32, y int32) {
    rect := towerAsRect(t)
    if IsInsideRect(rect, x, y) {
        renderer.SetDrawColor(255, 255, 255, 255)
        renderer.DrawRect(&rect)
        renderTowerInfo(renderer, t, x, y)
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

    CheckCursorHover(renderer *sdl.Renderer, x int32, y int32)

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

func (c *Cannon) CheckCursorHover(renderer *sdl.Renderer, x int32, y int32) {
    towerCheckCursorHover(c, renderer, x, y)
}

func (c *Cannon) Update(enemies []enemy.IEnemy) {
    towerUpdateRotation(c, enemies)
}

func (c *Cannon) Render(renderer *sdl.Renderer) {
    renderTower(renderer, c, TEXTURE_FILENAME_CANNON_BASE, TEXTURE_FILENAME_CANNON_HEAD)
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

func (t *RocketTower) GetFieldCol() int32 { return t.FieldCol }
func (t *RocketTower) GetFieldRow() int32 { return t.FieldRow }
func (t *RocketTower) GetHP() int { return t.Hp; }
func (t *RocketTower) IsReal() bool { return t.IsReal_ }
func (t *RocketTower) GetRotationDeg() float64 { return t.RotationDeg }

func (t *RocketTower) SetReal(val bool) { t.IsReal_ = val }
func (t *RocketTower) SetRotationDeg(val float64) { t.RotationDeg = val }

func (t *RocketTower) CheckCursorHover(renderer *sdl.Renderer, x int32, y int32) {
    towerCheckCursorHover(t, renderer, x, y)
}

func (t *RocketTower) Update(enemies []enemy.IEnemy) {
    towerUpdateRotation(t, enemies)
}

func (t *RocketTower) Render(renderer *sdl.Renderer) {
    renderTower(renderer, t, TEXTURE_FILENAME_ROCKETTOWER_BASE, TEXTURE_FILENAME_ROCKETTOWER_HEAD)
}


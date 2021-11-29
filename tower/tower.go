package tower

import (
    "github.com/veandco/go-sdl2/sdl"
    "math"
    "math/rand"
    "fmt"
    . "TowerDefense/common"
    "TowerDefense/enemy"
    "TowerDefense/missile"
)

//------------------------------------------------------------------------------

type TowerType int
const (
    TOWER_TYPE_NONE TowerType = iota
    TOWER_TYPE_CANNON
    TOWER_TYPE_ROCKETTOWER
    TOWER_TYPE__COUNT
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
    for _, char := range fmt.Sprintf("HP: %d", t.GetHP()) {
        tex := GetCharTex(char)
        rect := sdl.Rect{
            X: charX, Y: y,
            W: tex.Width, H: tex.Height}
        renderer.Copy(tex.Texture, nil, &rect)
        charX += tex.Width
    }
}

func renderTower(renderer *sdl.Renderer, t ITower, bodyTexName string, headTexName string) {
    // ----- Render body -----

    tex := TEXTURES[bodyTexName]

    if t.IsPreview() {
        tex.Texture.SetAlphaMod(128)
        tex.Texture.SetColorMod(50, 50, 50)
    } else {
        tex.Texture.SetAlphaMod(255)
        tex.Texture.SetColorMod(255, 255, 255)
    }

    rect := sdl.Rect{
        X: int32(float64(t.GetFieldCol())*FIELD_SIZE_PX), Y: int32(float64(t.GetFieldRow())*FIELD_SIZE_PX),
        W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
    renderer.Copy(tex.Texture, nil, &rect)


    // ----- Render head -----

    tex = TEXTURES[headTexName]

    if t.IsPreview() {
        tex.Texture.SetAlphaMod(128)
        tex.Texture.SetColorMod(50, 50, 50)
    } else {
        tex.Texture.SetAlphaMod(255)
        tex.Texture.SetColorMod(255, 255, 255)
    }

    rect = sdl.Rect{
        X: int32(float64(t.GetFieldCol())*FIELD_SIZE_PX), Y: int32(float64(t.GetFieldRow())*FIELD_SIZE_PX),
        W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
    renderer.CopyEx(tex.Texture, nil, &rect, t.GetRotationDeg(), nil, 0)
}

func towerDoUpdate(t ITower, enemies []enemy.IEnemy, missiles *[]missile.IMissile) {
    // ----- Update rotation -----
    cloCol, cloRow := enemy.GetClosestEnemyPos(enemies, float64(t.GetFieldCol()), float64(t.GetFieldRow()))
    rotRad := math.Atan2(float64(cloRow - t.GetFieldRow()), float64(cloCol - t.GetFieldCol()))
    rotDeg := RadToDeg(rotRad)
    t.SetRotationDeg(t.GetRotationDeg() + (rotDeg - t.GetRotationDeg()) / 10)

    // ----- Spawn missile if needed -----
    if rand.Int() % 120 == 0 {
        miss := missile.CannonBall{
                XPos: int32(float64(t.GetFieldCol())*FIELD_SIZE_PX+FIELD_SIZE_PX/4),
                YPos: int32(float64(t.GetFieldRow())*FIELD_SIZE_PX+FIELD_SIZE_PX/4),
                RotationRad: DegToRad(t.GetRotationDeg()),
                Speed: 10}
        *missiles = append(*missiles, &miss)
    }
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
    IsPreview() bool
    GetRotationDeg() float64

    SetRotationDeg(val float64)
    SetFieldCol(val int32)
    SetFieldRow(val int32)

    CheckCursorHover(renderer *sdl.Renderer, x int32, y int32)

    Update(enemies []enemy.IEnemy, missiles *[]missile.IMissile)

    Render(renderer *sdl.Renderer)
}

//------------------------------------------------------------------------------

type Cannon struct {
    FieldCol int32;
    FieldRow int32;

    IsPreview_ bool;

    RotationDeg float64

    Hp int
}
var _ ITower = (*Cannon)(nil)

func (c *Cannon) GetFieldCol() int32 { return c.FieldCol }
func (c *Cannon) GetFieldRow() int32 { return c.FieldRow }
func (c *Cannon) GetHP() int { return c.Hp; }
func (c *Cannon) IsPreview() bool { return c.IsPreview_ }
func (c *Cannon) GetRotationDeg() float64 { return c.RotationDeg }

func (c *Cannon) SetRotationDeg(val float64) { c.RotationDeg = val }

func (c *Cannon) SetFieldCol(val int32) { c.FieldCol = val }
func (c *Cannon) SetFieldRow(val int32) { c.FieldRow = val }

func (c *Cannon) CheckCursorHover(renderer *sdl.Renderer, x int32, y int32) {
    towerCheckCursorHover(c, renderer, x, y)
}

func (c *Cannon) Update(enemies []enemy.IEnemy, missiles *[]missile.IMissile) {
    towerDoUpdate(c, enemies, missiles)
}

func (c *Cannon) Render(renderer *sdl.Renderer) {
    renderTower(renderer, c, TEXTURE_FILENAME_CANNON_BASE, TEXTURE_FILENAME_CANNON_HEAD)
}

//------------------------------------------------------------------------------

type RocketTower struct {
    FieldCol int32;
    FieldRow int32;

    IsPreview_ bool;

    RotationDeg float64

    Hp int
}
var _ ITower = (*RocketTower)(nil)

func (t *RocketTower) GetFieldCol() int32 { return t.FieldCol }
func (t *RocketTower) GetFieldRow() int32 { return t.FieldRow }
func (t *RocketTower) GetHP() int { return t.Hp; }
func (t *RocketTower) IsPreview() bool { return t.IsPreview_ }
func (t *RocketTower) GetRotationDeg() float64 { return t.RotationDeg }

func (t *RocketTower) SetRotationDeg(val float64) { t.RotationDeg = val }
func (t *RocketTower) SetFieldCol(val int32) { t.FieldCol = val }
func (t *RocketTower) SetFieldRow(val int32) { t.FieldRow = val }

func (t *RocketTower) CheckCursorHover(renderer *sdl.Renderer, x int32, y int32) {
    towerCheckCursorHover(t, renderer, x, y)
}

func (t *RocketTower) Update(enemies []enemy.IEnemy, missiles *[]missile.IMissile) {
    towerDoUpdate(t, enemies, missiles)
}

func (t *RocketTower) Render(renderer *sdl.Renderer) {
    renderTower(renderer, t, TEXTURE_FILENAME_ROCKETTOWER_BASE, TEXTURE_FILENAME_ROCKETTOWER_HEAD)
}


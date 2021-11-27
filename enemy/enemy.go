package enemy

import (
    "github.com/veandco/go-sdl2/sdl"
    "math"
    . "TowerDefense/common"
)

//-------------------------------------------------------------------------------

func renderEnemy(renderer *sdl.Renderer, e IEnemy) {
    tex := TEXTURES[TEXTURE_FILENAME_TANK]
    rect := sdl.Rect{
        X: int32(float64(e.GetFieldCol())*FIELD_SIZE_PX), Y: int32(float64(e.GetFieldRow())*FIELD_SIZE_PX),
        W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
    renderer.Copy(tex.Texture, nil, &rect)
}

//-------------------------------------------------------------------------------

type IEnemy interface {
    GetFieldCol() int32
    GetFieldRow() int32
    GetHP() int

    Update()

    Render(renderer *sdl.Renderer)
}

//-------------------------------------------------------------------------------

type Tank struct {
    FieldCol int32;
    FieldRow int32;

    Hp int
}
var _ IEnemy = (*Tank)(nil)

func (t *Tank) GetFieldCol() int32 { return t.FieldCol }
func (t *Tank) GetFieldRow() int32 { return t.FieldRow }
func (t *Tank) GetHP() int { return t.Hp }

func (t *Tank) Update() {
    // TODO
}

func (t *Tank) Render(renderer *sdl.Renderer) {
    renderEnemy(renderer, t)
}


//-------------------------------------------------------------------------------

func calcDistance(a Vec2DF, b Vec2DF) float64 {
    xLen := math.Abs(a.X - b.X)
    yLen := math.Abs(a.Y - b.Y)
    return math.Sqrt(xLen*xLen + yLen*yLen)
}

func GetClosestEnemyPos(enemies []IEnemy, col float64, row float64) (int32, int32) {
    if len(enemies) == 0 {
        return -1, -1
    }

    closestDist := -1.0
    closestI := -1
    for i, enemy := range enemies {
        dist := calcDistance(
            Vec2DF{X: col, Y: row},
            Vec2DF{X: float64(enemy.GetFieldCol()), Y: float64(enemy.GetFieldRow())})
        // If this is the first checked or closer than the closest one
        if closestDist < 0 || dist < closestDist {
            closestDist = dist
            closestI = i
        }
    }
    return enemies[closestI].GetFieldCol(), enemies[closestI].GetFieldRow()
}


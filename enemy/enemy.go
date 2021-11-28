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
        X: e.GetXPos(), Y: e.GetYPos(),
        W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
    renderer.CopyEx(tex.Texture, nil, &rect, e.GetRotationDeg(), nil, 0)
}

//-------------------------------------------------------------------------------

type IEnemy interface {
    GetFieldCol() int32
    GetFieldRow() int32
    GetXPos() int32
    GetYPos() int32
    GetHP() int
    GetRotationDeg() float64
    HasArrivedToDestination() bool

    setRotationDeg(val float64)

    Update()

    Render(renderer *sdl.Renderer)
}

//-------------------------------------------------------------------------------

type Tank struct {
    Hp int
    RotationDeg float64

    roadI int
    roadOffset int

}
var _ IEnemy = (*Tank)(nil)

func (t *Tank) GetFieldCol() int32 { return ROAD_COORDS[t.roadI].X }
func (t *Tank) GetFieldRow() int32 { return ROAD_COORDS[t.roadI].Y }
func (t *Tank) GetHP() int { return t.Hp }
func (t *Tank) GetRotationDeg() float64 { return t.RotationDeg }
func (t *Tank) GetXPos() int32 {
    if t.roadI == len(ROAD_COORDS)-1 {
        return ROAD_COORDS[t.roadI].X
    } else {
        return int32(float64(
            Lerp(float64(ROAD_COORDS[t.roadI].X)*FIELD_SIZE_PX,
                 float64(ROAD_COORDS[t.roadI+1].X)*FIELD_SIZE_PX, float64(t.roadOffset)/100.0)))
    }
}

func (t *Tank) GetYPos() int32 {
    if t.roadI == len(ROAD_COORDS)-1 {
        return ROAD_COORDS[t.roadI].Y
    } else {
        return int32(float64(
            Lerp(float64(ROAD_COORDS[t.roadI].Y)*FIELD_SIZE_PX,
                 float64(ROAD_COORDS[t.roadI+1].Y)*FIELD_SIZE_PX, float64(t.roadOffset)/100.0)))
    }
}

func (t *Tank) HasArrivedToDestination() bool {
    return t.roadI >= len(ROAD_COORDS)-1
}

func (t *Tank) setRotationDeg(val float64) { t.RotationDeg = val }

func (t *Tank) Update() {
    t.roadOffset += 3
    if t.roadOffset >= 100 {
        t.roadI++
        if t.roadI >= len(ROAD_COORDS) {
            t.roadI = len(ROAD_COORDS)-1
            return
        }

        if t.roadI != len(ROAD_COORDS)-1 {
            col := ROAD_COORDS[t.roadI].X
            row := ROAD_COORDS[t.roadI].Y
   
            dir := t.GetRotationDeg()
            xDiff, yDiff := ROAD_COORDS[t.roadI+1].X-col, ROAD_COORDS[t.roadI+1].Y-row

            if xDiff > 0 {
                dir = 90
            } else if xDiff < 0 {
                dir = 270
            } else if yDiff > 0 {
                dir = 180
            } else if yDiff < 0 {
                dir = 0
            }
            t.setRotationDeg(dir)
        }

        t.roadOffset = 0
    }
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


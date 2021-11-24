package main

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/img"
    "github.com/veandco/go-sdl2/ttf"
    "fmt"
    "math"
    "os"
)

func ASSERT_TRUE(val bool) {
    if !val {
        panic("Assertion failed")
    }
}

const FONT_FILE_PATH = "/usr/share/fonts/truetype/freefont/FreeSans.ttf"
const DEF_FONT_SIZE = 40
var FIELD_SIZE_PX float64

var (
    COLOR_BG_1                  = sdl.Color{R:  64, G: 149, B:  64, A: 255}
    COLOR_BG_2                  = sdl.Color{R:  57, G: 161, B:  50, A: 255}
    COLOR_BG_ROAD               = sdl.Color{R: 131, G: 113, B:  95, A: 255}
    COLOR_BOTBAR_BG             = sdl.Color{R: 130, G: 130, B: 130, A: 255}
    COLOR_BOTBAR_BORDER         = sdl.Color{R: 100, G: 100, B: 100, A: 255}
)

const (
    MAP_WIDTH_FIELD         = 16
    MAP_HEIGHT_FIELD        = 8
    BOTTOM_BAR_HEIGHT_PX    = 100
)

type TowerType int
const (
    TOWER_TYPE_NONE TowerType = iota
    TOWER_TYPE_CANNON
)

func (t TowerType) getPrice() int {
    switch (t) {
    default: panic(t)
    case TOWER_TYPE_NONE:       panic(t)
    case TOWER_TYPE_CANNON:     return 10
    }
}

func (t TowerType) getInitialHP() int {
    switch (t) {
    default: panic(t)
    case TOWER_TYPE_NONE:       panic(t)
    case TOWER_TYPE_CANNON:     return 100
    }
}


type Vec2D struct {
    x int;
    y int;
}

type Vec2DF struct {
    x float64;
    y float64;
}

var ROAD_COORDS = [...]Vec2D{
    { 1,  0}, { 1,  1}, { 1,  2}, { 1,  3}, { 1,  6}, { 1,  7}, { 1,  8}, { 1,  9}, {10,  1},
    {10,  2}, {11, 11}, {11, 12}, {11, 13}, {11, 14}, {11, 15}, {11,  2}, {12, 11}, {12, 15},
    {12, 16}, {12, 17}, {12, 18}, {12,  2}, {12,  4}, {12,  5}, {12,  6}, {12,  7}, {13, 10},
    {13, 11}, {13, 18}, {13,  1}, {13,  2}, {13,  4}, {13,  7}, {13,  9}, {14, 18}, {14,  1},
    {14,  4}, {14,  7}, {14,  8}, {14,  9}, {15, 11}, {15, 12}, {15, 13}, {15, 16}, {15, 17},
    {15, 18}, {15,  1}, {15,  2}, {15,  3}, {15,  4}, {16, 11}, {16, 13}, {16, 14}, {16, 15},
    {16, 16}, {17, 10}, {17, 11}, {17,  4}, {17,  5}, {17,  6}, {17,  8}, {17,  9}, {18,  2},
    {18,  3}, {18,  4}, {18,  6}, {18,  7}, {18,  8}, {19,  2}, { 2, 10}, { 2, 11}, { 2, 12},
    { 2, 16}, { 2, 17}, { 2, 18}, { 2, 19}, { 2,  3}, { 2,  4}, { 2,  5}, { 2,  6}, { 2,  9},
    {20,  2}, {21,  2}, {21,  3}, {22, 11}, {22, 12}, {22, 13}, {22, 14}, {22, 15}, {22,  3},
    {23, 11}, {23, 15}, {23,  3}, {24, 10}, {24, 11}, {24, 15}, {24, 16}, {24, 17}, {24, 18},
    {24, 19}, {24,  1}, {24,  2}, {24,  3}, {24,  9}, {25,  1}, {25,  9}, {26,  1}, {26,  7},
    {26,  8}, {26,  9}, {27,  1}, {27,  2}, {27,  3}, {27,  4}, {27,  7}, {28,  4}, {28,  5},
    {28,  6}, {28,  7}, { 3, 12}, { 3, 13}, { 3, 15}, { 3, 16}, { 3, 19}, { 4, 13}, { 4, 14},
    { 4, 15}, { 4, 19}, { 5, 10}, { 5, 18}, { 5, 19}, { 5,  6}, { 5,  7}, { 5,  8}, { 5,  9},
    { 6, 10}, { 6, 11}, { 6, 16}, { 6, 17}, { 6, 18}, { 6,  3}, { 6,  4}, { 6,  5}, { 6,  6},
    { 7, 11}, { 7, 12}, { 7, 13}, { 7, 14}, { 7, 16}, { 7,  1}, { 7,  2}, { 7,  3}, { 8, 14},
    { 8, 15}, { 8, 16}, { 8,  1}, { 9,  1},
}

type Texture struct {
    texture *sdl.Texture
    width int32
    height int32
}

const TEXTURE_DIR_PATH = "img"

const (
    TEXTURE_FILENAME_TANK               = "tank/tank.png"
    TEXTURE_FILENAME_CANNON_BASE        = "cannon/base.png"
    TEXTURE_FILENAME_CANNON_HEAD        = "cannon/head.png"
    TEXTURE_FILENAME_COIN               = "coin/coin.png"
)

var TEXTURES = map[string]*Texture{
    TEXTURE_FILENAME_TANK:              nil,
    TEXTURE_FILENAME_CANNON_BASE:       nil,
    TEXTURE_FILENAME_CANNON_HEAD:       nil,
    TEXTURE_FILENAME_COIN:              nil,
}

var CHAR_TEXTURES = [94]*Texture{}

//-------------------------------------------------------------------------------

func UNUSED(x ...interface{}) {}

//-------------------------------------------------------------------------------

func drawCheckerBg(renderer *sdl.Renderer) {
    renderer.SetDrawColor(COLOR_BG_1.R, COLOR_BG_1.G, COLOR_BG_1.B, COLOR_BG_1.A)
    renderer.Clear()

    renderer.SetDrawColor(COLOR_BG_2.R, COLOR_BG_2.G, COLOR_BG_2.B, COLOR_BG_2.A)
    for y :=0; y < MAP_HEIGHT_FIELD; y++ {
        for x:=y%2; x < MAP_WIDTH_FIELD; x+=2 {
            rect := sdl.Rect{X: int32(float64(x)*FIELD_SIZE_PX), Y: int32(float64(y)*FIELD_SIZE_PX),
                             W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
            renderer.FillRect(&rect)
        }
    }
}

func drawRoad(renderer *sdl.Renderer) {
    renderer.SetDrawColor(COLOR_BG_ROAD.R, COLOR_BG_ROAD.G, COLOR_BG_ROAD.B, COLOR_BG_ROAD.A)

    for _, field := range ROAD_COORDS {
        rect := sdl.Rect{X: int32(float64(field.x)*FIELD_SIZE_PX), Y: int32(float64(field.y)*FIELD_SIZE_PX),
                         W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
        renderer.FillRect(&rect)
    }
}

func drawBottomBar(renderer *sdl.Renderer, winW int32, winH int32, coins int) {
    // Draw border
    renderer.SetDrawColor(
        COLOR_BOTBAR_BORDER.R, COLOR_BOTBAR_BORDER.G, COLOR_BOTBAR_BORDER.B, COLOR_BOTBAR_BORDER.A)
    rect := sdl.Rect{X: 0, Y: winH-BOTTOM_BAR_HEIGHT_PX, W: winW, H: BOTTOM_BAR_HEIGHT_PX}
    renderer.FillRect(&rect)

    // Fill
    renderer.SetDrawColor(COLOR_BOTBAR_BG.R, COLOR_BOTBAR_BG.G, COLOR_BOTBAR_BG.B, COLOR_BOTBAR_BG.A)
    rect = sdl.Rect{X: 8, Y: winH-BOTTOM_BAR_HEIGHT_PX+8, W: winW-16, H: BOTTOM_BAR_HEIGHT_PX-16}
    renderer.FillRect(&rect)

    // Draw coin texture
    tex := TEXTURES[TEXTURE_FILENAME_COIN]
    rect = sdl.Rect{X: 18, Y: winH-BOTTOM_BAR_HEIGHT_PX+BOTTOM_BAR_HEIGHT_PX/2-32, W: 64, H: 64}
    renderer.Copy(tex.texture, nil, &rect)

    var offs int32 = 0
    for _, char := range fmt.Sprint(coins) {
        tex := CHAR_TEXTURES[char-'!']
        rect = sdl.Rect{
            X: 100+offs, Y: winH-BOTTOM_BAR_HEIGHT_PX+BOTTOM_BAR_HEIGHT_PX/2-DEF_FONT_SIZE/2,
            W: tex.width, H: tex.height}
        renderer.Copy(tex.texture, nil, &rect)
        offs += tex.width
    }
}

//-------------------------------- Enemy ----------------------------------------

type IEnemy interface {
    getFieldCol() int32
    getFieldRow() int32
    getHP() int

    update()

    render(renderer *sdl.Renderer)
}

type Tank struct {
    fieldCol int32;
    fieldRow int32;

    hp int
}
var _ IEnemy = (*Tank)(nil)

func (t *Tank) getFieldCol() int32 { return t.fieldCol }
func (t *Tank) getFieldRow() int32 { return t.fieldRow }
func (t *Tank) getHP() int { return t.hp }

func (t *Tank) update() {
    // TODO
}

func (t *Tank) render(renderer *sdl.Renderer) {
    tex := TEXTURES[TEXTURE_FILENAME_TANK]
    rect := sdl.Rect{
        X: int32(float64(t.getFieldCol())*FIELD_SIZE_PX), Y: int32(float64(t.getFieldRow())*FIELD_SIZE_PX),
        W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
    renderer.Copy(tex.texture, nil, &rect)
}

//-------------------------------- Tower ----------------------------------------

type ITower interface {
    getFieldCol() int32
    getFieldRow() int32
    getHP() int
    /*
     * If the tower is phisically on the map, this is true.
     * If the tower is used on the bottom bar or otherwise an indicator, this is false.
    */
    isReal() bool
    getRotationDeg() float64

    setReal(val bool)
    setRotationDeg(val float64)

    update(enemies []IEnemy)

    render(renderer *sdl.Renderer)
}

type Cannon struct {
    fieldCol int32;
    fieldRow int32;

    isReal_ bool;

    rotationDeg float64

    hp int
}
var _ ITower = (*Cannon)(nil)

func (c *Cannon) getFieldCol() int32 { return c.fieldCol }
func (c *Cannon) getFieldRow() int32 { return c.fieldRow }
func (c *Cannon) getHP() int { return c.hp; }
func (c *Cannon) isReal() bool { return c.isReal_ }
func (c *Cannon) getRotationDeg() float64 { return c.rotationDeg }

func (c *Cannon) setReal(val bool) { c.isReal_ = val }
func (c *Cannon) setRotationDeg(val float64) { c.rotationDeg = val }

func (c *Cannon) update(enemies []IEnemy) {
    cloCol, cloRow := getClosestEnemyPos(enemies, float64(c.fieldCol), float64(c.fieldRow))
    rotRad := math.Atan2(float64(cloRow - c.fieldRow), float64(cloCol - c.fieldCol))
    rotDeg := rotRad * (180/math.Pi) + 90
    c.rotationDeg += (rotDeg - c.rotationDeg) / 10
}

func (c *Cannon) render(renderer *sdl.Renderer) {
    // Render body
    tex := TEXTURES[TEXTURE_FILENAME_CANNON_BASE]
    rect := sdl.Rect{
        X: int32(float64(c.getFieldCol())*FIELD_SIZE_PX), Y: int32(float64(c.getFieldRow())*FIELD_SIZE_PX),
        W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
    renderer.Copy(tex.texture, nil, &rect)

    // Render head
    tex = TEXTURES[TEXTURE_FILENAME_CANNON_HEAD]
    rect = sdl.Rect{
        X: int32(float64(c.getFieldCol())*FIELD_SIZE_PX), Y: int32(float64(c.getFieldRow())*FIELD_SIZE_PX),
        W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
    renderer.CopyEx(tex.texture, nil, &rect, c.getRotationDeg(), nil, 0)
}

//-------------------------------------------------------------------------------

func calcDistance(a Vec2DF, b Vec2DF) float64 {
    xLen := math.Abs(a.x - b.x)
    yLen := math.Abs(a.y - b.y)
    return math.Sqrt(xLen*xLen + yLen*yLen)
}

func getClosestEnemyPos(enemies []IEnemy, col float64, row float64) (int32, int32) {
    if len(enemies) == 0 {
        return -1, -1
    }

    closestDist := -1.0
    closestI := -1
    for i, enemy := range enemies {
        dist := calcDistance(
            Vec2DF{col, row},
            Vec2DF{float64(enemy.getFieldCol()), float64(enemy.getFieldRow())})
        // If this is the first checked or closer than the closest one
        if closestDist < 0 || dist < closestDist {
            closestDist = dist
            closestI = i
        }
    }
    return enemies[closestI].getFieldCol(), enemies[closestI].getFieldRow()
}


//-------------------------------------------------------------------------------

func main() {
    cVer := sdl.Version{}
    sdl.VERSION(&cVer)
    fmt.Printf("Compiled SDL version: %d.%d.%d\n", cVer.Major, cVer.Minor, cVer.Patch)

    lVer := sdl.Version{}
    sdl.GetVersion(&lVer)
    fmt.Printf("Linked SDL version: %d.%d.%d\n", lVer.Major, lVer.Minor, lVer.Patch)

    var err error;

    //------------------------------ Init --------------------------------------

    err = sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO)
    if err != nil {
        fmt.Printf("Failed to initialize SDL2: %s\n", err.Error())
        panic(err)
    }

    err = ttf.Init()
    if err != nil {
        panic(err)
    }

    const maxWinHeight = 900
    window, err := sdl.CreateWindow(
            "Tower Defense", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
            int32((maxWinHeight-BOTTOM_BAR_HEIGHT_PX)*float64(MAP_WIDTH_FIELD)/float64(MAP_HEIGHT_FIELD)), maxWinHeight,
            sdl.WINDOW_RESIZABLE)
    if err != nil {
        fmt.Printf("Failed to create window: %s\n", err.Error())
        panic(err)
    }

    renderer, err := sdl.CreateRenderer(window, 0, 0)
    renderer.SetDrawColor(100, 100, 100, 255)
    renderer.Clear()
    renderer.Present()

    fmt.Printf("Loading %d textures\n", len(TEXTURES))
    i := 0
    for fileName := range TEXTURES {
        path := TEXTURE_DIR_PATH+string(os.PathSeparator)+fileName
        fmt.Printf("[%d/%d] Loading \"%s\"\n", i+1, len(TEXTURES), path)
        surface, err := img.Load(path)
        if err != nil {
            panic(err)
        }
        texture, err := renderer.CreateTextureFromSurface(surface)
        if err != nil {
            panic(err)
        }
        tex := Texture{texture: texture, width: surface.W, height: surface.H}
        surface.Free()
        TEXTURES[fileName] = &tex
        fmt.Printf("Loaded \"%s\", size: %dx%d\n", fileName, tex.width, tex.height)
        i++
    }
    fmt.Println("Textures:", TEXTURES)

    fmt.Println("Rendering font")
    font, error := ttf.OpenFont(FONT_FILE_PATH, DEF_FONT_SIZE)
    if error != nil {
        panic(error)
    }
    for i := range CHAR_TEXTURES {
        surface, error := font.RenderGlyphBlended(rune('!'+i), sdl.Color{R: 255, G: 255, B: 255, A: 255})
        if error != nil {
            panic(error)
        }
        tex, error := renderer.CreateTextureFromSurface(surface)
        if error != nil {
            panic(error)
        }
        texture := Texture{texture: tex, width: surface.W, height: surface.H}
        surface.Free()
        CHAR_TEXTURES[i] = &texture
    }
    font.Close()
    font = nil

    //--------------------------- Variables ------------------------------------

    coins := 100
    var towers []ITower
    var enemies []IEnemy
    placedTowerType := TOWER_TYPE_CANNON

    isTowerAt := func(col int32, row int32) bool {
        for _, tower := range towers {
            if tower.getFieldCol() == col && tower.getFieldRow() == row {
                return true
            }
        }
        return false
    }

    isRoadAt := func(col int32, row int32) bool {
        for _, coord := range ROAD_COORDS {
            if int32(coord.x) == col && int32(coord.y) == row {
                return true
            }
        }
        return false
    }

    // TODO: Test -- Remove later
    tank1 := Tank{3, 3, 10}
    enemies = append(enemies, &tank1)
    tank2 := Tank{10, 15, 10}
    enemies = append(enemies, &tank2)

    //--------------------------- Main loop ------------------------------------

    fmt.Println("Setup done")

    done := false
    var startTime uint32 = 1
    var frameTime uint32 = 1
    for {
        startTime = sdl.GetTicks()
        winW, winH := window.GetSize()
        mouseX, mouseY, mouseState := sdl.GetMouseState()
        UNUSED(mouseState)

        for {
            var event = sdl.PollEvent()
            if event == nil { // No more events in the queue
                break
            }

            switch event.GetType() {
            case sdl.QUIT:
                done = true
                fmt.Println("Window close requested")

            case sdl.MOUSEBUTTONDOWN:
                col := int32(float64(mouseX)/FIELD_SIZE_PX)
                row := int32(float64(mouseY)/FIELD_SIZE_PX)
                if placedTowerType != TOWER_TYPE_NONE && coins >= placedTowerType.getPrice() &&
                        !isTowerAt(col, row) && !isRoadAt(col, row) {

                    coins -= placedTowerType.getPrice()

                    var tower ITower

                    switch (placedTowerType) {
                    case TOWER_TYPE_CANNON:
                        tower = &Cannon{
                            fieldCol: col,
                            fieldRow: row,
                            isReal_: true,
                            hp: placedTowerType.getInitialHP()}
                    default: panic(placedTowerType)
                    }

                    towers = append(towers, tower)

                    fmt.Printf("Placed a tower at {%d, %d}\n", col, row)
                }
                //fmt.Printf("{%d, %d}\n", col, row);
            }
        }
        if done {
            break
        }

        window.SetTitle(
            fmt.Sprintf("Tower Defense :: FT: %dms, FPS: %f", frameTime, 1000/float32(frameTime)))

        // TODO: Only do it on window resizing
        FIELD_SIZE_PX = math.Min(float64(winW)/MAP_WIDTH_FIELD,
                                float64(winH-BOTTOM_BAR_HEIGHT_PX)/MAP_HEIGHT_FIELD)

        // Render environment
        drawCheckerBg(renderer)
        drawRoad(renderer)
        ASSERT_TRUE(coins >= 0)
        drawBottomBar(renderer, winW, winH, coins)

        // Update entities
        for _, enemy := range enemies { enemy.update() }
        for _, tower := range towers { tower.update(enemies) }

        // Render entities
        for _, enemy := range enemies { enemy.render(renderer) }
        for _, tower := range towers { tower.render(renderer) }

        renderer.Present()
        sdl.Delay(16)

        frameTime = sdl.GetTicks() - startTime
    }

    //----------------------------- Cleanup ------------------------------------

    for _, texture := range TEXTURES {
        texture.texture.Destroy()
    }

    for _, texture := range CHAR_TEXTURES {
        texture.texture.Destroy()
    }

    renderer.Destroy()
    window.Destroy()
    sdl.Quit()
    ttf.Quit()

    fmt.Println("Window closed")
}

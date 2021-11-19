package main

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/img"
    "fmt"
    "math"
    "os"
)

var (
    COLOR_BG_1                  = sdl.Color{R:  64, G: 149, B:  64, A: 255}
    COLOR_BG_2                  = sdl.Color{R:  57, G: 161, B:  50, A: 255}
    COLOR_BG_ROAD               = sdl.Color{R: 131, G: 113, B:  95, A: 255}
    COLOR_BOTBAR_BG             = sdl.Color{R: 130, G: 130, B: 130, A: 255}
    COLOR_BOTBAR_BORDER         = sdl.Color{R: 100, G: 100, B: 100, A: 255}
)

const (
    MAP_WIDTH_FIELD = 30
    MAP_HEIGHT_FIELD = 20
    BOTTOM_BAR_HEIGHT_PX = 100
)

type Vec2D struct {
    x int;
    y int;
}

var ROAD_COORDS = [...]Vec2D{
    { 2,  0}, { 2,  1}, { 8,  1}, { 9,  1}, {10,  1}, {11,  1}, {12,  1}, {13,  1},
    { 2,  2}, { 7,  2}, { 8,  2}, {13,  2}, {14,  2}, {24,  2}, {25,  2}, {26,  2},
    { 2,  3}, { 7,  3}, {14,  3}, {21,  3}, {22,  3}, {23,  3}, {24,  3}, {26,  3},
    { 2,  4}, { 7,  4}, {14,  4}, {20,  4}, {21,  4}, {26,  4}, { 2,  5}, { 7,  5},
    {14,  5}, {19,  5}, {20,  5}, {26,  5}, { 2,  6}, { 6,  6}, { 7,  6}, {14,  6},
    {19,  6}, {26,  6}, { 2,  7}, { 6,  7}, {14,  7}, {19,  7}, {26,  7}, { 2,  8},
    { 6,  8}, {14,  8}, {19,  8}, {26,  8}, { 2,  9}, { 6,  9}, {14,  9}, {19,  9},
    {26,  9}, { 2, 10}, { 6, 10}, {14, 10}, {19, 10}, {26, 10}, {27, 10}, { 2, 11},
    { 3, 11}, { 6, 11}, { 7, 11}, {13, 11}, {14, 11}, {19, 11}, {27, 11}, { 3, 12},
    { 7, 12}, {12, 12}, {13, 12}, {19, 12}, {27, 12}, { 3, 13}, { 7, 13}, {11, 13},
    {12, 13}, {18, 13}, {19, 13}, {27, 13}, { 3, 14}, { 7, 14}, {11, 14}, {18, 14},
    {27, 14}, { 3, 15}, { 7, 15}, {11, 15}, {12, 15}, {18, 15}, {27, 15}, {28, 15},
    { 3, 16}, { 7, 16}, {12, 16}, {17, 16}, {18, 16}, {28, 16}, { 3, 17}, { 6, 17},
    { 7, 17}, {12, 17}, {13, 17}, {16, 17}, {17, 17}, {28, 17}, { 3, 18}, { 4, 18},
    { 5, 18}, { 6, 18}, {13, 18}, {14, 18}, {15, 18}, {16, 18}, {28, 18}, {29, 18},
    {29, 19},
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

//-------------------------------------------------------------------------------

func UNUSED(x ...interface{}) {}

//-------------------------------------------------------------------------------

func drawCheckerBg(renderer *sdl.Renderer, fieldSizePx float64) {
    renderer.SetDrawColor(COLOR_BG_1.R, COLOR_BG_1.G, COLOR_BG_1.B, COLOR_BG_1.A)
    renderer.Clear()

    renderer.SetDrawColor(COLOR_BG_2.R, COLOR_BG_2.G, COLOR_BG_2.B, COLOR_BG_2.A)
    for y :=0; y < MAP_HEIGHT_FIELD; y++ {
        for x:=y%2; x < MAP_WIDTH_FIELD; x+=2 {
            rect := sdl.Rect{X: int32(float64(x)*fieldSizePx), Y: int32(float64(y)*fieldSizePx),
                             W: int32(fieldSizePx), H: int32(fieldSizePx)}
            renderer.FillRect(&rect)
        }
    }
}

func drawRoad(renderer *sdl.Renderer, fieldSizePx float64) {
    renderer.SetDrawColor(COLOR_BG_ROAD.R, COLOR_BG_ROAD.G, COLOR_BG_ROAD.B, COLOR_BG_ROAD.A)

    for _, field := range ROAD_COORDS {
        rect := sdl.Rect{X: int32(float64(field.x)*fieldSizePx), Y: int32(float64(field.y)*fieldSizePx),
                         W: int32(fieldSizePx), H: int32(fieldSizePx)}
        renderer.FillRect(&rect)
    }
}

func drawBottomBar(renderer *sdl.Renderer, winW int32, winH int32) {
    // Draw border
    renderer.SetDrawColor(COLOR_BOTBAR_BORDER.R, COLOR_BOTBAR_BORDER.G, COLOR_BOTBAR_BORDER.B, COLOR_BOTBAR_BORDER.A)
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

}

//-------------------------------- Enemy ----------------------------------------

type IEnemy interface {
    getXPos() int32
    getYPos() int32
    getHP() int

    getTextureName() string

    render(renderer *sdl.Renderer)
}

type Tank struct {
    xPos int32;
    yPos int32;

    hp int
}
var _ IEnemy = (*Tank)(nil)

func (t *Tank) getXPos() int32 { return t.xPos; }
func (t *Tank) getYPos() int32 { return t.yPos; }
func (t *Tank) getTextureName() string { return TEXTURE_FILENAME_TANK; }
func (t *Tank) getHP() int { return t.hp; }

func (t *Tank) render(renderer *sdl.Renderer) {
    tex := TEXTURES[t.getTextureName()]
    rect := sdl.Rect{X: t.getXPos(), Y: t.getYPos(), W: tex.width, H: tex.height}
    renderer.Copy(tex.texture, nil, &rect)
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

    window, err := sdl.CreateWindow(
            "Tower Defense", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
            MAP_WIDTH_FIELD*40, MAP_HEIGHT_FIELD*40+BOTTOM_BAR_HEIGHT_PX, sdl.WINDOW_RESIZABLE)
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
        TEXTURES[fileName] = &tex
        fmt.Printf("Loaded \"%s\", size: %dx%d\n", fileName, tex.width, tex.height)
        i++
    }
    fmt.Println("Textures:", TEXTURES)

    //--------------------------- Main loop ------------------------------------

    done := false
    var startTime uint32 = 1
    var frameTime uint32 = 1
    for {
        startTime = sdl.GetTicks()

        for {
            var event = sdl.PollEvent()
            if event == nil { // No more events in the queue
                break
            }

            switch event.GetType() {
            case sdl.QUIT:
                done = true
                fmt.Println("Window close requested")
            }
        }
        if done {
            break
        }
        winW, winH := window.GetSize()

        window.SetTitle(
            fmt.Sprintf("Tower Defense :: FT: %dms, FPS: %f", frameTime, 1000/float32(frameTime)))


        fieldSizePx := math.Min(float64(winW)/MAP_WIDTH_FIELD,
                                float64(winH-BOTTOM_BAR_HEIGHT_PX)/MAP_HEIGHT_FIELD)
        drawCheckerBg(renderer, fieldSizePx)
        drawRoad(renderer, fieldSizePx)
        drawBottomBar(renderer, winW, winH)


        renderer.Present()
        sdl.Delay(16)

        frameTime = sdl.GetTicks() - startTime
    }

    //----------------------------- Cleanup ------------------------------------

    for _, texture := range TEXTURES {
        texture.texture.Destroy()
    }

    renderer.Destroy()
    window.Destroy()
    sdl.Quit()
    ttf.Quit()

    fmt.Println("Window closed")
}

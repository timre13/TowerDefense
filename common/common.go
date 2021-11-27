package common

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/ttf"
    "runtime"
    "strings"
)

//-------------------------------------------------------------------------------

func UNUSED(x ...interface{}) {}

func ASSERT_TRUE(val bool) {
    if !val {
        panic("Assertion failed")
    }
}

//-------------------------------------------------------------------------------

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

//-------------------------------------------------------------------------------

type Vec2D struct {
    X int;
    Y int;
}

type Vec2DF struct {
    X float64;
    Y float64;
}

//-------------------------------------------------------------------------------

var ROAD_COORDS = [...]Vec2D{
    { 1,  0}, { 1,  1}, { 1,  2}, { 1,  3}, { 1,  4}, { 1,  5}, { 1,  6}, { 2,  6}, { 3,  6}, { 3,  5},
    { 3,  4}, { 3,  3}, { 3,  2}, { 4,  2}, { 4,  1}, { 5,  1}, { 6,  1}, { 6,  2}, { 6,  3}, { 6,  4},
    { 5,  4}, { 5,  5}, { 5,  6}, { 6,  6}, { 7,  6}, { 8,  6}, { 8,  5}, { 8,  4}, { 8,  3}, { 8,  2},
    { 9,  2}, { 9,  1}, {10,  1}, {11,  1}, {11,  2}, {11,  3}, {11,  4}, {10,  4}, {10,  5}, {10,  6},
    {11,  6}, {12,  6}, {13,  6}, {13,  5}, {13,  4}, {13,  3}, {14,  3}, {14,  2}, {14,  1}, {14,  0},
}

//-------------------------------------------------------------------------------

type Texture struct {
    Texture *sdl.Texture
    Width int32
    Height int32
}

const TEXTURE_DIR_PATH = "img"

const (
    TEXTURE_FILENAME_COIN               = "coin/coin.png"
    TEXTURE_FILENAME_HP                 = "hp/hp.png"
    TEXTURE_FILENAME_CANNON_BASE        = "cannon/base.png"
    TEXTURE_FILENAME_CANNON_HEAD        = "cannon/head.png"
    TEXTURE_FILENAME_ROCKETTOWER_BASE   = "rocket_tower/base.png"
    TEXTURE_FILENAME_ROCKETTOWER_HEAD   = "rocket_tower/head.png"
    TEXTURE_FILENAME_TANK               = "tank/tank.png"
    TEXTURE_FILENAME_CANNONBALL         = "cannonball/cannonball.png"
)

var TEXTURES = map[string]*Texture{
    TEXTURE_FILENAME_COIN:              nil,
    TEXTURE_FILENAME_HP:                nil,
    TEXTURE_FILENAME_CANNON_BASE:       nil,
    TEXTURE_FILENAME_CANNON_HEAD:       nil,
    TEXTURE_FILENAME_ROCKETTOWER_BASE:  nil,
    TEXTURE_FILENAME_ROCKETTOWER_HEAD:  nil,
    TEXTURE_FILENAME_TANK:              nil,
    TEXTURE_FILENAME_CANNONBALL:        nil,
}

//-------------------------------------------------------------------------------

var charTextures = [95]*Texture{}


func OpenFont(renderer *sdl.Renderer, path string) {
    err := ttf.Init()
    CheckErr(err)

    font, err := ttf.OpenFont(path, DEF_FONT_SIZE)
    CheckErr(err)

    for i := range charTextures {
        surface, err := font.RenderGlyphBlended(rune(' '+i), sdl.Color{R: 255, G: 255, B: 255, A: 255})
        CheckErr(err)
        tex, err := renderer.CreateTextureFromSurface(surface)
        CheckErr(err)
        texture := Texture{Texture: tex, Width: surface.W, Height: surface.H}
        surface.Free()
        charTextures[i] = &texture
    }
    font.Close()
    font = nil

    ttf.Quit()
}

func GetCharTex(c rune) *Texture {
    return charTextures[c-' ']
}

func FreeCharTextures() {
    for _, texture := range charTextures {
        texture.Texture.Destroy()
    }
}

//-------------------------------------------------------------------------------

func IsInsideRect(r sdl.Rect, x int32, y int32) bool {
    return x >= r.X && x < r.X+r.W && y >= r.Y && y < r.Y+r.H
}

func IsInsideWorld(x int32, y int32) bool {
    return x >= 0 && float64(x) < MAP_WIDTH_FIELD*FIELD_SIZE_PX &&
           y >= 0 && float64(y) < MAP_HEIGHT_FIELD*FIELD_SIZE_PX
}

//-------------------------------------------------------------------------------

func ShowErrAndPanic(err string) {
    // Get the stack trace
    buf := make([]byte, 1 << 16)
    runtime.Stack(buf, true)

    // Show error dialog
    msg := "Error: " + err + "\n\nStack trace: \n" + strings.ReplaceAll(string(buf), "\t", "    ")
    sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error", msg, nil)

    panic(err)
}

func CheckErr(err error) {
    if err != nil {
        ShowErrAndPanic(err.Error())
    }
}

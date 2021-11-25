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

var CHAR_TEXTURES = [94]*Texture{}


func OpenFont(renderer *sdl.Renderer, path string) {
    err := ttf.Init()
    CheckErr(err)

    font, err := ttf.OpenFont(path, DEF_FONT_SIZE)
    CheckErr(err)
    for i := range CHAR_TEXTURES {
        surface, err := font.RenderGlyphBlended(rune('!'+i), sdl.Color{R: 255, G: 255, B: 255, A: 255})
        CheckErr(err)
        tex, err := renderer.CreateTextureFromSurface(surface)
        CheckErr(err)
        texture := Texture{Texture: tex, Width: surface.W, Height: surface.H}
        surface.Free()
        CHAR_TEXTURES[i] = &texture
    }
    font.Close()
    font = nil

    ttf.Quit()
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

package main

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/img"
    "fmt"
    "math"
    "os"
    . "TowerDefense/common"
    "TowerDefense/tower"
    "TowerDefense/enemy"
)


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
        rect := sdl.Rect{X: int32(float64(field.X)*FIELD_SIZE_PX), Y: int32(float64(field.Y)*FIELD_SIZE_PX),
                         W: int32(FIELD_SIZE_PX), H: int32(FIELD_SIZE_PX)}
        renderer.FillRect(&rect)
    }
}

func drawBottomBar(renderer *sdl.Renderer, winW int32, winH int32, coins int, hp int) {
    // Draw border
    renderer.SetDrawColor(
        COLOR_BOTBAR_BORDER.R, COLOR_BOTBAR_BORDER.G, COLOR_BOTBAR_BORDER.B, COLOR_BOTBAR_BORDER.A)
    rect := sdl.Rect{X: 0, Y: winH-BOTTOM_BAR_HEIGHT_PX, W: winW, H: BOTTOM_BAR_HEIGHT_PX}
    renderer.FillRect(&rect)

    // Fill
    renderer.SetDrawColor(COLOR_BOTBAR_BG.R, COLOR_BOTBAR_BG.G, COLOR_BOTBAR_BG.B, COLOR_BOTBAR_BG.A)
    rect = sdl.Rect{X: 8, Y: winH-BOTTOM_BAR_HEIGHT_PX+8, W: winW-16, H: BOTTOM_BAR_HEIGHT_PX-16}
    renderer.FillRect(&rect)

    var offs int32 = 16

    // Draw coin texture
    tex := TEXTURES[TEXTURE_FILENAME_COIN]
    rect = sdl.Rect{X: offs, Y: winH-BOTTOM_BAR_HEIGHT_PX+BOTTOM_BAR_HEIGHT_PX/2-32, W: 64, H: 64}
    renderer.Copy(tex.Texture, nil, &rect)

    // Render coin value
    offs += 64+5
    for _, char := range fmt.Sprint(coins) {
        tex := CHAR_TEXTURES[char-'!']
        rect = sdl.Rect{
            X: offs, Y: winH-BOTTOM_BAR_HEIGHT_PX+BOTTOM_BAR_HEIGHT_PX/2-DEF_FONT_SIZE/2,
            W: tex.Width, H: tex.Height}
        renderer.Copy(tex.Texture, nil, &rect)
        offs += tex.Width
    }

    // Draw HP texture
    offs += 50
    tex = TEXTURES[TEXTURE_FILENAME_HP]
    rect = sdl.Rect{X: offs, Y: winH-BOTTOM_BAR_HEIGHT_PX+BOTTOM_BAR_HEIGHT_PX/2-32, W: 64, H: 64}
    renderer.Copy(tex.Texture, nil, &rect)

    // Render HP value
    offs += 64+5
    for _, char := range fmt.Sprint(hp) {
        tex := CHAR_TEXTURES[char-'!']
        rect = sdl.Rect{
            X: offs, Y: winH-BOTTOM_BAR_HEIGHT_PX+BOTTOM_BAR_HEIGHT_PX/2-DEF_FONT_SIZE/2,
            W: tex.Width, H: tex.Height}
        renderer.Copy(tex.Texture, nil, &rect)
        offs += tex.Width
    }
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
    CheckErr(err)

    const maxWinHeight = 900
    window, err := sdl.CreateWindow(
            "Tower Defense", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
            int32((maxWinHeight-BOTTOM_BAR_HEIGHT_PX)*float64(MAP_WIDTH_FIELD)/float64(MAP_HEIGHT_FIELD)), maxWinHeight,
            sdl.WINDOW_RESIZABLE)
    CheckErr(err)

    renderer, err := sdl.CreateRenderer(window, 0, 0)
    CheckErr(err)
    renderer.SetDrawColor(100, 100, 100, 255)
    renderer.Clear()
    renderer.Present()

    fmt.Printf("Loading %d textures\n", len(TEXTURES))
    i := 0
    for fileName := range TEXTURES {
        path := TEXTURE_DIR_PATH+string(os.PathSeparator)+fileName
        fmt.Printf("[%d/%d] Loading \"%s\"\n", i+1, len(TEXTURES), path)
        surface, err := img.Load(path)
        CheckErr(err)

        if surface.W != surface.H {
            ShowErrAndPanic(fmt.Sprintf("Non-rectangular texture: %s: %dx%d", path, surface.W, surface.H))
        }

        texture, err := renderer.CreateTextureFromSurface(surface)
        CheckErr(err)
        tex := Texture{Texture: texture, Width: surface.W, Height: surface.H}
        surface.Free()
        TEXTURES[fileName] = &tex
        fmt.Printf("Loaded \"%s\", size: %dx%d\n", fileName, tex.Width, tex.Height)
        i++
    }
    fmt.Println("Textures:", TEXTURES)

    fmt.Println("Loading font: "+FONT_FILE_PATH)
    OpenFont(renderer, FONT_FILE_PATH)
    fmt.Println("Font loaded")

    //--------------------------- Variables ------------------------------------

    coins := 100
    hp := 100
    var towers []tower.ITower
    var enemies []enemy.IEnemy

    placedTowerType := tower.TOWER_TYPE_NONE
    var previewTower tower.ITower

    var mouseX, mouseY int32
    var mouseState uint32

    isTowerAt := func(col int32, row int32) bool {
        for _, tower := range towers {
            if tower.GetFieldCol() == col && tower.GetFieldRow() == row {
                return true
            }
        }
        return false
    }

    isRoadAt := func(col int32, row int32) bool {
        for _, coord := range ROAD_COORDS {
            if int32(coord.X) == col && int32(coord.Y) == row {
                return true
            }
        }
        return false
    }

    posToField := func(x int32, y int32) (int32, int32) {
        col := int32(float64(mouseX)/FIELD_SIZE_PX)
        row := int32(float64(mouseY)/FIELD_SIZE_PX)
        return col, row
    }

    switchPlacedTowerType := func(typ tower.TowerType) {
        if placedTowerType == typ {
            return
        }

        placedTowerType = typ

        switch typ {
        case tower.TOWER_TYPE_NONE:
            previewTower = nil

        case tower.TOWER_TYPE_CANNON:
            tow := tower.Cannon{FieldCol: 0, FieldRow: 0, IsPreview_: true}
            previewTower = &tow

        case tower.TOWER_TYPE_ROCKETTOWER:
            tow := tower.RocketTower{FieldCol: 0, FieldRow: 0, IsPreview_: true}
            previewTower = &tow

        default: panic(typ)
        }
    }

    // TODO: Test -- Remove later
    tank1 := enemy.Tank{FieldCol: 3, FieldRow: 1, Hp: 10}
    enemies = append(enemies, &tank1)
    tank2 := enemy.Tank{FieldCol: 14, FieldRow: 5, Hp: 10}
    enemies = append(enemies, &tank2)

    //--------------------------- Main loop ------------------------------------

    fmt.Println("Setup done")

    done := false
    var startTime uint32 = 1
    var frameTime uint32 = 1
    for {
        startTime = sdl.GetTicks()
        winW, winH := window.GetSize()
        mouseX, mouseY, mouseState = sdl.GetMouseState()
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
                col, row := posToField(mouseX, mouseY)
                if placedTowerType != tower.TOWER_TYPE_NONE && coins >= placedTowerType.GetPrice() &&
                        IsInsideWorld(mouseX, mouseY) &&
                        !isTowerAt(col, row) && !isRoadAt(col, row) {

                    coins -= placedTowerType.GetPrice()

                    var tower_ tower.ITower

                    switch (placedTowerType) {
                    case tower.TOWER_TYPE_CANNON:
                        tower_ = &tower.Cannon{
                            FieldCol: col,
                            FieldRow: row,
                            IsPreview_: false,
                            Hp: placedTowerType.GetInitialHP()}

                    case tower.TOWER_TYPE_ROCKETTOWER:
                        tower_ = &tower.RocketTower{
                            FieldCol: col,
                            FieldRow: row,
                            IsPreview_: false,
                            Hp: placedTowerType.GetInitialHP()}

                    default: panic(placedTowerType)
                    }

                    towers = append(towers, tower_)

                    fmt.Printf("Placed a tower at {%d, %d}\n", col, row)
                }
                //fmt.Printf("{%d, %d}\n", col, row);

            case sdl.MOUSEWHEEL:
                newType := placedTowerType
                if event.(*sdl.MouseWheelEvent).Y > 0 {
                    newType++
                } else {
                    newType--
                }

                if newType >= tower.TOWER_TYPE__COUNT {
                    newType = tower.TOWER_TYPE__COUNT-1
                } else if newType < 0 {
                    newType = 0
                }
                switchPlacedTowerType(newType)
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
        ASSERT_TRUE(hp >= 0) // TODO: Handle death
        drawBottomBar(renderer, winW, winH, coins, hp)

        // Update entities
        for _, enemy := range enemies { enemy.Update() }
        for _, tower := range towers { tower.Update(enemies) }

        // Render entities
        for _, enemy := range enemies { enemy.Render(renderer) }
        for _, tower := range towers { tower.Render(renderer) }

        for _, t := range towers {
            t.CheckCursorHover(renderer, mouseX, mouseY)
        }

        if previewTower != nil && IsInsideWorld(mouseX, mouseY) {
            col, row := posToField(mouseX, mouseY)
            previewTower.SetFieldCol(col)
            previewTower.SetFieldRow(row)
            previewTower.Render(renderer)
        }

        renderer.Present()
        sdl.Delay(16)

        frameTime = sdl.GetTicks() - startTime
    }

    //----------------------------- Cleanup ------------------------------------

    for _, texture := range TEXTURES {
        texture.Texture.Destroy()
    }

    for _, texture := range CHAR_TEXTURES {
        texture.Texture.Destroy()
    }

    renderer.Destroy()
    window.Destroy()
    sdl.Quit()

    fmt.Println("Window closed")
}

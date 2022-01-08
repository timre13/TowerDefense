package particle_source

import (
    "github.com/veandco/go-sdl2/sdl"
    . "TowerDefense/common"
    "math"
    "math/rand"
)

type Particle struct {
    XPos    int32
    YPos    int32
    MovRad  float64
    Life    int
    Color   sdl.Color
    Type    int
}

func (p* Particle) IsAlive() bool {
    return p.Life > 0
}

func (p* Particle) Update() {
    if (!p.IsAlive()) {
        return;
    }
    p.Life--

    var speed float64
    switch (p.Type) {
    case PARTSRC_TYPE_FIRE:
        speed = 8

    case PARTSRC_TYPE_SMOKE:
        speed = 1

    default:
        panic(p.Type)
    }

    // Move in the direction of the last movement
    p.XPos += int32(math.Round(math.Sin(p.MovRad)*speed))
    p.YPos += int32(math.Round(math.Cos(p.MovRad)*speed))
    // Randomly rotate the movement vector a bit
    if rand.Float64() > 0.6 {
        p.MovRad += (rand.Float64()-0.5)
    }
}

func (p* Particle) Render(rend *sdl.Renderer) {
    if (!p.IsAlive()) {
        return;
    }

    var size int32 
    switch (p.Type) {
    case PARTSRC_TYPE_FIRE:
        size = 8

    case PARTSRC_TYPE_SMOKE:
        size = 16

    default:
        panic(p.Type)
    }


    // Fade away before death
    alpha := p.Life*30
    if alpha > 255 {
        alpha = 255
    }
    TEXTURES[TEXTURE_FILENAME_PARTICLE].Texture.SetAlphaMod(uint8(alpha))

    TEXTURES[TEXTURE_FILENAME_PARTICLE].Texture.SetColorMod(p.Color.R, p.Color.G, p.Color.B)

    rect := sdl.Rect{X: p.XPos, Y: p.YPos, W: size, H: size}
    rend.Copy(TEXTURES[TEXTURE_FILENAME_PARTICLE].Texture, nil, &rect)
}

//-------------------------------------------------------------------------------

const (
    PARTSRC_TYPE_FIRE  = iota
    PARTSRC_TYPE_SMOKE = iota
)

type ParticleSource struct {
    XPos            int32
    YPos            int32
    Life            int
    Type            int
    spawnProb       float64
    partLife        int
    particles       *[]*Particle
}

func NewParticleSource (
    xpos            int32,
    ypos            int32,
    life            int,
    typ            int) ParticleSource {

    var spawnProb float64
    var partLife int
    switch (typ) {
    case PARTSRC_TYPE_FIRE:
        spawnProb = 0.5
        partLife = 30

    case PARTSRC_TYPE_SMOKE:
        spawnProb = 10
        partLife = 120

    default:
        panic(typ)
    }

    return ParticleSource{
        XPos: xpos,
        YPos: ypos,
        Life: life,
        Type: typ,
        spawnProb: spawnProb,
        partLife: partLife,
        particles: &[]*Particle{}}
}

func (ps* ParticleSource) IsAlive() bool {
    return ps.Life > 0
}

func (ps *ParticleSource) HasParticles() bool {
    return len(*ps.particles) > 0
}

func (ps* ParticleSource) Update() {
    ps.Life--

    for i := 0; i < len(*ps.particles); i++ {
        (*ps.particles)[i].Update()

        // Remove particle if expired
        if !(*ps.particles)[i].IsAlive() {
            *ps.particles = append((*ps.particles)[:i], (*ps.particles)[i+1:]...)
        }
    }

    if ps.Life > 0 && rand.Float64() <= ps.spawnProb {
        appendCount := 1
        if ps.spawnProb > 1 {
            appendCount = int(ps.spawnProb)
        }
        
        for i:=0; i < appendCount; i++ {
            var color sdl.Color
            switch (ps.Type) {
            case PARTSRC_TYPE_SMOKE:
                plusCol := uint8(rand.Int31n(100))
                color = sdl.Color{R: 10+plusCol, G: 10+plusCol, B: 10+plusCol, A: 255}

            default:
                panic(ps.Type)
            }
            *ps.particles = append(*ps.particles, &Particle{
                XPos: ps.XPos+int32(rand.Int31n(41)-20), YPos: ps.YPos+int32(rand.Int31n(41)-20), Life: ps.partLife, MovRad: math.Pi*2*rand.Float64(),
                Color: color, Type: PARTSRC_TYPE_SMOKE})
        }
    }
}

func (ps* ParticleSource) Render(rend *sdl.Renderer) {
    for _, p := range *ps.particles {
        p.Render(rend)
    }
}

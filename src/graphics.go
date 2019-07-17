package chip8

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	//ScreenWidth is the width of the screen.
	ScreenWidth = 64
	//ScreenHeight is the height of the screen.
	ScreenHeight = 32

	//ScreenScale specifies how the screen will be scaled.
	ScreenScale = 20
)

//Graphics handles the window for the Chip8.
type Graphics struct {
	m *Memory

	screen    [ScreenWidth][ScreenHeight]uint8
	screenmux sync.Mutex

	window  *sdl.Window
	surface *sdl.Surface

	w, h, scale int32

	fgColour color.Color
	bgColour color.Color
}

//NewGraphics returns a new graphics struct with initialised values.
func NewGraphics(mem *Memory) *Graphics {
	return &Graphics{
		m:        mem,
		w:        ScreenWidth,
		h:        ScreenHeight,
		scale:    ScreenScale,
		fgColour: color.RGBA{R: 255, G: 255, B: 255, A: 1},
		bgColour: color.RGBA{R: 0, G: 0, B: 0, A: 1},
	}
}

//Init initalises the sdl window.
func (g *Graphics) Init() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("could not initialise sdl: %v", err)
	}

	// Stop "non-name on left side of :=" error
	var err error
	g.window, err = sdl.CreateWindow("Chip8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		g.w*g.scale, g.h*g.scale, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create sdl window: %v", err)
	}

	/*
		for i := 0; i < int(g.w); i++ {
			for j := 0; j < int(g.h); j++ {
				g.screen[i][j] = 1
			}
		}
	*/

	g.surface, err = g.window.GetSurface()
	if err != nil {
		return fmt.Errorf("could not get sdl window surface: %v", err)
	}
	return nil
}

//Run starts the main loop for the graphics struct.
func (g *Graphics) Run() error {
	defer sdl.Quit()
	defer g.window.Destroy()

	running := true
	for running {
		err := g.PaintSurface()
		if err != nil {
			return fmt.Errorf("could not paint surface: %v", err)
		}
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}

	return nil
}

//PaintSurface takes the screen array and translates it into pixels on the window surface.
func (g *Graphics) PaintSurface() error {
	g.screenmux.Lock()
	defer g.screenmux.Unlock()
	for h := 0; h < int(g.h); h++ {
		for w := 0; w < int(g.w); w++ {
			if g.screen[w][h] == 1 {
				pixel := sdl.Rect{X: int32(w) * g.scale, Y: int32(h) * g.scale, W: g.scale, H: g.scale}
				g.surface.FillRect(&pixel, 0x0FFFFFFFF)
			}
		}
	}
	g.window.UpdateSurface()

	return nil
}

//Draw sprites onto the screen.
func (g *Graphics) Draw(w int32, h int32, n uint8, addr uint16) (bool, error) {
	a := addr & 0x0FFF

	var sprite []uint8
	carry := false

	//Read the sprite data from memory.
	for i := 0; i < int(n); i++ {
		spriteline, err := g.m.Read(a + uint16(i))

		if err != nil {
			return false, fmt.Errorf("could not read sprite data: %v", err)
		}

		sprite = append(sprite, uint8(spriteline))
	}

	g.screenmux.Lock()
	defer g.screenmux.Unlock()
	for screenh := 0; screenh < len(sprite); screenh++ {
		for screenw := 0; screenw < len(sprite); screenw++ {
			pixel := (sprite[screenh] >> uint8((7 - screenw))) & 0x01

			//Check if any bits will get flipped.
			if pixel != g.screen[(w+int32(screenw))%g.w][(h+int32(screenh))%g.h] {
				carry = true
			}

			//Xor screen pixel with sprite pixel, modulo there for wrap around.
			g.screen[(w+int32(screenw))%g.w][(h+int32(screenh))%g.h] ^= pixel
		}
	}

	return carry, nil
}

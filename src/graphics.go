package chip8

import (
	"fmt"
	"math/rand"

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
	m Memory

	screen [ScreenWidth][ScreenHeight]uint8

	window  *sdl.Window
	surface *sdl.Surface

	w, h, scale int32
}

//NewGraphics returns a new graphics struct with initialised values.
func NewGraphics(mem Memory) *Graphics {
	return &Graphics{
		m:     mem,
		w:     ScreenWidth,
		h:     ScreenHeight,
		scale: ScreenScale,
	}
}

//Init initalises the sdl window.
func (g *Graphics) Init() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("could not initialise sdl: %v", err)
	}

	// Stop "non-name on left side of :=" error
	var err error
	g.window, err = sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		g.w*g.scale, g.h*g.scale, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create sdl window: %v", err)
	}

	for i := 0; i < int(g.w); i++ {
		for j := 0; j < int(g.h); j++ {
			g.screen[i][j] = 1
		}
	}

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

	//g.surface.FillRect(nil, 0)

	/*
		rect := sdl.Rect{X: 0, Y: 0, W: 200, H: 200}
		g.surface.FillRect(&rect, 0xffff0000)
		g.window.UpdateSurface()
	*/

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
	for w := 0; w < int(g.w); w++ {
		for h := 0; h < int(g.h); h++ {
			if g.screen[w][h] == 1 {
				pixel := sdl.Rect{X: int32(w) * g.scale, Y: int32(h) * g.scale, W: g.scale, H: g.scale}
				g.surface.FillRect(&pixel, rand.Uint32())
			}
		}
	}
	g.window.UpdateSurface()

	return nil
}

package chip8

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type Graphics struct {
	m Memory

	screen []int

	window  sdl.Window
	surface sdl.Surface

	w, h int
}

func (g *Graphics) Init() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("could not initialise sdl: %v", err)
	}

	// TODO: Move this: defer sdl.Quit()

	g.window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create sdl window: %v", err)
	}
	defer window.Destroy()

	g.surface, err := window.GetSurface()
	if err != nil {
		return fmt.Errorf("could not get sdl window surface: %v", err)
	}
}

func (g *Graphics) Run() error {
	defer sdl.Quit()

	g.surface.FillRect(nil, 0)

	rect := sdl.Rect{0, 0, 200, 200}
	g.surface.FillRect(&rect, 0xffff0000)
	g.window.UpdateSurface()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}
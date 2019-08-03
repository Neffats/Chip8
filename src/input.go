package chip8

import (
	"fmt"

	"github.com/Neffats/bimap"
	"github.com/veandco/go-sdl2/sdl"
)

//Input handles the state info of the keyboard.
type Input struct {
	keys []uint8

	keymap *bimap.Uint8
}

//NewInput returns an empty uninitialised Input struct.
//Should I init() here?
func NewInput() *Input {
	return &Input{}
}

//Init sets up the keys array.
func (i *Input) Init() error {
	//Only have to call once.
	i.keys = sdl.GetKeyboardState()

	//Set up the keyboard translation mappings.
	i.keymap = bimap.NewUint8()

	i.keymap.Put(0x1, uint8(sdl.K_1))
	i.keymap.Put(0x2, uint8(sdl.K_2))
	i.keymap.Put(0x3, uint8(sdl.K_3))
	i.keymap.Put(0xC, uint8(sdl.K_4))
	i.keymap.Put(0x4, uint8(sdl.K_q))
	i.keymap.Put(0x5, uint8(sdl.K_w))
	i.keymap.Put(0x6, uint8(sdl.K_e))
	i.keymap.Put(0xD, uint8(sdl.K_r))
	i.keymap.Put(0x7, uint8(sdl.K_a))
	i.keymap.Put(0x8, uint8(sdl.K_s))
	i.keymap.Put(0x9, uint8(sdl.K_d))
	i.keymap.Put(0xE, uint8(sdl.K_f))
	i.keymap.Put(0xA, uint8(sdl.K_BACKSLASH))
	i.keymap.Put(0x0, uint8(sdl.K_z))
	i.keymap.Put(0xB, uint8(sdl.K_x))
	i.keymap.Put(0xF, uint8(sdl.K_c))

	return nil
}

//IsPressed returns true if the specified key is currently pressed.
func (i *Input) IsPressed(key uint8) (bool, error) {
	if key > 0xF {
		return false, fmt.Errorf("key out of bounds: %x", key)
	}
	k, exists := i.keymap.GetByKey(key)
	if !exists {
		return false, fmt.Errorf("key received does not exist in keymap: %d", key)
	}
	if i.keys[k] == 1 {
		return true, nil
	}
	return false, nil
}

//WaitForKey will loop and do nothing until specified key is pressed.
func (i *Input) WaitForKey() (uint8, error) {

	for event := sdl.WaitEvent(); event != nil; event = sdl.WaitEvent() {
		switch t := event.(type) {
		case *sdl.KeyboardEvent:
			if k := t.GetType(); k == sdl.KEYDOWN {
				key, exists := i.keymap.GetByValue(uint8(t.Keysym.Sym))
				if !exists {
					continue
				}
				fmt.Printf("SDL Key Pressed: %d\n", uint8(t.Keysym.Sym))
				fmt.Printf("Local mapping Key Pressed: %d\n", key)
				return key, nil
			}
		}
	}
	return 0, nil
}

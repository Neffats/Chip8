package chip8

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

//Input handles the state info of the keyboard.
type Input struct {
	keys []uint8

	keymap map[uint8]uint8
}

//Init sets up the keys array.
func (i *Input) Init() error {
	//Only have to call once.
	i.keys = sdl.GetKeyboardState()

	//Set up the keyboard translation mappings.
	i.keymap = map[uint8]uint8{
		0x1: sdl.K_1,
		0x2: sdl.K_2,
		0x3: sdl.K_3,
		0xC: sdl.K_4,
		0x4: sdl.K_q,
		0x5: sdl.K_w,
		0x6: sdl.K_e,
		0xD: sdl.K_r,
		0x7: sdl.K_a,
		0x8: sdl.K_s,
		0x9: sdl.K_d,
		0xE: sdl.K_f,
		0xA: sdl.K_BACKSLASH,
		0x0: sdl.K_z,
		0xB: sdl.K_x,
		0xF: sdl.K_c,
	}

	return nil
}

//IsPressed returns true if the specified key is currently pressed.
func (i *Input) IsPressed(key uint8) (bool, error) {
	if key > 0xF {
		return false, fmt.Errorf("key out of bounds: %x", key)
	}
	if i.keys[i.keymap[key]] == 1 {
		return true, nil
	}
	return false, nil
}

//WaitForKey will loop and do nothing until specified key is pressed.
func (i *Input) WaitForKey() (uint8, error) {
	for i.keys[i.keymap[key]] != 1 {
	}
	return 0, nil
}

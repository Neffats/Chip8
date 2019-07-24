package chip8

import "github.com/veandco/go-sdl2/sdl"

//Input handles the state info of the keyboard.
type Input struct {
	keys []uint8
}

//Init sets up the keys array.
func (i *Input) Init() error {
	//Only have to call once.
	i.keys = sdl.GetKeyboardState()
	return nil
}

//IsPressed returns true if the specified key is currently pressed.
func (i *Input) IsPressed() (bool, error) {
	return false, nil
}

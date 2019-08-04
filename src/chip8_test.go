package chip8

type treg struct {
	reg   uint8
	value uint8
}

func setup() *CPU {
	m := Memory{}
	g := NewGraphics(&m)
	i := NewInput()
	dt := NewTimer()
	return NewCPU(&m, g, i, dt)
}

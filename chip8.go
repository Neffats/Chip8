package chip8

type cpu interface {
	Get(reg string) (uint8, error)
	Set(reg string, data uint8) error
}

type memory interface {
	Write(data byte, addr uint16) error
	Read(addr uint16) (byte, error)
}

type Chip8 struct {
	Memory memory
	CPU    cpu
}

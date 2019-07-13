package chip8

type memory interface {
	Write(data byte, addr uint16) error
	Read(addr uint16) (byte, error)
}

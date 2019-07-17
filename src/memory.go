package chip8

import "fmt"

//Memory module
type Memory struct {
	memory [4096]byte
}

//Write writes a byte to a specified 12 bit address
func (m *Memory) Write(data byte, addr uint16) error {
	msb := addr & 0xF000
	if msb != 0 {
		return fmt.Errorf("address out of bounds: %x", addr)
	}
	m.memory[addr] = data

	return nil
}

//Read a byte from a specified 12 bit address
func (m *Memory) Read(addr uint16) (byte, error) {
	msb := addr & 0xF000
	if msb != 0 {
		return 0, fmt.Errorf("address out of bounds: %x", addr)
	}

	return m.memory[addr], nil
}

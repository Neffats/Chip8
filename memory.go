package chip8

type Memory struct {
	memory [4096]byte
}

func (m Memory) Write(data byte, addr uint16) error {
	m.memory[addr] = data
	return nil
}

func (m Memory) Read(addr uint16) (byte, error) {
	return m.memory[addr], nil
}

package main

import (
	chip8 "Github/Chip8/src"
	"fmt"
)

func main() {
	m := chip8.Memory{}
	g := chip8.NewGraphics(&m)

	err := setupMemory(&m)
	if err != nil {
		panic(err)
	}

	err = g.Init()
	if err != nil {
		panic(err)
	}
	_, err = g.Draw(10, 10, 5, 0x0000)
	_, err = g.Draw(15, 10, 5, 0x0005)
	_, err = g.Draw(20, 10, 5, 0x0005)
	if err != nil {
		panic(err)
	}

	err = g.Run()
	if err != nil {
		panic(err)
	}

}

func setupMemory(m *chip8.Memory) error {
	err := m.Write(0xE0, 0x0000)
	if err != nil {
		return fmt.Errorf("could not write memory: %v", err)
	}
	err = m.Write(0x90, 0x0001)
	if err != nil {
		return fmt.Errorf("could not write memory: %v", err)
	}
	err = m.Write(0xE0, 0x0002)
	if err != nil {
		return fmt.Errorf("could not write memory: %v", err)
	}
	err = m.Write(0x90, 0x0003)
	if err != nil {
		return fmt.Errorf("could not write memory: %v", err)
	}
	err = m.Write(0xE0, 0x0004)
	if err != nil {
		return fmt.Errorf("could not write memory: %v", err)
	}

	err = m.Write(0xF0, 0x0005)
	if err != nil {
		return fmt.Errorf("could not write memory: %v", err)
	}
	err = m.Write(0x80, 0x0006)
	if err != nil {
		return fmt.Errorf("could not write memory: %v", err)
	}
	err = m.Write(0xF0, 0x0007)
	if err != nil {
		return fmt.Errorf("could not write memory: %v", err)
	}
	err = m.Write(0x80, 0x0008)
	if err != nil {
		return fmt.Errorf("could not write memory: %v", err)
	}
	err = m.Write(0xF0, 0x0009)
	if err != nil {
		return fmt.Errorf("could not write memory: %v", err)
	}
	return nil
}

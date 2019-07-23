package main

import (
	"flag"
	"fmt"
	"os"

	chip8 "github.com/Neffats/Chip8/src"
)

func GetFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)

	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not get file info: %v", err)

	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("could not read bytes into buffer: %v", err)
	}

	return buffer, nil
}

func main() {
	program := flag.String("program", "", "Chip8 program file.")
	flag.Parse()

	ProgramData, err := GetFile(*program)

	m := chip8.Memory{}
	g := chip8.NewGraphics(&m)
	c := chip8.NewCPU(&m, g)

	err = setupMemory(&m)
	if err != nil {
		panic(err)
	}

	err = g.Init()
	if err != nil {
		panic(err)
	}
	defer g.Destroy()

	err = c.LoadProgram(ProgramData)
	if err != nil {
		panic(err)
	}

	err = c.Run()
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

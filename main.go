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
	program := flag.String("p", "", "Chip8 program file.")
	flag.Parse()

	ProgramData, err := GetFile(*program)

	m := chip8.Memory{}
	g := chip8.NewGraphics(&m)
	in := chip8.NewInput()
	dt := chip8.NewTimer()
	in.Init()
	c := chip8.NewCPU(&m, g, in, dt)

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
		c.Panic(err)
	}
}

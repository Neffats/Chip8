package main

import chip8 "Github/Chip8/src"

func main() {
	m := chip8.Memory{}
	g := chip8.NewGraphics(m)

	err := g.Init()
	if err != nil {
		panic(err)
	}

	err = g.Run()
	if err != nil {
		panic(err)
	}
}

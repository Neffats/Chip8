package main

import chip8 "Github/Chip8/src"

func main() {
	g := chip8.Graphics{}

	err := g.Init()
	if err != nil {
		panic(err)
	}

	err = g.Run()
	if err != nil {
		panic(err)
	}
}

package main

import (
	"github.com/CCAtAlvis/xorfall/src"
)

func main() {
	game := src.NewGame()
	defer game.Close()

	game.Start()
}

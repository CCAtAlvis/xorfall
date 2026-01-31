package main

import (
	"github.com/CCAtAlvis/xorfall/src/core"
)

func main() {
	game := core.NewGame()
	defer game.Close()

	game.Start()
}

package main

import (
	"log"

	"github.com/peterSZW/ebdemo/ebgame"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := &ebgame.Game{}
	// w := 1024
	// h := 768
	w := 400
	h := 880

	game.SetWindowSize(w, h)
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Hello, World!")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

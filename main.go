package main

import (
	"log"

	"github.com/peterSZW/ebdemo/ebgame"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(1024, 768)
	//ebiten.SetWindowSize(512, 384)
	ebiten.SetWindowTitle("Hello, World!")
	game := &ebgame.Game{}
	//game.SetWindowSize(512, 384)
	game.SetWindowSize(1024, 768)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

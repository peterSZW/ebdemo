package main

import (
	"log"

	"github.com/peterszw/ebdemo/ebgame"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&ebgame.Game{}); err != nil {
		log.Fatal(err)
	}
}

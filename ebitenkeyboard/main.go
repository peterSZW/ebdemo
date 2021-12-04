package main

import (
	//_ "image/png"
	"log"
	"mykb" ///localt GOROOT /usr/local/go/src/mykb

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(mykb.ScreenWidth, mykb.ScreenHeight)
	ebiten.SetWindowTitle("Keyboard (Ebiten Demo)")
	if err := ebiten.RunGame(&mykb.Game{}); err != nil {
		log.Fatal(err)
	}
}

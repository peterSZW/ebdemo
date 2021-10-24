package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/peterSZW/ebdemo/ebgame"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

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

	f2, err := os.Create("mem.pprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(f2)
	f2.Close()
}

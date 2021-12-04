package mobile

import (
	"mykb"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/mobile"
)

type size struct {
	width  int
	height int
}

var (
	window size
	game   *mykb.Game
)

func SetWindowSize(width, height int) {
	window.width = width
	window.height = height
	ebiten.SetWindowSize(width, height)
	game.SetWindowSize(width, height)
}

func init() {
	game = &mykb.Game{}
	mobile.SetGame(game)
}

func Dummy() {}

package mobile

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/peterSZW/ebdemo/ebgame"
)

type size struct {
	width  int
	height int
}

var (
	window size
	game   *ebgame.Game
)

func SetWindowSize(width, height int) {
	window.width = width
	window.height = height
	ebiten.SetWindowSize(width, height)
	game.SetWindowSize(width, height)
}

func init() {
	game = &ebgame.Game{}
	mobile.SetGame(game)
}

func Dummy() {}

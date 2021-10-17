package mobile

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"
	"eb-demo"
)

type size struct {
	width  int
	height int
}

var (
	window size
	game   *eb-demo.Game
)

func init() {
	mobile.SetGame(&Game{})
}
 
func Dummy() {}

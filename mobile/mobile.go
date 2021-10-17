package mobile

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/peterSZW/ebdemo/ebgame"
)

func init() {
	mobile.SetGame(&ebgame.Game{})
}

func Dummy() {}

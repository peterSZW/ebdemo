package mobile

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

func kkk() {
	var err error
	//读图片
	img, _, err = ebitenutil.NewImageFromFile("10.png")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(img)

}
func Dummy() {}

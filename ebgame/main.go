package ebgame

import (
	"fmt"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var img *ebiten.Image

var pointerImage = ebiten.NewImage(8, 8)

//初始化
func init() {
	// var err error
	// //读图片
	// img, _, err = ebitenutil.NewImageFromFile("10.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	pointerImage.Fill(color.RGBA{0xff, 0, 0, 0xff})

	xx = 5
	yy = 5

}

type Game struct{}

var x float64
var y float64
var xx float64
var yy float64

var r float64

//循环计算
func (g *Game) Update() error {
	x = x + xx
	y = y + yy

	if x > 1024-100 {
		xx = -5
	}

	if x < 0 {
		xx = 5
	}

	if y > 768-100 {
		yy = -5
	}
	if y < 0 {
		yy = 5
	}
	// r = r + 0.1
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//打印 hello world 加帧数

	mx, my := ebiten.CursorPosition()
	s := fmt.Sprintf("Hello, World! FPS : %f %d %d", ebiten.CurrentFPS(), mx, my)

	ebitenutil.DebugPrint(screen, s)

	//画图
	// op.GeoM.Rotate(r)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	op.GeoM.Scale(0.5, 0.5)

	//screen.DrawImage(img, op)
	screen.DrawImage(pointerImage, op)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 512, 384
}

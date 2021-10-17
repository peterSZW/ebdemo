package ebgame

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var img *ebiten.Image
var screenSize image.Point
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
var touchStr string

const (
	widthAsDots = 480.
)

//循环计算
func (g *Game) Update() error {
	x = x + xx
	y = y + yy

	if x > float64(screenSize.X)-2.0 {
		xx = -5
	}

	if x < 0 {
		xx = 5
	}

	if y > float64(screenSize.Y)-2.0 {
		yy = -5
	}
	if y < 0 {
		yy = 5
	}
	// r = r + 0.1

	touchStr = ""

	touches := ebiten.TouchIDs()

	if len(touches) > 0 {
		for _, id := range touches {
			x, y := ebiten.TouchPosition(id)
			touchStr = touchStr + "\n" + fmt.Sprintf("(%d,%d)", x, y)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//打印 hello world 加帧数

	mx, my := ebiten.CursorPosition()
	s := fmt.Sprintf("\n\n\nHello, World! FPS : %f %d %d %s", ebiten.CurrentFPS(), mx, my, touchStr)

	ebitenutil.DebugPrint(screen, s)

	//画图
	// op.GeoM.Rotate(r)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	//op.GeoM.Scale(0.5, 0.5)

	screen.DrawImage(pointerImage, op)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x+10, y)
	screen.DrawImage(pointerImage, op)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {

	return screenSize.X, screenSize.Y
}
func (g *Game) SetWindowSize(width, height int) {
	// screenSize.X = int(widthAsDots)
	// screenSize.Y = int(widthAsDots / float64(width) * float64(height))
	screenSize.X = width
	screenSize.Y = height

}

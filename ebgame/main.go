package ebgame

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/peterSZW/ebdemo/ebgame/gfx"
)

var img *ebiten.Image
var screenSize image.Point
var pointerImage = ebiten.NewImage(8, 8)

var home string
var curpath string
var errstr string

var (
	explosion1, explosion2, explosion3 *Sprite
)

//初始化
func init() {
	var err error
	//读图片
	img, _, err = ebitenutil.NewImageFromFile("10.png")
	if err != nil {
		log.Fatal(err)
	}

	pointerImage.Fill(color.RGBA{0xff, 0, 0, 0xff})

	xx = 5
	yy = 5

	home = os.Getenv("HOME")
	curpath = getCurrentDirectory()
	errstr = ""

	f, err := os.Create(home + "/Library/Caches/output3.txt") //创建文件
	defer f.Close()
	if err == nil {
		_, err = f.WriteString("writesn") //写入文件(字节数组)

		f.Sync()
	}
	if err != nil {
		errstr = err.Error()

	}

	explosion3 = NewSprite()
	explosion3.AddAnimationByte("default", &gfx.EXPLOSION3, 500, 9, ebiten.FilterNearest)
	//explosion3.AddAnimation("default", "gfx/explosion3.png", explosionDuration, 9, ebiten.FilterNearest)
	explosion3.Position(240-10-48, 400/3*2)
	explosion3.Start()

}
func file_exist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
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

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err.Error()

	}
	return strings.Replace(dir, "\\", "/", -1)
}

func (g *Game) Draw(screen *ebiten.Image) {
	//打印 hello world 加帧数

	mx, my := ebiten.CursorPosition()

	s := fmt.Sprintf("\n\n\nHello, World! FPS : %f %d %d %s\n%s\n %v\n%d",
		ebiten.CurrentFPS(), mx, my, touchStr, curpath+"\n"+home,
		file_exist(home+"/Library/Caches/output3.txt"), errstr)

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

	explosion3.Draw(screen)

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

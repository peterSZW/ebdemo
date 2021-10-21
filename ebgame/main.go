package ebgame

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

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

var joytouch JoyTouch

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

func limiXY(x, y float64) (float64, float64) {
	if x > float64(screenSize.X)-2.0 {
		x = float64(screenSize.X) - 2.0
	}

	if x < 0 {
		x = 0
	}

	if y > float64(screenSize.Y)-2.0 {
		y = float64(screenSize.Y) - 2.0
	}
	if y < 0 {
		y = 0
	}
	return x, y
}

//循环计算
func (g *Game) Update() error {

	joytouch.SetWH(screenSize.X, screenSize.Y)

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

	isstillpress := false

	xx := 0.
	yy := 0.
	if len(touches) > 0 {

		if joytouch.tid != 0 {
			//alread have last press, so we need to find is this touch still on screen
			id := touches[0]
			for _, id = range touches {
				if int(id) == int(joytouch.GetTid()) {
					isstillpress = true
					break

				}
			}

			if isstillpress {
				x, y := ebiten.TouchPosition(id)

				xx, yy = joytouch.Move(x, y)
			} else {
				joytouch.tid = 0
				joytouch.x = 0
				joytouch.y = 0
				//find new press

				for _, id := range touches {

					x, y := ebiten.TouchPosition(id)

					if joytouch.Press(x, y, int(id)) {
						isstillpress = true
						break
					}

				}

			}

		} else {
			//find new press

			for _, id := range touches {

				x, y := ebiten.TouchPosition(id)

				if joytouch.Press(x, y, int(id)) {
					isstillpress = true
					touchStr = touchStr + "\n" + "Press"
					break
				}

			}

		}

		for _, id := range touches {

			x, y := ebiten.TouchPosition(id)

			touchStr = touchStr + "\n" + fmt.Sprintf("(%d,%d)", x, y)
		}

		explosion3.X = explosion3.X + xx
		explosion3.Y = explosion3.Y + yy

		explosion3.X, explosion3.Y = limiXY(explosion3.X, explosion3.Y)

	} else {
		//all reased
		joytouch.tid = 0
		joytouch.x = 0
		joytouch.y = 0
	}
	touchStr = touchStr + "\n" + fmt.Sprintf("[%f,%f]", xx, yy)
	touchStr = touchStr + "\n" + fmt.Sprintf("[%d,%d]", joytouch.x, joytouch.y)
	touchStr = touchStr + "\n" + fmt.Sprintf("[%d]", joytouch.tid)
	if isstillpress {
		touchStr = touchStr + "\n" + "STILL PRESS"
	}

	time.Sleep(10 * time.Millisecond)
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

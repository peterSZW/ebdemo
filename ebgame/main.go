package ebgame

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/peterSZW/ebdemo/ebgame/resources/images"
	"github.com/peterSZW/go-sprite/example/gfx"
)

var screenSize image.Point

var home string
var curpath string
var errstr string

var (
	robot *Sprite
)

var joytouch JoyTouch
var joybutton1 JoyButton
var joybutton2 JoyButton

var spList map[string]*Sprite
var spCount int

//初始化
func init() {
	spList = make(map[string]*Sprite)

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

	robot = NewSprite()
	robot.AddAnimationByte("default", &images.E_ROBO2, 2000, 8, ebiten.FilterNearest)
	robot.Position(300-10-48, 400/3*2)
	robot.CenterCoordonnates = true

	robot.Start()

}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err.Error()

	}
	return strings.Replace(dir, "\\", "/", -1)
}
func file_exist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

type Game struct{}

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

func GetKeyBoardSpace() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}
func GetKeyBoard() (xx, yy float64) {
	xx = 0
	yy = 0
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		xx = 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		xx = -5
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		yy = 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		yy = -5
	}
	return
}
func GetDirectByXY(xx, yy float64) int {
	// 7 8 1
	// 6 0 2
	// 5 4 3
	if xx == 0 && yy == 0 {
		return 0
	}

	if xx > 0 && yy < 0 {
		return 1
	}
	if xx > 0 && yy == 0 {
		return 2
	}
	if xx > 0 && yy > 0 {
		return 3
	}
	if xx == 0 && yy > 0 {
		return 4
	}
	if xx < 0 && yy > 0 {
		return 5
	}

	if xx < 0 && yy == 0 {
		return 6
	}
	if xx < 0 && yy < 0 {
		return 7
	}
	if xx == 0 && yy < 0 {
		return 8
	}
	return 0
}

var isFirstUpdate = true
var lastBulletTime time.Time

//循环计算
func (g *Game) Update() error {
	//第一次设置
	if isFirstUpdate {
		joytouch.SetWH(screenSize.X, screenSize.Y)
		joybutton1.SetWH(screenSize.X, screenSize.Y)
		joybutton1.rect.x = joybutton1.rect.x - 35
		joybutton2.SetWH(screenSize.X, screenSize.Y)
		joybutton2.rect.x = joybutton2.rect.x + 35
		isFirstUpdate = false
	}

	//移动飞机
	xx, yy := joytouch.GetJoyTouchXY()
	if xx == 0 && yy == 0 {
		xx, yy = GetKeyBoard()
	}
	robot.Pause()
	if GetDirectByXY(xx, yy) > 0 {
		robot.Step(GetDirectByXY(xx, yy))

	}

	robot.X = robot.X + xx
	robot.Y = robot.Y + yy
	robot.X, robot.Y = limiXY(robot.X, robot.Y)

	//生成子弹
	if joybutton1.GetJoyButton() || GetKeyBoardSpace() {
		if time.Now().Sub(lastBulletTime) > time.Duration(100*time.Millisecond) {
			lastBulletTime = time.Now()

			newsprite := NewSprite()
			//newsprite.AddAnimationByte("default", &images.MISSILE, 2000, 8, ebiten.FilterNearest)
			newsprite.AddAnimationByte("default", &gfx.EXPLOSION2, 500, 7, ebiten.FilterNearest)
			newsprite.Position(robot.X, robot.Y)
			newsprite.AddEffect(&EffectOptions{Effect: Zoom, Zoom: 3, Duration: 6000, Repeat: false, GoBack: true})
			newsprite.CenterCoordonnates = true
			newsprite.Start()

			spCount++

			spList[strconv.Itoa(spCount)] = newsprite
		}
	}

	if joybutton2.GetJoyButton() {
		if time.Now().Sub(lastBulletTime) > time.Duration(100*time.Millisecond) {
			lastBulletTime = time.Now()

			newsprite := NewSprite()
			newsprite.AddAnimationByte("default", &gfx.EXPLOSION3, 500, 9, ebiten.FilterNearest)
			newsprite.Position(robot.X, robot.Y)
			newsprite.AddEffect(&EffectOptions{Effect: Zoom, Zoom: 3, Duration: 6000, Repeat: false, GoBack: true})
			newsprite.CenterCoordonnates = true
			newsprite.Start()

			spCount++

			spList[strconv.Itoa(spCount)] = newsprite

		}
	}

	//删除越界的对象
	for k, j := range spList {
		j.Y = j.Y - 5
		if j.Y < -20 {
			j.Hide()
			delete(spList, k)
		}
	}

	//生成字符串
	touchStr = touchStr + "\n" + fmt.Sprintf("[%f,%f]", xx, yy)
	touchStr = touchStr + "\n" + fmt.Sprintf("[%d,%d]", joytouch.x, joytouch.y)
	touchStr = touchStr + "\n" + fmt.Sprintf("[%d]", joytouch.tid)

	//延迟0.01毫秒
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//打印 hello world 加帧数

	mx, my := ebiten.CursorPosition()

	s := fmt.Sprintf("\n\n\n%s\n%s\nFPS : %f %d %d\n%v\n%s\n%s\n%s",
		curpath, home, ebiten.CurrentFPS(), mx, my,
		file_exist(home+"/Library/Caches/output3.txt"), errstr, runtime.GOOS, touchStr)
	ebitenutil.DebugPrint(screen, s)

	joytouch.DrawBorders(screen, color.Gray16{0x2222})
	joybutton1.DrawBorders(screen, color.Gray16{0x1111})
	joybutton2.DrawBorders(screen, color.Gray16{0x1111})

	robot.Draw(screen)

	for _, j := range spList {
		j.Draw(screen)
	}
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

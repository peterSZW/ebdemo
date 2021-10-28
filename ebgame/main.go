package ebgame

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/peterSZW/ebdemo/ebgame/gfx"
	"github.com/peterSZW/ebdemo/ebgame/resources/images"
	"github.com/peterSZW/ebdemo/ebgame/sound"
)

var screenSize image.Point

var home string
var curpath string
var errstr string

var (
	robot    *Sprite
	touchpad *Sprite
)

var joytouch JoyTouch
var joybutton1 JoyButton
var joybutton2 JoyButton
var joybutton3 JoyButton

var isShowText bool

var spList map[string]*Sprite
var enemyList map[string]*Sprite
var effList map[string]*Sprite
var spCount int

//初始化
func init() {
	spList = make(map[string]*Sprite)
	enemyList = make(map[string]*Sprite)
	effList = make(map[string]*Sprite)

	sound.Load()
	sound.PlayBgm(sound.BgmKindBattle)
	home = os.Getenv("HOME")
	curpath = getCurrentDirectory()
	errstr = ""

	f, err := os.Create(home + "/Library/Caches/output3.txt") //创建文件

	if err == nil {
		defer f.Close()
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

	touchpad = NewSprite()
	touchpad.AddAnimationByte("default", &gfx.TOUCHPAD, 2000, 1, ebiten.FilterNearest)

	touchpad.CenterCoordonnates = true

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

// const (
// 	widthAsDots = 480.
// )

func OutofScreen(x, y float64, size float64) bool {

	return x > float64(screenSize.X)+size || x < -size || y > float64(screenSize.Y)+size || y < -size
}
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

func GetDegreeByXY(xx, yy float64) float64 {

	// 45  0  315
	// 90  0  270
	// 135 180 225

	if yy == 0 && xx == 0 {
		return 0
	}
	deg := 180 - math.Atan2(yy, xx)/math.Pi*180
	if deg < 270 {
		deg = deg + 90
	} else {
		deg = deg - 270
	}
	return deg
}
func GetDirectByXY(xx, yy float64) int {
	// 7 8 1
	// 6 0 2
	// 5 4 3
	deg := GetDegreeByXY(xx, yy)

	if 22 > deg || deg > 360-22 {
		return 8
	}
	if deg >= 0+22 && deg < 45+22 {
		return 7
	}
	if deg >= 45+22 && deg < 90+22 {
		return 6
	}
	if deg >= 90+22 && deg < 135+22 {
		return 5
	}

	if deg >= 135+22 && deg < 180+22 {
		return 4
	}

	if deg >= 180+22 && deg < 225+22 {
		return 3
	}
	if deg >= 225+22 && deg < 270+22 {
		return 2
	}
	if deg >= 270+22 && deg < 315+22 {
		return 1
	}
	if deg >= 315+22 && deg < 360+22 {
		return 8
	}

	if xx == 0 && yy == 0 {
		return 0
	}
	return 0

	// if xx > 0 && yy < 0 {
	// 	return 1
	// }
	// if xx > 0 && yy == 0 {
	// 	return 2
	// }
	// if xx > 0 && yy > 0 {
	// 	return 3
	// }
	// if xx == 0 && yy > 0 {
	// 	return 4
	// }
	// if xx < 0 && yy > 0 {
	// 	return 5
	// }

	// if xx < 0 && yy == 0 {
	// 	return 6
	// }
	// if xx < 0 && yy < 0 {
	// 	return 7
	// }
	// if xx == 0 && yy < 0 {
	// 	return 8
	// }
	// return 0
}

var robot2 *Sprite
var lastGenEnemyTime time.Time

func GenEnemy() {

	if time.Since(lastGenEnemyTime) > time.Duration(70*time.Millisecond) {
		lastGenEnemyTime = time.Now()

		newsprite := NewSprite()

		newsprite.AddAnimationByteCol("default", &images.E_ROBO1, 100, 1, 8, ebiten.FilterNearest)
		newsprite.Name = "E_ROBO2"
		newsprite.Position(float64(rand.Intn(screenSize.X)), 0)
		newsprite.CenterCoordonnates = true
		newsprite.Pause()
		//newsprite.Step(18)
		newsprite.Speed = float64(7 + rand.Intn(6))
		//newsprite.Angle = robot2.Direction //+90 GetDegreeByXY(xx, yy)
		newsprite.Direction = float64(270 - 20 + rand.Intn(20))

		//newsprite.Start()

		spCount++

		enemyList[strconv.Itoa(spCount)] = newsprite

	}
}

var isFirstUpdate = true
var lastBulletTime time.Time
var lastLaserTime time.Time
var degree float64

//循环计算
func (g *Game) Update() error {
	//第一次设置
	if isFirstUpdate {
		joytouch.SetWH(screenSize.X, screenSize.Y)
		touchpad.X = float64(joytouch.rect.x + joytouch.rect.w/2)
		touchpad.Y = float64(joytouch.rect.y + joytouch.rect.w/2)

		joybutton1.SetWH(screenSize.X, screenSize.Y)
		joybutton1.rect.x = joybutton1.rect.x - 35
		joybutton2.SetWH(screenSize.X, screenSize.Y)
		joybutton2.rect.x = joybutton2.rect.x + 35
		joybutton3.SetWH(screenSize.X, screenSize.Y)
		joybutton3.rect.y = 35

		isFirstUpdate = false
		isShowText = false
	}

	GenEnemy()

	//移动飞机
	xx, yy, _ := joytouch.GetJoyTouchXY()

	if xx == 0 && yy == 0 {
		xx, yy = GetKeyBoard()
	}
	robot.Pause()
	if GetDirectByXY(xx, yy) > 0 {
		robot.Step(GetDirectByXY(xx, yy))
		degree = GetDegreeByXY(xx, yy)
	}

	// if ispress {
	// 	robot.Step(GetDirectByXY(xx, yy))

	// }

	robot.X = robot.X + xx
	robot.Y = robot.Y + yy
	robot.X, robot.Y = limiXY(robot.X, robot.Y)

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		robot.Angle = robot.Angle + 10
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		robot.Angle = robot.Angle - 10
	}

	//生成子弹
	if joybutton1.GetJoyButton() || ebiten.IsKeyPressed(ebiten.KeySpace) {

		if time.Since(lastBulletTime) > time.Duration(100*time.Millisecond) {
			lastBulletTime = time.Now()

			newsprite := NewSprite()
			//newsprite.AddAnimationByteCol("default", &images.RASER1, 2000, 4, 6, ebiten.FilterNearest)
			newsprite.AddAnimationByte("default", &gfx.EXPLOSION2, 500, 7, ebiten.FilterNearest)
			newsprite.Name = "RASER1"
			newsprite.Position(robot.X, robot.Y)
			newsprite.AddEffect(&EffectOptions{Effect: Zoom, Zoom: 3, Duration: 6000, Repeat: false, GoBack: true})
			newsprite.CenterCoordonnates = true

			newsprite.Direction = degree + 90 //GetDegreeByXY(xx, yy) + 90
			//float64(2-robot.GetStep()) * 45

			newsprite.Speed = 5
			newsprite.Start()

			spCount++

			spList[strconv.Itoa(spCount)] = newsprite
		}
	}

	if joybutton3.GetJoyButton() || ebiten.IsKeyPressed(ebiten.KeyT) {
		isShowText = !isShowText
	}

	if joybutton2.GetJoyButton() || ebiten.IsKeyPressed(ebiten.Key1) {
		if time.Since(lastLaserTime) > time.Duration(50*time.Millisecond) {
			lastLaserTime = time.Now()

			newsprite := NewSprite()
			//newsprite.AddAnimationByte("default", &gfx.EXPLOSION3, 500, 9, ebiten.FilterNearest)
			newsprite.AddAnimationByteCol("default", &images.RASER1, 200, 4, 6, ebiten.FilterNearest)
			newsprite.Name = "RASER1"
			newsprite.Position(robot.X, robot.Y)
			//newsprite.AddEffect(&EffectOptions{Effect: Zoom, Zoom: 3, Duration: 6000, Repeat: false, GoBack: true})
			newsprite.CenterCoordonnates = true
			newsprite.Pause()
			newsprite.Step(18)
			newsprite.Speed = 8
			newsprite.Angle = degree          //+90 GetDegreeByXY(xx, yy)
			newsprite.Direction = degree + 90 // GetDegreeByXY(xx, yy) + 90

			//newsprite.Start()

			spCount++

			spList[strconv.Itoa(spCount)] = newsprite

		}
	}
	touchpad.Angle = degree

	//删除越界的对象
	for k, j := range spList {

		if OutofScreen(j.X, j.Y, 20) {
			j.Hide()
			delete(spList, k)
		}
	}

	//删除越界的对象
	for k, j := range enemyList {

		if OutofScreen(j.X, j.Y, 20) {
			j.Hide()
			delete(enemyList, k)
		}
	}

	checkCollision()

	//生成字符串
	touchStr = touchStr + "\n" + fmt.Sprintf("[%f,%f]", xx, yy)
	touchStr = touchStr + "\n" + fmt.Sprintf("[%d,%d]", joytouch.x, joytouch.y)
	touchStr = touchStr + "\n" + fmt.Sprintf("[%d]", joytouch.tid)
	touchStr = touchStr + "\n" + fmt.Sprintf("[%f]", robot.Angle)
	touchStr = touchStr + "\n" + fmt.Sprintf("[%f]", GetDegreeByXY(xx, yy))
	touchStr = touchStr + "\n" + fmt.Sprintf("[%v]", len(spList))
	touchStr = touchStr + "\n" + fmt.Sprintf("[%v]", len(enemyList))

	if joytouch.x != 0 && joytouch.y != 0 {
		touchpad.X = float64(joytouch.x)
		touchpad.Y = float64(joytouch.y)
	}
	//延迟0.01毫秒
	//time.Sleep(10 * time.Millisecond)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//打印 hello world 加帧数

	if isShowText {

		mx, my := ebiten.CursorPosition()

		s := fmt.Sprintf("\n\n\n%s\n%s\nFPS : %f %d %d\n%v\n%s\n%s\n%s",
			curpath, home, ebiten.CurrentFPS(), mx, my,
			file_exist(home+"/Library/Caches/output3.txt"), errstr, runtime.GOOS, touchStr)
		ebitenutil.DebugPrint(screen, s)
	}

	joytouch.DrawBorders(screen, color.Gray16{0x1111})
	//touchpad.Draw(screen)

	joybutton1.DrawBorders(screen, color.Gray16{0x1111})
	joybutton2.DrawBorders(screen, color.Gray16{0x1111})
	joybutton3.DrawBorders(screen, color.Gray16{0x1111})
	robot.Draw(screen)

	for _, j := range spList {
		j.Draw(screen)
	}

	for _, j := range enemyList {
		j.Draw(screen)
	}
	for k, j := range effList {
		j.Draw(screen)
		if j.GetStep() >= 7 {
			delete(effList, k)
		}
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

//
func checkCollision() {

	for kshot, shot := range spList {

		for k, enemy := range enemyList {

			if IsCollideWith(enemy, shot) == false {

				continue
			}

			{
				newsprite := NewSprite()
				//newsprite.AddAnimationByteCol("default", &images.RASER1, 2000, 4, 6, ebiten.FilterNearest)
				newsprite.AddAnimationByte("default", &images.EXPLODE_MED, 500, 8, ebiten.FilterNearest)

				newsprite.Position(enemy.X, enemy.Y)
				newsprite.CenterCoordonnates = true
				newsprite.Speed = enemy.Speed
				newsprite.Direction = enemy.Direction
				newsprite.Start()

				spCount++

				effList[strconv.Itoa(spCount)] = newsprite
			}

			delete(enemyList, k)
			delete(spList, kshot)
			sound.PlaySe(sound.SeKindHit2)

		}
	}

}

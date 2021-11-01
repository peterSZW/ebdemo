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
	"github.com/peterSZW/ebdemo/ebgame/paint"
	"github.com/peterSZW/ebdemo/ebgame/resources/images"
	"github.com/peterSZW/ebdemo/ebgame/sound"
)

var screenSize image.Point

var home string
var curpath string
var errstr string

var (
	robot     *Sprite
	robotpath *Sprite
	touchpad  *Sprite
)

var joytouch JoyTouch
var btnFire JoyButton
var btnBullet JoyButton
var btnDebug JoyButton
var btnShowBox JoyButton
var btnStart JoyButton

var isShowText bool

var spList map[string]*Sprite
var enemyList map[string]*Sprite
var effList map[string]*Sprite
var spCount int

type Gloable struct {
	Life             int
	Level            int
	Score            int
	ShowCollisionBox bool
}

var path Path

var gv Gloable

//初始化
func init() {
	spList = make(map[string]*Sprite)
	enemyList = make(map[string]*Sprite)
	effList = make(map[string]*Sprite)

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
	robot.Name = "E_ROBO2"
	robot.Position(float64(screenSize.X/2), float64(screenSize.Y/2)+100)
	robot.CenterCoordonnates = true
	robot.Start()

	robotpath = NewSprite()
	robotpath.AddAnimationByte("default", &images.E_SHOOTER1, 2000, 8, ebiten.FilterNearest)
	robotpath.Name = "E_SHOOTER1"
	robotpath.Position(float64(screenSize.X/2), float64(screenSize.Y/2)+100)
	robotpath.CenterCoordonnates = true
	robotpath.Pause()

	path.Add(100, 100)
	path.Add(150, 50)
	path.Add(300, 100)
	path.Add(350, 350)
	path.Add(300, 600)
	path.Add(200, 650)
	path.Add(100, 600)
	path.Add(50, 350)
	path.Add(100, 100)

	path.PlayPath()
	path.Speed = 50

	p := path.Next()
	robotpath.Position(float64(p.x), float64(p.y))

	touchpad = NewSprite()
	touchpad.AddAnimationByte("default", &gfx.TOUCHPAD, 2000, 1, ebiten.FilterNearest)
	touchpad.CenterCoordonnates = true

	sound.Load()
	sound.PlayBgm(sound.BgmOutThere)

	paint.LoadFonts()

	gv.Life = 100

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

// var robot2 *Sprite
var lastGenEnemyTime time.Time

func GenEnemy() {

	if time.Since(lastGenEnemyTime) > time.Duration(70*time.Millisecond) {
		lastGenEnemyTime = time.Now()

		newsprite := NewSprite()

		newsprite.AddAnimationByteCol("default", &images.E_ROBO1, 100, 1, 8, ebiten.FilterNearest)
		newsprite.Name = "E_ROBO1"
		newsprite.Position(float64(rand.Intn(screenSize.X)), 0)
		newsprite.CenterCoordonnates = true
		newsprite.Pause()

		newsprite.Speed = float64(3 + rand.Intn(6))

		newsprite.Direction = float64(270 - 50 + rand.Intn(100))

		// newsprite.AddEffect(&EffectOptions{Effect: Move, X: 400, Y: 800, Duration: 2000, Repeat: false, GoBack: false})
		// newsprite.AddEffect(&EffectOptions{Effect: Move, X: 200, Y: 400, Duration: 2000, Repeat: false, GoBack: false})

		// newsprite.Start()

		spCount++

		enemyList[strconv.Itoa(spCount)] = newsprite

	}
}

func GenEnemy_level2() {

	if time.Since(lastGenEnemyTime) > time.Duration(210*time.Millisecond) {
		lastGenEnemyTime = time.Now()

		newsprite := NewSprite()
		newsprite.AddAnimationByteCol("default", &images.E_ROBO1, 100, 1, 8, ebiten.FilterNearest)
		newsprite.Name = "E_ROBO1"
		newsprite.Position(float64(rand.Intn(screenSize.X)), 0)
		newsprite.CenterCoordonnates = true
		newsprite.Pause()
		newsprite.Speed = float64(7 + rand.Intn(6))
		newsprite.Direction = float64(270 - 20 + rand.Intn(40))
		spCount++
		enemyList[strconv.Itoa(spCount)] = newsprite
		//===========
		newsprite = NewSprite()
		newsprite.AddAnimationByteCol("default", &images.E_ROBO1, 100, 1, 8, ebiten.FilterNearest)
		newsprite.Name = "E_ROBO1"
		newsprite.Position(0, float64(rand.Intn(screenSize.Y)))
		newsprite.CenterCoordonnates = true
		newsprite.Pause()
		newsprite.Speed = float64(2 + rand.Intn(6))
		newsprite.Direction = float64(0 - 20 + rand.Intn(40))
		spCount++
		enemyList[strconv.Itoa(spCount)] = newsprite

		//===========
		newsprite = NewSprite()
		newsprite.AddAnimationByteCol("default", &images.E_ROBO1, 100, 1, 8, ebiten.FilterNearest)
		newsprite.Name = "E_ROBO1"
		newsprite.Position(float64(screenSize.X), float64(rand.Intn(screenSize.Y)))
		newsprite.CenterCoordonnates = true
		newsprite.Pause()
		newsprite.Speed = float64(2 + rand.Intn(6))
		newsprite.Direction = float64(180 - 20 + rand.Intn(40))
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

		btnFire.SetWH(screenSize.X, screenSize.Y)
		btnFire.rect.x = btnFire.rect.x - 35
		btnBullet.SetWH(screenSize.X, screenSize.Y)
		btnBullet.rect.x = btnBullet.rect.x + 35
		btnDebug.SetWH(screenSize.X, screenSize.Y)
		btnDebug.rect.y = 35

		btnShowBox.SetWH(screenSize.X, screenSize.Y)
		btnShowBox.rect.y = 170

		btnStart.SetWH(screenSize.X, screenSize.Y)
		btnStart.rect.x = screenSize.X/2 - 100
		btnStart.rect.y = screenSize.Y/2 - 50
		btnStart.rect.w = 200
		btnStart.rect.h = 100

		robot.Position(float64(screenSize.X/2), float64(screenSize.Y/2)+100)

		isFirstUpdate = false
		isShowText = false

		gv.Level = 0
	}

	{
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
		if btnFire.GetJoyButton() || ebiten.IsKeyPressed(ebiten.KeySpace) {

			if time.Since(lastBulletTime) > time.Duration(100*time.Millisecond) {
				lastBulletTime = time.Now()

				newsprite := NewSprite()
				//newsprite.AddAnimationByteCol("default", &images.RASER1, 2000, 4, 6, ebiten.FilterNearest)
				newsprite.AddAnimationByte("default", &gfx.EXPLOSION2, 500, 7, ebiten.FilterNearest)
				newsprite.Name = "EXPLOSION2"
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

		if btnDebug.GetJoyButton() || ebiten.IsKeyPressed(ebiten.KeyT) {
			isShowText = !isShowText
		}

		if btnShowBox.GetClicked() || ebiten.IsKeyPressed(ebiten.KeyC) {
			gv.ShowCollisionBox = !gv.ShowCollisionBox
		}

		if btnBullet.GetJoyButton() || ebiten.IsKeyPressed(ebiten.Key1) {
			if time.Since(lastLaserTime) > time.Duration(50*time.Millisecond) {
				lastLaserTime = time.Now()

				newsprite := NewSprite()
				//newsprite.AddAnimationByte("default", &gfx.EXPLOSION3, 500, 9, ebiten.FilterNearest)
				newsprite.AddAnimationByteCol("default", &images.RASER1, 200, 4, 6, ebiten.FilterNearest)
				newsprite.Name = "RASER1"
				newsprite.Position(robot.X, robot.Y)
				//newsprite.AddEffect(&EffectOptions{Effect: Zoom, Zoom: 3, Duration: 6000, Repeat: false, GoBack: true})
				//newsprite.CenterCoordonnates = true
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
		if joytouch.x != 0 && joytouch.y != 0 {
			touchpad.X = float64(joytouch.x)
			touchpad.Y = float64(joytouch.y)
		}
		touchStr = touchStr + "\n" + fmt.Sprintf("PAD:[%f,%f] DEG:[%f]", xx, yy, GetDegreeByXY(xx, yy))

	}

	if gv.Level == 0 { //clickstart

		robot.Show()
		if btnStart.GetClicked() || btnShowBox.GetJoyButton() || ebiten.IsKeyPressed(ebiten.KeyS) {
			gv.Level = 1
			gv.Life = 100
			gv.Score = 0
			spList = make(map[string]*Sprite)
			sound.StopBgm(sound.BgmOutThere)
			sound.PlayBgm(sound.BgmKindBattle)
		}

	}
	if gv.Level == 1 { //level 1
		GenEnemy()

		//paint.DrawText(screen, "攻撃", screenSize.X/2, screenSize.Y/2, color.White, paint.FontSizeXLarge)

		if gv.Life <= 0 {
			gv.Level = 4
			spList = make(map[string]*Sprite)
			sound.StopBgm(sound.BgmKindBattle)
			sound.PlayBgm(sound.BgmOutThere)
		}

		if gv.Score > 200 {
			gv.Level = 2
		}

	}
	if gv.Level == 2 { //level 2
		GenEnemy_level2()
		if gv.Life <= 0 {
			gv.Level = 4
			spList = make(map[string]*Sprite)
			sound.StopBgm(sound.BgmKindBattle)
			sound.PlayBgm(sound.BgmOutThere)
		}
		if gv.Score > 300 {
			gv.Level = 5
		}

	}

	if gv.Level == 3 {
		GenEnemy_level2()

		if gv.Life <= 0 {
			gv.Level = 4
			spList = make(map[string]*Sprite)
			sound.StopBgm(sound.BgmKindBattle)
			sound.PlayBgm(sound.BgmOutThere)
		}

		if gv.Score > 300 {
			gv.Level = 5
		}

	}

	if gv.Level == 4 { //youlost
		robot.Hide()

		if btnStart.GetClicked() || btnShowBox.GetJoyButton() || ebiten.IsKeyPressed(ebiten.KeyS) {
			gv.Level = 0
			robot.Position(float64(screenSize.X/2), float64(screenSize.Y/2)+100)
			spList = make(map[string]*Sprite)
			sound.StopBgm(sound.BgmOutThere)
			sound.PlayBgm(sound.BgmKindBattle)

		}

	}

	if gv.Level == 5 { //youwin

		if btnStart.GetClicked() || btnShowBox.GetJoyButton() || ebiten.IsKeyPressed(ebiten.KeyS) {
			gv.Level = 0
			robot.Position(float64(screenSize.X/2), float64(screenSize.Y/2)+100)
			spList = make(map[string]*Sprite)
			sound.StopBgm(sound.BgmOutThere)
			sound.PlayBgm(sound.BgmKindBattle)

		}

	}

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

	touchStr = touchStr + "\n" + fmt.Sprintf("[%d,%d]", joytouch.x, joytouch.y)
	touchStr = touchStr + "\n" + fmt.Sprintf("[%d]", joytouch.tid)
	touchStr = touchStr + "\n" + fmt.Sprintf("[%f]", robot.Angle)

	touchStr = touchStr + "\n" + fmt.Sprintf("[%v]", len(spList))
	touchStr = touchStr + "\n" + fmt.Sprintf("[%v]", len(enemyList))

	//延迟0.01毫秒
	//time.Sleep(10 * time.Millisecond)
	return nil
}

func drawCollideBox(screen *ebiten.Image, sp *Sprite) {
	cb := sp.GetCollisionBox()
	for _, b := range cb {
		x := b.x + sp.X - sp.GetWidth()/2
		y := b.y + sp.Y - sp.GetHeight()/2
		w := b.w
		h := b.h

		//ebitenutil.DrawRect(screen, b.x+sp.X-sp.GetWidth()/2, b.y+sp.Y-sp.GetHeight()/2, b.w, b.h, color.White) // right

		ebitenutil.DrawLine(screen, x, y, x+w, y, color.White)
		ebitenutil.DrawLine(screen, x+w, y, x+w, y+h, color.White)
		ebitenutil.DrawLine(screen, x, y, x, y+h, color.White)
		ebitenutil.DrawLine(screen, x+w, y+h, x, y+h, color.White)

	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	//打印 hello world 加帧数

	if gv.Level == 1 {

		//paint.DrawText(screen, "Attack !", screenSize.X/2-100, screenSize.Y/2, color.White, paint.FontSizeXLarge)

	}

	if gv.Level == 0 {

		paint.DrawText(screen, "Click Fire to Start", screenSize.X/2-125, screenSize.Y/2, color.White, paint.FontSizeXLarge)

	}
	if gv.Level == 4 {

		paint.DrawText(screen, "YOU LOST！", screenSize.X/2-75, screenSize.Y/2, color.White, paint.FontSizeXLarge)

	}

	if gv.Level == 5 {

		paint.DrawText(screen, "YOU WIN！", screenSize.X/2-75, screenSize.Y/2, color.White, paint.FontSizeXLarge)

	}
	if isShowText {

		mx, my := ebiten.CursorPosition()

		s := fmt.Sprintf("\n\n\n%s\n%s\nFPS : %f %d %d\n%v\n%s\n%s\n%s",
			curpath, home, ebiten.CurrentFPS(), mx, my,
			file_exist(home+"/Library/Caches/output3.txt"), errstr, runtime.GOOS, touchStr)
		ebitenutil.DebugPrint(screen, s)
	}

	paint.DrawText(screen, fmt.Sprintf("Life: %d Score: %d", gv.Life, gv.Score), screenSize.X/2-125, 100, color.Gray{0x99}, paint.FontSizeXLarge)

	//ebitenutil.DebugPrintAt(screen, , screenSize.X-200, 100)

	joytouch.DrawBorders(screen, color.Gray16{0x1111})
	//touchpad.Draw(screen)

	btnFire.DrawBorders(screen, color.Gray16{0x1111})
	btnBullet.DrawBorders(screen, color.Gray16{0x1111})
	btnDebug.DrawBorders(screen, color.Gray16{0x1111})
	btnShowBox.DrawBorders(screen, color.Gray16{0x1111})
	btnStart.DrawBorders(screen, color.Gray16{0x1111})

	if gv.ShowCollisionBox {
		drawCollideBox(screen, robot)
	}
	robot.Draw(screen)

	if gv.Level != 4 {

		for _, j := range spList {
			if gv.ShowCollisionBox {
				drawCollideBox(screen, j)
			}
			j.Draw(screen)

		}
	}

	for _, j := range enemyList {
		if gv.ShowCollisionBox {
			drawCollideBox(screen, j)
		}
		j.Draw(screen)

	}
	for k, j := range effList {
		j.Draw(screen)
		if j.GetStep() >= j.GetTotalStep()-1 {
			delete(effList, k)
		}
	}

	p := path.Next()
	robotpath.Position(p.x, p.y)
	robotpath.Draw(screen)

	if path.LastProgress == path.Totallength {
		path.Reset()
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

	for k, enemy := range enemyList {

		if IsCollideWith(enemy, robot) {

			newsprite := NewSprite()

			newsprite.AddAnimationByte("default", &images.EXPLODE_MED, 500, 8, ebiten.FilterNearest)

			newsprite.Position(enemy.X, enemy.Y)
			newsprite.CenterCoordonnates = true
			newsprite.Speed = enemy.Speed
			newsprite.Direction = enemy.Direction
			newsprite.Start()

			spCount++

			effList[strconv.Itoa(spCount)] = newsprite

			delete(enemyList, k)

			sound.PlaySe(sound.SeKindHit2)

			if gv.Life > 0 {
				gv.Life = gv.Life - 5

			} else {
				gv.Life = 0

				newsprite := NewSprite()
				//newsprite.AddAnimationByteCol("default", &images.RASER1, 2000, 4, 6, ebiten.FilterNearest)

				newsprite.AddAnimationByte("default", &images.EXPLODE_BIG, 2000, 10, ebiten.FilterNearest)

				newsprite.Position(robot.X, robot.Y)
				newsprite.CenterCoordonnates = true

				newsprite.Start()
				spCount++

				effList[strconv.Itoa(spCount)] = newsprite
				sound.PlaySe(sound.SeKindBomb)
				break

			}

		}

		for kshot, shot := range spList {

			if !IsCollideWith(enemy, shot) {
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

			gv.Score++
			delete(enemyList, k)
			delete(spList, kshot)
			sound.PlaySe(sound.SeKindHit2)

		}
	}

}

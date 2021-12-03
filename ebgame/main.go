package ebgame

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/peterSZW/ebdemo/ebgame/gfx"
	"github.com/peterSZW/ebdemo/ebgame/paint"
	"github.com/peterSZW/ebdemo/ebgame/resources/images"
	"github.com/peterSZW/ebdemo/ebgame/sound"
	uuid "github.com/satori/go.uuid"
	"github.com/xiaomi-tc/log15"
	yaml "gopkg.in/yaml.v2"
)

var screenSize image.Point

var homePath string

//var curPath string
//var errstr string

var (
	robot     *Sprite
	robot2    *Sprite
	robotpath *Sprite
	touchpad  *Sprite
)

var joytouch JoyTouch
var btnFire JoyButton
var btnBullet JoyButton
var btnDebug JoyButton
var btnShowBox JoyButton
var btnStart JoyButton

var bulletList sync.Map
var bulletList_cnt int
var ebulletList sync.Map
var ebulletList_cnt int

//var bulletList map[int64]*Sprite
var enemyList map[int64]*Sprite
var effctList map[int64]*Sprite
var spriteCount int64

type Gloable struct {
	EnemyLife          int
	Life               int
	Level              int
	Score              int
	IsShowCollisionBox bool
	IsShowText         bool
}

type GameLogic struct {
}

var gamelogic GameLogic

var path Path

var gv Gloable

const yamlFile = "/Library/Caches/ebdemo.yaml"

type GameConfig struct {
	Account   string `yaml:"Account"`
	Uuid      string `yaml:"Uuid"`
	Token     string `yaml:"Token"`
	HighScore int    `yaml:"HighScore"`
}

var gamecfg GameConfig

func (gcfg *GameConfig) WriteToYaml(src string) {
	// id := uuid.NewV4()
	// ids := id.String()
	// gamecfg.Uuid = ids
	data, err := yaml.Marshal(gcfg) // 第二个表示每行的前缀，这里不用，第三个是缩进符号，这里用tab
	if err != nil {
		log15.Error("WriteToYaml", "err", err)
		return
	}
	err = ioutil.WriteFile(src, data, 0777)
	if err != nil {
		log15.Error("WriteFile", "err", err)
	}
}
func (gcfg *GameConfig) ReadFromYaml(src string) {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		log15.Error("ReadFile", "err", err)
		return
	}
	err = yaml.Unmarshal(content, &gcfg)
	if err != nil {
		log15.Error("Unmarshal", "err", err)
	}
}

const chan_name = "game_room_1"

//初始化
func init() {
	// 	gl.Init()
	// }
	// func (gl *GameLogic) Init() {
	//debug.SetGCPercent(-1)
	log15.Debug("init main")
	// bulletList = make(map[int64]*Sprite)
	enemyList = make(map[int64]*Sprite)
	effctList = make(map[int64]*Sprite)

	homePath = os.Getenv("HOME")
	//curPath = getCurrentDirectory()
	//errstr = ""

	if file_exist(homePath + yamlFile) {
		gamecfg.ReadFromYaml(homePath + yamlFile)

		log15.Debug("READ:", "gamecfg", gamecfg, "path", homePath+yamlFile)
		//读配置

	}

	// rsp, _ := gs.Signin("peta", "abc")
	// fmt.Println(rsp)
	// rsp, _ = gs.Join("peta", "abc")
	// fmt.Println(rsp)

	if gamecfg.Account == "" {
		rand.Seed(time.Now().UnixNano())
		gamecfg.Account = "peter" + strconv.Itoa(rand.Int()%10000)
	}
	if gamecfg.Uuid != "" {
		if beaver_enable {

			rsp, err := beaverChat.GetClient(gamecfg.Uuid)
			log15.Debug("GetClient:", rsp)
			if err != nil {
				log15.Error("", "err", err)
				gamecfg.Uuid = ""
				//log15.Debug(rsp)
			} else {

				gamecfg.Uuid = rsp.ID
				gamecfg.Token = rsp.Token
			}
		}
	}

	if gamecfg.Uuid == "" {
		if beaver_enable {
			_, err := beaverChat.CreateChannel(chan_name, "public")
			if err != nil {
				log15.Error("", "err", err)

			}

			rsp, err := beaverChat.CreateClient([]string{chan_name})
			log15.Debug("CreateClient", "rsp", rsp)
			if err == nil {
				gamecfg.Uuid = rsp.ID
				gamecfg.Token = rsp.ID
			} else {
				log15.Error("", "err", err)
			}
		}

	}
	if aroundus_enable {
		gamecfg.Uuid = uuid.NewV4().String()
	}
	aroundus_ip = "127.0.0.1"

	// NewUser()

	if gamecfg.Uuid != "" {
		gamecfg.WriteToYaml(homePath + yamlFile)
	}

	log15.Debug("write", "gamecfg", gamecfg, "path", homePath+yamlFile)

	rsp0, _ := gs.Signup(gamecfg.Account, "abc")
	log15.Debug("singup", "rsp", rsp0)

	rsp1, _ := gs.Signin(gamecfg.Account, "abc")
	log15.Debug("singin", "rsp", rsp1)

	if rsp1 != nil {
		gamecfg.Token = rsp1.Token
		rsp2, _ := gs.Joinroom(rsp1.Token, "myroom")
		log15.Debug("Joinroom", "rsp", rsp2)
	}

	gs_udp_client()
	gs_udp_Dial()

	//errstr = gamecfg.Uuid

	// f, err := os.Create(homePath + yamlFile) //创建文件

	// if err == nil {
	// 	defer f.Close()
	// 	id := uuid.NewV4()
	// 	ids := id.String()
	// 	_, err = f.WriteString(ids) //写入文件(字节数组)

	// 	f.Sync()
	// }
	// if err != nil {
	// 	errstr = err.Error()

	// }

	robot = NewSprite()
	robot.AddAnimationByte("default", &images.E_ROBO2, 2000, 8, ebiten.FilterNearest)
	robot.Name = "E_ROBO2"
	robot.Position(float64(screenSize.X/2), float64(screenSize.Y/2)+100)
	robot.CenterCoordonnates = true
	robot.Start()
	robot.Pause()

	robot2 = NewSprite()
	robot2.AddAnimationByte("default", &images.E_ROBO2, 2000, 8, ebiten.FilterNearest)
	robot2.Name = "E_ROBO2"
	robot2.Position(float64(screenSize.X-screenSize.X/2), float64(screenSize.Y-(screenSize.Y/2+100)))
	robot2.CenterCoordonnates = true
	robot2.Start()
	robot2.Pause()

	robotpath = NewSprite()
	robotpath.AddAnimationByte("default", &images.E_SHOOTER1, 2000, 8, ebiten.FilterNearest)
	robotpath.Name = "E_SHOOTER1"
	robotpath.Position(float64(screenSize.X/2), float64(screenSize.Y/2)+100)
	robotpath.CenterCoordonnates = true
	robotpath.Pause()

	// path.Add(100, 100)
	// path.Add(200, 50)
	// path.Add(300, 100)

	// path.Add(350, 350)
	// path.Add(300, 600)
	// path.Add(200, 650)
	// path.Add(100, 600)
	// path.Add(50, 350)
	// path.Add(100, 100)

	for r := 0; r <= 360; r = r + 10 {

		x := math.Sin(deg2rad(float64(r))) * 100
		y := math.Cos(deg2rad(float64(r))) * 100
		path.Add(x+200, y+200)

	}

	path.PlayPath()
	path.Speed = 4

	touchpad = NewSprite()
	touchpad.AddAnimationByte("default", &gfx.TOUCHPAD, 2000, 1, ebiten.FilterNearest)
	touchpad.CenterCoordonnates = true

	sound.Load()
	sound.PlayBgm(sound.BgmOutThere)

	paint.LoadFonts()

	gv.Life = 100

	if beaver_enable {
		go ws_client()
	}

	// client()
	// login()
	// go loopUpdate()

}

type Game struct{}

var touchStr string

// const (
// 	widthAsDots = 480.
// )

func (gl *GameLogic) OutofScreen(x, y float64, size float64) bool {

	return x > float64(screenSize.X)+size || x < -size || y > float64(screenSize.Y)+size || y < -size
}
func (gl *GameLogic) limiXY(x, y float64) (float64, float64) {
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

func (gl *GameLogic) GetKeyBoard() (xx, yy float64) {
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

// var robot2 *Sprite
var lastGenEnemyTime time.Time

func (gl *GameLogic) GenEnemy() {

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

		spriteCount++

		//enemyList[spriteCount] = newsprite
		enemyList[(spriteCount)] = newsprite

	}
}

func (gl *GameLogic) GenEnemy_level2() {

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
		spriteCount++
		enemyList[spriteCount] = newsprite
		//===========
		newsprite = NewSprite()
		newsprite.AddAnimationByteCol("default", &images.E_ROBO1, 100, 1, 8, ebiten.FilterNearest)
		newsprite.Name = "E_ROBO1"
		newsprite.Position(0, float64(rand.Intn(screenSize.Y)))
		newsprite.CenterCoordonnates = true
		newsprite.Pause()
		newsprite.Speed = float64(2 + rand.Intn(6))
		newsprite.Direction = float64(0 - 20 + rand.Intn(40))
		spriteCount++
		enemyList[spriteCount] = newsprite

		//===========
		newsprite = NewSprite()
		newsprite.AddAnimationByteCol("default", &images.E_ROBO1, 100, 1, 8, ebiten.FilterNearest)
		newsprite.Name = "E_ROBO1"
		newsprite.Position(float64(screenSize.X), float64(rand.Intn(screenSize.Y)))
		newsprite.CenterCoordonnates = true
		newsprite.Pause()
		newsprite.Speed = float64(2 + rand.Intn(6))
		newsprite.Direction = float64(180 - 20 + rand.Intn(40))
		spriteCount++
		enemyList[spriteCount] = newsprite

	}
}

var isFirstUpdate = true
var lastBulletTime time.Time
var lastLaserTime time.Time
var degree float64

//循环计算
var lastnetxx float64
var lastnetyy float64

var roundx, roundy float64

func (gl *GameLogic) InitSprite() {
	joytouch.SetWH(screenSize.X, screenSize.Y)
	touchpad.X = float64(joytouch.rect.x + joytouch.rect.w/2)
	touchpad.Y = float64(joytouch.rect.y + joytouch.rect.w/2)

	btnFire.SetWH(screenSize.X, screenSize.Y)
	btnFire.rect.x = btnFire.rect.x - 35
	btnBullet.SetWH(screenSize.X, screenSize.Y)
	btnBullet.rect.x = btnBullet.rect.x + 35

	btnDebug.SetPosition(screenSize.X/2-125, screenSize.Y/2-300, 50, 50)
	btnShowBox.SetPosition(screenSize.X/2+75, screenSize.Y/2-300, 50, 50)

	btnStart.SetPosition(screenSize.X/2-125, screenSize.Y/2-32, 240, 32)

	robot.Position(float64(screenSize.X/2), float64(screenSize.Y/2)+100)

	gv.IsShowText = false
	gv.Level = 0
	gv.Life = 100
	gv.EnemyLife = 100
}
func (gl *GameLogic) RemoveAllBullet() {
	bulletList.Range(func(k, v interface{}) bool {
		bulletList.Delete(k)
		return true
	})
	bulletList_cnt = 0

	ebulletList.Range(func(k, v interface{}) bool {
		ebulletList.Delete(k)
		return true
	})
	ebulletList_cnt = 0
}
func (gl *GameLogic) addbullet(x, y, degree float64) {

	newsprite := NewSprite()
	//newsprite.AddAnimationByte("default", &gfx.EXPLOSION3, 500, 9, ebiten.FilterNearest)
	newsprite.AddAnimationByteCol("default", &images.RASER1, 200, 4, 6, ebiten.FilterNearest)
	newsprite.Name = "RASER1"
	newsprite.Position(x, y)
	//newsprite.AddEffect(&EffectOptions{Effect: Zoom, Zoom: 3, Duration: 6000, Repeat: false, GoBack: true})
	//newsprite.CenterCoordonnates = true
	newsprite.Pause()
	newsprite.Step(18)
	newsprite.Speed = 8
	newsprite.Angle = degree          //+90 GetDegreeByXY(xx, yy)
	newsprite.Direction = degree + 90 // GetDegreeByXY(xx, yy) + 90

	//newsprite.Start()

	spriteCount++

	bulletList.Store(strconv.Itoa(int(spriteCount)), newsprite)
	bulletList_cnt++
	//bulletList[spriteCount] = newsprite
}

func (gl *GameLogic) addEnemybullet(x, y, degree float64) {

	newsprite := NewSprite()
	//newsprite.AddAnimationByte("default", &gfx.EXPLOSION3, 500, 9, ebiten.FilterNearest)
	newsprite.AddAnimationByteCol("default", &images.RASER2, 200, 4, 6, ebiten.FilterNearest)
	//newsprite.Name = "RASER2"
	newsprite.Name = "RASER1"
	newsprite.Position(x, y)
	//newsprite.AddEffect(&EffectOptions{Effect: Zoom, Zoom: 3, Duration: 6000, Repeat: false, GoBack: true})
	//newsprite.CenterCoordonnates = true
	newsprite.Pause()
	newsprite.Step(18)
	newsprite.Speed = 8
	newsprite.Angle = degree          //+90 GetDegreeByXY(xx, yy)
	newsprite.Direction = degree + 90 // GetDegreeByXY(xx, yy) + 90

	//newsprite.Start()

	spriteCount++

	ebulletList.Store(strconv.Itoa(int(spriteCount)), newsprite)
	ebulletList_cnt++

}
func (gl *GameLogic) addbullet2(x, y, degree float64) {

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

	spriteCount++

	bulletList.Store(strconv.Itoa(int(spriteCount)), newsprite)
	bulletList_cnt++
}

func (gl *GameLogic) drawCollideBox(screen *ebiten.Image, sp *Sprite) {
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

//
func (gl *GameLogic) checkCollision() {

	ebulletList.Range(func(ikshot, ishot interface{}) bool {
		kshot := ikshot.(string)
		shot := ishot.(*Sprite)

		if IsCollideWith(robot, ishot.(Collider)) {
			//effctList[spriteCount] = newsprite

			ebulletList.Delete(kshot)
			ebulletList_cnt--

			sound.PlaySe(sound.SeKindHit2)

			gv.Life = gv.Life - 5

			// add effect
			newsprite := NewSprite()
			newsprite.AddAnimationByte("default", &images.EXPLODE_MED, 500, 8, ebiten.FilterNearest)
			newsprite.Position(shot.X, shot.Y)
			newsprite.CenterCoordonnates = true
			newsprite.Start()
			spriteCount++
			effctList[spriteCount] = newsprite

			if gv.Life <= 0 {
				gv.Life = 0

				newsprite := NewSprite()

				newsprite.AddAnimationByte("default", &images.EXPLODE_BIG, 2000, 10, ebiten.FilterNearest)

				newsprite.Position(robot.X, robot.Y)
				newsprite.CenterCoordonnates = true

				newsprite.Start()
				spriteCount++

				effctList[spriteCount] = newsprite
				sound.PlaySe(sound.SeKindBomb)

			}
			return false
		}
		return true

	})

	bulletList.Range(func(ikshot, ishot interface{}) bool {
		kshot := ikshot.(string)
		shot := ishot.(*Sprite)

		if IsCollideWith(robot2, ishot.(Collider)) {
			//effctList[spriteCount] = newsprite
			bulletList.Delete(kshot)
			bulletList_cnt--

			sound.PlaySe(sound.SeKindHit2)

			gv.EnemyLife = gv.EnemyLife - 5
			// add effect
			newsprite := NewSprite()
			newsprite.AddAnimationByte("default", &images.EXPLODE_MED, 500, 8, ebiten.FilterNearest)
			newsprite.Position(shot.X, shot.Y)
			newsprite.CenterCoordonnates = true
			newsprite.Start()
			spriteCount++
			effctList[spriteCount] = newsprite

			if gv.EnemyLife <= 0 {
				gv.EnemyLife = 0

				robot2.Hide()

				newsprite := NewSprite()
				newsprite.AddAnimationByte("default", &images.EXPLODE_BIG, 2000, 10, ebiten.FilterNearest)
				newsprite.Position(robot2.X, robot2.Y)
				newsprite.CenterCoordonnates = true
				newsprite.Start()
				spriteCount++
				effctList[spriteCount] = newsprite
				sound.PlaySe(sound.SeKindBomb)

			}
		}
		return true

	})

	for k, enemy := range enemyList {

		if IsCollideWith(enemy, robot) {

			newsprite := NewSprite()

			newsprite.AddAnimationByte("default", &images.EXPLODE_MED, 500, 8, ebiten.FilterNearest)

			newsprite.Position(enemy.X, enemy.Y)
			newsprite.CenterCoordonnates = true
			newsprite.Speed = enemy.Speed
			newsprite.Direction = enemy.Direction
			newsprite.Start()

			spriteCount++

			effctList[spriteCount] = newsprite

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
				spriteCount++

				effctList[spriteCount] = newsprite
				sound.PlaySe(sound.SeKindBomb)
				break

			}

		}

		bulletList.Range(func(ikshot, ishot interface{}) bool {
			kshot := ikshot.(string)
			shot := ishot.(Collider)

			if !IsCollideWith(enemy, shot) {
				//continue
				return true
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

				spriteCount++

				effctList[spriteCount] = newsprite
			}

			gv.Score++
			delete(enemyList, k)
			//delete(bulletList, kshot)
			bulletList.Delete(kshot)
			bulletList_cnt--
			sound.PlaySe(sound.SeKindHit2)
			return true
		})

		//for kshot, shot := range bulletList {

		//}
	}

}

func (gl *GameLogic) FireBullet() {

	//生成子弹
	if btnFire.GetJoyButton() || ebiten.IsKeyPressed(ebiten.KeySpace) {
		if time.Since(lastBulletTime) > time.Duration(100*time.Millisecond) {
			lastBulletTime = time.Now()
			gs_UpdateFire(robot.X, robot.Y, (degree))
			gamelogic.addbullet2(robot.X, robot.Y, degree)

		}
	}

	if btnBullet.GetJoyButton() || ebiten.IsKeyPressed(ebiten.Key1) {
		if time.Since(lastLaserTime) > time.Duration(50*time.Millisecond) {
			lastLaserTime = time.Now()
			gs_UpdateFire(robot.X, robot.Y, degree)
			gamelogic.addbullet(robot.X, robot.Y, degree)
		}
	}
	touchpad.Angle = degree
	if joytouch.x != 0 && joytouch.y != 0 {
		touchpad.X = float64(joytouch.x)
		touchpad.Y = float64(joytouch.y)
	}

	//sent to network

}

func (gl *GameLogic) ShowHideDebug() {
	if btnDebug.GetJoyButton() || ebiten.IsKeyPressed(ebiten.KeyT) {
		gv.IsShowText = !gv.IsShowText
	}

	if btnShowBox.GetClicked() || ebiten.IsKeyPressed(ebiten.KeyC) {
		gv.IsShowCollisionBox = !gv.IsShowCollisionBox
	}
}
func (gl *GameLogic) movePlan() {
	//移动飞机
	xx, yy, _ := joytouch.GetJoyTouchXY()

	if xx == 0 && yy == 0 {
		xx, yy = gamelogic.GetKeyBoard()
	}
	robot.Pause()
	if GetDirectByXY(xx, yy) > 0 {
		robot.Step(GetDirectByXY(xx, yy))

		degree = GetDegreeByXY(xx, yy)
	}

	robot.X = robot.X + xx
	robot.Y = robot.Y + yy
	robot.X, robot.Y = gamelogic.limiXY(robot.X, robot.Y)

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		robot.Angle = robot.Angle + 10
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		robot.Angle = robot.Angle - 10
	}

	roundx = math.Round(robot.X)
	roundy = math.Round(robot.Y)

	if roundx != lastnetxx || roundy != lastnetyy {
		lastnetxx = roundx
		lastnetyy = roundy
		gs_UpdatePosNow(robot.X, robot.Y, robot.GetStep())

		// if beaver_enable {
		// 	beaverChat.PublishChannel(chan_name, fmt.Sprintf(`{"message":"%f,%f","id":"%s"}`, lastnetxx, lastnetyy, gamecfg.Uuid))
		// }
	}
	touchStr = touchStr + "\n" + fmt.Sprintf("PAD:[%f,%f] DEG:[%f]", xx, yy, GetDegreeByXY(xx, yy))

}

func (gl *GameLogic) SetLevel(level int) {
	if level == 0 {
		gv.Life = 100
		gv.EnemyLife = 100

		robot.Position(float64(screenSize.X/2), float64(screenSize.Y/2)+200)
		robot.Show()

		robot2.Position(float64(screenSize.X-screenSize.X/2), float64(screenSize.Y-(screenSize.Y/2+200)))
		robot2.Show()
	}

	if level == 4 {

		robot.Hide()

	}

	if level == 5 {

		robot2.Hide()

	}
	gv.Level = level

}
func (g *Game) Update() error {
	//第一次设置
	if isFirstUpdate {
		isFirstUpdate = false
		gamelogic.InitSprite()
		robot.Show()
		robot2.Show()

	}

	touchStr = ""

	gamelogic.ShowHideDebug()

	if gv.Level == 0 { //clickstart

		robot.Show()
		robot2.Show()
		gamelogic.movePlan()
		if btnStart.GetClicked() || btnShowBox.GetJoyButton() || ebiten.IsKeyPressed(ebiten.KeyS) {
			gv.Level = 1
			gv.Life = 100
			gv.EnemyLife = 100
			gv.Score = 0
			gs_UpdateGameStatus(1)

			gamelogic.RemoveAllBullet()
			sound.StopBgm(sound.BgmOutThere)
			sound.PlayBgm(sound.BgmKindBattle)

			//TODO:send game start PACKAGE
		}

	}
	if gv.Level == 1 { //level 1

		gamelogic.movePlan()
		gamelogic.FireBullet()
		//gamelogic.GenEnemy()

		if gv.Life <= 0 {
			gv.Level = 4
			gs_UpdateGameStatus(5)

			gamelogic.RemoveAllBullet()
			sound.StopBgm(sound.BgmKindBattle)
			sound.PlayBgm(sound.BgmOutThere)
		}

		if gv.EnemyLife <= 0 {
			gv.Level = 5
			gs_UpdateGameStatus(4)

			gamelogic.RemoveAllBullet()
			sound.StopBgm(sound.BgmKindBattle)
			sound.PlayBgm(sound.BgmOutThere)
		}

		// if gv.Score > 200 {
		// 	gv.Level = 2
		// }

	}
	// if gv.Level == 2 { //level 2
	// 	gamelogic.GenEnemy_level2()
	// 	if gv.Life <= 0 {
	// 		gv.Level = 4
	// 		gamelogic.RemoveAllBullet()
	// 		sound.StopBgm(sound.BgmKindBattle)
	// 		sound.PlayBgm(sound.BgmOutThere)
	// 	}
	// 	if gv.Score > 300 {
	// 		gv.Level = 5
	// 	}
	// }

	// if gv.Level == 3 {
	// 	gamelogic.GenEnemy_level2()

	// 	if gv.Life <= 0 {
	// 		gv.Level = 4
	// 		gamelogic.RemoveAllBullet()
	// 		sound.StopBgm(sound.BgmKindBattle)
	// 		sound.PlayBgm(sound.BgmOutThere)
	// 	}

	// 	if gv.Score > 300 {
	// 		gv.Level = 5
	// 	}

	// }

	if gv.Level == 4 { //youlost
		robot.Hide()

		if btnStart.GetClicked() || btnShowBox.GetJoyButton() || ebiten.IsKeyPressed(ebiten.KeyS) {
			gs_UpdateGameStatus(0)

			gamelogic.RemoveAllBullet()
			sound.StopBgm(sound.BgmOutThere)
			sound.PlayBgm(sound.BgmKindBattle)

			gamelogic.SetLevel(0)

		}

	}

	if gv.Level == 5 { //youwin
		gamelogic.movePlan()
		gamelogic.FireBullet()
		if btnStart.GetClicked() || btnShowBox.GetJoyButton() || ebiten.IsKeyPressed(ebiten.KeyS) {
			gs_UpdateGameStatus(0)

			bulletList.Range(func(k, v interface{}) bool {
				bulletList.Delete(k)
				bulletList_cnt--
				return true
			})
			sound.StopBgm(sound.BgmOutThere)
			sound.PlayBgm(sound.BgmKindBattle)

			gamelogic.SetLevel(0)

		}

	}

	//删除越界的子弹对象

	bulletList.Range(func(kk, vv interface{}) bool {
		k := kk.(string)
		v := vv.(*Sprite)
		if gamelogic.OutofScreen(v.X, v.Y, 20) {
			v.Hide()
			bulletList.Delete(k)
			bulletList_cnt--
		}
		return true
	})
	ebulletList.Range(func(kk, vv interface{}) bool {
		k := kk.(string)
		v := vv.(*Sprite)
		if gamelogic.OutofScreen(v.X, v.Y, 20) {
			v.Hide()
			ebulletList.Delete(k)
			ebulletList_cnt--
		}
		return true
	})

	//删除越界的敌人对象
	for k, j := range enemyList {

		if gamelogic.OutofScreen(j.X, j.Y, 20) {
			j.Hide()
			delete(enemyList, k)
		}
	}

	//检查碰撞
	gamelogic.checkCollision()

	//计算路径
	p := path.Next()
	robotpath.Position(p.x, p.y)

	if path.LastProgress == path.Totallength {
		path.Reset()
	}

	//生成字符串
	touchStr = touchStr + "\n" + fmt.Sprintf("joyxy:[%d,%d]", joytouch.x, joytouch.y)
	touchStr = touchStr + "\n" + fmt.Sprintf("tid:[%d]", joytouch.tid)
	touchStr = touchStr + "\n" + fmt.Sprintf("Angle:[%f]", robot.Angle)
	touchStr = touchStr + "\n" + fmt.Sprintf("Level:[%d]", gv.Level)

	touchStr = touchStr + "\n" + fmt.Sprintf("Blist[%d]", bulletList_cnt)
	touchStr = touchStr + "\n" + fmt.Sprintf("eBlist[%d]", ebulletList_cnt)
	touchStr = touchStr + "\n" + fmt.Sprintf("eff[%d]", len(effctList))
	touchStr = touchStr + "\n" + fmt.Sprintf("ene[%d]", len(enemyList))
	// touchStr = touchStr + "\n" + fmt.Sprintf("[%v]", len(enemyList))

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//打印 hello world 加帧数

	if gv.Level == 0 {
		paint.DrawText(screen, "Click Fire to Start", screenSize.X/2-125, screenSize.Y/2, color.White, paint.FontSizeXLarge)
	}
	if gv.Level == 4 {
		paint.DrawText(screen, "YOU LOST！", screenSize.X/2-75, screenSize.Y/2, color.White, paint.FontSizeXLarge)
	}

	if gv.Level == 5 {
		paint.DrawText(screen, "YOU WIN！", screenSize.X/2-75, screenSize.Y/2, color.White, paint.FontSizeXLarge)
	}
	if gv.IsShowText {
		mx, my := ebiten.CursorPosition()
		s := fmt.Sprintf("\n\n\n%s\nFPS : %f %d %d\n%v\n%s\n%s\n%s",
			homePath, ebiten.CurrentFPS(), mx, my,
			file_exist(homePath+yamlFile), gamecfg.Uuid, runtime.GOOS, touchStr)
		ebitenutil.DebugPrint(screen, s)
	}

	paint.DrawText(screen, fmt.Sprintf("Life: %d EnLife: %d Score: %d", gv.Life, gv.EnemyLife, gv.Score), screenSize.X/2-125, 100, color.Gray{0x99}, paint.FontSizeMedium)

	//ebitenutil.DebugPrintAt(screen, , screenSize.X-200, 100)

	joytouch.DrawBorders(screen, color.Gray16{0x1111})
	//touchpad.Draw(screen)

	btnFire.DrawBorders(screen, color.Gray16{0x1111})
	btnBullet.DrawBorders(screen, color.Gray16{0x1111})
	btnDebug.DrawBorders(screen, color.Gray16{0x1111})
	btnShowBox.DrawBorders(screen, color.Gray16{0x1111})

	if gv.IsShowCollisionBox {
		gamelogic.drawCollideBox(screen, robot)
		gamelogic.drawCollideBox(screen, robot2)
		gamelogic.drawCollideBox(screen, robotpath)
	}

	if gv.IsShowText {
		robotpath.Draw(screen)

	}
	robot.Draw(screen)
	robot2.Draw(screen)

	bulletList.Range(func(kk, vv interface{}) bool {

		v := vv.(*Sprite)
		if gv.IsShowCollisionBox {
			gamelogic.drawCollideBox(screen, v)
		}
		v.Draw(screen)
		return true
	})

	ebulletList.Range(func(kk, vv interface{}) bool {

		v := vv.(*Sprite)
		if gv.IsShowCollisionBox {
			gamelogic.drawCollideBox(screen, v)
		}
		v.Draw(screen)
		return true
	})

	for _, j := range enemyList {
		if gv.IsShowCollisionBox {
			gamelogic.drawCollideBox(screen, j)
		}
		j.Draw(screen)

	}

	//画爆炸
	for k, j := range effctList {
		j.Draw(screen)
		if j.GetStep() >= j.GetTotalStep()-1 {
			delete(effctList, k)
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

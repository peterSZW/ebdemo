package mykb

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/keyboard/keyboard"
	rkeyboard "github.com/hajimehoshi/ebiten/v2/examples/resources/images/keyboard"
)

var (
	ScreenWidth  = 400
	ScreenHeight = 880
)

var keyboardImage *ebiten.Image

func init() {
	img, _, err := image.Decode(bytes.NewReader(rkeyboard.Keyboard_png))
	if err != nil {
		log.Fatal(err)
	}

	keyboardImage = ebiten.NewImageFromImage(img)
}

var tx, ty int

type Game struct {
	keys []ebiten.Key
}

func (g *Game) SetWindowSize(width, height int) {
	//screenSize.X = width
	//screenSize.Y = height
	ScreenWidth = width
	ScreenHeight = height

}

type Rect struct {
	x int
	y int
	w int
	h int
}

func isInRect(x int, y int, r Rect) bool {
	if x > r.x && x < r.x+r.w && y > r.y && y < r.y+r.h {
		return true
	} else {
		return false
	}
}

var touchStr string

func (g *Game) Update() error {
	//	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	touches := ebiten.TouchIDs()

	touchStr = ""
	tx = 0
	ty = 0
	if len(touches) > 0 {

		for _, id := range touches {
			x, y := ebiten.TouchPosition(id)

			touchStr = touchStr + fmt.Sprintf("x %d y %d id %d\n", x, y, id)
			tx = x
			ty = y

			// if isInRect(x, y, joybtn.rect) {
			// 	joybtn.tid = tid
			// 	joybtn.x = x
			// 	joybtn.y = y

			// }

		}

	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	const (
		offsetX = 24
		offsetY = 400
	)

	// Draw the base (grayed) keyboard image.
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(offsetX, offsetY)

	op.ColorM.Scale(0.5, 0.5, 0.5, 1)
	screen.DrawImage(keyboardImage, op)
	mx, my := ebiten.CursorPosition()
	s := fmt.Sprintf(" FPS:%f  %d,%d \n %s", ebiten.CurrentFPS(), mx, my, touchStr)
	ebitenutil.DebugPrint(screen, s)

	op = &ebiten.DrawImageOptions{}

	if mx == 0 && my == 0 {
		mx = tx
		my = ty
	}
	for p := 0; p <= 109; p++ {
		op.GeoM.Reset()
		r, ok := keyboard.KeyRect(ebiten.Key(p))
		if !ok {
			continue
		}
		//fmt.Println(r)
		if isInRect(mx-offsetX, my-offsetY, Rect{x: r.Min.X, y: r.Min.Y, w: r.Max.X - r.Min.X, h: r.Max.Y - r.Min.Y}) {
			op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
			op.GeoM.Translate(offsetX, offsetY)
			screen.DrawImage(keyboardImage.SubImage(r).(*ebiten.Image), op)
		}

	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

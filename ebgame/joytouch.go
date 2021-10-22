package ebgame

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type JoyTouch struct {
	x      int
	y      int
	width  int
	height int
	rect   Rect
	tid    int
}

func (this *JoyTouch) GetTid() int {
	return this.tid
}

func (this *JoyTouch) SetWH(w, h int) {
	if h > w {
		this.width = w
		this.height = h
		this.rect.x = 0
		this.rect.w = w / 2
		this.rect.y = h - w/2
		this.rect.h = w / 2
	} else {
		this.width = w
		this.height = h
		this.rect.x = 0
		this.rect.w = w / 2
		this.rect.y = h / 2
		this.rect.h = h / 2
	}

}

func (this *JoyTouch) Press(x, y int, tid int) bool {
	if isInRect(x, y, this.rect) {
		this.tid = tid
		this.x = x
		this.y = y
		return true
	}
	return false
}

func (this *JoyTouch) Move(x, y int) (xx float64, yy float64) {

	xx = 0
	yy = 0
	dis := math.Sqrt(float64((x-this.x)*(x-this.x) + (y-this.y)*(y-this.y)))

	// if dis > 0 {
	// 	xx = float64(x-this.x) / dis
	// 	yy = float64(y-this.y) / dis
	// }
	if dis > 0 {
		xx = float64(x-this.x) / dis * 5
		yy = float64(y-this.y) / dis * 5
	}
	return

}

func (this *JoyTouch) DrawBorders(surface *ebiten.Image, c color.Color) {
	var x, y, x1, y1 float64

	x = float64(this.rect.x)
	y = float64(this.rect.y)

	x1 = x + float64(this.rect.w)
	y1 = y + float64(this.rect.h)

	ebitenutil.DrawLine(surface, x, y, x1, y, c)   // top
	ebitenutil.DrawLine(surface, x, y1, x1, y1, c) // bottom
	ebitenutil.DrawLine(surface, x, y, x, y1, c)   // left
	ebitenutil.DrawLine(surface, x1, y, x1, y1, c) // right
}

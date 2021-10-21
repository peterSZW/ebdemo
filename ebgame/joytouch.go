package ebgame

import (
	"math"
)

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
	this.width = w
	this.height = h
	this.rect.x = 0
	this.rect.w = w / 2
	this.rect.y = h - w/2
	this.rect.h = w / 2

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

	dis := math.Sqrt(float64((x - this.x) ^ 2 + (y - this.y) ^ 2))
	xx = float64(x-this.x) / dis
	yy = float64(y-this.y) / dis
	return

}

package ebgame

import (
        _ "image/png"
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

type Js struct {
        x      int
        y      int
        width  int
        height int
        rect   Rect
        tid    int
}

func (this *Js) GetTid() int {
        return this.tid
}

func (this *Js) SetWH(w, h int) {
        this.width = w
        this.height = h
        this.rect.x = 0
        this.rect.w = w / 2
        this.rect.y = h - w/2
        this.rect.h = w / 2

}
func (this *Js) Press(x, y int, tid int) {
        if isInRect(x, y, this.rect) {
                this.tid = tid
                this.x = x
                this.y = y
        }
}

func (this *Js) Move(x, y int) (xx float64, yy float64) {

        dis := math.Sqrt(float64((x - this.x) ^ 2 + (y - this.y) ^ 2))
        xx = float64(x-this.x) / dis
        yy = float64(y-this.y) / dis
        return

}

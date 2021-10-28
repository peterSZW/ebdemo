package ebgame

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

type JoyButton struct {
	x       int
	y       int
	width   int
	height  int
	rect    Rect
	tid     int
	clicked bool
}

func (this *JoyButton) GetClicked() bool {
	v := this.clicked

	this.clicked = false
	return v
}
func (this *JoyButton) GetTid() int {
	return this.tid
}

func (this *JoyButton) SetWH(w, h int) {
	if h > w {
		this.width = w
		this.height = h
		this.rect.x = w/2 + w/6
		this.rect.w = w / 6
		this.rect.y = h - w/2 + w/6
		this.rect.h = w / 6
	} else {
		this.width = w
		this.height = h
		this.rect.x = w / 2
		this.rect.w = w / 2
		this.rect.y = h / 2
		this.rect.h = h / 2
	}

}

func (this *JoyButton) Press(x, y int, tid int) bool {
	if isInRect(x, y, this.rect) {
		this.tid = tid
		this.x = x
		this.y = y
		this.clicked = true
		return true
	}
	return false
}

// func (this *JoyButton) Move(x, y int) (xx float64, yy float64) {

// 	xx = 0
// 	yy = 0
// 	dis := math.Sqrt(float64((x-this.x)*(x-this.x) + (y-this.y)*(y-this.y)))

// 	// if dis > 0 {
// 	// 	xx = float64(x-this.x) / dis
// 	// 	yy = float64(y-this.y) / dis
// 	// }
// 	if dis > 0 {
// 		xx = float64(x-this.x) / dis * 5
// 		yy = float64(y-this.y) / dis * 5
// 	}
// 	return

// }

func (this *JoyButton) DrawBorders(surface *ebiten.Image, c color.Color) {
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

func (this *JoyButton) GetJoyButton() bool {

	touches := ebiten.TouchIDs()

	isstillpress := false

	if len(touches) > 0 {

		if this.tid != 0 {
			//alread have last press, so we need to find is this touch still on screen
			id := touches[0]
			for _, id = range touches {
				if int(id) == int(this.GetTid()) {
					isstillpress = true
					break

				}
			}

			if isstillpress {

			} else {
				this.tid = 0

				for _, id := range touches {

					x, y := ebiten.TouchPosition(id)

					if this.Press(x, y, int(id)) {
						isstillpress = true
						break
					}

				}

			}

		} else {
			//find new press

			for _, id := range touches {
				x, y := ebiten.TouchPosition(id)
				if this.Press(x, y, int(id)) {
					isstillpress = true

					break
				}

			}

		}

	} else {
		//all reased
		this.tid = 0
		this.x = 0
		this.y = 0

	}
	if isstillpress {
		touchStr = touchStr + "\n" + "FIRE PRESS - " + strconv.Itoa(this.tid)
	}
	return isstillpress
}

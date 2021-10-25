package ebgame

import (
	"fmt"
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
		this.rect.x = 0 + 10
		this.rect.w = w / 2
		this.rect.y = h - w/2 - 10
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

//==================

func (this *JoyTouch) GetJoyTouchXY() (xx, yy float64) {
	touchStr = ""

	touches := ebiten.TouchIDs()

	isstillpress := false

	xx = 0.
	yy = 0.
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
				x, y := ebiten.TouchPosition(id)

				xx, yy = this.Move(x, y)
			} else {
				this.tid = 0
				this.x = 0
				this.y = 0
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
			//find new press

			for _, id := range touches {
				x, y := ebiten.TouchPosition(id)
				if this.Press(x, y, int(id)) {
					isstillpress = true

					break
				}

			}

		}

		for _, id := range touches {

			x, y := ebiten.TouchPosition(id)

			touchStr = touchStr + "\n" + fmt.Sprintf("(%d,%d)", x, y)
		}

	} else {
		//all reased
		this.tid = 0
		this.x = 0
		this.y = 0
		xx = 0
		yy = 0
	}
	if isstillpress {
		touchStr = touchStr + "\n" + "STILL PRESS"
	}

	return xx, yy
}

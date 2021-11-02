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

func (joytouch *JoyTouch) GetTid() int {
	return joytouch.tid
}

func (joytouch *JoyTouch) SetWH(w, h int) {
	if h > w {
		joytouch.width = w
		joytouch.height = h
		joytouch.rect.x = 0 + 20
		joytouch.rect.w = w / 2
		joytouch.rect.y = h - w/2 - 20
		joytouch.rect.h = w / 2
	} else {
		joytouch.width = w
		joytouch.height = h
		joytouch.rect.x = 0
		joytouch.rect.w = w / 2
		joytouch.rect.y = h / 2
		joytouch.rect.h = h / 2
	}

}

func (joytouch *JoyTouch) Press(x, y int, tid int) bool {
	if isInRect(x, y, joytouch.rect) {
		joytouch.tid = tid
		joytouch.x = x
		joytouch.y = y
		return true
	}
	return false
}

func (joytouch *JoyTouch) Move(x, y int) (xx float64, yy float64) {

	xx = 0
	yy = 0
	dis := math.Sqrt(float64((x-joytouch.x)*(x-joytouch.x) + (y-joytouch.y)*(y-joytouch.y)))

	// if dis > 0 {
	// 	xx = float64(x-joytouch.x) / dis
	// 	yy = float64(y-joytouch.y) / dis
	// }
	if dis > 0 {
		xx = float64(x-joytouch.x) / dis * 5
		yy = float64(y-joytouch.y) / dis * 5
	}
	return

}

func (joytouch *JoyTouch) DrawBorders(surface *ebiten.Image, c color.Color) {
	var x, y, x1, y1 float64

	x = float64(joytouch.rect.x)
	y = float64(joytouch.rect.y)

	x1 = x + float64(joytouch.rect.w)
	y1 = y + float64(joytouch.rect.h)

	ebitenutil.DrawLine(surface, x, y, x1, y, c)   // top
	ebitenutil.DrawLine(surface, x, y1, x1, y1, c) // bottom
	ebitenutil.DrawLine(surface, x, y, x, y1, c)   // left
	ebitenutil.DrawLine(surface, x1, y, x1, y1, c) // right
}

//==================

func (joytouch *JoyTouch) GetJoyTouchXY() (xx float64, yy float64, isstillpress bool) {
	touchStr = ""

	touches := ebiten.TouchIDs()

	isstillpress = false

	xx = 0.
	yy = 0.
	if len(touches) > 0 {

		if joytouch.tid != 0 {
			//alread have last press, so we need to find is joytouch touch still on screen
			id := touches[0]
			for _, id = range touches {
				if int(id) == int(joytouch.GetTid()) {
					isstillpress = true
					break

				}
			}

			if isstillpress {
				x, y := ebiten.TouchPosition(id)

				xx, yy = joytouch.Move(x, y)
			} else {
				joytouch.tid = 0
				joytouch.x = 0
				joytouch.y = 0
				//find new press

				for _, id := range touches {

					x, y := ebiten.TouchPosition(id)

					if joytouch.Press(x, y, int(id)) {
						isstillpress = true
						break
					}

				}

			}

		} else {
			//find new press

			for _, id := range touches {
				x, y := ebiten.TouchPosition(id)
				if joytouch.Press(x, y, int(id)) {
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
		joytouch.tid = 0
		joytouch.x = 0
		joytouch.y = 0
		xx = 0
		yy = 0
	}
	if isstillpress {
		touchStr = touchStr + "\n" + "STILL PRESS"
	}

	return xx, yy, isstillpress
}

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

func (joybtn *JoyButton) GetClicked() bool {
	v := joybtn.clicked

	joybtn.clicked = false
	return v
}
func (joybtn *JoyButton) GetTid() int {
	return joybtn.tid
}

func (joybtn *JoyButton) SetWH(w, h int) {
	if h > w {
		joybtn.width = w
		joybtn.height = h
		joybtn.rect.x = w/2 + w/6
		joybtn.rect.w = w / 6
		joybtn.rect.y = h - w/2 + w/6
		joybtn.rect.h = w / 6
	} else {
		joybtn.width = w
		joybtn.height = h
		joybtn.rect.x = w / 2
		joybtn.rect.w = w / 2
		joybtn.rect.y = h / 2
		joybtn.rect.h = h / 2
	}

}

func (joybtn *JoyButton) Press(x, y int, tid int) bool {
	if isInRect(x, y, joybtn.rect) {
		joybtn.tid = tid
		joybtn.x = x
		joybtn.y = y
		joybtn.clicked = true
		return true
	}
	return false
}

// func (joybtn *JoyButton) Move(x, y int) (xx float64, yy float64) {

// 	xx = 0
// 	yy = 0
// 	dis := math.Sqrt(float64((x-joybtn.x)*(x-joybtn.x) + (y-joybtn.y)*(y-joybtn.y)))

// 	// if dis > 0 {
// 	// 	xx = float64(x-joybtn.x) / dis
// 	// 	yy = float64(y-joybtn.y) / dis
// 	// }
// 	if dis > 0 {
// 		xx = float64(x-joybtn.x) / dis * 5
// 		yy = float64(y-joybtn.y) / dis * 5
// 	}
// 	return

// }

func (joybtn *JoyButton) DrawBorders(surface *ebiten.Image, c color.Color) {
	var x, y, x1, y1 float64

	x = float64(joybtn.rect.x)
	y = float64(joybtn.rect.y)

	x1 = x + float64(joybtn.rect.w)
	y1 = y + float64(joybtn.rect.h)

	ebitenutil.DrawLine(surface, x, y, x1, y, c)   // top
	ebitenutil.DrawLine(surface, x, y1, x1, y1, c) // bottom
	ebitenutil.DrawLine(surface, x, y, x, y1, c)   // left
	ebitenutil.DrawLine(surface, x1, y, x1, y1, c) // right
}

func (joybtn *JoyButton) GetJoyButton() bool {

	touches := ebiten.TouchIDs()

	isstillpress := false

	if len(touches) > 0 {

		if joybtn.tid != 0 {
			//alread have last press, so we need to find is joybtn touch still on screen
			//id := touches[0]
			for _, id := range touches {
				if int(id) == int(joybtn.GetTid()) {
					isstillpress = true
					break

				}
			}

			if isstillpress {

			} else {
				joybtn.tid = 0

				for _, id := range touches {

					x, y := ebiten.TouchPosition(id)

					if joybtn.Press(x, y, int(id)) {
						isstillpress = true
						break
					}

				}

			}

		} else {
			//find new press

			for _, id := range touches {
				x, y := ebiten.TouchPosition(id)
				if joybtn.Press(x, y, int(id)) {
					isstillpress = true

					break
				}

			}

		}

	} else {
		//all reased
		joybtn.tid = 0
		joybtn.x = 0
		joybtn.y = 0

	}
	if isstillpress {
		touchStr = touchStr + "\n" + "FIRE PRESS - " + strconv.Itoa(joybtn.tid)
	}
	return isstillpress
}

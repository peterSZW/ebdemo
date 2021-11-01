package ebgame

import (
	"fmt"
	"math"
	"time"
)

type Point struct {
	x float64
	y float64
}

type MinPath struct {
	start  float64
	end    float64
	length float64
}
type Path struct {
	points      []Point
	length      []float64
	Totallength float64

	minpath []MinPath

	Speed        float64
	LastIndex    int
	LastProgress float64
	lastCalc     time.Time
	isStarted    bool
}

func (p *Path) Add(x float64, y float64) {
	p.points = append(p.points, Point{x, y})
}
func (p *Path) PlayPath() {
	if len(p.points) < 2 {
		return
	}
	length := 0.
	lastx := 0.
	lasty := 0.
	for i := 0; i < len(p.points); i++ {
		point := p.points[i]
		fmt.Println(i, point.x, point.y, lastx, lasty)
		if i >= 1 {
			l := math.Sqrt(float64((point.x-lastx)*(point.x-lastx) + (point.y-lasty)*(point.y-lasty)))

			p.length = append(p.length, l)

			p.minpath = append(p.minpath, MinPath{length, length + l, l})

			fmt.Println("========", l, length, length+l)
			p.Totallength = length + l
			length = length + l
		}
		lastx = point.x
		lasty = point.y

	}

}
func (p *Path) Reset() {
	p.isStarted = false
	p.LastProgress = 0
	p.lastCalc = time.Now()

}
func (p *Path) Next() *Point {

	if len(p.points) < 2 {
		return nil
	}
	if p.LastIndex > len(p.points) {
		return nil
	}
	if !p.isStarted {
		p.isStarted = true
		//p.LastIndex = 0
		p.LastProgress = 0
		p.lastCalc = time.Now()
		return &p.points[0]
	}
	//fmt.Println(int(time.Since(p.lastCalc) / 100000000))
	//fmt.Println(time.Since(p.lastCalc))

	//p.LastProgress = p.LastProgress + 1
	p.LastProgress = p.Speed * float64(time.Since(p.lastCalc).Microseconds()) / 100000
	if p.LastProgress > p.Totallength {
		p.LastProgress = p.Totallength

	}

	for i := 0; i < len(p.minpath); i++ {
		mp := p.minpath[i]

		//fmt.Println(p.LastProgress)
		if (mp.start <= p.LastProgress) && (p.LastProgress <= mp.end) {
			pencent := (p.LastProgress - mp.start) / mp.length
			xx := p.points[i].x + ((p.points[i+1].x - p.points[i].x) * (pencent))
			yy := p.points[i].y + ((p.points[i+1].y - p.points[i].y) * (pencent))

			return &Point{xx, yy}
		}

	}
	return nil

}

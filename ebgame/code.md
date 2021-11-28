Draw dot


```
	pointerImage.Fill(color.RGBA{0xff, 0, 0, 0xff})
var pointerImage = ebiten.NewImage(8, 8)

```
```


	explosion2 = NewSprite()
	//	explosion3.AddAnimationByte("default", &gfx.EXPLOSION3, 500, 9, ebiten.FilterNearest)
	explosion2.AddAnimationByte("default", &gfx.EXPLOSION2, 500, 7, ebiten.FilterNearest)

	//explosion3.AddAnimation("default", "gfx/explosion3.png", explosionDuration, 9, ebiten.FilterNearest)
	explosion2.Position(240-10-48, 400/3*2)
	explosion2.Start()


var x float64
var y float64
var xxx float64
var yyy float64

	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(x, y)
	// screen.DrawImage(pointerImage, op)
	// op = &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(x+10, y)
	// screen.DrawImage(pointerImage, op)

```

```
func CalcBallPosition() {
	x = x + xxx
	y = y + yyy

	if x > float64(screenSize.X)-2.0 {
		xxx = -5
	}

	if x < 0 {
		xxx = 5
	}

	if y > float64(screenSize.Y)-2.0 {
		yyy = -5
	}
	if y < 0 {
		yyy = 5
	}
}
```
package ebgame

import (
	"fmt"
	"testing"
	"time"
)

func TestXXX(t *testing.T) {
	var a Path
	a.Add(100, 100)
	a.Add(300, 100)
	a.Add(300, 600)
	a.Add(100, 600)

	a.PlayPath()
	a.Speed = 50

	for a.LastProgress < a.Totallength {
		fmt.Println(a.Next(), a.LastProgress)
		time.Sleep(time.Duration(time.Millisecond * 30))

	}
}

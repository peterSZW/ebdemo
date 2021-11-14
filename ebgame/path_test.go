package ebgame

import (
	"testing"
	"time"

	"github.com/xiaomi-tc/log15"
)

func TestXXX(t *testing.T) {
	var path Path
	path.Add(100, 100)
	path.Add(200, 50)
	path.Add(300, 100)

	path.Add(350, 350)
	path.Add(300, 600)
	path.Add(200, 650)
	path.Add(100, 600)
	path.Add(50, 350)
	path.Add(100, 100)

	path.PlayPath()
	path.Speed = 50

	for path.LastProgress < path.Totallength {
		log15.Debug(path.Next(), path.LastProgress)
		time.Sleep(time.Duration(time.Millisecond * 30))

	}
}

package ebgame

import (
	"testing"

	"github.com/xiaomi-tc/log15"
)

func TestGSUrl(t *testing.T) {

	gs := NewGSConnect("", "http://127.0.0.1:7403")

	rsp0, _ := gs.Signup(gamecfg.Account, "abc")
	log15.Debug("singup", "rsp", rsp0)

	rsp1, _ := gs.Signin(gamecfg.Account, "abc")
	log15.Debug("singin", "rsp", rsp1)

	rsp2, _ := gs.Joinroom(rsp1.Token, "myroom")
	log15.Debug("singin", "rsp", rsp2)
	rsp2.Roomid = "myroom"

	rsp3, _ := gs.Leaveroom(rsp1.Token, rsp2.Roomid)
	log15.Debug("singin", "rsp", rsp3)
}

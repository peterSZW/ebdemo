package ebgame

import (
	"net/http"

	//"github.com/domgolonka/beavergo"
	"github.com/peterSZW/ebdemo/ebgame/internal/beavergo"
	"github.com/xiaomi-tc/log15"
)

var beaver_url string

var beaver_enable bool

func init_() {
	beaver_enable = true

	log15.Debug("init beaver_client")

	beaver_url = "http://192.168.2.218:7800"
	resp, err := http.Get(beaver_url + "/_health")
	if err != nil {
		log15.Error("", "err", err)
		beaver_url = "http://villa.tpddns.cn:7800"
		resp, err = http.Get(beaver_url + "/_health")
		if err != nil {
			log15.Error("", "err", err)
		} else {
			log15.Debug("", "statuscode", resp.StatusCode)
		}

		return
	}
	log15.Debug("", "statuscode", resp.StatusCode)

	token := "6c68b836-6f8e-465e-b59f-89c1db53afca"
	beaverChat = beavergo.NewConnect(token, beaver_url)
	//health, err := beaverChat.HealthCheck()

}

var beaverChat *beavergo.ChatClient

func testmain() {

	// token := "6c68b836-6f8e-465e-b59f-89c1db53afca"
	// chat = beavergo.NewConnect(token, beaver_url)
	health, err := beaverChat.HealthCheck()
	log15.Debug("", "health", health, "err", err)

	{
		key := "app_name"
		value := "val"
		v0, err := beaverChat.CreateConfig(key, value)
		log15.Debug("CreateConfig", v0, err)

		v1, err := beaverChat.GetConfig(key)
		log15.Debug("GetConfig", v1, err)

		v2, err := beaverChat.UpdateConfig(value)
		log15.Debug("UpdateConfig", v2, err)

		v3, err := beaverChat.DeleteConfig(key)
		log15.Debug("DeleteConfig", v3, err)
	}

	channame := "game_room_1"
	beaverChat.CreateChannel(channame, "public")
	channal, _ := beaverChat.GetChannel(channame)
	log15.Debug("", "chan", channal)

	// rsp, _ := beaverChat.CreateClient([]string{})

	// log15.Debug(rsp)

	v3, err := beaverChat.PublishChannel(channame, `{"message":"hello world!!!!!!!!!"}`)
	log15.Debug("", "v3", v3, "err", err)
	//9b7090b7-1028-4299-ba56-7a1423f6c545

	rsp, _ := beaverChat.GetClient("ca59cb90-7a43-4f96-acc3-086205136bf1")
	log15.Debug("", "", rsp)

	v3, _ = beaverChat.DeleteChannel(channame)
	log15.Debug("", "", v3, "", err)
	rsp, _ = beaverChat.GetClient("ca59cb90-7a43-4f96-acc3-086205136bf1")
	log15.Debug("", "rsp", rsp)
	rsp, _ = beaverChat.GetClient("9b7090b7-1028-4299-ba56-7a1423f6c545")
	log15.Debug("", "rsp", rsp)

}
func bv_getclient() string {

	rsp, _ := beaverChat.GetClient("ca59cb90-7a43-4f96-acc3-086205136bf1")
	log15.Debug("", "", rsp)

	return ""

}

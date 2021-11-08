package ebgame

import (
	"fmt"
	"net/http"

	//"github.com/domgolonka/beavergo"
	"github.com/peterSZW/ebdemo/ebgame/internal/beavergo"
)

var beaver_url string

func init() {
	beaver_url = "http://192.168.2.218:7800"
	resp, err := http.Get(beaver_url + "/_health")
	if err != nil {
		fmt.Println(err)
		beaver_url = "http://villa.tpddns.cn:7800"
		resp, err = http.Get(beaver_url + "/_health")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp.StatusCode)
		}

		return
	}
	fmt.Println(resp.StatusCode)

}

func testmain() {
	// url := "http://villa.tpddns.cn:7800"
	token := "6c68b836-6f8e-465e-b59f-89c1db53afca"
	chat := beavergo.NewConnect(token, beaver_url)
	health, err := chat.HealthCheck()
	fmt.Println(health, err)
	key := "app_name"
	value := "val"

	v0, err := chat.CreateConfig(key, value)
	fmt.Println("CreateConfig", v0, err)

	v1, err := chat.GetConfig(key)
	fmt.Println("GetConfig", v1, err)

	v2, err := chat.UpdateConfig(value)
	fmt.Println("UpdateConfig", v2, err)
	v3, err := chat.DeleteConfig(key)
	fmt.Println("DeleteConfig", v3, err)
	channame := "game_room_1"
	chat.CreateChannel(channame, "public")
	channal, _ := chat.GetChannel(channame)
	fmt.Println(channal)

	// rsp, _ := chat.CreateClient([]string{channame})

	// fmt.Println(rsp)

	v3, err = chat.PublishChannel(channame, `{"message":"hello world!!!!!!!!!"}`)
	fmt.Println(v3, err)

}

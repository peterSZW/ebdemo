package ebgame

import (
	"fmt"
	"net/http"

	//"github.com/domgolonka/beavergo"
	"github.com/peterSZW/ebdemo/ebgame/internal/beavergo"
)

var beaver_url string

func init() {
	fmt.Println("init beaver_client")

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

	token := "6c68b836-6f8e-465e-b59f-89c1db53afca"
	chat = beavergo.NewConnect(token, beaver_url)
	//health, err := chat.HealthCheck()

}

var chat *beavergo.ChatClient

func testmain() {

	// token := "6c68b836-6f8e-465e-b59f-89c1db53afca"
	// chat = beavergo.NewConnect(token, beaver_url)
	health, err := chat.HealthCheck()
	fmt.Println(health, err)

	{
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
	}

	channame := "game_room_1"
	chat.CreateChannel(channame, "public")
	channal, _ := chat.GetChannel(channame)
	fmt.Println(channal)

	// rsp, _ := chat.CreateClient([]string{})

	// fmt.Println(rsp)

	v3, err := chat.PublishChannel(channame, `{"message":"hello world!!!!!!!!!"}`)
	fmt.Println(v3, err)
	//9b7090b7-1028-4299-ba56-7a1423f6c545

	rsp, _ := chat.GetClient("ca59cb90-7a43-4f96-acc3-086205136bf1")
	fmt.Println(rsp)

	v3, _ = chat.DeleteChannel(channame)
	fmt.Println(v3, err)
	rsp, _ = chat.GetClient("ca59cb90-7a43-4f96-acc3-086205136bf1")
	fmt.Println(rsp)
	rsp, _ = chat.GetClient("9b7090b7-1028-4299-ba56-7a1423f6c545")
	fmt.Println(rsp)

}
func bv_getclient() string {

	rsp, _ := chat.GetClient("ca59cb90-7a43-4f96-acc3-086205136bf1")
	fmt.Println(rsp)

	return ""

}

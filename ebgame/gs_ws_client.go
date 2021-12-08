package ebgame

// gs websocket

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xiaomi-tc/log15"
)

var gs_addr = flag.String("addr", "gzfjsoft.com:7403", "http service address")

type gs_ClientResp struct {
	From_client string `json:"from_client"`
	To_client   string `json:"to_client"`
	Channel     string `json:"channel"`
	Data        string `json:"data"`
}

type gs_rspData struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}

func gs_handle(s []byte) {
	var msg gs_ClientResp
	var dt gs_rspData
	json.Unmarshal((s), &msg)

	json.Unmarshal([]byte(msg.Data), &dt)
	if dt.Id != gamecfg.Uuid {
		ss := strings.Split(dt.Message, ",")
		x, _ := strconv.ParseFloat(ss[0], 64)
		y, _ := strconv.ParseFloat(ss[1], 64)
		robot2.X = x
		robot2.Y = y

	}

}

var c *websocket.Conn

func gs_ws_send(msg string) {
	if c != nil {
		err := c.WriteMessage(websocket.TextMessage, []byte(msg))
		log15.Error("WriteMessage:", "err", err)
	}
}
func gs_ws_client() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/" + gamecfg.Uuid + "/" + gamecfg.Token}
	u := url.URL{Scheme: "ws", Host: gameserver_ip, Path: "/"}
	log.Printf("connecting to %s", u.String())

	var err error
	c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log15.Error("dial:", "err", err)
		return
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log15.Error("", "err", "read:", err)
				return
			}
			//log.Printf("recv: %s", message)
			gs_handle(message)

		}
	}()

	gs_ws_send(gamecfg.Token)

	// ticker := time.NewTicker(time.Second)
	// defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		// case t := <-ticker.C:
		// log15.Error("timer", "t", t)
		// err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		// if err != nil {
		// 	log15.Error("","err","write:", err)
		// 	return
		// }
		case <-interrupt:
			log15.Error("", "err", "interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log15.Error("", "err", "write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

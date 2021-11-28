package ebgame

import (
	// "github.com/peterSZW/ebdemo/ebgame/internal/gameserverServer/packet"
	// "github.com/peterSZW/ebdemo/ebgame/internal/gameserverServer/player"
	"encoding/json"
	"net"
	"net/http"
	"time"

	// // "github.com/peterSZW/ebdemo/ebgame/internal/gameserverServer/packet"
	// // "github.com/peterSZW/ebdemo/ebgame/internal/gameserverServer/player"
	// "github.com/peterSZW/ebdemo/ebgame/internal/aroundUsServer/packet"
	// "github.com/peterSZW/ebdemo/ebgame/internal/aroundUsServer/player"

	uuid "github.com/satori/go.uuid"
	"github.com/xiaomi-tc/log15"
)

var gameserver_url string
var gameserver_enable bool

func gs_TryApi(url string) bool {

	resp, err := http.Get(url + "/api")
	if err != nil {
		log15.Error("http.Get", "err", err, "url", url)

		return false
	}
	log15.Debug("http.Get", "resp", resp.StatusCode, "url", url)
	return resp.StatusCode == 200

}

var gameserver_ip string

func gs_init() {
	log15.Debug("init gameserver_url")

	beaver_url = "http://127.0.0.1:7403"
	gameserver_ip = "127.0.0.1"

	if !TryApi(beaver_url) {
		beaver_url = "http://192.168.2.218:7403"
		gameserver_ip = "192.168.2.218"

		if !TryApi(beaver_url) {
			beaver_url = "http://villa.tpddns.cn:7403"
			gameserver_ip = "villa.tpddns.cn"
			if !TryApi(beaver_url) {
				return
			}
		}

	}

	gameserver_enable = true
	if gameserver_enable {
		gamecfg.Uuid = uuid.NewV4().String()
	}

}

func gs_getIncomingClientUdp(gs_udpConnection *net.UDPConn) {
	errx := error(nil)
	log15.Debug("Client listen....")

	for errx == nil {
		buffer := make([]byte, 1024)

		size, _, errx := gs_udpConnection.ReadFromUDP(buffer)
		//addr
		if errx != nil {
			log15.Error("", "err", "Cant read packet!", "err", errx)
			time.Sleep(10 * time.Second)

			continue
		}
		//log15.Error("","err",addr)
		data := buffer[:size]

		var dataPacket ServerPacket
		err2 := json.Unmarshal(data, &dataPacket)
		if err2 != nil {
			log15.Error("Unmarshal", "err", err2)
			//log15.Error("","err","Couldn't parse json player data! Skipping iteration!")
			continue
		} else {
			if dataPacket.Data != nil {
				if dataPacket.Type == PositionBroadcast {
					//json.mas
					//log15.Debug(string(data))
					var dataPacket2 TPosReq
					err2 := json.Unmarshal(data, &dataPacket2)
					if err2 != nil {
						log15.Error("Unmarshal", "err", err2)
					}

					//log15.Debug((dataPacket2))

					robot2.X = float64(dataPacket2.X)
					robot2.Y = float64(dataPacket2.Y)
				}
			}

		}

	}

}

// var user Player

func gs_UpdatePosNow() {

	user.Uuid = gamecfg.Uuid
	user.PlayerPosition.X = float32(robot.X)
	user.PlayerPosition.Y = float32(robot.Y)
	packetToSend := StampPacket(user.Uuid, user.PlayerPosition, UpdatePos)

	_, err := packetToSend.SendUdpStream2(gs_udpConnection)
	if err != nil {
		log15.Error("", "err", err)
	}
}
func gs_loopUpdate() {
	for {
		UpdatePosNow()
		time.Sleep(time.Duration(200 * time.Millisecond))
	}

}
func gs_headtbeat() {
	for {
		user.Uuid = gamecfg.Uuid

		packetToSend := StampPacket(user.Uuid, nil, HeartBeat)

		_, err := packetToSend.SendUdpStream2(gs_udpConnection)
		if err != nil {
			log15.Error("", "err", err)
		}
		time.Sleep(time.Duration(10 * time.Second))
	}
}

func gs_udp_Dial() {

	user.Uuid = gamecfg.Uuid

	packetToSend := StampPacket(user.Uuid, user, DialAddr)

	_, err := packetToSend.SendUdpStream2(gs_udpConnection)
	if err != nil {
		log15.Error("", "err", err)
	}
}

var gs_udpConnection *net.UDPConn

func gs_udp_client() {
	// s := "192.168.2.218"
	// p := 7403
	// host = &s
	// port = &p

	// udpAddr, _ := net.ResolveUDPAddr("udp4", *host+":"+strconv.Itoa(*port))
	// user.UdpAddress = udpAddr

	// user.UdpAddress = &net.UDPAddr{
	// 	IP:   net.IPv4(192, 168, 2, 218),
	// 	Port: 7403,
	// }
	ipp, _ := net.ResolveIPAddr("ip", gameserver_ip)

	user.UdpAddress = &net.UDPAddr{
		IP:   ipp.IP,
		Port: 7403,
	}

	var err error

	gs_udpConnection, err = net.DialUDP("udp", nil, user.UdpAddress)

	if err != nil {
		log15.Error("", "err", err)
		time.Sleep(time.Duration(10 * time.Second))

	} else {
		go gs_getIncomingClientUdp(gs_udpConnection)
		go gs_headtbeat()
		//ClientConsoleCLI(gs_udpConnection)
	}

	// }

}

// func ClientConsoleCLI(gs_udpConnection *net.UDPConn) {

// 	for {
// 		var command, parameter string
// 		fmt.Scanln(&command, &parameter)
// 		//commands := strings.Split(strings.Trim(command, "\n\t/\\'\""), " ")
// 		//log15.Debug(command, "|", commands)
// 		switch command {
// 		case "help", "h":
// 			log15.Error("","err","help(h)")
// 			log15.Error("","err","login(lg)")
// 			log15.Error("","err","disconnet(dc) [id]")
// 		case "login", "lg":
// 			packetToSend := StampPacket(user, DialAddr)

// 			_, err := packetToSend.SendUdpStream2(gs_udpConnection)
// 			if err != nil {
// 				log15.Error("","err",err)
// 			}
// 		case "init", "it", "1":
// 			user.Name = "peter"
// 			user.Color = 1
// 			user.Id = 1
// 			packetToSend := StampPacket(user, InitUser)

// 			_, err := packetToSend.SendUdpStream2(gs_udpConnection)
// 			if err != nil {
// 				log15.Error("","err",err)
// 			}
// 		case "2":
// 			user.Name = "leo"
// 			user.Color = 2
// 			user.Id = 2
// 			packetToSend := StampPacket(user, InitUser)

// 			_, err := packetToSend.SendUdpStream2(gs_udpConnection)
// 			if err != nil {
// 				log15.Error("","err",err)
// 			}
// 		case "3":
// 			user.Name = "alex"
// 			user.Color = 3
// 			user.Id = 3
// 			packetToSend := StampPacket(user, InitUser)

// 			_, err := packetToSend.SendUdpStream2(gs_udpConnection)
// 			if err != nil {
// 				log15.Error("","err",err)
// 			}
// 		case "disconnet", "dc":
// 			i, err := strconv.Atoi(parameter)
// 			if err != nil {
// 				log15.Error("","err",err.Error() + "Cant convert to number position")
// 			}

// 			user := Player{Id: i}
// 			packetToSend := StampPacket(user, UserDisconnected)

// 			_, err = packetToSend.SendUdpStream2(gs_udpConnection)
// 			if err != nil {
// 				log15.Error("","err",err)
// 			}
// 		default:
// 			log15.Error("","err","Unknown command")
// 		}
// 	}
// }
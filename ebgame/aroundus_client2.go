package ebgame

import (
	// "github.com/peterSZW/ebdemo/ebgame/internal/aroundUsServer/packet"
	// "github.com/peterSZW/ebdemo/ebgame/internal/aroundUsServer/player"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/imroc/req"
	"github.com/peterSZW/ebdemo/ebgame/internal/aroundUsServer/packet"
	"github.com/peterSZW/ebdemo/ebgame/internal/aroundUsServer/player"
	uuid "github.com/satori/go.uuid"
	"github.com/xiaomi-tc/log15"
)

var host *string
var port *int

var aroundus_url string
var aroundus_enable bool

func TryApi(url string) bool {

	resp, err := http.Get(url + "/api")
	if err != nil {
		log15.Error("http.Get", "err", err, "url", url)

		return false
	}
	log15.Debug("http.Get", "resp", resp.StatusCode, "url", url)
	return resp.StatusCode == 200

}

var aroundus_ip string

func init() {
	log15.Debug("init aroundus_url")

	beaver_url = "http://127.0.0.1:7403"
	aroundus_ip = "127.0.0.1"

	if !TryApi(beaver_url) {
		beaver_url = "http://192.168.2.218:7403"
		aroundus_ip = "192.168.2.218"

		if !TryApi(beaver_url) {
			beaver_url = "http://villa.tpddns.cn:7403"
			aroundus_ip = "villa.tpddns.cn"
			if !TryApi(beaver_url) {
				return
			}
		}

	}

	aroundus_enable = true
	if aroundus_enable {
		gamecfg.Uuid = uuid.NewV4().String()
	}

}

func getIncomingClientUdp(udpConnection *net.UDPConn) {
	errx := error(nil)
	log15.Debug("Client listen....")

	for errx == nil {
		buffer := make([]byte, 1024)

		size, _, errx := udpConnection.ReadFromUDP(buffer)
		//addr
		if errx != nil {
			log15.Error("", "err", "Cant read packet!", "err", errx)
			time.Sleep(10 * time.Second)

			continue
		}
		//log15.Error("","err",addr)
		data := buffer[:size]

		var dataPacket packet.ServerPacket
		err2 := json.Unmarshal(data, &dataPacket)
		if err2 != nil {
			log15.Error("Unmarshal", "err", err2)
			//log15.Error("","err","Couldn't parse json player data! Skipping iteration!")
			continue
		} else {
			if dataPacket.Data != nil {
				if dataPacket.Type == packet.PositionBroadcast {
					//json.mas
					//log15.Debug(string(data))
					var dataPacket2 packet.TNewUserReq
					err2 := json.Unmarshal(data, &dataPacket2)
					if err2 != nil {
						log15.Error("Unmarshal", "err", err2)
					}

					//log15.Debug((dataPacket2))

					robot2.X = float64(dataPacket2.Data.PlayerPosition.X)
					robot2.Y = float64(dataPacket2.Data.PlayerPosition.Y)
				}
			}

		}

	}

}

var user player.Player

func UpdatePosNow() {

	user.Uuid = gamecfg.Uuid
	user.PlayerPosition.X = float32(robot.X)
	user.PlayerPosition.Y = float32(robot.Y)
	packetToSend := packet.StampPacket(user.Uuid, user.PlayerPosition, packet.UpdatePos)

	_, err := packetToSend.SendUdpStream2(udpConnection)
	if err != nil {
		log15.Error("", "err", err)
	}
}
func loopUpdate() {
	for {
		UpdatePosNow()
		time.Sleep(time.Duration(200 * time.Millisecond))
	}

}
func headtbeat() {
	for {
		user.Uuid = gamecfg.Uuid

		packetToSend := packet.StampPacket(user.Uuid, nil, packet.HeartBeat)

		_, err := packetToSend.SendUdpStream2(udpConnection)
		if err != nil {
			log15.Error("", "err", err)
		}
		time.Sleep(time.Duration(10 * time.Second))
	}
}

func Dial() {

	user.Uuid = gamecfg.Uuid

	packetToSend := packet.StampPacket(user.Uuid, user, packet.DialAddr)

	_, err := packetToSend.SendUdpStream2(udpConnection)
	if err != nil {
		log15.Error("", "err", err)
	}
}

func NewUser() {
	reqData := packet.TNewUserReq{Phone: "12"}
	reqData.Uuid = gamecfg.Uuid
	reqData.Type = packet.NewUser
	reqData.Data = &player.Player{Uuid: reqData.Uuid}

	data, _ := req.Post(beaver_url+"/api", req.BodyJSON(&reqData))

	fmt.Print(data, " ")
}

var udpConnection *net.UDPConn

func client() {
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
	ipp, _ := net.ResolveIPAddr("ip", aroundus_ip)

	user.UdpAddress = &net.UDPAddr{
		IP:   ipp.IP,
		Port: 7403,
	}

	var err error

	udpConnection, err = net.DialUDP("udp", nil, user.UdpAddress)

	if err != nil {
		log15.Error("", "err", err)
		time.Sleep(time.Duration(10 * time.Second))

	} else {
		go getIncomingClientUdp(udpConnection)
		go headtbeat()
		//ClientConsoleCLI(udpConnection)
	}

	// }

}

// func ClientConsoleCLI(udpConnection *net.UDPConn) {

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
// 			packetToSend := packet.StampPacket(user, packet.DialAddr)

// 			_, err := packetToSend.SendUdpStream2(udpConnection)
// 			if err != nil {
// 				log15.Error("","err",err)
// 			}
// 		case "init", "it", "1":
// 			user.Name = "peter"
// 			user.Color = 1
// 			user.Id = 1
// 			packetToSend := packet.StampPacket(user, packet.InitUser)

// 			_, err := packetToSend.SendUdpStream2(udpConnection)
// 			if err != nil {
// 				log15.Error("","err",err)
// 			}
// 		case "2":
// 			user.Name = "leo"
// 			user.Color = 2
// 			user.Id = 2
// 			packetToSend := packet.StampPacket(user, packet.InitUser)

// 			_, err := packetToSend.SendUdpStream2(udpConnection)
// 			if err != nil {
// 				log15.Error("","err",err)
// 			}
// 		case "3":
// 			user.Name = "alex"
// 			user.Color = 3
// 			user.Id = 3
// 			packetToSend := packet.StampPacket(user, packet.InitUser)

// 			_, err := packetToSend.SendUdpStream2(udpConnection)
// 			if err != nil {
// 				log15.Error("","err",err)
// 			}
// 		case "disconnet", "dc":
// 			i, err := strconv.Atoi(parameter)
// 			if err != nil {
// 				log15.Error("","err",err.Error() + "Cant convert to number position")
// 			}

// 			user := player.Player{Id: i}
// 			packetToSend := packet.StampPacket(user, packet.UserDisconnected)

// 			_, err = packetToSend.SendUdpStream2(udpConnection)
// 			if err != nil {
// 				log15.Error("","err",err)
// 			}
// 		default:
// 			log15.Error("","err","Unknown command")
// 		}
// 	}
// }

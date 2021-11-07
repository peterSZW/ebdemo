package ebgame

import (
	// "github.com/peterSZW/ebdemo/ebgame/internal/aroundUsServer/packet"
	// "github.com/peterSZW/ebdemo/ebgame/internal/aroundUsServer/player"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/peterSZW/ebdemo/ebgame/internal/aroundUsServer/packet"
	"github.com/peterSZW/ebdemo/ebgame/internal/aroundUsServer/player"

	"log"
)

var host *string
var port *int

func init() {

}

func getIncomingClientUdp(udpConnection *net.UDPConn) {
	err := error(nil)
	fmt.Println("Client listen....")

	for err == nil {
		buffer := make([]byte, 1024)

		size, addr, err := udpConnection.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Cant read packet!", err)
			time.Sleep(10 * time.Second)

			continue
		}
		log.Println(addr)
		data := buffer[:size]

		var dataPacket packet.ClientPacket
		err = json.Unmarshal(data, &dataPacket)
		if err != nil {
			log.Println("Couldn't parse json player data! Skipping iteration!")
			continue
		} else {
			fmt.Println(dataPacket)
		}

	}

}

var user player.Player

func loopUpdate() {
	for {
		user.Name = "leo"
		user.Color = 2
		user.Id = 0
		user.PlayerPosition.X = float32(robot.X)
		user.PlayerPosition.Y = float32(robot.Y)
		packetToSend := packet.StampPacket(user.PlayerPosition, packet.UpdatePos)

		_, err := packetToSend.SendUdpStream2(udpConnection)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Duration(200 * time.Millisecond))
	}

}
func login() {
	user.Name = "peter"
	user.Color = 1
	user.Id = 0
	user.PlayerPosition.X = float32(robot.X)
	user.PlayerPosition.Y = float32(robot.Y)
	packetToSend := packet.StampPacket(user, packet.InitUser)

	_, err := packetToSend.SendUdpStream2(udpConnection)
	if err != nil {
		log.Println(err)
	}
}

var udpConnection *net.UDPConn

func client() {
	// s := "192.168.2.218"
	// p := 7403
	// host = &s
	// port = &p

	// udpAddr, _ := net.ResolveUDPAddr("udp4", *host+":"+strconv.Itoa(*port))
	// user.UdpAddress = udpAddr

	user.UdpAddress = &net.UDPAddr{
		IP:   net.IPv4(192, 168, 2, 218),
		Port: 7403,
	}

	var err error

	udpConnection, err = net.DialUDP("udp", nil, user.UdpAddress)

	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(10 * time.Second))

	} else {
		go getIncomingClientUdp(udpConnection)
		//ClientConsoleCLI(udpConnection)
	}

	// }

}

func ClientConsoleCLI(udpConnection *net.UDPConn) {

	for {
		var command, parameter string
		fmt.Scanln(&command, &parameter)
		//commands := strings.Split(strings.Trim(command, "\n\t/\\'\""), " ")
		//fmt.Println(command, "|", commands)
		switch command {
		case "help", "h":
			log.Println("help(h)")
			log.Println("login(lg)")
			log.Println("disconnet(dc) [id]")
		case "login", "lg":
			packetToSend := packet.StampPacket(user, packet.DialAddr)

			_, err := packetToSend.SendUdpStream2(udpConnection)
			if err != nil {
				log.Println(err)
			}
		case "init", "it", "1":
			user.Name = "peter"
			user.Color = 1
			user.Id = 1
			packetToSend := packet.StampPacket(user, packet.InitUser)

			_, err := packetToSend.SendUdpStream2(udpConnection)
			if err != nil {
				log.Println(err)
			}
		case "2":
			user.Name = "leo"
			user.Color = 2
			user.Id = 2
			packetToSend := packet.StampPacket(user, packet.InitUser)

			_, err := packetToSend.SendUdpStream2(udpConnection)
			if err != nil {
				log.Println(err)
			}
		case "3":
			user.Name = "alex"
			user.Color = 3
			user.Id = 3
			packetToSend := packet.StampPacket(user, packet.InitUser)

			_, err := packetToSend.SendUdpStream2(udpConnection)
			if err != nil {
				log.Println(err)
			}
		case "disconnet", "dc":
			i, err := strconv.Atoi(parameter)
			if err != nil {
				log.Println(err.Error() + "Cant convert to number position")
			}

			user := player.Player{Id: i}
			packetToSend := packet.StampPacket(user, packet.UserDisconnected)

			_, err = packetToSend.SendUdpStream2(udpConnection)
			if err != nil {
				log.Println(err)
			}
		default:
			log.Println("Unknown command")
		}
	}
}

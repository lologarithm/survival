package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/lologarithm/survival/server"
	"github.com/lologarithm/survival/server/messages"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	exit := make(chan int, 1)
	mu := &MockUser{
		alive:           true,
		incoming:        make(chan messages.Packet, 100),
		outgoing:        make(chan messages.Packet, 100),
		partialMessages: map[uint32][]*messages.Multipart{},
	}
	connected := Connect(mu)
	if connected {
		go ReadMessages(mu)
		go RunUser(mu, exit)
	} else {
		c <- os.Interrupt
	}

	<-c
	exit <- 1

	time.Sleep(time.Second)
	log.Printf("Goodbye!")

}

type MockUser struct {
	alive           bool
	conn            *net.UDPConn
	incoming        chan messages.Packet
	outgoing        chan messages.Packet
	partialMessages map[uint32][]*messages.Multipart
}

func ReadMessages(mu *MockUser) {
	widx := 0
	buf := make([]byte, 1024)
	for mu.alive {
		n, err := mu.conn.Read(buf[widx:])
		if err != nil {
			fmt.Printf("Failed to read from conn.")
			fmt.Println(err)
			return
		}
		widx += n
		for {
			pack, ok := messages.NextPacket(buf[:widx])
			if !ok {
				break
			}
			copy(buf, buf[pack.Len():])
			widx -= pack.Len()
			mu.incoming <- pack
		}
	}
}

func Connect(mu *MockUser) bool {
	ra, err := net.ResolveUDPAddr("udp", "localhost:24816")
	if err != nil {
		fmt.Println(err)
		return false
	}
	mu.conn, err = net.DialUDP("udp", nil, ra)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func RunUser(mu *MockUser, exit chan int) {
	packet := messages.NewPacket(messages.CreateAcctMsgType, &messages.CreateAcct{
		Name:     "testuser",
		Password: "testpass",
		CharName: "mahuser",
	})
	_, err := mu.conn.Write(packet.Pack())
	if err != nil {
		fmt.Printf("Failed to write to connection.")
		fmt.Println(err)
	}

	go func() {
		<-exit
		mu.alive = false
	}()

	for mu.alive {
		select {
		case msg := <-mu.incoming:
			ProcessMessage(mu, msg)
		}
	}

	log.Printf("shutting down user.")
	disconn := messages.NewPacket(messages.DisconnectedMsgType, &messages.Disconnected{})
	disb := disconn.Pack()
	mu.conn.Write(disb)
}

func ProcessMessage(mu *MockUser, msg messages.Packet) {
	switch msg.Frame.MsgType {
	case messages.CreateAcctRespMsgType:
		sendmsg(mu, messages.NewPacket(messages.CreateGameMsgType, &messages.CreateGame{
			Name: "newgame",
		}))
	case messages.GameMasterFrameMsgType:
		tmsg := msg.NetMsg.(*messages.GameMasterFrame)
		for _, e := range tmsg.Entities {
			if e.EType == server.CreatureEType {
				fmt.Printf("Ent: %d @ (%d,%d)\n", e.ID, e.X, e.Y)
			}
		}
	case messages.CreateGameRespMsgType:
		sendmsg(mu, messages.NewPacket(messages.MovePlayerMsgType, &messages.MovePlayer{
			EntityID:  0,
			TickID:    0,
			Direction: 90,
		}))
	case messages.MultipartMsgType:
		handleMultipart(mu, msg)
	}

}

func sendmsg(mu *MockUser, msg *messages.Packet) {
	_, err := mu.conn.Write(msg.Pack())
	if err != nil {
		fmt.Printf("Failed to write to connection.")
		fmt.Println(err)
	}
}

func handleMultipart(mu *MockUser, packet messages.Packet) {
	netmsg := packet.NetMsg.(*messages.Multipart)
	// 1. Check if this group already exists
	if _, ok := mu.partialMessages[netmsg.GroupID]; !ok {
		mu.partialMessages[netmsg.GroupID] = make([]*messages.Multipart, netmsg.NumParts)
	}
	// 2. Insert into group
	mu.partialMessages[netmsg.GroupID][netmsg.ID] = netmsg
	// 3. See if group is ready to process
	isReady := true
	for _, p := range mu.partialMessages[netmsg.GroupID] {
		if p == nil {
			isReady = false
			break
		}
	}
	if isReady {
		buf := &bytes.Buffer{}
		for _, p := range mu.partialMessages[netmsg.GroupID] {
			buf.Write(p.Content)
		}
		packet, ok := messages.NextPacket(buf.Bytes())
		if !ok {
			fmt.Printf("lol, failed multipart.... %v", packet)
		}
		mu.incoming <- packet
	}
}

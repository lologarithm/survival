package server

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/lologarithm/survival/server/messages"
)

const (
	port string = ":24816"
)

type Server struct {
	conn             *net.UDPConn
	disconnectPlayer chan Client
	outToNetwork     chan OutgoingMessage
	toGameManager    chan GameMessage
	inputBuffer      []byte
	encryptionKey    *rsa.PrivateKey

	connections map[string]*Client
	gameManager *GameManager
	clientID    uint32
}

func (s *Server) handleMessage() {
	// TODO: Add timeout on read to check for stale connections and add new user connections.
	s.conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	n, addr, err := s.conn.ReadFromUDP(s.inputBuffer)

	if err != nil {
		return
	}
	addrkey := addr.String()
	if n == 0 {
		s.DisconnectConn(addrkey)
	}
	if _, ok := s.connections[addrkey]; !ok {
		s.clientID++
		fmt.Printf("New Connection: %v, ID: %d\n", addrkey, s.clientID)
		s.connections[addrkey] = &Client{
			address:         addr,
			FromNetwork:     NewBytePipe(0),
			FromGameManager: make(chan InternalMessage, 10),
			toGameManager:   s.toGameManager,
			ID:              s.clientID,
		}
		go s.connections[addrkey].ProcessBytes(s.outToNetwork, s.disconnectPlayer)
	}

	s.connections[addrkey].FromNetwork.Write(s.inputBuffer[0:n])
}

func (s *Server) DisconnectConn(addrkey string) {
	// close(s.connections[addrkey].FromNetwork)
	delete(s.connections, addrkey)
}

func (s *Server) sendMessages() {
	for {
		msg := <-s.outToNetwork
		numMsg := (len(msg.msg.RawBytes) / 512) + 1
		st := 0
		b := 512
		for i := 0; i < numMsg; i++ {
			if i == numMsg-1 {
				b = len(msg.msg.RawBytes) % 512
			}
			// TODO: write frames for each piece with a message type of 'continue' maybe. (SEQ IS IMPORTANT!) 		msg.dest.Seq
			if n, err := s.conn.WriteToUDP(msg.msg.RawBytes[st:st+b], msg.dest.address); err != nil {
				fmt.Println("Error: ", err, " Bytes Written: ", n)
			} else {
				log.Printf("Wrote message %d with %d bytes to %v.", msg.msg.Frame.MsgType, n, msg.dest.address)
			}
			st += b
		}
	}
}

func RunServer(exit chan int) {
	toGameManager := make(chan GameMessage, 1024)
	outToNetwork := make(chan OutgoingMessage, 1024)

	manager := NewGameManager(exit, toGameManager, outToNetwork)
	go manager.Run()

	udpAddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Printf("Failed to open UDP port: %s", err)
		os.Exit(1)
	}
	fmt.Println("Now listening on port", port)

	var s Server
	s.connections = make(map[string]*Client, 512)
	s.inputBuffer = make([]byte, 8092)
	s.toGameManager = toGameManager
	s.outToNetwork = outToNetwork
	s.disconnectPlayer = make(chan Client, 512)
	s.conn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Printf("Failed to open UDP port: %s", err)
		os.Exit(1)
	}

	go s.sendMessages()
	fmt.Println("Server Started!")

	run := true
	for run {
		select {
		case <-exit:
			fmt.Println("Killing Socket Server")
			s.conn.Close()
			run = false
		case client := <-s.disconnectPlayer:
			s.DisconnectConn(client.address.String())
		default:
			s.handleMessage()
		}
	}
}

type OutgoingMessage struct {
	dest *Client
	msg  messages.Message
}

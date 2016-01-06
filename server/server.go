package server

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	port string = ":24816"
)

type Server struct {
	conn             *net.UDPConn
	disconnectPlayer chan Client
	outToNetwork     chan NetMessage
	toGameManager    chan GameMessage
	inputBuffer      []byte
	encryptionKey    *rsa.PrivateKey

	connections map[string]*Client
	gameManager *GameManager
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
		fmt.Printf("New Connection: %v\n", addrkey)
		s.connections[addrkey] = &Client{address: addr, fromNetwork: make(chan []byte, 100), fromGameManager: make(chan GameMessage, 10)}
		go s.connections[addrkey].ProcessBytes(s.toGameManager, s.outToNetwork, s.disconnectPlayer)
	}
	s.connections[addrkey].fromNetwork <- s.inputBuffer[0:n]
}

func (s *Server) DisconnectConn(addrkey string) {
	close(s.connections[addrkey].fromNetwork)
	delete(s.connections, addrkey)
}

func (s *Server) sendMessages() {
	for {
		msg := <-s.outToNetwork
		if n, err := s.conn.WriteToUDP(msg.rawBytes, msg.destination.address); err != nil {
			fmt.Println("Error: ", err, " Bytes Written: ", n)
		}
	}
}

func RunServer(exit chan int) {
	toGameManager := make(chan GameMessage, 1024)
	outToNetwork := make(chan NetMessage, 1024)
	fmt.Println("Starting!")

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
	s.inputBuffer = make([]byte, 1024)
	s.toGameManager = toGameManager
	s.outToNetwork = outToNetwork
	s.disconnectPlayer = make(chan Client, 512)
	s.conn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Printf("Failed to open UDP port: %s", err)
		os.Exit(1)
	}

	go s.sendMessages()

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

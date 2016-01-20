package server

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/lologarithm/survival/server/messages"
)

func TestBasicServer(t *testing.T) {
	exit := make(chan int, 10)
	go RunServer(exit)
	time.Sleep(time.Millisecond * 100)
	ra, err := net.ResolveUDPAddr("udp", "localhost:24816")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	conn, err := net.DialUDP("udp", nil, ra)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	packet := messages.NewPacket(messages.LoginMsgType, &messages.Login{
		Name:     "testuser",
		Password: "testpass",
	})
	msgBytes := packet.Pack()
	_, err = conn.Write(msgBytes)
	if err != nil {
		fmt.Printf("Failed to write to connection.")
		fmt.Println(err)
		t.FailNow()
	}
	buf := make([]byte, 512)
	n, err := conn.Read(buf[0:])
	if err != nil {
		fmt.Printf("Failed to read from conn.")
		fmt.Println(err)
		t.FailNow()
	}
	if n < 5 || buf[0] != byte(messages.LoginRespMsgType) {
		fmt.Printf("Incorrect response message!")
		t.FailNow()
	}
	packet = messages.NewPacket(messages.DisconnectedMsgType, &messages.Disconnected{})
	conn.Write(packet.Pack())
	exit <- 1
	conn.Close()

}

func BenchmarkServerParsing(b *testing.B) {
	gamechan := make(chan GameMessage, 100)
	outchan := make(chan OutgoingMessage, 100)
	donechan := make(chan Client, 1)
	fakeClient := &Client{
		address:         &net.UDPAddr{},
		FromNetwork:     NewBytePipe(0),
		FromGameManager: make(chan InternalMessage, 10),
		toGameManager:   gamechan,
		ID:              1,
	}
	go fakeClient.ProcessBytes(outchan, donechan)

	packet := messages.NewPacket(messages.LoginMsgType, &messages.Login{
		Name:     "testuser",
		Password: "testpass",
	})
	msgBytes := packet.Pack()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fakeClient.FromNetwork.Write(msgBytes)
		<-gamechan
	}
}

package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/lologarithm/survival/server/messages"
)

func TestBasicServer(t *testing.T) {
	exit := make(chan int, 1)

	fmt.Println("Starting Server!")
	// Launch server manager
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
	//fmt.Println("Connection Complete")
	messageBytes := new(bytes.Buffer)
	messageBytes.WriteByte(byte(messages.LoginMsgType))
	binary.Write(messageBytes, binary.LittleEndian, uint16(0))

	tbuf := new(bytes.Buffer)
	msg := &messages.Login{
		Name:     "testuser",
		Password: "test",
	}
	msg.Serialize(tbuf)

	binary.Write(messageBytes, binary.LittleEndian, uint16(tbuf.Len()))
	tbuf.WriteTo(messageBytes)
	log.Printf("Writing Msg: %v", messageBytes.Bytes())
	_, err = conn.Write(messageBytes.Bytes())
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
	conn.Write([]byte{255, 0, 0, 0, 0})
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

	messageBytes := new(bytes.Buffer)
	messageBytes.WriteByte(byte(messages.LoginMsgType))
	binary.Write(messageBytes, binary.LittleEndian, uint16(0))
	tbuf := new(bytes.Buffer)
	msg := &messages.Login{
		Name:     "test",
		Password: "test",
	}
	msg.Serialize(tbuf)
	binary.Write(messageBytes, binary.LittleEndian, uint16(tbuf.Len()))
	tbuf.WriteTo(messageBytes)
	msgbytes := messageBytes.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fakeClient.FromNetwork.Write(msgbytes)
		<-gamechan
	}
	log.Printf("test complete!")
}

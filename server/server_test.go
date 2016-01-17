package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"runtime"
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
	netmsg := &messages.Login{
		Name:     "testuser",
		Password: "testpass",
	}
	packet := messages.Packet{
		Frame: messages.Frame{
			MsgType:       messages.LoginMsgType,
			ContentLength: uint16(netmsg.Len()),
		},
		NetMsg: netmsg,
	}
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
	conn.Write([]byte{255, 0, 0, 0, 0})
	conn.Close()
}

func TestCrazyLoad(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	time.Sleep(time.Millisecond * 100)

	for i := 0; i < 3000; i++ {
		go sendMessages()
		time.Sleep(time.Millisecond * 1)
	}
	// run for 2 minutes
	time.Sleep(time.Second * 60)
}

func sendMessages() {
	ra, err := net.ResolveUDPAddr("udp", "localhost:24816")
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.DialUDP("udp", nil, ra)
	if err != nil {
		fmt.Println(err)
		return
	}

	netmsg := &messages.Login{
		Name:     "testuser",
		Password: "testpass",
	}
	packet := messages.Packet{
		Frame: messages.Frame{
			MsgType:       messages.LoginMsgType,
			ContentLength: uint16(netmsg.Len()),
		},
		NetMsg: netmsg,
	}
	msgbytes := packet.Pack()

	go func() {
		widx := 0
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf[widx:])
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
			}
		}
	}()

	for {
		_, err = conn.Write(msgbytes)
		if err != nil {
			fmt.Printf("Failed to write to connection.")
			fmt.Println(err)
		}
		time.Sleep(time.Millisecond * 100)
	}
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

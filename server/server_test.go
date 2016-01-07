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
	conn.Write(messageBytes.Bytes())
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	buf := make([]byte, 512)
	n, err := conn.Read(buf[0:])
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	if n < 5 || buf[0] != 2 {
		t.FailNow()
	}
	conn.Write([]byte{255, 0, 0, 0, 0})
	conn.Close()

}

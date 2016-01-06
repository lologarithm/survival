package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"testing"
	"time"
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
	messageBytes.WriteByte(1)
	binary.Write(messageBytes, binary.LittleEndian, uint16(0))
	binary.Write(messageBytes, binary.LittleEndian, uint16(3))
	messageBytes.WriteByte(97)
	messageBytes.WriteByte(58)
	messageBytes.WriteByte(97)
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

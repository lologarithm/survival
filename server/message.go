package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const frameLen int = 5

type NetMessageType byte

type NetMessage struct {
	rawBytes    []byte
	frame       MessageFrame
	destination *Client
}

func (m *NetMessage) Content() []byte {
	return m.rawBytes[frameLen : frameLen+int(m.frame.length)]
}

func (m *NetMessage) CreateMessageBytes(content []byte) []byte {
	buf := new(bytes.Buffer)
	buf.Grow(5 + len(content))
	buf.WriteByte(byte(m.frame.msgType))
	binary.Write(buf, binary.LittleEndian, m.frame.seq)
	binary.Write(buf, binary.LittleEndian, m.frame.length)
	binary.Write(buf, binary.LittleEndian, content)
	return buf.Bytes()
}

type MessageFrame struct {
	msgType NetMessageType // byte 0, type
	seq     uint16         // byte 1-2, order of message
	length  uint16         // byte 3-4, content length
	from    uint32         // Determined by net addr the request came on.
}

func (mf MessageFrame) String() string {
	return fmt.Sprintf("Type: %d, Seq: %d, CL: %d\n", mf.msgType, mf.seq, mf.length)
}

func ParseFrame(rawBytes []byte) (mf MessageFrame, ok bool) {
	if len(rawBytes) < 5 {
		return
	}
	mf.msgType = NetMessageType(rawBytes[0])
	mf.seq = binary.LittleEndian.Uint16(rawBytes[1:3])
	mf.length = binary.LittleEndian.Uint16(rawBytes[3:5])
	return mf, true
}

type GameMessage interface {
}

type GameMessageValues struct {
	Client *Client
}

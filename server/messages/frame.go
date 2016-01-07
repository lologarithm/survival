package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const FrameLen int = 5

type Message struct {
	RawBytes []byte
	Frame    Frame
}

func (m *Message) Content() []byte {
	return m.RawBytes[FrameLen : FrameLen+int(m.Frame.ContentLength)]
}

func (m *Message) CreateMessageBytes(content []byte) []byte {
	buf := new(bytes.Buffer)
	buf.Grow(5 + len(content))
	buf.WriteByte(byte(m.Frame.MsgType))
	binary.Write(buf, binary.LittleEndian, m.Frame.Seq)
	binary.Write(buf, binary.LittleEndian, m.Frame.ContentLength)
	binary.Write(buf, binary.LittleEndian, content)
	m.RawBytes = buf.Bytes()
	return m.RawBytes
}

type Frame struct {
	MsgType       MessageType // byte 0, type
	Seq           uint16      // byte 1-2, order of message
	ContentLength uint16      // byte 3-4, content length
	From          uint32      // Determined by net addr the request came on.
}

func (mf Frame) String() string {
	return fmt.Sprintf("Type: %d, Seq: %d, CL: %d\n", mf.MsgType, mf.Seq, mf.ContentLength)
}

func ParseFrame(rawBytes []byte) (mf Frame, ok bool) {
	if len(rawBytes) < 5 {
		return
	}
	mf.MsgType = MessageType(rawBytes[0])
	mf.Seq = binary.LittleEndian.Uint16(rawBytes[1:3])
	mf.ContentLength = binary.LittleEndian.Uint16(rawBytes[3:5])
	return mf, true
}

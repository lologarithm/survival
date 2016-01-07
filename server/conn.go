package server

import (
	"net"

	"github.com/lologarithm/survival/server/messages"
)

// TODO: Track 'reliable' messages. Decide which need to be resent.

type Client struct {
	buffer          []byte
	address         *net.UDPAddr
	fromNetwork     chan []byte      // Bytes from client to server
	fromGameManager chan GameMessage // GameMessages from GameManger to client
	toServerManager chan GameMessage // Messages to server manager to join a game
	toGameManager   chan GameMessage // Messages to the game the client is connected to.

	// User  *Client // User attached to this network client
	Seq   uint16
	Alive bool
}

// ProcessBytes accepts raw bytes from a socket and turns them into NetMessage objects and then
// later into GameMessages. These are passed into the GameManager. This function also
// accepts outgoing messages from the GameManager to the client.
func (client *Client) ProcessBytes(toGameManager chan GameMessage, toClient chan OutgoingMessage, disconClient chan Client) {
	client.Alive = true
	for client.Alive {
		select {
		case bytes, ok := <-client.fromNetwork:
			if !ok {
				break
			} else {
				client.buffer = append(client.buffer, bytes...)
				msgFrame, ok := messages.ParseFrame(client.buffer)
				// Only try to parse if we have collected enough bytes.
				if ok && messages.FrameLen+int(msgFrame.Length) <= len(client.buffer) {
					netMsg := messages.ParseNetMessage(msgFrame, client.buffer[messages.FrameLen:messages.FrameLen+int(msgFrame.Length)])
					toGameManager <- GameMessage{net: netMsg, client: client, mtype: msgFrame.MsgType}
					// Remove the used bytes from the buffer.
					newBuffer := make([]byte, len(client.buffer)-messages.FrameLen+int(msgFrame.Length))
					copy(newBuffer, client.buffer[messages.FrameLen+int(msgFrame.Length):])
					client.buffer = newBuffer
				}
			}
		}
	}
}

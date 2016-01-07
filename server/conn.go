package server

import (
	"log"
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
				client.Alive = false
				break
			} else {
				client.buffer = append(client.buffer, bytes...)
				msgFrame, ok := messages.ParseFrame(client.buffer)
				numMsgBytes := messages.FrameLen + int(msgFrame.ContentLength)
				if msgFrame.MsgType == 255 {
					// TODO: this should probably not be a random 1off?
					client.Alive = false
					break
				}
				// Only try to parse if we have collected enough bytes.
				if ok && numMsgBytes <= len(client.buffer) {
					log.Printf("Bytes: %v", client.buffer[:numMsgBytes])
					netMsg := messages.ParseNetMessage(msgFrame, client.buffer[messages.FrameLen:numMsgBytes])
					toGameManager <- GameMessage{net: netMsg, client: client, mtype: msgFrame.MsgType}
					// Remove the used bytes from the buffer.
					newBuffer := make([]byte, len(client.buffer)-numMsgBytes)
					copy(newBuffer, client.buffer[numMsgBytes:])
					client.buffer = newBuffer
				}
			}
		}
	}
}

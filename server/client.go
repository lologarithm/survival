package server

import (
	"log"
	"net"

	"github.com/lologarithm/survival/server/messages"
)

// TODO: Track 'reliable' messages. Decide which need to be resent.

type Client struct {
	ID          uint32 // Unique ID for this session
	buffer      []byte
	wIdx        int
	address     *net.UDPAddr
	fromNetwork chan []byte // Bytes from client to server

	toGameManager   chan GameMessage     // Messages to the main game manager.
	fromGameManager chan InternalMessage //
	toActiveGame    chan GameMessage     // Messages to the current game

	Seq   uint16
	Alive bool
}

// ProcessBytes accepts raw bytes from a socket and turns them into NetMessage objects and then
// later into GameMessages. These are passed into the GameManager. This function also
// accepts outgoing messages from the GameManager to the client.
func (client *Client) ProcessBytes(toClient chan OutgoingMessage, disconClient chan Client) {
	client.toGameManager <- GameMessage{
		client: client,
		net: &messages.Connected{
			IsConnected: 1,
		},
		mtype: messages.ConnectedMsgType,
	}
	client.Alive = true

	var toGame chan GameMessage // used once client is connected to a game. TODO: Shoudl this be cached on the cilent struct?
	for client.Alive {
		select {
		case bytes, ok := <-client.fromNetwork:
			if !ok {
				client.Alive = false
				break
			} else {
				if len(client.buffer) < client.wIdx+len(bytes) {
					newBuffer := make([]byte, client.wIdx+len(bytes))
					copy(newBuffer, client.buffer)
					client.buffer = newBuffer
				}
				copy(client.buffer[client.wIdx:], bytes)
				client.wIdx += len(bytes)
				msgFrame, ok := messages.ParseFrame(client.buffer[:client.wIdx])
				numMsgBytes := messages.FrameLen + int(msgFrame.ContentLength)
				if msgFrame.MsgType == 255 {
					// TODO: this should probably not be a random 1off?
					client.Alive = false
					break
				}
				// Only try to parse if we have collected enough bytes.
				if ok && numMsgBytes <= len(client.buffer) {
					netMsg := messages.ParseNetMessage(msgFrame, client.buffer[messages.FrameLen:numMsgBytes])
					switch msgFrame.MsgType {
					case messages.CreateAcctMsgType, messages.LoginMsgType, messages.CreateCharMsgType, messages.DeleteCharMsgType, messages.ListGamesMsgType, messages.JoinGameMsgType:
						client.toGameManager <- GameMessage{net: netMsg, client: client, mtype: msgFrame.MsgType}
					default:
						if toGame == nil {
							log.Printf("Client sent message type %d(%v) before in a game!", msgFrame.MsgType, netMsg)
							break
						}
						toGame <- GameMessage{net: netMsg, client: client, mtype: msgFrame.MsgType}
					}

					// Remove the used bytes from the buffer.
					copy(client.buffer, client.buffer[numMsgBytes:])
					client.wIdx -= numMsgBytes
				}
			}
		case gmsg := <-client.fromGameManager:
			toGame = gmsg.ToGame
		}
	}
	client.toGameManager <- GameMessage{
		client: client,
		net: &messages.Connected{
			IsConnected: 0,
		},
		mtype: messages.ConnectedMsgType,
	}
}
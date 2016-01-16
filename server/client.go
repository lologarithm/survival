package server

import (
	"log"
	"net"

	"github.com/lologarithm/survival/server/messages"
)

// TODO: Track 'reliable' messages. Decide which need to be resent.

type Client struct {
	ID      uint32 // Unique ID for this session
	buffer  []byte
	wIdx    int
	address *net.UDPAddr

	// These channels are written to by another process
	FromNetwork     chan []byte          // Bytes from client to server
	FromGameManager chan InternalMessage //

	// These channels can be written to in the client but not read from.
	toGameManager chan<- GameMessage // Messages to the main game manager.
	toActiveGame  chan<- GameMessage // Messages to the current game

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

	var toGame chan<- GameMessage // used once client is connected to a game. TODO: Shoudl this be cached on the cilent struct?
	for client.Alive {
		select {
		case bytes, ok := <-client.FromNetwork:
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
				packet, ok := messages.NextPacket(client.buffer[:client.wIdx])
				// Only try to parse if we have collected enough bytes.
				if ok {
					switch packet.Frame.MsgType {
					case messages.CreateAcctMsgType, messages.LoginMsgType, messages.CreateCharMsgType, messages.DeleteCharMsgType, messages.ListGamesMsgType, messages.JoinGameMsgType, messages.CreateGameMsgType:
						client.toGameManager <- GameMessage{net: packet.NetMsg, client: client, mtype: packet.Frame.MsgType}
					default:
						if toGame == nil {
							log.Printf("Client sent message type %d(%v) before in a game!", packet.Frame.MsgType, packet.NetMsg)
							break
						}
						toGame <- GameMessage{net: packet.NetMsg, client: client, mtype: packet.Frame.MsgType}
					}

					// Remove the used bytes from the buffer.
					copy(client.buffer, client.buffer[packet.Len():])
					client.wIdx -= packet.Len()
				}
			}
		case gmsg := <-client.FromGameManager:
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

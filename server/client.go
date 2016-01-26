package server

import (
	"bytes"
	"log"
	"net"
	"sync/atomic"
	"unsafe"

	"github.com/lologarithm/survival/server/messages"
)

// Client represents a single connection to the server.
// Theoretically this could support multiple accounts logged in together (local coop)
type Client struct {
	ID      uint32 // Unique ID for this session
	buffer  []byte
	wIdx    int
	address *net.UDPAddr

	// These channels are written to by another process
	FromNetwork     *BytePipe            // Bytes from client to server
	FromGameManager chan InternalMessage //

	// These channels can be written to in the client but not read from.
	toGameManager chan<- GameMessage // Messages to the main game manager.
	toActiveGame  chan<- GameMessage // Messages to the current game

	Seq     uint16
	GroupID uint32
	Alive   bool
}

// ProcessBytes accepts raw bytes from a socket and turns them into NetMessage objects and then
// later into GameMessages. These are passed into the GameManager. This function also
// accepts outgoing messages from the GameManager to the client.
func (client *Client) ProcessBytes(toClient chan OutgoingMessage, disconClient chan Client) {
	client.toGameManager <- GameMessage{
		client: client,
		net:    &messages.Connected{},
		mtype:  messages.ConnectedMsgType,
	}
	client.Alive = true

	// Used to cache parts of a message.
	// TODO: When should this be cleaned out?
	partialMessages := map[uint32][]*messages.Multipart{}

	var toGame chan<- GameMessage = nil // used once client is connected to a game. TODO: Shoudl this be cached on the cilent struct?

	go func() {
		msg := <-client.FromGameManager
		tmsg := msg.(*ConnectedGame)
		log.Printf("got connected, hooked up toGame channel!")
		atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&toGame)), unsafe.Pointer(&tmsg.ToGame))
	}()
	for client.Alive {
		packet, ok := messages.NextPacket(client.buffer[:client.wIdx])

		if len(client.buffer) < packet.Len() {
			newBuffer := make([]byte, packet.Len()*2)
			copy(newBuffer, client.buffer)
			client.buffer = newBuffer
		}

		if packet.Frame.MsgType == messages.DisconnectedMsgType {
			client.Alive = false
			break
		} else if ok && packet.Frame.MsgType == messages.MultipartMsgType {
			netmsg := packet.NetMsg.(*messages.Multipart)
			// 1. Check if this group already exists
			if _, ok := partialMessages[netmsg.GroupID]; !ok {
				partialMessages[netmsg.GroupID] = make([]*messages.Multipart, netmsg.NumParts)
			}
			// 2. Insert into group
			partialMessages[netmsg.GroupID][netmsg.ID] = netmsg
			// 3. See if group is ready to process
			isReady := true
			for _, p := range partialMessages[netmsg.GroupID] {
				if p == nil {
					isReady = false
					break
				}
			}
			if isReady {
				buf := &bytes.Buffer{}
				for _, p := range partialMessages[netmsg.GroupID] {
					buf.Write(p.Content)
				}
				packet, ok = messages.NextPacket(buf.Bytes())
			}
		} else if !ok || packet.Len() > client.wIdx {
			// This means we need more data still.
			n := client.FromNetwork.Read(client.buffer[client.wIdx:])
			client.wIdx += n
			continue
		}
		// Only try to parse if we have collected enough bytes.
		if ok {
			switch packet.Frame.MsgType {
			case messages.CreateAcctMsgType, messages.LoginMsgType, messages.ListGamesMsgType, messages.JoinGameMsgType, messages.CreateGameMsgType:
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
	client.toGameManager <- GameMessage{
		client: client,
		net:    &messages.Disconnected{},
		mtype:  messages.ConnectedMsgType,
	}
}

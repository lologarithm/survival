package server

import (
	"bytes"
	"log"
	"net"
	"sync/atomic"
	"time"
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
	lastMsg int64

	// These channels are written to by another process
	FromNetwork     *BytePipe // Bytes from client to server
	ToNetwork       chan OutgoingMessage
	FromGameManager chan InternalMessage //

	// These channels can be written to in the client but not read from.
	toGameManager chan<- GameMessage // Messages to the main game manager.
	activeGame    *clientGame

	Seq     uint16
	GroupID uint32
	Alive   bool
}

type clientGame struct {
	toGame chan<- GameMessage
	id     uint32
}

// ProcessBytes accepts raw bytes from a socket and turns them into NetMessage objects and then
// later into GameMessages. These are passed into the GameManager. This function also
// accepts outgoing messages from the GameManager to the client.
func (client *Client) ProcessBytes(disconClient chan Client) {
	client.toGameManager <- GameMessage{
		client: client,
		net:    &messages.Connected{},
		mtype:  messages.ConnectedMsgType,
	}
	client.Alive = true
	client.lastMsg = time.Now().UTC().Unix()
	// Used to cache parts of a message.
	// TODO: When should this be cleaned out?
	partialMessages := map[uint32][]*messages.Multipart{}

	go func() {
		for {
			select {
			case msg := <-client.FromGameManager:
				switch tmsg := msg.(type) {
				case ConnectedGame:
					activeGame := &clientGame{
						toGame: tmsg.ToGame,
						id:     tmsg.ID,
					}
					atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&client.activeGame)), unsafe.Pointer(activeGame))
					log.Printf("Got connected, hooked up toGame channel!")
				}
			case <-time.After(time.Second * 10):
				if !client.Alive {
					log.Printf("Client %d: no longer alive.", client.ID)
					return
				}
				client.ToNetwork <- NewOutgoingMsg(client, messages.HeartbeatMsgType, &messages.Heartbeat{
					Time: time.Now().UTC().UnixNano(),
				})
				// If after 60 seconds we haven't gotten any messages, shut er down!
				lastMsg := time.Unix(atomic.LoadInt64(&client.lastMsg), 0)
				if time.Now().UTC().Sub(lastMsg).Seconds() >= 60 {
					client.FromNetwork.Close()
					return
				}
			}
		}

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
			if n == 0 {
				log.Printf("got 0 byte message from client, shutten er down!")
				client.Alive = false
				break // Break out of alive!
			}
			atomic.StoreInt64(&client.lastMsg, time.Now().UTC().Unix())
			client.wIdx += n
			continue
		}
		// Only try to parse if we have collected enough bytes.
		if ok {
			switch packet.Frame.MsgType {
			case messages.CreateAcctMsgType, messages.LoginMsgType, messages.ListGamesMsgType, messages.JoinGameMsgType, messages.CreateGameMsgType:
				client.toGameManager <- GameMessage{net: packet.NetMsg, client: client, mtype: packet.Frame.MsgType}
			default:
				if client.activeGame == nil {
					log.Printf("Client sent message (%d:%v) before in a game!", packet.Frame.MsgType, packet.NetMsg)
					break
				}
				client.activeGame.toGame <- GameMessage{net: packet.NetMsg, client: client, mtype: packet.Frame.MsgType}
			}

			// Remove the used bytes from the buffer.
			copy(client.buffer, client.buffer[packet.Len():])
			client.wIdx -= packet.Len()
		}
	}
	log.Printf("  shutdown client msg parser: %d\n", client.ID)
	client.toGameManager <- GameMessage{
		client: client,
		net:    &messages.Disconnected{},
		mtype:  messages.DisconnectedMsgType,
	}
	disconClient <- *client
	close(client.FromGameManager)
}

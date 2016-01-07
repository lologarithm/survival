package server

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/lologarithm/survival/server/messages"
)

// GameStatus type used for setting status of a game
type GameStatus byte

// Game statuses
const (
	Unknown GameStatus = 0
	Running GameStatus = iota
)

// GameManager manages all connected users and games.
type GameManager struct {
	// Player data
	Clients     map[uint32]*Client
	Games       map[uint32]*Game
	NextGameID  uint32
	FromNetwork chan GameMessage
	FromGames   chan GameMessage
	ToNetwork   chan OutgoingMessage
	Exit        chan int
}

func NewGameManager(exit chan int, fromNetwork chan GameMessage, toNetwork chan OutgoingMessage) *GameManager {
	gm := &GameManager{
		Clients:     make(map[uint32]*Client, 1),
		Games:       make(map[uint32]*Game, 1),
		FromGames:   make(chan GameMessage, 100),
		FromNetwork: fromNetwork,
		ToNetwork:   toNetwork,
		Exit:        exit,
	}
	return gm
}

// Game represents a single game
type Game struct {
	// Player data
	Clients           map[uint32]*Client
	IntoServerManager chan GameMessage
	FromNetwork       chan GameMessage // Messages from players.
	Exit              chan int
	Status            GameStatus
}

func (gm *GameManager) Run() {
	for {
		select {
		case netMsg := <-gm.FromNetwork:
			gm.ProcessNetMsg(netMsg)
		case gMsg := <-gm.FromGames:
			gm.ProcessGameMsg(gMsg)
		case <-gm.Exit:
			for _, game := range gm.Games {
				game.Exit <- 1
			}
			return
		}
	}
}

func (gm *GameManager) ProcessNetMsg(msg GameMessage) {
	switch msg.mtype {
	case messages.LoginMsgType:
		tmsg := msg.net.(*messages.Login)
		log.Printf("Logging in user: %s", tmsg.Name)
		lr := messages.LoginResponse{
			Success: 1,
		}
		buf := new(bytes.Buffer)
		lr.Serialize(buf)
		frame := messages.Frame{
			MsgType:       messages.LoginResponseMsgType,
			Seq:           1,
			ContentLength: uint16(buf.Len()),
		}
		// TODO: automate making this.
		resp := &OutgoingMessage{
			dest: msg.client,
			msg: messages.Message{
				Frame: frame,
			},
		}
		resp.msg.CreateMessageBytes(buf.Bytes())
		gm.ToNetwork <- *resp
	}
}

func (gm *GameManager) ProcessGameMsg(msg GameMessage) {
	switch msg.mtype {
	}
}

// Run starts the game!
func (gm *Game) Run() {
	waiting := true
	for {
		timeout := time.Millisecond * 50
		waiting = true
		for waiting {
			select {
			case <-time.After(timeout):
				waiting = false
				break
			case msg := <-gm.FromNetwork:
				fmt.Printf("GameManager: Received message: %T\n", msg)
				switch msg.mtype {
				default:
					fmt.Println("GameManager.go:RunGame(): UNKNOWN MESSAGE TYPE: %T", msg)
				}
			case <-gm.Exit:
				fmt.Println("EXITING Game Manager")
				return
			}
		}
		fmt.Printf("Sending client update!\n")
	}
}

// Setup sets up the game!
func (gm *Game) Setup() {
}

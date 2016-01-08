package server

import (
	"bytes"
	"log"

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
	Users       map[uint32]*User
	Games       map[uint32]*Game
	NextGameID  uint32
	FromNetwork chan GameMessage
	FromGames   chan GameMessage
	ToNetwork   chan OutgoingMessage
	Exit        chan int

	// Temp junk to make this crap work
	Accounts map[uint32]*Account
}

func NewGameManager(exit chan int, fromNetwork chan GameMessage, toNetwork chan OutgoingMessage) *GameManager {
	gm := &GameManager{
		Users:       make(map[uint32]*User, 1),
		Games:       make(map[uint32]*Game, 1),
		FromGames:   make(chan GameMessage, 100),
		FromNetwork: fromNetwork,
		ToNetwork:   toNetwork,
		Exit:        exit,
	}
	return gm
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
	case messages.CreateAccountMsgType:
		// TODO: Make a temp account here.
	case messages.CreateCharacterMsgType:
		// TODO: make a character for given account.
	case messages.LoginMsgType:
		gm.loginUser(msg)
	case messages.JoinGameMsgType:
		// TODO: for now just join the default game?
	case messages.CreateGameMsgType:
		// TODO: Make a new game!
		// Then the user that created it joins it!
	default:
		// These messages probably go to a game?
	}

}

func (gm *GameManager) loginUser(msg GameMessage) {
	tmsg := msg.net.(*messages.Login)
	log.Printf("Logging in user: %s", tmsg.Name)
	lr := messages.LoginResponse{
		Success: 1,
	}
	// TODO: automate making this.
	buf := new(bytes.Buffer)
	lr.Serialize(buf)
	frame := messages.Frame{
		MsgType:       messages.LoginResponseMsgType,
		Seq:           1,
		ContentLength: uint16(buf.Len()),
	}
	resp := &OutgoingMessage{
		dest: msg.client,
		msg: messages.Message{
			Frame: frame,
		},
	}
	resp.msg.CreateMessageBytes(buf.Bytes())
	gm.ToNetwork <- *resp
}

func (gm *GameManager) ProcessGameMsg(msg GameMessage) {
	switch msg.mtype {
	}
}

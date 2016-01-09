package server

import (
	"bytes"
	"log"
	"math"

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
	Users       []*User
	Games       []*Game
	NextGameID  uint32
	FromNetwork chan GameMessage
	FromGames   chan GameMessage
	ToNetwork   chan OutgoingMessage
	Exit        chan int

	// Temp junk to make this crap work
	Accounts   []*Account
	AcctByName map[string]*Account
}

func NewGameManager(exit chan int, fromNetwork chan GameMessage, toNetwork chan OutgoingMessage) *GameManager {
	gm := &GameManager{
		Users:       make([]*User, math.MaxUint16),
		Games:       make([]*Game, math.MaxUint16),
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

// ProcessNetMsg is the method by which the game manager can deal with incoming messages from the network.
func (gm *GameManager) ProcessNetMsg(msg GameMessage) {
	switch msg.mtype {
	case messages.ConnectedMsgType:
		gm.handleConnection(msg)
	case messages.CreateAcctMsgType:
		gm.createAccount(msg)
	case messages.CreateCharMsgType:
		gm.createCharacter(msg)
	case messages.LoginMsgType:
		gm.loginUser(msg)
	case messages.JoinGameMsgType:
		// TODO: for now just join the default game?
	case messages.CreateGameMsgType:
		// TODO: Make a new game!
		// Then the user that created it joins it!
	case messages.ListGamesMsgType:
		gameList := &messages.ListGamesResp{
			IDs:   []uint32{},
			Names: []string{},
		}
		for idx, g := range gm.Games {
			gameList.IDs = append(gameList.IDs, uint32(idx))
			gameList.Names = append(gameList.Names, g.Name)
		}
		resp := NewOutgoingMsg(msg.client, messages.ListGamesRespMsgType, gameList)
		gm.ToNetwork <- resp

	default:
		// These messages probably go to a game?
		// TODO: Probably have a direct conn to a game from the *Client
	}

}

func (gm *GameManager) handleConnection(msg GameMessage) {
	netmsg := msg.net.(*messages.Connected)
	// First make sure this is a new connection.
	isNew := gm.Users[msg.client.ID] == nil
	if netmsg.IsConnected == 0 && !isNew {
		gm.Users[msg.client.ID] = nil
	} else if netmsg.IsConnected == 1 && isNew {
		gm.Users[msg.client.ID] = &User{
			Client: msg.client,
		}
	}
}

func (gm *GameManager) createCharacter(msg GameMessage) {
	netmsg := msg.net.(*messages.CreateChar)
	ac := &messages.CreateCharResp{
		AccountID: netmsg.AccountID,
		Name:      netmsg.Name,
		ID:        0,
	}
	var acct *Account
	// 1. Validate user is logged in as account specified.
	for _, acct := range gm.Users[msg.client.ID].Accounts {
		if acct.ID == netmsg.AccountID {

		}
	}
	if acct != nil {
		char := &Character{
			Name:  netmsg.Name,
			Items: []*Item{},
		}
		acct.Characters = append(acct.Characters, char)
	}
	resp := NewOutgoingMsg(msg.client, messages.CreateCharRespMsgType, ac)
	gm.ToNetwork <- resp
}

func (gm *GameManager) createAccount(msg GameMessage) {
	netmsg := msg.net.(*messages.CreateAcct)
	ac := &messages.CreateAcctResp{
		AccountID: 0,
		Name:      netmsg.Name,
	}
	found := false
	for _, acc := range gm.Accounts {
		if acc.Name == netmsg.Name {
			found = true
			break
		}
	}
	if !found {
		ac.AccountID = uint32(len(gm.Accounts))
		gm.Accounts[ac.AccountID] = &Account{
			ID:         ac.AccountID,
			Name:       netmsg.Name,
			Password:   netmsg.Password,
			Characters: []*Character{},
		}
		// TODO: login the new account.
	}

	resp := NewOutgoingMsg(msg.client, messages.CreateAcctRespMsgType, ac)
	gm.ToNetwork <- resp
}

func (gm *GameManager) loginUser(msg GameMessage) {
	tmsg := msg.net.(*messages.Login)
	log.Printf("Logging in account: %s", tmsg.Name)
	lr := messages.LoginResp{
		Success: 0,
		Name:    tmsg.Name,
	}
	for _, acct := range gm.Accounts {
		if acct.Name == tmsg.Name && acct.Password == tmsg.Password {
			lr.AccountID = acct.ID
			lr.Characters = make([]*messages.Character, len(acct.Characters))
			for idx, ch := range acct.Characters {
				lr.Characters[idx] = &messages.Character{
					Name: ch.Name,
					ID:   ch.ID,
				}
			}
		}
	}
	resp := NewOutgoingMsg(msg.client, messages.LoginRespMsgType, &lr)
	gm.ToNetwork <- resp
}

func (gm *GameManager) ProcessGameMsg(msg GameMessage) {
	switch msg.mtype {
	}
}

// NewOutgoingMsg creates a new message that can be sent to a specific client.
func NewOutgoingMsg(dest *Client, tp messages.MessageType, msg messages.Net) OutgoingMessage {
	buf := new(bytes.Buffer)
	msg.Serialize(buf)
	frame := messages.Frame{
		MsgType:       tp,
		Seq:           1,
		ContentLength: uint16(buf.Len()),
	}
	resp := OutgoingMessage{
		dest: dest,
		msg: messages.Message{
			Frame: frame,
		},
	}
	resp.msg.CreateMessageBytes(buf.Bytes())
	return resp
}

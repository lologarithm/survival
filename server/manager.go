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
	Users      []*User
	Games      map[uint32]*Game
	NextGameID uint32 // TODO: this shouldn't just be a number..

	FromGames   chan GameMessage // Manager reads this only, all games created write only
	FromNetwork <-chan GameMessage
	ToNetwork   chan<- OutgoingMessage
	Exit        chan int

	// Temp junk to make this crap work
	Accounts   []*Account
	Characters []*Character
	CharID     uint32
	AccountID  uint32
	AcctByName map[string]*Account
}

func NewGameManager(exit chan int, fromNetwork chan GameMessage, toNetwork chan OutgoingMessage) *GameManager {
	gm := &GameManager{
		Users:       make([]*User, math.MaxUint16),
		Games:       map[uint32]*Game{},
		FromGames:   make(chan GameMessage, 100),
		FromNetwork: fromNetwork,
		ToNetwork:   toNetwork,
		Exit:        exit,
		Accounts:    make([]*Account, math.MaxUint16),
		Characters:  make([]*Character, math.MaxUint16),
		AcctByName:  map[string]*Account{},
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
		gm.createGame(msg)
		// TODO: the user that created it joins it!
	case messages.ListGamesMsgType:
		gameList := &messages.ListGamesResp{
			IDs:   []uint32{},
			Names: []string{},
		}
		for key, g := range gm.Games {
			gameList.IDs = append(gameList.IDs, uint32(key))
			gameList.Names = append(gameList.Names, g.Name)
		}
		resp := NewOutgoingMsg(msg.client, messages.ListGamesRespMsgType, gameList)
		gm.ToNetwork <- resp

	default:
		// These messages probably go to a game?
		// TODO: Probably have a direct conn to a game from the *Client
	}

}

func (gm *GameManager) createGame(msg GameMessage) {
	cgm := msg.net.(*messages.CreateGame)

	netchan := make(chan GameMessage, 100)
	g := NewGame(cgm.Name, gm.FromGames, netchan)
	g.SpawnChunk(0, 0)
	gm.NextGameID++
	gm.Games[gm.NextGameID] = g
	cgr := &messages.CreateGameResp{
		Name:     cgm.Name,
		ID:       gm.NextGameID,
		Seed:     g.Seed,
		Entities: g.World.EntitiesMsg(),
	}
	msg.client.FromGameManager <- InternalMessage{
		ToGame: netchan,
	}
	resp := NewOutgoingMsg(msg.client, messages.CreateGameRespMsgType, cgr)
	gm.ToNetwork <- resp
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
		Character: &messages.Character{
			Name: netmsg.Name,
			ID:   0,
		},
	}
	var acct *Account
	// 1. Validate user is logged in as account specified.
	for _, uact := range gm.Users[msg.client.ID].Accounts {
		if uact.ID == netmsg.AccountID {
			acct = uact
		}
	}
	gm.CharID++
	if acct != nil {
		char := &Character{
			ID:    gm.CharID,
			Name:  netmsg.Name,
			Items: []*Item{},
		}
		acct.Characters = append(acct.Characters, char)
		gm.Characters[gm.CharID] = char
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

	if _, ok := gm.AcctByName[netmsg.Name]; !ok {
		gm.AccountID++
		ac.AccountID = gm.AccountID
		gm.Accounts[ac.AccountID] = &Account{
			ID:         ac.AccountID,
			Name:       netmsg.Name,
			Password:   netmsg.Password,
			Characters: []*Character{},
		}
		gm.AcctByName[netmsg.Name] = gm.Accounts[ac.AccountID]
		gm.Users[msg.client.ID].Accounts = append(gm.Users[msg.client.ID].Accounts, gm.Accounts[ac.AccountID])
	}

	resp := NewOutgoingMsg(msg.client, messages.CreateAcctRespMsgType, ac)
	gm.ToNetwork <- resp
}

func (gm *GameManager) loginUser(msg GameMessage) {
	tmsg := msg.net.(*messages.Login)
	lr := messages.LoginResp{
		Success: 0,
		Name:    tmsg.Name,
	}

	if acct, ok := gm.AcctByName[tmsg.Name]; ok {
		if acct.Password == tmsg.Password {
			log.Printf("Logging in account: %s", tmsg.Name)
			lr.AccountID = acct.ID
			lr.Characters = make([]*messages.Character, len(acct.Characters))
			for idx, ch := range acct.Characters {
				lr.Characters[idx] = &messages.Character{
					Name: ch.Name,
					ID:   ch.ID,
				}
			}
			gm.Users[msg.client.ID].Accounts = append(gm.Users[msg.client.ID].Accounts, acct)
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

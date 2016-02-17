package server

import (
	"fmt"
	"log"
	"math"

	"github.com/lologarithm/survival/server/messages"
)

// GameStatus type used for setting status of a game
type GameStatus byte

// Game statuses
const (
	UnknownStatus GameStatus = 0
	RunningStatus GameStatus = iota
)

// GameManager manages all connected users and games.
type GameManager struct {
	// Player data
	Users      []*User
	Games      map[uint32]*GameSession
	NextGameID uint32 // TODO: this shouldn't just be a number..

	FromGames   chan GameMessage // Manager reads this only, all games created write only
	FromNetwork <-chan GameMessage
	ToNetwork   chan<- OutgoingMessage
	Exit        chan int

	// Temp junk to make this crap work
	Accounts   []*Account
	CharID     uint32
	AccountID  uint32
	AcctByName map[string]*Account
}

// NewGameManager is the constructor for the main game manager.
// This should only be called once on a single server.
func NewGameManager(exit chan int, fromNetwork chan GameMessage, toNetwork chan OutgoingMessage) *GameManager {
	gm := &GameManager{
		Users:       make([]*User, math.MaxUint16),
		Games:       map[uint32]*GameSession{},
		FromGames:   make(chan GameMessage, 100),
		FromNetwork: fromNetwork,
		ToNetwork:   toNetwork,
		Exit:        exit,
		Accounts:    make([]*Account, math.MaxUint16),
		AcctByName:  map[string]*Account{},
	}
	return gm
}

// Run launches the game manager.
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
	case messages.DisconnectedMsgType:
		gm.handleDisconnect(msg)
	case messages.ConnectedMsgType:
		gm.handleConnection(msg)
	case messages.CreateAcctMsgType:
		gm.createAccount(msg)
	case messages.LoginMsgType:
		gm.loginUser(msg)
	case messages.JoinGameMsgType:
		// TODO: make this work
	case messages.CreateGameMsgType:
		gm.createGame(msg)
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
	case messages.EndGameMsgType:
		tmsg := msg.net.(*messages.EndGame)
		gameid := tmsg.GameID
		if msg.client != nil {
			gameid = gm.Users[msg.client.ID].GameID
		}
		fmt.Printf("Ended game: %d", gameid)
		gm.Games[gameid] = nil

	default:
		// These messages probably go to a game?
		// TODO: Probably have a direct conn to a game from the *Client
	}
}

func (gm *GameManager) createGame(msg GameMessage) {
	cgm := msg.net.(*messages.CreateGame)
	gm.NextGameID++

	netchan := make(chan GameMessage, 100)
	g := NewGame(cgm.Name, gm.FromGames, netchan, gm.ToNetwork)
	g.ID = gm.NextGameID
	g.SpawnChunk(0, 0)
	go g.Run()

	for _, a := range gm.Users[msg.client.ID].Accounts {
		g.FromGameManager <- AddPlayer{
			Entity: &Entity{
				Name: a.Character.Name,
			},
			Client: msg.client,
		}
	}

	gm.Games[gm.NextGameID] = g
	cgr := &messages.CreateGameResp{
		Name: cgm.Name,
		Game: &messages.GameConnected{
			ID:       gm.NextGameID,
			Seed:     g.Seed,
			Entities: g.World.EntitiesMsg(),
		},
	}
	gm.Users[msg.client.ID].GameID = msg.client.ID
	msg.client.FromGameManager <- ConnectedGame{
		ToGame: netchan,
		ID:     msg.client.ID,
	}
	resp := NewOutgoingMsg(msg.client, messages.CreateGameRespMsgType, cgr)
	gm.ToNetwork <- resp
}

func (gm *GameManager) handleConnection(msg GameMessage) {
	// First make sure this is a new connection.
	if gm.Users[msg.client.ID] == nil {
		gm.Users[msg.client.ID] = &User{
			Client:   msg.client,
			Accounts: []*Account{},
		}
	}
}

func (gm *GameManager) handleDisconnect(msg GameMessage) {
	// message active game that player disconnected.
	gameid := gm.Users[msg.client.ID].GameID
	if gm.Games[gameid] != nil {
		gm.Games[gameid].FromGameManager <- RemovePlayer{Client: msg.client}
	}
	// Then clear out the user.
	gm.Users[msg.client.ID] = nil
}

func (gm *GameManager) createAccount(msg GameMessage) {
	netmsg := msg.net.(*messages.CreateAcct)
	ac := &messages.CreateAcctResp{
		AccountID: 0,
		Name:      netmsg.Name,
	}

	if _, ok := gm.AcctByName[netmsg.Name]; !ok {
		gm.AccountID++
		gm.CharID++
		gm.Accounts[gm.AccountID] = &Account{
			ID:       gm.AccountID,
			Name:     netmsg.Name,
			Password: netmsg.Password,
			Character: &Character{
				Name:           netmsg.CharName,
				EquippedItems:  make([]*Item, 6),
				InventoryItems: []*Item{},
			},
		}

		ac.AccountID = gm.AccountID
		ac.Character = &messages.Character{
			Name: netmsg.CharName,
			ID:   gm.Accounts[gm.AccountID].Character.ID,
		}

		gm.AcctByName[netmsg.Name] = gm.Accounts[gm.AccountID]
		gm.Users[msg.client.ID].Accounts = append(gm.Users[msg.client.ID].Accounts, gm.Accounts[gm.AccountID])
	}

	resp := NewOutgoingMsg(msg.client, messages.CreateAcctRespMsgType, ac)
	gm.ToNetwork <- resp
}

func (gm *GameManager) loginUser(msg GameMessage) {
	tmsg := msg.net.(*messages.Login)
	lr := messages.LoginResp{
		Success:   0,
		Name:      tmsg.Name,
		Character: &messages.Character{},
	}
	if acct, ok := gm.AcctByName[tmsg.Name]; ok {
		if acct.Password == tmsg.Password {
			log.Printf("Logging in account: %s", tmsg.Name)
			lr.AccountID = acct.ID
			lr.Character = &messages.Character{
				Name: acct.Character.Name,
				ID:   acct.Character.ID,
			}
			gm.Users[msg.client.ID].Accounts = append(gm.Users[msg.client.ID].Accounts, acct)
		}
	}
	resp := NewOutgoingMsg(msg.client, messages.LoginRespMsgType, &lr)
	gm.ToNetwork <- resp
}

// ProcessGameMsg is used to process messages from an individual game to the main server controller.
func (gm *GameManager) ProcessGameMsg(msg GameMessage) {
	switch msg.mtype {
	}
}

// NewOutgoingMsg creates a new message that can be sent to a specific client.
func NewOutgoingMsg(dest *Client, tp messages.MessageType, msg messages.Net) OutgoingMessage {
	frame := messages.Frame{
		MsgType:       tp,
		Seq:           1,
		ContentLength: uint16(msg.Len()),
	}
	resp := OutgoingMessage{
		dest: dest,
		msg: messages.Packet{
			Frame:  frame,
			NetMsg: msg,
		},
	}
	return resp
}

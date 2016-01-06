package server

import (
	"fmt"
	"time"
)

type GameStatus byte

const (
	Unknown GameStatus = 1
	Running GameStatus = iota
)

// GameManager manages all connected users and games.
type GameManager struct {
	// Player data
	Clients     map[uint32]*Client
	Games       map[uint32]*Game
	NextGameId  uint32
	FromNetwork chan GameMessage
	FromGames   chan GameMessage
	ToNetwork   chan NetMessage
	Exit        chan int
}

func NewGameManager(exit chan int, fromNetwork chan GameMessage, toNetwork chan NetMessage) *GameManager {
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
	switch msg.(type) {
	}
}

func (gm *GameManager) ProcessGameMsg(msg GameMessage) {
	switch msg.(type) {
	}
}

// Run starts the game!
func (gameManager *Game) Run() {
	wait_for_timeout := true
	for {
		timeout := time.Millisecond * 50
		wait_for_timeout = true
		for wait_for_timeout {
			select {
			case <-time.After(timeout):
				wait_for_timeout = false
				break
			case msg := <-gameManager.FromNetwork:
				fmt.Printf("GameManager: Received message: %T\n", msg)
				switch msg.(type) {
				default:
					fmt.Println("GameManager.go:RunGame(): UNKNOWN MESSAGE TYPE: %T", msg)
				}
			case <-gameManager.Exit:
				fmt.Println("EXITING Game Manager")
				return
			}
		}
		fmt.Printf("Sending client update!\n")
	}
}

// Setup sets up the game!
func (g *Game) Setup() {
}

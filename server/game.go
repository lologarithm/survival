package server

import (
	"fmt"
	"time"

	"github.com/lologarithm/survival/server/forestGen"
	"github.com/lologarithm/survival/server/messages"
)

// Game represents a single game
type Game struct {
	Name string
	// Player data
	Clients         map[uint32]*Client
	IntoGameManager chan GameMessage
	FromNetwork     chan GameMessage // Messages from players.
	Exit            chan int
	Status          GameStatus
}

type GameState struct {
	Map *forestGen.Map
}

// TODO: Structure tiles?

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
					fmt.Printf("GameManager.go:RunGame(): UNKNOWN MESSAGE TYPE: %T\n", msg)
				}
			case <-gm.Exit:
				fmt.Println("EXITING Game Manager")
				return
			}
		}
		fmt.Printf("Sending client update!\n")
	}
}

// NewGame constructs a new game and starts it.
func NewGame(name string, toGameManager chan GameMessage) *Game {
	g := &Game{
		Name:            name,
		IntoGameManager: toGameManager,
		FromNetwork:     make(chan GameMessage, 100),
	}
	g.Run()
	return g
}

type GameMessage struct {
	net    messages.Net
	client *Client
	mtype  messages.MessageType
}

// TODO: Is this needed?
type InternalMessage struct {
	ToGame chan GameMessage
}

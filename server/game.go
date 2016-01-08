package server

import (
	"fmt"
	"time"

	"github.com/lologarithm/survival/server/messages"
)

// Game represents a single game
type Game struct {
	// Player data
	Clients           map[uint32]*Client
	IntoServerManager chan GameMessage
	FromNetwork       chan GameMessage // Messages from players.
	Exit              chan int
	Status            GameStatus
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

type GameMessage struct {
	net    messages.Net
	client *Client
	mtype  messages.MessageType
}

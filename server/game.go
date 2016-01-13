package server

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"

	xxhash "github.com/OneOfOne/xxhash/native"
	"github.com/lologarithm/survival/server/messages"
)

// Game represents a single game
type Game struct {
	Name string
	Seed uint64 // Select a seed when starting the game!

	// Player data
	Clients map[uint32]*Client

	// Game can only write to this channel, not read.
	IntoGameManager chan<- GameMessage
	// FromNetwork is read here and written elsewhere.
	FromNetwork <-chan GameMessage // Messages from players.

	Exit   chan int
	Status GameStatus

	// Private
	World GameWorld
}

type GameWorld struct {
	Entities []*Entity
	Chunks   map[uint32]map[uint32]bool // list of chunks that have been already created.
}

func (gw *GameWorld) EntitiesMsg() []*messages.Entity {
	es := make([]*messages.Entity, len(gw.Entities))
	for idx, e := range gw.Entities {
		es[idx] = e.toMsg()
	}

	return es
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

func (gm *Game) SpawnChunk(x, y uint32) {
	h := xxhash.New64()
	tb := make([]byte, 8)

	binary.LittleEndian.PutUint64(tb[:8], gm.Seed)
	h.Write(tb[:8])

	binary.LittleEndian.PutUint32(tb[:4], x)
	h.Write(tb[:4])

	binary.LittleEndian.PutUint32(tb[:4], y)
	h.Write(tb[:4])

	binary.LittleEndian.PutUint32(tb[:4], 1)
	h.Write(tb[:4])

	chunkSeed := h.Sum64()
	numRocks := chunkSeed >> 60
	// numBush := (chunkSeed << 4) >> 60
	numTrees := (chunkSeed << 8) >> 56

	for i := 0; i < int(numRocks); i++ {
		oh := xxhash.New64() // (worldseed, chunkX, chunkY, 10, rock#)
		binary.LittleEndian.PutUint64(tb[:8], gm.Seed)
		oh.Write(tb[:8])

		binary.LittleEndian.PutUint32(tb[:4], x)
		oh.Write(tb[:4])

		binary.LittleEndian.PutUint32(tb[:4], y)
		oh.Write(tb[:4])

		binary.LittleEndian.PutUint32(tb[:4], 10)
		oh.Write(tb[:4])

		binary.LittleEndian.PutUint32(tb[:4], uint32(i))
		oh.Write(tb[:4])

		oSeed := oh.Sum64()
		// floor(bits 0:8 / 2.57) = rock X position relative to chunk
		ox := uint32(float64(oSeed>>56) / 2.57)
		// floor(bits 8:16 / 2.57) = rock Y position relative to chunk
		oy := uint32(float64((oSeed<<8)>>56) / 2.57)

		te := &Entity{
			X:      ox,
			Y:      oy,
			Seed:   oSeed,
			Height: 2,
			Width:  2,
			EType:  0,
		}
		// Check if existing rock overlaps this rock, if so, make old rock bigger!
		intersected := false
		for _, t := range gm.World.Entities {
			if t.Intersects(te) {
				if t.EType == te.EType {
					t.Height++
					t.Width++
				}
				intersected = true
				break
			}
		}
		if !intersected {
			gm.World.Entities = append(gm.World.Entities, te)
		}
	}

	for i := 0; i < int(numTrees); i++ {
		oh := xxhash.New64() // (worldseed, chunkX, chunkY, 10, rock#)
		binary.LittleEndian.PutUint64(tb[:8], gm.Seed)
		oh.Write(tb[:8])

		binary.LittleEndian.PutUint32(tb[:4], x)
		oh.Write(tb[:4])

		binary.LittleEndian.PutUint32(tb[:4], y)
		oh.Write(tb[:4])

		binary.LittleEndian.PutUint32(tb[:4], 12)
		oh.Write(tb[:4])

		binary.LittleEndian.PutUint32(tb[:4], uint32(i))
		oh.Write(tb[:4])

		oSeed := oh.Sum64()
		// floor(bits 0:8 / 2.57) = tree X position relative to chunk
		ox := uint32(float64(oSeed>>56) / 2.57)
		// floor(bits 8:16 / 2.57) = tree Y position relative to chunk
		oy := uint32(float64((oSeed<<8)>>56) / 2.57)

		te := &Entity{
			X:      ox,
			Y:      oy,
			Seed:   oSeed,
			Height: 3,
			Width:  3,
			EType:  2,
		}
		// Check if existing tree overlaps this tree, if so, make old tree bigger!
		intersected := false
		for _, t := range gm.World.Entities {
			if t.Intersects(te) {
				if t.EType == te.EType {
					t.Height += 2
					t.Width += 2
				}
				intersected = true
				break
			}
		}
		if !intersected {
			gm.World.Entities = append(gm.World.Entities, te)
		}
	}

	if gm.World.Chunks == nil {
		gm.World.Chunks = map[uint32]map[uint32]bool{}
	}
	if gm.World.Chunks[x] == nil {
		gm.World.Chunks[x] = map[uint32]bool{}
	}
	gm.World.Chunks[x][y] = true
}

// NewGame constructs a new game and starts it.
func NewGame(name string, toGameManager chan<- GameMessage, fromNetwork <-chan GameMessage) *Game {
	seed := uint64(rand.Uint32())
	seed = seed << 32
	seed += uint64(rand.Uint32())
	g := &Game{
		Name:            name,
		IntoGameManager: toGameManager,
		FromNetwork:     fromNetwork,
		Seed:            seed,
		World:           GameWorld{},
	}
	// go g.Run()
	return g
}

type Entity struct {
	ID     uint32
	EType  uint16
	Seed   uint64
	X      uint32
	Y      uint32
	Height uint32
	Width  uint32
}

func (e *Entity) toMsg() *messages.Entity {
	o := &messages.Entity{
		ID:     e.ID,
		X:      e.X,
		Y:      e.Y,
		Height: e.Height,
		Width:  e.Width,
		EType:  e.EType,
		Seed:   e.Seed,
	}

	return o
}

func (e *Entity) Intersects(o *Entity) bool {
	if e.X > o.X+o.Width || e.X+e.Width < o.X || e.Y > o.Y+o.Height || e.Y+e.Height < o.Y {
		return false
	}
	return true
}

type GameMessage struct {
	net    messages.Net
	client *Client
	mtype  messages.MessageType
}

// TODO: Is this needed?
type InternalMessage struct {
	ToGame chan<- GameMessage
}

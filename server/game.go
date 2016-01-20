package server

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"

	xxhash "github.com/OneOfOne/xxhash/native"
	"github.com/lologarithm/survival/physics"
	"github.com/lologarithm/survival/server/messages"
)

const ChunkSize int32 = 2000

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
	World *GameWorld
}

// GameWorld represents all the data in the world.
// Physical entities and the physics simulation.
type GameWorld struct {
	Entities []*Entity
	Chunks   map[uint32]map[uint32]bool // list of chunks that have been already created.
	Space    *physics.SimulatedSpace
}

// EntitiesMsg converts all entities in the world to a network message.
func (gw *GameWorld) EntitiesMsg() []*messages.Entity {
	es := make([]*messages.Entity, len(gw.Entities))
	for idx, e := range gw.Entities {
		es[idx] = e.toMsg()
	}

	return es
}

// Run starts the game!
func (g *Game) Run() {
	waiting := true
	for {
		timeout := time.Millisecond * 50
		waiting = true
		for waiting {
			select {
			case <-time.After(timeout):
				waiting = false
				break
			case msg := <-g.FromNetwork:
				fmt.Printf("GameManager: Received message: %T\n", msg)
				switch msg.mtype {
				case messages.JoinGameMsgType:
					tmsg := msg.net.(*messages.JoinGame)

					player := &Entity{
						ID:     tmsg.CharID,
						EType:  3, // TODO: make constnats
						Seed:   0, // players dont need a seed?
						Height: 6,
						Width:  6,
						Body:   physics.NewRigidBody(tmsg.CharID, physics.Vect2{1000, 1000}, physics.Vect2{}, 0, 100),
					}
					g.World.Space.AddEntity(player.Body, false)
				case messages.EntityMoveMsgType:
				default:
					fmt.Printf("GameManager.go:RunGame(): UNKNOWN MESSAGE TYPE: %T\n", msg)
				}
			case <-g.Exit:
				fmt.Println("EXITING Game Manager")
				return
			}
		}
		g.World.Space.Tick(true)
		// TODO: send updates from the tick?
		fmt.Printf("Sending client update!\n")
	}
}

func (g *Game) MoveEntity() {

}

// SpawnChunk creates all the entities for a chunk at the given x/y
func (g *Game) SpawnChunk(x, y uint32) {
	h := xxhash.New64()
	tb := make([]byte, 8)

	binary.LittleEndian.PutUint64(tb[:8], g.Seed)
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
		binary.LittleEndian.PutUint64(tb[:8], g.Seed)
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
		ox := int32(oSeed>>48) / 33
		oy := int32((oSeed<<16)>>48) / 33

		te := &Entity{
			Body: &physics.RigidBody{
				Position: physics.Vect2{
					X: ox,
					Y: oy,
				},
			},
			Seed:   oSeed,
			Height: 5,
			Width:  5,
			EType:  0,
		}
		// Check if existing rock overlaps this rock, if so, make old rock bigger!
		intersected := false
		for _, t := range g.World.Entities {
			if t.Intersects(te) {
				if t.EType == te.EType {
					t.Height += 3
					t.Width += 3
				}
				intersected = true
				break
			}
		}
		if !intersected {
			g.World.Entities = append(g.World.Entities, te)
		}
	}

	for i := 0; i < int(numTrees); i++ {
		oh := xxhash.New64() // (worldseed, chunkX, chunkY, 12, tree#)
		binary.LittleEndian.PutUint64(tb[:8], g.Seed)
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
		ox := int32(oSeed>>48) / 33
		oy := int32((oSeed<<16)>>48) / 33
		te := &Entity{
			Body: &physics.RigidBody{
				Position: physics.Vect2{
					X: ox,
					Y: oy,
				},
			},
			Seed:   oSeed,
			Height: 20,
			Width:  20,
			EType:  2,
		}

		// Check if existing tree overlaps this tree, if so, make old tree bigger!
		intersected := false
		for _, t := range g.World.Entities {
			if t.Intersects(te) {
				if t.EType == te.EType {
					t.Height += 5
					t.Width += 5
				}
				intersected = true
				break
			}
		}
		if !intersected {
			g.World.Entities = append(g.World.Entities, te)
		}
	}

	if g.World.Chunks == nil {
		g.World.Chunks = map[uint32]map[uint32]bool{}
	}
	if g.World.Chunks[x] == nil {
		g.World.Chunks[x] = map[uint32]bool{}
	}
	g.World.Chunks[x][y] = true
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
		World:           &GameWorld{},
	}
	// go g.Run()
	return g
}

// Entity represents a single object in the game.
type Entity struct {
	ID     uint32
	EType  uint16
	Seed   uint64
	Height int32
	Width  int32
	Body   *physics.RigidBody
}

func (e *Entity) toMsg() *messages.Entity {
	o := &messages.Entity{
		ID:     e.ID,
		X:      e.Body.Position.X,
		Y:      e.Body.Position.Y,
		Height: e.Height,
		Width:  e.Width,
		EType:  e.EType,
		Seed:   e.Seed,
	}

	return o
}

// Intersects calculates if two entities overlap -- used currently for chunk generation.
func (e *Entity) Intersects(o *Entity) bool {
	if e.Body.Position.X > o.Body.Position.X+o.Width || e.Body.Position.X+e.Width < o.Body.Position.X || e.Body.Position.Y > o.Body.Position.Y+o.Height || e.Body.Position.Y+e.Height < o.Body.Position.Y {
		return false
	}
	return true
}

// GameMessage is a message from a client to a game.
type GameMessage struct {
	net    messages.Net
	client *Client
	mtype  messages.MessageType
}

// InternalMessage TODO: Is this needed?
type InternalMessage struct {
	ToGame chan<- GameMessage
}

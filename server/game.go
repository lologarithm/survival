package server

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	xxhash "github.com/OneOfOne/xxhash/native"
	"github.com/lologarithm/survival/physics"
	"github.com/lologarithm/survival/server/messages"
)

const ChunkSize int32 = 10000

// Game represents a single game
type Game struct {
	Name string
	Seed uint64 // Select a seed when starting the game!

	// map character ID to client
	Clients map[uint32]*Client

	// Game can only write to this channel, not read.
	IntoGameManager chan<- GameMessage
	FromGameManager chan InternalMessage
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
				case messages.MovePlayerMsgType:
					tmsg := msg.net.(*messages.MovePlayer)
					log.Printf("Moving character: %v", tmsg)
					// TODO: go back in time and apply at tick!
					for _, ent := range g.World.Entities {
						if ent.ID == tmsg.EntityID {
							ent.Body.Angle = float64(tmsg.Direction%365) * math.Pi / 180.0
							// TODO: Replace hardcoded 50 with 'speed' setting of the character.
							ent.Body.Velocity = physics.Vect2{X: int32(math.Cos(ent.Body.Angle) * 50), Y: int32(math.Sin(ent.Body.Angle)) * 50}
							break
						}
					}
				default:
					fmt.Printf("GameManager.go:RunGame(): UNKNOWN MESSAGE TYPE: %T\n", msg)
				}
			case imsg := <-g.FromGameManager:
				switch timsg := imsg.(type) {
				case AddPlayer:
					player := &Entity{
						ID:    timsg.Entity.ID,
						EType: 3, // TODO: make constants
						Seed:  timsg.Entity.Seed,
						Body:  physics.NewRigidBody(timsg.Entity.ID, 22, 46, physics.Vect2{X: 5000, Y: 5000}, physics.Vect2{}, 0, 100),
					}
					g.World.Space.AddEntity(player.Body, false)
					g.Clients[timsg.Entity.ID] = timsg.Client
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
		ox := int32(oSeed>>48) / (math.MaxUint16/ChunkSize + 1)
		oy := int32((oSeed<<16)>>48) / (math.MaxUint16/ChunkSize + 1)

		te := &Entity{
			Body: &physics.RigidBody{
				Position: physics.Vect2{
					X: ox,
					Y: oy,
				},
				Height: 10,
				Width:  10,
			},
			Seed:  oSeed,
			EType: 0,
		}
		// Check if existing rock overlaps this rock, if so, make old rock bigger!
		intersected := false
		for _, t := range g.World.Entities {
			if t.Intersects(te) {
				if t.EType == te.EType {
					t.Body.Height += 7
					t.Body.Width += 7
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
		ox := int32(oSeed>>48) / (math.MaxUint16/ChunkSize + 1)
		oy := int32((oSeed<<16)>>48) / (math.MaxUint16/ChunkSize + 1)
		te := &Entity{
			Body: &physics.RigidBody{
				Position: physics.Vect2{
					X: ox,
					Y: oy,
				},
				Height: 100,
				Width:  100,
			},
			Seed:  oSeed,
			EType: 2,
		}

		// Check if existing tree overlaps this tree, if so, make old tree bigger!
		intersected := false
		for _, t := range g.World.Entities {
			if t.Intersects(te) {
				if t.EType == te.EType {
					t.Body.Height += 75
					t.Body.Width += 75
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
		FromGameManager: make(chan InternalMessage, 100),
		FromNetwork:     fromNetwork,
		Seed:            seed,
		World:           &GameWorld{},
	}
	// go g.Run()
	return g
}

// Entity represents a single object in the game.
type Entity struct {
	ID    uint32
	EType uint16
	Seed  uint64
	Body  *physics.RigidBody
}

func (e *Entity) toMsg() *messages.Entity {
	o := &messages.Entity{
		ID:     e.ID,
		X:      e.Body.Position.X,
		Y:      e.Body.Position.Y,
		Height: e.Body.Height,
		Width:  e.Body.Width,
		Angle:  int16(e.Body.Angle * 180 / math.Pi),
		EType:  e.EType,
		Seed:   e.Seed,
	}

	return o
}

// Intersects calculates if two entities overlap -- used currently for chunk generation.
func (e *Entity) Intersects(o *Entity) bool {
	if e.Body.Position.X > o.Body.Position.X+o.Body.Width || e.Body.Position.X+e.Body.Width < o.Body.Position.X || e.Body.Position.Y > o.Body.Position.Y+o.Body.Height || e.Body.Position.Y+e.Body.Height < o.Body.Position.Y {
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

// InternalMessage is for messages between components that never leaves the server.
type InternalMessage interface {
}

type ConnectedGame struct {
	ToGame chan<- GameMessage
}

type AddPlayer struct {
	Entity *Entity
	Client *Client
}

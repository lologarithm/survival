package server

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"time"

	xxhash "github.com/OneOfOne/xxhash/native"
	"github.com/lologarithm/survival/physics"
	"github.com/lologarithm/survival/server/messages"
)

const ChunkSize int32 = 10000

const (
	UnknownEType uint16 = iota
	RockEType
	BushEType
	TreeEType
	CreatureEType
	ProjectileEType
)

// Game represents a single game
type Game struct {
	ID   uint32
	Name string
	Seed uint64 // Select a seed when starting the game!

	// map character ID to client
	Clients map[uint32]*User

	IntoGameManager chan<- GameMessage     // Game can only write to this channel, not read.
	FromGameManager chan InternalMessage   // Messages from the game Manager.
	FromNetwork     <-chan GameMessage     // FromNetwork is read only here, messages from players.
	ToNetwork       chan<- OutgoingMessage // Messages to players!

	Exit   chan int
	Status GameStatus

	// Private
	World          *GameWorld    // Current world state
	prevWorlds     []*GameWorld  // Last X seconds of game states
	commandHistory []interface{} // Last X seconds of commands
}

// GameWorld represents all the data in the world.
// Physical entities and the physics simulation.
type GameWorld struct {
	Entities map[uint32]*Entity
	Chunks   map[uint32]map[uint32]bool // list of chunks that have been already created.
	Space    *physics.SimulatedSpace
}

// Clone returns a deep copy of the game world at this time.
func (gw *GameWorld) Clone() *GameWorld {
	// TODO: make this work.
	return &GameWorld{}
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
		timeout := time.Millisecond * 33
		waiting = true
		for waiting {
			select {
			case <-time.After(timeout):
				waiting = false
				break
			case msg := <-g.FromNetwork:
				switch msg.mtype {
				case messages.MovePlayerMsgType:
					g.MoveEntity(msg.client, msg.net.(*messages.MovePlayer))
				default:
					fmt.Printf("game.go:Run(): UNKNOWN MESSAGE TYPE: %T\n", msg)
				}
			case imsg := <-g.FromGameManager:
				switch timsg := imsg.(type) {
				case AddPlayer:
					newid := uint32(len(g.World.Entities))
					player := &Entity{
						ID:    newid,
						Name:  timsg.Entity.Name,
						EType: CreatureEType,
						Seed:  timsg.Entity.Seed,
						Body:  physics.NewRigidBody(newid, 22, 46, physics.Vect2{X: 5000, Y: 5000}, physics.Vect2{}, 0, 100),
					}
					g.World.Space.AddEntity(player.Body, false)
					g.World.Entities[newid] = player
					g.Clients[timsg.Client.ID] = &User{
						Client: timsg.Client,
						Accounts: []*Account{
							{
								Character: &Character{
									ID: newid,
								},
							},
						},
					}
				case RemovePlayer:
					id := g.Clients[timsg.Client.ID].Accounts[0].Character.ID
					ent := g.World.Entities[id]
					if ent == nil {
						return
					}
					// TODO: remove player from game after timeout?
					g.World.Space.RemoveEntity(ent.Body, false)
					delete(g.Clients, timsg.Client.ID)
					if len(g.Clients) == 0 {
						fmt.Printf("All clients disconnected, closing game %d.", g.ID)
						g.IntoGameManager <- GameMessage{
							net:   &messages.EndGame{},
							mtype: messages.EndGameMsgType,
						}
						return
					}
				}
			case <-g.Exit:
				fmt.Print("EXITING: Run in Game.go\n")
				return
			}
		}
		collisions := g.World.Space.Tick(true)
		for _, col := range collisions {
			var ent *Entity
			for _, e := range g.World.Entities {
				if e.ID == col.Body.ID {
					ent = e
					break
				}
			}
			if ent.EType == ProjectileEType {
				// TODO: remove the projectile
				// TODO: resolve hit to target.
			}
		}
		if g.World.Space.TickID%20 == 0 {
			g.SendMasterFrame()
		}
	}
}

func (g *Game) SendMasterFrame() {
	mf := &messages.GameMasterFrame{
		ID:       g.ID,
		Entities: g.World.EntitiesMsg(),
	}

	frame := messages.Frame{
		MsgType:       messages.GameMasterFrameMsgType,
		Seq:           1,
		ContentLength: uint16(mf.Len()),
	}
	msg := OutgoingMessage{
		msg: messages.Packet{
			Frame:  frame,
			NetMsg: mf,
		},
	}

	for _, c := range g.Clients {
		msg.dest = c.Client
		g.ToNetwork <- msg
	}
}

func (g *Game) MoveEntity(c *Client, tmsg *messages.MovePlayer) {
	// TODO: go back in time and apply at tick!
	id := g.Clients[c.ID].Accounts[0].Character.ID
	ent := g.World.Entities[id]
	if ent == nil {
		return
	}
	dirVect := physics.Vect2{
		X: int32(tmsg.X),
		Y: int32(tmsg.Y),
	}
	ent.Body.Angle = physics.Angle(dirVect, physics.Vect2{X: 0, Y: 1})
	// TODO: Replace hardcoded 50 with 'speed' setting of the character.
	ent.Body.Velocity = physics.Normalize(dirVect, 50)
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
			EType: RockEType,
			ID:    uint32(len(g.World.Entities)),
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
			g.World.Entities[te.ID] = te
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
			EType: TreeEType,
			ID:    uint32(len(g.World.Entities)),
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
			g.World.Entities[te.ID] = te
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
func NewGame(name string, toGameManager chan<- GameMessage, fromNetwork <-chan GameMessage, toNetwork chan<- OutgoingMessage) *Game {
	seed := uint64(rand.Uint32())
	seed = seed << 32
	seed += uint64(rand.Uint32())
	g := &Game{
		Name:            name,
		IntoGameManager: toGameManager,
		FromGameManager: make(chan InternalMessage, 100),
		FromNetwork:     fromNetwork,
		ToNetwork:       toNetwork,
		Seed:            seed,
		World: &GameWorld{
			Space:    physics.NewSimulatedSpace(),
			Entities: map[uint32]*Entity{},
			Chunks:   map[uint32]map[uint32]bool{}, // list of chunks that have been already created.
		},
		Exit:    make(chan int, 1),
		Clients: make(map[uint32]*User, 16),
	}
	return g
}

// Entity represents a single object in the game.
type Entity struct {
	ID    uint32
	Name  string
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

type RemovePlayer struct {
	Client *Client
}

type AddPlayer struct {
	Entity *Entity
	Client *Client
}

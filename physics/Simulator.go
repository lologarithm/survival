package physics

import (
	"math"
	"time"

	"github.com/lologarithm/survival/physics/quadtree"
)

const (
	SimUpdatesPerSecond = 50.0
	SimUpdateSleep      = 1000.0 / SimUpdatesPerSecond
	FullCircle          = math.Pi * 2

	AddEntity       = byte(1)
	UpdateForces    = byte(3)
	UpdatePosition  = byte(4)
	UpdateCollision = byte(5)
)

// Simulator design:
//  1. Needs to be able to represent position of each thing in time correctly.
//  2. Probably want a simplified 2d physics simulator running to allow for things with velocity?
//  3. Each tick should have an ID and should be rewindable (so we can insert updates in the past)
//  4.

func NewSimulatedSpace() *SimulatedSpace {
	world := quadtree.NewBoundingBox(-10000.0, 10000.0, -10000.0, 10000.0)
	return &SimulatedSpace{
		tree:     quadtree.NewQuadTree(world),
		Entities: make([]*RigidBody, 1000),
		Fixed:    make([]*RigidBody, 1000),
	}
}

type SimulatedSpace struct {
	tree     quadtree.QuadTree
	Entities []*RigidBody // Anything that can collide in the playspace
	Fixed    []*RigidBody // Anything that can collide but is fixed in place.
	entidx   int
	fixidx   int

	lastUpdate time.Time
	TickID     uint32
}

func (ss *SimulatedSpace) AddEntity(body *RigidBody, fixed bool) {
	if fixed {
		ss.Fixed[ss.fixidx] = body
		ss.fixidx++
		return
	}
	ss.Entities[ss.entidx] = body
	ss.entidx++
	ss.tree.Add(body)
}

func (ss *SimulatedSpace) RemoveEntity(body *RigidBody, fixed bool) {
	if fixed {
		for cidx, f := range ss.Fixed {
			if f.ID == body.ID {
				ss.Fixed[cidx] = nil
				break
			}
		}
	} else {
		for cidx, f := range ss.Entities {
			if f.ID == body.ID {
				ss.Entities[cidx] = nil
				break
			}
		}
	}
	ss.tree.Remove(body)
}

func (ss *SimulatedSpace) Tick(sendUpdate bool) []PhysicsEntityUpdate {
	ss.TickID++
	ss.lastUpdate = time.Now()
	var changeList []PhysicsEntityUpdate
	if sendUpdate {
		changeList = make([]PhysicsEntityUpdate, len(ss.Entities)*2)
	}
	cidx := 0
	changed := false
	for _, rigid := range ss.Entities {
		if rigid == nil {
			continue
		}
		changed = false
		rigid.Velocity = AddVect2(rigid.Velocity, MultVect2(rigid.Force, rigid.InvMass/SimUpdatesPerSecond))
		rigid.AngularVelocity += (rigid.Torque * float64(rigid.InvInertia)) / SimUpdatesPerSecond

		if rigid.Velocity.X != 0.0 {
			rigid.Position.X += rigid.Velocity.X / SimUpdatesPerSecond
			changed = true
		}
		if rigid.Velocity.Y != 0.0 {
			rigid.Position.Y += rigid.Velocity.Y / SimUpdatesPerSecond
			changed = true
		}
		if rigid.AngularVelocity != 0.0 {
			rigid.Angle += rigid.AngularVelocity / SimUpdatesPerSecond
			for rigid.Angle > FullCircle {
				rigid.Angle -= FullCircle
			}
			for rigid.Angle < -FullCircle {
				rigid.Angle += FullCircle
			}
			changed = true
		}

		if !changed {
			continue
		}
		// if changed && sendUpdate {
		// 	changeList[cidx].UpdateType = UpdatePosition
		// 	changeList[cidx].Body = *rigid
		// 	cidx++
		// 	if cidx == len(changeList) {
		// 		newlist := make([]PhysicsEntityUpdate, len(changeList)*2)
		// 		copy(newlist, changeList)
		// 		changeList = newlist
		// 	}
		// }

		collisions := ss.tree.Query(rigid.BoundingBox())

		for _, collbox := range collisions {
			other := collbox.(*RigidBody)
			if other == nil || other.ID == rigid.ID {
				continue
			}

			changeList[cidx].UpdateType = UpdateCollision
			changeList[cidx].Body = rigid
			changeList[cidx].Other = other

			cidx++
			if cidx == len(changeList) {
				newlist := make([]PhysicsEntityUpdate, len(changeList)*2)
				copy(newlist, changeList)
				changeList = newlist
			}
		}
	}

	return changeList[:cidx]
}

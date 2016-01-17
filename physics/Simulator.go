package physics

import (
	"math"
	"time"
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

type SimulatedSpace struct {
	Entities   map[uint32]*RigidBody // Anything that can collide in the playspace
	Fixed      map[uint32]*RigidBody // Anything that can collide but is fixed in place.
	lastUpdate time.Time
	TickID     uint32
}

func (ss *SimulatedSpace) AddEntity(body *RigidBody, fixed bool) {
	if fixed {
		ss.Fixed[body.ID] = body
		return
	}
	ss.Entities[body.ID] = body
}

func (ss *SimulatedSpace) Tick(sendUpdate bool) []PhysicsEntityUpdate {
	ss.TickID++
	ss.lastUpdate = time.Now()
	var changeList []PhysicsEntityUpdate
	if sendUpdate {
		changeList = make([]PhysicsEntityUpdate, len(ss.Entities))
	}
	cidx := 0
	changed := false
	for _, rigid := range ss.Entities {
		changed = false

		rigid.Velocity = rigid.Velocity.Add(MultVect2(rigid.Force, rigid.InvMass/SimUpdatesPerSecond))
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
		if changed && sendUpdate {
			changeList[cidx].UpdateType = UpdatePosition
			changeList[cidx].Body = *rigid
			cidx++
		}
	}
	// Check for collisions?

	return changeList
}

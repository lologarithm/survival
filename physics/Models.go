package physics

import "github.com/lologarithm/survival/physics/quadtree"

func CrossProduct(a Vect2, b Vect2) int32 {
	return a.X*b.Y - a.Y*b.X
}

func CrossScalar(v Vect2, s int32) Vect2 {
	return Vect2{v.Y * s, -s * v.X}
}

func CrossScalarFirst(s int32, v Vect2) Vect2 {
	return Vect2{v.Y * -s, s * v.X}
}

func MultVect2(a Vect2, s int32) Vect2 {
	return Vect2{a.X * s, a.Y * s}
}

type Vect2 struct {
	X, Y int32
}

func (v Vect2) Add(v2 Vect2) Vect2 {
	return Vect2{v.X + v2.X, v.Y + v2.Y}
}

// RigidBody is an object in the physics simulation
type RigidBody struct {
	ID uint32 // Unique ID for this rigidbody

	Position Vect2 // coords x,y of entity  (meters)
	Velocity Vect2 // speed in vector format (m/s)
	Force    Vect2 // Force to apply each tick.

	Angle           float64 // Current heading (radians)
	AngularVelocity float64 // speed of rotation around the Z axis (radians/sec)
	Torque          float64 // Torque to apply each tick

	Mass       int32 // mass of the object, (kg)
	InvMass    int32 // Inverted mass for physics calcs
	Inertia    int32 // Inertia of the ship
	InvInertia int32 // Inverted Inertia for physics calcs

	Height int32
	Width  int32
}

func (rb RigidBody) BoundingBox() quadtree.BoundingBox {
	return quadtree.NewBoundingBox(rb.Position.X, rb.Position.X+rb.Width, rb.Position.Y, rb.Position.Y+rb.Height)
}

func NewRigidBody(id uint32, pos Vect2, vel Vect2, angle float64, mass int32) *RigidBody {
	return &RigidBody{
		ID:       id,
		Position: pos,
		Velocity: vel,
		Angle:    angle,
		Mass:     mass,
		InvMass:  1 / mass,
	}
}

// PhysicsEntityUpdate message linked to an Entity.
type PhysicsEntityUpdate struct {
	UpdateType byte      // 2 == add, 3 == remove, 4 == physics update
	Body       RigidBody // Passed by value through channels
}

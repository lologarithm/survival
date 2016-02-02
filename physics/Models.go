package physics

import (
	"math"

	"github.com/lologarithm/survival/physics/quadtree"
)

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

func Angle(a Vect2, b Vect2) float64 {
	alpha := float64(a.X*a.X+a.Y*b.Y) / (a.Magnitude() * b.Magnitude())
	return math.Acos(alpha)
}

type Vect2 struct {
	X, Y int32
}

func (v Vect2) Add(v2 Vect2) Vect2 {
	return Vect2{v.X + v2.X, v.Y + v2.Y}
}

func (v Vect2) Magnitude() float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y))
}

// RigidBody is an object in the physics simulation
type RigidBody struct {
	ID uint32 // Unique ID for this rigidbody

	Position Vect2 // coords x,y of center of entity  (arbitrary units)
	Velocity Vect2 // speed in vector format (units/sec)
	Force    Vect2 // Force to apply each tick.

	Angle           float64 // Current heading (radians)
	AngularVelocity float64 // speed of rotation around the Z axis (radians/sec)
	Torque          float64 // Torque to apply each tick

	Mass       int32 // Mass of the object, (kg)
	InvMass    int32 // Inverted mass for physics calcs
	Inertia    int32 // Inertia of the object
	InvInertia int32 // Inverted Inertia for physics calcs

	Height int32
	Width  int32
}

func (rb RigidBody) BoxID() uint32 {
	return rb.ID
}

func (rb RigidBody) BoundingBox() quadtree.BoundingBox {
	hh := rb.Height / 2
	hw := rb.Width / 2
	if rb.Angle != 0 && rb.Angle != math.Pi && rb.Angle != math.Pi*2 {
		s := math.Sin(rb.Angle)
		c := math.Cos(rb.Angle)
		if s < 0 {
			s = -s
		}
		if c < 0 {
			c = -c
		}
		wn := int32(float64(rb.Height)*s) + int32(float64(rb.Width)*c) // width of AABB
		hwn := wn / 2
		hn := int32(float64(rb.Height)*c) + int32(float64(rb.Width)*s) // height of AABB
		hhn := hn / 2
		return quadtree.NewBoundingBox(rb.Position.X-hwn, rb.Position.X+hwn, rb.Position.Y-hhn, rb.Position.Y+hhn)
	}
	return quadtree.NewBoundingBox(rb.Position.X-hw, rb.Position.X+hw, rb.Position.Y-hh, rb.Position.Y+hh)
}

func NewRigidBody(id uint32, h int32, w int32, pos Vect2, vel Vect2, angle float64, mass int32) *RigidBody {
	return &RigidBody{
		ID:       id,
		Position: pos,
		Velocity: vel,
		Angle:    angle,
		Mass:     mass,
		InvMass:  1 / mass,
		Height:   h,
		Width:    w,
	}
}

// PhysicsEntityUpdate message linked to an Entity.
type PhysicsEntityUpdate struct {
	UpdateType byte // 2 == add, 3 == remove, 4 == physics update
	Body       *RigidBody
	Other      *RigidBody
}

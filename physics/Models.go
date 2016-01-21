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

func (rb RigidBody) BoundingBox() quadtree.BoundingBox {
	// if rb.Angle != 0 && rb.Angle != math.Pi && rb.Angle != math.Pi*2 {
	// 	ct := math.Cos(rb.Angle)
	// 	st := math.Sin(rb.Angle)
	//
	// 	A_y := (rb.Height / 2)
	// 	A_x := (rb.Width / 2)
	//
	// 	hct := int32(float64(rb.Height) * ct)
	// 	wct := int32(float64(rb.Width) * ct)
	// 	hst := int32(float64(rb.Height) * st)
	// 	wst := int32(float64(rb.Width) * st)
	//
	// 	var y_min, y_max, x_min, x_max int32
	// 	if rb.Angle > 0 {
	// 		if rb.Angle < math.Pi/2 {
	// 			y_min = A_y
	// 			y_max = A_y + hct + wst
	// 			x_min = A_x - hst
	// 			x_max = A_x + wct
	// 		} else {
	// 			// 90 <= theta <= 180
	// 			y_min = A_y + hct
	// 			y_max = A_y + wst
	// 			x_min = A_x - hst + wct
	// 			x_max = A_x
	// 		}
	// 	} else {
	// 		if rb.Angle > -math.Pi/2 {
	// 			y_min = A_y + wst
	// 			y_max = A_y + hct
	// 			x_min = A_x
	// 			x_max = A_x + wct - hst
	// 		} else {
	// 			// -180 <= theta <= -90
	// 			y_min = A_y + wst + hct
	// 			y_max = A_y
	// 			x_min = A_x + wct
	// 			x_max = A_x - hst
	// 		}
	// 	}
	// 	quadtree.NewBoundingBox(x_min, x_max, y_min, y_max)
	// }
	return quadtree.NewBoundingBox(rb.Position.X, rb.Position.X+rb.Width, rb.Position.Y, rb.Position.Y+rb.Height)
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
	UpdateType byte      // 2 == add, 3 == remove, 4 == physics update
	Body       RigidBody // Passed by value through channels
}

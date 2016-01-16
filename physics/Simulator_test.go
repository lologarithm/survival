package physics

import (
	"fmt"
	"math"
	"testing"
)

func TestTick(t *testing.T) {
	// 1. Create simple scene
	ss := &SimulatedSpace{Entities: map[uint32]*RigidBody{}}
	o1 := &RigidBody{}
	o1.Velocity = Vect2{1, 1}
	o1.Position = Vect2{0, 0}
	o1.Force = Vect2{0, 0}
	ss.Entities[1] = o1
	// 2. Make sure single tick correctly ticks.
	for i := float64(1); i < 50.0; i += 1 {
		ss.Tick(false)
		if !FloatCompare(o1.Position.X, i/50.0) {
			fmt.Printf("Incorrect X position after physics update. Expected: %f Actual: %f\n", i/50.0, o1.Position.X)
			t.FailNow()
		}
	}
	fmt.Println("Tick test passed.")
}

func FloatCompare(a float64, b float64) bool {
	if math.Abs(a-b) < 0.00001 {
		return true
	}
	return false
}

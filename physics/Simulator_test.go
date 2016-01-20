package physics

import (
	"fmt"
	"math"
	"math/rand"
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
	for i := int32(1); i < 50; i++ {
		ss.Tick(false)
		if o1.Position.X != i/50 {
			fmt.Printf("Incorrect X position after physics update. Expected: %d Actual: %d\n", i/50, o1.Position.X)
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

func BenchmarkTick(b *testing.B) {
	ss := &SimulatedSpace{Entities: map[uint32]*RigidBody{}}
	for i := uint32(0); i < 1000; i++ {
		o1 := &RigidBody{}
		o1.ID = i
		o1.Velocity = Vect2{1, 1}
		o1.Position = Vect2{int32(rand.Intn(1000)), int32(rand.Intn(1000))}
		o1.Force = Vect2{0, 0}
		ss.Entities[i] = o1
	}
	b.ResetTimer()
	// 2. Make sure single tick correctly ticks.
	for i := 0; i < b.N; i++ {
		ss.Tick(true)
	}
}

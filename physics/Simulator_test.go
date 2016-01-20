package physics

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestTick(t *testing.T) {
	// 1. Create simple scene
	ss := &SimulatedSpace{Entities: make([]*RigidBody, 2)}
	o1 := &RigidBody{}
	o1.Velocity = Vect2{1, 1}
	o1.Position = Vect2{0, 0}
	o1.Force = Vect2{0, 0}
	ss.AddEntity(o1, false)
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

func BenchmarkTick(b *testing.B) {
	ss := &SimulatedSpace{Entities: make([]*RigidBody, 10000)}
	for i := uint32(0); i < 10000; i++ {
		o1 := &RigidBody{}
		o1.ID = i
		o1.Velocity = Vect2{1, 1}
		o1.Position = Vect2{int32(rand.Intn(4000)), int32(rand.Intn(4000))}
		o1.Force = Vect2{0, 0}
		o1.Height = 10
		o1.Width = 10
		ss.AddEntity(o1, false)
	}
	b.ResetTimer()
	// 2. Make sure single tick correctly ticks.
	for i := 0; i < b.N; i++ {
		ss.Tick(true)
	}
}

package server

import "testing"

func TestGameSpawning(t *testing.T) {
	g := NewGame("A", nil, nil)
	g.Seed = 10

	g.SpawnChunk(0, 0)

	// for i, t := range g.World.Entities {
	// 	name := "Rock"
	// 	if t.EType == 2 {
	// 		name = "Tree"
	// 	}
	// 	fmt.Printf("Entity %s %3d: Size: %2d X:%3d Y:%3d\n", name, i, t.Height, t.X, t.Y)
	// }
}
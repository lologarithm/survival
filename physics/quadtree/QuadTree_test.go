/*
Copyright 2013 Volker Poplawski
*/

package quadtree

import (
	"log"
	"math"
	"testing"

	_ "fmt"
	"math/rand"
)

// Generates n BoundingBoxes in the range of frame with average width and height avgSize
func randomBoundingBoxes(n int, frame BoundingBox, avgSize float64) []BoundingBox {
	ret := make([]BoundingBox, n)

	for i := 0; i < len(ret); i++ {
		w := int32(rand.NormFloat64() * avgSize)
		h := int32(rand.NormFloat64() * avgSize)
		x := int32(rand.Float64()*float64(frame.SizeX())) + frame.MinX
		y := int32(rand.Float64()*float64(frame.SizeY())) + frame.MinY
		ret[i] = NewBoundingBox(x, int32(math.Min(float64(frame.MaxX), float64(x+w))), y, int32(math.Min(float64(frame.MaxY), float64(y+h))))
	}

	return ret
}

// Returns all elements of data which intersect query
func queryLinear(data []BoundingBox, query BoundingBox) (ret []BoundingBoxer) {
	for _, v := range data {
		if query.Intersects(v.Bounds()) {
			ret = append(ret, v)
		}
	}

	return ret
}

func compareBoundingBoxer(v1, v2 BoundingBoxer) bool {
	b1 := v1.Bounds()
	b2 := v2.Bounds()

	return b1.MinX == b2.MinX && b1.MaxX == b2.MaxX &&
		b1.MinY == b2.MinY && b2.MaxY == b2.MaxY
}

func lookupResults(r1, r2 []BoundingBoxer) int {
	for i, v1 := range r1 {
		found := false

		for _, v2 := range r2 {
			if compareBoundingBoxer(v1, v2) {
				found = true
				break
			}
		}

		if !found {
			return i
		}
	}

	return -1
}

// World-space extends from -1000..1000 in X and Y direction
var world BoundingBox = NewBoundingBox(-1000000, 1000000, -1000000, 1000000)

// Compary correctness of quad-tree results vs simple look-up on set of random rectangles
func TestQuadTreeRects(t *testing.T) {
	var rects []BoundingBox = randomBoundingBoxes(100*1000, world, 5)
	qt := NewQuadTree(world)

	for _, v := range rects {
		qt.Add(v)
	}

	queries := randomBoundingBoxes(1000, world, 10)

	for _, q := range queries {
		r1 := queryLinear(rects, q)
		r2 := qt.Query(q)

		if len(r1) != len(r2) {
			t.Errorf("r1 and r2 differ: %v   %v\n", r1, r2)
		}

		if i := lookupResults(r1, r2); i != -1 {
			t.Errorf("%v was not returned by QT\n", r1[i])
		}

		if i := lookupResults(r2, r1); i != -1 {
			t.Errorf("%v was not returned by brute-force\n", r2[i])
		}

	}
}

type mockBox struct {
	BoundingBox
	ID uint32
}

func (mb mockBox) BoxID() uint32 {
	return mb.ID
}

func TestQuadRemove(t *testing.T) {
	points := randomBoundingBoxes(10, world, 10)
	qt := NewQuadTree(world)

	boxes := make([]mockBox, len(points))
	for idx, v := range points {
		boxes[idx] = mockBox{BoundingBox: v, ID: uint32(idx + 1)}
		qt.Add(boxes[idx])
	}

	childs := qt.Query(world)
	if len(childs) != 10 {
		log.Printf("Failed to find all children in the world! expected: %d, actual: %d", 10, len(childs))
		t.FailNow()
	}

	found := qt.Remove(boxes[0])
	if !found {
		log.Printf("Box not found!? %d", boxes[0].ID)
	}
}

// Compary correctness of quad-tree results vs simple look-up on set of random points
func TestQuadTreePoints(t *testing.T) {
	var points []BoundingBox = randomBoundingBoxes(100*1000, world, 0)
	qt := NewQuadTree(world)

	for _, v := range points {
		qt.Add(v)
	}

	queries := randomBoundingBoxes(1000, world, 10)

	for _, q := range queries {
		r1 := queryLinear(points, q)
		r2 := qt.Query(q)

		if len(r1) != len(r2) {
			t.Errorf("r1 and r2 differ: %v   %v\n", r1, r2)
		}

		if i := lookupResults(r1, r2); i != -1 {
			t.Errorf("%v was not returned by QT\n", r1[i])
		}

		if i := lookupResults(r2, r1); i != -1 {
			t.Errorf("%v was not returned by brute-force\n", r2[i])
		}

	}
}

// Benchmark insertion into quad-tree
func BenchmarkInsert(b *testing.B) {
	b.StopTimer()

	var values []BoundingBox = randomBoundingBoxes(b.N, world, 5)
	qt := NewQuadTree(world)

	b.StartTimer()

	for _, v := range values {
		qt.Add(v)
	}
}

// A set of 10 million randomly distributed rectangles of avg size 5
var boxes10M []BoundingBox

// A set of 10 million randomly distributed points
var points10M []BoundingBox

// Benchmark quad-tree on set of rectangles
func BenchmarkRectsQuadtree(b *testing.B) {
	b.StopTimer()
	if boxes10M == nil {
		boxes10M = randomBoundingBoxes(10*1000*1000, world, 5)
		points10M = randomBoundingBoxes(10*1000*1000, world, 0)
	}
	rand.Seed(1)
	qt := NewQuadTree(world)

	for _, v := range boxes10M {
		qt.Add(v)
	}

	queries := randomBoundingBoxes(b.N, world, 10)

	b.StartTimer()
	for _, q := range queries {
		qt.Query(q)
	}
}

// Benchmark simple look up on set of rectangles
func BenchmarkRectsLinear(b *testing.B) {
	b.StopTimer()
	if boxes10M == nil {
		boxes10M = randomBoundingBoxes(10*1000*1000, world, 5)
	}
	rand.Seed(1)
	queries := randomBoundingBoxes(b.N, world, 10)

	b.StartTimer()
	for _, q := range queries {
		queryLinear(boxes10M, q)
	}
}

// Benchmark quad-tree on set of points
func BenchmarkPointsQuadtree(b *testing.B) {
	b.StopTimer()
	if points10M == nil {
		points10M = randomBoundingBoxes(10*1000*1000, world, 0)
	}
	rand.Seed(1)
	qt := NewQuadTree(world)

	for _, v := range points10M {
		qt.Add(v)
	}

	queries := randomBoundingBoxes(b.N, world, 10)

	b.StartTimer()
	for _, q := range queries {
		qt.Query(q)
	}
}

// Benchmark simple look-up on set of points
func BenchmarkPointsLinear(b *testing.B) {
	b.StopTimer()
	if points10M == nil {
		points10M = randomBoundingBoxes(10*1000*1000, world, 0)
	}
	rand.Seed(1)
	queries := randomBoundingBoxes(b.N, world, 10)

	b.StartTimer()
	for _, q := range queries {
		queryLinear(points10M, q)
	}
}

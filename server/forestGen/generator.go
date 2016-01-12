package forestGen

import (
	"math/rand"
	"time"
)

// Generate creates a new map.
func Generate(h, w int) *Map {
	m := NewMap()
	m.Tiles = make([][]Tile, h)
	for idx := range m.Tiles {
		m.Tiles[idx] = make([]Tile, w)
	}

	// big tree is diameter 7 circle
	// medium tree is diameter 5 circle
	// small tree is diameter 3 circle
	// sapling is 1

	rand.Seed(time.Now().UnixNano())

	numBigTree := ((h * w) / 49) / 25

	for i := 0; i < numBigTree; i++ {
		x := rand.Intn(w-10) + 5
		y := rand.Intn(h-10) + 5

		found := false
		for r := 0; r < 4; r++ {
			found = m.FindTree(x, y, r)

			if found {
				break
			}
		}

		if found {
			continue
		}

		for r := 0; r < 4; r++ {
			m.DrawCircle(x, y, r)
		}
	}
	return m
}

func (m *Map) FindTree(x, y, r int) bool {
	if r < 0 {
		return false
	}
	// Bresenham algorithm
	x1, y1, err := -r, 0, 2-2*r
	for {
		if m.Tiles[x-x1][y+y1] == Tree {
			return true
		}
		if m.Tiles[x-y1][y-x1] == Tree {
			return true
		}
		if m.Tiles[x+x1][y-y1] == Tree {
			return true
		}
		if m.Tiles[x+y1][y+x1] == Tree {
			return true
		}

		r = err
		if r > x1 {
			x1++
			err += x1*2 + 1
		}
		if r <= y1 {
			y1++
			err += y1*2 + 1
		}
		if x1 >= 0 {
			break
		}
	}
	return false
}

func (m *Map) DrawCircle(x, y, r int) {
	if r < 0 {
		return
	}
	// Bresenham algorithm
	x1, y1, err := -r, 0, 2-2*r
	for {
		m.Tiles[x-x1][y+y1] = Tree
		m.Tiles[x-y1][y-x1] = Tree
		m.Tiles[x+x1][y-y1] = Tree
		m.Tiles[x+y1][y+x1] = Tree

		r = err
		if r > x1 {
			x1++
			err += x1*2 + 1
		}
		if r <= y1 {
			y1++
			err += y1*2 + 1
		}
		if x1 >= 0 {
			break
		}
	}
}

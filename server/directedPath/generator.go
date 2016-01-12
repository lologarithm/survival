package directedPath

import (
	"math"
	"math/rand"
	"time"
)

// Generate creates a new map.
func Generate(h, w int) *Map {
	m := NewMap()
	m.Tiles = make([][]Tile, h) // 150 tall
	for idx := range m.Tiles {
		m.Tiles[idx] = make([]Tile, w) // 75 wide
	}
	rand.Seed(time.Now().UnixNano())
	// Setup starting room
	startX := rand.Intn(65) + 5
	startY := rand.Intn(25) + 5

	for x := startX - 5; x <= startX+5; x++ {
		m.Tiles[startY-5][x] = Wall
		m.Tiles[startY+5][x] = Wall
		for y := startY - 4; y < startY+5; y++ {
			m.Tiles[y][x] = Flat
		}
	}

	for y := startY - 5; y <= startY+5; y++ {
		m.Tiles[y][startX-5] = Wall
		m.Tiles[y][startX+5] = Wall
	}

	newX := startX
	newY := startY

	done := false
	for !done {
		if newY+(h/10) >= (h - 2) {
			done = true
			break
		}
		// 1. pick random direction & length
		dir := rand.Intn(w-20) + 10
		if newX+dir < w-1 {
			newX += dir
		} else if newX-dir > 1 {
			newX -= dir
		} else {
			// This option doesnt work!
			continue
		}

		newY += rand.Intn((h/10)-3) + 3
		// Now we make a 'path' from startX/Y to newX/Y
		angle := math.Atan2(float64(newY-startY), float64(newX-startX))
		flen := math.Sqrt(math.Pow(float64(dir), 2) + math.Pow(float64(newY-startY), 2))
		fy := startY - int(math.Sin(angle)*0.5)

		for l := float64(0.0); l <= flen; l += 0.5 {
			x := startX + int(math.Cos(angle)*l+0.5)
			y := startY + int(math.Sin(angle)*l+0.5)
			ey := startY + int(math.Sin(angle)*(l+0.25)+0.5)
			for i := fy - 2; i < ey+2; i++ {
				m.Tiles[i][x] = Flat
				if x == newX {
					t := x - 1
					if math.Cos(angle) > 0 {
						t = x + 1
					}
					m.Tiles[i][t] = Wall
				} else if x == startX {
					t := x - 1
					if math.Cos(angle) < 0 {
						t = x + 1
					}
					m.setIfEmpty(Wall, t, y)
				}
			}
			m.setIfEmpty(Wall, x, fy-3)
			m.setIfEmpty(Wall, x, ey+2)
			fy = y
		}
		startX = newX
		startY = newY
	}

	return m
}

func (m *Map) setIfEmpty(t Tile, x, y int) {
	if m.Tiles[y][x] == Empty {
		m.Tiles[y][x] = t
	}
}

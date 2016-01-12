package directedPath

import "bytes"

type Tile byte

const (
	Empty Tile = iota
	Flat  Tile = iota // 3 kinds of flat tiles allows for differences (varied) tiles.
	Flat2 Tile = iota
	Flat3 Tile = iota
	Wall  Tile = iota // same as flat, multiple blocked types
	Wall2 Tile = iota
	Wall3 Tile = iota
)

func (t Tile) String() string {
	switch t {
	case Empty:
		return " "
	case Flat, Flat2, Flat3:
		return "_"
	case Wall, Wall2, Wall3:
		return "#"
	}
	return "?"
}

type Map struct {
	Tiles [][]Tile
}

func NewMap() *Map {

	return &Map{}
}

func (m *Map) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString("")
	for y := len(m.Tiles) - 1; y > -1; y-- {
		for x := range m.Tiles[y] {
			buf.WriteString(m.Tiles[y][x].String())
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

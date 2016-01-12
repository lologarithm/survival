package forestGen

import "bytes"

type Tile byte

const (
	Flat Tile = iota
	Bush Tile = iota
	Tree Tile = iota
)

func (t Tile) String() string {
	switch t {
	case Flat:
		return "_"
	case Bush:
		return "#"
	case Tree:
		return "^"
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

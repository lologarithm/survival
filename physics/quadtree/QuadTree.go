/*
Copyright 2013 Volker Poplawski
Modified 2015 Ben Echols -- changed from float64 to int32 space
*/

// Package quadtree is a simple 2D implementation of the Quad-Tree data structure.
package quadtree

// MaxEntriesPerTile is the number of entries until a quad is split
const MaxEntriesPerTile = 16

// MaxLevels is the maximum depth of quad-tree (not counting the root node)
const MaxLevels = 10

// some constants for tile-indeces, for clarity
const (
	topRightTile    = 0
	topLeftTile     = 1
	bottomLeftTile  = 2
	bottomRightTile = 3
)

// BoundingBoxer interface allows arbitrary objects to be inserted into the quadtree since it only requires a bounding box.
type BoundingBoxer interface {
	Bounds() BoundingBox
	BoxID() uint32
	Clone() BoundingBoxer
}

// QuadTree is the core tree structure.
type QuadTree struct {
	root qtile
}

// NewQuadTree constructs an empty quad-tree
// bbox specifies the extends of the coordinate system.
func NewQuadTree(bbox BoundingBox) QuadTree {
	qt := QuadTree{qtile{BoundingBox: bbox}}

	return qt
}

// quad-tile / node of the quad-tree
type qtile struct {
	BoundingBox
	level    int             // level this tile is at. root is level 0
	contents []BoundingBoxer // values stored in this tile
	childs   [4]*qtile       // sub-tiles. none or four.
}

// Add a value to the quad-tree by trickle down from the root node.
func (qb *QuadTree) Add(v BoundingBoxer) {
	qb.root.add(v)
}

// Remove a value from the quad-tree by trickle down from the root node.
func (qb *QuadTree) Remove(v BoundingBoxer) bool {
	return qb.root.remove(v.Bounds(), v.BoxID())
}

// Move a box to new box
func (qb *QuadTree) Move(v BoundingBoxer, oldloc BoundingBox) int {
	return qb.root.move(v.BoxID(), v, oldloc, v.Bounds())
}

// Clone will create a full copy of this quadtree.
func (qb *QuadTree) Clone() *QuadTree {
	return &QuadTree{
		root: qb.root.clone(),
	}
}

// Query will return all objects which intersect the query box
func (qb *QuadTree) Query(bbox BoundingBox) []BoundingBoxer {
	return qb.root.query(bbox)
}

func (tile qtile) clone() qtile {
	ntile := qtile{
		BoundingBox: tile.BoundingBox,
		level:       tile.level,
		contents:    make([]BoundingBoxer, len(tile.contents)),
		childs:      [4]*qtile{},
	}
	for idx, bb := range tile.contents {
		ntile.contents[idx] = bb.Clone()
	}
	for idx, c := range tile.childs {
		if c != nil {
			cClone := c.clone()
			ntile.childs[idx] = &cClone
		}
	}
	return ntile
}

func (tile *qtile) add(v BoundingBoxer) {
	// look for sub-tile directly below this tile to accomodate value.
	if i := tile.findChildIndex(v.Bounds()); i < 0 {
		// no suitable sub-tile for value found.
		// either this tile has no childs or
		// value does not fit in any subtile.
		// store value at this level.
		tile.contents = append(tile.contents, v)

		// tile is split if exceeds it max number of entries and
		// has not childs already and max tree depth for this sub-tree not reached.
		if len(tile.contents) > MaxEntriesPerTile && tile.childs[topRightTile] == nil && tile.level < MaxLevels {
			tile.split()
		}
	} else {
		// suitable sub-tile for value found at index i.
		// recursivly add value.
		tile.childs[i].add(v)
	}
}

// return child index for BoundingBox
// returns -1 if quad has no children or BoundingBox does not fit into any child
func (tile *qtile) findChildIndex(bbox BoundingBox) int {
	if tile.childs[topRightTile] == nil {
		return -1
	}

	for i, child := range tile.childs {
		if child.Contains(bbox) {
			return i
		}
	}

	return -1
}

// create four child quads.
// distribute contents of this tiles on newly created childs.
func (tile *qtile) split() {
	mx := tile.MaxX/2.0 + tile.MinX/2.0
	my := tile.MaxY/2.0 + tile.MinY/2.0

	tile.childs[topRightTile] = &qtile{BoundingBox: NewBoundingBox(mx, tile.MaxX, my, tile.MaxY), level: tile.level + 1}
	tile.childs[topLeftTile] = &qtile{BoundingBox: NewBoundingBox(tile.MinX, mx, my, tile.MaxY), level: tile.level + 1}
	tile.childs[bottomLeftTile] = &qtile{BoundingBox: NewBoundingBox(tile.MinX, mx, tile.MinY, my), level: tile.level + 1}
	tile.childs[bottomRightTile] = &qtile{BoundingBox: NewBoundingBox(mx, tile.MaxX, tile.MinY, my), level: tile.level + 1}

	tempList := tile.contents

	// clear values on this tile
	tile.contents = []BoundingBoxer{}

	// reinsert from parent slice
	for _, v := range tempList {
		tile.add(v)
	}

}

// 0 == not found, 1 == found, needs home, 2 == completed
func (tile *qtile) move(id uint32, bb BoundingBoxer, old, newb BoundingBox) int {
	// end recursion if this tile does not intersect the query range
	if !tile.Intersects(old) {
		return 0
	}

	for vidx, v := range tile.contents {
		if id == v.BoxID() {
			if tile.Intersects(newb) {
				return 2
			}
			// REMOVE IT
			tile.contents = append(tile.contents[:vidx], tile.contents[vidx+1:]...)
			return 1
		}
	}

	// recurse into childs (if any)
	if tile.childs[topRightTile] != nil {
		for _, child := range tile.childs {
			cr := child.move(id, bb, old, newb)
			// Not found, continue searching
			if cr == 0 {
				continue
			}
			// Found&Handled, return now!
			if cr == 2 {
				return 2
			}
			// Found but not inserted
			if !tile.Intersects(newb) {
				// New loc doesn't fit in here, pass it up!
				return 1
			}
			// Add new loc here!
			tile.add(bb)
			return 2
		}
	}

	return 0
}

func (tile *qtile) remove(qbox BoundingBox, id uint32) bool {
	// end recursion if this tile does not intersect the query range
	if !tile.Intersects(qbox) {
		return false
	}

	// return candidates at this tile
	for vidx, v := range tile.contents {
		if id == v.BoxID() {
			// REMOVE IT
			tile.contents = append(tile.contents[:vidx], tile.contents[vidx+1:]...)
			return true
		}
	}

	// recurse into childs (if any)
	if tile.childs[topRightTile] != nil {
		for _, child := range tile.childs {
			if child.remove(qbox, id) {
				return true
			}
		}
	}

	return false
}

func (tile *qtile) query(qbox BoundingBox) []BoundingBoxer {
	ret := []BoundingBoxer{}
	// end recursion if this tile does not intersect the query range
	if !tile.Intersects(qbox) {
		return ret
	}

	// return candidates at this tile
	for _, v := range tile.contents {
		if qbox.Intersects(v.Bounds()) {
			ret = append(ret, v)
		}
	}

	// recurse into childs (if any)
	if tile.childs[topRightTile] != nil {
		for _, child := range tile.childs {
			ret = append(ret, child.query(qbox)...)
		}
	}

	return ret
}

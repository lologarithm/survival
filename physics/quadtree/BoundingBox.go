/*
Copyright 2013 Volker Poplawski
Modified 2015 Ben Echols -- changed float64 to int32 world space.
*/

package quadtree

// Use NewBoundingBox() to construct a BoundingBox object
type BoundingBox struct {
	MinX, MaxX, MinY, MaxY int32
}

func NewBoundingBox(xa, xb, ya, yb int32) BoundingBox {
	return BoundingBox{xa, xb, ya, yb}
}

// Make BoundingBox implement the BoundingBoxer interface
func (b BoundingBox) BoundingBox() BoundingBox {
	return b
}

func (b BoundingBox) SizeX() int32 {
	return b.MaxX - b.MinX
}

func (b BoundingBox) SizeY() int32 {
	return b.MaxY - b.MinY
}

// Returns true if o intersects this
func (b BoundingBox) Intersects(o BoundingBox) bool {
	return b.MinX < o.MaxX && b.MinY < o.MaxY &&
		b.MaxX > o.MinX && b.MaxY > o.MinY
}

// Returns true if o is within this
func (b BoundingBox) Contains(o BoundingBox) bool {
	return b.MinX <= o.MinX && b.MinY <= o.MinY &&
		b.MaxX >= o.MaxX && b.MaxY >= o.MaxY
}

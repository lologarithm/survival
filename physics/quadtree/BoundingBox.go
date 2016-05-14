/*
Copyright 2013 Volker Poplawski
Modified 2015 Ben Echols -- changed float64 to int32 world space.
*/

package quadtree

// BoundingBox represents a square area of possible intersection.
// Use NewBoundingBox() to construct a BoundingBox object
type BoundingBox struct {
	MinX, MaxX, MinY, MaxY int32
}

// NewBoundingBox is the constructor for the BoundingBox struct.
func NewBoundingBox(xa, xb, ya, yb int32) BoundingBox {
	return BoundingBox{xa, xb, ya, yb}
}

// BoundingBox should implement the BoundingBoxer interface.
func (b BoundingBox) Bounds() BoundingBox {
	return b
}

// SizeX returns width of the box.
func (b BoundingBox) SizeX() int32 {
	return b.MaxX - b.MinX
}

// SizeY returns height of the box
func (b BoundingBox) SizeY() int32 {
	return b.MaxY - b.MinY
}

// Intersects returns true if o intersects this
func (b BoundingBox) Intersects(o BoundingBox) bool {
	if b.MaxX < o.MinX {
		return false // a is left of b
	}

	if b.MinX > o.MaxX {
		return false // a is right of b
	}
	if b.MaxY < o.MinY {
		return false // a is above b
	}

	if b.MinY > o.MaxY {
		return false // a is below b
	}

	return true // boxes overlap
}

// Contains returns true if o is within this
func (b BoundingBox) Contains(o BoundingBox) bool {
	return b.MinX <= o.MinX && b.MinY <= o.MinY &&
		b.MaxX >= o.MaxX && b.MaxY >= o.MaxY
}

// BoxID returns 0 for bounding box because the box itself doesnt have an ID
func (b BoundingBox) BoxID() uint32 {
	return 0
}

// Clone exists to fulfill the interface
func (b BoundingBox) Clone() BoundingBoxer {
	return b
}

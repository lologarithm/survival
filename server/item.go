package server

// ItemSlot is the constant for slot type
type ItemSlot uint8

// Item slot constants
const (
	OtherSlot ItemSlot = iota
	MainHandSlot
	OffHandSlot
	BodySlot
	FeetSlot
	SpiritSlot
)

// Item represents
type Item struct {
	Stats // Item stats -- added to the bearer's own stats when equipped.

	ID   uint32   // Unique ID for this item
	Slot ItemSlot // Slot this item can be equipped in.
}

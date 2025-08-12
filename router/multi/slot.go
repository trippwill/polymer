package multi

// Slot identifies a [Routable] in a [Routed] or [Routed3] struct.
//
//go:generate stringer -type=Slot
type Slot uint8

const (
	SlotSkip Slot = iota // Skip: no slot handles the message.
	SlotT
	SlotU
	SlotV
)

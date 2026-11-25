// Package binpack provides generic 1-dimensional bin packing algorithms.
// It operates on ephemeral types that wrap the actual domain objects with
// their packing-relevant properties (length, spacing).
package binpack

// Item represents an ephemeral item to be packed.
// It contains the length needed for packing and the original item of type T.
// Spacing is handled at the domain level, not in the generic packer.
type Item[T any] struct {
	// Original is the actual domain object being packed (typically a pointer)
	Original T
	// Length is the length required for this item
	Length float64
}

// NewItem creates a new ephemeral packing item
func NewItem[T any](original T, length float64) Item[T] {
	return Item[T]{
		Original: original,
		Length:   length,
	}
}

// Bin represents an ephemeral bin (stock) that can hold items.
// It contains the capacity information and the original bin of type B.
type Bin[B any, I any] struct {
	// Original is the actual domain object representing the bin/stock (typically a pointer)
	Original B
	// Capacity is the total length capacity of this bin
	Capacity float64
	// Items contains the items packed into this bin
	Items []Item[I]
	// usedLength tracks the total used length (including spacing)
	usedLength float64
}

// NewBin creates a new ephemeral packing bin
func NewBin[B any, I any](original B, capacity float64) *Bin[B, I] {
	return &Bin[B, I]{
		Original:   original,
		Capacity:   capacity,
		Items:      make([]Item[I], 0),
		usedLength: 0,
	}
}

// AvailableLength returns the remaining available length in the bin
func (b *Bin[B, I]) AvailableLength() float64 {
	return b.Capacity - b.usedLength
}

// UsedLength returns the total used length in the bin
func (b *Bin[B, I]) UsedLength() float64 {
	return b.usedLength
}

// CanFit checks if the given item can fit in the remaining space
func (b *Bin[B, I]) CanFit(item Item[I]) bool {
	return item.Length <= b.AvailableLength()
}

// AddItem adds an item to the bin if it fits
// Returns true if successful, false if the item doesn't fit
func (b *Bin[B, I]) AddItem(item Item[I]) bool {
	if !b.CanFit(item) {
		return false
	}
	b.Items = append(b.Items, item)
	b.usedLength += item.Length
	return true
}

// Waste returns the unused capacity in the bin
func (b *Bin[B, I]) Waste() float64 {
	return b.AvailableLength()
}

// Utilization returns the percentage of capacity used (0-100)
func (b *Bin[B, I]) Utilization() float64 {
	if b.Capacity == 0 {
		return 0
	}
	return (b.usedLength / b.Capacity) * 100
}

// Reset clears all items from the bin
func (b *Bin[B, I]) Reset() {
	b.Items = b.Items[:0]
	b.usedLength = 0
}

// Clone creates a deep copy of the bin with the same original reference
func (b *Bin[B, I]) Clone() *Bin[B, I] {
	cloned := NewBin[B, I](b.Original, b.Capacity)
	cloned.Items = make([]Item[I], len(b.Items))
	copy(cloned.Items, b.Items)
	cloned.usedLength = b.usedLength
	return cloned
}

// PackingResult contains the result of a packing operation
type PackingResult[B any, I any] struct {
	// PackedBins contains all bins that have at least one item
	PackedBins []*Bin[B, I]
	// UnpackedItems contains items that could not be packed
	UnpackedItems []Item[I]
	// Statistics about the packing operation
	Statistics Statistics
}

// Statistics provides packing performance metrics
type Statistics struct {
	TotalItems          int
	PackedItems         int
	UnpackedItems       int
	TotalBinsUsed       int
	MaterialUtilization float64 // overall percentage
	TotalWaste          float64
	ExecutionTime       string
	Algorithm           string
}

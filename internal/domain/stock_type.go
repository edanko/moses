package domain

// StockType distinguishes stock bars, consumed remnants and explicit
// over-length bars created for parts that do not fit the configured stock.
type StockType string

const (
	StockTypeBar        StockType = "BAR"
	StockTypeRemnant    StockType = "REMNANT"
	StockTypeOverlength StockType = "OVERLENGTH"
)

// Values returns all supported stock types.
func (StockType) Values() []string {
	return []string{
		string(StockTypeBar),
		string(StockTypeRemnant),
		string(StockTypeOverlength),
	}
}

// IsValid reports whether the type is known.
func (k StockType) IsValid() bool {
	switch k {
	case StockTypeBar, StockTypeRemnant, StockTypeOverlength:
		return true
	default:
		return false
	}
}

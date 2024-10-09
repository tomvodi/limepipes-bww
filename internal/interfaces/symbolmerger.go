package interfaces

import "github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"

type SymbolMerger interface {
	// MergeSymbols merges the right symbol into the left symbol
	// If the symbols can be merged, the function returns true and the right
	// symbol can be discarded. If the symbols cannot be merged, the function
	// returns false and the right symbol should be kept.
	MergeSymbols(
		left *symbols.Symbol,
		right *symbols.Symbol,
	) bool
}

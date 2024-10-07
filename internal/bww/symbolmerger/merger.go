package symbolmerger

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-bww/internal/interfaces"
)

var mergers []interfaces.SymbolMerger

type CollectedMerger struct {
}

func (c *CollectedMerger) MergeSymbols(
	left *symbols.Symbol,
	right *symbols.Symbol,
) bool {
	for _, merger := range mergers {
		if merger.MergeSymbols(left, right) {
			return true
		}
	}

	return false
}

func NewCollectedMerger() *CollectedMerger {
	return &CollectedMerger{}
}

package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
)

type SymbolMapper interface {
	IsTimeSignature(token string) bool
	TimeSigForToken(token string) (*measure.TimeSignature, error)
	SymbolForToken(token string) (*symbols.Symbol, error)
}

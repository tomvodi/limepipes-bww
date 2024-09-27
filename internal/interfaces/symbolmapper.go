package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/barline"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
)

type SymbolMapper interface {
	IsTimeSignature(token string) bool
	TimeSigForToken(token string) (*measure.TimeSignature, error)
	BarlineForToken(token string) (*barline.Barline, error)
	SymbolForToken(token string) (*symbols.Symbol, error)
}

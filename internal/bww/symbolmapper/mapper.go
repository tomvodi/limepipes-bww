package symbolmapper

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"golang.org/x/exp/maps"
	"slices"
)

type Mapper struct {
	symbols     map[string]*symbols.Symbol
	timeSigKeys []string
}

func (m *Mapper) IsTimeSignature(token string) bool {
	return slices.Contains(m.timeSigKeys, token)
}

func (m *Mapper) TimeSigForToken(token string) (*measure.TimeSignature, error) {
	sig, ok := TimeSignatureMap[token]
	if !ok {
		return nil, common.ErrSymbolNotFound
	}

	return sig, nil
}

func (m *Mapper) SymbolForToken(token string) (*symbols.Symbol, error) {
	sym, ok := m.symbols[token]
	if !ok {
		return nil, common.ErrSymbolNotFound
	}
	return sym, nil
}

func (m *Mapper) init() {
	m.timeSigKeys = make([]string, 0, len(TimeSignatureMap))
	for _, k := range maps.Keys(TimeSignatureMap) {
		m.timeSigKeys = append(m.timeSigKeys, k)
	}

	m.symbols = make(map[string]*symbols.Symbol)
	maps.Copy(m.symbols, NewMelodyNotesMap())
}

func New() *Mapper {
	m := &Mapper{}
	m.init()
	return m
}

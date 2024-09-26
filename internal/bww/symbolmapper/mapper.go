package symbolmapper

import (
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"slices"
)

var symbolsMap = map[string]*symbols.Symbol{}
var timeSignatureMap = map[string]*measure.TimeSignature{}
var timeSignatureKeys = []string{}

type Mapper struct {
}

func (m *Mapper) IsTimeSignature(token string) bool {
	return slices.Contains(timeSignatureKeys, token)
}

func (m *Mapper) TimeSigForToken(token string) (*measure.TimeSignature, error) {
	sig, ok := timeSignatureMap[token]
	if !ok {
		return nil, common.ErrSymbolNotFound
	}

	return sig, nil
}

func (m *Mapper) SymbolForToken(token string) (*symbols.Symbol, error) {
	sym, ok := symbolsMap[token]
	if !ok {
		return nil, common.ErrSymbolNotFound
	}
	// If nil symbol is found, Symbols should be skipped
	if sym == nil {
		return nil, common.ErrSymbolSkip
	}

	symCopy := &symbols.Symbol{}
	err := copier.Copy(symCopy, sym)
	if err != nil {
		return nil, fmt.Errorf("failed creating symbol copy: %v", err)
	}
	return symCopy, nil
}

func New() *Mapper {
	m := &Mapper{}
	return m
}

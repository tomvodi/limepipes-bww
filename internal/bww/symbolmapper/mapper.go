package symbolmapper

import (
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/barline"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"slices"
)

var symbolsMap = map[string]*symbols.Symbol{}
var piobSymbols []string
var timeSignatureMap = map[string]*measure.TimeSignature{}
var barlineMap = map[string]*barline.Barline{}
var timeSignatureKeys = []string{}

type Mapper struct {
}

func (m *Mapper) BarlineForToken(
	token string,
) (*barline.Barline, error) {
	bl, ok := barlineMap[token]
	if !ok {
		return nil, common.ErrSymbolNotFound
	}

	return bl, nil
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
		if slices.Contains(piobSymbols, token) {
			return nil, common.ErrPiobNotSupported
		}

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

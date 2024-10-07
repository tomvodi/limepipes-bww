// Package bww contains the definition for a bww file structure and a
// converter to convert the file structure tune into a music model tune.
package bww

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/tune"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"github.com/tomvodi/limepipes-plugin-bww/internal/filestructure"
	"github.com/tomvodi/limepipes-plugin-bww/internal/interfaces"
)

type Converter struct {
	mapper interfaces.SymbolMapper
	merger interfaces.SymbolMerger
}

func (c *Converter) Convert(
	fst *filestructure.Tune,
) (*tune.Tune, error) {
	t := &tune.Tune{}
	fillTuneWithHeader(t, fst.Header)

	for _, m := range fst.Measures {
		meas := &measure.Measure{}
		err := c.fillMeasure(meas, m)
		if err != nil {
			return nil, err
		}
		t.Measures = append(t.Measures, meas)
	}

	return t, nil
}

func fillTuneWithHeader(
	t *tune.Tune,
	h *filestructure.TuneHeader,
) {
	t.Title = string(h.Title)
	t.Type = string(h.Type)
	t.Composer = string(h.Composer)
	t.Tempo = uint32(h.Tempo)

	fillTuneFooter(t, h.Footer)
	fillTuneInlineTexts(t, h.InlineTexts)
	fillTuneComments(t, h.Comments)
}

func fillTuneFooter(
	t *tune.Tune,
	f []filestructure.TuneFooter,
) {
	if len(f) == 0 {
		return
	}

	t.Footer = make([]string, len(f))
	for i, ft := range f {
		t.Footer[i] = string(ft)
	}
}

func fillTuneInlineTexts(
	t *tune.Tune,
	f []filestructure.TuneInline,
) {
	if len(f) == 0 {
		return
	}

	t.InlineTexts = make([]string, len(f))
	for i, ft := range f {
		t.InlineTexts[i] = string(ft)
	}
}

func fillTuneComments(
	t *tune.Tune,
	f []filestructure.TuneComment,
) {
	if len(f) == 0 {
		return
	}

	t.Comments = make([]string, len(f))
	for i, ft := range f {
		t.Comments[i] = string(ft)
	}
}

func (c *Converter) fillMeasure(
	dest *measure.Measure,
	src *filestructure.Measure,
) error {
	fillInlineTextAndComments(dest, src)
	c.setMeasureBarlines(dest, src)

	for _, s := range src.Symbols {
		err := c.addSymbolToMeasure(dest, s)
		if err != nil {
			return err
		}
	}

	return nil
}

// addSymbolToMeasure adds a symbol to the measure. If the symbol is a time signature,
// if it is a time signature symbol, it sets that to the measure.
// If the symbol
func (c *Converter) addSymbolToMeasure(
	dest *measure.Measure,
	s *filestructure.MusicSymbol,
) error {
	timeSigHandled, err := c.setPossibleTimeSignature(dest, s)
	if err != nil {
		return err
	}
	if timeSigHandled {
		return nil
	}

	sym, err := c.getConvertedSymbol(s)
	if errors.Is(err, common.ErrSymbolSkip) {
		return nil
	}

	prevSym := dest.LastSymbol()
	if prevSym != nil && c.merger.MergeSymbols(prevSym, sym) {
		return nil
	}

	dest.Symbols = append(dest.Symbols, sym)

	return nil
}

// setPossibleTimeSignature checks if the symbol is a time signature and sets it
// if it is. If it is not, it returns false.
func (c *Converter) setPossibleTimeSignature(
	dest *measure.Measure,
	s *filestructure.MusicSymbol,
) (bool, error) {
	if !c.mapper.IsTimeSignature(s.Text) {
		return false, nil
	}

	ts, err := c.mapper.TimeSigForToken(s.Text)
	if err != nil {
		return false, err
	}
	dest.Time = ts
	return true, nil
}

func (c *Converter) getConvertedSymbol(
	s *filestructure.MusicSymbol,
) (*symbols.Symbol, error) {
	sym, err := c.convertSymbol(s)
	if err != nil {
		return nil, err
	}

	if len(s.InlineTexts) > 0 {
		it := toStringSlice[filestructure.InlineText](s.InlineTexts)
		sym.InlineTexts = append(sym.InlineTexts, it...)
	}

	if len(s.Comments) > 0 {
		c := toStringSlice[filestructure.InlineComment](s.Comments)
		sym.Comments = append(sym.Comments, c...)
	}

	return sym, nil
}

func fillInlineTextAndComments(
	dest *measure.Measure,
	src *filestructure.Measure,
) {
	addInlineTexts(dest, toStringSlice[filestructure.InlineText](src.InlineTexts))
	addInlineTexts(dest, toStringSlice[filestructure.StaffInline](src.StaffInlineTexts))

	addComments(dest, toStringSlice[filestructure.InlineComment](src.InlineComments))
	addComments(dest, toStringSlice[filestructure.StaffComment](src.StaffComments))
}

func (c *Converter) setMeasureBarlines(
	dest *measure.Measure,
	src *filestructure.Measure,
) {
	if src.LeftBarline != "" {
		bl, err := c.mapper.BarlineForToken(string(src.LeftBarline))
		// there should be only common.ErrSymbolNotFound possible
		if errors.Is(err, common.ErrSymbolNotFound) {
			log.Fatal().Msgf(
				"barline %s not found",
				src.LeftBarline,
			)
		}

		dest.LeftBarline = bl
	}

	if src.RightBarline == "" {
		return
	}

	bl, err := c.mapper.BarlineForToken(string(src.RightBarline))
	// there should be only common.ErrSymbolNotFound possible
	if errors.Is(err, common.ErrSymbolNotFound) {
		log.Fatal().Msgf(
			"barline %s not found",
			src.RightBarline,
		)
	}

	dest.RightBarline = bl
}

type MeasureText interface {
	filestructure.InlineText |
		filestructure.InlineComment |
		filestructure.StaffInline |
		filestructure.StaffComment
}

func toStringSlice[T MeasureText](
	texts []T,
) []string {
	if len(texts) == 0 {
		return nil
	}

	strs := make([]string, len(texts))
	for i, t := range texts {
		strs[i] = string(t)
	}

	return strs
}

func addInlineTexts(
	staff *measure.Measure,
	texts []string,
) {
	if len(texts) == 0 {
		return
	}

	if len(staff.InlineTexts) == 0 {
		staff.InlineTexts = make([]string, 0)
	}

	staff.InlineTexts = append(staff.InlineTexts, texts...)
}

func addComments(
	staff *measure.Measure,
	texts []string,
) {
	if len(texts) == 0 {
		return
	}

	if len(staff.Comments) == 0 {
		staff.Comments = make([]string, 0)
	}

	staff.Comments = append(staff.Comments, texts...)
}

func (c *Converter) convertSymbol(
	ms *filestructure.MusicSymbol,
) (*symbols.Symbol, error) {
	if ms.TempoChange > 0 {
		tc := uint64(ms.TempoChange)
		return &symbols.Symbol{
			TempoChange: &tc,
		}, nil
	}

	if ms.Text == "" {
		return nil, fmt.Errorf("empty symbol text for %+v", ms)
	}

	sym, err := c.mapper.SymbolForToken(ms.Text)
	if errors.Is(err, common.ErrSymbolNotFound) {
		return nil, fmt.Errorf(
			"symbol %s not found: line %d, column %d",
			ms.Text,
			ms.Pos.Line,
			ms.Pos.Column,
		)
	}
	if err != nil {
		return nil, err
	}

	return sym, nil
}

func NewConverter(
	mapper interfaces.SymbolMapper,
	merger interfaces.SymbolMerger,
) *Converter {
	return &Converter{
		mapper: mapper,
		merger: merger,
	}
}

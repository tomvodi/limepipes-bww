// Package filestructure contains the definition for a bww file structure and a
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

	if len(h.Footer) > 0 {
		t.Footer = make([]string, len(h.Footer))
		for i, f := range h.Footer {
			t.Footer[i] = string(f)
		}
	}

	if len(h.InlineTexts) > 0 {
		t.InlineText = make([]string, len(h.InlineTexts))
		for i, it := range h.InlineTexts {
			t.InlineText[i] = string(it)
		}
	}

	if len(h.Comments) > 0 {
		t.Comments = make([]string, len(h.Comments))
		for i, c := range h.Comments {
			t.Comments[i] = string(c)
		}
	}
}

func (c *Converter) fillMeasure(
	dest *measure.Measure,
	src *filestructure.Measure,
) error {
	fillInlineTextAndComments(dest, src)
	c.setBarlines(dest, src)

	if len(src.Symbols) == 0 {
		return nil
	}

	for _, s := range src.Symbols {
		if c.mapper.IsTimeSignature(s.Text) {
			ts, err := c.mapper.TimeSigForToken(s.Text)
			if err != nil {
				return err
			}
			dest.Time = ts
			continue
		}

		sym, err := c.convertSymbol(s)
		if errors.Is(err, common.ErrSymbolSkip) {
			continue
		}

		if err != nil {
			return err
		}

		if len(s.InlineTexts) > 0 {
			it := toStringSlice[filestructure.InlineText](s.InlineTexts)
			for _, t := range it {
				sym.InlineTexts = append(sym.InlineTexts, t)
			}
		}

		if len(s.Comments) > 0 {
			c := toStringSlice[filestructure.InlineComment](s.Comments)
			for _, t := range c {
				sym.Comments = append(sym.Comments, t)
			}
		}

		prevSym := dest.LastSymbol()
		if prevSym != nil && prevSym.CanBeMergedWith(sym) {
			prevSym.MergeWith(sym)
			continue
		}

		dest.Symbols = append(dest.Symbols, sym)
	}

	return nil
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

func (c *Converter) setBarlines(
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

	for _, it := range texts {
		staff.InlineTexts = append(staff.InlineTexts, it)
	}
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

	for _, it := range texts {
		staff.Comments = append(staff.Comments, it)
	}
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

/*
func (c *Converter) convertGrammarToModel(

	t *Tune,

	) (*tune.Tune, error) {
		var tunes musicmodel.MusicModel

		var newTune *tune.Tune
		staffCtx := &staffContext{
			PendingOldTie: pitch.Pitch_NoPitch,
		}
		for i, t := range t.TuneDefs {
			if t.Tune.Header.HasTitle() {
				if newTune != nil {
					tunes = append(tunes, newTune)
				}
				newTune = &tune.Tune{}
				if err := fillTuneWithParameter(newTune, t.Header.TuneParameter); err != nil {
					return nil, err
				}
			} else {
				if i == 0 && newTune == nil {
					log.Warn().Msgf("first tune doesn't have a title. Setting it to 'no title'")
					newTune = &tune.Tune{}
					if err := fillTuneWithParameter(newTune, t.Header.TuneParameter); err != nil {
						return nil, err
					}
					newTune.Title = "no title"
				} else {
					staffCtx.NextMeasureComments = t.Header.GetComments()
					staffCtx.NextMeasureInLineText = t.Header.GetInlineTexts()
				}
			}

			// TODO when tempo only of first t, set to other tunes as well?
			if err := c.fillTunePartsFromStaves(newTune, t.Body.Staffs, staffCtx); err != nil {
				return nil, err
			}
		}
		tunes = append(tunes, newTune)

		return tunes, nil
	}

func fillTuneWithParameter(

	tune *tune.Tune,
	params []*common.TuneParameter,

	) error {
		for _, param := range params {
			if param.Tempo != nil {
				tempo, err := strconv.ParseUint(param.Tempo.Tempo, 10, 64)
				if err != nil {
					return fmt.Errorf("failed parsing tune tempo: %s", err.Error())
				}

				tune.Tempo = uint32(tempo)
			}

			if param.Description != nil {
				firstParam := param.Description.ParamList[0]
				text := param.Description.Text
				if firstParam == common.TitleParameter {
					tune.Title = text
				}
				if firstParam == common.TypeParameter {
					tune.Type = text
				}
				if firstParam == common.ComposerParameter {
					tune.Composer = text
				}
				if firstParam == common.FooterParameter {
					tune.Footer = append(tune.Footer, text)
				}
				if firstParam == common.InlineParameter {
					tune.InlineText = append(tune.InlineText, text)
				}
			}

			if param.Comment != "" {
				tune.Comments = append(tune.Comments, param.Comment)
			}
		}

		return nil
	}

func (c *Converter) fillTunePartsFromStaves(

	tune *tune.Tune,
	staves []*common.Staff,
	staffCtx *staffContext,

	) error {
		var measures []*measure.Measure

		for _, stave := range staves {
			staveMeasures, err := c.getMeasuresFromStave(stave, staffCtx)
			if err != nil {
				return err
			}
			staffCtx.PreviousStaveMeasures = staveMeasures

			measures = append(measures, staveMeasures...)
			if stave.End != "" &&
				stave.Dalsegno == nil &&
				stave.DacapoAlFine == nil &&
				len(measures) >= 1 {
				lastMeasure := measures[len(measures)-1]
				var rightBarline *barline.Barline
				if stave.End == "''!I" {
					rightBarline = &barline.Barline{
						Type: barline.Type_Heavy,
						Time: barline.Time_Repeat,
					}
				}
				if stave.End == "!I" {
					rightBarline = &barline.Barline{
						Type: barline.Type_Heavy,
					}
				}
				lastMeasure.RightBarline = rightBarline
			}
		}
		tune.Measures = append(tune.Measures, measures...)

		return nil
	}

func (c *Converter) getMeasuresFromStave(

	stave *common.Staff,
	ctx *staffContext,

	) ([]*measure.Measure, error) {
		var measures []*measure.Measure
		currMeasure := &measure.Measure{}
		currMeasure.InlineText = ctx.NextMeasureInLineText
		currMeasure.Comments = ctx.NextMeasureComments
		ctx.NextMeasureInLineText = nil
		ctx.NextMeasureComments = nil

		for _, staffSym := range stave.Symbols {
			// if staffSym bar or part start => new measure
			// currMeasure to return measures
			if staffSym.Barline != nil ||
				staffSym.PartStart != nil ||
				staffSym.NextStaffStart != nil {
				measures = cleanupAndAppendMeasure(measures, currMeasure)

				var leftBarline *barline.Barline
				if staffSym.PartStart != nil {
					leftBarline = &barline.Barline{
						Type: barline.Type_Heavy,
					}
					ps := *staffSym.PartStart
					if ps == "I!''" {
						leftBarline.Time = barline.Time_Repeat
					}
				}

				currMeasure = &measure.Measure{
					LeftBarline: leftBarline,
				}

				continue
			}

			if staffSym.Space != nil {
				continue
			}

			if staffSym.Comment != nil {
				handleCommentForMeasure(currMeasure, *staffSym.Comment)
				continue
			}
			if staffSym.Description != nil {
				handleCommentForMeasure(currMeasure, staffSym.Description.Text)
				continue
			}

			if staffSym.Tempo != nil {
				t, err := strconv.ParseUint(staffSym.Tempo.Tempo, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("failed parsing tempo at line %d column %d: %s",
						staffSym.Pos.Line,
						staffSym.Pos.Column,
						err.Error(),
					)
				}

				currMeasure.Symbols = append(currMeasure.Symbols, &symbols.Symbol{
					TempoChange: &t,
				})
				continue
			}

			if staffSym.MusicSymbol == nil {
				log.Fatal().Msgf("staff symbol not handled: %v", staffSym)
			}

			if c.mapper.IsTimeSignature(*staffSym.MusicSymbol) {
				ts, err := c.mapper.TimeSigForToken(*staffSym.MusicSymbol)
				if err != nil {
					return nil, err
				}
				currMeasure.Time = ts
				continue
			}

			if staffSym.MusicSymbol != nil {
				sym, err := c.mapper.SymbolForToken(*staffSym.MusicSymbol)
				if errors.Is(err, common.ErrSymbolNotFound) {
					return nil, fmt.Errorf(
						"symbol %s not found: line %d, column %d",
						*staffSym.MusicSymbol,
						staffSym.Pos.Line,
						staffSym.Pos.Column,
					)
				}
				if err != nil {
					return nil, err
				}

				prevSym := currMeasure.LastSymbol()
				if prevSym != nil && prevSym.CanBeMergedWith(sym) {
					prevSym.MergeWith(sym)
					continue
				}

				currMeasure.Symbols = append(currMeasure.Symbols, sym)
				continue
			}
		}

		// append time line end to measure symbols if it was outside of stave
		if stave.TimelineEnd != nil {
			tl := newTimeLineEnd(stave.TimelineEnd)
			currMeasure.Symbols = append(currMeasure.Symbols, tl)
		}
		if stave.Dalsegno != nil {
			currMeasure.RightBarline = &barline.Barline{
				Type: barline.Type_Regular,
				Time: barline.Time_Dalsegno,
			}
		}
		if stave.DacapoAlFine != nil {
			currMeasure.RightBarline = &barline.Barline{
				Type: barline.Type_Regular,
				Time: barline.Time_DacapoAlFine,
			}
		}

		measures = cleanupAndAppendMeasure(measures, currMeasure)
		return measures, nil
	}

func handleCommentForMeasure(

	measure *measure.Measure,
	comment string,

	) {
		if measure == nil {
			return
		}

		if len(measure.Symbols) == 0 {
			measure.Comments = append(measure.Comments, comment)
			return
		}

		if measure.LastSymbol().IsNote() {
			measure.LastSymbol().Note.Comment = comment
			return
		}
	}

func cleanupAndAppendMeasure(

	measures []*measure.Measure,
	m *measure.Measure,

	) []*measure.Measure {
		cleanupMeasure(m)
		if len(m.Symbols) == 0 {
			m.Symbols = nil
		}

		return append(measures, m)
	}

// cleanupMeasure removes invalid symbols from the measure
// this may be the case for the accidentals at the beginning of the measure which are
// indicating the key of the measure. For bagpipes the key is always sharpf sharpc,
// so we delete these symbols here.

	func cleanupMeasure(meas *measure.Measure) {
		for _, symbol := range meas.Symbols {
			if symbol.Note == nil {
				continue
			}

			if symbol.Note.IsOnlyAccidental() {
				idx := symbolIndexOf(meas.Symbols, symbol)
				if idx == -1 {
					log.Error().Msgf("symbol index could not be found for cleanup of measure")
				} else {
					meas.Symbols = removeSymbol(meas.Symbols, idx)
				}
			}
		}
	}

	func symbolIndexOf(symbols []*symbols.Symbol, findSym *symbols.Symbol) int {
		for i, symbol := range symbols {
			if reflect.DeepEqual(symbol, findSym) {
				return i
			}
		}
		return -1
	}

	func removeSymbol(symbols []*symbols.Symbol, idx int) []*symbols.Symbol {
		return append(symbols[:idx], symbols[idx+1:]...)
	}

	func newTimeLineEnd(sym *string) *symbols.Symbol {
		ttype := timeline.Type_NoType
		if *sym == "bis_'" {
			ttype = timeline.Type_Bis
		}
		return &symbols.Symbol{
			Timeline: &timeline.TimeLine{
				BoundaryType: boundary.Boundary_End,
				Type:         ttype,
			},
		}
	}
*/
func NewConverter(
	mapper interfaces.SymbolMapper,
) *Converter {
	return &Converter{
		mapper: mapper,
	}
}

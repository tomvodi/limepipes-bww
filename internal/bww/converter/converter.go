// Package converter contains code to convert the parsed data into a music model.
package converter

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/barline"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/boundary"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/musicmodel"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/pitch"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/timeline"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/tune"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"github.com/tomvodi/limepipes-plugin-bww/internal/interfaces"
	"reflect"
	"strconv"
)

type staffContext struct {
	PendingOldTie         pitch.Pitch
	PendingNewTie         bool
	NextMeasureComments   []string
	NextMeasureInLineText []string
	PreviousStaveMeasures []*measure.Measure
}

type Converter struct {
	mapper interfaces.SymbolMapper
}

func (c *Converter) Convert(
	grammar *common.BwwStructure,
) (musicmodel.MusicModel, error) {
	return c.convertGrammarToModel(grammar)
}

func (c *Converter) convertGrammarToModel(
	grammar *common.BwwStructure,
) (musicmodel.MusicModel, error) {
	var tunes musicmodel.MusicModel

	var newTune *tune.Tune
	staffCtx := &staffContext{
		PendingOldTie: pitch.Pitch_NoPitch,
	}
	for i, t := range grammar.Tunes {
		if t.Header.HasTitle() {
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

func cleanupAndAppendMeasure(
	measures []*measure.Measure,
	measure *measure.Measure,
) []*measure.Measure {
	cleanupMeasure(measure)
	if len(measure.Symbols) == 0 {
		measure.Symbols = nil
	}

	return append(measures, measure)
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

func New(
	mapper interfaces.SymbolMapper,
) *Converter {
	return &Converter{
		mapper: mapper,
	}
}

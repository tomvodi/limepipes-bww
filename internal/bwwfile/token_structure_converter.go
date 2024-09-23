package bwwfile

import (
	"fmt"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"github.com/tomvodi/limepipes-plugin-bww/internal/structure"
	"strings"
)

var structuredTextTemplate = "\"%s\",(%s,L,0,0,Times New Roman,11,700,0,0,0,0,0,0)"

type TuneTokens []*common.Token

type TokenConverter struct {
}

func (tc *TokenConverter) Convert(
	tokens []*common.Token,
) (*structure.BwwFile, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens to convert")
	}

	bf := &structure.BwwFile{}
	tt := getTuneTokens(tokens)

	bv, err := getBagpipePlayerVersion(tokens)
	if err != nil {
		return nil, err
	}
	bf.BagpipePlayerVersion = bv

	for _, t := range tt {
		td := structure.TuneDefinition{}
		td.Tune = getTuneFromTokens(t)
		td.Data = getTuneDataFromTokens(bv, t)
		bf.TuneDefs = append(bf.TuneDefs, td)
	}

	return bf, nil
}

// getBagpipePlayerVersion returns the first Bagpipe Player Version from the tokens
func getBagpipePlayerVersion(
	t []*common.Token,
) (structure.BagpipePlayerVersion, error) {
	for _, token := range t {
		if bp, ok := token.Value.(structure.BagpipePlayerVersion); ok {
			return bp, nil
		}
	}

	return "", fmt.Errorf("no Bagpipe Player Version found")
}

func getTuneTokens(
	tokens []*common.Token,
) []TuneTokens {
	var tuneTokens []TuneTokens
	var currTuneTokens TuneTokens
	for _, t := range tokens {
		switch t.Value.(type) {
		case structure.BagpipePlayerVersion:
			// skipped as it is a file related definition
		case structure.TuneTitle:
			if tuneTokensHaveTitle(currTuneTokens) {
				tuneTokens = append(tuneTokens, currTuneTokens)
				currTuneTokens = make(TuneTokens, 0)
			}
			currTuneTokens = append(currTuneTokens, t)
		default:
			currTuneTokens = append(currTuneTokens, t)
		}
	}

	if len(currTuneTokens) > 0 {
		tuneTokens = append(tuneTokens, currTuneTokens)
	}

	return tuneTokens
}

func tuneTokensHaveTitle(tokens TuneTokens) bool {
	for _, t := range tokens {
		if _, ok := t.Value.(structure.TuneTitle); ok {
			return true
		}
	}

	return false
}

func getTuneFromTokens(
	tt TuneTokens,
) structure.Tune {
	t := structure.Tune{
		Header:   &structure.TuneHeader{},
		Measures: make([]*structure.Measure, 0),
	}

	fillTuneHeader(t.Header, tt)

	t.Measures = measuresForTokens(tt)

	return t
}

func getTuneDataFromTokens(
	bv structure.BagpipePlayerVersion,
	tt TuneTokens,
) []byte {
	data := []byte(bv + "\n")

	for _, token := range tt {
		// if token is line token, add a newline character
		switch v := token.Value.(type) {
		case structure.TuneTitle,
			structure.TuneType,
			structure.TuneComposer,
			structure.TuneFooter,
			structure.TuneInline,
			structure.TuneComment:
			t := fmt.Sprintf("%#v", v)
			data = append(data, []byte(t+"\n")...)
		case structure.StaffStart:
			t := fmt.Sprintf("%#v", v)
			t = strings.Trim(t, "\"")
			data = append(data, []byte(t)...)
		case structure.StaffEnd:
			t := fmt.Sprintf("%#v", v)
			t = strings.Trim(t, "\"")
			data = append(data, []byte(" "+t+"\n")...)
		case structure.Barline:
			t := fmt.Sprintf("%#v", v)
			t = strings.Trim(t, "\"")
			data = append(data, []byte("\n"+t)...)
		case structure.InlineText:
			t := fmt.Sprintf(structuredTextTemplate, string(v), "I")
			data = append(data, []byte(" "+t)...)
		case structure.StaffInline:
			t := fmt.Sprintf(structuredTextTemplate, string(v), "I")
			data = append(data, []byte(t+"\n")...)
		case structure.StaffComment:
			t := fmt.Sprintf("%#v", v)
			data = append(data, []byte(t+"\n")...)
		case structure.InlineComment:
			t := fmt.Sprintf("%#v", v)
			data = append(data, []byte(" "+t)...)
		default:
			t := fmt.Sprintf("%#v", v)
			t = strings.Trim(t, "\"")
			data = append(data, []byte(" "+t)...)
		}
	}

	return data
}

func fillTuneHeader(
	h *structure.TuneHeader,
	tt TuneTokens,
) {
	for _, token := range tt {
		switch token.Value.(type) {
		case structure.TuneTitle:
			h.Title = token.Value.(structure.TuneTitle)
		case structure.TuneType:
			h.Type = token.Value.(structure.TuneType)
		case structure.TuneComposer:
			h.Composer = token.Value.(structure.TuneComposer)
		case structure.TuneFooter:
			h.Footer = token.Value.(structure.TuneFooter)
		case structure.TuneInline:
			h.InlineTexts = append(h.InlineTexts, token.Value.(structure.TuneInline))
		case structure.TuneComment:
			h.Comments = append(h.Comments, token.Value.(structure.TuneComment))
		}
	}
}

func measuresForTokens(
	tt TuneTokens,
) []*structure.Measure {
	var m []*structure.Measure
	currMeasure := &structure.Measure{}

	for _, t := range tt {
		switch v := t.Value.(type) {
		case structure.StaffStart:
			// skipped
		case structure.StaffEnd,
			structure.Barline:
			m = append(m, currMeasure)
			currMeasure = &structure.Measure{}
		case structure.StaffInline:
			currMeasure.StaffInlineTexts = append(currMeasure.StaffInlineTexts, v)
		case structure.StaffComment:
			currMeasure.StaffComments = append(currMeasure.StaffComments, v)
		case structure.InlineComment:
			if len(currMeasure.Symbols) > 0 {
				sym := currMeasure.Symbols[len(currMeasure.Symbols)-1]
				sym.Comments = append(sym.Comments, v)
			} else {
				currMeasure.InlineComments = append(currMeasure.InlineComments, v)
			}
		case structure.InlineText:
			if len(currMeasure.Symbols) > 0 {
				sym := currMeasure.Symbols[len(currMeasure.Symbols)-1]
				sym.InlineTexts = append(sym.InlineTexts, v)
			} else {
				currMeasure.InlineTexts = append(currMeasure.InlineTexts, v)
			}
		case string:
			sym := &structure.MusicSymbol{
				Pos: structure.Position{
					Line:   t.Line,
					Column: t.Col,
				},
				Text: v,
			}
			currMeasure.Symbols = append(currMeasure.Symbols, sym)
		}
	}

	// if last symbol wasn't a barline or staff end, add the current measure
	if len(currMeasure.Symbols) > 0 {
		m = append(m, currMeasure)
	}

	return m
}

func NewTokenConverter() *TokenConverter {
	return &TokenConverter{}
}

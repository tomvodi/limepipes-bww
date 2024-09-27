package bwwfile

import (
	"fmt"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"github.com/tomvodi/limepipes-plugin-bww/internal/filestructure"
	"strings"
)

var structuredTextTemplate = "\"%s\",(%s,L,0,0,Times NewConverter Roman,11,700,0,0,0,0,0,0)"

const staffEnd = "!t"
const simpleBarline = "!"

type TuneTokens []*common.Token

type TokenConverter struct {
}

func (tc *TokenConverter) Convert(
	tokens []*common.Token,
) (*filestructure.BwwFile, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens to convert")
	}

	bf := &filestructure.BwwFile{}

	bv, err := getBagpipePlayerVersion(tokens)
	if err != nil {
		return nil, err
	}
	bf.BagpipePlayerVersion = bv

	tt := getTuneTokens(tokens)
	for _, t := range tt {
		td := filestructure.TuneDefinition{}
		td.Tune = getTuneFromTokens(t)
		td.Data = getTuneDataFromTokens(bv, t)
		bf.TuneDefs = append(bf.TuneDefs, td)
	}

	return bf, nil
}

// getBagpipePlayerVersion returns the first Bagpipe Player Version from the tokens
func getBagpipePlayerVersion(
	t []*common.Token,
) (filestructure.BagpipePlayerVersion, error) {
	for _, token := range t {
		if bp, ok := token.Value.(filestructure.BagpipePlayerVersion); ok {
			return bp, nil
		}
	}

	return "", fmt.Errorf("no Bagpipe Player Version found")
}

// getTuneTokens gets the tokens from a file and splits them up into tokens for each tune
func getTuneTokens(
	tokens []*common.Token,
) []TuneTokens {
	var tuneTokens []TuneTokens
	var currTuneTokens TuneTokens
	for _, t := range tokens {
		switch t.Value.(type) {
		case filestructure.BagpipePlayerVersion:
			// skipped as it is a file related definition
		case filestructure.TuneTitle:
			if tuneTokensHaveTitle(currTuneTokens) {
				tuneTokens = append(tuneTokens, currTuneTokens)
				currTuneTokens = make(TuneTokens, 0)
			} else {
				if tuneTokensHaveStaff(currTuneTokens) {
					currTuneTokens = prependNoNameTitle(currTuneTokens)
					tuneTokens = append(tuneTokens, currTuneTokens)
					currTuneTokens = make(TuneTokens, 0)
				}
			}
			currTuneTokens = append(currTuneTokens, t)
		default:
			currTuneTokens = append(currTuneTokens, t)
		}
	}

	if len(currTuneTokens) > 0 {
		if !tuneTokensHaveTitle(currTuneTokens) {
			currTuneTokens = prependNoNameTitle(currTuneTokens)
		}
		tuneTokens = append(tuneTokens, currTuneTokens)
	}

	return tuneTokens
}

func tuneTokensHaveTitle(tokens TuneTokens) bool {
	for _, t := range tokens {
		if _, ok := t.Value.(filestructure.TuneTitle); ok {
			return true
		}
	}

	return false
}

func tuneTokensHaveStaff(tokens TuneTokens) bool {
	for _, t := range tokens {
		if _, ok := t.Value.(filestructure.StaffStart); ok {
			return true
		}
	}

	return false
}

func prependNoNameTitle(
	tt TuneTokens,
) TuneTokens {
	return append(
		[]*common.Token{
			{
				Value: filestructure.TuneTitle("No Name"),
				Line:  0,
				Col:   0,
			},
		},
		tt...,
	)
}

func getTuneFromTokens(
	tt TuneTokens,
) *filestructure.Tune {
	t := &filestructure.Tune{
		Header:   &filestructure.TuneHeader{},
		Measures: make([]*filestructure.Measure, 0),
	}

	fillTuneHeader(t.Header, tt)

	t.Measures = measuresForTokens(tt)

	return t
}

func getTuneDataFromTokens(
	bv filestructure.BagpipePlayerVersion,
	tt TuneTokens,
) []byte {
	data := []byte(bv + "\n")

	for _, token := range tt {
		// if token is line token, add a newline character
		switch v := token.Value.(type) {
		case filestructure.TuneTitle:
			t := fmt.Sprintf(structuredTextTemplate, string(v), "T")
			data = append(data, []byte(t+"\n")...)
		case filestructure.TuneType:
			t := fmt.Sprintf(structuredTextTemplate, string(v), "Y")
			data = append(data, []byte(t+"\n")...)
		case filestructure.TuneComposer:
			t := fmt.Sprintf(structuredTextTemplate, string(v), "M")
			data = append(data, []byte(t+"\n")...)
		case filestructure.TuneFooter:
			t := fmt.Sprintf(structuredTextTemplate, string(v), "F")
			data = append(data, []byte(t+"\n")...)
		case filestructure.TuneInline:
			t := fmt.Sprintf(structuredTextTemplate, string(v), "I")
			data = append(data, []byte(t+"\n")...)
		case filestructure.TuneComment:
			t := fmt.Sprintf("%#v", v)
			data = append(data, []byte(t+"\n")...)
		case filestructure.StaffStart:
			t := fmt.Sprintf("%#v", v)
			t = strings.Trim(t, "\"")
			data = append(data, []byte(t)...)
		case filestructure.StaffEnd:
			t := fmt.Sprintf("%#v", v)
			t = strings.Trim(t, "\"")
			data = append(data, []byte(" "+t+"\n")...)
		case filestructure.Barline:
			t := fmt.Sprintf("%#v", v)
			t = strings.Trim(t, "\"")
			data = append(data, []byte("\n"+t)...)
		case filestructure.InlineText:
			t := fmt.Sprintf(structuredTextTemplate, string(v), "I")
			data = append(data, []byte(" "+t)...)
		case filestructure.StaffInline:
			t := fmt.Sprintf(structuredTextTemplate, string(v), "I")
			data = append(data, []byte(t+"\n")...)
		case filestructure.StaffComment:
			t := fmt.Sprintf("%#v", v)
			data = append(data, []byte(t+"\n")...)
		case filestructure.InlineComment:
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
	h *filestructure.TuneHeader,
	tt TuneTokens,
) {
	for _, token := range tt {
		switch token.Value.(type) {
		case filestructure.TuneTitle:
			h.Title = token.Value.(filestructure.TuneTitle)
		case filestructure.TuneType:
			h.Type = token.Value.(filestructure.TuneType)
		case filestructure.TuneComposer:
			h.Composer = token.Value.(filestructure.TuneComposer)
		case filestructure.TuneFooter:
			h.Footer = append(h.Footer, token.Value.(filestructure.TuneFooter))
		case filestructure.TuneInline:
			h.InlineTexts = append(h.InlineTexts, token.Value.(filestructure.TuneInline))
		case filestructure.TuneComment:
			h.Comments = append(h.Comments, token.Value.(filestructure.TuneComment))
		}
	}
}

func measuresForTokens(
	tt TuneTokens,
) []*filestructure.Measure {
	var m []*filestructure.Measure
	currMeasure := &filestructure.Measure{}

	for _, t := range tt {
		switch v := t.Value.(type) {
		case filestructure.StaffStart:
			// a new staff started though the current staff wasn't finished
			if len(currMeasure.Symbols) > 0 {
				m = append(m, currMeasure)
				currMeasure = &filestructure.Measure{}
			}
		case filestructure.StaffEnd:
			if v != staffEnd {
				currMeasure.RightBarline = filestructure.Barline(v)
			}
			m = append(m, currMeasure)
			currMeasure = &filestructure.Measure{}
		case filestructure.Barline:
			m = append(m, currMeasure)
			currMeasure = &filestructure.Measure{}
			if v != simpleBarline {
				currMeasure.LeftBarline = v
			}
		case filestructure.StaffInline:
			currMeasure.StaffInlineTexts = append(currMeasure.StaffInlineTexts, v)
		case filestructure.StaffComment:
			currMeasure.StaffComments = append(currMeasure.StaffComments, v)
		case filestructure.InlineComment:
			if len(currMeasure.Symbols) > 0 {
				sym := currMeasure.Symbols[len(currMeasure.Symbols)-1]
				sym.Comments = append(sym.Comments, v)
			} else {
				currMeasure.InlineComments = append(currMeasure.InlineComments, v)
			}
		case filestructure.InlineText:
			if len(currMeasure.Symbols) > 0 {
				sym := currMeasure.Symbols[len(currMeasure.Symbols)-1]
				sym.InlineTexts = append(sym.InlineTexts, v)
			} else {
				currMeasure.InlineTexts = append(currMeasure.InlineTexts, v)
			}
		case string:
			sym := &filestructure.MusicSymbol{
				Pos: filestructure.Position{
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

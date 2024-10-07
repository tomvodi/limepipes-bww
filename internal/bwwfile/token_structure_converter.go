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
type MeasureTokens []*common.Token

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

// getTuneTokens gets the tokens from a file and splits them up into tokens for each tune.
// If the first tune doesn't have a title, a TuneTitle token "No Name" is added.
func getTuneTokens(
	tokens []*common.Token,
) []TuneTokens {
	var allTuneTokens []TuneTokens
	var currTuneTokens TuneTokens

	for {
		currTuneTokens, tokens = extractFirstTuneTokens(tokens)

		if !tuneTokensHaveTitle(currTuneTokens) {
			currTuneTokens = prependNoNameTitle(currTuneTokens)
		}

		allTuneTokens = append(allTuneTokens, currTuneTokens)

		if tokens == nil {
			break
		}
	}

	return allTuneTokens
}

// extractFirstTuneTokens extracts the TuneTokens for the first tune from the tokens
// and returns the remaining tokens.
func extractFirstTuneTokens(
	tokens []*common.Token,
) (TuneTokens, []*common.Token) {
	tt := make(TuneTokens, 0)
	titleAdded := false

	for i, t := range tokens {
		switch t.Value.(type) {
		case filestructure.BagpipePlayerVersion:
			// skipped as it is a file related definition
		case filestructure.TuneTitle:
			// When tokens have a staff, there was a tune without a title before this title
			// was found.
			if tuneTokensHaveStaff(tt) || titleAdded {
				return tt, tokens[i:]
			}

			tt = append(tt, t)
			titleAdded = true
		default:
			tt = append(tt, t)
		}
	}

	return tt, nil
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

	tt = fillTuneHeader(t.Header, tt)

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
		case filestructure.TuneTempo:
			t := fmt.Sprintf("TuneTempo,%d\n", v)
			data = append(data, []byte(t)...)
		case filestructure.TempoChange:
			t := fmt.Sprintf(" TuneTempo,%d", v)
			data = append(data, []byte(t)...)
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

// fillTuneHeader fills the tune header with the tokens and returns the remaining tokens
func fillTuneHeader(
	h *filestructure.TuneHeader,
	tt TuneTokens,
) TuneTokens {
	staffStartIdx := -1
	for i, token := range tt {
		switch token.Value.(type) {
		case filestructure.StaffStart,
			filestructure.StaffComment,
			filestructure.StaffInline:
			// if staff starts, the header is finished
			staffStartIdx = i
			goto headerFinished
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
		case filestructure.TuneTempo:
			h.Tempo = token.Value.(filestructure.TuneTempo)
		}
	}

headerFinished:
	if staffStartIdx == -1 {
		return nil // only a header was found
	}

	return tt[staffStartIdx:]
}

func measuresForTokens(
	tt TuneTokens,
) []*filestructure.Measure {
	var m []*filestructure.Measure

	mtks := measureTokensForTuneTokens(tt)
	for _, mt := range mtks {
		m = append(m, measureForTokens(mt))
	}

	return m
}

// measureTokensForTuneTokens converts the tokens for a whole tune into
// tokens for each measure
func measureTokensForTuneTokens(
	tt TuneTokens,
) []MeasureTokens {
	var mt []MeasureTokens
	var currMeasureTokens MeasureTokens

	for _, t := range tt {
		switch t.Value.(type) {
		case filestructure.StaffStart:
			if measureTokensAreComplete(currMeasureTokens) {
				mt = append(mt, currMeasureTokens)
				currMeasureTokens = make(MeasureTokens, 0)
			}
		case filestructure.StaffEnd:
			// add staff end to current measure and start new measure
			currMeasureTokens = append(currMeasureTokens, t)
			mt = append(mt, currMeasureTokens)
			currMeasureTokens = make(MeasureTokens, 0)
		case filestructure.Barline:
			// start new measure and add barline to new measure
			mt = append(mt, currMeasureTokens)
			currMeasureTokens = make(MeasureTokens, 0)
			currMeasureTokens = append(currMeasureTokens, t)
		default:
			currMeasureTokens = append(currMeasureTokens, t)
		}
	}

	// if last symbol wasn't a barline or staff end, add the current measure
	if len(currMeasureTokens) > 0 {
		mt = append(mt, currMeasureTokens)
	}

	return mt
}

// returns true, if the measure tokens are complete, i.e. contains symbols.
// StaffInline and StaffComment are not considered symbols.
func measureTokensAreComplete(
	mt MeasureTokens,
) bool {
	for _, t := range mt {
		switch t.Value.(type) {
		case filestructure.StaffInline,
			filestructure.StaffComment:
			continue
		default:
			return true
		}
	}

	return false
}

func measureForTokens(
	mt MeasureTokens,
) *filestructure.Measure {
	m := &filestructure.Measure{}

	for _, t := range mt {
		processTokenForMeasure(m, t)
	}

	return m
}

// revive:disable:cognitive-complexity this method has a high complexity due to the number of different token types
// but it is easy to understand and maintain
func processTokenForMeasure(
	m *filestructure.Measure,
	t *common.Token,
) {
	newSym := &filestructure.MusicSymbol{
		Pos: filestructure.Position{
			Line:   t.Line,
			Column: t.Col,
		},
	}

	switch v := t.Value.(type) {
	case filestructure.StaffEnd:
		if v == staffEnd {
			break
		}

		m.RightBarline = filestructure.Barline(v)
	case filestructure.Barline:
		if v == simpleBarline {
			break
		}

		m.LeftBarline = v
	case filestructure.StaffInline:
		m.StaffInlineTexts = append(m.StaffInlineTexts, v)
	case filestructure.StaffComment:
		m.StaffComments = append(m.StaffComments, v)
	case filestructure.InlineComment:
		if len(m.Symbols) == 0 {
			m.InlineComments = append(m.InlineComments, v)
			break
		}

		sym := m.Symbols[len(m.Symbols)-1]
		sym.Comments = append(sym.Comments, v)
	case filestructure.InlineText:
		if len(m.Symbols) == 0 {
			m.InlineTexts = append(m.InlineTexts, v)
			break
		}

		sym := m.Symbols[len(m.Symbols)-1]
		sym.InlineTexts = append(sym.InlineTexts, v)
	case filestructure.TempoChange:
		newSym.TempoChange = v
		m.Symbols = append(m.Symbols, newSym)
	case string:
		newSym.Text = v
		m.Symbols = append(m.Symbols, newSym)
	}
}

// revive:enable:cognitive-complexity

func NewTokenConverter() *TokenConverter {
	return &TokenConverter{}
}

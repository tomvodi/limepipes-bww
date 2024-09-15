package common

import (
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
)

const TitleParameter = "T"
const TypeParameter = "Y"
const ComposerParameter = "M"
const FooterParameter = "F"
const InlineParameter = "I"

type BwwStructure struct {
	Tunes []*Tune `@@*`
}

type Tune struct {
	Pos    lexer.Position
	EndPos lexer.Position

	BagpipePlayerVersion string      `(BagpipeReader VERSION_SEP @VersionNumber)*`
	Header               *TuneHeader `@@+`
	Body                 *TuneBody   `@@`
}

type TuneHeader struct {
	TuneParameter []*TuneParameter `@@+`
}

func (t *TuneHeader) HasTitle() bool {
	for _, param := range t.TuneParameter {
		desc := param.Description
		if desc != nil && desc.FirstParameter() == TitleParameter {
			return true
		}
	}

	return false
}

func (t *TuneHeader) GetComments() []string {
	var comments []string
	for _, param := range t.TuneParameter {
		if param.Comment != "" {
			comments = append(comments, param.Comment)
		}
	}
	return comments
}

func (t *TuneHeader) GetInlineTexts() []string {
	var inlineTexts []string
	for _, param := range t.TuneParameter {
		desc := param.Description
		if desc != nil && desc.FirstParameter() == InlineParameter {
			inlineTexts = append(inlineTexts, desc.Text)
		}
	}
	return inlineTexts
}

type TuneParameter struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Config      *TuneConfig      `@@`
	Tempo       *TuneTempo       `| @@`
	Description *TuneDescription `| @@`
	Comment     string           `| @STRING`
}

// TuneConfig like page layout or MIDI note mappings
// these lines start with a defined word e.g. MIDINoteMappings,(...)
type TuneConfig struct {
	Name      string   `@PARAM_DEF PARAM_SEP`
	ParamList []string `PARAM_START @PARAM? (PARAM_SEP @PARAM?)* PARAM_END`
}

type TuneTempo struct {
	Tempo string `TEMPO_DEF PARAM_SEP @TEMPO_VALUE`
}

// TuneDescription like title, composer, arranger
// they all start with a string "title",(...)
type TuneDescription struct {
	Text      string   `@STRING PARAM_SEP`
	ParamList []string `PARAM_START @PARAM? (PARAM_SEP @PARAM?)* PARAM_END`
}

func (t *TuneDescription) FirstParameter() string {
	if len(t.ParamList) > 0 {
		return t.ParamList[0]
	}

	return ""
}

type TuneBody struct {
	Staffs []*Staff `@@*`
}

type Staff struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Start        string          `@STAFF_START`
	Symbols      []*StaffSymbols `@@*`
	End          string          `@(STAFF_END | EOF)`
	TimelineEnd  *string         `@TIMELINE_END?`
	Dalsegno     *string         `@DALSEGNO?`
	DacapoAlFine *string         `@DACAPOALFINE?`
}

type StaffSymbols struct {
	Pos    lexer.Position
	EndPos lexer.Position

	PartStart      *string `@PART_START`
	Barline        *string `| @BARLINE`
	Space          *string `| @SPACE`
	NextStaffStart *string `| @NEXT_STAFF_START`
	Comment        *string `| @STRING`
	MusicSymbol    *string `| @MUSIC_SYMBOL`

	Description *TuneDescription `| @@`
	Tempo       *TuneTempo       `| @@`
}

func (s StaffSymbols) String() string {
	if s.PartStart != nil {
		return fmt.Sprintf("PartStart(%s)", *s.PartStart)
	}
	if s.Barline != nil {
		return fmt.Sprintf("Barline(%s)", *s.Barline)
	}

	return ""
}

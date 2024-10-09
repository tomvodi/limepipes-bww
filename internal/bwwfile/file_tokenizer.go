package bwwfile

import (
	"errors"
	"fmt"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"github.com/tomvodi/limepipes-plugin-bww/internal/filestructure"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type ParserState uint8

const (
	FileState  ParserState = iota // Top level state
	StaffState                    // State for music symbols between & and!t
)

const StaffStart = "&"

var bpDefRegex = regexp.MustCompile(`(Bagpipe Reader|Bagpipe Music Writer Gold|Bagpipe Musicworks Gold):\d+\.\d+`)
var descRegex = regexp.MustCompile(`"([^"]*)",\(([TYMFI])[^)\n]+\)`)
var tokenRegex = regexp.MustCompile(`"([^"]*)",\(I(,[^,)]+)+\)|"([^"]*)"|\S+`)
var commentRegex = regexp.MustCompile(`"([^"]*)"$`)
var staffEndRegex = regexp.MustCompile(`^(''!I|!t|!I)$`)
var barlineRegex = regexp.MustCompile(`^(!|I!''|I!)$`)
var metaRegex = regexp.MustCompile(`^(MIDINoteMappings|FrequencyMappings|InstrumentMappings|GracenoteDurations|FontSizes|TuneFormat)`)
var tuneTempoRegex = regexp.MustCompile(`^TuneTempo,(\d+)$`)

type Tokenizer struct {
	state    ParserState
	currLine int
}

func (t *Tokenizer) Tokenize(
	data []byte,
) ([]*common.Token, error) {
	t.state = FileState
	t.currLine = 0

	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	// for windows file endings
	cleanData := strings.ReplaceAll(string(data), "\r\n", "\n")

	lines := strings.Split(cleanData, "\n")

	var allTokens []*common.Token
	for _, line := range lines {
		lineTokens, err := t.TokenizeLine(line)
		if errors.Is(err, common.ErrLineSkip) {
			t.currLine++
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("error tokenizing line %d: %w", t.currLine, err)
		}

		allTokens = append(allTokens, lineTokens...)

		t.currLine++
	}

	allTokens = t.checkAndModifyTokensForStaffComments(allTokens)

	return allTokens, nil
}

func (t *Tokenizer) TokenizeLine(
	line string,
) ([]*common.Token, error) {
	trimLine := strings.TrimSpace(line)
	if trimLine == "" {
		return nil, common.ErrLineSkip
	}

	if strings.HasPrefix(trimLine, StaffStart) {
		t.state = StaffState
	}

	lineTokens, err := t.getTokensFromLine(line)
	if err != nil {
		return nil, err
	}

	if lineTokensEndStaff(lineTokens) {
		t.state = FileState
	}

	return lineTokens, nil
}

// lineTokensEndStaff checks if the last token in the slice is a StaffEnd token or
// if the last token is a dalsegno or dacapoalfine and the penultimate token is a StaffEnd token.
func lineTokensEndStaff(tokens TuneTokens) bool {
	if len(tokens) == 0 {
		return false
	}

	lastToken := tokens[len(tokens)-1].Value
	_, lastTokenIsStaffEnd := lastToken.(filestructure.StaffEnd)

	if lastTokenIsStaffEnd {
		return true
	}

	// check for staff end followed by dalsegno or dacapoalfine
	if len(tokens) <= 1 {
		return false
	}

	penultimateToken := tokens[len(tokens)-2].Value
	_, penultimateTokenIsStaffEnd := penultimateToken.(filestructure.StaffEnd)

	_, lastTokenIsDalSegno := lastToken.(filestructure.DalSegno)
	_, lastTokenIsDaCapoAlFine := lastToken.(filestructure.DacapoAlFine)

	if penultimateTokenIsStaffEnd && (lastTokenIsDalSegno || lastTokenIsDaCapoAlFine) {
		return true
	}

	return false
}

func containsStaffEnd(tokens []*common.Token) bool {
	for _, tok := range tokens {
		if _, ok := tok.Value.(filestructure.StaffEnd); ok {
			return true
		}
	}

	return false
}

// checkAndModifyTokensForStaffComments
// A comment and/or inline text right before a staff start is considered a staff comment
// or staff inline text.
// This method checks the last tokens before a StaffStart for TuneComment and TuneInline
// types and changes them to StaffComment and StaffInline respectively.
func (t *Tokenizer) checkAndModifyTokensForStaffComments(
	tokens []*common.Token,
) []*common.Token {
	for i, tok := range slices.Backward(tokens) {
		if _, ok := tok.Value.(filestructure.StaffStart); ok {
			tokens = t.modifyTokensForStaffCommentsAtStaffStart(tokens, i)
		}
	}

	return tokens
}

// modifyTokensForStaffCommentsAtStaffStart modifies the tokens right before a StaffStart token
// to StaffComment or StaffInline if they are TuneComment or TuneInline.
func (t *Tokenizer) modifyTokensForStaffCommentsAtStaffStart(
	tokens []*common.Token,
	staffStartIndex int,
) []*common.Token {
	x := staffStartIndex - 1
	y := staffStartIndex - 2

	if x >= 0 {
		switch v := tokens[x].Value.(type) {
		case filestructure.TuneComment:
			tokens[x].Value = filestructure.StaffComment(v)
		case filestructure.TuneInline:
			tokens[x].Value = filestructure.StaffInline(v)
		}
	}
	if y >= 0 {
		switch v := tokens[y].Value.(type) {
		case filestructure.TuneComment:
			tokens[y].Value = filestructure.StaffComment(v)
		case filestructure.TuneInline:
			tokens[y].Value = filestructure.StaffInline(v)
		}
	}

	return tokens
}

func (t *Tokenizer) getTokensFromLine(
	line string,
) ([]*common.Token, error) {
	switch t.state {
	case FileState:
		return t.getFileTokensFromLine(line)
	case StaffState:
		tokens, err := t.getStaffTokensFromLine(line)
		if err != nil {
			return nil, err
		}

		if containsStaffEnd(tokens) && !lineTokensEndStaff(tokens) {
			return nil, fmt.Errorf("staff end token is not at the end of line %d", t.currLine)
		}

		return tokens, nil
	default:
		panic("tokenizer: unhandled parser state")
	}
}

func (t *Tokenizer) getFileTokensFromLine(
	line string,
) ([]*common.Token, error) {
	bpDef := t.isBagpipeDefinition(line)
	if bpDef != nil {
		return []*common.Token{bpDef}, nil
	}

	desc := t.isTuneDescription(line)
	if desc != nil {
		return []*common.Token{desc}, nil
	}

	comment := t.isComment(line)
	if comment != nil {
		return []*common.Token{comment}, nil
	}

	tute, err := getTuneTempo(line)
	if err != nil && !errors.Is(err, common.ErrSymbolNotFound) {
		return nil, err
	}
	if err == nil {
		return []*common.Token{
			{
				Value: filestructure.TuneTempo(tute),
				Line:  t.currLine,
				Col:   0,
			},
		}, nil
	}

	if t.isMetaData(line) {
		return nil, common.ErrLineSkip
	}

	return nil, fmt.Errorf("no file token found for line: '%s'", line)
}

func (t *Tokenizer) isBagpipeDefinition(text string) *common.Token {
	idx := bpDefRegex.FindIndex([]byte(text))
	if idx == nil {
		return nil
	}
	bpDef := text[idx[0]:idx[1]]
	return &common.Token{
		Value: filestructure.BagpipePlayerVersion(bpDef),
		Line:  t.currLine,
		Col:   idx[0],
	}
}

func (t *Tokenizer) isTuneDescription(text string) *common.Token {
	idx := descRegex.FindAllSubmatchIndex([]byte(text), -1)
	if idx == nil {
		return nil
	}

	loc := idx[0]
	var val any
	desc := text[loc[2]:loc[3]]
	descType := text[loc[4]:loc[5]]
	if descType == "T" {
		val = filestructure.TuneTitle(desc)
	}
	if descType == "Y" {
		val = filestructure.TuneType(desc)
	}
	if descType == "M" {
		val = filestructure.TuneComposer(desc)
	}
	if descType == "F" {
		val = filestructure.TuneFooter(desc)
	}
	if descType == "I" {
		val = filestructure.TuneInline(desc)
	}

	return &common.Token{
		Value: val,
		Line:  t.currLine,
		Col:   loc[0],
	}
}

func (t *Tokenizer) isInlineText(text string) *common.Token {
	idx := descRegex.FindAllSubmatchIndex([]byte(text), -1)
	if idx == nil {
		return nil
	}

	loc := idx[0]
	var val any
	desc := text[loc[2]:loc[3]]
	descType := text[loc[4]:loc[5]]
	if descType == "I" {
		val = filestructure.InlineText(desc)
	}

	return &common.Token{
		Value: val,
		Line:  t.currLine,
		Col:   loc[0],
	}
}

func (t *Tokenizer) isComment(
	text string,
) *common.Token {
	idx := commentRegex.FindAllSubmatchIndex([]byte(text), -1)
	if idx == nil {
		return nil
	}

	loc := idx[0]
	comment := text[loc[2]:loc[3]]
	tok := &common.Token{
		Value: filestructure.InlineComment(comment),
		Line:  t.currLine,
		Col:   loc[0],
	}

	if t.state == FileState {
		tok.Value = filestructure.TuneComment(comment)
	}

	return tok
}

func (t *Tokenizer) isMetaData(
	text string,
) bool {
	trimmed := strings.TrimSpace(text)
	return metaRegex.MatchString(trimmed)
}

func (t *Tokenizer) getStaffTokensFromLine(
	line string,
) (tokens []*common.Token, err error) {
	idxs := tokenRegex.FindAllIndex([]byte(line), -1)
	for _, idx := range idxs {
		tokStr := line[idx[0]:idx[1]]
		currTok := &common.Token{
			Line: t.currLine,
			Col:  idx[0],
		}

		val, err := t.getTokenValueForString(tokStr)
		if err != nil {
			return nil, err
		}
		currTok.Value = val

		tokens = append(tokens, currTok)
	}

	return tokens, nil
}

// revive:disable:cognitive-complexity this method has a high complexity due to the number of different token types
// but it is easy to understand and maintain
func (t *Tokenizer) getTokenValueForString(
	tokStr string,
) (any, error) {
	if barlineRegex.MatchString(tokStr) {
		return filestructure.Barline(tokStr), nil
	}

	if tokStr == StaffStart {
		return filestructure.StaffStart(tokStr), nil
	}

	if tokStr == "dalsegno" {
		return filestructure.DalSegno(tokStr), nil
	}

	if tokStr == "dacapoalfine" {
		return filestructure.DacapoAlFine(tokStr), nil
	}

	if staffEndRegex.MatchString(tokStr) {
		return filestructure.StaffEnd(tokStr), nil
	}

	tute, err := getTuneTempo(tokStr)
	if err != nil && !errors.Is(err, common.ErrSymbolNotFound) {
		return nil, err
	}
	if err == nil {
		return filestructure.TempoChange(tute), nil
	}

	if ct := t.isComment(tokStr); ct != nil {
		return ct.Value, nil
	}

	if it := t.isInlineText(tokStr); it != nil {
		return it.Value, nil
	}

	return tokStr, nil
}

// revive:enable:cognitive-complexity

func getTuneTempo(text string) (uint32, error) {
	idx := tuneTempoRegex.FindAllSubmatchIndex([]byte(text), -1)
	if len(idx) == 0 {
		return 0, common.ErrSymbolNotFound
	}

	loc := idx[0]
	tt := text[loc[2]:loc[3]]
	tempo, err := strconv.ParseUint(tt, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(tempo), nil
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{}
}

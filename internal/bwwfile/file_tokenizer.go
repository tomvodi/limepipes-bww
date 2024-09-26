package bwwfile

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"github.com/tomvodi/limepipes-plugin-bww/internal/filestructure"
	"regexp"
	"slices"
	"strings"
)

type ParserState uint8

const (
	FileState  ParserState = iota // Top level state
	StaffState                    // State for music symbols between & and!t
)

var bpDefRegex = regexp.MustCompile(`(Bagpipe Reader|Bagpipe Music Writer Gold|Bagpipe Musicworks Gold):\d+\.\d+`)
var descRegex = regexp.MustCompile(`"([^"]*)",\(([TYMFI])(,[^,)]+)+\)`)
var tokenRegex = regexp.MustCompile(`"([^"]*)",\(I(,[^,)]+)+\)|"([^"]*)"|\S+`)
var commentRegex = regexp.MustCompile(`"([^"]*)"$`)
var staffEndRegex = regexp.MustCompile(`^(''!I|!t|!I)$`)
var barlineRegex = regexp.MustCompile(`^(!|I!''|I!)$`)

type Tokenizer struct {
	state    ParserState
	currLine int
	currCol  int
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
		trimLine := strings.TrimSpace(line)
		if trimLine == "" {
			t.currLine++
			continue
		}

		if strings.HasPrefix(trimLine, "&") {
			allTokens = t.checkAndModifyLastTokensForStaffComments(allTokens)
			t.state = StaffState
		}

		tokens, err := t.getTokensFromLine(line)
		if err != nil {
			return nil, err
		}

		lastToken := tokens[len(tokens)-1].Value
		_, ok := lastToken.(filestructure.StaffEnd)
		if ok {
			t.state = FileState
		}

		allTokens = append(allTokens, tokens...)

		t.currLine++
	}

	for _, tok := range allTokens {
		log.Info().Msgf("Token: '%s' at Line %d in column %d of type %T",
			tok.Value, tok.Line, tok.Col, tok.Value,
		)
	}

	return allTokens, nil
}

// checkAndModifyLastTokensForStaffComments checks the last tokens for TuneComment and TuneInline
// types and changes them to StaffComment and StaffInline respectively.
// As a comment or inline text right before a staff start is considered a staff comment
// or staff inline text.
func (t *Tokenizer) checkAndModifyLastTokensForStaffComments(
	tokens []*common.Token,
) []*common.Token {
	comment, inline := false, false
	for i, tok := range slices.Backward(tokens) {
		// maximum of 2 lines can be changed when there is a
		// staff comment and an staff inline text
		if i < len(tokens)-2 {
			break
		}

		switch v := tok.Value.(type) {
		case filestructure.TuneComment:
			if !comment {
				tokens[i].Value = filestructure.StaffComment(v)
				comment = true
			}
		case filestructure.TuneInline:
			if !inline {
				tokens[i].Value = filestructure.StaffInline(v)
				inline = true
			}
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
		return t.getStaffTokensFromLine(line)
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
	for _, loc := range idx {
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

	return nil
}

func (t *Tokenizer) isInlineText(text string) *common.Token {
	idx := descRegex.FindAllSubmatchIndex([]byte(text), -1)
	if idx == nil {
		return nil
	}
	for _, loc := range idx {
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

	return nil
}

func (t *Tokenizer) isComment(
	text string,
) *common.Token {
	idx := commentRegex.FindAllSubmatchIndex([]byte(text), -1)
	if idx == nil {
		return nil
	}
	for _, loc := range idx {
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

	return nil
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

		if barlineRegex.MatchString(tokStr) {
			currTok.Value = filestructure.Barline(tokStr)
			tokens = append(tokens, currTok)
			continue
		}

		if tokStr == "&" {
			currTok.Value = filestructure.StaffStart(tokStr)
			tokens = append(tokens, currTok)
			continue
		}

		if staffEndRegex.MatchString(tokStr) {
			currTok.Value = filestructure.StaffEnd(tokStr)
			tokens = append(tokens, currTok)
			continue
		}

		if ct := t.isComment(tokStr); ct != nil {
			currTok.Value = ct.Value
			tokens = append(tokens, currTok)
			continue
		}

		if it := t.isInlineText(tokStr); it != nil {
			currTok.Value = it.Value
			tokens = append(tokens, currTok)
			continue
		}

		currTok.Value = tokStr
		tokens = append(tokens, currTok)
	}

	return tokens, nil
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{}
}

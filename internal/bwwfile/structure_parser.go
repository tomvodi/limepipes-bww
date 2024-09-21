package bwwfile

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes-plugin-bww/internal/bwwfile/interfaces"
	"github.com/tomvodi/limepipes-plugin-bww/internal/bwwfile/parserstate"
	"github.com/tomvodi/limepipes-plugin-bww/internal/structure"
	"regexp"
	"strings"
)

var bpDefRegex = regexp.MustCompile(`(Bagpipe Reader|Bagpipe Music Writer Gold|Bagpipe Musicworks Gold):\d+\.\d+`)
var descRegex = regexp.MustCompile(`"([^"]*)",\(([TYMFI])(,[^,)]+)+\)`)
var tokenRegex = regexp.MustCompile(`\S+`)

type token struct {
	value any
	line  int
	col   int
}

type StructureParser struct {
	state    parserstate.ParserState
	currLine int
	currCol  int
}

func (s *StructureParser) ParseDocumentStructure(
	data []byte,
) (*structure.BwwFile, error) {
	s.state = parserstate.File
	s.currLine = 1
	s.currCol = 1

	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	// for windows file endings
	cleanData := strings.ReplaceAll(string(data), "\r\n", "\n")

	lines := strings.Split(cleanData, "\n")

	bwFile := &structure.BwwFile{}
	for _, line := range lines {
		trimLine := strings.TrimSpace(line)
		if trimLine == "" {
			s.currLine++
			continue
		}

		if strings.HasPrefix(trimLine, "&") {
			s.state = parserstate.Stave
		}

		tokens, err := s.getTokensFromLine(line)
		if err != nil {
			return nil, err
		}

		for _, tok := range tokens {
			log.Info().Msgf("Token: %s is a %T\n", tok.value, tok.value)
			log.Info().Msgf("Token: %s at line %d in column %d\n", tok.value, tok.line, tok.col)
		}

		s.currLine++
	}

	return bwFile, nil
}

func (s *StructureParser) getTokensFromLine(
	line string,
) ([]*token, error) {
	switch s.state {
	case parserstate.File:
		return s.getFileTokensFromLine(line)
	case parserstate.Stave:
		return s.getStaffTokensFromLine(line)
	default:
		panic("unhandled default case")
	}
}

func (s *StructureParser) getFileTokensFromLine(
	line string,
) ([]*token, error) {
	bpDef := s.isBagpipeDefinition(line)
	if bpDef != nil {
		return []*token{bpDef}, nil
	}

	desc := s.isTuneDescription(line)
	if desc != nil {
		return []*token{desc}, nil
	}

	return nil, fmt.Errorf("no token found")
}

func (s *StructureParser) isBagpipeDefinition(line string) *token {
	idx := bpDefRegex.FindIndex([]byte(line))
	if idx == nil {
		return nil
	}
	bpDef := line[idx[0]:idx[1]]
	return &token{
		value: bpDef,
		line:  s.currLine,
		col:   idx[0],
	}
}

func (s *StructureParser) isTuneDescription(line string) *token {
	idx := descRegex.FindAllSubmatchIndex([]byte(line), -1)
	if idx == nil {
		return nil
	}
	for _, loc := range idx {
		var val any
		desc := line[loc[2]:loc[3]]
		descType := line[loc[4]:loc[5]]
		if descType == "T" {
			val = structure.TuneTitle(desc)
		}
		if descType == "Y" {
			val = structure.TuneType(desc)
		}
		if descType == "M" {
			val = structure.TuneComposer(desc)
		}
		if descType == "F" {
			val = structure.TuneFooter(desc)
		}
		if descType == "I" {
			val = structure.TuneInline(desc)
		}

		return &token{
			value: val,
			line:  s.currLine,
			col:   loc[0],
		}
	}

	return nil
}

func (s *StructureParser) getStaffTokensFromLine(
	line string,
) (tokens []*token, err error) {
	idxs := tokenRegex.FindAllIndex([]byte(line), -1)
	for _, idx := range idxs {
		tokStr := line[idx[0]:idx[1]]
		currTok := &token{
			line: s.currLine,
			col:  idx[0],
		}

		if tokStr == "!" {
			currTok.value = structure.Barline(tokStr)
			tokens = append(tokens, currTok)
			continue
		}
		currTok.value = tokStr
		tokens = append(tokens, currTok)
	}

	return tokens, nil
}

func NewStructureParser() interfaces.StructureParser {
	return &StructureParser{}
}

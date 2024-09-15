package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
)

type DocumentStructureParser struct {
}

func (d *DocumentStructureParser) ParseDocumentStructure(
	data []byte,
) (*common.BwwStructure, error) {
	parser, err := participle.Build[common.BwwStructure](
		participle.Elide("WHITESPACE"),
		participle.Lexer(Lexer),
		participle.Unquote("STRING"),
	)
	if err != nil {
		return nil, err
	}

	bwwStruct, err := parser.ParseBytes("", data)
	if err != nil {
		return nil, err
	}

	return bwwStruct, nil
}

func NewDocumentStructureParser() *DocumentStructureParser {
	return &DocumentStructureParser{}
}

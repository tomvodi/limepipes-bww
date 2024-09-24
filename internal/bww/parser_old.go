package bww

import (
	"github.com/alecthomas/participle/v2"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/musicmodel"
)

type Parser struct {
}

func (b *Parser) ParseBwwData(
	data []byte,
) (musicmodel.MusicModel, error) {
	parser, err := participle.Build[BwwDocument](
		participle.Elide("WHITESPACE"),
		participle.Lexer(BwwLexer),
		participle.Unquote("STRING"),
	)
	if err != nil {
		return nil, err
	}

	var bwwDoc *BwwDocument
	bwwDoc, err = parser.ParseBytes("", data)
	if err != nil {
		return nil, err
	}

	return convertGrammarToModel(bwwDoc)
}

func NewBwwParser() *Parser {
	return &Parser{}
}

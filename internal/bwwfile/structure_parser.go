package bwwfile

import (
	"github.com/tomvodi/limepipes-plugin-bww/internal/filestructure"
	"github.com/tomvodi/limepipes-plugin-bww/internal/interfaces"
)

type StructureParser struct {
	tokenizer interfaces.FileTokenizer
	conv      interfaces.TokenStructureConverter
}

func (t *StructureParser) ParseDocumentStructure(
	data []byte,
) (*filestructure.BwwFile, error) {
	tokens, err := t.tokenizer.Tokenize(data)
	if err != nil {
		return nil, err
	}

	return t.conv.Convert(tokens)
}

func NewStructureParser(
	tokenizer interfaces.FileTokenizer,
	conv interfaces.TokenStructureConverter,
) *StructureParser {
	return &StructureParser{
		tokenizer: tokenizer,
		conv:      conv,
	}
}

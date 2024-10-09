package parser

import (
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes-plugin-bww/internal/interfaces"
)

type Parser struct {
	structureParser interfaces.StructureParser
	gConverter      interfaces.StructureToModelConverter
}

func (p *Parser) ParseBwwData(
	data []byte,
) ([]*messages.ParsedTune, error) {
	bd, err := p.structureParser.ParseDocumentStructure(data)
	if err != nil {
		return nil, err
	}

	var pt []*messages.ParsedTune
	for _, def := range bd.TuneDefs {
		ct, err := p.gConverter.Convert(def.Tune)
		if err != nil {
			return nil, err
		}

		pt = append(pt, &messages.ParsedTune{
			Tune:         ct,
			TuneFileData: def.Data,
		})
	}

	return pt, nil
}

func New(
	structureParser interfaces.StructureParser,
	gConverter interfaces.StructureToModelConverter,
) *Parser {
	return &Parser{
		structureParser: structureParser,
		gConverter:      gConverter,
	}
}

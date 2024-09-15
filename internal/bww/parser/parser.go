package parser

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/musicmodel"
	"github.com/tomvodi/limepipes-plugin-bww/internal/interfaces"
)

type Parser struct {
	fileSplitter    interfaces.BwwFileByTuneSplitter
	structureParser interfaces.StructureParser
	gConverter      interfaces.GrammarConverter
}

func (p *Parser) ParseBwwData(data []byte) (musicmodel.MusicModel, error) {
	// 1. Split data by single tunes
	td, err := p.fileSplitter.SplitFileData(data)
	if err != nil {
		return nil, err
	}

	// 2. Parse structure of each tune and convert it into music model
	var mumo musicmodel.MusicModel
	for i := range td.TuneTitles() {
		bd, err := p.structureParser.ParseDocumentStructure(td.Data(i))
		if err != nil {
			return nil, err
		}

		m, err := p.gConverter.Convert(bd)
		if err != nil {
			return nil, err
		}

		mumo = append(mumo, m...)
	}

	return mumo, nil
}

func New(
	fileSplitter interfaces.BwwFileByTuneSplitter,
	structureParser interfaces.StructureParser,
	gConverter interfaces.GrammarConverter,
) *Parser {
	return &Parser{
		fileSplitter:    fileSplitter,
		structureParser: structureParser,
		gConverter:      gConverter,
	}
}

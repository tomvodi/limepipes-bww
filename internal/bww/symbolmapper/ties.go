package symbolmapper

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/tie"
)

func init() {
	tieStart := &symbols.Symbol{
		Note: &symbols.Note{
			Tie: tie.Tie_Start,
		},
	}

	tieEnd := &symbols.Symbol{
		Note: &symbols.Note{
			Tie: tie.Tie_End,
		},
	}

	symbolsMap["^ts"] = tieStart
	symbolsMap["^te"] = tieEnd
}

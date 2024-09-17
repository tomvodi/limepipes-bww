package symbolmapper

import (
	"fmt"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/accidental"
)

var accidentals = []string{"sharp", "natural", "flat"}
var accMap = map[string]accidental.Accidental{
	"sharp":   accidental.Accidental_Sharp,
	"natural": accidental.Accidental_Natural,
	"flat":    accidental.Accidental_Flat,
}

func init() {
	for _, a := range accidentals {
		for _, p := range lowPitchesLgToHA {
			symbolsMap[fmt.Sprintf("%s%s", a, p)] = &symbols.Symbol{
				Note: &symbols.Note{
					Accidental: accMap[a],
				},
			}
		}
	}
}

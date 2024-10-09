package symbolmapper

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/pitch"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"
)

func init() {
	dbl := newEmbellishment(embellishment.Type_DoubleGrace, pitch.Pitch_D)
	for _, s := range lowPitchesLgToC {
		symbolsMap["d"+s] = dbl
	}

	dbl = newEmbellishment(embellishment.Type_DoubleGrace, pitch.Pitch_E)
	for _, s := range lowPitchesLgToD {
		symbolsMap["e"+s] = dbl
	}

	dbl = newEmbellishment(embellishment.Type_DoubleGrace, pitch.Pitch_F)
	for _, s := range lowPitchesLgToE {
		symbolsMap["f"+s] = dbl
	}

	dbl = newEmbellishment(embellishment.Type_DoubleGrace, pitch.Pitch_HighG)
	for _, s := range lowPitchesLgToF {
		symbolsMap["g"+s] = dbl
	}

	dbl = newEmbellishment(embellishment.Type_DoubleGrace, pitch.Pitch_HighA)
	for _, s := range lowPitchesLgToHG {
		symbolsMap["t"+s] = dbl
	}
}

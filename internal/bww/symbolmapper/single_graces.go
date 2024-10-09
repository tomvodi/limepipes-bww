package symbolmapper

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/pitch"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"
)

var singleGraceMap = map[string]*embellishment.Embellishment{
	"ag": {
		Type:  embellishment.Type_SingleGrace,
		Pitch: pitch.Pitch_LowA,
	},
	"bg": {
		Type:  embellishment.Type_SingleGrace,
		Pitch: pitch.Pitch_B,
	},
	"cg": {
		Type:  embellishment.Type_SingleGrace,
		Pitch: pitch.Pitch_C,
	},
	"dg": {
		Type:  embellishment.Type_SingleGrace,
		Pitch: pitch.Pitch_D,
	},
	"eg": {
		Type:  embellishment.Type_SingleGrace,
		Pitch: pitch.Pitch_E,
	},
	"fg": {
		Type:  embellishment.Type_SingleGrace,
		Pitch: pitch.Pitch_F,
	},
	"gg": {
		Type:  embellishment.Type_SingleGrace,
		Pitch: pitch.Pitch_HighG,
	},
	"tg": {
		Type:  embellishment.Type_SingleGrace,
		Pitch: pitch.Pitch_HighA,
	},
}

func init() {
	for k, e := range singleGraceMap {
		symbolsMap[k] = &symbols.Symbol{
			Note: &symbols.Note{
				Embellishment: e,
			},
		}
	}
}
